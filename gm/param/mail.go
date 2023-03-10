package param

import (
	"shared/common"
)

type SendMailParam struct {
	UserIds    []*UserId       `json:"user_ids"`
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	StartTime  string          `json:"start_time"`
	Attachment []common.Reward `json:"attachment"`
	ExpireDay  int32           `json:"expire_day"`
}

type UserId struct {
	Id int64 `json:"uid"`
}

type SendWholeServerMailParam struct {
	Title      string          `json:"title"`
	Content    string          `json:"content"`
	StartTime  string          `json:"start_time"`
	Attachment []common.Reward `json:"attachment"`
	ExpireDay  int32           `json:"expire_day"`
}
