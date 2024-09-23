package controllers

import (
	"awesomeProject1/middleware"
	"awesomeProject1/pkg/models"
	"awesomeProject1/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
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
		//判断表中是否存在
		if ctrl.DB.Where("name = ?", tag.Name).First(&tag).RowsAffected > 0 {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签已存在"))
			return
		}

		result := ctrl.DB.Create(&tag)
		if result.Error != nil {
			response.WriteJSON(c, response.NewResponse(2, nil, "创建标签失败"))
		} else {
			response.WriteJSON(c, response.NewResponse(0, tag, "标签创建成功"))
		}
	} else {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的输入"))
	}
}

// GetTags 获取标签列表
func (ctrl *TagController) GetTags(c *gin.Context) {
	var tags []models.Tag
	result := ctrl.DB.Find(&tags)
	if result.Error != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "获取标签列表失败"))
	} else {
		response.WriteJSON(c, response.NewResponse(0, tags, "获取标签列表成功"))
	}
}

// GetTag 根据ID获取单个标签
func (ctrl *TagController) GetTag(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Param("id"))
	var tag models.Tag
	result := ctrl.DB.Where("id = ?", tagID).First(&tag)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
		} else {
			response.WriteJSON(c, response.NewResponse(2, nil, "获取标签失败"))
		}
		return
	}
	response.WriteJSON(c, response.NewResponse(0, tag, "获取标签成功"))
}

// UpdateTag 更新标签
func (ctrl *TagController) UpdateTag(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Param("id"))
	var tag models.Tag
	result := ctrl.DB.Where("id = ?", tagID).First(&tag)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
			return
		}
		response.WriteJSON(c, response.NewResponse(2, nil, "更新标签失败"))
		return
	}
	if err := c.ShouldBindJSON(&tag); err == nil {
		result = ctrl.DB.Save(&tag)
		if result.Error != nil {
			response.WriteJSON(c, response.NewResponse(2, nil, "更新标签失败"))
		} else {
			response.WriteJSON(c, response.NewResponse(0, tag, "标签更新成功"))
		}
	} else {
		response.WriteJSON(c, response.NewResponse(1, nil, "无效的输入"))
	}
}

// DeleteTag 删除标签（逻辑删除）
func (ctrl *TagController) DeleteTag(c *gin.Context) {
	tagID, _ := strconv.Atoi(c.Param("id"))
	var tag models.Tag
	result := ctrl.DB.Where("id = ?", tagID).First(&tag)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到"))
			return
		}
		response.WriteJSON(c, response.NewResponse(2, nil, "删除标签时数据库错误"))
		return
	}

	// 执行软删除
	tag.DeletedAt.Time = time.Now() // 设置 DeletedAt 字段为当前时间
	tag.DeletedAt.Valid = true      // 标记为有效，表示记录被软删除
	result = ctrl.DB.Save(&tag)
	if result.Error != nil {
		response.WriteJSON(c, response.NewResponse(2, nil, "删除标签失败"))
	} else if result.RowsAffected == 0 {
		response.WriteJSON(c, response.NewResponse(1, nil, "标签未找到或已被删除"))
	} else {
		response.WriteJSON(c, response.NewResponse(0, nil, "标签删除成功"))
	}
}
