package common

import (
	"time"

	"shared/protobuf/pb"
	"shared/utility/number"
)

var EquipmentEXPItems = []int32{131001, 131002, 131003, 131004}

const (
	EquipmentRecastCampInit = 0
)

type Equipment struct {
	ID         int64             `json:"id"`          // 唯一ID
	EID        int32             `json:"eid"`         // 种类ID
	CID        int32             `json:"cid"`         // 装备的角色ID
	EXP        *number.CalNumber `json:"exp"`         // 经验值
	Level      *number.CalNumber `json:"level"`       // 等级
	Stage      *number.CalNumber `json:"stage"`       // 阶级
	Rarity     int8              `json:"rarity"`      // 稀有度，使用频率非常高，所以存一下
	IsLocked   bool              `json:"is_locked"`   // 是否上锁
	Camp       int8              `json:"camp"`        // 种族
	CTime      int64             `json:"ctime"`       // 获得时间
	Attrs      []EquipmentAttr   `json:"attrs"`       // 随机属性
	RecastCamp int8              `json:"recast_camp"` // 重铸随机的阵营
}

func NewEquipment(id int64, eid int32) *Equipment {
	equipment := &Equipment{
		ID:         id,
		EID:        eid,
		CID:        0,
		EXP:        number.NewCalNumber(0),
		Level:      number.NewCalNumber(1),
		Stage:      number.NewCalNumber(0),
		Rarity:     0,
		IsLocked:   false,
		Camp:       0,
		CTime:      time.Now().Unix(),
		Attrs:      []EquipmentAttr{},
		RecastCamp: EquipmentRecastCampInit,
	}

	return equipment
}

type EquipmentAttr struct {
	Attr   int8  `json:"attr"`   // 属性枚举
	Value  int32 `json:"value"`  // 属性数值
	Unlock bool  `json:"unlock"` // 属性是否解锁
}

func NewEquipmentAttr(attr int8, value int32) *EquipmentAttr {
	return &EquipmentAttr{
		Attr:   attr,
		Value:  value,
		Unlock: false,
	}
}

// 解锁新属性
func (e *Equipment) SetAttr(attrs []EquipmentAttr) {
	e.Attrs = attrs
}

func (e *Equipment) UnlockNextAttr() {
	for i, attr := range e.Attrs {
		if !attr.Unlock {
			e.Attrs[i].Unlock = true
			break
		}
	}
}

// 属性数量
func (e *Equipment) UnlockAttrNum() int {
	num := 0

	for _, attr := range e.Attrs {
		if attr.Unlock {
			num++
		}
	}

	return num
}

// 确认上次重铸是否完成
func (e *Equipment) IsLastRecastCampConfirmed() bool {
	return e.RecastCamp == EquipmentRecastCampInit
}

// 保存重铸结果等待确认
func (e *Equipment) SaveRecastCamp(camp int8) {
	e.RecastCamp = camp
}

// 确认重铸结果
func (e *Equipment) ConfirmRecastCamp(confirm bool) {
	if confirm {
		e.Camp = e.RecastCamp
	}

	// 清除重铸数据
	e.RecastCamp = EquipmentRecastCampInit
}

func (e *Equipment) VOUserEquipment() *pb.VOUserEquipment {
	attrs := make([]*pb.VOUEquipmentAttr, 0, len(e.Attrs))
	for _, attr := range e.Attrs {
		voAttr := &pb.VOUEquipmentAttr{Attr: int32(attr.Attr)}
		// 解锁的属性才返回数值
		if attr.Unlock {
			voAttr.Value = attr.Value
		}

		attrs = append(attrs, voAttr)
	}

	return &pb.VOUserEquipment{
		EquipmentUId: e.ID,
		EquipmentId:  e.EID,
		CharacterId:  e.CID,
		Exp:          e.EXP.Value(),
		Level:        e.Level.Value(),
		Stage:        e.Stage.Value(),
		LockFlag:     e.IsLocked,
		Camp:         int32(e.Camp),
		CampRandom:   0,
		CreateAt:     e.CTime,
		Attrs:        attrs,
		RecastCamp:   int32(e.RecastCamp),
	}
}

// 检查材料，材料不能是自己
func CheckEquipmentMaterial(targetID int64, materialIDs []int64) error {
	for _, v := range materialIDs {
		if targetID == v {
			return ErrEquipmentCantUseSelf
		}
	}

	return nil
}

// 检查材料是否被穿戴
func CheckEquipmentWear(materials []*Equipment) error {
	for _, material := range materials {
		if material.CID != 0 {
			return ErrEquipmentWearing
		}
	}

	return nil
}
