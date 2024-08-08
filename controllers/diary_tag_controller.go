// Package controllers Package controllers Package controllers controllers/diary_tag_controller.go
package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"awesomeProject1/models" // 确保这个路径是正确的
	"gorm.io/gorm"
)

type DiaryTagController struct {
	DB *gorm.DB
}

// AttachTagToDiary 用于将标签附加到日记
func (ctrl *DiaryTagController) AttachTagToDiary(c *gin.Context) {
	diaryIDStr := c.Param("diary_id")
	tagIDStr := c.Param("tag_id")

	diaryID, err := strconv.Atoi(diaryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
		return
	}

	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	// 创建 DiaryTag 关联对象
	diaryTag := models.DiaryTag{
		DiaryID: uint(diaryID),
		TagID:   uint(tagID),
	}

	// 检查日记和标签是否已经关联
	if result := ctrl.DB.Where("diary_id = ? AND tag_id = ?", diaryID, tagID).First(&models.DiaryTag{}); result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Tag is already attached to the diary"})
		return
	}

	// 关联日记和标签
	if result := ctrl.DB.Create(&diaryTag); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag attached to diary"})
}

// DetachTagFromDiary 用于从日记中移除标签
func (ctrl *DiaryTagController) DetachTagFromDiary(c *gin.Context) {
	diaryIDStr := c.Param("diary_id")
	tagIDStr := c.Param("tag_id")

	diaryID, err := strconv.Atoi(diaryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
		return
	}

	tagID, err := strconv.Atoi(tagIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag ID"})
		return
	}

	// 删除 DiaryTag 关联对象
	result := ctrl.DB.Where("diary_id = ? AND tag_id = ?", diaryID, tagID).Delete(&models.DiaryTag{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tag is not attached to the diary"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tag detached from diary"})
}
