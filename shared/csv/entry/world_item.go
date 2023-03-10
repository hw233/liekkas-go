package entry

import (
	"math"
	"sort"
	"sync"

	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
)

const (
	CfgWorldItemData                = "cfg_world_item_data"
	CfgWorldItemLevelUp             = "cfg_world_item_level_up"
	CfgWorldItemAdvance             = "cfg_world_item_advance"
	CfgWorldItemRandAttributes      = "cfg_world_item_rand_attributes"
	CfgWorldItemRandAttributesValue = "cfg_world_item_rand_attributes_value"
)

type WorldItemEntry struct {
	sync.RWMutex

	worldItems         map[int32]WorldItemData
	levelEXP           map[int32]int32                // 等级经验 key: 等级*10+稀有度(makeLevelEXPKey), val: 经验值
	strengthenDiscount []WorldItemMaterialEXPDiscount // 强化折扣
	advance            map[int32]WorldItemAdvance     // 进阶 key: 阶级*10+稀有度(makeAdvanceKey)
	maxLevel           map[int8]int32                 // 最高等级
}

type WorldItemAdvance struct {
	LevelLimit int32 // 等级限制
	GoldCost   int32 // 金币消耗
	// ItemCost   int32 // 素材消耗
	Stage int32 // 最高阶级
}

type WorldItemRandAttr struct {
	Probs    []int32 // 权重
	ValueIDs []int32 // 数值ID
}

type WorldItemRandAttrValue struct {
	Probs  []int32   // 概率
	Ranges [][]int32 // 范围
}

type WorldItemData struct {
	WID int32 `src:"Id"`

	Rarity int32 // 品质
	// UseType      int32    // 使用类型 1可使用 2可出售 3可使用且可出售 4不可使用不可出售
	EXP int32 `src:"UseParam"` // 用作材料的强化经验
	// SellPrice    string   // 出售价格
	// Part         int32    // 装备部位
	// Careers    []int32 `src:"WorldItemLmt"` // 穿戴职业限制
	GoldPerExp float64 `src:"GoldperExp"` // 强化1经验消耗的金币
	// IsLocked     bool //
	// PassiveID    string  // 被动技能和升级规则
	// AttributeID  string  // 被动属性和升级规则
	// CanProduced  bool    // 是否能在建筑产出
	AdvanceWID map[int32]bool `ignore:"true"` // 可以用作突破的世界级道具ID
}

type WorldItemMaterialEXPDiscount struct {
	Level    int32   `json:"level"`
	Discount float64 `json:"discount"`
}

func NewWorldItemEntry() *WorldItemEntry {
	return &WorldItemEntry{
		worldItems:         map[int32]WorldItemData{},
		levelEXP:           map[int32]int32{},
		strengthenDiscount: []WorldItemMaterialEXPDiscount{},
		advance:            map[int32]WorldItemAdvance{},
		// randAttrs:          map[int32]WorldItemRandAttr{},
		// randAttrsValue:     map[int32]WorldItemRandAttrValue{},
	}
}

func (e *WorldItemEntry) Check(config *Config) error {
	// err := e.Reload(config, global)
	// if err != nil {
	// 	return err
	// }
	//
	// // 检查是否有0经验的等级
	// for _, exp := range e.levelEXP {
	// 	// 缺少exp
	// 	if exp == 0 {
	// 		// todo: 定义err
	// 		return err
	// 	}
	// }

	return nil
}

func (e *WorldItemEntry) Reload(config *Config) error {
	e.Lock()
	defer e.Unlock()

	// ------------------------------------------------------------------------------------
	worldItems := map[int32]WorldItemData{}

	for _, breakItems := range config.CfgWorldItemDataConfig.GetAllData() {
		worldItem := WorldItemData{
			AdvanceWID: map[int32]bool{},
		}

		err := transfer.Transfer(breakItems, &worldItem)
		if err != nil {
			return errors.WrapTrace(err)
		}

		for _, v := range breakItems.BreakItem {
			worldItem.AdvanceWID[v] = true
		}

		worldItems[worldItem.WID] = worldItem
	}

	// ------------------------------------------------------------------------------------
	levelEXP := map[int32]int32{}

	for _, levelUp := range config.CfgWorldItemLevelUpDataConfig.GetAllData() {
		levelEXP[e.makeLevelEXPKey(levelUp.Level, int8(levelUp.Rarity))] = levelUp.Exp
	}

	// ------------------------------------------------------------------------------------
	advance := map[int32]WorldItemAdvance{}

	for rarity, adv := range config.CfgWorldItemAdvanceConfig.GetAllData() {
		if len(adv.GoldCost) != int(adv.Stage) ||
			// len(adv.ItemCost) != int(adv.Stage) ||
			len(adv.LevelLimit) != int(adv.Stage) {
			return errors.Swrapf(common.ErrCSVFormatInvalid, "CfgWorldItemRandAttributes", rarity)
		}

		for i := int32(0); i < adv.Stage; i++ {
			advance[e.makeAdvanceKey(i+1, int8(rarity))] = WorldItemAdvance{
				LevelLimit: adv.LevelLimit[i],
				GoldCost:   adv.GoldCost[i],
				// ItemCost:   adv.ItemCost[i],
				Stage: adv.Stage,
			}
		}
	}

	// ------------------------------------------------------------------------------------
	e.worldItems = worldItems
	e.levelEXP = levelEXP
	e.advance = advance
	e.strengthenDiscount = config.WorldItemMaterialEXPDiscount
	sort.Slice(e.strengthenDiscount, func(i, j int) bool {
		return e.strengthenDiscount[i].Level < e.strengthenDiscount[i].Level
	})

	e.maxLevel = config.WorldItemMaxLevel

	return nil
}

func (e *WorldItemEntry) makeLevelEXPKey(level int32, rarity int8) int32 {
	return level*10 + int32(rarity)
}

func (e *WorldItemEntry) makeAdvanceKey(stage int32, rarity int8) int32 {
	return stage*10 + int32(rarity)
}

// 执行下面依赖wid的查询需要先检查wid是否存在，下面使用wid的函数不包含error的必须调用
func (e *WorldItemEntry) CheckWIDExist(wid int32) error {
	e.RLock()
	defer e.RUnlock()

	_, ok := e.worldItems[wid]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, wid)
	}

	return nil
}

// 检查职业是否符合
// func (e *WorldItemEntry) CheckCareer(wid, career int32) error {
// 	e.RLock()
// 	defer e.RUnlock()
//
// 	worldItem, ok := e.worldItems[wid]
// 	if !ok {
// 		return errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, wid)
// 	}
//
// 	for _, v := range worldItem.Careers {
// 		if career == v {
// 			return nil
// 		}
// 	}
//
// 	return errors.Swrapf(common.ErrWorldItemNotMatchCareer, wid, career)
// }

// 检查进阶素材
func (e *WorldItemEntry) CheckAdvanceMaterials(target *common.WorldItem, materials []*common.WorldItem, itemID int32) error {
	e.RLock()
	defer e.RUnlock()

	data, ok := e.worldItems[target.WID]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, target.WID)
	}

	// 检查世界级道具
	for _, worldItem := range materials {
		// 消耗世界级道具只能用本体
		if target.WID != worldItem.WID {
			return errors.Swrapf(common.ErrWorldItemAdvanceMaterialNotMatch, worldItem, worldItem.WID)
		}
	}

	// 检查升星道具
	if itemID != 0 && !data.AdvanceWID[itemID] {
		return errors.Swrapf(common.ErrWorldItemAdvanceMaterialNotMatch, target.WID, itemID)
	}

	// 下一阶级
	nextStage := target.Stage.Value() + 1

	// 检查材料数量，世界级道具和升星道具有一样就可以了
	if int32(len(materials)) < 1 && itemID == 0 {
		return errors.Swrapf(common.ErrWorldItemAdvanceMaterialNoEnough, target.WID, target.Rarity, nextStage, materials)
	}

	return nil
}

// 检查是否满足进阶条件
func (e *WorldItemEntry) CheckAdvanceLevel(worldItem *common.WorldItem) error {
	e.RLock()
	defer e.RUnlock()

	needLevel := e.advance[e.makeAdvanceKey(worldItem.Stage.Value(), worldItem.Rarity)].LevelLimit

	if !worldItem.Level.Enough(needLevel) {
		return errors.Swrapf(common.ErrWorldItemAdvanceLevelNotEnough, worldItem.WID, worldItem.Level, needLevel)
	}

	return nil
}

// 检查是否满足进阶条件
func (e *WorldItemEntry) CheckStageUpToLimit(worldItem *common.WorldItem) error {
	e.RLock()
	defer e.RUnlock()

	key := e.makeAdvanceKey(1, worldItem.Rarity)

	advance, ok := e.advance[key]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentAdvance, key)
	}

	if worldItem.Stage.Value() >= advance.Stage {
		return common.ErrWorldItemStageUpToLimit
	}

	return nil
}

// 检查是否满足进阶条件
func (e *WorldItemEntry) AdvanceCostGold(worldItem *common.WorldItem) (int32, error) {
	e.RLock()
	defer e.RUnlock()

	key := e.makeAdvanceKey(worldItem.Stage.Value()+1, worldItem.Rarity)

	advance, ok := e.advance[key]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemAdvance, key)
	}

	return advance.GoldCost, nil
}

// 初始化装备表格相关数据
func (e *WorldItemEntry) NewWorldItem(id int64, wid int32) (*common.WorldItem, error) {
	e.RLock()
	defer e.RUnlock()

	worldItem := common.NewWorldItem(id, wid)

	data, ok := e.worldItems[worldItem.WID]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, worldItem.WID)
	}

	// 随机阵营
	// worldItem.Camp = int8(rand.SinglePerm(data.CampProbs))
	// 稀有度
	worldItem.Rarity = int8(data.Rarity)

	return worldItem, nil
}

// 同步等级并解锁随机属性
func (e *WorldItemEntry) SyncLevelAndUnlockAttr(worldItem *common.WorldItem) error {
	e.RLock()
	defer e.RUnlock()

	maxLevel := e.maxLevel[worldItem.Rarity]

	// 当装备等级小于最大等级，循环判断经验是否足够
	for worldItem.Level.Value() < maxLevel {
		// 获取下一个等级升级经验
		nextLevelEXPKey := e.makeLevelEXPKey(worldItem.Level.Value()+1, worldItem.Rarity)

		nextLevelEXP, ok := e.levelEXP[nextLevelEXPKey]
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemLevelUp, nextLevelEXPKey)
		}

		// 升级经验不足，跳出
		if !worldItem.EXP.Enough(nextLevelEXP) {
			break
		}

		// 升级了！
		worldItem.Level.Plus(1)
	}

	return nil
}

// 计算装备强化经验
func (e *WorldItemEntry) StrengthenEXP(target *common.WorldItem, itemEXP int32, materials []*common.WorldItem, maxLevelByTeam int32) (int32, error) {
	e.RLock()
	defer e.RUnlock()

	// 计算材料的强化经验值
	exps, err := e.strengthenMaterialsEXP(materials)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	// 世界级道具等级最大值上限，经验溢出不保留，世界级道具最大等级 = 阶级加成 + 队伍等级加成
	// maxLevel := e.maxLevel[target.Rarity]
	maxLevel := e.advance[e.makeAdvanceKey(target.Stage.Value(), target.Rarity)].LevelLimit + maxLevelByTeam

	expNow := target.EXP.Value()

	var addEXP int32 = itemEXP // 强化实际增加的经验

	levelEXPKey := e.makeLevelEXPKey(maxLevel, target.Rarity)

	// 装备最大经验，不可溢出
	maxEXP, ok := e.levelEXP[levelEXPKey]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemLevelUp, levelEXPKey)
	}

	// 装备等级已经可以升到上限
	for _, exp := range exps {
		if expNow+exp > maxEXP {
			// 检查经验溢出了还有后续装备的话就报错
			// if i != len(exps)-1 {
			// 	return 0, errors.Swrapf(common.ErrWorldItemEXPUpToLimit, i)
			// }

			// 计算溢出后真正加上的经验
			addEXP += maxEXP - expNow
			break
		} else {
			addEXP += exp
		}

		expNow += exp
	}

	// 注释是因为当初装备强化溢出经验也要保留，现在不要了
	// // 装备等级上限受到账号等级制约
	// limitEXP, err := manager.CSV.WorldItem.LevelEXP(worldItem.WID, limitLevel)
	// if err != nil {
	// 	return 0, errors.WrapTrace(err)
	// }
	//
	// for i, exp := range exps {
	// 	if expNow+exp > limitEXP {
	// 		// 检查经验溢出了还有后续装备的话就报错
	// 		if i != len(exps)-1 {
	// 			return 0, errors.WrapTrace(common.ErrWorldItemEXPUpToLimit)
	// 		}
	//
	// 		// 同时超过账号限制和最大等级限制的情况
	// 		if expNow+exp > maxEXP {
	// 			// 计算溢出后真正加上的经验
	// 			addEXP += maxEXP - expNow
	// 			break
	// 		}
	//
	// 		// 可以溢出
	// 	}
	//
	// 	addEXP += exp
	// 	expNow += exp
	// }

	return addEXP, nil

}

// 计算装备强化经验
func (e *WorldItemEntry) strengthenMaterialsEXP(materials []*common.WorldItem) ([]int32, error) {
	exps := make([]int32, 0, len(materials))

	for _, worldItem := range materials {
		data, ok := e.worldItems[worldItem.WID]
		if !ok {
			return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, worldItem.WID)
		}

		level := worldItem.Level.Value() // 装备等级

		var sum float64     // 这个装备强化经验折扣后的和\
		var lastEXP float64 // 上一级强化等级经验

		// 计算强化经验的折扣
		for _, val := range e.strengthenDiscount {
			if level >= val.Level {
				// 全额折扣
				levelEXPKey := e.makeLevelEXPKey(val.Level, worldItem.Rarity)

				exp, ok := e.levelEXP[levelEXPKey] // 当前等级经验
				if !ok {
					return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemLevelUp, levelEXPKey)
				}

				sum += (float64(exp) - lastEXP) * (val.Discount / 10000)

				lastEXP = float64(exp)
			} else {
				// 部分折扣
				sum += (float64(worldItem.EXP.Value()) - lastEXP) * (val.Discount / 10000)
				break
			}
		}

		// 基本经验 + 强化经验
		exps = append(exps, data.EXP+int32(math.Trunc(sum)))
	}

	return exps, nil
}

// func (e *WorldItemEntry) StrengthenGoldCost(worldItem *common.WorldItem, addEXP int32) (int32, error) {
// 	e.RLock()
// 	defer e.RUnlock()
//
// 	data, ok := e.worldItems[worldItem.WID]
// 	if !ok {
// 		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, worldItem.WID)
// 	}
//
// 	return data.GoldPerExp * addEXP, nil
// }

// 装备强化金币只计算装备固定值，不计算折扣金币
func (e *WorldItemEntry) StrengthenGoldCost(worldItem *common.WorldItem, itemEXP int32, materials []*common.WorldItem) (int32, error) {
	e.RLock()
	defer e.RUnlock()

	data, ok := e.worldItems[worldItem.WID]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, worldItem.WID)
	}

	costGolds := data.GoldPerExp * float64(itemEXP)

	for _, material := range materials {
		md, ok := e.worldItems[material.WID]
		if !ok {
			return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemData, material.WID)
		}

		costGolds += data.GoldPerExp * float64(md.EXP)
	}

	return int32(math.Floor(costGolds)), nil
}

// 计算强化的阶级
func (e *WorldItemEntry) CalAddStage(worldItem *common.WorldItem, materials []*common.WorldItem, itemID int32) (int32, error) {
	e.RLock()
	defer e.RUnlock()

	key := e.makeAdvanceKey(1, worldItem.Rarity)

	advance, ok := e.advance[key]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgWorldItemAdvance, key)
	}

	var allMaterialStage int32 = 0
	for _, v := range materials {
		allMaterialStage += v.Stage.Value() + 1
	}

	// 升星道具
	if itemID != 0 {
		allMaterialStage++
	}

	var retStage = worldItem.Stage.Value() + allMaterialStage

	// 溢出
	if retStage > advance.Stage {
		retStage = advance.Stage
	}

	return retStage - worldItem.Stage.Value(), nil
}
