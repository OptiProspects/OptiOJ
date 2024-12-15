package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"context"
	"errors"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// GetUserList 获取用户列表
func GetUserList(req *models.UserListRequest) ([]models.UserListItem, int64, error) {
	var users []models.UserListItem
	var total int64

	// 创建基础查询，使用 JOIN 连接 profiles 表获取创建时间和 admins 表获取角色信息
	query := config.DB.Model(&models.User{}).
		Select(`
			users.id, 
			users.username, 
			users.email, 
			users.phone, 
			profiles.create_at as created_at, 
			profiles.update_at as updated_at,
			last_login.login_time as last_login_time,
			last_login.ip_address as last_login_ip,
			COALESCE(admins.role, 'user') as role
		`).
		Joins("LEFT JOIN profiles ON users.id = profiles.user_id").
		Joins("LEFT JOIN admins ON users.id = admins.user_id").
		// 使用子查询获取每个用户最后一次成功登录的记录
		Joins(`
			LEFT JOIN (
				SELECT 
					user_id, 
					MAX(login_time) as login_time,
					ip_address
				FROM user_logins
				WHERE login_status = 'success'
				GROUP BY user_id
			) AS last_login ON users.id = last_login.user_id
		`)

	// 应用筛选条件
	if req.Username != "" {
		query = query.Where("users.username LIKE ?", "%"+req.Username+"%")
	}
	if req.Email != "" {
		query = query.Where("users.email LIKE ?", "%"+req.Email+"%")
	}
	if req.Phone != "" {
		query = query.Where("users.phone LIKE ?", "%"+req.Phone+"%")
	}

	// 处理状态筛选
	if req.Status != "" {
		if req.Status == "banned" {
			query = query.Joins("JOIN user_bans ON users.id = user_bans.user_id").
				Where("user_bans.is_active = ? AND (user_bans.expires_at > ? OR user_bans.expires_at IS NULL)",
					true, time.Now())
		} else if req.Status == "normal" {
			query = query.Where("NOT EXISTS (SELECT 1 FROM user_bans WHERE user_bans.user_id = users.id AND user_bans.is_active = ? AND (user_bans.expires_at > ? OR user_bans.expires_at IS NULL))",
				true, time.Now())
		}
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = query.Order("profiles.create_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	// 查询每个用户的封禁状态
	for i := range users {
		banned, reason, _ := IsUserBanned(uint(users[i].ID))
		if banned {
			users[i].Status = "banned"
			users[i].BanReason = reason
		} else {
			users[i].Status = "normal"
		}
	}

	return users, total, nil
}

// UpdateUserInfo 更新用户信息
func UpdateUserInfo(userID uint, req *models.UserUpdateRequest) error {
	updates := make(map[string]interface{})

	if req.Username != "" {
		// 检查用户名是否已存在
		var count int64
		if err := config.DB.Model(&models.User{}).Where("username = ? AND id != ?", req.Username, userID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("用户名已存在")
		}
		updates["username"] = req.Username
	}

	if req.Email != "" {
		// 检查邮箱是否已存在
		var count int64
		if err := config.DB.Model(&models.User{}).Where("email = ? AND id != ?", req.Email, userID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("邮箱已存在")
		}
		updates["email"] = req.Email
	}

	if req.Phone != "" {
		// 检查手机号是否已存在
		var count int64
		if err := config.DB.Model(&models.User{}).Where("phone = ? AND id != ?", req.Phone, userID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return errors.New("手机号已存在")
		}
		updates["phone"] = req.Phone
	}

	if len(updates) == 0 {
		return nil
	}

	return config.DB.Model(&models.User{}).Where("id = ?", userID).Updates(updates).Error
}

// IsUserBanned 检查用户是否被封禁
func IsUserBanned(userID uint) (bool, string, error) {
	var ban models.UserBan
	err := config.DB.Where("user_id = ? AND is_active = ? AND (expires_at > ? OR expires_at IS NULL)",
		userID, true, time.Now()).
		Order("created_at DESC").
		First(&ban).Error

	if err == gorm.ErrRecordNotFound {
		return false, "", nil
	}
	if err != nil {
		return false, "", err
	}

	return true, ban.Reason, nil
}

// BanUser 封禁用户
func BanUser(req *models.UserBanRequest, adminID uint) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {
		// 创建封禁记录
		ban := models.UserBan{
			UserID:    uint64(req.UserID),
			Reason:    req.BanReason,
			ExpiresAt: req.BanExpires,
			CreatedBy: uint64(adminID),
			IsActive:  true,
		}

		if err := tx.Create(&ban).Error; err != nil {
			return err
		}

		// 删除用户的所有token
		patterns := []string{"access_token:*", "refresh_token:*"}
		for _, pattern := range patterns {
			var cursor uint64
			for {
				keys, nextCursor, err := config.RedisClient.Scan(context.Background(), cursor, pattern, 100).Result()
				if err != nil {
					return err
				}

				for _, key := range keys {
					userIDStr, err := config.RedisClient.Get(context.Background(), key).Result()
					if err != nil {
						continue
					}
					if userIDStr == strconv.FormatUint(uint64(req.UserID), 10) {
						config.RedisClient.Del(context.Background(), key)
					}
				}

				cursor = nextCursor
				if cursor == 0 {
					break
				}
			}
		}

		return nil
	})
}

// UnbanUser 解封用户
func UnbanUser(userID uint, adminID uint) error {
	// 将所有该用户的活跃封禁记录设置为非活跃
	return config.DB.Model(&models.UserBan{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Updates(map[string]interface{}{
			"is_active": false,
		}).Error
}
