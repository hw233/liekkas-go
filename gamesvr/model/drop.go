package model

import (
	"errors"

	"shared/utility/number"
	"shared/utility/rand"
)

type Drop map[int32]*number.BitNumber

func NewDrop() *Drop {
	return (*Drop)(&map[int32]*number.BitNumber{})
}

// 给定长度, 随机下标，不会随机到标记过的
func (d *Drop) RandIndex(rewardID int32, length int) (int, error) {
	if length <= 0 {
		return 0, errors.New("length invalid")
	}

	bn, ok := (*d)[rewardID]
	if !ok {
		return rand.RangeInt(0, length-1), nil
	}

	is := rand.Perm(length)
	for _, i := range is {
		if !bn.IsMarked(i) {
			return i, nil
		}
	}

	return 0, errors.New("all dropped")
}

// 标记下标，下次不会随机到
func (d *Drop) MarkIndex(rewardID int32, index int) {
	b, ok := (*d)[rewardID]
	if !ok {
		b = number.NewBitNumber()
		b.Mark(index)
		(*d)[rewardID] = b
	}

	b.Mark(index)
}

// 如果全部随机到，清空标记
func (d *Drop) ClearIfAllDropped(rewardID int32, length int) error {
	if length <= 0 {
		return errors.New("length invalid")
	}

	bn, ok := (*d)[rewardID]
	if !ok {
		// no need refresh
		return nil
	}

	if bn.Counts() >= length {
		// all dropped
		bn.Clear()
	}

	return nil
}
