package models

import (
	"gorm.io/gorm"
)

// DiaryTag 定义日记和标签的多对多关联数据模型
type DiaryTag struct {
	gorm.Model
	DiaryID uint `gorm:"not null;index:idx_diary_tag_diary_id,unique" json:"diary_id"`
	TagID   uint `gorm:"not null;index:idx_diary_tag_tag_id,unique" json:"tag_id"`
	/*
		Diary   Diary `gorm:"foreignKey:DiaryID"` // 如果需要自动加载日记信息
		Tag     Tag   `gorm:"foreignKey:TagID"`   // 如果需要自动加载标签信息
	*/
}

/*
// BeforeDelete 删除 DiaryTag 之前的钩子
// 当 Diary 或 Tag 被删除时，自动删除 DiaryTag 中的对应记录
func (dt *DiaryTag) BeforeDelete(*gorm.DB) error {
	// 可以在这里添加额外的删除前逻辑
	log.Printf("Deleting DiaryTag with DiaryID: %d and TagID: %d\n", dt.DiaryID, dt.TagID)
	return nil
}

// String 返回 DiaryTag 的字符串表示形式
func (dt *DiaryTag) String() string {
	return "DiaryTag[DiaryID:" + strconv.Itoa(int(dt.DiaryID)) + ", TagID:" + strconv.Itoa(int(dt.TagID)) + "]"
}

// Validate 验证 DiaryTag 的数据
func (dt *DiaryTag) Validate() error {
	// 可以添加更复杂的验证逻辑
	if dt.DiaryID == 0 || dt.TagID == 0 {
		return gorm.ErrInvalidData
	}
	return nil
}

// ... 其他可能的方法 ...
*/
