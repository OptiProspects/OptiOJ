package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/location"
	"OptiOJ/src/models"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterUser(req *models.RegisterRequest) (uint, error) {
	var userID uint
	// 开启事务
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// 检查用户名是否已存在
		var existingUser models.User
		if err := tx.Where("username = ?", req.UserName).First(&existingUser).Error; err == nil {
			return errors.New("用户名已存在")
		}

		// 检查邮箱是否已存在
		if req.RequestEmail != "" {
			if err := tx.Where("email = ?", req.RequestEmail).First(&existingUser).Error; err == nil {
				return errors.New("邮箱已被注册")
			}
		}

		// 检查手机号是否已存在
		if req.RequestPhone != "" {
			if err := tx.Where("phone = ?", req.RequestPhone).First(&existingUser).Error; err == nil {
				return errors.New("手机号已被注册")
			}
		}

		// 密码加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.PassWord), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 创建用户
		user := &models.User{
			Username: req.UserName,
			Password: string(hashedPassword),
			Email:    req.RequestEmail,
		}

		// 只有当请求中包含手机号时才设置
		if req.RequestPhone != "" {
			user.Phone = &req.RequestPhone
		}

		// 存储用户信息到数据库
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		userID = uint(user.ID)

		// 创建用户资料
		now := time.Now()
		profile := models.Profile{
			UserID:   int(userID),
			Bio:      "",
			Gender:   "",
			School:   "",
			Birthday: nil,
			Location: "",
			RealName: "",
			CreateAt: now,
			UpdateAt: now,
		}

		// 如果提供了省份和城市，组合成location字段
		if req.Province != "" && req.City != "" {
			if !location.IsValidProvince(req.Province) {
				return errors.New("无效的省份")
			}
			if !location.IsValidCity(req.Province, req.City) {
				return errors.New("无效的城市")
			}
			profile.Location = req.Province + "-" + req.City
		}

		if err := tx.Create(&profile).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return userID, nil
}

func ValidateLogin(accountInfo string, password string) (*models.User, error) {
	var user models.User

	// 尝试使用用户名、邮箱或手机号查找用户
	result := config.DB.Where("username = ? OR email = ? OR phone = ?",
		accountInfo, accountInfo, accountInfo).First(&user)

	if result.Error != nil {
		return nil, errors.New("用户不存在")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误")
	}

	return &user, nil
}

// HashPassword 辅助函数，用于密码加密
func HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// GenerateUsers 批量生成用户
func GenerateUsers(req *models.GenerateUsersRequest) (*models.GenerateUsersResponse, error) {
	users := make([]models.GeneratedUserInfo, 0, req.Count)
	createdUsers := make([]models.User, 0, req.Count)

	// 生成一个随机字符串作为批次标识
	batchID := generateRandomString(6)

	// 开启事务
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < req.Count; i++ {
			// 生成用户名和邮箱，加入批次标识避免重复
			username := fmt.Sprintf("%s_%s_%d_%s", req.Prefix, batchID, i+1, req.Suffix)
			email := fmt.Sprintf("%s@%s", username, req.Domain)

			// 生成随机密码
			randomPassword := generateRandomPassword(10)

			hashedPassword, err := HashPassword(randomPassword)
			if err != nil {
				return fmt.Errorf("密码加密失败: %v", err)
			}

			// 创建用户记录，只设置必要的字段
			user := models.User{
				Username: username,
				Password: hashedPassword,
				Email:    email,
			}

			// 创建用户
			if err := tx.Select("Username", "Password", "Email").Create(&user).Error; err != nil {
				if strings.Contains(err.Error(), "duplicate") {
					continue // 如果用户名重复，跳过当前用户
				}
				return fmt.Errorf("创建用户失败: %v", err)
			}

			// 创建用户资料，只设置必要的字段
			now := time.Now()
			profile := models.Profile{
				UserID:   int(user.ID),
				CreateAt: now,
				UpdateAt: now,
			}

			// 只插入指定的字段
			if err := tx.Select("user_id", "create_at", "update_at").Create(&profile).Error; err != nil {
				return fmt.Errorf("创建用户资料失败: %v", err)
			}

			createdUsers = append(createdUsers, user)
			users = append(users, models.GeneratedUserInfo{
				Username: username,
				Password: randomPassword,
				Email:    email,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &models.GenerateUsersResponse{
		Users: users,
		Total: len(users),
	}, nil
}

// generateRandomPassword 生成随机密码
func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"
	rand.Seed(time.Now().UnixNano())
	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}
	return string(password)
}

// generateRandomString 生成指定长度的随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

// GetUserByID 根据用户ID获取用户信息
func GetUserByID(userID uint, user *models.User) error {
	return config.DB.First(user, userID).Error
}
