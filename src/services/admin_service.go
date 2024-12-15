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

// 获取所有管理员列表
func GetAllAdmins() ([]models.Admin, error) {
	var admins []models.Admin
	if err := config.DB.Find(&admins).Error; err != nil {
		return nil, errors.New("获取管理员列表失败")
	}
	return admins, nil
}
