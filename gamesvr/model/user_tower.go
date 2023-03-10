package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/slice"
)

func (u *User) TowerDailyRefresh(refreshTime int64) {
	u.TowerInfo.ResetDaily()
}

func (u *User) checkTowerLevel(towerId, towerStage, LevelId int32, charas []int32) error {
	towerCfg, err := manager.CSV.Tower.GetTower(towerId)
	if err != nil {
		return err
	}

	todayRefreshTime := TodayRefreshTime()
	todayWeekDay := int32((todayRefreshTime.Weekday()+6)%7 + 1)

	actived := false
	for _, weekDay := range towerCfg.ActiveDate {
		if todayWeekDay == weekDay {
			actived = true
			break
		}
	}
	if !actived {
		return errors.Swrapf(common.ErrTowerNotActived, towerId)
	}

	towerStageCfg, err := manager.CSV.Tower.GetTowerStage(towerId, towerStage)
	if err != nil {
		return err
	}

	if towerStageCfg.LevelId != LevelId {
		return errors.Swrapf(common.ErrTowerInvalidLevel, towerId, LevelId)
	}

	tower, ok := u.TowerInfo.GetTower(towerId)
	if !ok {
		if towerStage != 1 {
			return errors.Swrapf(common.ErrTowerStageNotArrival, towerId, towerStage)
		}
		tower = NewTower(towerId)
	}

	nextStage := tower.GetCurStage() + 1
	if nextStage < towerStage {
		return errors.Swrapf(common.ErrTowerStageNotArrival, towerId, towerStage)
	}

	if nextStage == towerStage && towerCfg.GoUpLimit != 0 &&
		tower.GetTodayGoUpTimes() >= towerCfg.GoUpLimit {
		return errors.Swrapf(common.ErrTowerGoUpLimited, towerId)
	}

	if len(towerCfg.Camp) > 0 {
		for _, charaId := range charas {
			camp, err := manager.CSV.Character.Camp(charaId)
			if err != nil {
				return err
			}

			if !slice.SliceInt32HasEle(towerCfg.Camp, camp) {
				return errors.Swrapf(common.ErrTowerCharaCampLimited, towerId, charaId)
			}
		}
	}

	return nil
}

func (u *User) onTowerLevelPass(towerId, stage int32) {
	tower := u.TowerInfo.GetOrCreateTower(towerId)
	if stage > tower.GetCurStage() {
		tower.RecordGoUp()
	}

	u.AddTowerUpdateNotify(towerId)

	u.TriggerQuestUpdate(static.TaskTypeTowerStagePassed, towerId, stage)
}
