package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/number"
)

func (u *User) ActiveManual(manualType, relatedId int32) {
	manualId, ok := manager.CSV.Manual.FindManualIdByRelatedId(manualType, relatedId)
	// 配置不存在
	if !ok {
		return
	}

	if u.ManualInfo.IsActive(manualId) {
		return
	}

	u.ManualInfo.Active(manualId)

	u.TriggerQuestUpdate(static.TaskTypeManualCollect, manualId)
}

func (u *User) ManualGetReward(manualIds *number.NonRepeatableArr) ([]*pb.VOManualInfo, error) {

	err := u.CheckActionUnlock(static.ActionIdTypeManualunlock)
	if err != nil {
		return nil, err
	}

	// 判断是否获得和是否已领奖
	for _, manualId := range manualIds.Values() {
		rewarded, err := u.ManualInfo.IsRewarded(manualId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		if rewarded {
			return nil, errors.WrapTrace(common.ErrManualRewarded)
		}
	}
	vos := make([]*pb.VOManualInfo, 0, len(*manualIds))
	// 获得掉落
	dropIds, err := manager.CSV.Manual.GetManualDrops(manualIds)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 发奖
	reason := logreason.NewReason(logreason.Manual)
	_, err = u.AddRewardsByDropIds(dropIds, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 设置已领奖
	u.ManualInfo.SetRewarded(manualIds)
	u.Info.ManualCount += int32(len(manualIds.Values()))

	// manualInfo变动
	for _, manualId := range manualIds.Values() {
		vos = append(vos, &pb.VOManualInfo{
			ManualId: manualId,
			Rewarded: true,
		})
	}
	return vos, nil

}
