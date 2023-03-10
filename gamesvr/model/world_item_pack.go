package model

import (
	"encoding/json"

	"gamesvr/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/uid"
)

type WorldItemPack struct {
	RecastID   int64                       `json:"recast_id"`   // 重铸随机的阵营
	RecastCamp int8                        `json:"recast_camp"` // 重铸随机的阵营
	UID        *uid.UID                    `json:"uid"`         // 生成唯一ID
	WorldItems map[int64]*common.WorldItem `json:"equipments"`  // 装备数据
}

func NewWorldItemPack() *WorldItemPack {
	return &WorldItemPack{
		RecastID:   0,
		RecastCamp: 0,
		UID:        uid.NewUID(),
		WorldItems: map[int64]*common.WorldItem{},
	}
}

func (wp *WorldItemPack) NewWorldItem(wid int32) (*common.WorldItem, error) {
	// 检查装备ID是否配置
	err := manager.CSV.WorldItem.CheckWIDExist(wid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 初始化数据
	worldItem, err := manager.CSV.WorldItem.NewWorldItem(wp.UID.Gen(), wid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	wp.WorldItems[worldItem.ID] = worldItem

	return worldItem, nil
}

func (wp *WorldItemPack) Equip(id int64, cid int32) error {
	e, ok := wp.WorldItems[id]
	if !ok {
		return errors.Swrapf(common.ErrWorldItemNotExist, id)
	}

	e.CID = cid

	return nil
}

func (wp *WorldItemPack) Get(id int64) (*common.WorldItem, error) {
	e, ok := wp.WorldItems[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrWorldItemNotExist, id)
	}
	return e, nil
}

// 批量获取装备
func (wp *WorldItemPack) BatchGet(ids []int64) ([]*common.WorldItem, error) {
	equipments := make([]*common.WorldItem, 0, len(ids))

	// 取出装备信息
	for _, id := range ids {
		equipment, err := wp.Get(id)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

// 获取未上锁的装备，如果上锁就会报错，适用于强化进阶等会消耗装备的逻辑
func (wp *WorldItemPack) BatchGetUnlocked(ids []int64) ([]*common.WorldItem, error) {
	equipments := make([]*common.WorldItem, 0, len(ids))

	// 取出装备信息
	for _, id := range ids {
		equipment, err := wp.Get(id)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		if equipment.IsLock {
			return nil, errors.Swrapf(common.ErrWorldItemLocked, id)
		}

		equipments = append(equipments, equipment)
	}

	return equipments, nil
}

// 销毁装备
func (wp *WorldItemPack) Destroy(id int64, u *User) {
	delete(wp.WorldItems, id)
	if id == u.Info.CardShow.WorldItemUId {
		u.Info.CardShow.WorldItemId = 0
		u.Info.CardShow.WorldItemUId = 0
	}
}

// 批量销毁装备
func (wp *WorldItemPack) BatchDestroy(ids []int64, u *User) {
	for _, id := range ids {
		wp.Destroy(id, u)
	}
}

// 检查有无确认上次重铸结果
// func (ep *WorldItemPack) CheckConfirmRecastCamp() error {
// 	if ep.RecastID != 0 {
// 		return errors.Swrapf(common.ErrWorldItemHasNotConfirm, ep.RecastID, ep.RecastCamp)
// 	}
//
// 	return nil
// }

// 保存重铸结果等待确认
// func (ep *WorldItemPack) SaveRecastCamp(id int64, camp int8) {
// 	ep.RecastID = id
// 	ep.RecastCamp = camp
// }
//
// // 确认重铸结果
// func (ep *WorldItemPack) ConfirmRecastCamp(confirm bool) (*common.WorldItem, error) {
// 	equipment, err := ep.Get(ep.RecastID)
// 	if err != nil {
// 		return nil, errors.WrapTrace(err)
// 	}
//
// 	if confirm {
// 		equipment.Camp = ep.RecastCamp
// 	}
//
// 	// 清除重铸数据
// 	ep.RecastID = 0
// 	ep.RecastCamp = 0
//
// 	return equipment, nil
// }

func (wp *WorldItemPack) Count() int {
	return len(wp.WorldItems)
}

func (wp *WorldItemPack) VOUserWorldItem() []*pb.VOUserWorldItem {
	vos := make([]*pb.VOUserWorldItem, 0, len(wp.WorldItems))

	for _, v := range wp.WorldItems {
		vos = append(vos, v.VOUserWorldItem())
	}

	return vos
}

func (wp *WorldItemPack) Marshal() ([]byte, error) {
	return json.Marshal(wp)
}

func (wp *WorldItemPack) Unmarshal(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	err := json.Unmarshal(data, wp)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}
