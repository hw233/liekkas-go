package rand

import "testing"

func TestUniquePerm(t *testing.T) {
	r := UniquePerm(5, []int32{100, 100, 0, 0, 100, 100, 100})
	t.Logf("r: %v", r)
}
