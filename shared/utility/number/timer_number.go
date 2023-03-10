package number

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"shared/utility/servertime"
)

// TimerNumber 支持时间回复事件，适用于体力，能量等可以自己随着时间回复的数值类型需求
// TODO： 注意测试，容易出问题，尤其是玩家下线的情况
type TimerNumber struct {
	*EventNumber
	last  int64 // 最后通过时间回复数值的时间戳，需要落地！
	upper int32 // 时间回复数值的上限
}

func NewTimerNumber(init, min, max int32) *TimerNumber {
	return &TimerNumber{
		EventNumber: NewEventNumber(init, min, max),
		last:        servertime.Now().Unix(),
		upper:       -1, // 未初始化，不会修改last的值
	}
}

// SetTimerUpper 设置时间回复数值上限
// Design: 体力等上限会根据等级不断变化的需求情况，所以单独进行更新
func (tn *TimerNumber) SetTimerUpper(upper int32) {
	tn.execPreparedEvents()
	tn.upper = upper
}

// RegisterInterval 注册时间回复事件
// Parameter: interval 数值回复的时间间隔，单位秒
// Notice: 更新interval，需要调用ClearBeforeEvents再注册，否则会重复注册事件
func (tn *TimerNumber) RegisterInterval(interval time.Duration) {
	tn.RegisterPreparedEvent(func() {
		now := servertime.Now()

		if tn.upper == -1 {
			return
		}

		if tn.CalNumber.Value() >= tn.upper {
			tn.last = now.Unix()
			return
		}

		lastTime := time.Unix(tn.last, 0)

		truncate := now.Sub(lastTime).Truncate(interval)
		if truncate == 0 {
			return
		}

		addition := truncate / interval

		if addition == 0 {
			return
		}

		if tn.CalNumber.Value()+int32(addition) >= tn.upper {
			tn.CalNumber.SetValue(tn.upper)
			tn.last = now.Unix()
		} else {
			tn.CalNumber.Plus(int32(addition))
			tn.last = lastTime.Add(interval * addition).Unix()
		}
	})
}

type TimerNumberJSON struct {
	Num  int32 `json:"num"`
	Last int64 `json:"last"`
}

func (tn *TimerNumber) MarshalJSON() ([]byte, error) {
	tmJSON := TimerNumberJSON{
		Num:  tn.Value(),
		Last: tn.Last(),
	}

	return json.Marshal(&tmJSON)
}

func (tn *TimerNumber) UnmarshalJSON(b []byte) error {
	tmJSON := TimerNumberJSON{}

	err := json.Unmarshal(b, &tmJSON)
	if err != nil {
		return err
	}
	if tn.EventNumber == nil {
		tn.EventNumber = NewEventNumber(0, math.MinInt32, math.MaxInt32)
	}
	tn.CalNumber.SetValue(tmJSON.Num)
	tn.last = tmJSON.Last

	return nil
}

func (tn *TimerNumber) Format(s fmt.State, verb rune) {
	_, _ = s.Write([]byte(fmt.Sprintf("num: %d, last: %d", tn.CalNumber, tn.last)))
}

func (tn *TimerNumber) Last() int64 {
	return tn.last
}

func (tn *TimerNumber) Marshal() ([]byte, error) {
	return json.Marshal(tn)
}

func (tn *TimerNumber) Unmarshal(b []byte) error {
	err := json.Unmarshal(b, tn)
	if err != nil {
		return err
	}

	return nil
}
