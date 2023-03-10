package entry

import (
	"shared/common"
	"shared/csv/base"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sort"
	"sync"
)

const CfgNameGuildLevel = "cfg_guide_level"

type GuildEntry struct {
	sync.RWMutex
	levels                    []GuildLevel
	exps                      []int32
	guildContributionExpRatio int32
	tasks                     map[int32]GuildTaskConfig
}
type GuildLevel struct {
	Id                              int32
	Exp                             int32
	SingleHelpAddIntimacy           int32
	SingleHelpAddGold               int32
	SingleHelpAddActivation         int32
	DailyGoldUpperLimitByHelp       int32
	DailyActivationUpperLimitByHelp int32
}

type GuildTaskConfig struct {
	Id         int32
	Count      int32
	Activation int32
	DropId     int32
	IsSeparate bool
}

func NewGuildEntry() *GuildEntry {
	return &GuildEntry{
		tasks: map[int32]GuildTaskConfig{},
	}
}

func (g *GuildEntry) Check(config *Config) error {
	return nil
}

func (g *GuildEntry) Reload(config *Config) error {
	g.Lock()
	defer g.Unlock()

	var levels []GuildLevel

	var exps []int32

	tasks := map[int32]GuildTaskConfig{}

	var cfgs []*base.CfgGuildLevel
	for _, v := range config.CfgGuildLevelConfig.GetAllData() {
		cfgs = append(cfgs, v)
	}
	less := func(i, j int) bool {
		return cfgs[i].Id < cfgs[j].Id
	}
	sort.Slice(cfgs, less)

	for _, v := range cfgs {
		guildLevel := &GuildLevel{}
		err := transfer.Transfer(v, guildLevel)
		if err != nil {
			return errors.WrapTrace(err)
		}
		levels = append(levels, *guildLevel)
	}

	for _, level := range levels {
		exps = append(exps, level.Exp)
	}

	for _, v := range config.CfgGuildTaskConfig.GetAllData() {
		task := &GuildTaskConfig{}
		err := transfer.Transfer(v, task)
		if err != nil {
			return errors.WrapTrace(err)
		}
		tasks[task.Id] = *task
	}

	g.levels = levels
	g.exps = exps
	g.tasks = tasks
	g.guildContributionExpRatio = config.GuildContributionExpRatio
	return nil
}

func (g *GuildEntry) GetExpArr() []int32 {
	g.RLock()
	defer g.RUnlock()
	return g.exps
}

func (g *GuildEntry) GetGuildLevel(lv int32) (*GuildLevel, error) {
	g.RLock()
	defer g.RUnlock()
	if lv <= 0 {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGuildLevel, lv)
	}
	if lv > int32(len(g.levels)) {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGuildLevel, lv)

	}
	return &g.levels[lv-1], nil
}

func (g *GuildEntry) GetGuildContributionExpRatio() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.guildContributionExpRatio
}

func (g *GuildEntry) GetAllTaskIds() ([]int32, []int32) {
	g.RLock()
	defer g.RUnlock()

	var normal []int32
	var separate []int32

	for id, task := range g.tasks {
		if task.IsSeparate {
			separate = append(separate, id)
		} else {
			normal = append(normal, id)
		}
	}
	return normal, separate
}

func (g *GuildEntry) GetTaskConfig(taskId int32) (*GuildTaskConfig, error) {
	g.RLock()
	defer g.RUnlock()

	task, ok := g.tasks[taskId]
	if !ok {
		return nil, errors.Swrapf(common.ErrGuildTaskConfigNotFound, taskId)
	}
	return &task, nil
}
