package main

import (
	"awesomeProject1/config"
	"awesomeProject1/pkg"
	"awesomeProject1/router"
	"log"
)

func main() {
	// 加载配置
	cfg := config.NewDatabaseConfig()

	// 连接数据库
	_, err := pkg.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// 初始化路由器
	r := router.SetupRouter()

	// 启动服务器
	log.Printf("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
