package main

import (
	"awesomeProject1/config"
	"awesomeProject1/pkg"
	"awesomeProject1/router"
	"log"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// 连接数据库
	db, err := pkg.Connect(cfg)
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}
	//不加关闭会怎样
	defer db.Close()

	// 初始化路由器
	appRouter := router.NewAppRouter(db)

	// 启动服务器
	log.Printf("Server is running on :8080")
	if err := appRouter.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
