package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"time"

	"github.com/gin-gonic/gin"
)

// RecordLogin 记录用户登录
func RecordLogin(c *gin.Context, userID uint, status string, failReason string) error {
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()

	login := models.UserLogin{
		UserID:      uint64(userID),
		LoginTime:   time.Now(),
		IPAddress:   ipAddress,
		UserAgent:   userAgent,
		LoginStatus: status,
		FailReason:  failReason,
		// Location 字段可以通过 IP 地理位置服务获取，这里暂时留空
	}

	return config.DB.Create(&login).Error
}

// GetLoginHistory 获取登录历史
func GetLoginHistory(req *models.LoginHistoryRequest) ([]models.UserLogin, int64, error) {
	var logins []models.UserLogin
	var total int64

	query := config.DB.Model(&models.UserLogin{})

	if req.UserID != 0 {
		query = query.Where("user_id = ?", req.UserID)
	}

	if !req.StartTime.IsZero() {
		query = query.Where("login_time >= ?", req.StartTime)
	}

	if !req.EndTime.IsZero() {
		query = query.Where("login_time <= ?", req.EndTime)
	}

	if req.Status != "" {
		query = query.Where("login_status = ?", req.Status)
	}

	// 获取总数
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	err = query.Order("login_time DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&logins).Error

	if err != nil {
		return nil, 0, err
	}

	return logins, total, nil
}
