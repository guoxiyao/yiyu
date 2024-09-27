package response

type Page struct {
	Total       uint          `json:"total"`       // 总数据量
	CurrentPage uint          `json:"currentPage"` // 当前页
	TotalPages  uint          `json:"totalPages"`  // 总页数
	PageSize    uint          `json:"pageSize"`    // 每页大小
	Records     []interface{} `json:"records"`     // 数据列表
}

// TODO: 分页组件
