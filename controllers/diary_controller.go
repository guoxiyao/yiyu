package controllers

import (
	"awesomeProject1/middleware"
	"awesomeProject1/pkg/models"
	"awesomeProject1/request"
	"awesomeProject1/response"
	"awesomeProject1/service"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
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
	var diaryDto request.DiaryDto

	// 绑定请求中的 JSON 数据到 diary 变量中
	if err := c.ShouldBindJSON(&diaryDto); err != nil {
		response.WriteJSON(c, response.UserErrorNoMsgResponse("无效的输入"))
		//response.WriteBadRequestJSON(c, "无效的输入")
		return
	}

	// 检查日记内容是否为空
	if diaryDto.Content == "" {
		response.WriteJSON(c, response.NewResponse(1, nil, "日记内容不能为空"))
		return
	}
	diary.Content = diaryDto.Content

	// 检查 user_id 是否在 users 表中存在
	var user models.User
	if err := ctrl.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		response.WriteJSON(c, response.NewResponse(1, nil, "用户未找到"))
		return
	}
	diary.UserID = userID

	// 查找标签,未创建的接口返回错误
	var tags []models.Tag
	for _, tagId := range diaryDto.TagIds {
		var tag models.Tag
		if err := ctrl.DB.First(&tag, "id = ?", tagId).Error; err != nil {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
			return
		}
		tags = append(tags, tag)
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
	diary.User = user

	// 如果创建成功，返回创建的日记和 HTTP 201 状态码
	response.WriteJSON(c, response.NewResponse(0, diary, "日记创建成功"))
}

// GetDiaries 获取用户的日记列表
func (ctrl *DiaryController) GetDiaries(c *gin.Context) {
	// 绑定查询参数
	var queryParams struct {
		Page      int    `form:"page" binding:"required,gte=1"`
		PageSize  int    `form:"pageSize" binding:"required,gte=1"`
		TagID     string `form:"tagId"`
		Content   string `form:"content"`
		StartTime string `form:"startTime"`
		EndTime   string `form:"endTime"`
		SortField string `form:"sortField"`
		SortBy    string `form:"sortBy"`
	}
	if err := c.ShouldBindQuery(&queryParams); err != nil {
		response.WriteJSON(c, response.NewResponse(http.StatusBadRequest, nil, "无效的查询参数"))
		return
	}

	// 调用分页服务
	paginationResult, err := service.Paginate(ctrl.DB, queryParams.Page, queryParams.PageSize, &models.Diary{}, queryParams.SortField, queryParams.SortBy)
	if err != nil {
		response.WriteJSON(c, response.NewResponse(http.StatusInternalServerError, nil, "获取日记列表失败"))
		return
	}

	// 转换日记数据为DiaryVo
	diaryVos := make([]response.DiaryVo, len(paginationResult.Records.([]models.Diary)))
	for i, diary := range paginationResult.Records.([]models.Diary) {
		diaryVos[i].Copy(diary)
	}

	// 调用响应服务发送分页响应
	response.PageResponse(c, paginationResult, func(model interface{}) []response.DiaryVo {
		records, _ := model.([]models.Diary)
		diaryVos := make([]response.DiaryVo, len(records))
		for i, diary := range records {
			diaryVos[i].Copy(diary)
		}
		return diaryVos
	})
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
