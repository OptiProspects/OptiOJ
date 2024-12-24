package controllers

import (
	"OptiOJ/src/models"
	"OptiOJ/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UpdateProfile(c *gin.Context) {
	// 从请求头获取访问令牌
	accessToken := c.GetHeader("Authorization")

	// 验证访问令牌并获取用户ID
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求格式错误"})
		return
	}

	// 更新用户资料
	if err := services.UpdateProfile(userID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "个人资料更新成功",
	})
}

// GetUserActivity 获取用户活跃度
func GetUserActivity(c *gin.Context) {
	// 获取目标用户ID
	targetUserID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	// 获取天数参数
	days, err := strconv.Atoi(c.DefaultQuery("days", "365"))
	if err != nil || days <= 0 {
		days = 365 // 默认获取一年的数据
	}

	// 获取活跃度数据
	activity, err := services.GetUserActivity(uint(targetUserID), days)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": activity,
	})
}
