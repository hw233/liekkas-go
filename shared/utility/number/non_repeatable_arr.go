package number

type NonRepeatableArr []int32

func NewNonRepeatableArr() *NonRepeatableArr {
	return (*NonRepeatableArr)(&[]int32{})
}

func (n *NonRepeatableArr) Contains(value int32) bool {
	for _, i := range *n {
		if i == value {
			return true
		}
	}
	return false
}

func (n *NonRepeatableArr) Append(values ...int32) {
	for _, value := range values {
		if n.Contains(value) {
			return
		}
		*n = append(*n, value)
	}

}

func (n *NonRepeatableArr) Values() []int32 {
	return *n
}

func (n *NonRepeatableArr) IsEmpty() bool {
	return len(*n) == 0
}

// NonRepeatableArrInt64 int64的slice (需要按插入顺序)
type NonRepeatableArrInt64 []int64

func NewNonRepeatableArrInt64() *NonRepeatableArrInt64 {
	return (*NonRepeatableArrInt64)(&[]int64{})
}

func (n *NonRepeatableArrInt64) Contains(value int64) bool {
	for _, i := range *n {
		if i == value {
			return true
		}
	}
	return false
}
func (n *NonRepeatableArrInt64) Clear() {
	*n = []int64{}
}
func (n *NonRepeatableArrInt64) Append(values ...int64) {
	for _, value := range values {
		if n.Contains(value) {
			return
		}
		*n = append(*n, value)
	}

}
func (n *NonRepeatableArrInt64) Remove(value int64) {
	var newList []int64
	for i, v := range *n {
		if v == value {
			newList = append((*n)[:i], (*n)[i+1:]...)
			break
		}
	}
	*n = newList
}
func (n *NonRepeatableArrInt64) Values() []int64 {
	return *n
}

// ReverseValues 倒序值
func (n *NonRepeatableArrInt64) ReverseValues() []int64 {
	length := len(*n)
	for i := 0; i < length/2; i++ {
		(*n)[i], (*n)[length-1-i] = (*n)[length-1-i], (*n)[i]
	}

	return *n
}

func (n *NonRepeatableArrInt64) IsEmpty() bool {
	return len(*n) == 0
}
