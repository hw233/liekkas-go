package common

import (
	"testing"
)

func TestGachaResultRecords_RemoveIfOutOfDate(t *testing.T) {

	records := NewGachaResultRecords()
	var g []*GachaReward
	for i := 0; i < 365000; i++ {
		g = append(g, NewGachaReward(NewReward(1, 1), 1))
	}
	records.Add(1, g)
	records.RemoveIfOutOfDate(1)
	t.Log(records)

}

func BenchmarkGachaResultRecords_RemoveIfOutOfDate(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		records := NewGachaResultRecords()

		var g []*GachaReward
		for i := 0; i < 365000; i++ {
			g = append(g, NewGachaReward(NewReward(1, 1), 1))
		}
		b.ReportAllocs()

		records.Add(1, g)

		records.RemoveIfOutOfDate(0)
	}
	b.ReportAllocs()

}
