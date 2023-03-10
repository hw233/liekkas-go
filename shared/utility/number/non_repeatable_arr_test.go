package number

import "testing"

func TestRepeatableArr(t *testing.T) {
	arr := NewNonRepeatableArr()
	arr.Append(1)
	arr.Append(1)
	t.Log(arr)
}
