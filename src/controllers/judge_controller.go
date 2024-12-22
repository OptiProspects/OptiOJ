package controllers

import (
	"OptiOJ/src/models"
	"OptiOJ/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SubmitCode 提交代码
func SubmitCode(c *gin.Context) {
	// 获取当前用户ID
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.SubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	submissionID, err := services.CreateSubmission(&req, uint64(currentUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    gin.H{"submission_id": submissionID},
		"message": "提交成功",
	})
}

// GetSubmissionList 获取提交记录列表
func GetSubmissionList(c *gin.Context) {
	// 获取当前用户ID
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.SubmissionListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查是否为管理员
	isAdmin, _ := services.IsAdmin(uint(currentUserID))
	if !isAdmin {
		// 非管理员只能查看自己的提交记录
		userID := uint64(currentUserID) // 转换为 uint64
		req.UserID = &userID
	}

	response, err := services.GetSubmissionList(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}

// GetSubmissionDetail 获取提交记录详情
func GetSubmissionDetail(c *gin.Context) {
	// 获取当前用户ID
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	submissionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的提交ID"})
		return
	}

	detail, err := services.GetSubmissionDetail(submissionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 检查访问权限
	isAdmin, _ := services.IsAdmin(uint(currentUserID))
	if !isAdmin && detail.UserID != uint64(currentUserID) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权访问该提交记录"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": detail,
	})
}
