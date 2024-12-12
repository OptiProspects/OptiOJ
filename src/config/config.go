package config

import (
	"context"
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
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
}

type GeetestConfig struct {
	CaptchaURL string
	CaptchaID  string
	CaptchaKey string
}

var DB *gorm.DB
var SMTP SMTPConfig
var Geetest GeetestConfig
var RedisClient *redis.Client
var logger = logrus.New()
var ctx = context.Background()

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
			return
		}

		logger.Errorf("数据库连接失败: %v", err)
		logger.Info("等待 5 秒后重试...")
		time.Sleep(5 * time.Second)
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
	Geetest = config.Geetest
}
