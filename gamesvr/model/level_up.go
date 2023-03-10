package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/statistic/logreason"
	"shared/utility/glog"
)

type LevelExp interface {
	GetLevel() int32
	SetLevel(level int32, u *User)
	GetExp() int32
	SetExp(int32)
	GetMaxLevel(*User) int32
}

func (u *User) AddExpCommon(addNum int32, lvExp LevelExp, cache entry.LevelCfgCache) {
	if addNum <= 0 {
		return
	}

	expArr := cache.GetExpArr()
	nowLevel := lvExp.GetLevel()
	maxLevel := lvExp.GetMaxLevel(u)
	expLen := int32(len(expArr))
	if maxLevel > expLen {
		maxLevel = expLen
	}
	if nowLevel >= maxLevel {
		return
	}

	nowExp := lvExp.GetExp() + addNum

	for i := nowLevel - 1; i < maxLevel; i++ {
		exp := expArr[i]
		if exp == nowExp {
			nowLevel = i + 1
			break
		} else if exp > nowExp {
			break
		}
		nowLevel = i + 1
	}

	if nowLevel >= maxLevel {
		// 满级后经验就不会再加了,停在满级0经验
		lvExp.SetExp(expArr[maxLevel-1])
		lvExp.SetLevel(maxLevel, u)

	} else {
		lvExp.SetExp(nowExp)
		lvExp.SetLevel(nowLevel, u)

	}

}

type UserExpLevelCalculator struct {
	UserInfo *UserInfo
	Reason   *logreason.Reason
}

func (uelc *UserExpLevelCalculator) GetLevel() int32 {
	return uelc.UserInfo.Level.Value()
}

func (uelc *UserExpLevelCalculator) SetLevel(level int32, user *User) {
	if uelc.UserInfo.Level.Value() < level {
		for i := uelc.UserInfo.Level.Value() + 1; i <= level; i++ {
			levelConfig, ok := manager.CSV.TeamLevelCache.GetByLv(i)
			if ok {
				reward := common.NewReward(static.CommonResourceTypeEnergy, levelConfig.StaminaRecover)

				reason := logreason.NewReason(logreason.LevelUp)
				_, _, err := user.addReward(reward, reason)
				if err != nil {
					glog.Errorf("UserInfo SetLevel, addReward err:%v", err)
				}
			}
		}
	}
	oldLevel := uelc.UserInfo.Level.Value()
	uelc.UserInfo.Level.SetValue(level, uelc.Reason)
	levelConfig, ok := manager.CSV.TeamLevelCache.GetByLv(uelc.UserInfo.Level.Value())
	if ok {
		uelc.UserInfo.Energy.SetTimerUpper(levelConfig.MaxStamina)
		uelc.UserInfo.Ap.SetTimerUpper(levelConfig.MaxAp)
	}

	user.TriggerQuestUpdate(static.TaskTypeAccountLevel, level)
	user.QuestCheckUnlock(static.ConditionTypeUserLevel, oldLevel, level)

}

func (uelc *UserExpLevelCalculator) GetExp() int32 {
	return uelc.UserInfo.Exp.Value()
}

func (uelc *UserExpLevelCalculator) SetExp(exp int32) {
	uelc.UserInfo.Exp.SetValue(exp, uelc.Reason)
}

func (uelc *UserExpLevelCalculator) GetMaxLevel(*User) int32 {
	return manager.CSV.TeamLevelCache.GetMaxLv()
}
