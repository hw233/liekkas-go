package common

import (
	"encoding/json"
	"math"

	"shared/statistic/logreason"
	"shared/utility/number"
)

// BIEventNumber 支持绑定事件在数值参与计算前后，适用于记录数值变化流水或者更新数值排行榜等需求
type BIEventNumber struct {
	*number.LimitedNumber
	preparedEvents []func()                 // execute before function call
	changedEvents  []func(*ChangedEventOpt) // execute after function call if number changed, will be not executed in prepared event
}

type ChangedEventOpt struct {
	NowValue int32
	Diff     int32
	Reason   *logreason.Reason
}

func NewBIEventNumber(init, min, max int32) *BIEventNumber {
	return &BIEventNumber{
		LimitedNumber:  number.NewLimitedNumber(init, min, max),
		preparedEvents: []func(){},
		changedEvents:  []func(*ChangedEventOpt){},
	}
}

func (en *BIEventNumber) execPreparedEvents() {
	for _, event := range en.preparedEvents {
		event()
	}
}

func (en *BIEventNumber) execChangedEvents(diff int32, reason *logreason.Reason) {
	for _, event := range en.changedEvents {
		event(&ChangedEventOpt{
			NowValue: en.Value(),
			Diff:     diff,
			Reason:   reason,
		})
	}
}

func (en *BIEventNumber) RegisterPreparedEvent(event func()) {
	en.preparedEvents = append(en.preparedEvents, event)
}

func (en *BIEventNumber) RegisterChangedEvent(event func(*ChangedEventOpt)) {
	en.changedEvents = append(en.changedEvents, event)
}

func (en *BIEventNumber) ClearPreparedEvents() {
	en.preparedEvents = []func(){}
}

func (en *BIEventNumber) ClearChangedEvents() {
	en.changedEvents = []func(*ChangedEventOpt){}
}

func (en *BIEventNumber) calculate(num int32, reason *logreason.Reason, f func(int32)) {
	en.execPreparedEvents()

	before := en.LimitedNumber.Value()

	f(num)

	after := en.LimitedNumber.Value()

	diff := after - before
	if diff != 0 {
		en.execChangedEvents(diff, reason)
	}
}

func (en *BIEventNumber) Value() int32 {
	en.execPreparedEvents()
	value := en.LimitedNumber.Value()

	return value
}

func (en *BIEventNumber) SetValue(num int32, reason *logreason.Reason) {
	en.calculate(num, reason, en.LimitedNumber.SetValue)
}

func (en *BIEventNumber) Plus(num int32, reason *logreason.Reason) {
	en.calculate(num, reason, en.LimitedNumber.Plus)
}

func (en *BIEventNumber) Minus(num int32, reason *logreason.Reason) {
	en.calculate(num, reason, en.LimitedNumber.Minus)
}

func (en *BIEventNumber) Multi(num int32, reason *logreason.Reason) {
	en.calculate(num, reason, en.LimitedNumber.Multi)
}

func (en *BIEventNumber) Divide(num int32, reason *logreason.Reason) {
	en.calculate(num, reason, en.LimitedNumber.Divide)
}

func (en *BIEventNumber) Enough(num int32) bool {
	en.execPreparedEvents()
	ok := en.LimitedNumber.Enough(num)

	return ok
}

func (en *BIEventNumber) Equal(num int32) bool {
	en.execPreparedEvents()
	ok := en.LimitedNumber.Equal(num)

	return ok
}

func (en *BIEventNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(en.LimitedNumber)
}

func (en *BIEventNumber) UnmarshalJSON(b []byte) error {
	if en.LimitedNumber == nil {
		LimitedNumber := number.NewLimitedNumber(0, math.MinInt32, math.MaxInt32)
		err := json.Unmarshal(b, LimitedNumber)
		if err != nil {
			return err
		}
		en.LimitedNumber = LimitedNumber
		return nil
	} else {
		return json.Unmarshal(b, en.LimitedNumber)
	}

}

// func (en *BIEventNumber) Format(s fmt.State, verb rune) {
// 	_, _ = s.Write([]byte(fmt.Sprint(en.LimitedNumber)))
// }

func (en *BIEventNumber) Marshal() ([]byte, error) {
	return json.Marshal(en.LimitedNumber)
}

func (en *BIEventNumber) Unmarshal(b []byte) error {
	if en.LimitedNumber == nil {
		LimitedNumber := number.NewLimitedNumber(0, math.MinInt32, math.MaxInt32)
		err := json.Unmarshal(b, LimitedNumber)
		if err != nil {
			return err
		}
		en.LimitedNumber = LimitedNumber
		return nil
	} else {
		return json.Unmarshal(b, en.LimitedNumber)
	}
}
