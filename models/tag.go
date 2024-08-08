package models

import (
	"gorm.io/gorm"
	"time"
)

// Tag 定义标签数据模型

type Tag struct {
	gorm.Model
	Name      string    `gorm:"unique;not null;index" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Diaries   []Diary   `gorm:"many2many:diary_tags"` // 多对多关系
}

/*
// BeforeDelete 删除标签之前的钩子

func (t *Tag) BeforeDelete(tx *gorm.DB) error {
	log.Printf("Deleting Tag with ID: %d and Name: %s\n", t.ID, t.Name)

	// 如果设置了外键的级联删除，则不需要手动删除 DiaryTag 记录
	// 以下代码是手动删除 DiaryTag 记录的示例
	if err := tx.Where("tag_id = ?", t.ID).Delete(&DiaryTag{}).Error; err != nil {
		log.Println("Error occurred while deleting related DiaryTag records:", err)
	}
	// 注意：这里不需要返回错误，因为 BeforeDelete 钩子不返回任何值

	// 如果需要执行其他删除前的操作，可以在这里添加
	return nil
}

// AfterFind 标签记录找到后的钩子

func (t *Tag) AfterFind(*gorm.DB) error {
	log.Printf("Retrieved Tag with ID: %d and Name: %s\n", t.ID, t.Name)
	// 如果需要执行其他查找后的操作，可以在这里添加
	return nil // AfterFind 钩子需要返回 error
}

// String 返回 Tag 的字符串表示形式

func (t *Tag) String() string {
	// 使用 fmt.Sprintf 直接格式化 uint 类型的 ID
	return fmt.Sprintf("Tag[ID:%d, Name:%s]", t.ID, t.Name)
}

// Validate 验证标签模型的数据

func (t *Tag) Validate() error {
	// 添加验证逻辑，例如检查名称是否为空
	if t.Name == "" {
		return gorm.ErrInvalidData
	}
	// 可以添加更多验证规则
	return nil
}

// ... 其他可能的方法 ...

/*
package models

import (
	"errors"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"time"
)

// Tag 定义标签数据模型
type Tag struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);unique;not null" json:"name"` // 标签名，唯一且非空
	// 可以添加其他字段，如标签描述、创建时间等
	Description string    `gorm:"type:text" json:"description,omitempty"` // 标签描述
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate 在创建标签之前的钩子
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	// 可以在此处添加创建前的逻辑，如格式化数据等
	return nil
}

// Validate 验证标签模型的数据
func (t Tag) Validate() error {
	// 检查标签名是否符合预期的格式
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(t.Name) {
		return errors.New("invalid tag name format")
	}
	// 可以添加更多验证规则
	return nil
}

// String 返回 Tag 的字符串表示形式
func (t Tag) String() string {
	return "Tag[ID:" + strconv.Itoa(int(t.ID)) + ", Name:" + t.Name + "]"
}

// Migrate 为 Tag 执行数据库迁移
func (t *Tag) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(t)
}

// 其他业务逻辑和方法可以根据需要添加
*/
