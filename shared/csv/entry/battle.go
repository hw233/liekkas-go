package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type BattleConfig struct {
	Id         int32
	Npc        []int32
	ControlNpc []int32
}

type BattleNpc struct {
	Id         int32
	CharaID    int32
	Level      int32
	Star       int32
	Stage      int32
	Equipments []*common.Equipment `ignore:"true"`
	SkillLv    map[int32]int32     `rule:"strsToSkills"`
	WorldItem  int32
}

type Battle struct {
	sync.RWMutex

	Battles    map[int32]*BattleConfig
	BattleNpcs map[int32]*BattleNpc
}

func NewBattle() *Battle {
	return &Battle{}
}

func (b *Battle) Check(config *Config) error {
	return nil
}

func (b *Battle) Reload(config *Config) error {
	b.Lock()
	defer b.Unlock()

	battles := map[int32]*BattleConfig{}
	battleNpcs := map[int32]*BattleNpc{}

	for _, csv := range config.CfgBattleLevelConfig.GetAllData() {
		cfg := &BattleConfig{}
		err := transfer.Transfer(csv, cfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		battles[cfg.Id] = cfg
	}

	for _, csv := range config.CfgBattleNpcConfig.GetAllData() {
		cfg := &BattleNpc{}
		err := transfer.Transfer(csv, cfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		cfg.Equipments = make([]*common.Equipment, 0, 4)

		appendEquipment := func(equipments []*common.Equipment, cfgEquip []int32) []*common.Equipment {
			if len(cfgEquip) > 0 {
				equipments = append(equipments, intsToEquipment(cfgEquip))
			}

			return equipments
		}

		cfg.Equipments = appendEquipment(cfg.Equipments, csv.Equip1)
		cfg.Equipments = appendEquipment(cfg.Equipments, csv.Equip2)
		cfg.Equipments = appendEquipment(cfg.Equipments, csv.Equip3)
		cfg.Equipments = appendEquipment(cfg.Equipments, csv.Equip4)

		battleNpcs[cfg.Id] = cfg
	}

	b.Battles = battles
	b.BattleNpcs = battleNpcs

	return nil
}

func (b *Battle) GetBattleConfig(id int32) (*BattleConfig, error) {
	cfg, ok := b.Battles[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrBattleCfgNotFound, id)
	}

	return cfg, nil
}

func (b *Battle) GetBattleNPC(id int32) (*BattleNpc, error) {
	cfg, ok := b.BattleNpcs[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrBattleNPCCfgNotFound, id)
	}

	return cfg, nil
}

func intsToEquipment(ints []int32) *common.Equipment {
	eid := ints[0]
	level := ints[1]
	stage := ints[2]
	camp := ints[3]

	equipment := common.NewEquipment(0, eid)
	equipment.Level.SetValue(level)
	equipment.Stage.SetValue(stage)
	equipment.Camp = int8(camp)

	return equipment
}
