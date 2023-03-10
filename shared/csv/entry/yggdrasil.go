package entry

import (
	"fmt"
	"math"
	"shared/common"
	"shared/csv/static"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/number"
	"shared/utility/rand"
	"shared/utility/servertime"
	"shared/utility/transfer"
	"sync"
	"time"
)

const CfgNameYggdrasilAreaPos = "cfg_yggdrasil_area_pos"
const CfgNameYggdrasilArea = "cfg_yggdrasil_area"
const CfgNameYggdrasilObject = "cfg_yggdrasil_object"
const CfgNameYggdrasilObjectState = "cfg_yggdrasil_object_state"
const CfgNameYggdrasilAreaCity = "cfg_yggdrasil_area_city"
const CfgNameYggdrasilTask = "cfg_yggdrasil_task"
const CfgNameYggdrasilSubTask = "cfg_yggdrasil_sub_task"
const CfgNameYggdrasilBuilding = "cfg_yggdrasil_building"
const CfgNameYggdrasilCoBuilding = "cfg_yggdrasil_co_building"
const CfgNameYggdrasilMark = "cfg_yggdrasil_mark"
const CfgNameYggdrasilDispatch = "cfg_yggdrasil_dispatch"
const CfgNameYggdrasilDailyMonster = "cfg_yggdrasil_daily_monster"
const CfgNameYggdrasilSubTaskEnv = "cfg_yggdrasil_sub_task_env"

var PosTypePriority map[int32]int32

const (
	PosTypeWalkable    = 0 // 可行走区域
	PosTypeUnWalkable  = 1 // 不可行走区域
	PosTypeUnBuildable = 2 // 可行走不可建造

	DefaultHeight = 1
)

func init() {
	/**
	PosTypeWalkable    = 0 // 可行走区域
	PosTypeUnWalkable  = 1 // 不可行走区域
	PosTypeUnBuildable = 2 // 可行走不可建造

	重叠的情况下1优先级最高 其次是2 然后是0

	*/
	PosTypePriority = map[int32]int32{}
	PosTypePriority[PosTypeUnWalkable] = 3
	PosTypePriority[PosTypeUnBuildable] = 2
	PosTypePriority[PosTypeWalkable] = 1
}

type YggdrasilEntry struct {
	sync.RWMutex

	areas          map[int32]YggdrasilAreaConfig
	posInitHeight  map[coordinate.Position]int32
	posInitType    map[coordinate.Position]int32
	objects        map[int32]YggdrasilObjectConfig
	cities         map[int32]YggdrasilAreaCity
	subTasks       map[int32]CfgYggdrasilSubTask
	tasks          map[int32]CfgYggdrasilTask
	buildings      map[int32]YggdrasilBuilding
	coBuildings    map[int32]YggdrasilCoBuilding
	objectInitArea map[int32][]int32
	posAreaId      map[coordinate.Position]int32 // 坐标对应areaId的缓存
	posLightCityId map[coordinate.Position]int32 // 坐标对应cityId光照半径的缓存
	posApCost      map[coordinate.Position]int32 // 坐标对应消耗
	marks          map[int32]CfgYggdrasilMark
	dispatches     map[int32]CfgYggdrasilDispatch
	dailyMonsters  map[int32]CfgYggdrasilDailyMonster
	envs           map[int32]*Env

	yggBlockLength             int32
	yggBlockWidth              int32
	yggBagAllCount             int32
	yggApRestoreTime           time.Duration
	yggOtherMaxTaskCount       int32
	yggBuildHelpCount          int32
	yggInitCity                int32
	yggCharacterRestSec        int64
	yggLightRadius             int32
	yggMarkBarrageLimit        int32
	yggCityBuildLimitRadius    int32
	yggMatchGrabAlgorithmParam float64
	yggInitPos                 coordinate.Position
	yggDailyTravelTime         int32
	yggMarkTotalCount          int32
	yggMailMaxCount            int32
	yggEditTeamLevel           []int32
	yggBagLockLevel            []int32
	yggEditTeamMaxCount        int32
	mailReceiveMaxCount        int
}

func NewYggdrasilEntry() *YggdrasilEntry {
	return &YggdrasilEntry{
		areas:          map[int32]YggdrasilAreaConfig{},
		posInitHeight:  map[coordinate.Position]int32{},
		posInitType:    map[coordinate.Position]int32{},
		objects:        map[int32]YggdrasilObjectConfig{},
		cities:         map[int32]YggdrasilAreaCity{},
		subTasks:       map[int32]CfgYggdrasilSubTask{},
		tasks:          map[int32]CfgYggdrasilTask{},
		buildings:      map[int32]YggdrasilBuilding{},
		coBuildings:    map[int32]YggdrasilCoBuilding{},
		objectInitArea: map[int32][]int32{},
		posAreaId:      map[coordinate.Position]int32{},
		posLightCityId: map[coordinate.Position]int32{},
		posApCost:      map[coordinate.Position]int32{},
		marks:          map[int32]CfgYggdrasilMark{},
		dispatches:     map[int32]CfgYggdrasilDispatch{},
		dailyMonsters:  map[int32]CfgYggdrasilDailyMonster{},
		envs:           map[int32]*Env{},
	}
}

type YggdrasilAreaConfig struct {
	Id                      int32
	WorldId                 int32
	AreaCityIds             []int32
	Area                    *common.Area `src:"AreaPosId" rule:"area"`
	UnlockCondition         string
	ItemMaxCount            int32
	MessageMaxMatchCount    int32
	MessageMaxCreateCount   int32
	MoveCost                int32
	BuildMaxBuildCount      []int32               `src:"AMaxBuildCount,BMaxBuildCount,CMaxBuildCount,DMaxBuildCount,EMaxBuildCount,FMaxBuildCount" rule:"mergeInt32"`
	BuildMaxMatchCount      []int32               `src:"AMaxMatchCount,BMaxMatchCount,CMaxMatchCount,DMaxMatchCount,EMaxMatchCount,FMaxMatchCount" rule:"mergeInt32"`
	InitObjs                []YggdrasilPosInitObj `ignore:"true"`
	ExploredProgressPercent []int32
	ExploredProgressDrop    []int32
	PosCount                int32 `ignore:"true"` // area的所有点总数
	DispatchUnlock          int32
	DailyStar               map[int32]int32 `rule:"int32Map"`
	GuildStar               []int32
	GuildNum                []int32
	PrestigeItemID          int32
	SafePos                 []coordinate.Position `rule:"stringsToPositions"`
}

func (y YggdrasilAreaConfig) GetMaxBuildBuildCount(buildType int32) int32 {
	switch buildType {
	case static.YggBuildingTypeApSpring:
		if len(y.BuildMaxMatchCount) >= 1 {
			return y.BuildMaxBuildCount[0]
		} else {
			return 0
		}
	case static.YggBuildingTypeHpSpring:
		if len(y.BuildMaxMatchCount) >= 2 {
			return y.BuildMaxBuildCount[1]
		} else {
			return 0
		}
	case static.YggBuildingTypeMagicPeeping:
		if len(y.BuildMaxMatchCount) >= 3 {
			return y.BuildMaxBuildCount[2]
		} else {
			return 0
		}
	case static.YggBuildingTypeMagicTransport:
		if len(y.BuildMaxMatchCount) >= 4 {
			return y.BuildMaxBuildCount[3]
		} else {
			return 0
		}
	case static.YggBuildingTypeMagicBuff:
		if len(y.BuildMaxMatchCount) >= 5 {
			return y.BuildMaxBuildCount[4]
		} else {
			return 0
		}
	case static.YggBuildingTypeStepladder:
		if len(y.BuildMaxMatchCount) >= 6 {
			return y.BuildMaxBuildCount[5]
		} else {
			return 0
		}
	}
	return 0
}

func (y YggdrasilAreaConfig) GetBuildMatchCount(buildType int32) int32 {
	switch buildType {
	case static.YggBuildingTypeApSpring:
		if len(y.BuildMaxMatchCount) >= 1 {
			return y.BuildMaxMatchCount[0]
		} else {
			return 0
		}
	case static.YggBuildingTypeHpSpring:
		if len(y.BuildMaxMatchCount) >= 2 {
			return y.BuildMaxMatchCount[1]
		} else {
			return 0
		}
	case static.YggBuildingTypeMagicPeeping:
		if len(y.BuildMaxMatchCount) >= 3 {
			return y.BuildMaxMatchCount[2]
		} else {
			return 0
		}
	case static.YggBuildingTypeMagicTransport:
		if len(y.BuildMaxMatchCount) >= 4 {
			return y.BuildMaxMatchCount[3]
		} else {
			return 0
		}
	case static.YggBuildingTypeMagicBuff:
		if len(y.BuildMaxMatchCount) >= 5 {
			return y.BuildMaxMatchCount[4]
		} else {
			return 0
		}
	case static.YggBuildingTypeStepladder:
		if len(y.BuildMaxMatchCount) >= 6 {
			return y.BuildMaxMatchCount[5]
		} else {
			return 0
		}
	}
	return 0
}

type YggdrasilBuilding struct {
	Id               int32
	BuildingType     int32
	UnlockConditions *common.Conditions        `rule:"conditions"`
	SingleCost       map[int32]*common.Rewards `rule:"toYggdrasilAreaBuildCosts"`
	UsingCost        *common.Rewards           `rule:"rewards"`
	UsingParam       []int32
	UsingTimes       int32
	BuildingR        int32
	MatchParam       int32
	Intimacy         int32
}

type YggdrasilCoBuilding struct {
	Id           int32
	BuildingType int32
	PosX         int32
	PosY         int32
	Pos          *coordinate.Position `rule:"position" src:"PosX,PosY"`
	BuildCost    *common.Rewards      `rule:"rewards"`
	UsingCost    *common.Rewards      `rule:"rewards"`
	Intimacy     int32
}

type YggdrasilPosInitObj struct {
	Position *coordinate.Position `rule:"position" src:"AreaPosX,AreaPosY"`
	ObjectId int32
}

type YggdrasilObjectConfig struct {
	Id           int32
	DefaultState int32                          `src:"ObjectState"`
	Sates        map[int32]YggdrasilObjectState `ignore:"true"`
}

type YggdrasilObjectState struct {
	Id          int32
	ObjectType  int32
	PosType     int32
	NextState   int32
	SubTaskID   int32 `src:"ServerSubTaskId"` // 携带子任务ID,才能与交互物交互
	ObjectParam int32 `rule:"stringToInt32"`
}

type YggdrasilAreaCity struct {
	Id            int32
	CostAp        int32
	WorldId       int32
	CityPos       []coordinate.Position `rule:"stringsToPositions"`
	CityCenterPos *coordinate.Position  `rule:"position" src:"CityCenterPosX,CityCenterPosY"`
	CityExitPos   *coordinate.Position  `rule:"position" src:"CityExitPosX,CityExitPosY"`
	CityRadius    int32
	CityBanR      int32
}

type CfgYggdrasilTask struct {
	Id                int32
	FinishTaskCity    int32
	TaskType          int32
	UnlockCondition   *common.Conditions `rule:"conditions"`
	TaskGroup         int32
	NextSubTaskId     int32
	DropId            int32
	EnableMatchAreaId int32
	AreaId            int32 `ignore:"true"`
	AbandonTask       bool
}

type CfgYggdrasilSubTask struct {
	Id             int32
	TaskId         int32
	Envs           *Envs                    `ignore:"true"`
	NextSubTaskIds *number.NonRepeatableArr `rule:"toNonRepeatableArr"`
	DropId         int32
	YggDropId      int32
}

type Env struct {
	Id                int32
	SubTaskId         int32
	SubTaskType       int32
	SubTaskTargets    *common.YggdrasilSubTaskTargets `rule:"toYggdrasilSubTaskTargets" src:"Id,SubTaskType,Target"`
	ChangeObjectState []int32
	CreateObjects     *common.EnvObjs     `rule:"stringsToEnvObjs" src:"Id,CreateGroupObject,CreateDeleteAtGroupObj,DeleteAt"`
	ClearPosGroup     *common.Area        ` rule:"area"`
	DeleteTaskItem    *common.TaskItems   `rule:"stringsToTaskItems"`
	AddTaskItem       *common.TaskItems   `rule:"stringsToTaskItems"`
	Terrains          *common.EnvTerrains `rule:"stringsToEnvTerrains" src:"Id,ChangeTerrainState,TerrainStateDeleteAt"`
}

type Envs map[int32]*Env

func NewEnvs() *Envs {
	return (*Envs)(&map[int32]*Env{})
}

func (e *Envs) Env(envId int32) *Env {
	return (*e)[envId]

}
func (e *Envs) RandomEnv() (*Env, error) {
	if len(*e) == 0 {
		return nil, errors.WrapTrace(errors.New("RandomEnv"))
	}
	rangeInt := rand.RangeInt(0, len(*e)-1)
	for _, v := range *e {
		if rangeInt == 0 {
			return v, nil
		}
		rangeInt--
	}

	return nil, errors.WrapTrace(errors.New("RandomEnv"))
}

type CfgYggdrasilMark struct {
	Id int32
}

type CfgYggdrasilDispatch struct {
	Id                  int32
	AreaId              int32
	Type                int32 `src:"Type"`
	TimeCost            int64 `src:"TimeCost" rule:"intMinuteToSecond"`
	TaskStarsCnt        int32 `src:"TaskStarsCount"`
	TeamSize            int32
	GuildCharacNum      int32              `src:"GuildCharacterNum"`
	GuildCharacId       int32              `src:"GuildCharacterId"`
	NecessaryConditions *common.Conditions `rule:"conditions"`
	ExtraConditions     *common.Conditions `rule:"conditions"`
	BaseRewards         int32              `src:"BaseDropId"`
	ExtraRewards        []int32            `src:"ExtraDropId"`
	CloseTime           int64              `src:"CloseTime" rule:"intMinuteToSecond"`
}

type CfgYggdrasilDailyMonster struct {
	Id          int32
	ObjectId    int32
	Pos         *coordinate.Position `rule:"position" src:"PosX,PosY"`
	Radius      int32
	IntervalDay int32
}

func (y *YggdrasilEntry) Check(config *Config) error {
	return nil
}

func (y *YggdrasilEntry) Reload(config *Config) error {
	y.Lock()
	defer y.Unlock()
	areas := map[int32]YggdrasilAreaConfig{}
	var pAreas []*YggdrasilAreaConfig
	posInitObj := map[int32][]YggdrasilPosInitObj{}
	posInitHeight := map[coordinate.Position]int32{}
	posInitType := map[coordinate.Position]int32{}
	objects := map[int32]YggdrasilObjectConfig{}
	cities := map[int32]YggdrasilAreaCity{}
	subTasks := map[int32]CfgYggdrasilSubTask{}
	tasks := map[int32]CfgYggdrasilTask{}
	builds := map[int32]YggdrasilBuilding{}
	coBuilds := map[int32]YggdrasilCoBuilding{}
	marks := map[int32]CfgYggdrasilMark{}
	objectInitArea := map[int32][]int32{}
	dispatches := map[int32]CfgYggdrasilDispatch{}
	dailyMonsters := map[int32]CfgYggdrasilDailyMonster{}
	posAreaId := map[coordinate.Position]int32{}
	posLightCityId := map[coordinate.Position]int32{}
	posApCost := map[coordinate.Position]int32{}
	envs := map[int32]*Env{}

	for _, v := range config.CfgYggdrasilDispatchConfig.GetAllData() {
		dispatchCfg := &CfgYggdrasilDispatch{}
		err := transfer.Transfer(v, dispatchCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}
		if len(*dispatchCfg.ExtraConditions) != len(dispatchCfg.ExtraRewards) {
			return errors.Swrapf(common.ErrYggDispatchCSVExtraConditionNotMatchRewards, dispatchCfg.Id)
		}
		if dispatchCfg.Type != static.YggdrasilDispatchTypeDaily && dispatchCfg.Type != static.YggdrasilDispatchTypeGuild {
			return errors.Swrapf(common.ErrYggDispatchCSVWrongDispatchType, dispatchCfg.Id)
		}
		if dispatchCfg.Type == static.YggdrasilDispatchTypeDaily && (dispatchCfg.CloseTime > 0) {
			return errors.Swrapf(common.ErrYggDispatchCSVDispatchTypeNotMatchCloseTime, dispatchCfg.Id)
		}
		if dispatchCfg.Type == static.YggdrasilDispatchTypeGuild && dispatchCfg.CloseTime == 0 {
			return errors.Swrapf(common.ErrYggDispatchCSVDispatchTypeNotMatchCloseTime, dispatchCfg.Id)
		}
		if dispatchCfg.Type == static.YggdrasilDispatchTypeGuild && (dispatchCfg.GuildCharacId <= 0 || dispatchCfg.GuildCharacNum > 1) {
			return errors.Swrapf(common.ErrYggDispatchCSVWrongGuildCharacter, dispatchCfg.Id, dispatchCfg.GuildCharacNum, dispatchCfg.GuildCharacId)
		}
		dispatches[v.Id] = *dispatchCfg
	}

	for _, v := range config.CfgYggdrasilAreaConfig.GetAllData() {
		yggAreaConfig := &YggdrasilAreaConfig{}
		err := transfer.Transfer(v, yggAreaConfig)
		if err != nil {
			return errors.WrapTrace(err)
		}
		if len(yggAreaConfig.GuildNum) != 2 {
			return errors.Swrapf(common.ErrYggDispatchCSVWrongLengthOfGuildNum, yggAreaConfig.Id)
		}
		yggAreaConfig.PosCount = yggAreaConfig.Area.Count()
		pAreas = append(pAreas, yggAreaConfig)
	}

	findArea := func(p coordinate.Position) (int32, bool) {
		var findAreaId int32
		for _, a := range pAreas {
			if a.Area.Contains(p) {
				findAreaId = a.Id
				break
			}
		}
		if findAreaId == 0 {
			return 0, false
		}
		return findAreaId, true
	}

	for _, v := range config.CfgYggdrasilAreaPosConfig.GetAllData() {
		//验证初始点位配置正确性
		p := *coordinate.NewPosition(v.AreaPosX, v.AreaPosY)
		findAreaId, ok := findArea(p)

		if !ok {
			//glog.Error(errors.Swrapf(common.ErrYggdrasilNoAreaForPos, p))

			//return errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilNoAreaForPos, p))
		}

		if findAreaId != v.AreaGroupId {
			//glog.Error(errors.Swrapf(common.ErrYggdrasilInitPosAreaError))

			//return errors.WrapTrace(common.ErrYggdrasilInitPosAreaError)
		}

		if v.ObjectId > 0 {
			yggdrasilPosInitObj := &YggdrasilPosInitObj{}
			err := transfer.Transfer(v, yggdrasilPosInitObj)
			if err != nil {
				return errors.WrapTrace(err)
			}

			initObjs := posInitObj[v.AreaGroupId]
			initObjs = append(initObjs, *yggdrasilPosInitObj)
			posInitObj[v.AreaGroupId] = initObjs

			objectInitArea[v.ObjectId] = append(objectInitArea[v.ObjectId], v.AreaGroupId)
		}

		position := *coordinate.NewPosition(v.AreaPosX, v.AreaPosY)
		posInitHeight[position] = v.AreaPosHeight
		posInitType[position] = v.PosType
	}

	for _, pArea := range pAreas {
		pArea.InitObjs = posInitObj[pArea.Id]
		areas[pArea.Id] = *pArea
	}
	// 各区域不能重合
	for _, pArea := range pAreas {
		for otherAreaId, otherArea := range areas {
			if otherAreaId == pArea.Id {
				continue
			}
			if pArea.Area.CoincidenceArea(otherArea.Area) {
				return errors.WrapTrace(common.ErrYggdrasilAreaCoincidence)
			}

		}

	}

	for _, v := range config.CfgYggdrasilObjectConfig.GetAllData() {
		yggdrasilObjectConfig := &YggdrasilObjectConfig{}
		err := transfer.Transfer(v, yggdrasilObjectConfig)
		if err != nil {
			return errors.WrapTrace(err)
		}
		yggdrasilObjectConfig.Sates = map[int32]YggdrasilObjectState{}
		objects[v.Id] = *yggdrasilObjectConfig
	}
	totalStates := map[int32]YggdrasilObjectState{}
	for _, v := range config.CfgYggdrasilObjectStateConfig.GetAllData() {
		yggdrasilObjectState := &YggdrasilObjectState{}
		err := transfer.Transfer(v, yggdrasilObjectState)
		if err != nil {
			return errors.WrapTrace(err)
		}
		if yggdrasilObjectState.ObjectType == static.YggObjectTypeChest {
			if yggdrasilObjectState.ObjectParam == 0 {
				return errors.New("CfgYggdrasilObjectState Type Chest param cannot =0, CfgYggdrasilObjectState Id = %d", yggdrasilObjectState.Id)
			}
			_, ok := config.CfgDropDataConfig.Find(yggdrasilObjectState.ObjectParam)
			if !ok {
				return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgDropDataConfig, yggdrasilObjectState.ObjectParam))

			}
		}
		totalStates[v.Id] = *yggdrasilObjectState
	}
	for _, object := range objects {
		state, ok := totalStates[object.DefaultState]
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObjectState, object.DefaultState)
		}
		for {
			object.Sates[state.Id] = state
			if state.NextState == 0 {
				break
			}
			state, ok = totalStates[state.NextState]
			if !ok {
				return errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObjectState, state.NextState)
			}

		}
	}

	for _, v := range config.CfgYggdrasilAreaCityConfig.GetAllData() {

		city := &YggdrasilAreaCity{}
		err := transfer.Transfer(v, city)
		if err != nil {
			return errors.WrapTrace(err)
		}
		cities[v.Id] = *city
	}

	for _, v := range config.CfgYggdrasilTaskConfig.GetAllData() {
		task := &CfgYggdrasilTask{}
		err := transfer.Transfer(v, task)
		if err != nil {
			return errors.WrapTrace(err)
		}
		city, ok := config.CfgYggdrasilAreaCityConfig.Find(task.FinishTaskCity)
		if !ok {
			return errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilAreaCity, task.FinishTaskCity)

		}
		p := *coordinate.NewPosition(city.CityCenterPosX, city.CityCenterPosY)
		findAreaId, ok := findArea(p)
		if !ok {
			return errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilNoAreaForPos, p))

		}
		task.AreaId = findAreaId
		tasks[v.Id] = *task
	}

	envsCfg := map[int32][]*Env{}
	for _, cfg := range config.CfgYggdrasilSubTaskEnvConfig.GetAllData() {
		env := &Env{}
		err := transfer.Transfer(cfg, env)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// check
		for _, obj := range *env.CreateObjects {
			_, ok := objects[obj.ObjectId]
			if !ok {
				return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObject, obj.ObjectId))
			}
		}
		for _, objId := range env.ChangeObjectState {
			_, ok := objects[objId]
			if !ok {
				return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObject, objId))
			}
		}
		if env.SubTaskType == static.YggdrasilSubTaskTypeTypeOwn {
			for _, target := range *env.SubTaskTargets {
				itemId := target.FilterConditions[0]
				itemData, ok := config.CfgItemDataConfig.Find(itemId)
				if !ok {
					return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgItemData, itemId))
				}
				if itemData.ItemType != static.ItemTypeYggItemSpecial {
					return errors.WrapTrace(errors.New("Yggdrasil SubTaskTypeTypeOwn item type err,envId:%d,itemId:%d", env.Id, itemId))
				}
			}
		}

		envsCfg[env.SubTaskId] = append(envsCfg[env.SubTaskId], env)
		envs[env.Id] = env
	}
	for _, v := range config.CfgYggdrasilSubTaskConfig.GetAllData() {
		subTask := &CfgYggdrasilSubTask{}
		err := transfer.Transfer(v, subTask)
		if err != nil {
			return errors.WrapTrace(err)
		}
		num := len(envsCfg[v.Id])
		if num == 0 {
			return errors.WrapTrace(errors.New(fmt.Sprintf("env len = 0,subtaskId:%d", v.Id)))
		}
		subTask.Envs = NewEnvs()
		for _, env := range envsCfg[v.Id] {

			(*subTask.Envs)[env.Id] = env
		}
		subTasks[v.Id] = *subTask
	}
	for _, buildCsv := range config.CfgYggdrasilBuildingConfig.GetAllData() {
		buildCfg := &YggdrasilBuilding{}

		err := transfer.Transfer(buildCsv, buildCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		builds[buildCsv.Id] = *buildCfg
	}

	for _, coBuildCsv := range config.CfgYggdrasilCobuildingConfig.GetAllData() {
		coBuildCfg := &YggdrasilCoBuilding{}

		err := transfer.Transfer(coBuildCsv, coBuildCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		coBuilds[coBuildCsv.Id] = *coBuildCfg
	}

	for _, v := range config.CfgYggdrasilMarkConfig.GetAllData() {
		mark := &CfgYggdrasilMark{}
		err := transfer.Transfer(v, mark)
		if err != nil {
			return errors.WrapTrace(err)
		}
		marks[v.Id] = *mark
	}

	for _, v := range config.CfgYggdrasilDailyMonsterConfig.GetAllData() {
		cfg := &CfgYggdrasilDailyMonster{}
		err := transfer.Transfer(v, cfg)
		if err != nil {
			return errors.WrapTrace(err)
		}
		_, ok := objects[cfg.ObjectId]
		if !ok {
			return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObject, cfg.ObjectId))
		}
		dailyMonsters[v.Id] = *cfg
	}

	// 检查区域初始obj是否存在
	for _, area := range areas {
		for _, obj := range area.InitObjs {
			_, ok := objects[obj.ObjectId]
			if !ok {
				return errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObject, obj.ObjectId))
			}
		}
	}
	// todo:检查env的点有没有覆盖联合建筑

	y.areas = areas
	y.posInitHeight = posInitHeight
	y.posInitType = posInitType
	y.objects = objects
	y.cities = cities
	y.subTasks = subTasks
	y.tasks = tasks
	y.buildings = builds
	y.coBuildings = coBuilds
	y.objectInitArea = objectInitArea
	y.posAreaId = posAreaId
	y.posLightCityId = posLightCityId
	y.posApCost = posApCost
	y.marks = marks
	y.dispatches = dispatches
	y.dailyMonsters = dailyMonsters
	y.envs = envs

	y.yggBlockLength = config.GlobalEntry.YggBlockLengthAndWidth[0]
	y.yggBlockWidth = config.GlobalEntry.YggBlockLengthAndWidth[1]
	y.yggBagAllCount = config.GlobalEntry.YggBagAllCount
	y.yggApRestoreTime = time.Duration(config.GlobalEntry.YggApRestoreTime) * time.Second
	y.yggOtherMaxTaskCount = config.GlobalEntry.YggOtherMaxTaskCount
	y.yggBuildHelpCount = config.GlobalEntry.YggBuildHelpCount
	y.yggInitCity = config.GlobalEntry.YggInitCity
	y.yggCharacterRestSec = int64(config.GlobalEntry.YggCharacterRestHour) * servertime.SecondPerHour
	y.yggLightRadius = config.GlobalEntry.YggLightRadius
	y.yggMarkBarrageLimit = config.GlobalEntry.YggMarkBarrageLimit
	y.yggCityBuildLimitRadius = config.GlobalEntry.YggCityBuildLimitRadius
	y.yggMatchGrabAlgorithmParam = float64(config.GlobalEntry.YggMatchGrabAlgorithmParam) / 100
	y.yggInitPos = config.GlobalEntry.YggInitPos
	y.yggDailyTravelTime = config.GlobalEntry.YggDailyTravelTime
	y.yggMarkTotalCount = config.GlobalEntry.YggMarkTotalCount
	y.yggMailMaxCount = config.GlobalEntry.YggMailMaxCount
	y.yggEditTeamLevel = config.GlobalEntry.YggEditTeamLevel
	y.yggBagLockLevel = config.GlobalEntry.YggBagLockLevel
	y.yggEditTeamMaxCount = config.GlobalEntry.YggEditTeamMaxCount
	y.mailReceiveMaxCount = int(config.GlobalEntry.MailReceiveMaxCount)
	return nil
}

func (y *YggdrasilEntry) GetAllAreaConfig() map[int32]YggdrasilAreaConfig {
	y.RLock()
	defer y.RUnlock()
	return y.areas
}

func (y *YggdrasilEntry) GetArea(areaId int32) (*YggdrasilAreaConfig, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.areas[areaId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilArea, areaId)
	}
	return &config, nil
}

func (y *YggdrasilEntry) GetPosAreaId(position coordinate.Position) (int32, error) {
	y.RLock()
	defer y.RUnlock()
	return y.getPosAreaId(position)
}

func (y *YggdrasilEntry) getPosAreaId(position coordinate.Position) (int32, error) {
	areaId, ok := y.posAreaId[position]
	if ok {
		return areaId, nil
	}

	for _, config := range y.areas {
		if config.Area.Contains(position) {
			y.posAreaId[position] = config.Id
			return config.Id, nil
		}
	}
	return 0, errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilNoAreaForPos, position))
}

func (y *YggdrasilEntry) GetPosInitHeightType() (map[coordinate.Position]int32, map[coordinate.Position]int32) {
	y.RLock()
	defer y.RUnlock()
	return y.posInitHeight, y.posInitType
}

func (y *YggdrasilEntry) GetYggBlockLengthAndWidth() (int32, int32) {
	y.RLock()
	defer y.RUnlock()
	return y.yggBlockLength, y.yggBlockWidth
}

func (y *YggdrasilEntry) GetYggBagAllCount(lv int32) int32 {
	y.RLock()
	defer y.RUnlock()
	var delta int32
	for i, v := range y.yggBagLockLevel {
		if lv < v {
			delta = int32(len(y.yggBagLockLevel) - i)
			break
		}
	}

	return y.yggBagAllCount - delta
}

func (y *YggdrasilEntry) GetYggEditTeamCount(lv int32) int32 {
	y.RLock()
	defer y.RUnlock()
	var delta int32
	for i, v := range y.yggEditTeamLevel {
		if lv < v {
			delta = int32(len(y.yggEditTeamLevel) - i)
			break
		}
	}

	return y.yggEditTeamMaxCount - delta
}

func (y *YggdrasilEntry) GetMailReceiveMaxCount() int {
	y.RLock()
	defer y.RUnlock()
	return y.mailReceiveMaxCount
}
func (y *YggdrasilEntry) GetYggApRestoreTime() time.Duration {
	y.RLock()
	defer y.RUnlock()
	return y.yggApRestoreTime
}

func (y *YggdrasilEntry) GetYggOtherMaxTaskCount() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggOtherMaxTaskCount
}

func (y *YggdrasilEntry) GetYggBuildHelpCount() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggBuildHelpCount
}

func (y *YggdrasilEntry) GetYggInitCity() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggInitCity
}

func (y *YggdrasilEntry) GetYggCharacterRestSec() int64 {
	y.RLock()
	defer y.RUnlock()
	return y.yggCharacterRestSec
}

func (y *YggdrasilEntry) GetYggLightRadius() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggLightRadius
}
func (y *YggdrasilEntry) GetYggMarkBarrageLimit() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggMarkBarrageLimit
}
func (y *YggdrasilEntry) GetYggCityBuildLimitRadius() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggCityBuildLimitRadius
}
func (y *YggdrasilEntry) GetYggMatchGrabAlgorithmParam() float64 {
	y.RLock()
	defer y.RUnlock()
	return y.yggMatchGrabAlgorithmParam
}

func (y *YggdrasilEntry) GetYggInitPos() coordinate.Position {
	y.RLock()
	defer y.RUnlock()
	return y.yggInitPos
}

func (y *YggdrasilEntry) GetYggDailyTravelTime() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggDailyTravelTime
}
func (y *YggdrasilEntry) GetYggMarkTotalCount() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggMarkTotalCount
}
func (y *YggdrasilEntry) GetYggMailMaxCount() int32 {
	y.RLock()
	defer y.RUnlock()
	return y.yggMailMaxCount
}

func (y *YggdrasilEntry) GetYggdrasilObjectConfig(objectId int32) (*YggdrasilObjectConfig, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.objects[objectId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObject, objectId)
	}
	return &config, nil

}

func (y *YggdrasilEntry) GetYggdrasilObjectType(objectId, state int32) (int32, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.objects[objectId]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObject, objectId)
	}
	stateConfig, ok := config.Sates[state]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObjectState, state)
	}
	return stateConfig.ObjectType, nil

}

func (y *YggdrasilEntry) GetYggCityById(cityId int32) (*YggdrasilAreaCity, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.cities[cityId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilAreaCity, cityId)
	}
	return &config, nil
}

// IsCityEntrance 是否是城市入口
func (y *YggdrasilEntry) IsCityEntrance(position *coordinate.Position) int32 {
	y.RLock()
	defer y.RUnlock()
	for _, city := range y.cities {
		for _, po := range city.CityPos {
			if po == *position {
				return city.Id
			}
		}
	}
	return 0
}

func (y *YggdrasilEntry) GetSubTaskConfig(subTaskId int32) (*CfgYggdrasilSubTask, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.subTasks[subTaskId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilSubTask, subTaskId)
	}
	return &config, nil
}
func (y *YggdrasilEntry) GetEnv(envId int32) (*Env, error) {
	y.RLock()
	defer y.RUnlock()

	env, ok := y.envs[envId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilSubTaskEnv, envId)
	}
	return env, nil
}
func (y *YggdrasilEntry) GetTaskConfig(taskId int32) (*CfgYggdrasilTask, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.tasks[taskId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilTask, taskId)
	}
	return &config, nil
}

func (y *YggdrasilEntry) GetYggdrasilBuilding(id int32) (*YggdrasilBuilding, error) {
	y.RLock()
	defer y.RUnlock()

	build, ok := y.buildings[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilBuilding, id)
	}

	return &build, nil
}

// 返回object默认状态下的物品类型和下一状态
func (y *YggdrasilEntry) GetYggdrasilObjectTypeAndNextState(objectId int32, curState int32) (int32, int32, error) {
	y.RLock()
	defer y.RUnlock()

	config, ok := y.objects[objectId]
	if !ok {
		return 0, 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObject, objectId)
	}

	state, ok := config.Sates[curState]
	if !ok {
		return 0, 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObjectState, config.Id)
	}

	return state.ObjectType, state.NextState, nil
}

func (y *YggdrasilEntry) GetPosInitHeightAndType(position coordinate.Position) (int32, int32) {
	y.RLock()
	defer y.RUnlock()

	height, ok := y.posInitHeight[position]
	if !ok {
		height = DefaultHeight
	}

	posType, ok := y.posInitType[position]
	if !ok {
		posType = PosTypeWalkable
	}

	return height, posType
}

func (y *YggdrasilEntry) GetObjectState(objectId, state int32) (*YggdrasilObjectState, error) {

	config, err := y.GetYggdrasilObjectConfig(objectId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	objState, ok := config.Sates[state]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilObjectState, state)
	}
	return &objState, nil

}

func (y *YggdrasilEntry) GetYggCoBuilding(buildId int32) (*YggdrasilCoBuilding, error) {
	y.RLock()
	defer y.RUnlock()

	coBuild, ok := y.coBuildings[buildId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilCoBuilding, buildId)
	}

	return &coBuild, nil

}

func (y *YggdrasilEntry) GetAllYggCoBuild(targetAreaId int32) []*YggdrasilCoBuilding {
	y.RLock()
	defer y.RUnlock()

	result := []*YggdrasilCoBuilding{}

	for _, build := range y.coBuildings {
		areaId, ok := y.posAreaId[*build.Pos]
		if !ok {
			for _, config := range y.areas {
				if config.Area.Contains(*build.Pos) {
					y.posAreaId[*build.Pos] = config.Id
					areaId = config.Id
				}
			}
		}
		if areaId == targetAreaId {
			result = append(result, &build)
		}
	}

	return result
}

func (y *YggdrasilEntry) GetObjectInitArea(objectId int32) []int32 {
	y.RLock()
	defer y.RUnlock()
	return y.objectInitArea[objectId]
}

func (y *YggdrasilEntry) GetExploreProcessIndex(areaId, percent int32) (int, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.areas[areaId]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilArea, areaId)
	}
	for i, progress := range config.ExploredProgressPercent {

		if percent >= progress {
			return i, nil
		}
	}
	return 0, errors.WrapTrace(common.ErrParamError)
}

func (y *YggdrasilEntry) GetMark(markId int32) (*CfgYggdrasilMark, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.marks[markId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilMark, markId)
	}
	return &config, nil
}

// CityLightPos 在城市光照半径内，返回城市id
func (y *YggdrasilEntry) CityLightPos(position coordinate.Position) (int32, bool) {
	y.RLock()
	defer y.RUnlock()
	cityId, ok := y.posLightCityId[position]
	if ok {
		return cityId, true
	}
	for _, city := range y.cities {
		distance := coordinate.CubeDistance(*city.CityCenterPos, position)
		if distance <= city.CityRadius {
			y.posLightCityId[position] = city.Id
			return city.Id, true
		}
	}
	return -1, false
}

func (y *YggdrasilEntry) GetMostCost(position coordinate.Position) (int32, error) {
	y.RLock()
	defer y.RUnlock()

	cost, ok := y.posApCost[position]
	if ok {
		return cost, nil
	}

	areaId, err := y.getPosAreaId(position)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	config, ok := y.areas[areaId]
	if !ok {
		return 0, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilArea, areaId)
	}
	for _, city := range y.cities {
		distance := coordinate.CubeDistance(*city.CityCenterPos, position)
		if distance <= city.CityBanR {
			y.posApCost[position] = 0
			return 0, nil
		}
	}
	y.posApCost[position] = config.MoveCost
	return config.MoveCost, nil
}

func (y *YggdrasilEntry) GetAllDailyMonster() map[int32]CfgYggdrasilDailyMonster {
	y.RLock()
	defer y.RUnlock()
	return y.dailyMonsters
}

func (y *YggdrasilEntry) GetDailyMonster(id int32) (*CfgYggdrasilDailyMonster, error) {
	y.RLock()
	defer y.RUnlock()
	config, ok := y.dailyMonsters[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilDailyMonster, id)
	}
	return &config, nil
}

func (y *YggdrasilEntry) GetDispatchTask(taskId int32) (*CfgYggdrasilDispatch, error) {
	y.RLock()
	defer y.RUnlock()

	task, ok := y.dispatches[taskId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilDispatch, taskId)
	}

	return &task, nil
}

func (y *YggdrasilEntry) GetSpecificTasks(areaId, dispatchType, star int32) (map[int32]CfgYggdrasilDispatch, error) {
	y.RLock()
	defer y.RUnlock()

	result := map[int32]CfgYggdrasilDispatch{}
	for _, task := range y.dispatches {
		if (task.AreaId == areaId) && (task.Type == dispatchType) && (task.TaskStarsCnt == star) {
			result[task.Id] = task
		}
	}

	return result, nil
}

func (y *YggdrasilEntry) GetAllAreaIDs() ([]int32, error) {
	y.RLock()
	defer y.RUnlock()

	result := []int32{}
	for id, _ := range y.areas {
		result = append(result, id)
	}

	if len(result) == 0 {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilArea, -1)
	}

	return result, nil
}

func (y *YggdrasilEntry) GetAllDispatchTask() map[int32]CfgYggdrasilDispatch {
	y.RLock()
	defer y.RUnlock()

	return y.dispatches
}

func (y *YggdrasilEntry) GetClosestSafePos(position coordinate.Position) (*coordinate.Position, error) {
	y.RLock()
	defer y.RUnlock()

	areaId, err := y.getPosAreaId(position)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	config, ok := y.areas[areaId]
	if !ok {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilArea, areaId)
	}
	var min int32 = math.MaxInt32
	var ret *coordinate.Position
	for _, to := range config.SafePos {
		distance := coordinate.CubeDistance(position, to)
		if distance < min {
			min = distance
			ret = to.Clone()
		}

	}
	if ret == nil {
		return nil, errors.Swrapf(common.ErrNotFoundInCSV, CfgNameYggdrasilArea, areaId)

	}

	return ret, nil
}
