package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/servertime"
	"time"
)

func (u *User) StartScorePass(seasonId int32, startTime int64) {
	_, has := u.ScorePassInfo.GetSeason(seasonId)
	if has {
		return
	}

	refreshTm := DailyRefreshTime(time.Unix(startTime, 0))

	season := u.ScorePassInfo.StartSeason(seasonId, refreshTm.Unix())
	u.checkScoreSeasonUpdate(season, startTime)
	u.AddScorePassNotify(seasonId)
}

func (u *User) EndScorePass(seasonId int32) {
	season, has := u.ScorePassInfo.GetSeason(seasonId)
	if !has {
		return
	}

	if season.IsEnd() {
		return
	}

	for _, phase := range season.Phases {
		for groupId := range phase.Groups {
			groupCfg, err := manager.CSV.ScorePasses.GetGroup(groupId)
			if err != nil {
				glog.Error(err.Error())
				return
			}

			u.RemoveQuestByGroup(groupCfg.TaskGroupId)
		}
	}

	season.End()
}

func (u *User) ReceiveScorePassReward(rewardId int32) error {
	rewardCfg, err := manager.CSV.ScorePasses.GetReward(rewardId)
	if err != nil {
		return err
	}

	phaseId := rewardCfg.PhaseId
	phaseCfg, err := manager.CSV.ScorePasses.GetPhase(phaseId)
	if err != nil {
		return errors.Swrapf(common.ErrScorePassPhaseNotStart, phaseId)
	}

	seasonId := phaseCfg.SeasonId
	season, ok := u.ScorePassInfo.GetSeason(seasonId)
	if !ok || !season.IsProgressing() {
		return errors.Swrapf(common.ErrScorePassPhaseNotStart, phaseId)
	}

	phase, ok := season.GetPhase(phaseId)
	if !ok {
		return errors.Swrapf(common.ErrScorePassPhaseNotStart, phaseId)
	}

	err = u.CheckUserConditions(rewardCfg.ReceiveConditions)
	if err != nil {
		return err
	}

	reason := logreason.NewReason(logreason.ScorePassReward)
	u.AddRewardsByDropId(rewardCfg.DropId, reason)
	phase.RecordReward(rewardId)

	return nil
}

func (u *User) CheckScorePassUpdate(checkTime int64) {
	for _, season := range u.ScorePassInfo.Seasons {
		if !season.IsProgressing() {
			continue
		}

		u.checkScoreSeasonUpdate(season, checkTime)
	}
}

func (u *User) checkScoreSeasonUpdate(season *ScorePassSeason, checkTime int64) {
	if !season.IsProgressing() {
		return
	}

	phaseIds, err := manager.CSV.ScorePasses.GetPhaseIdsBySeason(season.Id)
	if err != nil {
		glog.Error(err.Error())
		return
	}

	update := false
	for _, phaseId := range phaseIds {
		phase, has := season.GetPhase(phaseId)
		if has {
			u.checkScorePassPhaseUpdate(phase, checkTime)
			continue
		}

		phaseCfg, err := manager.CSV.ScorePasses.GetPhase(phaseId)
		if err != nil {
			glog.Error(err.Error())
			continue
		}

		startTime := season.GetStartTime() + int64(phaseCfg.StartDayOffset)*servertime.SecondPerDay
		if checkTime >= startTime {
			u.startScorePassPhase(season, phaseId, startTime)
			update = true
		}
	}

	if update {
		u.AddScorePassNotify(season.Id)
	}
}

func (u *User) startScorePassPhase(season *ScorePassSeason, phaseId int32, startTime int64) {
	phase := season.StartPhase(phaseId, startTime)
	u.checkScorePassPhaseUpdate(phase, startTime)
}

func (u *User) checkScorePassPhaseUpdate(phase *ScorePassPhase, checkTime int64) {
	phaseId := phase.Id
	groupIds, err := manager.CSV.ScorePasses.GetGroupIdsByPhaseId(phaseId)
	if err != nil {
		glog.Error(err.Error())
		return
	}

	for groupId := range groupIds {
		if phase.IsGroupStarted(groupId) {
			continue
		}

		groupCfg, err := manager.CSV.ScorePasses.GetGroup(groupId)
		if err != nil {
			glog.Error(err.Error())
			return
		}

		startTime := phase.GetStartTime() + int64(groupCfg.StartDayOffset)*servertime.SecondPerDay
		if checkTime >= startTime {
			u.startScorePassGroup(phase, groupCfg)
		}
	}

}

func (u *User) startScorePassGroup(phase *ScorePassPhase, groupCfg *entry.ScorePassGroup) {
	phase.StartGroup(groupCfg.Id)

	questGroupId := groupCfg.TaskGroupId
	if questGroupId <= 0 {
		return
	}

	u.TryAcceptQuestByGroup(questGroupId)
}
