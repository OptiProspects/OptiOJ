package controllers

import (
	"OptiOJ/src/models"
	"OptiOJ/src/services"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateTeam 创建团队
func CreateTeam(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.CreateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	teamID, err := services.CreateTeam(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建团队成功",
		"data": gin.H{
			"team_id": teamID,
		},
	})
}

// UpdateTeam 更新团队信息
func UpdateTeam(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	var req models.UpdateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := services.UpdateTeam(teamID, &req, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新团队成功",
	})
}

// DeleteTeam 删除团队
func DeleteTeam(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	if err := services.DeleteTeam(teamID, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除团队成功",
	})
}

// GetTeamDetail 获取团队详情
func GetTeamDetail(c *gin.Context) {
	// 获取用户身份（可选）
	var userID uint
	accessToken := c.GetHeader("Authorization")
	if accessToken != "" {
		id, err := services.ValidateAccessToken(accessToken)
		if err == nil {
			userID = id
		}
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	detail, err := services.GetTeamDetail(teamID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": detail,
	})
}

// GetTeamList 获取团队列表
func GetTeamList(c *gin.Context) {
	// 获取用户身份（可选）
	var userID uint
	accessToken := c.GetHeader("Authorization")
	if accessToken != "" {
		id, err := services.ValidateAccessToken(accessToken)
		if err == nil {
			userID = id
		}
	}

	var req models.TeamListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	response, err := services.GetTeamList(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}

// CreateTeamInvitation 创建团队邀请
func CreateTeamInvitation(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	invitation, err := services.CreateTeamInvitation(teamID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": invitation,
	})
}

// JoinTeam 加入团队
func JoinTeam(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req struct {
		Code string `json:"code" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := services.JoinTeamByInvitation(req.Code, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "加入团队成功",
	})
}

// CreateAssignment 创建团队作业
func CreateAssignment(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.CreateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	assignmentID, err := services.CreateAssignment(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建作业成功",
		"data": gin.H{
			"assignment_id": assignmentID,
		},
	})
}

// UpdateAssignment 更新团队作业
func UpdateAssignment(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的作业ID"})
		return
	}

	var req models.UpdateAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := services.UpdateAssignment(assignmentID, &req, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新作业成功",
	})
}

// CreateProblemList 创建团队题单
func CreateProblemList(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.CreateProblemListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	listID, err := services.CreateProblemList(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "创建题单成功",
		"data": gin.H{
			"list_id": listID,
		},
	})
}

// UpdateProblemList 更新团队题单
func UpdateProblemList(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	listID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题单ID"})
		return
	}

	var req models.UpdateProblemListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := services.UpdateProblemList(listID, &req, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新题单成功",
	})
}

// UpdateTeamMemberRole 更新团队成员角色
func UpdateTeamMemberRole(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	teamID, err := strconv.ParseUint(c.Param("team_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	var req struct {
		UserID uint64 `json:"user_id" binding:"required"`
		Role   string `json:"role" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := services.UpdateTeamMemberRole(teamID, req.UserID, req.Role, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成员角色成功",
	})
}

// RemoveTeamMember 移除团队成员
func RemoveTeamMember(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	teamID, err := strconv.ParseUint(c.Param("team_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	targetUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	if err := services.RemoveTeamMember(teamID, targetUserID, uint64(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "移除成员成功",
	})
}

// GetAssignmentDetail 获取作业详情
func GetAssignmentDetail(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	assignmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的作业ID"})
		return
	}

	assignment, err := services.GetAssignmentDetail(assignmentID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": assignment,
	})
}

// GetAssignmentList 获取作业列表
func GetAssignmentList(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	teamID, err := strconv.ParseUint(c.Query("team_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	assignments, err := services.GetAssignmentList(teamID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": assignments,
	})
}

// GetProblemListDetail 获取题单详情
func GetProblemListDetail(c *gin.Context) {
	// 验证用户身份（可选）
	var userID uint
	accessToken := c.GetHeader("Authorization")
	if accessToken != "" {
		id, err := services.ValidateAccessToken(accessToken)
		if err == nil {
			userID = id
		}
	}

	listID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题单ID"})
		return
	}

	list, err := services.GetProblemListDetail(listID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": list,
	})
}

// GetProblemListList 获取题单列表
func GetProblemListList(c *gin.Context) {
	// 验证用户身份（可选）
	var userID uint
	accessToken := c.GetHeader("Authorization")
	if accessToken != "" {
		id, err := services.ValidateAccessToken(accessToken)
		if err == nil {
			userID = id
		}
	}

	teamID, err := strconv.ParseUint(c.Query("team_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	lists, err := services.GetProblemListList(teamID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": lists,
	})
}

// UploadTeamAvatar 上传团队头像
func UploadTeamAvatar(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 获取团队ID
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	// 检查用户权限
	role, err := services.GetTeamUserRole(teamID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if role != "owner" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}

	// 保存头像
	avatar, err := services.SaveTeamAvatar(teamID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "头像上传成功",
		"data": gin.H{
			"filename": avatar.Filename,
		},
	})
}

// GetTeamAvatar 获取团队头像
func GetTeamAvatar(c *gin.Context) {
	// 获取头像文件名
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文件名"})
		return
	}

	// 获取头像信息
	avatar, err := services.GetTeamAvatar(filename)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到头像"})
		return
	}

	// 构建文件路径
	filePath := filepath.Join(services.GetTeamAvatarPath(), avatar.Filename)

	// 返回文件
	c.File(filePath)
}

// RemoveTeamAvatar 删除团队头像
func RemoveTeamAvatar(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 获取团队ID
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	// 检查用户权限
	role, err := services.GetTeamUserRole(teamID, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if role != "owner" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	// 删除头像
	if err := services.RemoveTeamAvatar(teamID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "头像删除成功",
	})
}

// GetTeamMemberList 获取团队成员列表
func GetTeamMemberList(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 获取团队ID
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	// 绑定请求参数
	var req models.TeamMemberListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取成员列表
	response, err := services.GetTeamMemberList(teamID, &req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}

// UpdateTeamNickname 更新团队内名称
func UpdateTeamNickname(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 获取团队ID
	teamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的团队ID"})
		return
	}

	// 绑定请求参数
	var req models.UpdateTeamNicknameRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 更新团队内名称
	if err := services.UpdateTeamNickname(teamID, uint64(userID), req.Nickname); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新团队内名称成功",
	})
}

// GetAvailableProblemList 获取可用题目列表
func GetAvailableProblemList(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.AvailableProblemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	response, err := services.GetAvailableProblemList(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}

// GetAssignmentProblems 获取作业题目列表
func GetAssignmentProblems(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.GetAssignmentProblemsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	response, err := services.GetAssignmentProblems(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}

// GetAssignmentProblemDetail 获取作业题目详情
func GetAssignmentProblemDetail(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.GetAssignmentProblemDetailRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	detail, err := services.GetAssignmentProblemDetail(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": detail,
	})
}

// SubmitAssignmentCode 提交作业代码
func SubmitAssignmentCode(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.SubmitAssignmentCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	submissionID, err := services.SubmitAssignmentCode(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"submission_id": submissionID,
		},
	})
}

// GetAssignmentSubmissions 获取作业提交记录
func GetAssignmentSubmissions(c *gin.Context) {
	// 验证用户身份
	accessToken := c.GetHeader("Authorization")
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.GetAssignmentSubmissionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	response, err := services.GetAssignmentSubmissions(&req, uint64(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}
