package models

import (
	"errors"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"time"
)

// User 定义用户数据模型

type User struct {
	gorm.Model
	PhoneNumber string    `gorm:"unique;not null" json:"phone_number"`
	Password    string    `gorm:"not null" json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// 实现软删除
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 确保DeletedAt被标记为-，以避免序列化
}

// BeforeCreate 在创建用户之前的钩子

func (u *User) BeforeCreate(*gorm.DB) error {
	// 密码加密处理
	u.Password = encryptPassword(u.Password)
	return nil
}

// Validate 验证用户模型的数据

func (u User) Validate() error {
	// 检查电话号码是否符合预期的格式
	if !regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(u.PhoneNumber) {
		return errors.New("invalid phone number format")
	}
	// 检查密码长度
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	// 可以添加更多验证规则
	return nil
}

// String 返回 User 的字符串表示形式

func (u User) String() string {
	return "User[ID:" + strconv.Itoa(int(u.ID)) + ", PhoneNumber:" + u.PhoneNumber + "]"
}

// 辅助函数，用于密码加密

func encryptPassword(password string) string {
	// 这里应实现密码加密逻辑，返回加密后的密码
	// 请使用合适的加密库和算法
	// 例如: return bcrypt.HashPassword(password)
	return password // 注意：这里仅作示例，实际应用中不应直接返回密码
}

/*
package models

import (
	"errors"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"time"
)

// User 定义用户数据模型
type User struct {
	gorm.Model
	PhoneNumber string    `gorm:"unique;not null" json:"phone_number"`
	Password    string    `gorm:"not null" json:"password"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	// 实现软删除
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // 确保DeletedAt被标记为-，以避免序列化
}

// BeforeCreate 在创建用户之前的钩子，用于加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.Password = encryptPassword(u.Password)
	return nil
}

// Validate 验证用户模型的数据
func (u User) Validate() error {
	// 检查电话号码是否符合预期的格式
	if !regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).MatchString(u.PhoneNumber) {
		return errors.New("invalid phone number format")
	}
	// 检查密码长度
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	// 可以添加更多验证规则
	return nil
}

// String 返回 User 的字符串表示形式
func (u User) String() string {
	return "User[ID:" + strconv.Itoa(int(u.ID)) + ", PhoneNumber:" + u.PhoneNumber + "]"
}

// Migrate 为 User 执行数据库迁移
func (u *User) Migrate(db *gorm.DB) error {
	return db.AutoMigrate(u)
}

// 辅助函数，用于密码加密
// 注意：实际应用中应使用安全的加密库，例如 bcrypt 或 argon2
func encryptPassword(password string) string {
	// 这里仅作示例，实际应用中不应直接返回密码
	// 应替换为安全的密码散列函数
	return "encrypted_password" // 请替换为实际的加密逻辑
}

// 这里可以添加其他模型定义，如 Diary, Tag 等，以及它们的 Migrate 方法
*/
