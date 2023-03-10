package number

import "encoding/json"

type LimitedNumber struct {
	*CalNumber
	min int32
	max int32
}

func NewLimitedNumber(init, min, max int32) *LimitedNumber {
	ln := &LimitedNumber{
		CalNumber: NewCalNumber(init),
		min:       min,
		max:       max,
	}
	ln.limit()

	return ln
}

func (ln *LimitedNumber) limit() {
	if ln.CalNumber.Value() > ln.max {
		ln.CalNumber.SetValue(ln.max)
	} else if ln.CalNumber.Value() < ln.min {
		ln.CalNumber.SetValue(ln.min)
	}
}

func (ln *LimitedNumber) SetLimit(min int32, max int32) {
	ln.min = min
	ln.max = max

	ln.limit()
}

func (ln *LimitedNumber) Plus(num int32) {
	ln.CalNumber.Plus(num)
	ln.limit()
}

func (ln *LimitedNumber) Minus(num int32) {
	ln.CalNumber.Minus(num)
	ln.limit()
}

func (ln *LimitedNumber) Multi(num int32) {
	ln.CalNumber.Multi(num)
	ln.limit()
}

func (ln *LimitedNumber) Divide(num int32) {
	ln.CalNumber.Multi(num)
	ln.limit()
}

func (ln *LimitedNumber) Marshal() ([]byte, error) {
	return json.Marshal(ln.CalNumber)
}

func (ln *LimitedNumber) Unmarshal(b []byte) error {
	if ln.CalNumber == nil {
		calNumber := NewCalNumber(0)
		err := json.Unmarshal(b, calNumber)
		if err != nil {
			return err
		}
		ln.CalNumber = calNumber
		return nil
	} else {
		return json.Unmarshal(b, ln.CalNumber)
	}
}

func (ln *LimitedNumber) MarshalJSON() ([]byte, error) {
	return json.Marshal(ln.CalNumber)
}

func (ln *LimitedNumber) UnmarshalJSON(b []byte) error {
	if ln.CalNumber == nil {
		calNumber := NewCalNumber(0)
		err := json.Unmarshal(b, calNumber)
		if err != nil {
			return err
		}
		ln.CalNumber = calNumber
		return nil
	} else {
		return json.Unmarshal(b, ln.CalNumber)
	}
}
