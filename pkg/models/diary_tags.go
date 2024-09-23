package models

type DiaryTags struct {
	DiaryID uint `gorm:"primaryKey;not null"`
	TagID   uint `gorm:"primaryKey;not null"`
}
