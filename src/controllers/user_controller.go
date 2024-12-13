package controllers

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"OptiOJ/src/services"

	"net/http"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
)

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
	var req models.RegisterRequest
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
	userID, err := services.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 生成访问令牌和刷新令牌
	accessToken, refreshToken, err := services.GenerateTokenPair(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 存储令牌信息到 Redis
	accessSessionKey := "access_token:" + accessToken
	refreshSessionKey := "refresh_token:" + refreshToken

	// 访问令牌有效期2小时
	if err := config.RedisClient.Set(c, accessSessionKey, userID, 2*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储访问令牌失败"})
		return
	}

	// 刷新令牌有效期30天
	if err := config.RedisClient.Set(c, refreshSessionKey, userID, 30*24*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储刷新令牌失败"})
		return
	}

	// 注册成功后删除验证码
	config.RedisClient.Del(c, redisKey)

	c.JSON(http.StatusOK, gin.H{
		"message":       "用户注册成功",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func LoginUser(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 验证用户名和密码
	user, err := services.ValidateLogin(req.AccountInfo, req.PassWord)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// 生成访问令牌和刷新令牌
	accessToken, refreshToken, err := services.GenerateTokenPair(uint(user.ID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成令牌失败"})
		return
	}

	// 存储令牌信息到 Redis
	accessSessionKey := "access_token:" + accessToken
	refreshSessionKey := "refresh_token:" + refreshToken

	// 访问令牌有效期2小时
	if err := config.RedisClient.Set(c, accessSessionKey, user.ID, 2*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储访问令牌失败"})
		return
	}

	// 刷新令牌有效期30天
	if err := config.RedisClient.Set(c, refreshSessionKey, user.ID, 30*24*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储刷新令牌失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func GetGlobalData(c *gin.Context) {
	// 从请求头获取访问令牌
	accessToken := c.GetHeader("Authorization")

	// 验证访问令牌并获取用户ID
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 查询用户信息
	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取头像信息
	avatar, _ := services.GetAvatarByUserID(userID)
	var avatarFilename string
	if avatar != nil {
		avatarFilename = avatar.Filename
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"phone":    user.Phone,
			"avatar":   avatarFilename,
		},
	})
}
