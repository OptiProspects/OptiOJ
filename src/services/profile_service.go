package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/location"
	"OptiOJ/src/models"
	"errors"
	"fmt"
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
			profile.Birthday = &birthday
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

// GetUserActivity 获取用户活跃度
func GetUserActivity(userID uint, days int) (*models.ActivityResponse, error) {
	// 默认获取过去一年的数据
	if days <= 0 {
		days = 365
	}

	// 计算开始日期
	now := time.Now()
	startDate := now.AddDate(0, 0, -days+1).Format("2006-01-02")

	// 获取每日提交次数和状态
	var dailySubmissions []struct {
		Date  time.Time `gorm:"column:date"`
		Count int       `gorm:"column:count"`
	}

	err := config.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM submissions
		WHERE user_id = ?
		AND created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date ASC
	`, userID, startDate).Scan(&dailySubmissions).Error

	if err != nil {
		return nil, fmt.Errorf("获取提交记录失败: %v", err)
	}

	// 获取90天内的通过率
	var acceptStats struct {
		Total    int64
		Accepted int64
	}
	err = config.DB.Raw(`
		SELECT 
			COUNT(*) as total,
			SUM(CASE WHEN status = 'accepted' THEN 1 ELSE 0 END) as accepted
		FROM submissions
		WHERE user_id = ?
		AND created_at >= DATE_SUB(NOW(), INTERVAL 90 DAY)
	`, userID).Scan(&acceptStats).Error

	if err != nil {
		return nil, fmt.Errorf("获取通过率统计失败: %v", err)
	}

	// 计算通过率
	var acceptRate float64
	if acceptStats.Total > 0 {
		acceptRate = float64(acceptStats.Accepted) / float64(acceptStats.Total) * 100
	}

	// 构建日期到提交次数的映射
	submissionMap := make(map[string]int)
	maxCount := 0
	totalCount := 0
	activeDays := 0 // 有提交的天数
	for _, ds := range dailySubmissions {
		dateStr := ds.Date.Format("2006-01-02")
		submissionMap[dateStr] = ds.Count
		if ds.Count > maxCount {
			maxCount = ds.Count
		}
		if ds.Count > 0 {
			activeDays++
		}
		totalCount += ds.Count
	}

	// 计算平均每天的提交次数（只考虑有提交的天数）
	var avgCount float64
	if activeDays > 0 {
		avgCount = float64(totalCount) / float64(activeDays)
	}

	// 根据平均提交次数计算等级阈值
	lowThreshold := avgCount * 0.5
	mediumThreshold := avgCount
	highThreshold := avgCount * 1.5

	// 只生成有提交记录的日期的活跃度数据
	var activities []models.DailyActivity
	for _, ds := range dailySubmissions {
		dateStr := ds.Date.Format("2006-01-02")
		count := ds.Count

		// 计算活跃度等级
		var level models.ActivityLevel
		switch {
		case count == 0:
			level = models.ActivityLevelNone
		case float64(count) <= lowThreshold:
			level = models.ActivityLevelLow
		case float64(count) <= mediumThreshold:
			level = models.ActivityLevelMedium
		case float64(count) <= highThreshold:
			level = models.ActivityLevelHigh
		default:
			level = models.ActivityLevelSuper
		}

		activities = append(activities, models.DailyActivity{
			Date:  dateStr,
			Count: count,
			Level: level,
		})
	}

	return &models.ActivityResponse{
		Activities: activities,
		TotalDays:  days,
		MaxCount:   maxCount,
		TotalCount: totalCount,
		AcceptRate: acceptRate,
	}, nil
}
