package number

import (
	"encoding/json"
	"math"
)

// EventNumber 支持绑定事件在数值参与计算前后，适用于记录数值变化流水或者更新数值排行榜等需求
type EventNumber struct {
	*LimitedNumber
	preparedEvents []func()      // execute before function call
	changedEvents  []func(int32) // execute after function call if number changed, will be not executed in prepared event
}

func NewEventNumber(init, min, max int32) *EventNumber {
	return &EventNumber{
		LimitedNumber:  NewLimitedNumber(init, min, max),
		preparedEvents: []func(){},
		changedEvents:  []func(int32){},
	}
}

func (en *EventNumber) execPreparedEvents() {
	for _, event := range en.preparedEvents {
		event()
	}
}

func (en *EventNumber) execChangedEvents(diff int32) {
	for _, event := range en.changedEvents {
		event(diff)
	}
}

func (en *EventNumber) RegisterPreparedEvent(event func()) {
	en.preparedEvents = append(en.preparedEvents, event)
}

func (en *EventNumber) RegisterChangedEvent(event func(int32)) {
	en.changedEvents = append(en.changedEvents, event)
}

func (en *EventNumber) ClearPreparedEvents() {
	en.preparedEvents = []func(){}
}

func (en *EventNumber) ClearChangedEvents() {
	en.changedEvents = []func(int32){}
}

func (en *EventNumber) calculate(num int32, f func(int32)) {
	en.execPreparedEvents()

	before := en.LimitedNumber.Value()

	f(num)

	after := en.LimitedNumber.Value()

	diff := after - before
	if diff != 0 {
		en.execChangedEvents(diff)
	}
}

func (en *EventNumber) Value() int32 {
	en.execPreparedEvents()
	value := en.LimitedNumber.Value()

	return value
}

func (en *EventNumber) SetValue(num int32) {
	en.calculate(num, en.LimitedNumber.SetValue)
}

func (en *EventNumber) Plus(num int32) {
	en.calculate(num, en.LimitedNumber.Plus)
}

func (en *EventNumber) Minus(num int32) {
	en.calculate(num, en.LimitedNumber.Minus)
}

func (en *EventNumber) Multi(num int32) {
	en.calculate(num, en.LimitedNumber.Multi)
}

func (en *EventNumber) Divide(num int32) {
	en.calculate(num, en.LimitedNumber.Divide)
}

func (en *EventNumber) Enough(num int32) bool {
	en.execPreparedEvents()
	ok := en.LimitedNumber.Enough(num)

	return ok
}

func (en *EventNumber) Equal(num int32) bool {
	en.execPreparedEvents()
	ok := en.LimitedNumber.Equal(num)

	return ok
}

func (en *EventNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(en.LimitedNumber)
}

func (en *EventNumber) UnmarshalJSON(b []byte) error {
	if en.LimitedNumber == nil {
		LimitedNumber := NewLimitedNumber(0, math.MinInt32, math.MaxInt32)
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

// func (en *EventNumber) Format(s fmt.State, verb rune) {
// 	_, _ = s.Write([]byte(fmt.Sprint(en.LimitedNumber)))
// }

func (en *EventNumber) Marshal() ([]byte, error) {
	return json.Marshal(en.LimitedNumber)
}

func (en *EventNumber) Unmarshal(b []byte) error {
	if en.LimitedNumber == nil {
		LimitedNumber := NewLimitedNumber(0, math.MinInt32, math.MaxInt32)
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
