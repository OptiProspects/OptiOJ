package controllers

import (
	"OptiOJ/src/models"
	"OptiOJ/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserList 获取用户列表
func GetUserList(c *gin.Context) {
	// 验证管理员权限
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

	var req models.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	users, total, err := services.GetUserList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"users":     users,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}

// UpdateUser 更新用户信息
func UpdateUser(c *gin.Context) {
	// 验证管理员权限
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

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	var req models.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := services.UpdateUserInfo(uint(id), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户信息更新成功"})
}

// BanUser 封禁用户
func BanUser(c *gin.Context) {
	// 验证管理员权限
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

	var req models.UserBanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := services.BanUser(&req, currentUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已被封禁"})
}

// UnbanUser 解封用户
func UnbanUser(c *gin.Context) {
	// 验证管理员权限
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

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户ID不能为空"})
		return
	}

	id, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := services.UnbanUser(uint(id), currentUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "用户已解封"})
}

// GenerateUsers 批量生成用户
func GenerateUsers(c *gin.Context) {
	// 验证管理员权限
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

	var req models.GenerateUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 调用服务生成用户
	response, err := services.GenerateUsers(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    response,
		"message": "批量生成用户成功",
	})
}
