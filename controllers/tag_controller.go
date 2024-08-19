package controllers

import (
	"awesomeProject1/middleware"
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
	//添加jwt校验
	tags := r.Group("/tags", middleware.JwtMiddleware())
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

		//判断表中是否存在
		if ctrl.DB.Where("name = ?", tag.Name).First(&tag).RowsAffected > 0 {
			c.JSON(http.StatusConflict, gin.H{"message": "Tag already exists"})
			return
		}
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

	//查询标签下的日记id
	var diaryTags []models.DiaryTags
	tx := ctrl.DB.Find(&diaryTags, "tag_id = ?", tagID)
	if tx.Error != nil { //删除标签下的日记失败
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete diary under the tag"})
		return
	}

	//根据日记id软删除日记
	var err2 error
	for _, diaryTag := range diaryTags {
		var diary models.Diary
		result2 := ctrl.DB.Where("id = ? AND deleted_at IS NULL", diaryTag.DiaryID).First(&diary)
		if result2.Error != nil {
			err2 = result2.Error
			continue
		}
		if err := ctrl.DB.Delete(&models.Diary{}, diaryTag.DiaryID).Error; err != nil {
			err2 = err
		}
	}
	if err2 != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete diary under the tag"})
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
