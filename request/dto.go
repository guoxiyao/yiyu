package request

// DiaryDto 代表日记的请求模型
// Dto Data Transform Object
// TODO: Dto 以及数据转换
type DiaryDto struct {
	Content string `json:"content"` // 日记内容
	TagIds  []uint `json:"tagIds"`
}
