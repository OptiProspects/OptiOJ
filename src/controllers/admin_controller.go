package controllers

import (
	"OptiOJ/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 添加管理员
func AddAdmin(c *gin.Context) {
	// 验证当前用户是否为超级管理员
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isSuperAdmin, _ := services.IsSuperAdmin(currentUserID)
	if !isSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取请求参数
	var req struct {
		UserID uint   `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 添加管理员
	if err := services.AddAdmin(req.UserID, req.Role); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加管理员成功"})
}

// 移除管理员
func RemoveAdmin(c *gin.Context) {
	// 验证当前用户是否为超级管理员
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isSuperAdmin, _ := services.IsSuperAdmin(currentUserID)
	if !isSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取请求参数
	var req struct {
		UserID uint `json:"user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 移除管理员
	if err := services.RemoveAdmin(req.UserID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "移除管理员成功"})
}

// 获取管理员列表
func GetAdminList(c *gin.Context) {
	// 验证当前用户是否为管理员
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isAdmin, _ := services.IsAdmin(currentUserID)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取管理员列表
	admins, err := services.GetAllAdmins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"admins": admins,
			"total":  len(admins),
		},
		"message": "获取管理员列表成功",
	})
}
