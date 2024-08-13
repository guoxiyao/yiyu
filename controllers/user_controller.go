package controllers

import (
	"awesomeProject1/pkg/jwt"
	"awesomeProject1/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// UserController 用户控制器
type UserController struct {
	DB *gorm.DB
}

// NewUserController 创建用户控制器实例
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	var newUser models.User
	if err := c.ShouldBindJSON(&newUser); err == nil {
		//if err := newUser.Validate(ctrl.DB); err != nil {
		//	c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
		//	return
		//}
		// 这里应该添加密码加密逻辑
		newUser.Password = "encrypted_password" // 用实际的加密密码替换
		result := ctrl.DB.Create(&newUser)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var loginUser models.User
	if err := c.ShouldBindJSON(&loginUser); err == nil {
		var user models.User
		result := ctrl.DB.Where("phone_number = ? AND is_deleted = false", loginUser.PhoneNumber).First(&user)
		if result.Error == nil && user.CheckPassword(loginUser.Password) {
			// 登录成功，生成JWT令牌
			token, err := jwt.GenerateToken(user.ID, user.PhoneNumber) // 确保您接收返回的令牌和错误
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			} else {
				// 确保您发送了生成的令牌
				c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
			}
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or password"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}
}
