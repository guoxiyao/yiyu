package response

import (
	"awesomeProject1/pkg"
	"awesomeProject1/pkg/models"
)

// DiaryVo 用于API响应的日记信息
// Vo View Object
// TODO: 出参Vo 以及数据转换
type DiaryVo struct {
	ID        uint    `json:"id"`
	CreatedAt string  `json:"createdAt,omitempty"`
	UpdatedAt string  `json:"updatedAt,omitempty"`
	UserID    uint    `json:"userID"`
	Content   string  `json:"content"`
	Tags      []TagVo `json:"tags,omitempty"`
}

func (dairyVo *DiaryVo) Copy(diary models.Diary) {
	dairyVo = &DiaryVo{
		ID:        diary.ID,
		Content:   diary.Content,
		CreatedAt: diary.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: diary.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// UserVo 用于API响应的用户信息
type UserVo struct {
	ID          uint   `json:"id"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

func (userVo *UserVo) Copy(user models.User) error {
	err := pkg.SimpleCopyProperties(userVo, user)
	if err != nil {
		return err
	}
	return nil
}

// TagVo 用于API响应的标签信息
type TagVo struct {
	ID        uint   `json:"id"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Name      string `json:"name"`
}
