package models

import (
	"gorm.io/gorm"
	"time"
)

// Tag 代表标签的数据库模型
type Tag struct {
	gorm.Model
	Name string `gorm:"unique;not null"` // 标签名称，需要唯一
	//ID        uint   `gorm:"not null"`        // 关联的用户ID
	CreatedAt time.Time      `json:"created_at,omitempty"`              // 创建时间，使用 omitempty 标签
	UpdatedAt time.Time      `json:"updated_at,omitempty"`              // 更新时间，使用 omitempty 标签
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 假删除字段，使用 omitempty 标签

	// 多对多关系，不返回空值
	Diaries []Diary `gorm:"many2many:diary_tags" json:"diaries,omitempty"` // 与其他日记的关联
}
