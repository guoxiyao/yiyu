// Package controllers Package controllers Package controllers Package controllers controllers/tag_controller.go
package controllers

import (
	"awesomeProject1/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// TagController 结构体，用于包含数据库连接
type TagController struct {
	DB *gorm.DB
}

// NewTagController 构造函数，初始化 TagController 实例

// GetTagByID 根据ID获取单个标签
func (ctrl *TagController) GetTagByID(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil || tagID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}
	var tag models.Tag
	result := ctrl.DB.Find(&tag, tagID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"tag": tag})
}

// UpdateTag 更新标签名称
func (ctrl *TagController) UpdateTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil || tagID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}
	var updatedTag models.Tag
	if err := c.ShouldBindJSON(&updatedTag); err == nil {
		result := ctrl.DB.Model(&models.Tag{}).Where("id = ?", tagID).Updates(&updatedTag)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
		} else {
			c.JSON(http.StatusOK, gin.H{"message": "Tag updated"})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// DeleteTag 删除标签
func (ctrl *TagController) DeleteTag(c *gin.Context) {
	tagIDStr := c.Param("id")
	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil || tagID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}
	result := ctrl.DB.Delete(&models.Tag{}, tagID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Tag deleted"})
	}
}

// ListTags 列出所有标签
func (ctrl *TagController) ListTags(c *gin.Context) {
	var tags []models.Tag
	result := ctrl.DB.Find(&tags)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// CreateTag 创建新的标签
func (ctrl *TagController) CreateTag(c *gin.Context) {
	var tag models.Tag
	if err := c.ShouldBindJSON(&tag); err == nil {
		result := ctrl.DB.Create(&tag)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "Tag created", "tag": tag})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}
