package math

import "testing"

func TestLeftBound(t *testing.T) {
	t.Log(LeftBound([]int64{}, 1))
	t.Log(LeftBound([]int64{1}, 1))
	t.Log(LeftBound([]int64{0}, 1))
	t.Log(LeftBound([]int64{2}, 1))
	t.Log(LeftBound([]int64{1, 2, 3, 4}, 1))
	t.Log(LeftBound([]int64{0, 2, 3, 4}, 1))

}

func TestRightBound(t *testing.T) {
	t.Log(RightBound([]int64{}, 1))
	t.Log(RightBound([]int64{1}, 1))
	t.Log(RightBound([]int64{0}, 1))
	t.Log(RightBound([]int64{2}, 1))
	t.Log(RightBound([]int64{1, 2, 3, 4}, 1))
	t.Log(RightBound([]int64{0, 2, 3, 4}, 5))
	t.Log(RightBound([]int64{0, 2, 3, 4}, -1))

}
