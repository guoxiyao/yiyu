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
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//IsDeleted bool `gorm:"default:false"` // 假删除字段
	// 其他字段，如关联的日记等...
	Diaries []Diary `gorm:"many2many:diary_tags"` // 多对多关系
}

/*
// Validate 验证标签数据
func (t *Tag) Validate() error {
	// 确保名称不为空
	if t.Name == "" {
		return errors.New("tag name cannot be empty")
	}
	// 可以添加更多的验证逻辑
	return nil
}
*/
