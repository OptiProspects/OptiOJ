package controllers

import (
	"OptiOJ/src/models"
	"OptiOJ/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetLoginHistory 获取登录历史
func GetLoginHistory(c *gin.Context) {
	// 验证管理员权限
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isAdmin, _ := services.IsAdmin(currentUserID)

	var req models.LoginHistoryRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 非管理员只能查看自己的登录历史
	if !isAdmin && req.UserID != 0 && req.UserID != currentUserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 非管理员查询时，强制只能查看自己的记录
	if !isAdmin {
		req.UserID = currentUserID
	}

	logins, total, err := services.GetLoginHistory(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logins":    logins,
		"total":     total,
		"page":      req.Page,
		"page_size": req.PageSize,
	})
}
