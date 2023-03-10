package uuid

import (
	"testing"
)

func TestGenUUID(t *testing.T) {
	t.Run("Conflict", func(t *testing.T) {
		m := map[string]bool{}

		for i := 0; i < 10000000; i++ {
			if !m[GenUUID()] {
				m[GenUUID()] = true
			} else {
				t.Error("uuid conflict")
			}
		}
	})
}

func BenchmarkGenUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenUUID()
	}
}
