package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type StaminaEntry struct {
	sync.RWMutex

	StaminaCosts []*StaminaCost
	MaxTimes     int32
}

type StaminaCost struct {
	Id    int32
	Count int32
	Cost  int32
}

func NewStaminaEntry() *StaminaEntry {
	return &StaminaEntry{}
}

func (se *StaminaEntry) Check(config *Config) error {
	return nil
}

func (se *StaminaEntry) Reload(config *Config) error {
	se.Lock()
	defer se.Unlock()

	costs := make([]*StaminaCost, len(config.CfgQuickPurchaseStaminaConfig.GetAllData()))

	for _, staminaCsv := range config.CfgQuickPurchaseStaminaConfig.GetAllData() {
		staminaCfg := &StaminaCost{}
		err := transfer.Transfer(staminaCsv, staminaCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		//costs = append(costs, staminaCfg)
		if staminaCfg.Id < 1 {
			return errors.WrapTrace(err)
		}
		costs[staminaCfg.Id-1] = staminaCfg
	}

	se.StaminaCosts = costs
	se.MaxTimes = int32(len(costs))

	return nil
}

func (se *StaminaEntry) GetMaxTimes() int32 {
	se.RLock()
	defer se.RUnlock()

	return se.MaxTimes
}

func (se *StaminaEntry) GetStaminaData(id int32) (*StaminaCost, error) {
	se.RLock()
	defer se.RUnlock()

	if int(id) > len(se.StaminaCosts) {
		return nil, errors.Swrapf(common.ErrQuickPurchaseStaminaIndexOutOfRangeForData, id)
	}

	if int(id) == len(se.StaminaCosts) {
		return se.StaminaCosts[id-1], nil
	}
	return se.StaminaCosts[id], nil
}
