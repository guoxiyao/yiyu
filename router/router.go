// router.go
package router

import (
	"awesomeProject1/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AppRouter 结构体，包含 Gin 路由器和数据库连接
type AppRouter struct {
	DB     *gorm.DB
	Engine *gin.Engine
}

// NewAppRouter 初始化并返回 AppRouter 实例
func NewAppRouter(db *gorm.DB) *AppRouter {
	r := gin.Default() // 创建 Gin 路由器实例

	// 初始化中间件
	// r.Use(middleware.Logger())

	// 初始化控制器并设置路由
	controllers.NewDiaryController(db).Routes(r)
	//controllers.NewUserController(db).Routes(r)
	controllers.NewTagController(db).Routes(r)
	// ... 其他控制器的注册 ...

	user := r.Group("/diary")
	{
		user.POST("/login", controllers.NewUserController(db).Login)
	}

	return &AppRouter{
		DB:     db,
		Engine: r,
	}
}

// Run 启动 Gin 服务器
func (ar *AppRouter) Run(port string) error {
	return ar.Engine.Run(port) // 使用 Gin 路由器的 Run 方法
}
