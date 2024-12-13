package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user *models.User) (uint, error) {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return 0, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if user.Email != "" {
		if err := config.DB.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
			return 0, errors.New("邮箱已被注册")
		}
	}

	// 检查手机号是否已存在
	if user.Phone != "" {
		if err := config.DB.Where("phone = ?", user.Phone).First(&existingUser).Error; err == nil {
			return 0, errors.New("手机号已被注册")
		}
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	// 存储用户信息到数据库
	user.Password = string(hashedPassword)
	if err := config.DB.Create(user).Error; err != nil {
		return 0, err
	}

	return uint(user.ID), nil
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
