package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/statistic/logreason"
	"shared/utility/errors"
)

func (u *User) VnGetReward(vnId int32) error {
	contains := u.Info.RewardedVnIds.Contains(vnId)
	if contains {
		return errors.WrapTrace(common.ErrVnRewarded)
	}
	config, err := manager.CSV.VN.GetConfigById(vnId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if config.Reward.IsEmpty() {
		return errors.WrapTrace(common.ErrVnRewardIsEmpty)
	}
	u.Info.RewardedVnIds.Append(vnId)

	reason := logreason.NewReason(logreason.VisualNovel)
	_, err = u.addRewards(config.Reward, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (u *User) VnRead(ctx context.Context, vnId int32) {
	// 任务监听
	u.Yggdrasil.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeVn, 1, vnId)
}
