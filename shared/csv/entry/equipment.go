package entry

import (
	"log"
	"math"
	"sort"
	"sync"

	"shared/common"
	"shared/utility/errors"
	"shared/utility/rand"
	"shared/utility/transfer"
)

const (
	CfgEquipmentData                = "cfg_equip_data"
	CfgEquipmentLevelUp             = "cfg_equip_level_up"
	CfgEquipmentAdvance             = "cfg_equip_advance"
	CfgEquipmentRandAttributes      = "cfg_equip_rand_attributes"
	CfgEquipmentRandAttributesValue = "cfg_equip_rand_attributes_value"

	// EquipmentID = "id"
)

type EquipmentEntry struct {
	sync.RWMutex

	equipments         map[int32]EquipmentData
	levelEXP           map[int32]int32                  // 等级经验 key: 等级*10+稀有度(makeLevelEXPKey), val: 经验值
	strengthenDiscount []EquipmentMaterialEXPDiscount   // 强化折扣
	advance            map[int32]EquipmentAdvance       // 进阶 key: 阶级*10+稀有度(makeAdvanceKey)
	randAttrs          map[int32]EquipmentRandAttr      // 随机属性
	randAttrsValue     map[int32]EquipmentRandAttrValue // 随机属性数值
	maxLevel           map[int8]int32                   // 最高等级
	attrUnlockLevel    []int32                          // 随机属性解锁等级等级
}

type EquipmentAdvance struct {
	LevelLimit int32 // 等级限制
	GoldCost   int32 // 金币消耗
	ItemCost   int32 // 素材消耗
	Stage      int32 // 最高阶级
}

type EquipmentRandAttr struct {
	Probs    []int32 // 权重
	ValueIDs []int32 // 数值ID
}

type EquipmentRandAttrValue struct {
	Probs  []int32   // 概率
	Ranges [][]int32 // 范围
}

type EquipmentData struct {
	EID int32 `src:"Id"`

	Rarity int32 // 品质
	// UseType      int32    // 使用类型 1可使用 2可出售 3可使用且可出售 4不可使用不可出售
	EXP int32 `src:"UseParam"` // 用作材料的强化经验
	// SellPrice    string   // 出售价格
	Part         int32           // 装备部位
	Careers      []int32         `src:"EquipLmt"`                  // 穿戴职业限制
	RecastCost   *common.Rewards `src:"RepickCost" rule:"rewards"` // 重铸消耗
	CampAddition int32           // 阵营加成（万分比）
	CampProbs    []int32         `src:"CampWeight"` // 阵营权重，长度=5
	GoldPerExp   float64         `src:"GoldperExp"` // 强化1经验消耗的金币
	// IsLocked     bool //
	// PassiveID    string  // 被动技能和升级规则
	// AttributeID  string  // 被动属性和升级规则
	// CanProduced  bool    // 是否能在建筑产出
	AdvanceEID  map[int32]bool `ignore:"true"` // 可以用作突破的装备ID
	RandEntries []int32        // 随机属性词条
}

type EquipmentMaterialEXPDiscount struct {
	Level    int32   `json:"level"`
	Discount float64 `json:"discount"`
}

func NewEquipmentEntry() *EquipmentEntry {
	return &EquipmentEntry{
		equipments:         map[int32]EquipmentData{},
		levelEXP:           map[int32]int32{},
		strengthenDiscount: []EquipmentMaterialEXPDiscount{},
		advance:            map[int32]EquipmentAdvance{},
		randAttrs:          map[int32]EquipmentRandAttr{},
		randAttrsValue:     map[int32]EquipmentRandAttrValue{},
		attrUnlockLevel:    []int32{},
	}
}

func (e *EquipmentEntry) Check(config *Config) error {
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

func (e *EquipmentEntry) Reload(config *Config) error {
	e.Lock()
	defer e.Unlock()

	// ------------------------------------------------------------------------------------
	equipments := map[int32]EquipmentData{}

	for _, equip := range config.CfgEquipDataConfig.GetAllData() {
		equipment := EquipmentData{
			AdvanceEID: map[int32]bool{},
		}

		err := transfer.Transfer(equip, &equipment)
		if err != nil {
			return errors.WrapTrace(err)
		}

		for _, v := range equip.BreakEquipId {
			equipment.AdvanceEID[v] = true
		}

		equipments[equipment.EID] = equipment
	}

	// ------------------------------------------------------------------------------------
	levelEXP := map[int32]int32{}

	for _, levelUp := range config.CfgEquipLevelUpDataConfig.GetAllData() {
		levelEXP[e.makeLevelEXPKey(levelUp.Level, int8(levelUp.Rarity))] = levelUp.Exp
	}
	// ------------------------------------------------------------------------------------
	advance := map[int32]EquipmentAdvance{}

	for rarity, adv := range config.CfgEquipAdvanceConfig.GetAllData() {
		if len(adv.GoldCost) != int(adv.Stage) ||
			len(adv.ItemCost) != int(adv.Stage) ||
			len(adv.LevelLimit) != int(adv.Stage) {
			return errors.Swrapf(common.ErrCSVFormatInvalid, CfgEquipmentAdvance, rarity)
		}

		for i := int32(0); i < adv.Stage; i++ {
			advance[e.makeAdvanceKey(i+1, int8(rarity))] = EquipmentAdvance{
				LevelLimit: adv.LevelLimit[i],
				GoldCost:   adv.GoldCost[i],
				ItemCost:   adv.ItemCost[i],
				Stage:      adv.Stage,
			}
		}
	}
	// ------------------------------------------------------------------------------------
	randAttrs := map[int32]EquipmentRandAttr{}

	for id, attr := range config.CfgEquipRandAttributesConfig.GetAllData() {
		// 检查数组长度
		if len(attr.HpAddPercent) != 2 ||
			len(attr.AtkAddPercent) != 2 ||
			len(attr.MAtkAddPercent) != 2 ||
			len(attr.DefAddPercent) != 2 ||
			len(attr.MDefAddPercent) != 2 ||
			len(attr.CritAddPercent) != 2 ||
			len(attr.CritDamPercent) != 2 ||
			len(attr.CritDamReducePercent) != 2 ||
			len(attr.CureAddPercent) != 2 ||
			len(attr.PhyDamAdd) != 2 ||
			len(attr.PhyDamReduce) != 2 ||
			len(attr.MagDamAdd) != 2 ||
			len(attr.MagDamReduce) != 2 ||
			len(attr.PhyPen) != 2 ||
			len(attr.MagPen) != 2 {
			return errors.Swrapf(common.ErrCSVFormatInvalid, CfgEquipmentRandAttributes, id)
		}

		randAttrs[id] = EquipmentRandAttr{
			Probs: []int32{
				attr.HpAddPercent[0],
				attr.AtkAddPercent[0],
				attr.MAtkAddPercent[0],
				attr.DefAddPercent[0],
				attr.MDefAddPercent[0],
				attr.CritAddPercent[0],
				attr.CritDamPercent[0],
				attr.CritDamReducePercent[0],
				attr.CureAddPercent[0],
				attr.PhyDamAdd[0],
				attr.PhyDamReduce[0],
				attr.MagDamAdd[0],
				attr.MagDamReduce[0],
				attr.PhyPen[0],
				attr.MagPen[0],
			},
			ValueIDs: []int32{
				attr.HpAddPercent[1],
				attr.AtkAddPercent[1],
				attr.MAtkAddPercent[1],
				attr.DefAddPercent[1],
				attr.MDefAddPercent[1],
				attr.CritAddPercent[1],
				attr.CritDamPercent[1],
				attr.CritDamReducePercent[1],
				attr.CureAddPercent[1],
				attr.PhyDamAdd[1],
				attr.PhyDamReduce[1],
				attr.MagDamAdd[1],
				attr.MagDamReduce[1],
				attr.PhyPen[1],
				attr.MagPen[1],
			},
		}
	}
	// ------------------------------------------------------------------------------------
	randAttrsValue := map[int32]EquipmentRandAttrValue{}

	for id, value := range config.CfgEquipRandAttributesValueConfig.GetAllData() {
		if len(value.Range1) != 2 ||
			len(value.Range2) != 2 ||
			len(value.Range3) != 2 {
			return errors.Swrapf(common.ErrCSVFormatInvalid, CfgEquipmentRandAttributesValue, id)
		}

		randAttrsValue[id] = EquipmentRandAttrValue{
			Probs:  []int32{value.Prob1, value.Prob2, value.Prob3},
			Ranges: [][]int32{value.Range1, value.Range2, value.Range3},
		}
	}

	// ------------------------------------------------------------------------------------
	e.equipments = equipments
	e.levelEXP = levelEXP
	e.advance = advance
	e.randAttrs = randAttrs
	e.randAttrsValue = randAttrsValue

	e.strengthenDiscount = config.EquipmentMaterialEXPDiscount
	sort.Slice(e.strengthenDiscount, func(i, j int) bool {
		return e.strengthenDiscount[i].Level < e.strengthenDiscount[i].Level
	})

	e.maxLevel = config.EquipmentMaxLevel
	e.attrUnlockLevel = config.EquipmentAttrUnlockLevel

	return nil
}

func (e *EquipmentEntry) makeLevelEXPKey(level int32, rarity int8) int32 {
	return level*10 + int32(rarity)
}

func (e *EquipmentEntry) makeAdvanceKey(stage int32, rarity int8) int32 {
	return stage*10 + int32(rarity)
}

// 执行下面依赖eid的查询需要先检查eid是否存在，下面使用eid的函数不包含error的必须调用
func (e *EquipmentEntry) CheckEIDExist(eid int32) error {
	e.RLock()
	defer e.RUnlock()

	_, ok := e.equipments[eid]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, eid)
	}

	return nil
}

// 检查职业是否符合
func (e *EquipmentEntry) CheckCareer(eid, career int32) error {
	e.RLock()
	defer e.RUnlock()

	equipment, ok := e.equipments[eid]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, eid)
	}

	for _, v := range equipment.Careers {
		if career == v {
			return nil
		}
	}

	return errors.Swrapf(common.ErrEquipmentNotMatchCareer, eid, career)
}

// 检查部位是否符合
func (e *EquipmentEntry) Part(eid int32) (int32, error) {
	e.RLock()
	defer e.RUnlock()

	equipment, ok := e.equipments[eid]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, eid)
	}

	// 检查part是否合法
	if equipment.Part < 1 || equipment.Part > 4 {
		return 0, errors.Swrapf(common.ErrCSVFormatInvalid, CfgEquipmentData, eid)
	}

	return equipment.Part, nil
}

// 检查进阶素材
func (e *EquipmentEntry) CheckAdvanceMaterials(target *common.Equipment, materials []*common.Equipment) error {
	e.RLock()
	defer e.RUnlock()

	data, ok := e.equipments[target.EID]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, target.EID)
	}

	// 检查材料ID
	for _, equipment := range materials {
		// 同名也可以消耗
		if target.EID == equipment.EID {
			continue
		}

		if !data.AdvanceEID[equipment.EID] {
			return errors.Swrapf(common.ErrEquipmentAdvanceMaterialNotMatch, equipment, equipment.EID)
		}
	}

	nextStage := target.Stage.Value() + 1

	// 检查材料数量
	if int32(len(materials)) < e.advance[e.makeAdvanceKey(nextStage, target.Rarity)].ItemCost {
		return errors.Swrapf(common.ErrEquipmentAdvanceMaterialNoEnough, target.EID, target.Rarity, nextStage, materials)
	}

	return nil
}

// 检查是否满足进阶条件
func (e *EquipmentEntry) CheckAdvanceLevel(equipment *common.Equipment) error {
	e.RLock()
	defer e.RUnlock()

	needLevel := e.advance[e.makeAdvanceKey(equipment.Stage.Value(), equipment.Rarity)].LevelLimit

	if !equipment.Level.Enough(needLevel) {
		return errors.Swrapf(common.ErrEquipmentAdvanceLevelNotEnough, equipment.EID, equipment.Level, needLevel)
	}

	return nil
}

// 检查是否满足进阶条件
func (e *EquipmentEntry) CheckStageUpToLimit(equipment *common.Equipment) error {
	e.RLock()
	defer e.RUnlock()

	key := e.makeAdvanceKey(1, equipment.Rarity)

	advance, ok := e.advance[key]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentAdvance, key)
	}

	if equipment.Stage.Value() >= advance.Stage {
		return common.ErrEquipmentStageUpToLimit
	}

	return nil
}

// 检查是否满足进阶条件
func (e *EquipmentEntry) AdvanceCostGold(equipment *common.Equipment) (int32, error) {
	e.RLock()
	defer e.RUnlock()

	key := e.makeAdvanceKey(equipment.Stage.Value()+1, equipment.Rarity)

	advance, ok := e.advance[key]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentAdvance, key)
	}

	return advance.GoldCost, nil
}

// 初始化装备表格相关数据
func (e *EquipmentEntry) NewEquipment(id int64, eid int32) (*common.Equipment, error) {
	e.RLock()
	defer e.RUnlock()

	equipment := common.NewEquipment(id, eid)

	data, ok := e.equipments[equipment.EID]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, equipment.EID)
	}

	// 随机阵营
	equipment.Camp = int8(rand.SinglePerm(data.CampProbs))
	// 稀有度
	equipment.Rarity = int8(data.Rarity)

	// 随机出所有随机属性，但是不解锁
	// 随机属性数量
	attrNum := len(data.RandEntries)

	// 不能超过通用配置
	if attrNum > len(e.attrUnlockLevel) {
		attrNum = len(e.attrUnlockLevel)
	}

	attrs, err := e.randAllAttr(equipment.EID, attrNum)
	if err != nil {
		return equipment, errors.WrapTrace(err)
	}

	log.Printf("num: %d, attr: %v", len(data.RandEntries), attrs)

	equipment.Attrs = attrs

	return equipment, nil
}

// 重铸阵营
func (e *EquipmentEntry) RecastCamp(equipment *common.Equipment) (int8, error) {
	e.RLock()
	defer e.RUnlock()

	data, ok := e.equipments[equipment.EID]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, equipment.EID)
	}

	// 不可以重铸
	if len(data.CampProbs) == 0 {
		return 0, errors.Swrapf(common.ErrEquipmentNoCamp, equipment.EID)
	}

	return int8(rand.SinglePerm(data.CampProbs)) + 1, nil
}

// 重铸阵营
func (e *EquipmentEntry) RecastCost(equipment *common.Equipment) (*common.Rewards, error) {
	e.RLock()
	defer e.RUnlock()

	data, ok := e.equipments[equipment.EID]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, equipment.EID)
	}

	return data.RecastCost, nil
}

// 同步等级并解锁随机属性
func (e *EquipmentEntry) SyncLevelAndUnlockAttr(equipment *common.Equipment) error {
	e.RLock()
	defer e.RUnlock()

	maxLevel := e.maxLevel[equipment.Rarity]

	// 当装备等级小于最大等级，循环判断经验是否足够
	for equipment.Level.Value() < maxLevel {
		// 获取下一个等级升级经验
		nextLevelEXPKey := e.makeLevelEXPKey(equipment.Level.Value()+1, equipment.Rarity)

		nextLevelEXP, ok := e.levelEXP[nextLevelEXPKey]
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentLevelUp, nextLevelEXPKey)
		}

		// 升级经验不足，跳出
		if !equipment.EXP.Enough(nextLevelEXP) {
			break
		}

		// 升级了！
		equipment.Level.Plus(1)

		// 检查是否解锁随机属性
		data, ok := e.equipments[equipment.EID]
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, equipment.EID)
		}

		// 未解锁的随机属性数量
		unlockAttrNum := equipment.UnlockAttrNum()
		// 最大随机属性数量
		maxAttrNum := len(data.RandEntries)

		// 不能超过通用配置
		if maxAttrNum > len(e.attrUnlockLevel) {
			maxAttrNum = len(e.attrUnlockLevel)
		}

		// 还有未解锁的属性
		if unlockAttrNum < maxAttrNum {
			// 判断装备的最大属性数量和通用配置的最大属性数量
			// if attrNum < len(data.RandEntries) && attrNum < len(e.attrUnlockLevel) {
			// 检查等级是否足够
			if equipment.Level.Enough(e.attrUnlockLevel[unlockAttrNum]) {
				// 满足条件，解锁新的随机属性
				// attr, value, err := e.randAttr(equipment.EID, attrNum)
				// if err != nil {
				// 	return errors.WrapTrace(err)
				// }

				equipment.UnlockNextAttr()
			}
			// }
		}
	}

	return nil
}

// 计算装备强化经验
func (e *EquipmentEntry) StrengthenEXP(target *common.Equipment, itemEXP int32, materials []*common.Equipment, maxLevelByTeam int32) (int32, error) {
	e.RLock()
	defer e.RUnlock()

	// 计算材料的强化经验值
	exps, err := e.strengthenMaterialsEXP(materials)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}

	// 装备等级最大值上限，经验溢出不保留
	maxLevel := e.maxLevel[target.Rarity]

	// 比较账号等级限制和装备稀有度限制哪个小用哪个
	if maxLevelByTeam < maxLevel {
		maxLevel = maxLevelByTeam
	}

	expNow := target.EXP.Value()

	var addEXP int32 = itemEXP // 强化实际增加的经验

	levelEXPKey := e.makeLevelEXPKey(maxLevel, target.Rarity)

	// 装备最大经验，不可溢出
	maxEXP, ok := e.levelEXP[levelEXPKey]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentLevelUp, levelEXPKey)
	}

	// 装备等级已经可以升到上限
	for _, exp := range exps {
		if expNow+exp > maxEXP {
			// 检查经验溢出了还有后续装备的话就报错
			// if i != len(exps)-1 {
			// 	return 0, errors.Swrapf(common.ErrEquipmentEXPUpToLimit, i)
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
	// limitEXP, err := manager.CSV.Equipment.LevelEXP(equipment.WID, limitLevel)
	// if err != nil {
	// 	return 0, errors.WrapTrace(err)
	// }
	//
	// for i, exp := range exps {
	// 	if expNow+exp > limitEXP {
	// 		// 检查经验溢出了还有后续装备的话就报错
	// 		if i != len(exps)-1 {
	// 			return 0, errors.WrapTrace(common.ErrEquipmentEXPUpToLimit)
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
func (e *EquipmentEntry) strengthenMaterialsEXP(materials []*common.Equipment) ([]int32, error) {
	exps := make([]int32, 0, len(materials))

	for _, equipment := range materials {
		data, ok := e.equipments[equipment.EID]
		if !ok {
			return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, equipment.EID)
		}

		// 升阶通用道具的基本经验为0，不允许作为升级材料
		if data.EXP == 0 {
			return nil, errors.Swrapf(common.ErrEquipmentStrengthenMaterialInvalid, equipment.EID)
		}

		level := equipment.Level.Value() // 装备等级

		var sum float64     // 这个装备强化经验折扣后的和
		var lastEXP float64 // 上一级强化等级经验

		// 计算强化经验的折扣
		for _, val := range e.strengthenDiscount {
			if level >= val.Level {
				// 全额折扣
				levelEXPKey := e.makeLevelEXPKey(val.Level, equipment.Rarity)

				exp, ok := e.levelEXP[levelEXPKey] // 当前等级经验
				if !ok {
					return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentLevelUp, levelEXPKey)
				}

				sum += (float64(exp) - lastEXP) * (val.Discount / 10000)

				lastEXP = float64(exp)
			} else {
				// 部分折扣
				sum += (float64(equipment.EXP.Value()) - lastEXP) * (val.Discount / 10000)
				break
			}
		}

		// 基本经验 + 强化经验
		exps = append(exps, data.EXP+int32(math.Trunc(sum)))
	}

	return exps, nil
}

// // 最高等级经验，用来判断强化经验溢出
// func (e *EquipmentEntry) MaxLevel(equipment *common.Equipment) (int32, error) {
// 	e.RLock()
// 	defer e.RUnlock()
//
// 	return e.maxLevel[equipment.Rarity], nil
// }
//
// // 最高等级经验，用来判断强化经验溢出
// func (e *EquipmentEntry) LevelEXP(equipment *common.Equipment) (int32, error) {
// 	e.RLock()
// 	defer e.RUnlock()
//
// 	exp, ok := e.levelEXP[e.makeLevelEXPKey(equipment)]
// 	if !ok {
// 		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentLevelUp, equipment.Rarity)
// 	}
//
// 	return exp, nil
// }

// func (e *EquipmentEntry) StrengthenGoldCost(equipment *common.Equipment, addEXP int32) (int32, error) {
// 	e.RLock()
// 	defer e.RUnlock()
//
// 	data, ok := e.equipments[equipment.EID]
// 	if !ok {
// 		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, equipment.EID)
// 	}
//
// 	return data.GoldPerExp * addEXP, nil
// }

// 装备强化金币只计算装备固定值，不计算折扣金币
func (e *EquipmentEntry) StrengthenGoldCost(equipment *common.Equipment, itemEXP int32, materials []*common.Equipment) (int32, error) {
	e.RLock()
	defer e.RUnlock()

	data, ok := e.equipments[equipment.EID]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, equipment.EID)
	}

	costGolds := data.GoldPerExp * float64(itemEXP)

	for _, material := range materials {
		md, ok := e.equipments[material.EID]
		if !ok {
			return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, material.EID)
		}

		costGolds += data.GoldPerExp * float64(md.EXP)
	}

	return int32(math.Floor(costGolds)), nil
}

// 随机词条数量
// func (e *EquipmentEntry) RandEntriesNum(eid int32) (int, error) {
// 	e.RLock()
// 	defer e.RUnlock()
//
// 	equipment, ok := e.equipments[eid]
// 	if !ok {
// 		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, eid)
// 	}
//
// 	return len(equipment.RandEntries), nil
// }

// 生成随机词条，返回值 1.属性 2.数值
func (e *EquipmentEntry) randAllAttr(eid int32, num int) ([]common.EquipmentAttr, error) {
	// 最终属性
	allAttrs := make([]common.EquipmentAttr, 0, num)
	allAttrsM := map[int8]bool{}

	for i := 0; i < num; i++ {
		// 临时属性，因为不能重复，又得按顺序随机，所以随机n+1个不同属性，保证有一个不重复的
		attrs, err := e.randAttr(eid, i)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		isAttrOK := false

		log.Printf("allAttrs: %v, attrs: %v allAttrsM: %v", allAttrs, attrs, allAttrsM)

		for _, v := range attrs {

			if !allAttrsM[v.Attr] {
				allAttrs = append(allAttrs, v)
				allAttrsM[v.Attr] = true
				isAttrOK = true
				break
			}
		}

		if !isAttrOK {
			// 属性不够，必须重复
			return nil, errors.Swrapf(common.ErrCSVFormatInvalid, CfgEquipmentRandAttributes, eid)
		}
	}

	return allAttrs, nil
}

// 生成随机词条
func (e *EquipmentEntry) randAttr(eid int32, i int) ([]common.EquipmentAttr, error) {
	equipment, ok := e.equipments[eid]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgEquipmentData, eid)
	}

	if i > len(equipment.RandEntries)-1 {
		return nil, errors.Swrapf(common.ErrCSVFormatInvalid, CfgEquipmentData, eid)
	}

	entry, ok := e.randAttrs[equipment.RandEntries[i]]
	if !ok {
		return nil, errors.Swrapf(common.ErrCSVFormatInvalid, CfgEquipmentRandAttributes, equipment.RandEntries[i])
	}

	// 随机属性 1~16
	attrs := make([]common.EquipmentAttr, 0, i+1)

	attrIndexRs := rand.UniquePerm(i+1, entry.Probs)
	for _, attrIndex := range attrIndexRs {
		attr := attrIndex + 1

		valID := entry.ValueIDs[attrIndex]

		val, ok := e.randAttrsValue[valID]
		if !ok {
			return nil, errors.Swrapf(common.ErrCSVFormatInvalid, CfgEquipmentRandAttributes, valID)
		}

		r := rand.SinglePerm(val.Probs)
		ret := rand.RangeInt32(val.Ranges[r][0], val.Ranges[r][1])

		remain := ret % 100
		if remain > 0 {
			ret -= remain
			ret += 100
		}

		attrs = append(attrs, *common.NewEquipmentAttr(int8(attr), ret))
	}

	// 随机范围
	return attrs, nil
}
