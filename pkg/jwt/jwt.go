package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// CustomClaims 自定义的 JWT 声明结构体
type CustomClaims struct {
	UserID      uint   `json:"user_id"`
	PhoneNumber string `json:"phone_number"`
	jwt.StandardClaims
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(userID uint, phoneNumber string) (string, error) {
	key := []byte("your_jwt_secret_key")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID:      userID,
		PhoneNumber: phoneNumber,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	return token.SignedString(key)
}

// VerifyToken 校验 JWT 令牌
func VerifyToken(tokenString string) (*CustomClaims, error) {
	key := []byte("your_jwt_secret_key")

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
