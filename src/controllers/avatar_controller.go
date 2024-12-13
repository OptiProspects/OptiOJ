package controllers

import (
	"OptiOJ/src/services"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func UploadAvatar(c *gin.Context) {
	// 从请求头获取访问令牌
	accessToken := c.GetHeader("Authorization")

	// 验证访问令牌并获取用户ID
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 获取上传的文件
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "获取文件失败"})
		return
	}

	// 保存头像
	avatar, err := services.SaveAvatar(userID, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "头像上传成功",
		"filename": avatar.Filename,
	})
}

func GetAvatar(c *gin.Context) {
	// 从请求头获取访问令牌
	accessToken := c.GetHeader("Authorization")

	// 验证访问令牌并获取用户ID
	userID, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的访问令牌"})
		return
	}

	// 获取头像信息
	avatar, err := services.GetAvatarByUserID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到头像"})
		return
	}

	// 构建文件路径
	filePath := filepath.Join(services.GetAvatarPath(), avatar.Filename)

	// 返回文件
	c.File(filePath)
}
