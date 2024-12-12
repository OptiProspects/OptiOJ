package controllers

import (
	"OptiOJ/src/config"
	"OptiOJ/src/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type VerificationRequest struct {
	RequestValue string `json:"requestValue"`
	RequestType  string `json:"requestType"`
	CaptchaID    string `json:"captchaID"`
	UserExist    bool   `json:"userExist"`
}

func RequestVerification(c *gin.Context) {
	var req VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误"})
		return
	}

	if req.RequestValue == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "值不能为空"})
		return
	}

	// 验证 captchaID
	val, err := config.RedisClient.Get(c, req.CaptchaID).Result()
	if err != nil || val != "geetest:result:success" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的 captchaID"})
		return
	}

	// 检查用户是否存在
	userExists := services.CheckUserExist(req.RequestValue, req.RequestType)

	var code string
	code = services.GenerateVerificationCode()

	if req.RequestType == "email" {
		if err := services.SendVerificationCode(req.RequestValue, code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "验证码已发送到邮箱", "userExist": userExists})
	} else if req.RequestType == "phone" {
		if err := services.SendVerificationCodeToPhone(req.RequestValue, code); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "验证码已发送到手机号", "userExist": userExists})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的类型"})
		return
	}

	// 存储验证码到 Redis
	redisKey := "verification:" + req.RequestValue + ":" + req.RequestType
	if err := config.RedisClient.Set(c, redisKey, code, 5*time.Minute).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储验证码失败"})
		return
	}
}
