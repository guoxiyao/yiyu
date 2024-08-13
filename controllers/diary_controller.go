// file: diary_controller.go
package controllers

import (
	"awesomeProject1/pkg/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// DiaryController 处理日记相关的请求
type DiaryController struct {
	DB *gorm.DB
}

// NewDiaryController 创建日记控制器实例
func NewDiaryController(db *gorm.DB) *DiaryController {
	return &DiaryController{DB: db}
}

/*
// Routes 注册日记相关的路由
func (ctrl *DiaryController) Routes(r *gin.Engine) {
	diaries := r.Group("/diaries")
	{
		diaries.POST("", ctrl.CreateDiary)           // 创建日记
		diaries.GET("", ctrl.GetDiaries)             // 获取日记列表
		diaries.GET("/:id", ctrl.GetDiary)           // 根据 ID 获取日记
		diaries.PUT("/:id", ctrl.UpdateDiary)        // 更新日记
		diaries.DELETE("/:id", ctrl.DeleteDiary)     // 删除日记
	}
}
*/
// CreateDiary 创建新的日记
func (ctrl *DiaryController) CreateDiary(c *gin.Context) {
	// 实现创建日记逻辑
}

// GetDiaries 获取用户的日记列表
func (ctrl *DiaryController) GetDiaries(c *gin.Context) {
	// 实现获取日记列表逻辑
	var diaries []models.Diary // 假设 models.Diary 是您的日记模型
	result := ctrl.DB.Find(&diaries)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"diaries": diaries})
	}
}

func (ctrl *DiaryController) Routes(r *gin.Engine) {

}
