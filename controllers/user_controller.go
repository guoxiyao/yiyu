package controllers

import (
	"awesomeProject1/pkg/jwt"
	"awesomeProject1/pkg/models"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

// UserController 用户控制器
type UserController struct {
	DB *gorm.DB
}

// Login 用户注册或登录
func (ctrl *UserController) Login(c *gin.Context) {
	var credentials models.User
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "0", "msg": "Invalid input"})
		return
	}

	// 检查用户是否存在
	var user models.User
	result := ctrl.DB.Where("phone_number = ?", credentials.PhoneNumber).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 用户不存在，执行注册逻辑
			hashPassword, err := bcrypt.GenerateFromPassword([]byte(credentials.Password), bcrypt.DefaultCost)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "0", "msg": "Failed to encrypt password"})
				return
			}
			user.PhoneNumber = credentials.PhoneNumber
			user.Password = string(hashPassword)
			// 保存新用户
			createResult := ctrl.DB.Create(&user)
			if createResult.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"status": "0", "msg": "Failed to register user"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"status": "0", "msg": "Database error"})
			return
		}
	} else {
		// 用户存在，执行登录逻辑
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "0", "msg": "Invalid phone number or password"})
			return
		}
	}

	// 登录成功，生成JWT令牌
	token, err := jwt.GenerateToken(user.ID, user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "0", "msg": "Failed to generate token"})
	} else {
		c.JSON(http.StatusOK, gin.H{"status": "1", "msg": "注册或登录成功！", "token": token})
	}
}

// 确保在 UserController 的 Routes 方法中注册 Login 路由
func (ctrl *UserController) Routes(r *gin.Engine) {
	r.POST("/diary/login", ctrl.Login)
}

// NewUserController 创建 UserController 的新实例
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}
