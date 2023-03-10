package number

import "testing"

func TestBitNumberBase(t *testing.T) {
	bitNum := NewBitNumber()
	bitNum.Mark(1)
	bitNum.IsMarked(1)
	t.Logf("%v", bitNum)
}
