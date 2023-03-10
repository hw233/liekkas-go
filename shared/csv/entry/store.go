package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type Store struct {
	Id              int32
	UpdateRule      int32
	UnlockCondition *common.Conditions `rule:"conditions"`
	Currencies      []int32
	SubStores       []int32
	StoreType       int32 `src:"SubStoreType"`
}

type Cell struct {
	Goods           []int32
	UnlockCondition *common.Conditions `rule:"conditions"`
}

type SubStoreOrigin struct {
	Id                       int32
	UnlockCondition          *common.Conditions   `rule:"conditions"`
	GoodsIndexs              [][]int32            `rule:"mergeIntSlice" src:"Cell0,Cell1,Cell2,Cell3,Cell4,Cell5,Cell6,Cell7,Cell8,Cell9,Cell10,Cell11,Cell12,Cell13,Cell14,Cell15,Cell16,Cell17,Cell18,Cell19"`
	UnlockConditionsForCells []*common.Conditions `rule:"mergeConditions" src:"UnlockCondition0,UnlockCondition1,UnlockCondition2,UnlockCondition3,UnlockCondition4,UnlockCondition5,UnlockCondition6,UnlockCondition7,UnlockCondition8,UnlockCondition9,UnlockCondition10,UnlockCondition11,UnlockCondition12,UnlockCondition13,UnlockCondition14,UnlockCondition15,UnlockCondition16,UnlockCondition17,UnlockCondition18,UnlockCondition19"`
}

type SubStore struct {
	Id              int32
	UnlockCondition *common.Conditions
	Cells           []Cell
}

type RewardsIdAndCount struct {
	ItemID int32
	Cnt    int32
}

type Goods struct {
	Id              int32
	Gain            *common.Rewards    `rule:"rewards" src:"RewardsIdAndCnt"` // 商品对应的可以获得的道具以及道具数量
	UnlockCondition *common.Conditions `rule:"conditions"`
	Times           int32
	Probability     int32
	Currencies      []int32
	Price           []int32 `src:"RealPrice"`
}

type UpdateRule struct {
	Id       int32
	Period   int32 // 自动刷新周期,为0代表不会自动刷新
	Times    int32 `src:"TimesLimit"` // 周期内可强制刷新的次数
	Currency int32
	Cnt      []int32 // 强制刷新每次的花费数额
}

type StoreEntry struct {
	sync.RWMutex

	Stores      map[int32]Store
	Substores   map[int32]SubStore
	GoodsData   map[int32]Goods
	UpdateRules map[int32]UpdateRule
}

func NewStoreEntry() *StoreEntry {
	return &StoreEntry{}
}

func (se *StoreEntry) Check(config *Config) error {

	return nil
}

func (se *StoreEntry) Reload(config *Config) error {
	se.Lock()
	defer se.Unlock()

	stores := map[int32]Store{}
	subStores := map[int32]SubStore{}
	goods := map[int32]Goods{}
	updateRules := map[int32]UpdateRule{}

	for _, storeCsv := range config.CfgStoreGeneralConfig.GetAllData() {
		storeCfg := &Store{}

		err := transfer.Transfer(storeCsv, storeCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		if storeCfg.StoreType != 1 && storeCfg.StoreType != 2 {
			return errors.Swrapf(common.ErrStoreWrongSubStoreTypeForData, storeCsv.Id)
		}
		stores[storeCsv.Id] = *storeCfg
	}

	for _, v := range config.CfgStoreSubstoreConfig.GetAllData() {
		subCfgOrigin := &SubStoreOrigin{}

		err := transfer.Transfer(v, subCfgOrigin)
		if err != nil {
			return errors.WrapTrace(err)
		}

		subCfg, err := subCfgOrigin.ToSubStore()
		if err != nil {
			return errors.WrapTrace(err)
		}

		subStores[v.Id] = *subCfg
	}

	for _, goodsCsv := range config.CfgStoreGoodsConfig.GetAllData() {
		goodsCfg := &Goods{}

		err := transfer.Transfer(goodsCsv, goodsCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		if len(goodsCfg.Currencies) != len(goodsCfg.Price) {
			return errors.Swrapf(common.ErrStoreCurrencyNotMatchPriceForData, goodsCsv.Id)
		}

		goods[goodsCsv.Id] = *goodsCfg
	}

	for _, updateCsv := range config.CfgStoreUpdateConfig.GetAllData() {
		updateCfg := &UpdateRule{}

		err := transfer.Transfer(updateCsv, updateCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		if int(updateCfg.Times) > 0 && int(updateCfg.Times) < len(updateCfg.Cnt) {
			return errors.Swrapf(common.ErrStoreUpdateTimesNotMatchCntOfUpdateCost, updateCsv.Id)
		}

		// 如果强制刷新的次数上限为10次，但是只填了三次花费的数额，那么把最后一次的花费金额延续到后面的七次
		for i := 0; i < int(updateCfg.Times)-len(updateCfg.Cnt); i++ {
			toBeComplete := int32(0)
			if len(updateCfg.Cnt) > 0 {
				toBeComplete = updateCfg.Cnt[len(updateCfg.Cnt)-1]
			}
			updateCfg.Cnt = append(updateCfg.Cnt, toBeComplete)
		}

		updateRules[updateCsv.Id] = *updateCfg
	}

	se.Stores = stores
	se.Substores = subStores
	se.GoodsData = goods
	se.UpdateRules = updateRules

	return nil
}

func (so *SubStoreOrigin) ToSubStore() (*SubStore, error) {
	subStore := &SubStore{}
	var cells []Cell
	// cells := []Cell{}

	if len(so.GoodsIndexs) != len(so.UnlockConditionsForCells) {
		return nil, errors.WrapTrace(common.ErrCSVFormatInvalid)
	}

	subStore.Id = so.Id
	subStore.UnlockCondition = so.UnlockCondition
	for i := 0; i < len(so.GoodsIndexs); i++ {
		cell := Cell{}
		cell.Goods = so.GoodsIndexs[i]
		cell.UnlockCondition = so.UnlockConditionsForCells[i]
		cells = append(cells, cell)
	}
	subStore.Cells = cells

	return subStore, nil
}

func (se *StoreEntry) GetStoreGeneral(id int32) (*Store, error) {
	se.RLock()
	defer se.RUnlock()

	store, ok := se.Stores[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrStoreWrongStoreIDForData, id)
	}

	return &store, nil
}

func (se *StoreEntry) GetSubIDs(id int32) ([]int32, error) {
	se.RLock()
	defer se.RUnlock()

	store, ok := se.Stores[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrStoreWrongStoreIDForData, id)
	}

	return store.SubStores, nil
}

func (se *StoreEntry) GetSubStore(id int32) (*SubStore, error) {
	se.RLock()
	defer se.RUnlock()

	subStore, ok := se.Substores[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrStoreWrongSubStoreIDForData, id)
	}

	return &subStore, nil
}

func (se *StoreEntry) GetGoods(id int32) (*Goods, error) {
	se.RLock()
	defer se.RUnlock()

	goods, ok := se.GoodsData[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrStoreWrongGoodsIDForData, id)
	}

	return &goods, nil
}

func (se *StoreEntry) GetUpdateRule(id int32) (*UpdateRule, error) {
	se.RLock()
	defer se.RUnlock()

	updaterule, ok := se.UpdateRules[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrStoreWrongUpdateRuleIDForData, id)
	}

	return &updaterule, nil
}
