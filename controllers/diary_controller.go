package controllers

import (
	"awesomeProject1/middleware"
	"awesomeProject1/pkg/models"
	"awesomeProject1/response"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"math"
	"time"
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
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的输入"))
		return
	}

	// 检查日记内容是否为空
	if diary.Content == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "日记内容不能为空"))
		return
	}

	// 检查 user_id 是否在 users 表中存在
	var user models.User
	if err := ctrl.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "用户未找到"))
		return
	}
	diary.UserID = userID

	// 查找标签,未创建的接口返回错误
	var tags []models.Tag
	for _, tag := range diary.Tags {
		var tag1 models.Tag
		if err := ctrl.DB.First(&tag1, "name = ?", tag.Name).Error; err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
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

		response.WriteJSON(c, response.NewResponse(2, nil, "创建日记失败"))
		return
	}

	// 如果创建成功，返回创建的日记和 HTTP 201 状态码
	response.WriteJSON(c, response.NewResponse(0, diary, "日记创建成功"))
}

// GetDiaries 获取用户的日记列表
func (ctrl *DiaryController) GetDiaries(c *gin.Context) {
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)
	// QueryParams 用于绑定查询参数的结构体
	type QueryParams struct {
		Page      int
		PageSize  int
		SortField string
		SortBy    string
		QueryType string
		Tag       string
		Content   string
		StratTime string
		EndTime   string
	}

	// 绑定查询参数
	var queryParams QueryParams
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的查询参数"))
		return
	}

	//验证分页参数
	if queryParams.Page < 1 || queryParams.PageSize < 1 {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的页码或页面大小"))
		return
	}
	// 验证排序参数
	if queryParams.SortField != "created_at" && queryParams.SortField != "updated_at" {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的排序字段"))
		return
	}
	if queryParams.SortBy != "ASC" && queryParams.SortBy != "DESC" {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的排序顺序"))
		return
	}
	//构建查询
	var result *gorm.DB
	switch queryParams.QueryType {
	case "tag":
		// 按标签查询
		var tagData models.Tag
		if err := ctrl.DB.Where("name = ?", queryParams.Tag).First(&tagData).Error; err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
			return
		}
		// 通过id进行筛选
		result = ctrl.DB.Joins("JOIN diary_tags ON diary_tags.diary_id = diaries.id").
			Where("diary_tags.tag_id = ?", tagData.ID)

	case "content":
		// 按内容模糊查询
		result = ctrl.DB.Where("content LIKE ?", "%"+queryParams.Content+"%")

	case "time":
		// 按创建时间范围查询
		startTime, err := time.Parse("2006-01-02/15:04:05", queryParams.StratTime)
		if err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "无效的开始时间"))
			return
		}
		endTime, err := time.Parse("2006-01-02/15:04:05", queryParams.EndTime)
		if err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "无效的结束时间"))
			return
		}
		result = ctrl.DB.Where("created_at BETWEEN ? AND ?", startTime, endTime)

	default:
		// 其他类型返回错误
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的查询类型"))
		return
	}

	// 应用分页和排序
	offset := (queryParams.Page - 1) * queryParams.PageSize
	var diaries []models.Diary
	var count int64
	result = ctrl.DB.Model(&models.Diary{}).Where("user_id = ?", userID).Count(&count)
	if result.Error != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "获取日记总数失败"))
		return
	}
	result = ctrl.DB.Where("user_id = ?", userID).Preload("Tags").Offset(offset).Limit(queryParams.PageSize).Order(queryParams.SortField + " " + queryParams.SortBy).Find(&diaries)
	if result.Error != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "获取日记列表失败"))
	} else {
		// 计算总页数
		totalPages := int(math.Ceil(float64(count) / float64(queryParams.PageSize)))
		response.WriteJSON(c, response.NewResponse(0, gin.H{"diaries": diaries, "currentPage": queryParams.Page, "totalPages": totalPages}, "获取日记列表成功"))
	}
}

// UpdateDiary 更新日记内容
func (ctrl *DiaryController) UpdateDiary(c *gin.Context) {
	diaryID := c.Param("id")
	if diaryID == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "缺少日记ID"))
		return
	}
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	var update models.Diary
	if err := c.BindJSON(&update); err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的更新数据"))
		return
	}

	var diary models.Diary
	result := ctrl.DB.First(&diary, "id = ? AND user_id = ?", diaryID, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.WriteJSON(c, response.NewResponse(1, nil, "日记未找到"))
			return
		}
		response.WriteJSON(c, response.NewResponse(2, nil, "更新日记时数据库错误"))
		return
	}

	//先判断标签是否存在
	var tags []models.Tag
	for _, tag := range update.Tags {
		var tag1 models.Tag
		if err := ctrl.DB.First(&tag1, "name = ?", tag.Name).Error; err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
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
		response.WriteJSON(c, response.NewResponse(2, nil, "更新日记失败"))
		return
	}
	//更新标签
	err := ctrl.DB.Model(&diary).Association("Tags").Replace(&tags)
	if err != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "更新标签失败"))
		return
	}
	response.WriteJSON(c, response.NewResponse(0, diary, "日记更新成功"))
}

// GetDiary 根据ID获取单个日记
func (ctrl *DiaryController) GetDiary(c *gin.Context) {
	diaryID := c.Param("id") // 从URL参数中获取日记ID
	if diaryID == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "缺少日记ID"))
		return
	}
	//用户id
	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	var diary models.Diary
	result := ctrl.DB.Where("user_id = ?", userID).Preload("Tags").Take(&diary, diaryID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.WriteJSON(c, response.NewResponse(1, nil, "日记未找到"))
		} else {
			response.WriteJSON(c, response.NewResponse(2, nil, "获取日记时数据库错误"))
		}
		return
	}
	response.WriteJSON(c, response.NewResponse(0, diary, "获取日记成功"))
}

// DeleteDiary 软删除日记（标记为已删除）
func (ctrl *DiaryController) DeleteDiary(c *gin.Context) {
	diaryID := c.Param("id")
	if diaryID == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "缺少日记ID"))
		return
	}

	userIDAny, _ := c.Get("user_id")
	userID := userIDAny.(uint)

	// 尝试获取日记记录
	var diary models.Diary
	result := ctrl.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", diaryID, userID).First(&diary)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			response.WriteJSON(c, response.NewResponse(1, nil, "日记未找到或已被删除"))
			return
		}
		response.WriteJSON(c, response.NewResponse(2, nil, "删除日记时数据库错误"))
		return
	}

	// 执行软删除
	updateResult := ctrl.DB.Model(&diary).Update("deleted_at", time.Now())
	if updateResult.Error != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "删除日记失败"))
	} else if updateResult.RowsAffected == 0 {
		response.WriteJSON(c, response.NewResponse(1, nil, "日记未找到或已被删除"))
	} else {
		response.WriteJSON(c, response.NewResponse(0, nil, "日记删除成功"))
	}
}
