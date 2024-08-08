package models

import (
	"gorm.io/gorm"
	"time"
)

// Diary 定义日记数据模型
type Diary struct {
	gorm.Model
	UserID    uint      `gorm:"not null" json:"user_id"`
	Content   string    `gorm:"type:text" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    int       `gorm:"default:0" json:"status"`
	// 实现软删除
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Tags      []Tag          `gorm:"many2many:diary_tags"` // 多对多关系
	DiaryTags []DiaryTag     `gorm:"foreignKey:DiaryID"`
}

/*
// Validate 验证日记模型的数据
func (d *Diary) Validate() error {
	// 添加验证逻辑，例如检查内容是否为空
	if d.Content == "" {
		return errors.New("content cannot be empty")
	}
	// 可以添加更多验证规则
	return nil
}

// String 返回日记的字符串表示形式
func (d *Diary) String() string {
	// 格式化输出日记内容和时间
	return d.Content + " - " + d.CreatedAt.String()
}

// ... 其他可能的方法 ...
*/
