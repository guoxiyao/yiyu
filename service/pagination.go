package service

import (
	"gorm.io/gorm"
	"math"
)

type PaginationResult struct {
	Records    interface{} // 记录数据
	Total      int64       // 总记录数
	Page       int         // 当前页码
	PageSize   int         // 每页记录数
	TotalPages int         // 总页数
}

func Paginate(db *gorm.DB, page, pageSize int, model interface{}, sortField string, sortBy string) (*PaginationResult, error) {
	var count int64
	var err error

	// 构建基础查询
	query := db.Model(model)

	// 应用排序条件
	if sortField != "" && sortBy != "" {
		query = query.Order(sortField + " " + sortBy)
	}

	err = query.Count(&count).Error
	if err != nil {
		return nil, err
	}

	// 计算总页数和偏移量
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))
	offset := (page - 1) * pageSize

	// 执行分页查询
	var records []interface{}
	err = query.Offset(offset).Limit(pageSize).Find(&records).Error
	if err != nil {
		return nil, err
	}

	// 创建分页响应
	return &PaginationResult{
		Records:    records,
		Total:      count,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
