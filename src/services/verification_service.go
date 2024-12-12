package services

import (
	"OptiOJ/src/config"
	"OptiOJ/src/models"
	"errors"
	"fmt"
	"strconv"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/wneessen/go-mail"
	"golang.org/x/exp/rand"
)

// 生成验证码
func GenerateVerificationCode() string {
	rand.Seed(uint64(time.Now().UnixNano()))
	code := rand.Intn(999999-100000) + 100000 // 生成六位随机数
	return strconv.Itoa(code)
}

// 发送验证码到邮箱
func SendVerificationCode(email string, code string) error {
	smtpHost := config.SMTP.Host
	smtpPort := config.SMTP.Port
	user := config.SMTP.User
	sender := config.SMTP.Sender
	password := config.SMTP.Password
	useTLS := config.SMTP.UseTLS

	m := mail.NewMsg()
	if err := m.From(sender); err != nil {
		return errors.New("设置发件人失败: " + err.Error())
	}
	if err := m.To(email); err != nil {
		return errors.New("设置收件人失败: " + err.Error())
	}
	m.Subject("验证码")
	m.SetBodyString("text/plain", "您的验证码是: "+code)

	c, err := mail.NewClient(smtpHost, mail.WithTLSPortPolicy(mail.TLSMandatory), mail.WithUsername(user), mail.WithPassword(password), mail.WithPort(smtpPort))
	if err != nil {
		return errors.New("创建邮件客户端失败: " + err.Error())
	}

	// 设置 TLS 选项
	if useTLS {
		c.SetSSLPort(true, true) // 启用 SSL/TLS
	} else {
		c.SetSSLPort(false, false) // 禁用 SSL/TLS
	}

	// 发送邮件
	if err := c.DialAndSend(m); err != nil {
		return errors.New("发送验证码失败: " + err.Error())
	}

	return nil
}

// 创建阿里云短信客户端
func CreateClient() (*dysmsapi20170525.Client, error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(config.Aliyun.AccessKeyId),
		AccessKeySecret: tea.String(config.Aliyun.AccessKeySecret),
		Endpoint:        tea.String("dysmsapi.aliyuncs.com"),
	}
	client, err := dysmsapi20170525.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("创建短信客户端失败: %v", err)
	}

	if client.Credential == nil {
		return nil, fmt.Errorf("短信初始化失败，请检查 config.toml 文件")
	}

	return client, nil
}

// 发送验证码到手机号
func SendVerificationCodeToPhone(phone string, code string) error {
	client, err := CreateClient()
	if err != nil {
		return fmt.Errorf("创建短信客户端失败: %v", err)
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  tea.String(phone),
		SignName:      tea.String(config.Aliyun.SignName),
		TemplateCode:  tea.String(config.Aliyun.TemplateCode),
		TemplateParam: tea.String(fmt.Sprintf(`{"code":"%s"}`, code)),
	}

	runtime := &util.RuntimeOptions{}
	_, err = client.SendSmsWithOptions(sendSmsRequest, runtime)
	if err != nil {
		return fmt.Errorf("发送短信失败: %v", err)
	}

	return nil
}

// 检查用户是否存在
func CheckUserExist(value string, userType string) bool {
	var count int64

	if userType == "email" {
		// 检查邮箱是否存在
		result := config.DB.Model(&models.User{}).Where("email = ?", value).Count(&count)
		if result.Error != nil {
			return false // 查询出错，返回 false
		}
	} else if userType == "phone" {
		// 检查手机号是否存在
		result := config.DB.Model(&models.User{}).Where("phone = ?", value).Count(&count)
		if result.Error != nil {
			return false // 查询出错，返回 false
		}
	} else {
		return false // 无效的类型
	}

	return count > 0 // 如果计数大于 0，表示用户存在
}
