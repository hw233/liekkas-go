package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/statistic/logreason"
	"shared/utility/glog"
	"shared/utility/servertime"
)

func (u *User) ActivityPreSetting() {
	allActivity := manager.CSV.Activities.GetAllActivities()
	now := servertime.Now().Unix()

	for activityId, activityCfg := range allActivity {
		var activityStartTime int64 = 0
		var activityEndTime int64 = 0

		switch activityCfg.TimeType {
		case static.ActivityTimeTypeDate:
			activityStartTime = activityCfg.StartTimeArg
			activityEndTime = activityCfg.EndTimeArg
		case static.ActivityTimeTypeUserCreate:
			activityStartTime = u.Info.RegisterAt + activityCfg.StartTimeArg
			activityEndTime = u.Info.RegisterAt + activityCfg.EndTimeArg
		}

		if now >= activityEndTime {
			continue
		}

		funcIds, err := manager.CSV.Activities.GetActivityFuncIdsByActivity(activityId)
		if err != nil {
			glog.Error(err.Error())
			continue
		}

		for funcId := range funcIds {
			activityFunc, has := u.ActivityInfo.GetActivityFunc(funcId)
			if has && activityFunc.IsEnded() {
				continue
			}

			funcCfg, err := manager.CSV.Activities.GetActivityFunc(funcId)
			if err != nil {
				glog.Error(err.Error())
				continue
			}

			startTime := activityStartTime + funcCfg.StartTimeSecOffset
			endTime := activityEndTime + funcCfg.EndTimeSecOffset

			if !has {
				u.ActivityInfo.AddActivityFunc(funcId, startTime, endTime)
			} else if !activityFunc.IsEnded() {
				activityFunc.ResetTime(startTime, endTime)
			}
		}
	}

	u.UpdateActivities(now)
}

func (u *User) UpdateActivities(timestamp int64) {
	for _, activitFuncId := range u.ActivityInfo.GetWaitingActivityFuncs() {
		if u.checkActivityFuncStart(activitFuncId, timestamp) {
			u.startActivityFunc(activitFuncId, timestamp)
		}
	}

	for _, activitFuncId := range u.ActivityInfo.GetProgressingActivityFuncs() {
		if u.checkActivityFuncEnd(activitFuncId, timestamp) {
			u.endActivityFunc(activitFuncId, timestamp)
		}
	}
}

func (u *User) checkActivityFuncStart(activityFuncId int32, timestamp int64) bool {
	activityFunc, ok := u.ActivityInfo.GetActivityFunc(activityFuncId)
	if !ok {
		return false
	}

	if !activityFunc.IsWaiting() {
		return false
	}

	if !activityFunc.IsArrivalStartTime(timestamp) {
		return false
	}

	activityCfg, err := manager.CSV.Activities.GetActivityByFuncId(activityFuncId)
	if err != nil {
		glog.Error(err.Error())
		return false
	}

	if u.CheckUserConditions(activityCfg.UnlockConditions) == nil {
		return true
	}

	return true
}

func (u *User) checkActivityFuncEnd(activityFuncId int32, timestamp int64) bool {
	activityFunc, ok := u.ActivityInfo.GetActivityFunc(activityFuncId)
	if !ok {
		return false
	}

	if !activityFunc.IsProgressing() && !activityFunc.IsWaiting() {
		return false
	}

	activityCfg, err := manager.CSV.Activities.GetActivityByFuncId(activityFuncId)
	if err != nil {
		glog.Error(err.Error())
		return false
	}

	if !activityCfg.CloseConditions.Empty() {
		if u.CheckUserConditions(activityCfg.CloseConditions) == nil {
			return true
		}
	}

	if !activityFunc.IsArrivalEndTime(timestamp) {
		return false
	}

	return true
}

// activity effect
func (u *User) startActivityFunc(activityFuncId int32, timestamp int64) {
	funcCfg, err := manager.CSV.Activities.GetActivityFunc(activityFuncId)
	if err != nil {
		glog.Error(err.Error())
		return
	}

	switch funcCfg.FuncType {
	case static.ActivityFuncTypeScorePass:
		seasonId := funcCfg.FuncArgs[0]
		u.StartScorePass(seasonId, timestamp)
	}

	u.ActivityInfo.StartActivityFunc(activityFuncId, timestamp)
}

func (u *User) endActivityFunc(activityFuncId int32, timestamp int64) {
	funcCfg, err := manager.CSV.Activities.GetActivityFunc(activityFuncId)
	if err != nil {
		glog.Error(err.Error())
		return
	}

	switch funcCfg.FuncType {
	case static.ActivityFuncTypeScorePass:
		seasonId := funcCfg.FuncArgs[0]
		u.EndScorePass(seasonId)
	}

	cost := common.NewRewards()
	for _, itemId := range funcCfg.EndItemRemove {
		itemCount := u.ItemPack.Count(itemId)
		if itemCount > 0 {
			cost.AddReward(common.NewReward(itemId, itemCount))
		}
	}

	reason := logreason.NewReason(logreason.ActiviyEnd)
	u.CostRewards(cost, reason)

	u.ActivityInfo.EndActivityFunc(activityFuncId, timestamp)
}
