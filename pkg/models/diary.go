package models

import (
	"gorm.io/gorm"
	"time"
)

// Diary 代表日记的数据库模型
type Diary struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`  // 用户ID，外键
	Content   string    `gorm:"type:text"` // 日记内容
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	// 假删除字段
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	// 关联用户
	User User `gorm:"foreignKey:ID"`
	// 多对多关联标签
	Tags []Tag `gorm:"many2many:diary_tags"`
}
