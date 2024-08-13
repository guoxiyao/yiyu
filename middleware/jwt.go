// middleware.go
package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// JWTMiddleware 用于验证JWT的中间件
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里实现JWT验证逻辑
		// ...
		// 如果验证失败
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
	}
}
