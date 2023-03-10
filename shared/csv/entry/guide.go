package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

const CfgNameGuide = "cfg_guide_module"

type GuideEntry struct {
	sync.RWMutex
	guideConfigs map[int32]GuideConfig
}
type GuideConfig struct {
	Id              int32
	UnlockCondition *common.CompoundConditions `rule:"compoundConditions" src:"UnlockCondition,UnlockConditionOr"`
	GuideOrder      int32
	DropId          int32
	Desc            string
}

func NewGuideEntry() *GuideEntry {
	return &GuideEntry{}
}

func (g *GuideEntry) Check(config *Config) error {
	return nil
}

func (g *GuideEntry) Reload(config *Config) error {
	g.Lock()
	defer g.Unlock()
	guideConfigs := map[int32]GuideConfig{}
	for _, cfg := range config.CfgGuideModuleConfig.GetAllData() {
		guideConfig := &GuideConfig{}
		err := transfer.Transfer(cfg, guideConfig)
		if err != nil {
			return errors.WrapTrace(err)
		}
		guideConfigs[cfg.Id] = *guideConfig
	}

	g.guideConfigs = guideConfigs
	return nil
}
func (g *GuideEntry) GetConfigById(guideId int32) (*GuideConfig, error) {
	g.RLock()
	defer g.RUnlock()
	config, ok := g.guideConfigs[guideId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGuide, guideId)
	}
	return &config, nil
}

func (g *GuideEntry) GetGuideBefore(order int32) []GuideConfig {
	g.RLock()
	defer g.RUnlock()

	ret := make([]GuideConfig, 0, len(g.guideConfigs))
	for _, config := range g.guideConfigs {
		if config.GuideOrder <= order {
			ret = append(ret, config)
		}
	}

	return ret
}
