package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
)

// 装备进阶
func (u *User) WorldItemAdvance(targetID int64, materialIDs []int64, itemID int32) (*common.WorldItem, error) {
	// 检查材料，不能吃自己
	err := common.CheckWorldItemMaterial(targetID, materialIDs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取没有上锁材料装备信息
	materials, err := u.WorldItemPack.BatchGetUnlocked(materialIDs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查材料是否被装备
	err = common.CheckWorldItemWear(materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取目标装备信息
	target, err := u.WorldItemPack.Get(targetID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查阶级是否已经满了
	err = manager.CSV.WorldItem.CheckStageUpToLimit(target)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查目标装备的强化材料是否匹配
	err = manager.CSV.WorldItem.CheckAdvanceMaterials(target, materials, itemID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 计算强化阶级
	addStage, err := manager.CSV.WorldItem.CalAddStage(target, materials, itemID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取强化消耗的金币
	goldCost, err := manager.CSV.WorldItem.AdvanceCostGold(target)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 消耗
	costRewards := common.NewRewards()
	costRewards.AddReward(common.NewReward(static.CommonResourceTypeMoney, goldCost*addStage))
	if itemID != 0 {
		// 升星道具
		costRewards.AddReward(common.NewReward(itemID, 1))
	}

	// 检查消耗是否足够
	err = u.CheckRewardsEnough(costRewards)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	oldStage := target.Stage.Value()

	// 进阶
	target.Stage.Plus(addStage)

	// 扣除金币
	reason := logreason.NewReason(logreason.WorldItemAdvance)
	err = u.CostRewards(costRewards, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 销毁装备材料
	u.WorldItemPack.BatchDestroy(materialIDs, u)

	nowStage := target.Stage.Value()
	u.TriggerQuestUpdate(static.TaskTypeWorldStarCount, target.Rarity, oldStage, nowStage)

	u.BIWorldItemOp(target, bilog.WorldItemOpStageUp, logreason.EmptyReason())

	return target, nil
}

// 装备强化
func (u *User) WorldItemStrengthen(targetID int64, expItems []int32, materialIDs []int64) (*common.WorldItem, error) {
	// 检查材料，不能吃自己
	err := common.CheckWorldItemMaterial(targetID, materialIDs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取没有上锁材料装备信息
	materials, err := u.WorldItemPack.BatchGetUnlocked(materialIDs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查材料是否被装备
	err = common.CheckWorldItemWear(materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取目标装备信息
	target, err := u.WorldItemPack.Get(targetID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	maxLevel, err := manager.CSV.TeamLevelCache.GetWorldItemLevelUp(u.Info.Level.Value())
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取经验药水的经验
	itemEXP := manager.CSV.Item.CalculateTotalWorldItemEXP(expItems)

	// 获取强化增加的经验
	addEXP, err := manager.CSV.WorldItem.StrengthenEXP(target, itemEXP, materials, maxLevel)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 计算金币消耗
	goldCost, err := manager.CSV.WorldItem.StrengthenGoldCost(target, itemEXP, materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查金币是否足够
	// 消耗
	costRewards := common.NewRewards()

	costRewards.AddReward(common.NewReward(static.CommonResourceTypeMoney, goldCost))

	for i, v := range common.WorldItemEXPItems {
		if expItems[i] > 0 {
			costRewards.AddReward(common.NewReward(v, expItems[i]))
		}
	}

	// 检查金币是否足够
	err = u.CheckRewardsEnough(costRewards)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	oldLevel := target.Level.Value()

	// 增加经验
	target.EXP.Plus(addEXP)

	// 同步等级并解锁随机属性
	err = manager.CSV.WorldItem.SyncLevelAndUnlockAttr(target)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 扣除金币
	reason := logreason.NewReason(logreason.WorldItemStrengthen)
	err = u.CostRewards(costRewards, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 销毁装备
	u.WorldItemPack.BatchDestroy(materialIDs, u)

	u.TriggerQuestUpdate(static.TaskTypeWorldItemLevelupTimes, int32(target.Rarity))

	nowLevel := target.Level.Value()
	if nowLevel > oldLevel {
		u.TriggerQuestUpdate(static.TaskTypeWorldLevelCount, target.Rarity, oldLevel, nowLevel)
	}

	u.BIWorldItemOp(target, bilog.WorldItemOpLevelUp, logreason.EmptyReason())

	return target, nil
}

// 装备加锁解锁
func (u *User) WorldItemLock(id int64, lock bool) error {
	// 获取目标装备信息
	equipment, err := u.WorldItemPack.Get(id)
	if err != nil {
		return errors.WrapTrace(err)
	}

	equipment.IsLock = lock

	return nil
}

// 装备重铸阵营
// func (u *User) WorldItemRecastCamp(id int64) (int8, error) {
// 	// 获取目标装备信息
// 	equipment, err := u.WorldItemPack.Get(id)
// 	if err != nil {
// 		return 0, errors.WrapTrace(err)
// 	}
//
// 	// 重铸消耗
// 	costRewards, err := manager.CSV.WorldItem.RecastCost(equipment)
// 	if err != nil {
// 		return 0, errors.WrapTrace(err)
// 	}
//
// 	err = u.CheckRewardsEnough(costRewards)
// 	if err != nil {
// 		return 0, errors.WrapTrace(err)
// 	}
//
// 	camp, err := manager.CSV.WorldItem.RecastCamp(equipment)
// 	if err != nil {
// 		return 0, errors.WrapTrace(err)
// 	}
//
// 	// 保存重铸结果
// 	u.WorldItemPack.SaveRecastCamp(id, camp)
//
// 	// 扣除消耗
// 	err = u.CostRewards(costRewards)
// 	if err != nil {
// 		return 0, errors.WrapTrace(err)
// 	}
//
// 	return camp, nil
// }
//
// // 装备重铸阵营确认选择
// func (u *User) WorldItemConfirmRecastCamp(confirm bool) (*common.WorldItem, error) {
// 	// 检查上次重铸是否确认
// 	err := u.WorldItemPack.CheckConfirmRecastCamp()
// 	if err != nil {
// 		return nil, errors.WrapTrace(err)
// 	}
//
// 	// 确认重铸结果
// 	equipment, err := u.WorldItemPack.ConfirmRecastCamp(confirm)
// 	if err != nil {
// 		return nil, errors.WrapTrace(err)
// 	}
//
// 	return equipment, nil
// }
