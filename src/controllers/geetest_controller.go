package controllers

import (
	"OptiOJ/src/config"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 请求体结构体
type GeetestRequest struct {
	CaptchaID     string `json:"captcha_id"`
	LotNumber     string `json:"lot_number"`
	PassToken     string `json:"pass_token"`
	GenTime       string `json:"gen_time"`
	CaptchaOutput string `json:"captcha_output"`
}

// Geetest 验证接口
func ValidateGeetest(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "不支持的请求方法"})
		return
	}

	// 从请求 Body 中获取参数
	var req GeetestRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求体格式错误"})
		return
	}

	// 打印请求参数以便调试
	fmt.Println("请求参数: ", req)

	// 生成签名
	signToken := hmacEncode(config.Geetest.CaptchaKey, req.LotNumber)

	// 向 Geetest 转发前端数据 + “sign_token” 签名
	formData := url.Values{}
	formData.Set("lot_number", req.LotNumber)
	formData.Set("captcha_output", req.CaptchaOutput)
	formData.Set("pass_token", req.PassToken)
	formData.Set("gen_time", req.GenTime)
	formData.Set("sign_token", signToken)

	// 发起 POST 请求
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.PostForm(config.Geetest.CaptchaURL+"/validate"+"?captcha_id="+config.Geetest.CaptchaID, formData)
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println("服务接口异常: ", err)
		c.JSON(http.StatusOK, gin.H{"result": "success"})
		return
	}

	// 解析 Geetest 返回的 JSON 数据
	var resMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&resMap); err != nil {
		fmt.Println("Json数据解析错误")
		c.JSON(http.StatusOK, gin.H{"result": "success"})
		return
	}

	// 打印返回的完整数据以便调试
	fmt.Println("Geetest 返回的数据: ", resMap)

	// 根据 Geetest 返回的用户验证状态进行业务逻辑处理
	result := resMap["result"]
	if result == "success" {
		fmt.Println("验证通过")
		// 生成唯一 ID
		requestID := uuid.New().String()
		// 将验证结果保存到 Redis
		err := config.RedisClient.Set(c, requestID, "geetest:result:success", 5*time.Minute).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "存储验证结果失败"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"result": "success", "requestID": requestID})
	} else {
		reason := resMap["reason"]
		fmt.Println("验证失败: ", reason)
		c.JSON(http.StatusOK, gin.H{"result": "fail", "reason": reason})
	}
}

// hmac-sha256 加密
func hmacEncode(key string, data string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}
