package uid

type UID int64

func NewUID() *UID {
	var init int64 = 1
	return (*UID)(&init)
}

func (u *UID) Gen() int64 {
	*u++
	return int64(*u - 1)
}
