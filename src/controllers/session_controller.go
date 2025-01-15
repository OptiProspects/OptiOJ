package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"OptiOJ/src/services"
)

// GetActiveSessions 获取当前用户的所有活跃会话
func GetActiveSessions(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	sessions, err := services.GetActiveSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"sessions": sessions,
	})
}

// RevokeSession 吊销指定会话
func RevokeSession(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供会话ID"})
		return
	}

	// 验证会话是否属于当前用户
	sessions, err := services.GetActiveSessions(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取会话信息失败"})
		return
	}

	var found bool
	for _, session := range sessions {
		if session.SessionID == sessionID {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作此会话"})
		return
	}

	if err := services.RevokeSession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "吊销会话失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成功吊销会话",
	})
}

// Logout 退出当前设备
func Logout(c *gin.Context) {
	refreshToken := c.GetHeader("Authorization")
	if refreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未提供刷新令牌"})
		return
	}

	// 验证刷新令牌有效性
	if _, err := services.ValidateRefreshToken(refreshToken); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的刷新令牌"})
		return
	}

	if err := services.Logout(refreshToken); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退出登录失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成功退出登录",
	})
}

// LogoutAllDevices 退出所有设备
func LogoutAllDevices(c *gin.Context) {
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	if err := services.LogoutAllDevices(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退出所有设备失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "成功退出所有设备",
	})
}
