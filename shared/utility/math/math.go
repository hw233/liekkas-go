package math

func MaxInt32Between(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func MinInt32Between(x, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func AbsInt32(i int32) int32 {
	if i > 0 {
		return i
	}
	return -i
}

// BinarySearch ,nums递增数组
func BinarySearch(nums []int64, target int64) (int, bool) {
	left := 0
	right := len(nums) - 1
	for left <= right {
		mid := (left + right) / 2
		if nums[mid] == target {
			return mid, true
		} else if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid - 1
		}
	}

	return -1, false
}

// LeftBound 返回最接近且<=target 的数组下标,nums递增数组
func LeftBound(nums []int64, target int64) (int, bool) {
	totalLen := len(nums)
	if totalLen == 0 {
		return -1, false
	}
	left := 0
	right := totalLen
	for left < right {
		mid := (left + right) / 2
		if nums[mid] == target {
			left = mid + 1
		} else if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid
		}
	}
	ret := left - 1
	if ret < 0 {
		return -1, false
	}
	return ret, true

}

// RightBound 返回最接近且>=target 的数组下标,nums递增数组
func RightBound(nums []int64, target int64) (int, bool) {

	totalLen := len(nums)
	if totalLen == 0 {
		return -1, false
	}
	left := 0
	right := totalLen
	for left < right {
		mid := (left + right) / 2
		if nums[mid] == target {
			right = mid
		} else if nums[mid] < target {
			left = mid + 1
		} else if nums[mid] > target {
			right = mid
		}
	}
	if left >= totalLen {
		return -1, false
	}
	return left, true

}
