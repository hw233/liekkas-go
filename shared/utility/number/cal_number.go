package number

import (
	"encoding/json"
	"fmt"
	"math"
)

// CalNumber 计算数据类型，支持加减乘除等运算，主要防止数值溢出
type CalNumber int32

func NewCalNumber(init int32) *CalNumber {
	return (*CalNumber)(&init)
}

func (cn *CalNumber) Value() int32 {
	return int32(*cn)
}

func (cn *CalNumber) SetValue(num int32) {
	*cn = CalNumber(num)
}

func (cn *CalNumber) Plus(num int32) {
	if int64(*cn)+int64(num) <= int64(math.MaxInt32) {
		*cn += CalNumber(num)
	} else {
		*cn = math.MaxInt32
	}
}

func (cn *CalNumber) Minus(num int32) {
	if int64(*cn)-int64(num) >= int64(math.MinInt32) {
		*cn -= CalNumber(num)
	} else {
		*cn = math.MinInt32
	}
}

func (cn *CalNumber) Multi(num int32) {
	if int64(*cn)*int64(num) <= int64(math.MaxInt32) {
		*cn *= CalNumber(num)
	} else {
		*cn = math.MaxInt32
	}
}

func (cn *CalNumber) Divide(num int32) {
	if num == 0 {
		return
	}

	*cn /= CalNumber(num)
}

func (cn *CalNumber) Enough(num int32) bool {
	return int32(*cn) >= num
}

func (cn *CalNumber) Equal(num int32) bool {
	return int32(*cn) == num
}

// func (cn *CalNumber) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(cn)
// }
//
// func (cn *CalNumber) UnmarshalJSON(b []byte) error {
// 	return json.Unmarshal(b, cn)
// }
//
func (cn *CalNumber) Format(s fmt.State, verb rune) {
	val, err := json.Marshal(cn)
	if err == nil {
		_, _ = s.Write(val)
	}
}

//
// func (cn *CalNumber) Marshal() ([]byte, error) {
// 	return json.Marshal(cn)
// }
//
// func (cn *CalNumber) Unmarshal(b []byte) error {
// 	return json.Unmarshal(b, &cn)
// }
