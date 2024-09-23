package middleware

import (
	"awesomeProject1/pkg/jwt"
	"github.com/gin-gonic/gin"
)

// 校验token中间件
func JwtMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		customClaims, err := jwt.VerifyToken(token)
		if err != nil {
			c.JSON(401, gin.H{"message": "JWT verification failed"})
			c.Abort()
		}
		c.Set("user_id", customClaims.UserID)
		c.Set("phone_number", customClaims.PhoneNumber)
		c.Next()
	}
}
