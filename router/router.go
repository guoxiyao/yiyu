package router

import (
	"awesomeProject1/config"
	"awesomeProject1/controllers"
	"awesomeProject1/pkg"
	"github.com/gin-gonic/gin"
	_ "gorm.io/gorm"
	"log"
	"net/http"
	_ "net/http"
)

// SetupRouter 初始化路由器和路由
func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 配置路由组
	userGroup := r.Group("/users")
	diaryGroup := r.Group("/diaries")

	// 连接到数据库
	db, err := pkg.Connect(config.NewDatabaseConfig())
	if err != nil {
		// 记录日志并处理错误
		log.Printf("Failed to connect to the database: %v", err)
		r.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		})
		return r
	}
	// 确保在函数返回前关闭数据库连接
	defer func() {
		sqlDB, err := db.DB()
		if err != nil {
			log.Printf("Error getting SQL DB: %v", err)
			return
		}

		// 关闭数据库连接
		err = sqlDB.Close()
		if err != nil {
			log.Printf("Error closing the database connection: %v", err)
		}
	}()

	// 初始化控制器
	userController := controllers.NewUserController(db)
	diaryController := controllers.NewDiaryController(db)

	// 设置用户相关路由
	userGroup.POST("/register", userController.Register)
	userGroup.POST("/login", userController.Login)

	// 设置日记相关路由
	diaryGroup.POST("", diaryController.CreateDiary)
	diaryGroup.GET("", diaryController.ListDiaries)

	// 使用日志中间件
	r.Use(Logger())

	// 启动服务器
	log.Printf("Server is running on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

	return r
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在这里添加日志逻辑
		c.Next() // 调用后续的处理函数
	}
}
