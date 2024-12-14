package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func GetAvatarPath() string {
	// 根据操作系统选择正确的路径分隔符
	baseDir := filepath.Join("data", "user", "avatar")
	// 确保目录存在
	os.MkdirAll(baseDir, 0755)
	return baseDir
}

func SaveAvatar(userID uint, file *multipart.FileHeader) (*models.Avatar, error) {
	// 获取文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return nil, fmt.Errorf("不支持的文件格式，仅支持 jpg、jpeg 和 png")
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	filePath := filepath.Join(GetAvatarPath(), filename)

	// 保存文件
	if err := os.MkdirAll(GetAvatarPath(), 0755); err != nil {
		return nil, fmt.Errorf("创建目录失败: %v", err)
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("创建目标文件失败: %v", err)
	}
	defer dst.Close()

	// 读取文件内容并写入新文件
	buffer := make([]byte, 1024*1024) // 1MB buffer
	for {
		n, err := src.Read(buffer)
		if err != nil {
			break
		}
		if _, err := dst.Write(buffer[:n]); err != nil {
			return nil, fmt.Errorf("写入文件失败: %v", err)
		}
	}

	// 更新数据库
	avatar := &models.Avatar{
		UserID:     int(userID),
		Filename:   filename,
		UploadTime: time.Now(),
	}

	// 删除旧头像
	var oldAvatar models.Avatar
	if err := config.DB.Where("user_id = ?", userID).First(&oldAvatar).Error; err == nil {
		// 删除旧文件
		oldPath := filepath.Join(GetAvatarPath(), oldAvatar.Filename)
		os.Remove(oldPath)
		// 更新记录
		if err := config.DB.Model(&oldAvatar).Updates(avatar).Error; err != nil {
			return nil, fmt.Errorf("更新头像记录失败: %v", err)
		}
		return avatar, nil
	}

	// 创建新记录
	if err := config.DB.Create(avatar).Error; err != nil {
		return nil, fmt.Errorf("保存头像记录失败: %v", err)
	}

	return avatar, nil
}

func GetAvatarByUserID(userID uint) (*models.Avatar, error) {
	var avatar models.Avatar
	if err := config.DB.Where("user_id = ?", userID).First(&avatar).Error; err != nil {
		return nil, err
	}
	return &avatar, nil
}

func RemoveAvatar(userID uint) error {
	var avatar models.Avatar
	// 查找用户的头像记录
	if err := config.DB.Where("user_id = ?", userID).First(&avatar).Error; err != nil {
		return fmt.Errorf("未找到头像记录")
	}

	// 删除文件
	filePath := filepath.Join(GetAvatarPath(), avatar.Filename)
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("删除头像文件失败: %v", err)
	}

	// 删除数据库记录
	if err := config.DB.Delete(&avatar).Error; err != nil {
		return fmt.Errorf("删除头像记录失败: %v", err)
	}

	return nil
}
