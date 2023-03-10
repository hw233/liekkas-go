package rand

import "testing"

func TestGrabRedPacket(t *testing.T) {
	for j := 0; j < 100; j++ {
		first := true
		var total int32 = 100
		var result []int32
		for i := 0; i < 100; i++ {
			grabNum := GrabRedPacket(total, 0.4, first)
			first = false
			total -= grabNum
			result = append(result, grabNum)
		}
		t.Log(result)
	}

}
