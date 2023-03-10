package slice

func SliceInt32HasEle(slice []int32, ele int32) bool {
	for _, v := range slice {
		if v == ele {
			return true
		}
	}

	return false
}

func SliceInt64HasEle(slice []int64, ele int64) bool {
	for _, v := range slice {
		if v == ele {
			return true
		}
	}

	return false
}

func SliceInt32Contain(mainSlice, subSlice []int32) bool {
	for _, ele := range subSlice {
		if !SliceInt32HasEle(mainSlice, ele) {
			return false
		}
	}

	return true
}

func SliceInt64Contain(mainSlice, subSlice []int64) bool {
	for _, ele := range subSlice {
		if !SliceInt64HasEle(mainSlice, ele) {
			return false
		}
	}

	return true
}

func GetSliceInt32Intersection(mainSlice, subSlice []int32) []int32 {
	record := map[int]bool{}
	intersection := []int32{}
	for _, ele := range subSlice {
		for idx, mainEle := range mainSlice {
			_, ok := record[idx]
			if !ok && ele == mainEle {
				intersection = append(intersection, ele)
				record[idx] = true
			}
		}
	}

	return intersection
}

func GetSliceInt64Intersection(mainSlice, subSlice []int64) []int64 {
	record := map[int]bool{}
	intersection := []int64{}
	for _, ele := range subSlice {
		for idx, mainEle := range mainSlice {
			_, ok := record[idx]
			if !ok && ele == mainEle {
				intersection = append(intersection, ele)
				record[idx] = true
			}
		}
	}

	return intersection
}
