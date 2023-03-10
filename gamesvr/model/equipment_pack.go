package model

import (
	"encoding/json"

	"gamesvr/manager"
	"shared/common"
	"shared/protobuf/pb"

	"shared/utility/errors"
	"shared/utility/uid"
)

type EquipmentPack struct {
	UID        *uid.UID                    `json:"uid"`        // 生成唯一ID
	Equipments map[int64]*common.Equipment `json:"equipments"` // 装备数据
}

func NewEquipmentPack() *EquipmentPack {
	return &EquipmentPack{
		UID:        uid.NewUID(),
		Equipments: map[int64]*common.Equipment{},
	}
}

func (ep *EquipmentPack) NewEquipment(eid int32) (*common.Equipment, error) {
	// 检查装备ID是否配置
	err := manager.CSV.Equipment.CheckEIDExist(eid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 初始化数据
	equipment, err := manager.CSV.Equipment.NewEquipment(ep.UID.Gen(), eid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ep.Equipments[equipment.ID] = equipment

	return equipment, nil
}

func (ep *EquipmentPack) Equip(id int64, cid int32) error {
	e, ok := ep.Equipments[id]
	if !ok {
		return errors.Swrapf(common.ErrEquipmentNotExist, id)
	}

	e.CID = cid

	return nil
}

func (ep *EquipmentPack) Get(id int64) (*common.Equipment, error) {
	e, ok := ep.Equipments[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrEquipmentNotExist, id)
	}
	return e, nil
}

// 批量获取装备
func (ep *EquipmentPack) BatchGet(ids []int64) ([]*common.Equipment, error) {
	equipments := make([]*common.Equipment, 0, len(ids))

	// 取出装备信息
	for _, id := range ids {
		equipment, err := ep.Get(id)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

// 获取未上锁的装备，如果上锁就会报错，适用于强化进阶等会消耗装备的逻辑
func (ep *EquipmentPack) BatchGetUnlocked(ids []int64) ([]*common.Equipment, error) {
	equipments := make([]*common.Equipment, 0, len(ids))

	// 取出装备信息
	for _, id := range ids {
		equipment, err := ep.Get(id)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		if equipment.IsLocked {
			return nil, errors.Swrapf(common.ErrEquipmentLocked, id)
		}

		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

// 销毁装备
func (ep *EquipmentPack) Destroy(id int64) {
	delete(ep.Equipments, id)
}

// 批量销毁装备
func (ep *EquipmentPack) BatchDestroy(ids []int64) {
	for _, id := range ids {
		ep.Destroy(id)
	}
}

func (ep *EquipmentPack) Count() int {
	return len(ep.Equipments)
}

func (ep *EquipmentPack) VOUserEquipment() []*pb.VOUserEquipment {
	vos := make([]*pb.VOUserEquipment, 0, len(ep.Equipments))

	for _, v := range ep.Equipments {
		vos = append(vos, v.VOUserEquipment())
	}

	return vos
}

func (ep *EquipmentPack) Marshal() ([]byte, error) {
	return json.Marshal(ep)
}

func (ep *EquipmentPack) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, ep)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}
