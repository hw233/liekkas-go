package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

const CfgNameVn = "cfg_vn"

type VnEntry struct {
	sync.RWMutex
	vnConfigs map[int32]VnConfig
}

type VnConfig struct {
	Id     int32
	Reward *common.Rewards `src:"Reward1" rule:"rewards"`
}

func NewVnEntry() *VnEntry {
	return &VnEntry{
		vnConfigs: map[int32]VnConfig{},
	}
}

func (v *VnEntry) Check(config *Config) error {
	return nil
}

func (v *VnEntry) Reload(config *Config) error {
	v.Lock()
	defer v.Unlock()
	vnConfigs := map[int32]VnConfig{}
	for _, cfg := range config.CfgVnConfig.GetAllData() {
		vnConfig := &VnConfig{}
		err := transfer.Transfer(cfg, vnConfig)
		if err != nil {
			return errors.WrapTrace(err)
		}
		vnConfigs[cfg.Id] = *vnConfig
	}

	v.vnConfigs = vnConfigs
	return nil
}

func (v *VnEntry) GetConfigById(vnId int32) (*VnConfig, error) {
	v.RLock()
	defer v.RUnlock()
	config, ok := v.vnConfigs[vnId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameVn, vnId)
	}
	return &config, nil
}
