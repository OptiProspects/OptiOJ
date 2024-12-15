package config

import (
	"OptiOJ/src/models"
	"context"
	"crypto/rand"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Database DatabaseConfig
	SMTP     SMTPConfig
	Redis    RedisConfig
	Aliyun   AliyunConfig
	Geetest  GeetestConfig
}

type DatabaseConfig struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
}

type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Sender   string
	UseTLS   bool
	UseSSL   bool
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type AliyunConfig struct {
	AccessKeyId     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

type GeetestConfig struct {
	CaptchaURL string
	CaptchaID  string
	CaptchaKey string
}

var DB *gorm.DB
var SMTP SMTPConfig
var Aliyun AliyunConfig
var Geetest GeetestConfig
var RedisClient *redis.Client
var logger = logrus.New()
var ctx = context.Background()
var JWTSecret []byte

func InitDB() {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		logger.Fatal(err)
	}

	dsn := config.Database.User + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port) + ")/" + config.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"

	var err error
	for {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			logger.Infof("数据库连接成功: %s@%s:%d/%s", config.Database.User, config.Database.Host, config.Database.Port, config.Database.DBName)
			break
		}

		logger.Errorf("数据库连接失败: %v", err)
		logger.Info("等待 5 秒后重试...")
		time.Sleep(5 * time.Second)
	}

	// 检查并添加第一个管理员用户
	var count int64
	err = DB.Model(&models.Admin{}).Count(&count).Error
	if err != nil {
		logger.Fatal("检查管理员用户失败:", err)
	}

	if count == 0 {
		// 创建第一个管理员用户
		admin := models.Admin{
			UserID:    1, // 假设第一个用户的 ID 为 1
			Role:      "super_admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := DB.Create(&admin).Error; err != nil {
			logger.Fatal("添加第一个管理员用户失败:", err)
		}
		logger.Info("第一个管理员用户已添加")
	}
}

func InitRedis() {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		logger.Fatal(err)
	}

	// 构建 Redis 连接字符串
	redisOptions := &redis.Options{
		Addr: config.Redis.Host + ":" + strconv.Itoa(config.Redis.Port),
		DB:   0, // 默认 DB
	}

	if config.Redis.Password != "" {
		redisOptions.Password = config.Redis.Password // 设置密码
	}

	var err error
	for {
		RedisClient = redis.NewClient(redisOptions)

		// 测试连接
		_, err = RedisClient.Ping(ctx).Result()
		if err == nil {
			logger.Infof("Redis 连接成功: %s", redisOptions.Addr)
			return
		}

		logger.Errorf("Redis 连接失败: %v", err)
		logger.Info("等待 5 秒后重试...")
		time.Sleep(5 * time.Second)
	}
}

func InitConfig() {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		logger.Fatal(err)
	}

	SMTP = config.SMTP
	Aliyun = config.Aliyun
	Geetest = config.Geetest

	// 初始化 JWT 密钥
	InitJWTSecret()
}

func InitJWTSecret() {
	keyFile := "jwtKey"

	// 尝试读取现有的 key 文件
	content, err := os.ReadFile(keyFile)
	if err == nil && len(content) > 0 {
		// 如果文件存在且不为空，使用现有密钥
		JWTSecret = content
		logger.Info("JWT密钥加载成功")
		return
	}

	// 如果文件不存在或为空，生成新的密钥
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		logger.Fatal("生成JWT密钥失败:", err)
	}

	// 将密钥写入文件
	if err := os.WriteFile(keyFile, secret, 0600); err != nil {
		logger.Fatal("写入密钥文件失败:", err)
	}

	JWTSecret = secret
	logger.Info("已生成新的JWT密钥并保存到key文件")
}
