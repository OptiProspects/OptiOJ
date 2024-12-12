package main

import (
	"OptiOJ/src/config"
	"OptiOJ/src/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	config.InitConfig()

	// 设置 logrus 日志格式
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	config.InitDB()
	config.InitRedis()

	r := gin.Default()

	// 配置 CORS 规则
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},                                       // 允许所有来源
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},            // 允许的方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
		ExposeHeaders:    []string{"Content-Length"},                          // 允许暴露的响应头
		AllowCredentials: true,                                                // 允许携带凭证
	}))

	routes.SetupRoutes(r)

	logrus.Info("服务器启动，监听端口 8080")
	if err := r.Run(":8080"); err != nil {
		logrus.Fatal("服务器启动失败: ", err)
	}
}
