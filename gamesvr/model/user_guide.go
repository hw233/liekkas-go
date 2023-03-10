package model

import (
	"gamesvr/manager"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
)

func (u *User) UserGuide(id int32) error {
	config, err := manager.CSV.Guide.GetConfigById(id)
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = u.CheckUserCompoundConditions(config.UnlockCondition)
	if err != nil {
		return errors.WrapTrace(err)
	}

	biCommon := u.CreateBICommon(bilog.EventNameGuide)

	// 兼容触发多次相同guide_module_id,但是触发多次只第一次发奖励
	if !u.Info.PassedGuideIds.Contains(id) {
		if config.DropId > 0 {
			reason := logreason.NewReason(logreason.Guide, logreason.AddRelateLog(biCommon.GetLogId()))
			_, err = u.AddRewardsByDropId(config.DropId, reason)
			if err != nil {
				return errors.WrapTrace(err)
			}
		}
		u.Info.PassedGuideIds.Append(id)
	}

	u.BIGuide(id, biCommon)

	return nil
}
