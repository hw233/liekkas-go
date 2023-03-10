package rand

import "math"

// GrabRedPacket 世界探索匹配 抢红包算法:
// total 总数，
// firstBonus ，小数 如果是第一个抢最多可以抢total*firstBonus
func GrabRedPacket(total int32, firstBonus float64, first bool) int32 {
	if total == 0 || total == 1 {
		return total
	}
	var min int32 = 1
	max := total
	if first {
		max = int32(math.Ceil(float64(total) * firstBonus))
	}
	// 生成随机值
	return RangeInt32(min, max)
}
