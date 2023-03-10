package rand

import (
	"math/rand"
	"time"
)

const (
	MaxWeight = 1000
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func Try(prob int32) bool {
	return prob >= rand.Int31n(MaxWeight)+1
}

// 可重复权重随机
// n = 3, probs = [1, 2, 3, 4, 5]
// ret = [1, 2, 2], [4, 2, 4]
func RepeatPerm(n int, probs []int32) []int {
	var ret []int

	total := int32(0)
	for _, prob := range probs {
		total += prob
	}

	for i := 0; i < n; i++ {
		if total <= 0 {
			return ret
		}

		rd := rand.Int31n(total) + 1

		for i, prob := range probs {
			rd -= prob
			if rd <= 0 {
				ret = append(ret, i)
				break
			}
		}
	}

	return ret
}

// 不重复权重随机
// n = 3, probs = [1, 2, 3, 4, 5]
// ret = [1, 2, 3], [1, 2, 4]
func UniquePerm(n int, probs []int32) []int {
	var ret []int

	if len(probs) == n {
		for i := 0; i < n; i++ {
			ret = append(ret, i)
		}

		return ret
	}

	total := int32(0)
	for _, prob := range probs {
		total += prob
	}

	record := map[int]bool{}

	for i := 0; i < n; i++ {
		if total <= 0 {
			return ret
		}

		rd := rand.Int31n(total) + 1

		for k, prob := range probs {
			if record[k] {
				continue
			} else {
				rd -= prob
			}

			if rd <= 0 {
				ret = append(ret, k)
				record[k] = true
				total -= prob
				break
			}
		}
	}

	return ret
}

// 权重随机一个
func SinglePerm(probs []int32) int {
	total := int32(0)
	for _, prob := range probs {
		total += prob
	}

	if total <= 0 {
		return 0
	}

	rd := rand.Int31n(total) + 1

	for i, prob := range probs {
		rd -= prob
		if rd <= 0 {
			return i
		}
	}

	return 0
}

// 纯概率随机
func TryPerm(probs []int32) []int {
	var ret []int

	for i, prob := range probs {
		if Try(prob) {
			ret = append(ret, i)
		}
	}

	return ret
}

func RangeInt32(low, up int32) int32 {
	if up < low {
		return 0
	}

	return rand.Int31n(up-low+1) + low
}

func RangeInt(low, up int) int {
	if up < low {
		return 0
	}

	return rand.Intn(up-low+1) + low
}

func Perm(length int) []int {
	return rand.Perm(length)
}
