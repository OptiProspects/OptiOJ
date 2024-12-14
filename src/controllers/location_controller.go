package controllers

import (
	"OptiOJ/src/location"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetProvinces(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"provinces": location.Provinces,
	})
}

func GetCities(c *gin.Context) {
	province := c.Query("province")
	if province == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "省份不能为空"})
		return
	}

	// 验证省份是否有效
	if !location.IsValidProvince(province) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的省份"})
		return
	}

	cities := location.GetCities(province)
	c.JSON(http.StatusOK, gin.H{
		"cities": cities,
	})
}
