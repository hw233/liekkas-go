package number

import "shared/utility/math"

type SortedInt64sSet []int64

func NewSortedInt64sSet() *SortedInt64sSet {
	return (*SortedInt64sSet)(&[]int64{})
}

func (s *SortedInt64sSet) Add(k int64) {
	_, ok := math.BinarySearch(*s, k)
	if ok {
		return
	}
	// 保持递增
	index, ok := math.RightBound(*s, k)
	if ok {
		*s = append((*s)[:index], append([]int64{k}, (*s)[index:]...)...)
	} else {
		*s = append(*s, k)

	}

}

// LeftBound 返回最接近且<target 的数组下标
func (s *SortedInt64sSet) LeftBound(target int64, count int) []int64 {
	var ret []int64
	index, ok := math.LeftBound(*s, target)
	if !ok {
		return ret
	}
	right := index
	left := right - count + 1
	if left < 0 {
		left = 0
	}
	for i := right; i >= left; i-- {
		ret = append(ret, (*s)[i])
	}
	return ret
}
func (s *SortedInt64sSet) FindIndex(target int64) (int, bool) {
	return math.BinarySearch(*s, target)

}

func (s *SortedInt64sSet) Delete(k int64) {
	findIndex, ok := s.FindIndex(k)
	if ok {
		*s = append((*s)[:findIndex], (*s)[findIndex+1:]...)
	}
}

// PagingSearch 返回key <offset 的n个数据
func (s *SortedInt64sSet) PagingSearch(offset int64, n int) []int64 {
	var values []int64

	index, ok := math.LeftBound(*s, offset)
	if !ok {
		return values
	}
	right := index
	if (*s)[index] == offset {
		right--
	}
	if right < 0 {
		return values
	}
	left := right - n + 1
	if left < 0 {
		left = 0
	}
	for i := right; i >= left; i-- {
		values = append(values, (*s)[i])
	}

	return values

}
