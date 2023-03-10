package slice

import (
	"fmt"
	"testing"
)

func TestSliceHasEle(t *testing.T) {
	testSlice := []int32{1, 2, 3}
	num := int32(3)

	fmt.Printf("test result: %+v", SliceInt32HasEle(testSlice, num))
}
