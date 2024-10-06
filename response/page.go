package response

import (
	"awesomeProject1/service"
	"github.com/gin-gonic/gin"
)

type Page struct {
	Total       uint      `json:"total"`       // 总数据量
	CurrentPage uint      `json:"currentPage"` // 当前页
	TotalPages  uint      `json:"totalPages"`  // 总页数
	PageSize    uint      `json:"pageSize"`    // 每页大小
	Records     []DiaryVo `json:"records"`     // 数据列表
}

// PageResponse 创建一个分页响应
func PageResponse(c *gin.Context, page *service.PaginationResult, recordVoConverter func(model interface{}) []DiaryVo) {
	// 将 model 转换为 DiaryVo
	diaryVos := recordVoConverter(page.Records)

	// 创建 Page 实例
	pageData := Page{
		Total:       uint(page.Total),
		CurrentPage: uint(page.Page),
		TotalPages:  uint(page.TotalPages),
		PageSize:    uint(page.PageSize),
		Records:     diaryVos,
	}

	// 使用 ApiResponse 封装分页响应
	apiResponse := ApiResponse{
		Code:    200,
		Data:    pageData,
		Message: "success",
	}

	// 发送响应
	WriteJSON(c, apiResponse)
}
