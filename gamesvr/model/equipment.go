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
func (u *User) EquipmentAdvance(targetID int64, materialIDs []int64) (*common.Equipment, error) {
	err := u.CheckActionUnlock(static.ActionIdTypeEquipbreakunlock)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查材料，不能吃自己
	err = common.CheckEquipmentMaterial(targetID, materialIDs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取没有上锁材料装备信息
	materials, err := u.EquipmentPack.BatchGetUnlocked(materialIDs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查材料是否被装备
	err = common.CheckEquipmentWear(materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取目标装备信息
	target, err := u.EquipmentPack.Get(targetID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查阶级是否已经满了
	err = manager.CSV.Equipment.CheckStageUpToLimit(target)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查目标装备等级是否满足进阶条件
	err = manager.CSV.Equipment.CheckAdvanceLevel(target)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查目标装备的强化材料是否匹配
	err = manager.CSV.Equipment.CheckAdvanceMaterials(target, materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取强化消耗的金币
	goldCost, err := manager.CSV.Equipment.AdvanceCostGold(target)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 消耗
	costGoldReward := common.NewReward(static.CommonResourceTypeMoney, goldCost)
	costRewards := common.NewRewards()
	costRewards.AddReward(costGoldReward)

	// 检查金币是否足够
	err = u.CheckRewardsEnough(costRewards)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 进阶
	target.Stage.Plus(1)

	// 扣除金币
	reason := logreason.NewReason(logreason.EquipmentAdvance)
	err = u.CostRewards(costRewards, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 销毁装备材料
	u.EquipmentPack.BatchDestroy(materialIDs)

	u.BIEquipmentOp(target, bilog.EquipmentOpStageUp, logreason.EmptyReason())

	return target, nil
}

// 装备强化
func (u *User) EquipmentStrengthen(targetID int64, expItems []int32, materialIDs []int64) (*common.Equipment, error) {
	err := u.CheckActionUnlock(static.ActionIdTypeEquipgrowthunlock)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查材料，不能吃自己
	err = common.CheckEquipmentMaterial(targetID, materialIDs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取没有上锁材料装备信息
	materials, err := u.EquipmentPack.BatchGetUnlocked(materialIDs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查材料是否被装备
	err = common.CheckEquipmentWear(materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取目标装备信息
	target, err := u.EquipmentPack.Get(targetID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	maxLevel, err := manager.CSV.TeamLevelCache.GetEquipmentMaxLevel(u.Info.Level.Value())
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 获取经验药水的经验
	itemEXP := manager.CSV.Item.CalculateTotalEquipmentEXP(expItems)

	// 获取强化增加的经验
	addEXP, err := manager.CSV.Equipment.StrengthenEXP(target, itemEXP, materials, maxLevel)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 计算金币消耗
	goldCost, err := manager.CSV.Equipment.StrengthenGoldCost(target, itemEXP, materials)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查金币是否足够
	// 消耗
	costRewards := common.NewRewards()

	costRewards.AddReward(common.NewReward(static.CommonResourceTypeMoney, goldCost))

	for i, v := range common.EquipmentEXPItems {
		if expItems[i] > 0 {
			costRewards.AddReward(common.NewReward(v, expItems[i]))
		}
	}

	// 检查金币和药水是否足够
	err = u.CheckRewardsEnough(costRewards)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 增加经验
	target.EXP.Plus(addEXP)

	// 同步等级并解锁随机属性
	oldLevel := target.Level.Value()
	err = manager.CSV.Equipment.SyncLevelAndUnlockAttr(target)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 扣除金币
	reason := logreason.NewReason(logreason.EquipmentStrengthen)
	err = u.CostRewards(costRewards, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 销毁装备
	u.EquipmentPack.BatchDestroy(materialIDs)

	newLevel := target.Level.Value()
	if oldLevel < newLevel {
		u.TriggerQuestUpdate(static.TaskTypeEquipmentLevelUpCount, int32(target.Rarity), oldLevel, newLevel)
		u.TriggerQuestUpdate(static.TaskTypeEquipmentLevelCount, int32(target.Rarity), oldLevel, newLevel)
	}

	u.TriggerQuestUpdate(static.TaskTypeEquipmentExpUpTimes, int32(target.Rarity))

	u.BIEquipmentOp(target, bilog.EquipmentOpLevelUp, logreason.EmptyReason())

	return target, nil
}

// 装备加锁解锁
func (u *User) EquipmentLock(id int64, lock bool) error {
	// 获取目标装备信息
	equipment, err := u.EquipmentPack.Get(id)
	if err != nil {
		return errors.WrapTrace(err)
	}

	equipment.IsLocked = lock

	return nil
}

// 装备重铸阵营
func (u *User) EquipmentRecastCamp(id int64) (int8, error) {
	err := u.CheckActionUnlock(static.ActionIdTypeEquiprecastunlock)
	if err != nil {
		return 0, err
	}

	// 获取目标装备信息
	equipment, err := u.EquipmentPack.Get(id)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	// 检查是否有重铸过
	if !equipment.IsLastRecastCampConfirmed() {
		return 0, common.ErrEquipmentHasNotConfirm
	}

	// 重铸消耗
	costRewards, err := manager.CSV.Equipment.RecastCost(equipment)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	// 判断消耗是否足够
	err = u.CheckRewardsEnough(costRewards)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	// 检查装备是否可以重铸并且获取重铸结果
	camp, err := manager.CSV.Equipment.RecastCamp(equipment)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	// 保存重铸结果
	equipment.SaveRecastCamp(camp)

	// 扣除消耗
	reason := logreason.NewReason(logreason.EquipmentRecastCamp)
	err = u.CostRewards(costRewards, reason)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	return camp, nil
}

// 装备重铸阵营确认选择
func (u *User) EquipmentConfirmRecastCamp(id int64, confirm bool) (*common.Equipment, error) {
	err := u.CheckActionUnlock(static.ActionIdTypeEquiprecastunlock)
	if err != nil {
		return nil, err
	}

	// 获取目标装备信息
	equipment, err := u.EquipmentPack.Get(id)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否有重铸过
	if equipment.IsLastRecastCampConfirmed() {
		return nil, errors.Swrapf(common.ErrEquipmentNotRecast, id)
	}

	// 确认重铸结果
	equipment.ConfirmRecastCamp(confirm)

	return equipment, nil
}
