// file: diary_controller.go
package controllers

import (
	"awesomeProject1/middleware"
	"awesomeProject1/pkg/models"
	"errors"
	"strconv"
	"time"

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

// Routes 注册日记相关的路由
func (ctrl *DiaryController) Routes(r *gin.Engine) {
	//添加jwt校验
	diaries := r.Group("/diaries", middleware.JwtMiddleware())
	{
		diaries.POST("", ctrl.CreateDiary)       // 创建日记
		diaries.GET("", ctrl.GetDiaries)         // 获取日记列表
		diaries.GET("/:id", ctrl.GetDiary)       // 根据 ID 获取日记
		diaries.PUT("/:id", ctrl.UpdateDiary)    // 更新日记
		diaries.DELETE("/:id", ctrl.DeleteDiary) // 删除日记
		//diaries.POST("/login", ctrl.Login)       // 用户登录
	}
}

// CreateDiary 创建新的日记
func (ctrl *DiaryController) CreateDiary(c *gin.Context) {
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)
	var diary models.Diary

	// 绑定请求中的 JSON 数据到 diary 变量中
	if err := c.ShouldBindJSON(&diary); err != nil {
		// 如果绑定 JSON 数据失败，返回错误信息和 HTTP 400 状态码
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// 检查日记内容是否为空
	if diary.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Diary content cannot be empty"})
		return
	}

	// 检查 user_id 是否在 users 表中存在
	var user models.User
	if err := ctrl.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User not found"})
		return
	}
	diary.UserID = userID

	// 查找标签,未创建的接口返回错误
	var tags []models.Tag
	for _, tag := range diary.Tags {
		var tag1 models.Tag
		if err := ctrl.DB.First(&tag1, "name = ?", tag.Name).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tag not found"})
			return
		}
		tags = append(tags, tag1)
	}

	diary.Tags = tags

	// 使用 GORM 创建日记记录
	result := ctrl.DB.Create(&diary)

	// 检查是否有错误发生
	if result.Error != nil {
		// 如果有错误，返回错误信息和 HTTP 500 状态码
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 如果创建成功，返回创建的日记和 HTTP 201 状态码
	c.JSON(http.StatusCreated, gin.H{"diary": diary})
}

// GetDiaries 获取用户的日记列表
func (ctrl *DiaryController) GetDiaries(c *gin.Context) {
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)
	pageStr, ok := c.GetQuery("page")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	pageSizeStr, ok := c.GetQuery("pageSize")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
		return
	}
	sortField, ok := c.GetQuery("sortField")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort field"})
		return
	}
	if sortField != "created_at" && sortField != "updated_at" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort field"})
	}
	sortBy, ok := c.GetQuery("sortBy")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort order"})
		return
	}
	if sortBy != "ASC" && sortBy != "DESC" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sort order"})
	}
	queryType, ok := c.GetQuery("queryType")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query type"})
		return
	}
	var result *gorm.DB
	if queryType == "tag" {
		//按标签查询
		tag, ok := c.GetQuery("tag")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid tag"})
		}
		//查询标签的id
		var tagData models.Tag
		tx := ctrl.DB.Where("name = ?", tag).First(&tagData)
		if tx.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tag not found"})
			return
		}
		//通过id进行筛选
		result = ctrl.DB.Joins("JOIN diary_tags ON diary_tags.diary_id = diaries.id").
			Where("diary_tags.tag_id = ?", tagData.ID)
	} else if queryType == "content" {
		//按内容模糊查询
		content, ok := c.GetQuery("content")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid content"})
		}
		result = ctrl.DB.Where("content LIKE ?", "%"+content+"%")
	} else if queryType == "time" {
		//按创建时间范围查询
		stratTime, ok := c.GetQuery("stratTime")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stratTime"})
		}
		endTime, ok := c.GetQuery("endTime")
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endTime"})
		}

		startTimeNew, err := time.Parse("2006-01-02/15:04:05", stratTime)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stratTime"})
			return
		}
		endTimeNew, err2 := time.Parse("2006-01-02/15:04:05", endTime)
		if err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid endTime"})
			return
		}
		result = ctrl.DB.Where("created_at BETWEEN ? AND ?", startTimeNew.Format("2006-01-02 15:04:05"), endTimeNew.Format("2006-01-02 15:04:05"))
	} else {
		//其他类型返回错误
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query type"})
		return
	}

	//校验参数
	if page < 1 || pageSize < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number or page size"})
	}
	offset := (page - 1) * pageSize
	// 实现获取日记列表逻辑
	var diaries []models.Diary // 假设 models.Diary 是您的日记模型
	result = result.Where("user_id = ?", userID).Preload("Tags").Offset(offset).Limit(pageSize).Order(sortField + " " + sortBy).Find(&diaries)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"diaries": diaries})
	}
}

// UpdateDiary 更新日记内容
func (ctrl *DiaryController) UpdateDiary(c *gin.Context) {
	diaryID := c.Param("id")
	if diaryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing diary ID"})
		return
	}
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	var update models.Diary
	if err := c.BindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid update data"})
		return
	}

	var diary models.Diary
	result := ctrl.DB.First(&diary, "id = ? AND user_id = ?", diaryID, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	//先判断标签是否存在
	var tags []models.Tag
	for _, tag := range update.Tags {
		var tag1 models.Tag
		if err := ctrl.DB.First(&tag1, "name = ?", tag.Name).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tag not found"})
			return
		}
		tags = append(tags, tag1)
	}
	diary.Tags = tags

	// 更新字段...
	diary.Content = update.Content // 假设只更新内容字段
	diary.UpdatedAt = time.Now()   // 更新更新时间

	result = ctrl.DB.Save(&diary)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	//更新标签
	err := ctrl.DB.Model(&diary).Association("Tags").Replace(&tags)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update label"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"diary": diary})
}

// GetDiary 根据ID获取单个日记
func (ctrl *DiaryController) GetDiary(c *gin.Context) {
	diaryID := c.Param("id") // 从URL参数中获取日记ID
	if diaryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing diary ID"})
		return
	}
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	var diary models.Diary
	result := ctrl.DB.Where("user_id = ?", userID).Preload("Tags").Take(&diary, diaryID)
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

// DeleteDiary 软删除日记（标记为已删除）
func (ctrl *DiaryController) DeleteDiary(c *gin.Context) {
	diaryID := c.Param("id")
	if diaryID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing diary ID"})
		return
	}
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	// 尝试获取日记记录
	var diary models.Diary
	result := ctrl.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", diaryID, userID).First(&diary)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found or already deleted"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 执行软删除
	updateResult := ctrl.DB.Model(&diary).Update("deleted_at", time.Now())
	if updateResult.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": updateResult.Error.Error()})
	} else if updateResult.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Diary not found or already deleted"})
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Diary deleted"})
	}
}
