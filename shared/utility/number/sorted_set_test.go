package number

import (
	"shared/utility/rand"
	"testing"
)

func TestNewYggdrasilMail(t *testing.T) {
	set := NewSortedInt64sSet()
	for i := 0; i < 1000; i++ {
		set.Add(int64(rand.RangeInt(0, 1000)))
	}
	t.Log(set)

}
