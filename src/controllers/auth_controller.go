package controllers

import (
	"OptiOJ/src/config"
	"OptiOJ/src/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthController() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取令牌
		accessToken := c.GetHeader("Access-Token")
		refreshToken := c.GetHeader("Refresh-Token")
		userID, err := services.ValidateAccessToken(accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
			return
		}

		// 存储令牌信息
		accessSessionKey := "access_token:" + accessToken
		refreshSessionKey := "refresh_token:" + refreshToken

		// 访问令牌有效期2小时
		if err := config.RedisClient.Set(c, accessSessionKey, userID, 2*time.Hour).Err(); err != nil {
			return
		}

		// 刷新令牌有效期30天
		if err := config.RedisClient.Set(c, refreshSessionKey, userID, 30*24*time.Hour).Err(); err != nil {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

// RefreshToken 刷新访问令牌
func RefreshToken(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供刷新令牌"})
		return
	}

	userID, err := services.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的刷新令牌"})
		return
	}

	// 生成新的访问令牌
	newAccessToken, err := services.GenerateToken(userID, 2*time.Hour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成访问令牌失败"})
		return
	}

	// 存储新的访问令牌
	accessSessionKey := "access_token:" + newAccessToken
	if err := config.RedisClient.Set(c, accessSessionKey, userID, 2*time.Hour).Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "存储访问令牌失败"})
		return
	}

	// 更新会话信息
	if err := services.UpdateSessionLastRefresh(c, refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新会话信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}
