package controllers

import (
	"OptiOJ/src/models"
	"OptiOJ/src/services"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProblem 创建题目
func CreateProblem(c *gin.Context) {
	// 验证管理员权限
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isAdmin, _ := services.IsAdmin(uint(currentUserID))
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	var req models.CreateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	problemID, err := services.CreateProblem(&req, uint64(currentUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    gin.H{"problem_id": problemID},
		"message": "创建题目成功",
	})
}

// UpdateProblem 更新题目
func UpdateProblem(c *gin.Context) {
	// 验证管理员权限
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isAdmin, _ := services.IsAdmin(uint(currentUserID))
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	problemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	var req models.UpdateProblemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if err := services.UpdateProblem(problemID, &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新题目成功",
	})
}

// DeleteProblem 删除题目
func DeleteProblem(c *gin.Context) {
	// 验证管理员权限
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isAdmin, _ := services.IsAdmin(uint(currentUserID))
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	problemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	if err := services.DeleteProblem(problemID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除题目成功",
	})
}

// GetProblemDetail 获取题目详情
func GetProblemDetail(c *gin.Context) {
	// 获取当前用户ID
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	problemID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	problem, err := services.GetProblemDetail(problemID, uint64(currentUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": problem,
	})
}

// GetProblemList 获取题目列表
func GetProblemList(c *gin.Context) {
	// 获取当前用户ID
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	var req models.ProblemListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	response, err := services.GetProblemList(&req, uint64(currentUserID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": response,
	})
}

// UploadTestCase 上传测试用例
func UploadTestCase(c *gin.Context) {
	// 验证管理员权限
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isAdmin, _ := services.IsAdmin(uint(currentUserID))
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	var req models.TestCaseUploadRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取上传的文件
	inputFile, err := c.FormFile("input")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传输入文件"})
		return
	}
	outputFile, err := c.FormFile("output")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请上传输出文件"})
		return
	}

	// 创建临时文件
	tempDir := filepath.Join(os.TempDir(), "optioj_testcases")
	os.MkdirAll(tempDir, 0755)

	inputTempFile := filepath.Join(tempDir, inputFile.Filename)
	if err := c.SaveUploadedFile(inputFile, inputTempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存输入文件失败"})
		return
	}
	defer os.Remove(inputTempFile)

	outputTempFile := filepath.Join(tempDir, outputFile.Filename)
	if err := c.SaveUploadedFile(outputFile, outputTempFile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存输出文件失败"})
		return
	}
	defer os.Remove(outputTempFile)

	// 打开临时文件
	input, err := os.Open(inputTempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开输入文件失败"})
		return
	}
	defer input.Close()

	output, err := os.Open(outputTempFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "打开输出文件失败"})
		return
	}
	defer output.Close()

	// 上传测试用例
	if err := services.UploadTestCase(req.ProblemID, input, output, req.IsSample); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "上传测试用例成功",
	})
}

// DeleteTestCase 删除测试用例
func DeleteTestCase(c *gin.Context) {
	// 验证管理员权限
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isAdmin, _ := services.IsAdmin(uint(currentUserID))
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	testCaseID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的测试用例ID"})
		return
	}

	if err := services.DeleteTestCase(testCaseID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除测试用例成功",
	})
}

// GetTestCases 获取题目的测试用例列表
func GetTestCases(c *gin.Context) {
	// 验证管理员权限
	accessToken := c.GetHeader("Authorization")
	currentUserID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	isAdmin, _ := services.IsAdmin(uint(currentUserID))
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "权限不足"})
		return
	}

	problemID, err := strconv.ParseUint(c.Param("problem_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的题目ID"})
		return
	}

	testCases, err := services.GetTestCases(problemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": testCases,
	})
}
