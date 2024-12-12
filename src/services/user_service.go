package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(user *models.User) error {
	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return errors.New("用户名已存在")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 存储用户信息到数据库
	user.Password = string(hashedPassword) // 确保存储的是哈希密码
	if err := config.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}
