package entry

import (
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type ScorePassPhase struct {
	Id             int32
	SeasonId       int32
	StartDayOffset int32
	ScoreItemId    int32
}

type ScorePassGroup struct {
	Id             int32
	PhaseId        int32
	StartDayOffset int32
	TaskGroupId    int32
}

type ScorePassReward struct {
	Id                int32
	PhaseId           int32
	Score             int32
	DropId            int32
	ReceiveConditions *common.Conditions `ignore:"true"`
}

type ScorePasses struct {
	sync.RWMutex

	ScorePassPhases  map[int32]*ScorePassPhase
	ScorePassGroups  map[int32]*ScorePassGroup
	ScorePassRewards map[int32]*ScorePassReward

	Seasons       map[int32]map[int32]int32
	PhaseToGroups map[int32]map[int32]int32
}

func NewScorePasses() *ScorePasses {
	return &ScorePasses{}
}

func (sp *ScorePasses) Check(config *Config) error {
	return nil
}

func (sp *ScorePasses) Reload(config *Config) error {
	sp.Lock()
	defer sp.Unlock()

	phases := map[int32]*ScorePassPhase{}
	groups := map[int32]*ScorePassGroup{}
	rewards := map[int32]*ScorePassReward{}

	seasons := map[int32]map[int32]int32{}
	phaseToGroups := map[int32]map[int32]int32{}

	for _, csv := range config.CfgScorePassPhaseConfig.GetAllData() {
		phaseCfg := &ScorePassPhase{}

		err := transfer.Transfer(csv, phaseCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		phaseId := phaseCfg.Id
		phases[phaseId] = phaseCfg

		seasonId := phaseCfg.SeasonId
		seasonPhases, ok := seasons[seasonId]
		if !ok {
			seasonPhases = map[int32]int32{}
			seasons[seasonId] = seasonPhases
		}

		seasonPhases[phaseId] = phaseId
	}

	for _, csv := range config.CfgScorePassGroupConfig.GetAllData() {
		groupCfg := &ScorePassGroup{}

		err := transfer.Transfer(csv, groupCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		groupId := groupCfg.Id
		phaseId := groupCfg.PhaseId
		groupIds, ok := phaseToGroups[phaseId]
		if !ok {
			groupIds = map[int32]int32{}
			phaseToGroups[phaseId] = groupIds
		}
		groupIds[groupId] = groupId

		groups[groupId] = groupCfg
	}

	for _, csv := range config.CfgScorePassRewardConfig.GetAllData() {
		rewardCfg := &ScorePassReward{
			ReceiveConditions: common.NewConditions(),
		}

		err := transfer.Transfer(csv, rewardCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		phaseCfg, ok := phases[rewardCfg.PhaseId]
		if ok {
			itemId := phaseCfg.ScoreItemId
			count := rewardCfg.Score

			condition := common.NewCondition(static.ConditionTypeHasItem, itemId, count)
			rewardCfg.ReceiveConditions.AddCondition(condition)
		}

		rewards[rewardCfg.Id] = rewardCfg
	}

	sp.ScorePassPhases = phases
	sp.ScorePassGroups = groups
	sp.ScorePassRewards = rewards
	sp.Seasons = seasons
	sp.PhaseToGroups = phaseToGroups

	return nil
}

func (sp *ScorePasses) GetPhase(id int32) (*ScorePassPhase, error) {
	phase, ok := sp.ScorePassPhases[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrScorePassPhaseCfgNotFound, id)
	}

	return phase, nil
}

func (sp *ScorePasses) GetGroup(id int32) (*ScorePassGroup, error) {
	group, ok := sp.ScorePassGroups[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrScorePassGroupCfgNotFound, id)
	}

	return group, nil
}

func (sp *ScorePasses) GetReward(id int32) (*ScorePassReward, error) {
	reward, ok := sp.ScorePassRewards[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrScorePassRewardCfgNotFound, id)
	}

	return reward, nil
}

func (sp *ScorePasses) GetPhaseIdsBySeason(id int32) (map[int32]int32, error) {
	phaseIds, ok := sp.Seasons[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrScorePassNoSuchSeason, id)
	}

	return phaseIds, nil
}

func (sp *ScorePasses) GetGroupIdsByPhaseId(id int32) (map[int32]int32, error) {
	groupIds, ok := sp.PhaseToGroups[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrScorePassPhaseCfgNotFound, id)
	}

	return groupIds, nil
}
