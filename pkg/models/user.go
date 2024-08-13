package models

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	gorm.Model
	PhoneNumber string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// 假删除字段
	IsDeleted bool `gorm:"default:false"`
}

// CheckPassword 检查密码是否正确
func (u *User) CheckPassword(inputPassword string) bool {
	// 这里应该是密码比较逻辑，示例中假设密码以明文存储
	return u.Password == inputPassword
}
func (u *User) BeforeCreate(tx *gorm.DB) error {
	// 在这里添加创建用户之前的逻辑
	return nil
}

func (u *User) AfterCreate(tx *gorm.DB) error {
	// 在这里添加创建用户之后的逻辑
	return nil
}

// 其他 GORM 钩子...
