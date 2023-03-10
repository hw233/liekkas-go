package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/rand"
)

/**
使用道具
*/
var (
	itemUses = map[int32]itemUse{}
)

type itemUse func(u *User, itemId, useNum, param int32, reason *logreason.Reason) error

func init() {
	itemUses[static.ItemTypeGraveyardGetProduct] = GraveyardGetProductByItem
	itemUses[static.ItemTypeGiftSelectOne] = GiftSelectOne
	itemUses[static.ItemTypeGiftRandomDrop] = GiftRandom
	itemUses[static.ItemTypeEnergyItem] = EnergyItemUse
}

func (u *User) ItemUse(reward *common.Reward, param int32) error {
	use, ok := itemUses[reward.Type]
	if !ok {
		return errors.Swrapf(common.ErrItemTypeCannotUse, reward.ID)
	}
	consume := common.NewRewards()
	consume.AddReward(reward)
	// 检查消耗
	err := u.CheckRewardsEnough(consume)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 使用道具
	reason := logreason.NewReason(logreason.ItemUse)
	err = use(u, reward.ID, reward.Num, param, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 消耗道具
	return u.CostRewards(consume, reason)

}

func GraveyardGetProductByItem(u *User, itemId, useNum, param int32, reason *logreason.Reason) error {
	item, ok := manager.CSV.Item.GetGraveyardGetProductItem(itemId)
	if !ok {
		return errors.WrapTrace(common.ErrParamError)
	}
	total := common.NewRewards()
	builds := u.Graveyard.GetBuildsByBuildId(item.BuildId)
	if len(builds) == 0 {
		return errors.WrapTrace(common.ErrParamError)
	}

	for _, build := range builds {
		// 没有建完成的没有产出
		if !build.BuildComplete() {
			continue
		}
	}
	for ; useNum > 0; useNum-- {
		// 所有该种建筑的sec秒内产出
		once, err := manager.CSV.GraveyardGetBuildsProductByItem(builds, item)
		if err != nil {
			return errors.WrapTrace(err)
		}
		total.AddRewards(once)
	}

	u.TriggerQuestUpdate(static.TaskTypeGraveyarAccelerateTimes, 0, 1)
	u.TriggerQuestUpdate(static.TaskTypeGraveyarAccelerateTime, 0, item.Sec*useNum)

	_, err := u.addRewards(total, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil

}

func GiftSelectOne(u *User, itemId, useNum, param int32, reason *logreason.Reason) error {
	item, ok := manager.CSV.Item.GetGiftSelectOne(itemId)
	if !ok {
		return errors.WrapTrace(common.ErrParamError)
	}
	if param < 0 || int(param) > len(item.Rewards)-1 {
		return common.ErrParamError
	}
	rewards := common.NewRewards()
	for ; useNum > 0; useNum-- {
		rewards.AddReward(&item.Rewards[param])
	}

	_, err := u.addRewards(rewards, reason)

	return err
}

func GiftRandom(u *User, itemId, useNum, param int32, reason *logreason.Reason) error {
	item, ok := manager.CSV.Item.GetGiftRandomDrop(itemId)
	if !ok {
		return errors.WrapTrace(common.ErrParamError)
	}

	dropIds := make([]int32, 0, useNum)
	for ; useNum > 0; useNum-- {
		index := rand.RangeInt(0, len(item.DropIds)-1)
		dropIds = append(dropIds, item.DropIds[index])
	}

	_, err := u.AddRewardsByDropIds(dropIds, reason)
	return err
}

func EnergyItemUse(u *User, itemId, useNum, param int32, reason *logreason.Reason) error {
	item, ok := manager.CSV.Item.GetEnergyItem(itemId)
	if !ok {
		return errors.WrapTrace(common.ErrParamError)
	}
	rewards := common.NewRewards()
	rewards.AddReward(common.NewReward(static.CommonResourceTypeEnergy, item.AddNum))
	rewards.Multiple(useNum)

	_, err := u.addRewards(rewards, reason)
	return err
}
