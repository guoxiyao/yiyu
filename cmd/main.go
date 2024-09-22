package main

import (
	"awesomeProject1/config"
	"awesomeProject1/pkg"
	"awesomeProject1/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	//实例化gin引擎
	r := gin.New()

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

	// 初始化路由器
	appRouter := router.NewAppRouter(db)

	// 配置 CORS
	config := cors.Config{
		AllowOrigins:  []string{"*"}, // 允许所有来源
		AllowMethods:  []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders:  []string{"Origin"},
		ExposeHeaders: []string{"Content-Length", "Authorization"}, // 暴露的头部
	}

	// 应用 CORS 中间件
	r.Use(cors.New(config))

	// 启动服务器
	log.Printf("Server is running on :8080")
	if err := appRouter.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
