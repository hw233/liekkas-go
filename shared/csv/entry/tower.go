package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type Tower struct {
	Id         int32
	Camp       []int32
	GoUpLimit  int32
	ActiveDate []int32
}

type TowerStage struct {
	Id      int32
	TowerId int32
	Stage   int32
	LevelId int32
}

type TowerEntry struct {
	sync.RWMutex

	Towers      map[int32]*Tower
	TowerStages map[int32]map[int32]*TowerStage
}

func NewTowerEntry() *TowerEntry {
	return &TowerEntry{
		Towers:      map[int32]*Tower{},
		TowerStages: map[int32]map[int32]*TowerStage{},
	}
}

func (te *TowerEntry) Check(config *Config) error {

	return nil
}

func (te *TowerEntry) Reload(config *Config) error {
	te.Lock()
	defer te.Unlock()

	for _, towerCSV := range config.CfgTowerConfig.GetAllData() {
		towerCfg := &Tower{}

		err := transfer.Transfer(towerCSV, towerCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		te.Towers[towerCfg.Id] = towerCfg
	}

	for _, stageCSV := range config.CfgTowerStageConfig.GetAllData() {
		stageCfg := &TowerStage{}

		err := transfer.Transfer(stageCSV, stageCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		towerStages, ok := te.TowerStages[stageCfg.TowerId]
		if !ok {
			towerStages = make(map[int32]*TowerStage)
			te.TowerStages[stageCfg.TowerId] = towerStages
		}

		towerStages[stageCfg.Stage] = stageCfg
	}

	return nil
}

func (te *TowerEntry) GetTower(id int32) (*Tower, error) {
	tower, ok := te.Towers[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrTowerConfigNotFound, id)
	}

	return tower, nil
}

func (te *TowerEntry) GetTowerStage(towerId, stage int32) (*TowerStage, error) {
	TowerStages, ok := te.TowerStages[towerId]
	if !ok {
		return nil, errors.Swrapf(common.ErrTowerStageConfigNotFound, towerId, stage)
	}

	stageCfg, ok := TowerStages[stage]
	if !ok {
		return nil, errors.Swrapf(common.ErrTowerStageConfigNotFound, towerId, stage)
	}

	return stageCfg, nil
}
