package controllers

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"OptiOJ/src/services"

	"net/http"
	"unicode"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	UserName         string `json:"userName"`
	PassWord         string `json:"passWord"`
	RequestEmail     string `json:"requestEmail"`
	RequestPhone     string `json:"requestPhone"`
	VerificationCode string `json:"verificationCode"`
	VerificationType string `json:"verificationType"` // "email" 或 "phone"
}

// 密码强度验证
func validatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}

	var hasNumber, hasLetter, hasSymbol, hasUpper, hasLower bool
	for _, char := range password {
		switch {
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsLetter(char):
			hasLetter = true
			if unicode.IsUpper(char) {
				hasUpper = true
			}
			if unicode.IsLower(char) {
				hasLower = true
			}
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSymbol = true
		}
	}

	// 计算满足的条件数
	conditions := 0
	if hasNumber && hasLetter {
		conditions++
	}
	if hasSymbol {
		conditions++
	}
	if hasUpper && hasLower {
		conditions++
	}

	return conditions >= 2
}

func RegisterUser(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 基本参数验证
	if req.UserName == "" || req.PassWord == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名和密码不能为空"})
		return
	}

	// 密码强度验证
	if !validatePassword(req.PassWord) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "密码强度不足，需要满足以下条件中的两项：1. 密码长度至少8位且同时包含数字和字母；2. 包含特殊符号；3. 同时包含大小写字母"})
		return
	}

	if req.VerificationType != "email" && req.VerificationType != "phone" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的验证类型"})
		return
	}

	// 根据验证类型验证对应字段
	var verificationValue string
	if req.VerificationType == "email" {
		if req.RequestEmail == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不能为空"})
			return
		}
		verificationValue = req.RequestEmail
	} else {
		if req.RequestPhone == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "手机号不能为空"})
			return
		}
		verificationValue = req.RequestPhone
	}

	// 验证验证码
	redisKey := "verification:" + verificationValue + ":" + req.VerificationType
	val, err := config.RedisClient.Get(c, redisKey).Result()
	if err != nil || val != req.VerificationCode {
		c.JSON(http.StatusBadRequest, gin.H{"error": "验证码无效或已过期"})
		return
	}

	// 创建用户对象
	user := &models.User{
		Username: req.UserName,
		Password: req.PassWord,
		Email:    req.RequestEmail,
		Phone:    req.RequestPhone,
	}

	// 注册用户
	if err := services.RegisterUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 注册成功后删除验证码
	config.RedisClient.Del(c, redisKey)

	c.JSON(http.StatusOK, gin.H{"message": "用户注册成功"})
}
