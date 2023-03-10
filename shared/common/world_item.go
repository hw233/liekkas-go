package common

import (
	"time"

	"shared/protobuf/pb"
	"shared/utility/number"
)

var WorldItemEXPItems = []int32{151001, 151002, 151003, 151004}

type WorldItem struct {
	ID     int64             `json:"id"`     // 唯一ID
	WID    int32             `json:"wid"`    // 种类ID
	CID    int32             `json:"cid"`    // 装备的角色ID
	EXP    *number.CalNumber `json:"exp"`    // 经验值
	Level  *number.CalNumber `json:"level"`  // 等级
	Stage  *number.CalNumber `json:"stage"`  // 阶级
	Rarity int8              `json:"rarity"` // 稀有度，使用频率非常高，所以存一下
	IsLock bool              `json:"lock"`   // 是否上锁
	CTime  int64             `json:"ctime"`  // 获得时间
}

func NewWorldItem(id int64, wid int32) *WorldItem {
	equipment := &WorldItem{
		ID:     id,
		WID:    wid,
		CID:    0,
		EXP:    number.NewCalNumber(0),
		Level:  number.NewCalNumber(1),
		Stage:  number.NewCalNumber(0),
		Rarity: 0,
		IsLock: false,
		CTime:  time.Now().Unix(),
	}

	return equipment
}

func (e *WorldItem) VOUserWorldItem() *pb.VOUserWorldItem {
	// attrs := make([]*pb.VOUWorldItemAttr, 0, len(e.Attrs))
	// for _, attr := range e.Attrs {
	// 	attrs = append(attrs, &pb.VOUWorldItemAttr{Attr: int32(attr.Attr), Value: attr.Value})
	// }

	return &pb.VOUserWorldItem{
		WorldItemUId: e.ID,
		WorldItemId:  e.WID,
		CharacterId:  e.CID,
		Exp:          e.EXP.Value(),
		Level:        e.Level.Value(),
		Star:         e.Stage.Value(),
		LockFlag:     e.IsLock,
		CreateAt:     e.CTime,
	}
}

// 检查材料，材料不能是自己
func CheckWorldItemMaterial(targetID int64, materialIDs []int64) error {
	for _, v := range materialIDs {
		if targetID == v {
			return ErrWorldItemCantUseSelf
		}
	}

	return nil
}

// 检查材料是否被穿戴
func CheckWorldItemWear(materials []*WorldItem) error {
	for _, material := range materials {
		if material.CID != 0 {
			return ErrWorldItemWearing
		}
	}

	return nil
}
