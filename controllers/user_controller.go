package controllers

import (
	"awesomeProject1/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	_ "strconv"
	_ "time"
)

type UserController struct {
	DB *gorm.DB
}

func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

// Register 用户注册
func (ctrl *UserController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err == nil {
		// 验证用户数据
		if err := user.Validate(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 检查手机号是否已存在
		var existingUser models.User
		result := ctrl.DB.Where("phone_number = ?", user.PhoneNumber).First(&existingUser)
		if result.Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Phone number already registered"})
			return
		}

		// 创建用户
		result = ctrl.DB.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			// 注册成功后，可以发送确认邮件或其他操作
			c.JSON(http.StatusCreated, gin.H{"message": "User registered", "user": user})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}
}

// Login 用户登录
func (ctrl *UserController) Login(c *gin.Context) {
	var loginInfo struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
		Password    string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&loginInfo); err == nil {
		var user models.User
		result := ctrl.DB.Where("phone_number = ? AND password = ?", loginInfo.PhoneNumber, loginInfo.Password).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "User not found or incorrect password"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
			}
		} else {
			// 登录成功后，可以生成JWT令牌或其他操作
			// token := generateJWTForUser(&user)
			// c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
			c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login information"})
	}
}

// generateJWTForUser is a placeholder function for generating a JWT token for a user.
// You would need to implement the actual JWT generation logic.

// ... 其他用户相关控制器方法 ...
