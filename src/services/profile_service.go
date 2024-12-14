package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/location"
	"OptiOJ/src/models"
	"errors"
	"strings"
	"time"
)

func GetProfile(userID uint) (*models.Profile, error) {
	var profile models.Profile
	if err := config.DB.Where("user_id = ?", userID).First(&profile).Error; err != nil {
		return nil, err
	}

	// 解析 Location 字段
	if profile.Location != "" {
		parts := strings.Split(profile.Location, "-")
		if len(parts) == 2 {
			profile.Province = parts[0]
			profile.City = parts[1]
		}
	}

	return &profile, nil
}

func UpdateProfile(userID uint, req *models.UpdateProfileRequest) error {
	var profile models.Profile

	// 检查性别值是否有效
	if req.Gender != "" && req.Gender != "male" && req.Gender != "female" && req.Gender != "other" {
		return errors.New("无效的性别值")
	}

	// 检查个人签名长度
	if len(req.Bio) > 255 {
		return errors.New("个人签名过长")
	}

	// 检查学校名称长度
	if len(req.School) > 100 {
		return errors.New("学校名称过长")
	}

	// 检查真实姓名长度
	if len(req.RealName) > 50 {
		return errors.New("真实姓名过长")
	}

	// 处理生日
	var birthday time.Time
	if req.Birthday != "" {
		// 尝试解析带时区的时间格式
		layouts := []string{
			"2006-01-02T15:04:05Z07:00", // ISO 8601 格式
			"2006-01-02 15:04:05Z07:00", // 带时区的标准格式
			"2006-01-02 15:04:05",       // 不带时区的标准格式
			"2006-01-02",                // 仅日期格式
		}

		var parseErr error
		for _, layout := range layouts {
			birthday, parseErr = time.Parse(layout, req.Birthday)
			if parseErr == nil {
				break
			}
		}

		if parseErr != nil {
			return errors.New("生日格式错误，支持的格式：YYYY-MM-DD, YYYY-MM-DD HH:mm:ss 或带时区的ISO 8601格式")
		}

		// 如果输入的时间没有时区信息，设置为本地时区
		if birthday.Location() == time.UTC {
			birthday = birthday.In(time.Local)
		}

		if birthday.After(time.Now()) {
			return errors.New("生日不能是未来日期")
		}
		minDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local)
		if birthday.Before(minDate) {
			return errors.New("生日日期无效")
		}
	}

	// 处理地址信息
	var locationStr string
	if req.Province != "" || req.City != "" {
		// 如果提供了省份或城市，则两者都必须提供
		if req.Province == "" || req.City == "" {
			return errors.New("省份和城市必须同时提供")
		}

		// 验证省份
		if !location.IsValidProvince(req.Province) {
			return errors.New("无效的省份")
		}

		// 验证城市
		if !location.IsValidCity(req.Province, req.City) {
			return errors.New("无效的城市")
		}

		// 组合地址
		locationStr = req.Province + "-" + req.City
	}

	result := config.DB.Where("user_id = ?", userID).First(&profile)
	if result.Error != nil {
		// 如果记录不存在,创建新记录
		profile = models.Profile{
			UserID: int(userID),
		}
		// 只设置请求中包含的字段
		if req.Bio != "" {
			profile.Bio = req.Bio
		}
		if req.Gender != "" {
			profile.Gender = req.Gender
		}
		if req.School != "" {
			profile.School = req.School
		}
		if !birthday.IsZero() {
			profile.Birthday = birthday
		}
		if locationStr != "" {
			profile.Location = locationStr
		}
		if req.RealName != "" {
			profile.RealName = req.RealName
		}
		return config.DB.Create(&profile).Error
	}

	// 构建更新字段的map
	updates := make(map[string]interface{})

	// 只更新请求中包含的非空字段
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}
	if req.Gender != "" {
		updates["gender"] = req.Gender
	}
	if req.School != "" {
		updates["school"] = req.School
	}
	if !birthday.IsZero() {
		updates["birthday"] = birthday
	}
	if locationStr != "" {
		updates["location"] = locationStr
	}
	if req.RealName != "" {
		updates["real_name"] = req.RealName
	}

	// 如果没有需要更新的字段，直接返回
	if len(updates) == 0 {
		return nil
	}

	// 更新现有记录
	return config.DB.Model(&profile).Updates(updates).Error
}
