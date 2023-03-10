package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type ActionUnlockEntry struct {
	sync.RWMutex

	ActionUnlocks map[int32]*ActionUnlock
}

type ActionUnlock struct {
	Id              int32
	UnlockCondition *common.Conditions `rule:"conditions"`
}

func NewActionUnlockEntry() *ActionUnlockEntry {
	return &ActionUnlockEntry{
		ActionUnlocks: map[int32]*ActionUnlock{},
	}
}

func (au *ActionUnlockEntry) Check(config *Config) error {
	return nil
}

func (an *ActionUnlockEntry) Reload(config *Config) error {
	an.Lock()
	defer an.Unlock()

	actionUnlocks := map[int32]*ActionUnlock{}

	for _, actionUnlockCsv := range config.CfgActionUnlockConfig.GetAllData() {
		actionUnlockCfg := &ActionUnlock{}

		err := transfer.Transfer(actionUnlockCsv, actionUnlockCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		actionUnlocks[actionUnlockCsv.Id] = actionUnlockCfg
	}

	an.ActionUnlocks = actionUnlocks

	return nil
}

func (an *ActionUnlockEntry) GetUnlockConditions(id int32) (*ActionUnlock, error) {
	an.RLock()
	defer an.RUnlock()

	actionUnlock, ok := an.ActionUnlocks[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrActionUnlockWrongActionID, id)
	}

	return actionUnlock, nil
}
