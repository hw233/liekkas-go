package entry

import (
	"math"
	"shared/common"
	"shared/csv/base"
	"shared/csv/static"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/rand"
	"shared/utility/servertime"
	"shared/utility/transfer"
	"sort"
	"sync"
)

/**

产出方面：

产出的随机
 产出的各物品随建筑等级变化
 产出各物品的概率权重随阶级随建筑阶级变化

持续性建筑的产出：
建筑产出 = 单位时间产出* 产出次数
产出次数 = int(已进行产出的时间/单位时间)
单位时间产出还和放入工人加成有关，每个工人放入加百分之7（配表）的产出

非持续性建筑的产出：
建筑产出 = 单位时间产出* 产出次数
产出次数 是开始产出时候选择的
单位时间还和派遣侍从星级有关，每派遣一个侍从减百分之X的单位生产时间（每个角色减的百分比不同，配在cfg_character_star的buildProducePercent中）

*/

const UnitTypeCharacter = 1
const UnitTypePopulation = 2
const MainTowerBuildId = 11

const CfgNameGraveyardBuild = "cfg_graveyard_build"
const CfgNameGraveyardBuildLv = "cfg_graveyard_build_lv"
const CfgNameGraveyardMainTowerLv = "cfg_graveyard_main_tower_lv"
const CfgNameGraveyardBuildStage = "cfg_graveyard_build_stage"
const CfgNameGraveyardProduceBuff = "cfg_graveyard_produce_buff"

type GraveyardEntry struct {
	sync.RWMutex

	buildingMap map[int32]Building

	buildLevelUps                map[int32][]BuildingLvUp
	buildStageUps                map[int32][]BuildingStageUp
	produceBuildLevels           map[int32][]ProduceBuildLevel
	buildingStages               map[int32][]BuildingStage
	characterStarProduces        map[int32][]CharacterStarProduce
	continuousProduceBuildings   map[int32][]ContinuousBuildingLevel
	consumeProduceBuildings      map[int32][]ConsumeBuildingLevel
	produceBuffs                 map[int32]ProduceBuff
	mainTowerLevels              []MainTowerLevel
	mainTowerInitPosition        coordinate.Position
	graveyardPopulationAddition  float64
	graveyardPlotRewardMaxNum    int32
	graveyardPlotRewardHourStart int32
	graveyardPlotRewardHourEnd   int32
	graveyardPlotRewardDropId    int32
	graveyardMinHelpSec          int32
	graveyardSingleHelpCount     int32
	graveyardDailyHelpSendCount  int32
	graveyardMinHelpPercent      float64
}

type CharacterStarProduce struct {
	BuildProducePercent float64 `rule:"tenThousandth"`
}

func NewGraveyardEntry() *GraveyardEntry {
	return &GraveyardEntry{
		buildingMap:                 map[int32]Building{},
		buildLevelUps:               map[int32][]BuildingLvUp{},
		buildStageUps:               map[int32][]BuildingStageUp{},
		produceBuildLevels:          map[int32][]ProduceBuildLevel{},
		buildingStages:              map[int32][]BuildingStage{},
		characterStarProduces:       map[int32][]CharacterStarProduce{},
		continuousProduceBuildings:  map[int32][]ContinuousBuildingLevel{},
		consumeProduceBuildings:     map[int32][]ConsumeBuildingLevel{},
		produceBuffs:                map[int32]ProduceBuff{},
		mainTowerLevels:             []MainTowerLevel{},
		mainTowerInitPosition:       coordinate.Position{},
		graveyardPopulationAddition: 0,
	}
}

type Building struct {
	common.BuildingArea `ignore:"true"`
	Type                int32
	MinProduceCount     int32              //至少生产满多少次后能领取奖励（对持续性产出类的建筑有效）
	UnlockCondition     *common.Conditions `rule:"conditions"`
}

type ProduceBuildLevel struct {
	ProduceTime   int32   // 单位产出时间
	UnitCount     []int32 // 侍从孔位
	LevelUpDropId int32   // 升到本级揭幕奖励
}

type BuildingLvUp struct {
	LevelUpConsume   *common.Rewards `rule:"rewards"` // 升级消耗
	LevelUpTime      int32           // 升级所需时间（秒）
	MainTowerLvUp    `ignore:"true"` // maintower需要的升级条件
	ProduceBuildLvUp `ignore:"true"` // 生产型建筑需要的升级条件
}
type ProduceBuildLvUp struct {
	MainTowerLevelLimit int32 // 主堡等级限制

}

type MainTowerLvUp struct {
	UnlockConditions    *common.Conditions   `ignore:"true"`
	BuildingLevelLimits *BuildingLevelLimits `src:"BuildsCondition" rule:"buildingLevelLimits"` //其他建筑限制

}

// ContinuousBuildingLevel 持续生产型建筑等级成长
type ContinuousBuildingLevel struct {
	StoreLimit int32 // 最多储存

}

// ConsumeBuildingLevel 消耗物品再产出的建筑等级成长
type ConsumeBuildingLevel struct {
	ProduceLimit    int32           // 一次最多生产
	ConsumeResource *common.Rewards `rule:"rewards"` // 每次生产消耗
}

type BuildingStage struct {
	DropId        int32 // 随机掉落
	StageUpDropId int32 //升到本阶揭幕奖励
}

type BuildingStageUp struct {
	MainTowerLevelLimit int32                // 主堡等级限制
	StageUpConsume      *common.Rewards      `rule:"rewards"` // 升级消耗
	StageUpTime         int32                // 升级所需时间（秒）
	BuildingLevelLimits *BuildingLevelLimits `src:"OtherBuildingLvLimit" rule:"buildingLevelLimits"` //其他建筑限制
	ProduceCountLimit   int32                // 升阶所需成功生产道具次数

}

type BuildingLevelLimits struct {
	Limits []BuildingLevelLimit //升阶所需其他建筑等级
}

type BuildingLevelLimit struct {
	BuildId int32
	Lv      int32
	NeedNum int32
}

type MainTowerLevel struct {
	PopulationCount  int32
	BuildsCountLimit map[int32]int32 ` rule:"int32Map"`
	UnlockArea       *common.Area    `src:"UnlockBuildArea" rule:"area"`
}

type ProduceBuff struct {
	Id          int32
	Type        int32
	TypeContent int32
	BuildId     int32
	StageDrop   []int32 `src:"Stage0drop,Stage1drop,Stage2drop,Stage3drop,Stage4drop,Stage5drop" rule:"mergeInt32"`
}

func (g *GraveyardEntry) Check(config *Config) error {

	return nil
}

func (g *GraveyardEntry) Reload(config *Config) error {
	g.Lock()
	defer g.Unlock()
	buildingMap := map[int32]Building{}
	buildLevelUps := map[int32][]BuildingLvUp{}
	buildStageUps := map[int32][]BuildingStageUp{}
	produceBuildLevels := map[int32][]ProduceBuildLevel{}
	buildingStages := map[int32][]BuildingStage{}
	characterStarProduces := map[int32][]CharacterStarProduce{}

	continuousProduceBuildings := map[int32][]ContinuousBuildingLevel{}
	consumeProduceBuildings := map[int32][]ConsumeBuildingLevel{}
	produceBuffs := map[int32]ProduceBuff{}

	var mainTowerLevels []MainTowerLevel
	mainTowerInitPosition := coordinate.Position{}

	for _, v := range config.CfgGraveyardBuildConfig.GetAllData() {
		building := &Building{}
		buildingArea := &common.BuildingArea{}
		err := transfer.Transfer(v, buildingArea)
		if err != nil {
			return errors.WrapTrace(err)
		}
		err = transfer.Transfer(v, building)
		if err != nil {
			return errors.WrapTrace(err)
		}
		building.BuildingArea = *buildingArea
		buildingMap[v.Id] = *building
	}

	// lv配置按lv排序
	var levelConfigs []*base.CfgGraveyardBuildLevel
	for _, v := range config.CfgGraveyardBuildLevelConfig.GetAllData() {
		levelConfigs = append(levelConfigs, v)
	}
	lvLess := func(i, j int) bool {
		return levelConfigs[i].Level < levelConfigs[j].Level
	}
	sort.Slice(levelConfigs, lvLess)

	for _, v := range levelConfigs {

		// 升级条件
		lvUps, ok := buildLevelUps[v.BuildId]
		if !ok {
			lvUps = []BuildingLvUp{}
		}
		lvUp := &BuildingLvUp{}
		err := transfer.Transfer(v, lvUp)
		if err != nil {
			return errors.WrapTrace(err)
		}
		produceBuildLvUp := &ProduceBuildLvUp{}
		err = transfer.Transfer(v, produceBuildLvUp)
		if err != nil {
			return errors.WrapTrace(err)
		}
		lvUp.ProduceBuildLvUp = *produceBuildLvUp
		lvUps = append(lvUps, *lvUp)
		buildLevelUps[v.BuildId] = lvUps

		// 生产型建筑的level各属性
		if v.Level > 0 {
			produceBuildLvs, ok := produceBuildLevels[v.BuildId]
			if !ok {
				produceBuildLvs = []ProduceBuildLevel{}
			}
			produceBuildLevel := &ProduceBuildLevel{}
			err = transfer.Transfer(v, produceBuildLevel)
			if err != nil {
				return errors.WrapTrace(err)
			}
			produceBuildLvs = append(produceBuildLvs, *produceBuildLevel)
			produceBuildLevels[v.BuildId] = produceBuildLvs
			// 手动开始建筑
			if buildingMap[v.BuildId].Type == static.GraveyardBuildTypeConsumeProduceItem {
				consumeProduceBuilding, ok := consumeProduceBuildings[v.BuildId]
				if !ok {
					consumeProduceBuilding = []ConsumeBuildingLevel{}
				}
				consumeBuildingLevel := &ConsumeBuildingLevel{}
				err = transfer.Transfer(v, consumeBuildingLevel)
				if err != nil {
					return errors.WrapTrace(err)
				}
				consumeProduceBuilding = append(consumeProduceBuilding, *consumeBuildingLevel)
				consumeProduceBuildings[v.BuildId] = consumeProduceBuilding
			}

			// 持续性产出建筑
			if buildingMap[v.BuildId].Type == static.GraveyardBuildTypeContinuous {
				continuousProduceBuilding, ok := continuousProduceBuildings[v.BuildId]
				if !ok {
					continuousProduceBuilding = []ContinuousBuildingLevel{}
				}
				continuousBuildingLevel := &ContinuousBuildingLevel{}
				err = transfer.Transfer(v, continuousBuildingLevel)
				if err != nil {
					return errors.WrapTrace(err)
				}
				continuousProduceBuilding = append(continuousProduceBuilding, *continuousBuildingLevel)
				continuousProduceBuildings[v.BuildId] = continuousProduceBuilding
			}

		}

	}
	// stage配置按stage排序

	var stageConfigs []*base.CfgGraveyardBuildStage
	for _, v := range config.CfgGraveyardBuildStageConfig.GetAllData() {
		stageConfigs = append(stageConfigs, v)

	}
	stageLess := func(i, j int) bool {
		return stageConfigs[i].Stage < stageConfigs[j].Stage
	}
	sort.Slice(stageConfigs, stageLess)

	tmpDrop := NewDrop()
	err := tmpDrop.Reload(config)
	if err != nil {
		return errors.WrapTrace(err)
	}

	for _, v := range stageConfigs {
		// 0阶数据不读
		if v.Stage > 0 {
			stageUps, ok := buildStageUps[v.BuildId]
			if !ok {
				stageUps = []BuildingStageUp{}
			}
			stageUp := &BuildingStageUp{}
			err := transfer.Transfer(v, stageUp)
			if err != nil {
				return errors.WrapTrace(err)
			}
			stageUps = append(stageUps, *stageUp)
			buildStageUps[v.BuildId] = stageUps
		}
		stages, ok := buildingStages[v.BuildId]
		if !ok {
			stages = []BuildingStage{}
		}
		stage := &BuildingStage{}
		err := transfer.Transfer(v, stage)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// 检查手动型产出建筑的产出配置
		if buildingMap[v.BuildId].Type == static.GraveyardBuildTypeConsumeProduceItem {
			err := g.checkConsumeBuildDrop(tmpDrop, stage.DropId)
			if err != nil {
				return errors.WrapTrace(err)
			}
		}
		stages = append(stages, *stage)
		buildingStages[v.BuildId] = stages
	}

	for _, v := range config.CfgGraveyardProduceBuffConfig.GetAllData() {
		produceBuff := &ProduceBuff{}
		err := transfer.Transfer(v, produceBuff)
		if err != nil {
			return errors.WrapTrace(err)
		}
		produceBuffs[v.Id] = *produceBuff
	}

	// 主堡配置按等级排序
	var mainTowerLvConfigs []*base.CfgGraveyardMainTowerLevel
	for _, v := range config.CfgGraveyardMainTowerLevelConfig.GetAllData() {
		mainTowerLvConfigs = append(mainTowerLvConfigs, v)
	}
	mainTowerLvLess := func(i, j int) bool {
		return mainTowerLvConfigs[i].Id < mainTowerLvConfigs[j].Id
	}
	sort.Slice(mainTowerLvConfigs, mainTowerLvLess)

	oldBuildLimitCount := map[int32]int32{}
	oldUnlockArea := common.NewArea()
	for _, v := range mainTowerLvConfigs {
		mainTowerLevel := &MainTowerLevel{}
		err := transfer.Transfer(v, mainTowerLevel)
		if err != nil {
			return errors.WrapTrace(err)
		}
		limit := mainTowerLevel.BuildsCountLimit
		mergeBuildLimitCount(limit, oldBuildLimitCount)
		mainTowerLevel.UnlockArea.MergeArea(oldUnlockArea)
		mainTowerLevels = append(mainTowerLevels, *mainTowerLevel)
		oldBuildLimitCount = limit
		oldUnlockArea = mainTowerLevel.UnlockArea

		mainTowerLvUp := &MainTowerLvUp{
			UnlockConditions: common.NewConditions(),
		}
		err = transfer.Transfer(v, mainTowerLvUp)
		if err != nil {
			return errors.WrapTrace(err)
		}

		if v.ExploreCondition > 0 {
			condition := common.NewCondition(static.ConditionTypePassLevel, v.ExploreCondition)
			mainTowerLvUp.UnlockConditions.AddCondition(condition)
		}

		mainTowerConfig, ok := buildLevelUps[MainTowerBuildId]
		if ok {
			mainTowerConfig[v.Id-1].MainTowerLvUp = *mainTowerLvUp

		}

	}

	for _, star := range config.CfgCharacterStarConfig.GetAllData() {
		if _, ok := characterStarProduces[star.CharID]; !ok {
			characterStarProduces[star.CharID] = make([]CharacterStarProduce, CharacterStarsCount, CharacterStarsCount)
		}

		starIndex := int(star.Star)
		if starIndex > len(characterStarProduces[star.CharID]) {
			return errors.Swrapf(common.ErrCSVFormatInvalid, CfgCharacter, star)
		}

		err := transfer.Transfer(star, &characterStarProduces[star.CharID][starIndex-1])
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	mainTowerInitPosition = config.GlobalEntry.MainTowerPosition

	g.buildingMap = buildingMap
	g.buildLevelUps = buildLevelUps
	g.buildStageUps = buildStageUps
	g.buildingStages = buildingStages
	g.characterStarProduces = characterStarProduces
	g.produceBuildLevels = produceBuildLevels
	g.continuousProduceBuildings = continuousProduceBuildings
	g.consumeProduceBuildings = consumeProduceBuildings
	g.produceBuffs = produceBuffs
	g.mainTowerLevels = mainTowerLevels
	g.mainTowerInitPosition = mainTowerInitPosition
	g.graveyardPopulationAddition = float64(config.GlobalEntry.GraveyardPopulationAddition) / 10000
	g.graveyardPlotRewardMaxNum = config.GlobalEntry.GraveyardPlotRewardMaxNum
	g.graveyardPlotRewardHourStart = config.GlobalEntry.GraveyardPlotRewardHour[0]
	g.graveyardPlotRewardHourEnd = config.GlobalEntry.GraveyardPlotRewardHour[1]
	g.graveyardPlotRewardDropId = config.GlobalEntry.GraveyardPlotRewardDropId
	g.graveyardMinHelpSec = config.GlobalEntry.GraveyardMinHelpSec
	g.graveyardSingleHelpCount = config.GlobalEntry.GraveyardSingleHelpCount
	g.graveyardDailyHelpSendCount = config.GlobalEntry.GraveyardDailyHelpSendCount
	g.graveyardMinHelpPercent = float64(config.GlobalEntry.GraveyardMinHelpPercent) / 10000
	return nil
}

// 消耗型建筑drop检查
func (g *GraveyardEntry) checkConsumeBuildDrop(drop *Drop, dropId int32) error {

	info, ok := drop.drops[dropId]
	if !ok {
		return errors.Swrapf(common.ErrNotFoundInCSV, CfgDropDataConfig, dropId)
	}
	if info.Type != dropTypeRelative {
		return errors.New("graveyard drop type error dropId:%d", dropId)
	}
	if info.Count != 1 {
		return errors.New("graveyard drop count error dropId:%d", dropId)
	}
	for _, detail := range info.Details {
		if detail.Up != 1 || detail.Low != 1 {
			return errors.New("graveyard drop Up Low error dropId:%d", dropId)
		}
	}
	return nil
}
func mergeBuildLimitCount(old, new map[int32]int32) {
	for k, v := range new {
		i, ok := old[k]
		if ok {
			old[k] = v + i
		} else {
			old[k] = v
		}
	}

}
func (g *GraveyardEntry) MainTowerInitPosition() *coordinate.Position {
	g.RLock()
	defer g.RUnlock()
	return &g.mainTowerInitPosition
}

func (g *GraveyardEntry) PopulationAddition() float64 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardPopulationAddition
}

func (g *GraveyardEntry) GraveyardPlotRewardMaxNum() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardPlotRewardMaxNum
}
func (g *GraveyardEntry) GraveyardPlotRewardHourStart() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardPlotRewardHourStart
}
func (g *GraveyardEntry) GraveyardPlotRewardHourEnd() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardPlotRewardHourEnd
}

func (g *GraveyardEntry) GraveyardPlotRewardDropId() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardPlotRewardDropId
}

func (g *GraveyardEntry) GetGraveyardMinHelpSec() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardMinHelpSec
}
func (g *GraveyardEntry) GetGraveyardSingleHelpCount() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardSingleHelpCount
}
func (g *GraveyardEntry) GetGraveyardDailyHelpSendCount() int32 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardDailyHelpSendCount
}
func (g *GraveyardEntry) GetGraveyardMinHelpPercent() float64 {
	g.RLock()
	defer g.RUnlock()
	return g.graveyardMinHelpPercent
}

// GetBuildById 获得建筑type和区域
func (g *GraveyardEntry) GetBuildById(buildId int32) (*Building, error) {
	g.RLock()
	defer g.RUnlock()
	building, ok := g.buildingMap[buildId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuild, buildId)
	}
	return &building, nil
}

// GetBuildType 获得建筑type
func (g *GraveyardEntry) GetBuildType(buildId int32) int32 {
	g.RLock()
	defer g.RUnlock()
	building, ok := g.buildingMap[buildId]
	if !ok {
		glog.Errorf("GetBuildType err:%v", errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuild, buildId))
		return 0
	}
	return building.Type
}

// CalUnlockArea 计算剩余面积
func (g *GraveyardEntry) CalUnlockArea(mainTowerLv int32, buildPositionListMap map[int32][]*coordinate.Position) (*common.Area, error) {
	g.RLock()
	defer g.RUnlock()
	if mainTowerLv == 0 {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardMainTowerLv, mainTowerLv)
	}
	if int(mainTowerLv) > len(g.mainTowerLevels) {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardMainTowerLv, mainTowerLv)
	}
	mainTowerLevel := g.mainTowerLevels[mainTowerLv-1]
	area := common.NewArea()
	area.MergeArea(mainTowerLevel.UnlockArea)

	for buildId, positions := range buildPositionListMap {
		buildingArea := g.buildingMap[buildId].BuildingArea
		for _, position := range positions {
			err := area.MinusBuildingArea(buildingArea, position)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
		}

	}
	return area, nil
}

// CanBuildCount 当前主堡等级最多可建造buildId的建筑个数
func (g *GraveyardEntry) CanBuildCount(mainTowerLevel, buildId int32) int32 {
	g.RLock()
	defer g.RUnlock()
	if mainTowerLevel == 0 {
		return 0
	}
	if len(g.mainTowerLevels) < int(mainTowerLevel) {
		return 0
	}
	num, ok := g.mainTowerLevels[mainTowerLevel-1].BuildsCountLimit[buildId]
	if ok {
		return num
	}
	return 0
}

// GetCreate 获得建造相关条件
func (g *GraveyardEntry) GetCreate(buildId int32) (*BuildingLvUp, error) {
	return g.GetLvUp(buildId, 1)
}

// GetLvUp 获得升级各配置
func (g *GraveyardEntry) GetLvUp(buildId, lv int32) (*BuildingLvUp, error) {
	g.RLock()
	defer g.RUnlock()
	levels, ok := g.buildLevelUps[buildId]
	if !ok {
		return nil, errors.Swrapf(common.ErrGraveyardBuildCanNotLvUp, buildId, lv)
	}
	if lv <= 0 {
		return nil, errors.Swrapf(common.ErrGraveyardBuildCanNotLvUp, buildId, lv)
	}
	if int(lv) > len(levels) {

		return nil, errors.Swrapf(common.ErrGraveyardBuildMaxLv, buildId, lv)
	}
	return &levels[lv-1], nil

}

// GetStageUp 获得升阶各配置
func (g *GraveyardEntry) GetStageUp(buildId, stage int32) (*BuildingStageUp, error) {
	g.RLock()
	defer g.RUnlock()
	stages, ok := g.buildStageUps[buildId]
	if !ok {
		return nil, errors.Swrapf(common.ErrGraveyardBuildCanNotStageUp, buildId, stage)
	}
	if stage <= 0 {
		return nil, errors.Swrapf(common.ErrGraveyardBuildCanNotStageUp, buildId, stage)
	}
	if int(stage) > len(stages) {
		return nil, errors.Swrapf(common.ErrGraveyardBuildMaxStage, buildId, stage)
	}
	return &stages[stage-1], nil
}

func (g *GraveyardEntry) GetCurtainDrop(build *common.UserGraveyardBuild) int32 {
	g.RLock()
	defer g.RUnlock()
	//
	switch build.CurtainType {
	case common.CurtainTypeStageUp:
		stages, ok := g.buildingStages[build.BuildId]
		if !ok {
			return 0
		}
		if int(build.Stage) > len(stages) {
			return 0
		}
		return stages[build.Stage].StageUpDropId
	case common.CurtainTypeCreate, common.CurtainTypeLvUp:
		levels, ok := g.produceBuildLevels[build.BuildId]
		if !ok {
			return 0
		}
		if int(build.Lv) > len(levels) {
			return 0
		}
		return levels[build.Lv-1].LevelUpDropId
	default:
		return 0
	}

}

func (g *GraveyardEntry) GetUnitCount(buildId, lv, unitType int32) int32 {
	g.RLock()
	defer g.RUnlock()
	/**
	数据格式:1,m及2,m
	1.侍从单位
	2.工人单位
	*/
	levels, ok := g.produceBuildLevels[buildId]
	if !ok {
		return 0
	}
	if lv <= 0 {
		return 0
	}
	if int(lv) > len(levels) {
		return 0
	}
	level := levels[lv-1]
	if len(level.UnitCount) == 2 {
		if level.UnitCount[0] == unitType {
			return level.UnitCount[1]
		}
	}
	return 0
}

func (g *GraveyardEntry) CanPopulationSet(buildId, level, population int32) error {
	if population < 0 {
		return errors.WrapTrace(common.ErrParamError)
	}
	buildingType := g.GetBuildType(buildId)
	if !(buildingType == static.GraveyardBuildTypeContinuous || buildingType == static.GraveyardBuildTypeConsumeProduceItem) {
		return errors.WrapTrace(common.ErrParamError)
	}

	maxCount := g.GetUnitCount(buildId, level, UnitTypePopulation)
	if maxCount < population {
		return errors.Swrapf(common.ErrGraveyardPopulationNumOverflow, population, maxCount)
	}
	return nil
}

func (g *GraveyardEntry) CanCharacterDispatch(buildId, level int32, characterPosition *common.CharacterPositions) error {
	count := g.GetUnitCount(buildId, level, UnitTypeCharacter)
	for _, v := range characterPosition.Positions {
		if v > count {
			return errors.Swrapf(common.ErrGraveyardCharacterNumOverflow, v, count)
		}
	}

	return nil
}

func (m *Manager) RefreshOutPut(build *common.UserGraveyardBuild, starMap map[int32]int32, manualFinish bool) {
	buildingType := m.GraveyardEntry.GetBuildType(build.BuildId)
	// 判断type
	if buildingType == static.GraveyardBuildTypeMaintower || buildingType == static.GraveyardBuildTypeDecoration {
		return
	}
	// 手动产出型建筑只在前端调获得产出接口（GraveyardProductionGet）时候刷新产出
	if buildingType == static.GraveyardBuildTypeConsumeProduceItem && !manualFinish {
		return
	}
	// 未建成的建筑不可刷新
	if !build.BuildComplete() {
		return
	}
	// 持续型产出，达到storeLimit
	if build.ReachLimit {
		// 到达生产上限后不再累计产生的余数秒数
		build.ProduceStartAt = servertime.Now().Unix()
		return

	}
	produceTime, err := m.GraveyardEntry.GetProduceTime(build, starMap)
	if err != nil {
		glog.Errorf("RefreshOutPut GetProduceTime err:%+v", err)
		return
	}

	dropNum, produceDiff, err := m.GraveyardEntry.calDropNumTickTime(build, produceTime)
	if err != nil {
		glog.Errorf("RefreshOutPut calDropNumTickTime, err:%+v", err)
		return
	}
	storeLimit := m.GraveyardEntry.getStoreLimit(build.BuildId, build.FetchLevel())

	var dropCount int32
	var reachLimit bool
	for _, dropId := range *dropNum {
		rewards, err := m.Drop.DropRewards(dropId)
		if err != nil {
			glog.Errorf("RefreshOutPut DropRewards, err:%+v", err)
			continue
		}
		glog.Debugf("RefreshOutPut before limit,buildID: %v, dropID: %v, rewards: %v, productions: %v", build.BuildId, dropId, rewards, build.Productions)
		reachLimit = productionsAddWithStoreLimit(buildingType, storeLimit, build.Productions, rewards)
		dropCount++
		if reachLimit {
			break
		}

	}
	// 设置build 生产状态
	switch buildingType {
	case static.GraveyardBuildTypeContinuous:
		// produceDiff 这次生产完了后产生的余数秒数
		build.ProduceStartAt = servertime.Now().Unix() - produceDiff
		build.CurrProduceNum += dropCount
		build.TotalProduceNum += dropCount
		build.ReachLimit = reachLimit
	case static.GraveyardBuildTypeConsumeProduceItem:
		// 产出型建筑按件收取
		build.CurrReceiveProduceNum += dropCount
		build.TotalProduceNum += dropCount
	}

}

// 收取产出的时候每次判断是否达到存储上限
func productionsAddWithStoreLimit(buildType, storeLimit int32, productions *common.Rewards, rewards *common.Rewards) bool {
	// 所有产出建筑都有存储上限
	if buildType == static.GraveyardBuildTypeContinuous {
		totalNum := calTotalNum(productions)

		if totalNum >= storeLimit {
			return true
		}

		addNum := calTotalNum(rewards)
		if totalNum+addNum >= storeLimit {
			for _, reward := range rewards.Value() {
				if totalNum+reward.Num > storeLimit {
					tmp := storeLimit - totalNum
					if tmp > 0 {
						productions.AddReward(common.NewReward(reward.ID, tmp))
						totalNum += tmp
					}
					// 达到存储上限break
					break
				} else {
					productions.AddReward(&reward)
					totalNum += reward.Num
				}
			}
			return true
		} else {
			productions.AddRewards(rewards)
			totalNum += addNum
		}

	} else {
		productions.AddRewards(rewards)
	}
	return false
}

func ToCeilRewards(rewardFloatNum map[int32]float64) *common.Rewards {
	result := common.NewRewards()

	for id, num := range rewardFloatNum {
		result.AddReward(common.NewReward(id, int32(math.Ceil(num))))
	}
	return result
}

func calTotalNumFloatCeil(rewardFloatNum map[int32]float64) int32 {
	var total int32
	for _, num := range rewardFloatNum {
		total += int32(math.Ceil(num))
	}
	return total

}

func calTotalNum(reward *common.Rewards) int32 {
	var total int32
	for _, rs := range reward.MultiValue() {
		for _, r := range rs {
			total += r.Num
		}
	}
	return total

}

// 结合正在使用的buff获得所有每次获得单位产出
func (g *GraveyardEntry) calDropNumTickTime(build *common.UserGraveyardBuild, produceTime float64) (*DropNum, int64, error) {
	buildStage, err := g.getBuildStage(build.BuildId, build.FetchStage())
	if err != nil {
		return nil, 0, err
	}
	normalDrop := buildStage.DropId
	// 每个单位产出都有产出时间点，如果这个时间点在buff的使用时间内的话就享受buff加成
	dropNum, err := g.tickWithBuff(build, produceTime, normalDrop)
	if err != nil {
		return nil, 0, err
	}
	total := build.CalProduceSecs()
	totalNum := int32(float64(total) / produceTime)
	return dropNum, int64(math.Ceil(float64(total) - float64(totalNum)*produceTime)), nil
}

// DropNumFloat 区别于 DropNum是带小数的
type DropNumFloat map[int32]float64

func NewDropNumFloat() *DropNumFloat {
	return (*DropNumFloat)(&map[int32]float64{})
}
func (d *DropNumFloat) isEmpty() bool {
	return len(*d) == 0
}

func (d *DropNumFloat) add(dropId int32, addNum float64) {
	dropNum, ok := (*d)[dropId]
	if ok {
		(*d)[dropId] = dropNum + addNum
	} else {
		(*d)[dropId] = addNum
	}
}

// 使用道具直接获得产出
func (g *GraveyardEntry) calDropNumNow(build *common.UserGraveyardBuild, produceTime int32, sec int32) (*DropNumFloat, error) {
	buildStage, err := g.getBuildStage(build.BuildId, build.FetchStage())
	if err != nil {
		return nil, err
	}
	normalDrop := buildStage.DropId
	num := float64(sec) / float64(produceTime)
	dropNum := NewDropNumFloat()
	dropNum.add(normalDrop, num)
	return dropNum, nil
}

// TimeSegments 时间段计算
type TimeSegments struct {
	segments []TimeSegment
}

func NewTimeSegments(produce *common.UserGraveyardProduce) *TimeSegments {
	segments := &TimeSegments{}
	startAt := produce.ProduceStartAt

	for _, record := range produce.AccRecords.Records {
		segments.add(NewNormalTimeSegment(startAt, int32(record.AccAt-startAt)))
		segments.add(NewAccTimeSegment(record.AccAt, record.AccSec))
		startAt = record.AccAt
	}
	now := servertime.Now().Unix()
	if startAt < now {
		segments.add(NewNormalTimeSegment(startAt, int32(now-startAt)))
	}
	return segments
}

func (t TimeSegments) getTickTime(extendSec int32) int64 {
	for _, segment := range t.segments {
		if segment.getLen() > extendSec {
			return segment.getTime(extendSec)
		} else {
			extendSec -= segment.getLen()
			continue
		}
	}

	return servertime.Now().Unix()
}

func (t *TimeSegments) add(segment TimeSegment) {
	t.segments = append(t.segments, segment)
}

type TimeSegment interface {
	getTime(extendSec int32) int64
	getLen() int32
}

type NormalTimeSegment struct {
	startAt int64
	Sec     int32
}

func NewNormalTimeSegment(startAt int64, Sec int32) *NormalTimeSegment {
	return &NormalTimeSegment{
		startAt: startAt,
		Sec:     Sec,
	}
}
func (n *NormalTimeSegment) getTime(extendSec int32) int64 {
	return n.startAt + int64(extendSec)
}

func (n *NormalTimeSegment) getLen() int32 {
	return n.Sec
}

type AccTimeSegment struct {
	startAt int64
	Sec     int32
}

func NewAccTimeSegment(startAt int64, Sec int32) *AccTimeSegment {
	return &AccTimeSegment{
		startAt: startAt,
		Sec:     Sec,
	}
}
func (a *AccTimeSegment) getTime(extendSec int32) int64 {
	return a.startAt
}

func (a *AccTimeSegment) getLen() int32 {
	return a.Sec
}

func (g *GraveyardEntry) tickWithBuff(build *common.UserGraveyardBuild, produceTime float64, normalDrop int32) (*DropNum, error) {
	dropNumMap := NewDropNum()
	total := build.CalProduceSecs()
	totalNum := int32(float64(total) / produceTime)

	if g.GetBuildType(build.BuildId) == static.GraveyardBuildTypeConsumeProduceItem {
		if totalNum >= build.CurrProduceNum {
			// 手动生产最多领取，开始生产时候选择的件数
			totalNum = build.CurrProduceNum
		}
	}

	timeSegments := NewTimeSegments(build.UserGraveyardProduce)
	for i := build.CurrReceiveProduceNum + 1; i <= totalNum; i++ {
		tickTime := timeSegments.getTickTime(int32(float64(i) * produceTime))
		// 此时是否有buff加成
		duringBuff, ok := build.DuringBuffWhen(tickTime)

		drop := normalDrop
		// 次数类buff如果在生产中使用，则下次才生效
		if ok && duringBuff.IsEffective(build) {
			buffDrop, err := g.getBuffDrop(duringBuff.GetBuffId(), build.FetchStage())
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			drop = buffDrop
		}
		dropNumMap.add(drop)
	}
	return dropNumMap, nil
}

type DropNum []int32

func NewDropNum() *DropNum {
	return (*DropNum)(&[]int32{})
}
func (d *DropNum) isEmpty() bool {
	return len(*d) == 0
}

func (d *DropNum) add(dropId int32) {
	*d = append(*d, dropId)
}
func (d *DropNum) calDropCount() int32 {

	return int32(len(*d))
}

// GetProduceBaseTime 计算建筑单位产出时间(不算工人和侍从加成)
func (g *GraveyardEntry) GetProduceBaseTime(build *common.UserGraveyardBuild) (int32, error) {
	g.RLock()
	defer g.RUnlock()
	buildLevels, ok := g.produceBuildLevels[build.BuildId]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, build.BuildId)
	}
	lv := build.FetchLevel()
	if lv <= 0 {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, build.BuildId)
	}
	if int(lv) > len(buildLevels) {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, build.BuildId)
	}
	produceTime := buildLevels[lv-1].ProduceTime

	return produceTime, nil
}

// GetProduceTime 计算建筑单位产出时间
func (g *GraveyardEntry) GetProduceTime(build *common.UserGraveyardBuild, starMap map[int32]int32) (float64, error) {
	g.RLock()
	defer g.RUnlock()
	buildLevels, ok := g.produceBuildLevels[build.BuildId]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, build.BuildId)
	}
	lv := build.FetchLevel()
	if lv <= 0 {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, build.BuildId)
	}
	if int(lv) > len(buildLevels) {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, build.BuildId)
	}
	produceTime := buildLevels[lv-1].ProduceTime
	var produceTimeReducePercent float64
	// characters 减少produceTime
	for characterId, star := range starMap {
		produceTimeReducePercent += g.getCharacterProduceTimeReduce(characterId, star)
	}
	produceTimeReducePercent += g.graveyardPopulationAddition * float64(build.PopulationCount)
	return (1.0 - produceTimeReducePercent) * float64(produceTime), nil
}

func (g *GraveyardEntry) getContinuousBuildLv(buildId, lv int32) (*ContinuousBuildingLevel, error) {
	g.RLock()
	defer g.RUnlock()
	continuousBuilding, ok := g.continuousProduceBuildings[buildId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, buildId)
	}
	if lv <= 0 {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, buildId)
	}
	if int(lv) > len(continuousBuilding) {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, buildId)
	}

	return &continuousBuilding[lv-1], nil
}

func (g *GraveyardEntry) getConsumeBuildLv(buildId, lv int32) (*ConsumeBuildingLevel, error) {
	g.RLock()
	defer g.RUnlock()

	consumeBuilding, ok := g.consumeProduceBuildings[buildId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, buildId)
	}
	if lv <= 0 {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, buildId)
	}
	if int(lv) > len(consumeBuilding) {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildLv, buildId)
	}
	return &consumeBuilding[lv-1], nil
}

func (g *GraveyardEntry) getBuildStage(buildId, stage int32) (*BuildingStage, error) {
	g.RLock()
	defer g.RUnlock()

	buildStages, ok := g.buildingStages[buildId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildStage, buildId)
	}
	if stage < 0 {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildStage, buildId)
	}
	if int(stage) >= len(buildStages) {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardBuildStage, buildId)
	}
	return &buildStages[stage], nil
}

func (g *GraveyardEntry) CheckProduceNum(build *common.UserGraveyardBuild, num int32) (*common.Rewards, error) {
	consumeBuildLv, err := g.getConsumeBuildLv(build.BuildId, build.FetchLevel())
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if consumeBuildLv.ProduceLimit < num {
		return nil, common.ErrGraveyardProduceNumLimit
	}
	consume := common.NewRewards()
	consume.AddRewards(consumeBuildLv.ConsumeResource)
	multiple := consume.Multiple(num)

	return multiple, nil
}

func (g *GraveyardEntry) getCharacterProduceTimeReduce(characterId, star int32) float64 {
	g.RLock()
	defer g.RUnlock()

	stars, ok := g.characterStarProduces[characterId]
	if !ok {
		return 0
	}
	if star <= 0 {
		return 0
	}
	if int(star) > CharacterStarsCount {
		return 0
	}
	return stars[star-1].BuildProducePercent
}

// Acc 加速建筑，返回实际消耗
func (m *Manager) Acc(build *common.UserGraveyardBuild, items map[*GraveyardAccItem]int32, characterStar map[int32]int32) (*common.Rewards, error) {
	building, err := m.GraveyardEntry.GetBuildById(build.BuildId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	consume := common.NewRewards()

	for item, num := range items {
		switch item.Acc.AccType {
		case static.GraveyardAccelerateTypeBuildAcc:
			consume.AddRewards(m.GraveyardEntry.buildAcc(build, item, num))
		case static.GraveyardAccelerateTypeReduceProduceTime:

			consume.AddRewards(m.GraveyardEntry.reduceProduceTimeAcc(build, item, building.Type, num, characterStar))

		}
	}
	return consume, nil
}

func (g *GraveyardEntry) buildAcc(build *common.UserGraveyardBuild, item *GraveyardAccItem, num int32) *common.Rewards {
	consume := common.NewRewards()
	if build.UserGraveyardTransition == nil {
		return consume
	}
	now := servertime.Now().Unix()
	ceil := math.Ceil(float64(build.EndAt-now) / float64(item.Acc.Sec))
	useNum := int32(math.Min(ceil, float64(num)))
	if useNum <= 0 {
		return consume
	}
	build.BuildAcc(item.Acc.Sec * useNum)

	consume.AddReward(common.NewReward(item.Id, useNum))
	return consume
}

func (m *Manager) productAcc(build *common.UserGraveyardBuild, item *GraveyardAccItem, buildType, num int32) *common.Rewards {
	consume := common.NewRewards()
	if buildType != static.GraveyardBuildTypeContinuous {
		return consume
	}
	nowStoreNum := calTotalNum(build.Productions)
	storeLimit := m.GraveyardEntry.getStoreLimit(build.BuildId, build.FetchLevel())
	// 达到存储上限
	if storeLimit <= nowStoreNum {
		return consume
	}

	buildStage, err := m.GraveyardEntry.getBuildStage(build.BuildId, build.Stage)
	if err != nil {
		glog.Errorf("productAcc getBuildStage err:%+v", errors.Format(err))
		return consume
	}
	dropId := buildStage.DropId
	onceDrop, err := m.Drop.DropRewards(dropId)
	if err != nil {
		glog.Errorf("productAcc DropRewards err:%+v", errors.Format(err))
		return consume
	}
	// 最多消耗
	ceil := math.Ceil(float64((storeLimit)-nowStoreNum) / float64(calTotalNum(onceDrop)))
	useNum := int32(math.Min(ceil, float64(num)))
	build.AccRecords.AddRecord(item.Acc.Sec*useNum, servertime.Now().Unix())
	consume.AddReward(common.NewReward(item.Id, useNum))
	return consume
}

func (g *GraveyardEntry) reduceProduceTimeAcc(build *common.UserGraveyardBuild, item *GraveyardAccItem, buildType, num int32, star map[int32]int32) *common.Rewards {
	consume := common.NewRewards()
	if buildType != static.GraveyardBuildTypeConsumeProduceItem {
		return consume
	}
	if build.CurrProduceNum == 0 {
		return consume
	}
	produceTime, err := g.GetProduceTime(build, star)
	if err != nil {
		glog.Errorf("reduceProduceTimeAcc GetProduceTime err:%+v", errors.Format(err))
		return consume
	}
	// 已生产完成
	needSec := int32(math.Floor(float64(build.CurrProduceNum) * produceTime))
	if needSec <= build.CalProduceSecs() {
		return consume
	}
	// 最多消耗
	ceil := math.Ceil(float64(needSec) / float64(item.Acc.Sec))
	useNum := int32(math.Min(ceil, float64(num)))
	if useNum <= 0 {
		return consume
	}
	build.AccRecords.AddRecord(item.Acc.Sec*useNum, servertime.Now().Unix())
	consume.AddReward(common.NewReward(item.Id, useNum))
	return consume
}

func (g *GraveyardEntry) getStoreLimit(buildId int32, level int32) int32 {
	continuousBuildLv, err := g.getContinuousBuildLv(buildId, level)
	if err != nil {
		return 0
	} else {
		return continuousBuildLv.StoreLimit
	}

}

// RewardsWithLowUp 如果小于low或者大于up 则等比缩放奖励
func (g *GraveyardEntry) RewardsWithLowUp(rewards *common.Rewards, low int32, up int32) *common.Rewards {
	totalNum := calTotalNum(rewards)
	if totalNum == 0 {
		return rewards
	}
	var ratio float64
	if totalNum < low {
		ratio = float64(low) / float64(totalNum)
		multiple := rewards.MultipleFloor(ratio)
		tmp := calTotalNum(multiple)
		if tmp < low {
			multiple.AddReward(common.NewReward(rewards.Value()[0].ID, low-tmp))
		}
		return multiple

	} else if totalNum > up {
		ratio = float64(up) / float64(totalNum)
		multiple := rewards.MultipleFloor(ratio)
		tmp := calTotalNum(multiple)

		if tmp < up {
			multiple.AddReward(common.NewReward(rewards.Value()[0].ID, up-tmp))
		}
		return multiple
	}
	return rewards
}

func (g *GraveyardEntry) GetBuffConfig(buffId int32) (*ProduceBuff, error) {
	g.RLock()
	defer g.RUnlock()
	buff, ok := g.produceBuffs[buffId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardProduceBuff, buffId)
	}
	return &buff, nil
}

func (g *GraveyardEntry) getBuffDrop(buffId, stage int32) (int32, error) {
	config, err := g.GetBuffConfig(buffId)
	if err != nil {
		return 0, err
	}
	if stage < 0 {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardProduceBuff, buffId)
	}
	if int(stage) >= len(config.StageDrop) {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameGraveyardProduceBuff, buffId)
	}
	return config.StageDrop[stage], nil
}

func (g *GraveyardEntry) CanProductionGet(build *common.UserGraveyardBuild) bool {
	building, err := g.GetBuildById(build.BuildId)
	if err != nil {
		return false
	}
	if building.Type == static.GraveyardBuildTypeConsumeProduceItem {
		return !build.Productions.IsEmpty()
	} else if building.Type == static.GraveyardBuildTypeContinuous {
		return build.CurrProduceNum >= building.MinProduceCount
	}
	return false
}

// CheckConsumeProduceInProduce 非可持续建筑在生产中
func (g *GraveyardEntry) CheckConsumeProduceInProduce(build *common.UserGraveyardBuild) error {
	if g.GetBuildType(build.BuildId) == static.GraveyardBuildTypeConsumeProduceItem {
		// 非可持续建筑在生产中
		if build.CurrProduceNum > 0 {
			return errors.WrapTrace(common.ErrGraveyardBuildInProduct)
		}
	}

	return nil
}

// CheckInNormalState 建筑是否不在生产且不在升级升阶建造中
func (g *GraveyardEntry) CheckInNormalState(build *common.UserGraveyardBuild) error {
	err := g.CheckConsumeProduceInProduce(build)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 未揭幕
	if build.UserGraveyardTransition != nil {
		return errors.WrapTrace(common.ErrGraveyardBuildInTransaction)
	}
	return nil
}

// CalCharacterPositionMap 入驻角色替换
func (g *GraveyardEntry) CalCharacterPositionMap(buildings map[int64]*common.UserGraveyardBuild, BuildUid int64, targetBuilding *common.UserGraveyardBuild, characters *common.CharacterPositions) (map[int64]*common.CharacterPositions, error) {
	characterPositionMap := map[int64]*common.CharacterPositions{}
	for uid, tmpBuild := range buildings {
		if uid != BuildUid {
			tmpCharacters, ok := tmpBuild.FetchCharacters()
			if !ok {
				continue
			}
			// 其他建筑的character移除
			newTmpCharacters, removed := tmpCharacters.RemoveIfExist(characters)
			if removed {
				err := g.CheckConsumeProduceInProduce(tmpBuild)
				if err != nil {
					return nil, errors.WrapTrace(err)
				}
				characterPositionMap[uid] = newTmpCharacters
			}
		} else {
			err := g.CanCharacterDispatch(targetBuilding.BuildId, targetBuilding.FetchLevel(), characters)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			err = g.CheckConsumeProduceInProduce(tmpBuild)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			characterPositionMap[uid] = characters

		}
	}
	return characterPositionMap, nil
}

// RandomRewardHours 在配置范围内随机n个整点时间
func (g *GraveyardEntry) RandomRewardHours() []int32 {

	start := g.GraveyardPlotRewardHourStart()
	end := g.graveyardPlotRewardHourEnd
	var weights []int32

	for i := 0; i < int(end-start); i++ {
		weights = append(weights, 1)
	}
	perm := rand.UniquePerm(int(g.GraveyardPlotRewardMaxNum()), weights)
	var randomHours []int32
	for _, i := range perm {
		randomHours = append(randomHours, int32(i)+start)
	}

	less := func(i, j int) bool {
		return randomHours[i] < randomHours[j]
	}
	sort.Slice(randomHours, less)

	return randomHours
}

func (g *GraveyardEntry) GetPlotRewardNum(hours []int32, nowHour int32) int32 {
	var rewardNum int32
	for _, hour := range hours {
		if hour > nowHour {
			break
		}
		rewardNum++
	}
	return rewardNum
}

// GraveyardGetBuildsProductByItem 使用道具直接获得产出
func (m *Manager) GraveyardGetBuildsProductByItem(builds map[int64]*common.UserGraveyardBuild, item *GraveyardGetProduct) (*common.Rewards, error) {
	once := common.NewRewards()
	for _, build := range builds {
		produceTime, err := m.GraveyardEntry.GetProduceBaseTime(build)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		rewards, err := m.GraveyardGetBuildProductByItem(build, produceTime, item.Sec)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		once.AddRewards(rewards)
	}
	// 如果小于low或者大于up 则等比缩放奖励
	onceWithLowUp := m.GraveyardEntry.RewardsWithLowUp(once, item.Low, item.Up)

	// 避免缩成0
	ret := common.NewRewards()
	for _, reward := range onceWithLowUp.Value() {
		if reward.Num != 0 {
			ret.AddReward(&reward)
		}
	}

	return ret, nil
}

// GraveyardGetBuildProductByItem 使用道具直接获得“单个建筑”产出
func (m *Manager) GraveyardGetBuildProductByItem(build *common.UserGraveyardBuild, produceTime int32, sec int32) (*common.Rewards, error) {
	dropNum, err := m.GraveyardEntry.calDropNumNow(build, produceTime, sec)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	rewardFloatNum := map[int32]float64{}
	for dropId, num := range *dropNum {
		var j float64 = 0
		for i := 0; i < int(math.Ceil(num)); i++ {
			// 计算DropNumFloat 带的小数
			if j+1 > num {
				j = num
			} else {
				j++
			}
			rewards, err := m.Drop.DropRewards(dropId)
			if err != nil {
				glog.Errorf("CalRewardsBySec DropRewards, err:%+v", errors.WrapTrace(err))
				continue
			}
			// 工人会增加产量
			for _, reward := range rewards.MergeValue() {
				rewardFloatNum[reward.ID] = rewardFloatNum[reward.ID] + float64(reward.Num)*(j-float64(i))

			}
			glog.Debugf("GraveyardGetBuildProductByItem, rewardFloatNum:%+v", rewardFloatNum)

		}
	}

	return ToCeilRewards(rewardFloatNum), nil
}
