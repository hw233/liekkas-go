package number

type BitNumber uint32

func NewBitNumber() *BitNumber {
	return new(BitNumber)
}

func (n *BitNumber) Get() uint32 {
	return uint32(*n)
}

func (n *BitNumber) Set(i uint32) {
	*n = BitNumber(i)
}

func (n *BitNumber) Clear() {
	n.Set(0)
}

func (n *BitNumber) Mark(i int) {
	*n |= 1 << i
}

func (n *BitNumber) IsMarked(i int) bool {
	return *n&(1<<i) != 0
}

func (n *BitNumber) IsLimited(i int) bool {
	return i < 32
}

func (n *BitNumber) Counts() int {
	count := 0
	for i := uint8(0); i < 32; i++ {
		if *n&(1<<i) != 0 {
			count++
		}
	}

	return count
}
