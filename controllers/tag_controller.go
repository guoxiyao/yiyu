package controllers

import (
	"awesomeProject1/pkg/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// TagController 处理标签相关的请求
type TagController struct {
	DB *gorm.DB
}

// NewTagController 创建标签控制器实例
func NewTagController(db *gorm.DB) *TagController {
	return &TagController{DB: db}
}

// Routes 注册标签相关的路由
func (ctrl *TagController) Routes(r *gin.Engine) {
	tags := r.Group("/tags")
	{
		tags.POST("", ctrl.CreateTag)       // 创建标签
		tags.GET("", ctrl.GetTags)          // 获取标签列表
		tags.GET("/:id", ctrl.GetTag)       // 根据ID获取单个标签
		tags.PUT("/:id", ctrl.UpdateTag)    // 更新标签
		tags.DELETE("/:id", ctrl.DeleteTag) // 删除标签
	}
}

// CreateTag 创建新的标签
func (ctrl *TagController) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err == nil {
		//if err := tag.Validate(); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}
		result := ctrl.DB.Create(&tag)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "Tag created", "tag": tag})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}
}

// GetTags 获取标签列表
func (ctrl *TagController) GetTags(c *gin.Context) {
	var tags []models.Tag
	result := ctrl.DB.Find(&tags)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"tags": tags})
	}
}

// GetTag 根据ID获取单个标签
func (ctrl *TagController) GetTag(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Param("id"))
	var tag models.Tag
	result := ctrl.DB.Where("id = ?", tagID).First(&tag)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag": tag})
}

// UpdateTag 更新标签
func (ctrl *TagController) UpdateTag(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Param("id"))
	var tag models.Tag
	result := ctrl.DB.Where("id = ?", tagID).First(&tag)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if err := c.ShouldBindJSON(&tag); err == nil {
		//if err := tag.Validate(); err != nil {
		//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//	return
		//}
		result = ctrl.DB.Save(&tag)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Tag updated", "tag": tag})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
	}
}

// DeleteTag 删除标签（逻辑删除）
func (ctrl *TagController) DeleteTag(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Param("id"))
	var tag models.Tag
	result := ctrl.DB.Where("id = ?", tagID).First(&tag)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	result = ctrl.DB.Delete(&tag)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found or already deleted"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Tag deleted"})
	}
}
