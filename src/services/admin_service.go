package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"errors"
	"time"
)

// 检查用户是否为管理员
func IsAdmin(userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Admin{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 检查用户是否为超级管理员
func IsSuperAdmin(userID uint) (bool, error) {
	var count int64
	err := config.DB.Model(&models.Admin{}).
		Where("user_id = ? AND role = ?", userID, "super_admin").
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// 添加管理员
func AddAdmin(userID uint, role string) error {
	// 验证角色类型
	if role != "admin" && role != "super_admin" {
		return errors.New("无效的角色类型")
	}

	// 检查用户是否存在
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 检查是否已经是管理员
	exists, _ := IsAdmin(userID)
	if exists {
		return errors.New("该用户已经是管理员")
	}

	// 创建管理员记录
	admin := models.Admin{
		UserID:    uint64(userID),
		Role:      role,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := config.DB.Create(&admin).Error; err != nil {
		return errors.New("添加管理员失败")
	}

	return nil
}

// 移除管理员
func RemoveAdmin(userID uint) error {
	result := config.DB.Where("user_id = ?", userID).Delete(&models.Admin{})
	if result.Error != nil {
		return errors.New("移除管理员失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("该用户不是管理员")
	}
	return nil
}

// 获取所有管理员列表（包含详细信息）
func GetAllAdmins() ([]models.AdminDetailItem, error) {
	var admins []models.AdminDetailItem

	err := config.DB.Model(&models.Admin{}).
		Select(`
			admins.id,
			admins.user_id,
			admins.role,
			admins.created_at,
			users.username,
			users.email,
			last_login.login_time as last_login_time,
			last_login.ip_address as last_login_ip
		`).
		Joins("LEFT JOIN users ON admins.user_id = users.id").
		Joins(`
			LEFT JOIN (
				SELECT 
					ul1.user_id, 
					ul1.login_time,
					ul1.ip_address
				FROM user_logins ul1
				INNER JOIN (
					SELECT user_id, MAX(login_time) as max_login_time
					FROM user_logins
					WHERE login_status = 'success'
					GROUP BY user_id
				) ul2 ON ul1.user_id = ul2.user_id AND ul1.login_time = ul2.max_login_time
			) last_login ON admins.user_id = last_login.user_id
		`).
		Find(&admins).Error

	if err != nil {
		return nil, errors.New("获取管理员列表失败")
	}

	return admins, nil
}
