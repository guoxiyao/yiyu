package controllers

import (
	"awesomeProject1/models"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type DiaryController struct {
	DB *gorm.DB
}

func NewDiaryController(db *gorm.DB) *DiaryController {
	return &DiaryController{DB: db}
}

// CreateDiary 创建新的日记
func (ctrl *DiaryController) CreateDiary(c *gin.Context) {
	var diary models.Diary
	if err := c.ShouldBindJSON(&diary); err == nil {
		result := ctrl.DB.Create(&diary)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		} else {
			c.JSON(http.StatusCreated, gin.H{"message": "Diary created", "diary": diary})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
}

// ListDiaries 列出所有日记
func (ctrl *DiaryController) ListDiaries(c *gin.Context) {
	var diaries []models.Diary
	result := ctrl.DB.Find(&diaries)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"diaries": diaries})
	}
}

// GetDiaryByID 根据ID获取单个日记
func (ctrl *DiaryController) GetDiaryByID(c *gin.Context) {
	diaryID, err := strconv.Atoi(c.Param("id"))
	if err != nil || diaryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
		return
	}
	var diary models.Diary
	result := ctrl.DB.Where("id = ?", diaryID).First(&diary)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"diary": diary})
}

// UpdateDiary 更新日记内容
// UpdateDiary 更新日记内容
func (ctrl *DiaryController) UpdateDiary(c *gin.Context) {
	diaryID, err := strconv.Atoi(c.Param("id"))
	if err != nil || diaryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
		return
	}

	// 从请求中获取更新的内容
	var updateContent struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateContent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content", "details": err.Error()})
		return
	}

	// 使用 GORM 更新日记内容
	result := ctrl.DB.Model(&models.Diary{}).Where("id = ?", diaryID).Updates(map[string]interface{}{"content": updateContent.Content})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Diary updated"})
	}
}

// DeleteDiary 删除日记（软删除）
func (ctrl *DiaryController) DeleteDiary(c *gin.Context) {
	diaryID, err := strconv.Atoi(c.Param("id"))
	if err != nil || diaryID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid diary ID"})
		return
	}
	result := ctrl.DB.Where("id = ?", diaryID).Delete(&models.Diary{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Diary deleted"})
	}
}
