package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/utility/errors"
)

//----------------------------------------
//Store UnlockCondition
//----------------------------------------
// 检查总商店的解锁条件是否达到
func (u *User) checkStoreUnlockCondition(id int32) bool {
	storeCfg, err := manager.CSV.Store.GetStoreGeneral(id)
	if err != nil {
		return false
	}

	if u.CheckUserConditions(storeCfg.UnlockCondition) != nil {
		return false
	}

	return true
}

func (u *User) checkStageStoreForSoldOut(store *Store, subStoreID int32) error {
	// 1.获取所有子商店ID
	subIDs := store.SubStores

	// 2. 遍历所有ID小于参数subID的子商店，看是否SoldOut
	for _, subID := range subIDs {
		if subID >= subStoreID {
			continue
		}
		subStore, ok := u.StoreInfo.SubStores[subID]
		if !ok {
			return errors.Swrapf(common.ErrStoreWrongSubStoreIDForInfo, subID)
		}

		if !subStore.checkForSoldOut() { // 前面阶段的子商店没有售罄，则现在阶段的子商店不应该解锁
			return errors.Swrapf(common.ErrStoreSoldOutConditionNotSatisfied, subStoreID)
		}
	}

	return nil
}

//----------------------------------------
//Subtore UnlockCondition
//----------------------------------------

// 对替换式商店id下的所有子商店检查，看哪个是当前适合展示的。
func (u *User) SelectSubStore(sid int32) ([]int32, error) {
	subStoreIDs, err := manager.CSV.Store.GetSubIDs(sid)
	if err != nil {
		return nil, err
	}

	qualified := make([]int32, 0, len(subStoreIDs))

	for _, subID := range subStoreIDs {
		//fmt.Printf("========>user_store_condition subID: %v\n", subID)
		if u.checkSubStoreUnlockCondition(subID) {
			//fmt.Printf("========>user_store_condition subID after check: %v\n", subID)
			qualified = append(qualified, subID)
		}
	}

	if len(qualified) == 0 {
		return qualified, nil
	}

	var result []int32
	var tmp int32
	// 对于替换式商店，id越大代表商店等级越高
	for _, id := range qualified {
		if tmp < id {
			tmp = id
		}
	}
	result = append(result, tmp)

	return result, nil
}

// 对于阶段式商店，返回所有子商店的id
func (u *User) SelectSubStores(sid int32) ([]int32, error) {
	subStoreIDs, err := manager.CSV.Store.GetSubIDs(sid)
	if err != nil {
		return nil, err
	}

	qualified := []int32{}

	for _, subID := range subStoreIDs {
		if u.checkSubStoreUnlockCondition(subID) {
			qualified = append(qualified, subID)
		}
	}

	if len(qualified) == 0 {
		return qualified, nil
	}

	return subStoreIDs, nil
}

// 检查子商店的解锁条件是否达到
func (u *User) checkSubStoreUnlockCondition(subID int32) bool {
	subtoreCfg, err := manager.CSV.Store.GetSubStore(subID)
	if err != nil {
		return false
	}

	if u.CheckUserConditions(subtoreCfg.UnlockCondition) != nil {
		return false
	}

	return true
}

//----------------------------------------
//Cells UnlockCondition
//----------------------------------------

// 检查某子商店某个栏位的解锁条件是否达到
func (u *User) checkCellUnlockCondition(id int32, index int32) bool {
	subtoreCfg, err := manager.CSV.Store.GetSubStore(id)
	if err != nil {
		return false
	}

	if u.CheckUserConditions(subtoreCfg.Cells[index].UnlockCondition) != nil {
		return false
	}

	return true
}

//----------------------------------------
//Goods UnlockCondition
//----------------------------------------

// 检查某个商品的解锁条件是否达到
func (u *User) checkGoodsUnlockCondition(id int32) bool {
	goodsCfg, err := manager.CSV.Store.GetGoods(id)
	if err != nil {
		return false
	}

	if u.CheckUserConditions(goodsCfg.UnlockCondition) != nil {
		return false
	}

	return true
}
