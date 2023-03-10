package number

import (
	"testing"
)

// TestCalNumberBase 测试常规的加减乘除
func TestCalNumberBase(t *testing.T) {
	calNumber := NewCalNumber(0)
	calNumber.SetValue(100) // value = 100
	if !calNumber.Equal(100) {
		t.Errorf("!calNumber.Equal(%d)", 100)
	}

	calNumber.Plus(100) // value = 100 + 100
	if !calNumber.Equal(200) {
		t.Errorf("!calNumber.Equal(%d)", 200)
	}

	calNumber.Minus(150) // value = 200 - 150
	if !calNumber.Equal(50) {
		t.Errorf("!calNumber.Equal(%d)", 50)
	}

	calNumber.Multi(10) // value = 50 * 10
	if !calNumber.Equal(500) {
		t.Errorf("!calNumber.Equal(%d)", 500)
	}

	calNumber.Divide(100) // value = 500 / 100
	if !calNumber.Equal(5) {
		t.Errorf("!calNumber.Equal(%d)", 5)
	}

	calNumber.Divide(0) // value = 5
	if !calNumber.Equal(5) {
		t.Errorf("!calNumber.Equal(%d)", 5)
	}

	if !calNumber.Enough(5) { // value >= 5
		t.Errorf("!calNumber.Equal(%d)", 5)
	}
}
