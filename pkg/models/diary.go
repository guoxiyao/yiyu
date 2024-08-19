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
	CreatedAt time.Time // 创建时间
	UpdatedAt time.Time // 更新时间

	// 假删除字段
	DeletedAt gorm.DeletedAt `gorm:"index"`
	//IsDeleted bool `gorm:"default:false"`
	// 关联用户
	User User `gorm:"foreignKey:ID"`
	// 多对多关联标签
	Tags []Tag `gorm:"many2many:diary_tags"`
}

/*
// Validate 验证日记数据
func (d *Diary) Validate() error {
	// 确保内容不为空
	if d.Content == "" {
		return errors.New("diary content cannot be empty")
	}
	// 可以添加更多的验证逻辑
	return nil
}
*/
