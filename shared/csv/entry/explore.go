package entry

import (
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

const (
	CfgExploreNPC     = "cfg_explore_npc"
	CfgExploreGather  = "cfg_explore_gather"
	CfgExploreMonster = "cfg_explore_monster"
)

type ExploreMapObject struct {
	Id        int32
	ChapterId int32
	Type      int32
	TypeParam int32
}

type ExploreEventPoint struct {
	Id               int32
	EventType        int32
	EventParam       int32
	UnlockConditions *common.Conditions `rule:"conditions"`
	RefreshTime      int64              `rule:"int32ToInt64"`
	RefreshLimit     int32
}

type ExploreNPC struct {
	Id           int32
	OptionCost   []*common.Rewards `ignore:"true"`
	OptionDrops  []int32           `ignore:"true"`
	EventPointId int32             `ignore:"true"`
}

type ExploreRewardPoint struct {
	Id           int32
	GatherDrop   int32
	EventPointId int32 `ignore:"true"`
}

type ExploreMonster struct {
	Id           int32
	Monstermap   int32
	EventPointId int32 `ignore:"true"`
}

type ExploreFog struct {
	Id               int32
	UnlockConditions *common.Conditions `rule:"conditions"`
}

type ExploreResource struct {
	Id                int32
	Monstermap        int32
	CollectConditions *common.Conditions `ignore:"true"`
	Time              int64              `rule:"int32ToInt64"`
	Drop              int32
	RefreshLimit      int32
}

type ExploreTransportGate struct {
	Id               int32
	UnlockConditions *common.Conditions `rule:"conditions"`
	CompleteLimit    int32
}

type ExploreEntry struct {
	sync.RWMutex

	Chapters       map[int32]*Chapter
	ChapterRewards map[int32]*ChapterReward
	MapObjects     map[int32]*ExploreMapObject
	EventPoints    map[int32]*ExploreEventPoint
	NPCs           map[int32]*ExploreNPC
	RewardPoints   map[int32]*ExploreRewardPoint
	Monsters       map[int32]*ExploreMonster
	Fogs           map[int32]*ExploreFog
	Resources      map[int32]*ExploreResource
	TransportGates map[int32]*ExploreTransportGate
}

func NewExploreEntry() *ExploreEntry {
	return &ExploreEntry{
		Chapters:       map[int32]*Chapter{},
		ChapterRewards: map[int32]*ChapterReward{},
		MapObjects:     map[int32]*ExploreMapObject{},
		EventPoints:    map[int32]*ExploreEventPoint{},
		NPCs:           map[int32]*ExploreNPC{},
		RewardPoints:   map[int32]*ExploreRewardPoint{},
		Monsters:       map[int32]*ExploreMonster{},
		Fogs:           map[int32]*ExploreFog{},
		Resources:      map[int32]*ExploreResource{},
		TransportGates: map[int32]*ExploreTransportGate{},
	}
}

func (ee *ExploreEntry) Check(config *Config) error {
	return nil
}

func (ee *ExploreEntry) Reload(config *Config) error {
	ee.Lock()
	defer ee.Unlock()

	reloadList := []func(config *Config) error{
		ee.reloadMap,
		ee.reloadNPC,
		ee.reloadRewardPoint,
		ee.reloadMonster,
		ee.reloadEventPoint,
		ee.reloadFog,
		ee.reloadResource,
		ee.reloadTransportGate,
	}

	for _, reloadFunc := range reloadList {
		err := reloadFunc(config)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil
}

func (ee *ExploreEntry) reloadMap(config *Config) error {
	for _, mapObjectCSV := range config.CfgExploreMapCoordinateConfig.GetAllData() {
		mapObj := &ExploreMapObject{}

		err := transfer.Transfer(mapObjectCSV, mapObj)
		if err != nil {
			return errors.WrapTrace(err)
		}

		ee.MapObjects[mapObj.Id] = mapObj
	}

	return nil
}

func (ee *ExploreEntry) reloadNPC(config *Config) error {
	for _, npcCSV := range config.CfgExploreNpcConfig.GetAllData() {
		npc := &ExploreNPC{}
		err := transfer.Transfer(npcCSV, npc)
		if err != nil {
			return errors.WrapTrace(err)
		}

		npc.OptionCost = make([]*common.Rewards, 0)
		npc.OptionDrops = make([]int32, 0)

		if len(npcCSV.NpcOpt1) > 0 {
			cost, _ := StringsToRewards(npcCSV.NpcOptCost1)
			npc.OptionCost = append(npc.OptionCost, cost)
			npc.OptionDrops = append(npc.OptionDrops, npcCSV.NpcOptDropId1)
		}

		if len(npcCSV.NpcOpt2) > 0 {
			cost, _ := StringsToRewards(npcCSV.NpcOptCost2)
			npc.OptionCost = append(npc.OptionCost, cost)
			npc.OptionDrops = append(npc.OptionDrops, npcCSV.NpcOptDropId2)
		}

		if len(npcCSV.NpcOpt3) > 0 {
			cost, _ := StringsToRewards(npcCSV.NpcOptCost3)
			npc.OptionCost = append(npc.OptionCost, cost)
			npc.OptionDrops = append(npc.OptionDrops, npcCSV.NpcOptDropId3)
		}

		ee.NPCs[npc.Id] = npc
	}

	return nil
}

func (ee *ExploreEntry) reloadRewardPoint(config *Config) error {
	for _, rewardPointCSV := range config.CfgExploreGatherConfig.GetAllData() {
		rewardPoint := &ExploreRewardPoint{}
		err := transfer.Transfer(rewardPointCSV, rewardPoint)
		if err != nil {
			return errors.WrapTrace(err)
		}

		ee.RewardPoints[rewardPoint.Id] = rewardPoint
	}

	return nil
}

func (ee *ExploreEntry) reloadMonster(config *Config) error {
	for _, monsterCSV := range config.CfgExploreMonsterConfig.GetAllData() {
		monster := &ExploreMonster{}
		err := transfer.Transfer(monsterCSV, monster)
		if err != nil {
			return err
		}

		ee.Monsters[monster.Id] = monster
	}

	return nil
}

func (ee *ExploreEntry) reloadEventPoint(config *Config) error {
	for _, eventPointCSV := range config.CfgExploreEventConfig.GetAllData() {
		eventPoint := &ExploreEventPoint{}
		err := transfer.Transfer(eventPointCSV, eventPoint)
		if err != nil {
			return errors.WrapTrace(err)
		}

		ee.EventPoints[eventPoint.Id] = eventPoint
	}

	for id, eventPoint := range ee.EventPoints {
		switch eventPoint.EventType {
		case static.ExploreEventTypeNpc:
			npcId := eventPoint.EventParam
			npc, ok := ee.NPCs[npcId]
			if !ok {
				return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgExploreNPC, npcId))
			}

			npc.EventPointId = id
		case static.ExploreEventTypeDoodad:
			rewardPointId := eventPoint.EventParam
			rewardPoint, ok := ee.RewardPoints[rewardPointId]
			if !ok {
				return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgExploreGather, rewardPointId))
			}

			rewardPoint.EventPointId = id
		case static.ExploreEventTypeMonster:
			monsterId := eventPoint.EventParam
			monster, ok := ee.Monsters[monsterId]
			if !ok {
				return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgExploreMonster, monsterId))
			}

			monster.EventPointId = id
		}
	}

	return nil
}

func (ee *ExploreEntry) reloadFog(config *Config) error {
	for _, fogCSV := range config.CfgExploreFogConfig.GetAllData() {
		fog := &ExploreFog{}
		err := transfer.Transfer(fogCSV, fog)
		if err != nil {
			return errors.WrapTrace(err)
		}

		ee.Fogs[fog.Id] = fog
	}

	return nil
}

func (ee *ExploreEntry) reloadResource(config *Config) error {
	for _, resourceCSV := range config.CfgExploreResourceConfig.GetAllData() {
		resource := &ExploreResource{}
		err := transfer.Transfer(resourceCSV, resource)
		if err != nil {
			return err
		}

		condition := common.NewCondition(static.ConditionTypePassLevel, resourceCSV.Monstermap)
		resource.CollectConditions = common.NewConditions()
		resource.CollectConditions.AddCondition(condition)

		ee.Resources[resource.Id] = resource
	}

	return nil
}

func (ee *ExploreEntry) reloadTransportGate(config *Config) error {
	for _, transportGateCSV := range config.CfgTpGateConfig.GetAllData() {
		transportGate := &ExploreTransportGate{}
		err := transfer.Transfer(transportGateCSV, transportGate)
		if err != nil {
			return err
		}

		ee.TransportGates[transportGate.Id] = transportGate
	}

	return nil
}

func (ee *ExploreEntry) GetExploreMapObject(id int32) (*ExploreMapObject, error) {
	obj, ok := ee.MapObjects[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrMapObjectConfigNotFound, id)
	}

	return obj, nil
}

func (ee *ExploreEntry) GetExploreEventPoint(id int32) (*ExploreEventPoint, error) {
	obj, ok := ee.EventPoints[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrExploreEventPointConfigNotFound, id)
	}

	return obj, nil
}

func (ee *ExploreEntry) GetExploreNPC(id int32) (*ExploreNPC, error) {
	npc, ok := ee.NPCs[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrExploreNPCConfigNotFound, id)
	}

	return npc, nil
}

func (ee *ExploreEntry) GetExploreRewardPoint(id int32) (*ExploreRewardPoint, error) {
	rewardPoint, ok := ee.RewardPoints[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrExploreRewardPointConfigNotFound, id)
	}

	return rewardPoint, nil
}

func (ee *ExploreEntry) GetExploreMonster(id int32) (*ExploreMonster, error) {
	monster, ok := ee.Monsters[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrExploreMonsterConfigNotFound, id)
	}

	return monster, nil
}

func (ee *ExploreEntry) GetFog(id int32) (*ExploreFog, error) {
	fog, ok := ee.Fogs[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrExploreFogConfigNotFound, id)
	}
	return fog, nil
}

func (ee *ExploreEntry) GetResource(id int32) (*ExploreResource, error) {
	resource, ok := ee.Resources[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrExploreResourceConfigNotFound, id)
	}
	return resource, nil
}

func (ee *ExploreEntry) GetTransportGate(id int32) (*ExploreTransportGate, error) {
	transportGate, ok := ee.TransportGates[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrExploreTPGateConfigNotFound, id)
	}

	return transportGate, nil
}
