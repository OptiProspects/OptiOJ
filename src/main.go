package main

import (
	"OptiOJ/src/config"
	"OptiOJ/src/routes"
	"OptiOJ/src/services"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 初始化配置
	config.InitConfig()
	config.InitDB()
	config.InitRedis()

	// 初始化 gRPC 判题客户端
	if err := services.InitJudgeGrpcClient(); err != nil {
		logrus.Fatalf("初始化判题服务失败: %v", err)
	}

	// 创建路由
	r := gin.Default()

	// 配置路由
	routes.SetupRoutes(r)

	// 启动服务器
	logrus.Info("服务器启动，监听端口 2550")
	if err := r.Run(":2550"); err != nil {
		log.Fatal(err)
	}
}
