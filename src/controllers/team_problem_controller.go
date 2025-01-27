package controllers

import (
	"OptiOJ/src/models"
	"OptiOJ/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTeamProblem 创建团队私有题目
func CreateTeamProblem(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.CreateTeamProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	problemID, err := services.CreateTeamProblem(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建题目成功",
		"data": gin.H{
			"problem_id": problemID,
		},
	})
}

// UpdateTeamProblem 更新团队私有题目
func UpdateTeamProblem(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	problemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	var req models.UpdateTeamProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := services.UpdateTeamProblem(problemID, &req, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新题目成功",
	})
}

// DeleteTeamProblem 删除团队私有题目
func DeleteTeamProblem(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	problemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	if err := services.DeleteTeamProblem(problemID, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除题目成功",
	})
}

// GetTeamProblemDetail 获取团队私有题目详情
func GetTeamProblemDetail(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	problemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	detail, err := services.GetTeamProblemDetail(problemID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": detail,
	})
}

// GetTeamProblemList 获取团队私有题目列表
func GetTeamProblemList(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.TeamProblemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	response, err := services.GetTeamProblemList(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}
