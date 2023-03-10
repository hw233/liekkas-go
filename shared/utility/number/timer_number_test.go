package number

import (
	"testing"
	"time"
)

func TestTimerNumber(t *testing.T) {
	timerNumber := NewTimerNumber(0)
	timerNumber.RegisterInterval(1 * time.Second)
	timerNumber.SetTimerUpper(8)

	time.Sleep(5 * time.Second)
	if !timerNumber.Equal(5) {
		t.Errorf("!timerNumber.Equal(%d), value: %d", 5, timerNumber.Value())
	}

	t.Logf("value: %d", timerNumber.Value())

	time.Sleep(5 * time.Second)
	if !timerNumber.Equal(8) {
		t.Errorf("!timerNumber.Equal(%d), value: %d", 8, timerNumber.Value())
	}

	t.Logf("value: %d", timerNumber.Value())

	time.Sleep(5 * time.Second)
	timerNumber.SetTimerUpper(20)
	time.Sleep(5 * time.Second)
	if !timerNumber.Equal(13) {
		t.Errorf("!timerNumber.Equal(%d), value: %d", 13, timerNumber.Value())
	}

	t.Logf("value: %d", timerNumber.Value())
}

func TestTimerNumber2(t *testing.T) {
	timerNumber := NewTimerNumber(0)
	timerNumber.RegisterInterval(1 * time.Second)
	timerNumber.SetTimerUpper(2)

	timerNumber.Plus(10)
	t.Logf("value: %d", timerNumber.Value())

	time.Sleep(5 * time.Second)
	t.Logf("value: %d", timerNumber.Value())
	timerNumber.SetTimerUpper(20)
	t.Logf("value: %d", timerNumber.Value())
	time.Sleep(5 * time.Second)
	t.Logf("value: %d", timerNumber.Value())

}

// goos: darwin
// goarch: amd64
// pkg: shared/utility/number
// cpu: Intel(R) Core(TM) i7-1068NG7 CPU @ 2.30GHz
// BenchmarkTimerNumber
// BenchmarkTimerNumber-8   	12047143	        98.98 ns/op
func BenchmarkTimerNumber(b *testing.B) {
	timerNumber := NewTimerNumber(0)
	timerNumber.RegisterInterval(1 * time.Second)
	timerNumber.SetTimerUpper(8)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		timerNumber.Plus(1)
	}
}
