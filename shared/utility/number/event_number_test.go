package number

import (
	"testing"
)

func TestEventNumber(t *testing.T) {
	count, difF := 0, int32(0)

	eventNumber := NewEventNumber(0)
	eventNumber.RegisterPreparedEvent(func() {
		count++
	})

	eventNumber.RegisterChangedEvent(func(diff int32) {
		difF = diff
	})

	eventNumber.Value() // value = 0
	if count != 1 || difF != 0 {
		t.Errorf("count != %d || difF != %d", 1, 0)
	}

	eventNumber.Plus(10) // value = 10
	if count != 2 || difF != 10 {
		t.Errorf("count != %d || difF != %d", 2, 10)
	}

	eventNumber.Minus(5) // value = 5
	if count != 3 || difF != -5 {
		t.Errorf("count != %d || difF != %d", 3, -5)
	}

	eventNumber.Multi(5) // value = 25
	if count != 4 || difF != 20 {
		t.Errorf("count != %d || difF != %d", 4, 20)
	}
}
