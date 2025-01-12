package config

import (
	"OptiOJ/src/models"
	"context"
	"crypto/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

func CheckAndInitializeDatabase() {
	// 检查数据库是否需要初始化
	var tableCount int
	row := DB.Raw("SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = ?", DB.Migrator().CurrentDatabase()).Row()
	if err := row.Scan(&tableCount); err != nil {
		logger.Fatal("检查数据库表失败:", err)
	}

	// 如果数据库中没有表，则执行初始化
	if tableCount == 0 {
		logger.Info("检测到新数据库，开始执行初始化脚本...")

		// 获取 sql 目录路径
		sqlDir := "sql"
		if _, err := os.Stat(sqlDir); os.IsNotExist(err) {
			// 如果当前目录下没有 sql 目录，尝试上级目录（针对 Docker 环境）
			sqlDir = "../sql"
			if _, err := os.Stat(sqlDir); os.IsNotExist(err) {
				logger.Fatal("找不到 SQL 脚本目录")
			}
		}

		// 读取 SQL 目录下的所有文件
		entries, err := os.ReadDir(sqlDir)
		if err != nil {
			logger.Fatalf("读取 SQL 目录失败: %v", err)
		}

		// 定义表的依赖关系和执行顺序
		sqlFileOrder := []string{
			"users.sql",        // 基础表，无依赖
			"profile.sql",      // 依赖 users
			"admin.sql",        // 依赖 users
			"banned.sql",       // 依赖 users
			"loginHistory.sql", // 依赖 users
			"problems.sql",     // 基础表，无依赖
			"teams.sql",        // 依赖 users, problems
			"judge.sql",        // 依赖 users, problems
			"messages.sql",     // 依赖 users
		}

		// 创建文件名到路径的映射
		fileMap := make(map[string]string)
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
				fileMap[entry.Name()] = filepath.Join(sqlDir, entry.Name())
			}
		}

		// 按照指定顺序执行 SQL 文件
		for _, fileName := range sqlFileOrder {
			filePath, exists := fileMap[fileName]
			if !exists {
				logger.Warnf("未找到 SQL 文件: %s，跳过", fileName)
				continue
			}

			sqlContent, err := os.ReadFile(filePath)
			if err != nil {
				logger.Fatalf("读取 SQL 文件 %s 失败: %v", fileName, err)
			}

			// 分割 SQL 语句并逐个执行
			statements := strings.Split(string(sqlContent), ";")
			for _, stmt := range statements {
				// 跳过空语句
				stmt = strings.TrimSpace(stmt)
				if stmt == "" {
					continue
				}

				if err := DB.Exec(stmt).Error; err != nil {
					logger.Fatalf("执行 SQL 语句失败 [%s]: %v\n语句内容: %s", fileName, err, stmt)
				}
			}
			logger.Infof("成功执行 SQL 文件: %s", fileName)
		}

		logger.Info("数据库初始化完成")
	} else {
		logger.Info("数据库已存在，跳过初始化")
	}
}

func InitDB() {
	var config Config
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		logger.Fatal(err)
	}

	// 首先尝试连接 MySQL 服务器（不指定数据库）
	rootDSN := config.Database.User + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port) + ")/?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"

	rootDB, err := gorm.Open(mysql.Open(rootDSN), &gorm.Config{})
	if err != nil {
		logger.Fatalf("连接 MySQL 服务器失败: %v", err)
	}

	// 检查数据库是否存在
	var count int64
	result := rootDB.Raw("SELECT COUNT(*) FROM information_schema.SCHEMATA WHERE SCHEMA_NAME = ?", config.Database.DBName).Scan(&count)
	if result.Error != nil {
		logger.Fatalf("检查数据库是否存在失败: %v", result.Error)
	}

	// 如果数据库不存在，创建它
	if count == 0 {
		logger.Infof("数据库 %s 不存在，正在创建...", config.Database.DBName)
		if err := rootDB.Exec("CREATE DATABASE " + config.Database.DBName + " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci").Error; err != nil {
			logger.Fatalf("创建数据库失败: %v", err)
		}
		logger.Infof("数据库 %s 创建成功", config.Database.DBName)
	}

	// 关闭 rootDB 连接
	sqlDB, err := rootDB.DB()
	if err != nil {
		logger.Fatalf("获取底层数据库连接失败: %v", err)
	}
	sqlDB.Close()

	// 连接到指定的数据库
	dsn := config.Database.User + ":" + config.Database.Password + "@tcp(" + config.Database.Host + ":" + strconv.Itoa(config.Database.Port) + ")/" + config.Database.DBName + "?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"

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

	// 检查并初始化数据库
	CheckAndInitializeDatabase()

	// 检查并添加第一个管理员用户
	var adminCount int64
	err = DB.Model(&models.Admin{}).Count(&adminCount).Error
	if err != nil {
		logger.Fatal("检查管理员用户失败:", err)
	}

	if adminCount == 0 {
		// 首先创建超级管理员用户
		superUser := models.User{
			Username: "admin",
			Password: "$2a$10$RIJBMqcxE/qi8Hs8YxjA1.ZFxiJGO6H.YxBXDDEoVxYYqZYxIrIGi", // 默认密码: admin
			Email:    "admin@optioj.com",
		}

		if err := DB.Create(&superUser).Error; err != nil {
			logger.Fatal("添加超级管理员用户失败:", err)
		}
		logger.Info("超级管理员用户已创建")

		// 创建管理员记录
		admin := models.Admin{
			UserID:    superUser.ID,
			Role:      "super_admin",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := DB.Create(&admin).Error; err != nil {
			logger.Fatal("添加管理员记录失败:", err)
		}
		logger.Info("管理员记录已添加")
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
