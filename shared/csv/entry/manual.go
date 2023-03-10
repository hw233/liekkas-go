package entry

import (
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/number"
	"sync"
)

const CfgNameManual = "cfg_manual"

type Manual struct {
	sync.RWMutex

	manualDrop       map[int32]int32
	characterManual  map[int32]int32
	worldItemManual  map[int32]int32
	manualTotalCount int32
	manualVersion    map[int32]int32
}

func NewManual() *Manual {
	return &Manual{
		manualDrop:       map[int32]int32{},
		characterManual:  map[int32]int32{},
		worldItemManual:  map[int32]int32{},
		manualTotalCount: 0,
		manualVersion:    map[int32]int32{},
	}
}

func (m *Manual) Reload(config *Config) error {
	m.Lock()
	defer m.Unlock()

	characterManual := map[int32]int32{}
	worldItemManual := map[int32]int32{}
	manualDrop := map[int32]int32{}
	manualVersion := map[int32]int32{}
	var manualTotalCount int32
	for _, manual := range config.CfgManualConfig.GetAllData() {
		manualTotalCount++

		manualDrop[manual.Id] = manual.DropId

		if manual.Type == static.ManualTypeCharacter {
			characterManual[manual.RelatedId] = manual.Id
		} else if manual.Type == static.ManualTypeWorldItem {
			worldItemManual[manual.RelatedId] = manual.Id
		}

		manualVersion[manual.Id] = manual.Version
	}
	arr := number.NewNonRepeatableArr()
	for _, cfg := range config.CfgExploreChapterLevelConfig.GetAllData() {
		if cfg.ChapterStory > 0 {
			arr.Append(cfg.ChapterStory)
		}
		if cfg.WinStory > 0 {
			arr.Append(cfg.WinStory)
		}

	}
	manualTotalCount += int32(len(arr.Values()))

	m.manualDrop = manualDrop
	m.characterManual = characterManual
	m.worldItemManual = worldItemManual
	m.manualTotalCount = manualTotalCount
	m.manualVersion = manualVersion
	return nil
}

func (m *Manual) FindManualIdByRelatedId(manualType, relatedId int32) (int32, bool) {
	m.RLock()
	defer m.RUnlock()
	if manualType == static.ManualTypeCharacter {
		manualId, ok := m.characterManual[relatedId]
		return manualId, ok
	} else if manualType == static.ManualTypeWorldItem {
		manualId, ok := m.worldItemManual[relatedId]
		return manualId, ok
	}
	return 0, false
}

func (m *Manual) getManualDrop(manualId int32) (int32, error) {
	m.RLock()
	defer m.RUnlock()
	dropId, ok := m.manualDrop[manualId]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameManual, manualId)
	}
	return dropId, nil
}

func (m *Manual) GetManualDrops(manualIds *number.NonRepeatableArr) ([]int32, error) {
	ret := make([]int32, 0, len(*manualIds))
	for _, id := range manualIds.Values() {
		drop, err := m.getManualDrop(id)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		ret = append(ret, drop)
	}
	return ret, nil
}

func (m *Manual) GetManualTotalCount() int32 {
	m.RLock()
	defer m.RUnlock()
	return m.manualTotalCount
}

func (m *Manual) GetManualVersion(manualId int32) (int32, error) {
	id, ok := m.manualVersion[manualId]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameManual, manualId)
	}

	return id, nil
}
