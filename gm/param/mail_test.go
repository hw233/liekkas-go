package param

import (
	"encoding/json"
	"shared/common"
	"testing"
)

func TestMail(t *testing.T) {
	param := &SendMailParam{
		UserIds:    []int64{1, 2, 3},
		Title:      "",
		Content:    "",
		StartTime:  "",
		Attachment: []common.Reward{*common.NewReward(1, 1)},
		ExpireDay:  8,
	}
	marshal, _ := json.Marshal(param)
	t.Log(string(marshal))
}
