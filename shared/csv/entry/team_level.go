package entry

import (
	"shared/csv/static"
	"sort"
	"sync"
	"time"

	"shared/common"
	"shared/csv/base"
	"shared/utility/errors"
	"shared/utility/transfer"
)

const (
	CfgTeamLevelConfig = "cfg_team_level"
)

type TeamLevelCache struct {
	sync.RWMutex

	exps                  []int32
	teamLevels            []TeamLevel
	changeNameConsume     *common.Rewards
	staminaRecoverSeconds time.Duration
}

func NewTeamLevelCache() *TeamLevelCache {
	return &TeamLevelCache{
		exps:                  []int32{},
		teamLevels:            []TeamLevel{},
		staminaRecoverSeconds: 300 * time.Second,
	}
}

type TeamLevel struct {
	CharLv           int32
	MaxStamina       int32
	StaminaRecover   int32
	MaxAp            int32
	ExploreAp        int32
	HeroMaxLevel     int32
	EquipMaxLevel    int32
	WorldItemLevelUp int32
}

type TeamLevelExp struct {
	LvExp int32
}

func NewTeamLevel() *TeamLevel {
	return &TeamLevel{}
}

func (t *TeamLevelCache) Reload(config *Config) error {
	t.Lock()
	defer t.Unlock()
	var exps []int32
	var teamLevels []TeamLevel
	var cfgs []*base.CfgTeamLevel
	for _, v := range config.CfgTeamLevelConfig.GetAllData() {
		cfgs = append(cfgs, v)
	}
	less := func(i, j int) bool {
		return cfgs[i].Id < cfgs[j].Id
	}
	sort.Slice(cfgs, less)

	var teamLevelExps []TeamLevelExp
	for _, v := range cfgs {
		teamLevel := NewTeamLevel()
		err := transfer.Transfer(v, teamLevel)
		if err != nil {
			return errors.WrapTrace(err)
		}
		teamLevels = append(teamLevels, *teamLevel)
		teamLevelExp := &TeamLevelExp{}
		err = transfer.Transfer(v, teamLevelExp)
		if err != nil {
			return errors.WrapTrace(err)
		}
		teamLevelExps = append(teamLevelExps, *teamLevelExp)
	}

	for _, teamLevelExp := range teamLevelExps {
		exps = append(exps, teamLevelExp.LvExp)
	}
	changeNameConsume := common.NewRewards()
	changeNameConsume.AddReward(common.NewReward(static.CommonResourceTypeDiamondGift, config.NicknameChangeConsume))
	t.exps = exps
	t.teamLevels = teamLevels
	t.staminaRecoverSeconds = time.Duration(config.StaminaRecoverSeconds) * time.Second
	t.changeNameConsume = changeNameConsume
	return nil
}

func (t *TeamLevelCache) Check(config *Config) error {
	var cfgs []*base.CfgTeamLevel
	for _, v := range config.CfgTeamLevelConfig.GetAllData() {
		cfgs = append(cfgs, v)
	}
	less := func(i, j int) bool {
		return cfgs[i].Id < cfgs[j].Id
	}
	sort.Slice(cfgs, less)

	for i, v := range cfgs {
		if int32(i+1) != v.Id {
			return errors.Swrapf(common.ErrNotFoundInCSV, "cfg_team_level", v.Id)
		}
	}
	return nil
}

func (t *TeamLevelCache) getByLv(lv int32) (*TeamLevel, bool) {
	index := lv - 1
	if int(index) >= len(t.teamLevels) {
		return nil, false
	}
	return &t.teamLevels[index], true
}

func (t *TeamLevelCache) GetByLv(lv int32) (*TeamLevel, bool) {
	t.RLock()
	defer t.RUnlock()

	index := lv - 1
	if int(index) >= len(t.teamLevels) {
		return nil, false
	}
	return &t.teamLevels[index], true
}

func (t *TeamLevelCache) GetMaxLv() int32 {
	t.RLock()
	defer t.RUnlock()
	return int32(len(t.teamLevels))
}

func (t *TeamLevelCache) GetExpArr() []int32 {
	t.RLock()
	defer t.RUnlock()
	return t.exps
}

func (t *TeamLevelCache) GetCharacterMaxLv(lv int32) (int32, error) {
	t.RLock()
	defer t.RUnlock()

	byLv, ok := t.getByLv(lv)
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgTeamLevelConfig, lv)
	}
	return byLv.CharLv, nil
}

func (t *TeamLevelCache) GetHeroMaxLevel(userLevel int32) (int32, error) {
	t.RLock()
	defer t.RUnlock()

	byLv, ok := t.getByLv(userLevel)
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgTeamLevelConfig, userLevel)
	}
	return byLv.HeroMaxLevel, nil
}

func (t *TeamLevelCache) GetEquipmentMaxLevel(userLevel int32) (int32, error) {
	t.RLock()
	defer t.RUnlock()

	byLv, ok := t.getByLv(userLevel)
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgTeamLevelConfig, userLevel)
	}
	return byLv.EquipMaxLevel, nil
}

func (t *TeamLevelCache) GetWorldItemLevelUp(userLevel int32) (int32, error) {
	t.RLock()
	defer t.RUnlock()

	byLv, ok := t.getByLv(userLevel)
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgTeamLevelConfig, userLevel)
	}
	return byLv.WorldItemLevelUp, nil
}

func (t *TeamLevelCache) GetStaminaRecoverSeconds() time.Duration {
	t.RLock()
	defer t.RUnlock()
	return t.staminaRecoverSeconds
}

func (t *TeamLevelCache) GetChangeNameConsume() *common.Rewards {
	t.RLock()
	defer t.RUnlock()
	return t.changeNameConsume
}
