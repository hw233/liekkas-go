package model

import (
	"encoding/json"
	"math"

	"shared/protobuf/pb"
	"shared/utility/glog"
	"shared/utility/number"
)

type ItemPack map[int32]*number.EventNumber

func NewItemPack() *ItemPack {
	return (*ItemPack)(&map[int32]*number.EventNumber{})
}

func (ip *ItemPack) Add(id, num int32) (before, after int32) {
	v, ok := (*ip)[id]
	if !ok {
		before = 0
		v = number.NewEventNumber(num, 0, math.MaxInt32)
		(*ip)[id] = v
		return before, v.Value()
	}

	before = v.Value()
	v.Plus(num)
	return before, v.Value()
}

func (ip *ItemPack) Minus(id, num int32) (before, after int32) {
	v, ok := (*ip)[id]
	if !ok {
		return 0, 0
	}

	before = v.Value()
	v.Minus(num)
	val := v.Value()
	if val == 0 {
		delete(*ip, id)
		return before, 0
	}

	return before, val
}

func (ip *ItemPack) Enough(id, num int32) bool {
	v, ok := (*ip)[id]
	if !ok {
		return false
	}

	return v.Enough(num)
}

func (ip *ItemPack) Count(id int32) int32 {
	v, ok := (*ip)[id]
	if !ok {
		return 0
	}

	return v.Value()
}

func (ip *ItemPack) SetLimit(id, limit int32) {
	v, ok := (*ip)[id]
	if !ok {
		return
	}

	v.SetLimit(0, limit)
}

func (ip *ItemPack) Marshal() ([]byte, error) {
	return json.Marshal(ip)
}

func (ip *ItemPack) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, &ip)
	if err != nil {
		glog.Errorf("json.Unmarshal error: %v", err)
		return err
	}

	return nil
}

func (ip *ItemPack) VOItemInfo() []*pb.VOItemInfo {
	vos := make([]*pb.VOItemInfo, 0, len(*ip))

	for k, v := range *ip {
		vos = append(vos, &pb.VOItemInfo{
			ItemId: k,
			Amount: v.Value(),
		})
	}

	return vos
}
