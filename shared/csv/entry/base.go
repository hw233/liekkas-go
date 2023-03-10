package entry

import (
	"log"
	"shared/csv/static"
	"shared/utility/coordinate"
	"shared/utility/number"
	"strconv"
	"strings"

	"shared/common"
	"shared/utility/errors"
	"shared/utility/servertime"
	"shared/utility/transfer"
)

const tenThousand = 10000

func init() {
	transfer.RegisterRule("reward", StringToReward)
	transfer.RegisterRule("rewards", StringsToRewards)
	transfer.RegisterRule("randReward", common.NewRandReward)
	transfer.RegisterRule("stringToInt32", stringToInt32)
	transfer.RegisterRule("conditions", StringsToConditions)
	transfer.RegisterRule("compoundConditions", StringsToCompoundConditions)

	transfer.RegisterRule("intToBool", IntToBool)
	transfer.RegisterRule("buildingLevelLimits", StringsToBuildingLevelLimits)
	transfer.RegisterRule("int32Map", StringsToInt32Map)
	transfer.RegisterRule("area", StringsToArea)
	transfer.RegisterRule("tenThousandth", Int32ToTenThousandth)
	transfer.RegisterRule("mergeInt32", MergeInt32)
	transfer.RegisterRule("floatsToVec2", FloatsToVec2)
	transfer.RegisterRule("int32ToInt64", Int32ToInt64)
	transfer.RegisterRule("mergeIntSlice", MergeInt32Slice)
	transfer.RegisterRule("mergeConditions", MergeConditions)
	transfer.RegisterRule("stringToTimeUnix", StringToTimeUnix)
	transfer.RegisterRule("stringsToPositions", StringsToPositions)
	transfer.RegisterRule("toYggdrasilSubTaskTargets", ToYggdrasilSubTaskTargets)
	transfer.RegisterRule("toNonRepeatableArr", ToNonRepeatableArr)
	transfer.RegisterRule("toYggdrasilAreaBuildCosts", ToYggdrasilAreaBuildCosts)
	transfer.RegisterRule("position", ToPosition)
	transfer.RegisterRule("intMinuteToSecond", IntMinuteToSecond)
	transfer.RegisterRule("dailyTimeToSec", DailyTimeToSec)
	transfer.RegisterRule("stringsToEnvObjs", StringsToEnvObjs)
	transfer.RegisterRule("stringsToTaskItems", StringsToTaskItems)
	transfer.RegisterRule("strsToSkills", StringsToCharacterSkill)
	transfer.RegisterRule("stringsToEnvTerrains", StringsToEnvTerrains)

	err := transfer.CheckRules()
	if err != nil {
		log.Fatalf("check rules error: %v", err)
	}
}

func StringsToRewards(ss []string) (*common.Rewards, error) {
	rewards := common.NewRewards()

	for _, s := range ss {
		reward, err := StringToReward(s)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		if reward.Type != 0 {
			rewards.AddReward(reward)
		}
	}

	return rewards, nil
}

func StringToReward(s string) (*common.Reward, error) {
	if len(s) <= 0 {
		return &common.Reward{
			ID:   0,
			Type: 0,
			Num:  0,
		}, nil
	}

	ss := strings.Split(s, ":")
	if len(ss) != 2 {
		return nil, errors.New("Format error. (string: %s", s)
	}

	id, err := strconv.ParseInt(ss[0], 10, 32)
	if err != nil {
		return nil, errors.New("Format error. (string: %s", s)
	}

	num, err := strconv.ParseInt(ss[1], 10, 32)
	if err != nil {
		return nil, errors.New("Format error. (string: %s", s)
	}

	return common.NewReward(int32(id), int32(num)), nil
}

func stringToInt32(s string) (int32, error) {
	if s == "" {
		return 0, nil
	}
	i, err := strconv.ParseInt(s, 10, 32)
	return int32(i), err
}

func StringsToCompoundConditions(and []string, or []string) (*common.CompoundConditions, error) {
	andConditions, err := StringsToConditions(and)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	var orConditions []*common.Conditions
	for _, str := range or {
		conditions, err := StringsToConditions([]string{str})
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		orConditions = append(orConditions, conditions)
	}
	return common.NewCompoundConditions(andConditions, orConditions...), nil

}
func StringsToConditions(strs []string) (*common.Conditions, error) {
	size := len(strs)
	conditions := common.NewConditions()

	for i := 0; i < size; i++ {
		if len(strs[i]) <= 0 {
			continue
		}

		condition, err := StringToCondition(strs[i])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		conditions.AddCondition(condition)
	}

	return conditions, nil
}

func StringToCondition(str string) (*common.Condition, error) {
	strSlice := strings.Split(str, ":")

	size := len(strSlice)
	if size <= 1 {
		return nil, errors.WrapTrace(common.ErrCSVFormatInvalid)
	}

	conditionType, err := strconv.ParseInt(strSlice[0], 10, 32)
	if err != nil {
		return nil, errors.WrapError(err, common.ErrCSVFormatInvalid)
	}

	params := make([]int32, 0, size-1)
	for i := 1; i < size; i++ {
		param, err := strconv.ParseInt(strSlice[i], 10, 32)
		if err != nil {
			return nil, errors.WrapError(err, common.ErrCSVFormatInvalid)
		}
		params = append(params, int32(param))
	}

	return &common.Condition{
		ConditionType: int32(conditionType),
		Params:        params,
	}, nil
}

func IntToBool(value int32) (bool, error) {
	if value == 0 {
		return false, nil
	}

	if value == 1 {
		return true, nil
	}

	return false, errors.WrapTrace(common.ErrCSVFormatInvalid)
}

func StringsToBuildingLevelLimits(v []string) (*BuildingLevelLimits, error) {
	if len(v) == 0 {
		return &BuildingLevelLimits{}, nil
	}
	var limits []BuildingLevelLimit
	for _, s := range v {
		if s == "" {
			continue
		}
		split := strings.Split(s, ":")
		limit, err := StringToBuildingLevelLimit(split)
		if err != nil {
			return nil, err
		}
		limits = append(limits, *limit)
	}
	return &BuildingLevelLimits{
		Limits: limits,
	}, nil
}

func StringToBuildingLevelLimit(split []string) (*BuildingLevelLimit, error) {
	if len(split) == 2 {
		buildId, err := stringToInt32(split[1])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		num, err := stringToInt32(split[0])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		return &BuildingLevelLimit{
			BuildId: buildId,
			NeedNum: num,
			Lv:      0,
		}, nil
	} else if len(split) == 3 {
		buildId, err := stringToInt32(split[1])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		num, err := stringToInt32(split[0])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		lv, err := stringToInt32(split[2])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		return &BuildingLevelLimit{
			BuildId: buildId,
			NeedNum: num,
			Lv:      lv,
		}, nil
	} else {
		return nil, errors.WrapTrace(common.ErrCSVFormatInvalid)
	}
}

func StringsToInt32Map(v []string) (map[int32]int32, error) {
	int32Map := map[int32]int32{}
	if len(v) == 0 {
		return int32Map, nil
	}
	for _, s := range v {
		if s == "" {
			continue
		}
		split := strings.Split(s, ":")
		if len(split) != 2 {
			return nil, errors.WrapTrace(common.ErrCSVFormatInvalid)
		}
		buildId, err := stringToInt32(split[0])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		num, err := stringToInt32(split[1])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		int32Map[buildId] = num
	}
	return int32Map, nil
}

func StringsToArea(v []string) (*common.Area, error) {
	area := common.NewArea()
	err := common.StringsToArea(v, area)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return area, nil
}

func Int32ToTenThousandth(v int32) float64 {
	return float64(v) / tenThousand
}

func MergeInt32(v ...int32) []int32 {
	var result []int32
	for _, i := range v {
		result = append(result, i)
	}
	return result
}

func FloatsToVec2(floats []float64) (*common.Vec2, error) {
	if len(floats) != 2 {
		return nil, errors.WrapTrace(common.ErrCSVFormatInvalid)
	}

	return common.NewVec2(floats[0], floats[1]), nil
}

func Int32ToInt64(value int32) int64 {
	return int64(value)
}

func MergeConditions(vals ...[]string) ([]*common.Conditions, error) {
	conditions := make([]*common.Conditions, 0, len(vals))
	for _, ss := range vals {
		condition, err := StringsToConditions(ss)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		conditions = append(conditions, condition)
	}

	return conditions, nil
}

func MergeInt32Slice(vals ...[]int32) ([][]int32, error) {
	goods := make([][]int32, 0, len(vals))

	for _, v := range vals {
		goods = append(goods, v)
	}

	return goods, nil
}

func StringToTimeUnix(timeStr string) (int64, error) {
	parse, err := servertime.ParseTime(timeStr)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	return parse, nil
}

func StringsToPositions(strs []string) ([]coordinate.Position, error) {
	var posList []coordinate.Position
	for _, str := range strs {
		ss := strings.Split(str, ":")
		if len(ss) != 2 {
			return nil, errors.New("Format error. (string: %s", ss)
		}
		x, err := stringToInt32(ss[0])
		if err != nil {
			return nil, errors.New("Format error. (string: %s", ss[0])
		}

		y, err := stringToInt32(ss[1])
		if err != nil {
			return nil, errors.New("Format error. (string: %s", ss[1])
		}
		posList = append(posList, *coordinate.NewPosition(x, y))
	}
	return posList, nil
}

func ToYggdrasilSubTaskTargets(envId, subTaskType int32, subTaskData []string) (*common.YggdrasilSubTaskTargets, error) {
	targets := &common.YggdrasilSubTaskTargets{}
	for _, s := range subTaskData {
		split := strings.Split(s, ":")
		var intArr []int32
		for _, sp := range split {
			toInt32, err := stringToInt32(sp)
			if err != nil {
				return nil, errors.New("ToYggdrasilSubTaskData Format error. (string: %s", sp)
			}
			intArr = append(intArr, toInt32)

		}

		switch subTaskType {

		case static.YggdrasilSubTaskTypeTypeVn,
			static.YggdrasilSubTaskTypeTypeMultiTask,
			static.YggdrasilSubTaskTypeTypeCity,
			static.YggdrasilSubTaskTypeTypeMultiSubTask:
			if len(intArr) != 1 {
				return nil, errors.New("ToYggdrasilSubTaskData len error1,envId:%d,subTaskType:%d ,subTaskData:%s ", envId, subTaskType, subTaskData)
			}
		case static.YggdrasilSubTaskTypeTypeMove,
			static.YggdrasilSubTaskTypeTypeObjectStateChange,
			static.YggdrasilSubTaskTypeTypeMonster,
			static.YggdrasilSubTaskTypeTypeHelpBuild,
			static.YggdrasilSubTaskTypeTypeBuild,
			static.YggdrasilSubTaskTypeTypeOwn,
			static.YggdrasilSubTaskTypeTypeDeliverItem,
			static.YggdrasilSubTaskTypeTypeDeliverItemSelectOne,
			static.YggdrasilSubTaskTypeTypeChapter:
			if len(intArr) != 2 {
				return nil, errors.New("ToYggdrasilSubTaskData len error2,envId:%d,subTaskType:%d ,subTaskData:%s ", envId, subTaskType, subTaskData)
			}
		case static.YggdrasilSubTaskTypeTypeConvoy,
			static.YggdrasilSubTaskTypeTypeLeadWay:
			if len(intArr) != 3 {
				return nil, errors.New("ToYggdrasilSubTaskData len error3,envId:%d,subTaskType:%d ,subTaskData:%s ", envId, subTaskType, subTaskData)
			}

		}

		*targets = append(*targets, common.NewYggdrasilSubTaskTarget(subTaskType, intArr))

	}

	return targets, nil
}

func ToNonRepeatableArr(arr []int32) *number.NonRepeatableArr {
	nonRepeatableArr := number.NewNonRepeatableArr()
	for _, i := range arr {
		nonRepeatableArr.Append(i)
	}
	return nonRepeatableArr

}

func ToYggdrasilAreaBuildCosts(ss []string) (map[int32]*common.Rewards, error) {

	costs := map[int32]*common.Rewards{}

	for _, s := range ss {
		split := strings.SplitN(s, ":", 2)
		// 第一项代表区域
		areaId, err := strconv.ParseInt(split[0], 10, 32)
		if err != nil {
			return nil, errors.New("ToYggdrasilAreaBuildCosts Format error. (string: %s", s)
		}

		rewards, ok := costs[int32(areaId)]
		if !ok {
			rewards = common.NewRewards()
		}

		reward, err := StringToReward(split[1])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		if reward.Type != 0 {
			rewards.AddReward(reward)
		}

		costs[int32(areaId)] = rewards
	}

	return costs, nil
}

func ToPosition(x, y int32) *coordinate.Position {
	return coordinate.NewPosition(x, y)
}

func IntMinuteToSecond(m int32) int64 {
	if m < 0 {
		return int64(-1)
	}
	return int64(m * servertime.SecondPerMinute)
}

func DailyTimeToSec(timeStr string) (int32, error) {

	splits := strings.Split(timeStr, ":")
	hour, err := stringToInt32(splits[0])
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	min, err := stringToInt32(splits[1])
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	return hour*servertime.SecondPerHour + min*servertime.SecondPerMinute, nil
}

func StringsToEnvObjs(Id int32, CreateGroupObject, CreateDeleteAtGroupObj []string, DeleteAt int32) (*common.EnvObjs, error) {
	en := common.NewEnvObjs()
	for _, s := range CreateGroupObject {
		if len(s) == 0 {
			continue
		}
		e := &common.EnvObj{}
		err := FillEnvObj(Id, s, e, 0)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		*en = append(*en, e)
	}
	if len(CreateDeleteAtGroupObj) == 1 && CreateDeleteAtGroupObj[0] == "" && DeleteAt != 0 {
		return nil, errors.New("CreateDeleteAtGroupObj is nil but DeleteAt!=0, envId:%d", Id)
	}
	for _, s := range CreateDeleteAtGroupObj {
		if len(s) == 0 {
			continue
		}
		e := &common.EnvObj{}
		err := FillEnvObj(Id, s, e, DeleteAt)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		*en = append(*en, e)
	}
	return en, nil
}

func FillEnvObj(envId int32, str string, en *common.EnvObj, at int32) error {
	split := strings.Split(str, ":")

	if len(split) == 4 {

		X, err := stringToInt32(split[1])
		if err != nil {
			return errors.WrapTrace(err)
		}
		Y, err := stringToInt32(split[2])
		if err != nil {
			return errors.WrapTrace(err)
		}
		ObjectId, err := stringToInt32(split[3])
		if err != nil {
			return errors.WrapTrace(err)
		}
		en.Relative = true
		en.ObjectId = ObjectId
		en.Position.X = X
		en.Position.Y = Y
		en.DeleteAt = at
	} else if len(split) == 3 {

		X, err := stringToInt32(split[0])
		if err != nil {
			return errors.WrapTrace(err)
		}
		Y, err := stringToInt32(split[1])
		if err != nil {
			return errors.WrapTrace(err)
		}
		ObjectId, err := stringToInt32(split[2])
		if err != nil {
			return errors.WrapTrace(err)
		}
		en.ObjectId = ObjectId
		en.Position.X = X
		en.Position.Y = Y
		en.DeleteAt = at

	} else {
		return errors.New("FillEnvObj len err,envId:%d, str:%s,", envId, str)
	}
	return nil
}

func StringsToTaskItems(split []string) (*common.TaskItems, error) {
	taskItems := common.NewTaskItems()
	for _, s := range split {
		if len(s) == 0 {
			continue
		}
		e := &common.TaskItem{}
		err := FillTaskItem(s, e)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		*taskItems = append(*taskItems, e)
	}
	return taskItems, nil
}

func FillTaskItem(s string, t *common.TaskItem) error {
	split := strings.Split(s, ":")

	ID, err := stringToInt32(split[0])
	if err != nil {
		return errors.WrapTrace(err)
	}
	Num, err := stringToInt32(split[1])
	if err != nil {
		return errors.WrapTrace(err)
	}

	t.ID = ID
	t.Num = Num
	return nil
}

func StringsToCharacterSkill(strs []string) (map[int32]int32, error) {
	skillMap := map[int32]int32{}

	for _, str := range strs {
		splits := strings.Split(str, ":")
		if len(splits) < 2 {
			return skillMap, errors.New("Format error. (string: %s", str)
		}

		skillId, err := stringToInt32(splits[0])
		if err != nil {
			return skillMap, errors.WrapTrace(err)
		}

		skillLevel, err := stringToInt32(splits[1])
		if err != nil {
			return skillMap, errors.WrapTrace(err)
		}

		skillMap[skillId] = skillLevel
	}

	return skillMap, nil
}

func StringsToEnvTerrains(envId int32, terrains string, DeleteAt int32) (*common.EnvTerrains, error) {

	/**
	改变地形状态(type 1-地形 2 -nature)
	(type:posX:startY:endY:state,type:posX:startY:endY:state)
	*/
	if terrains == "" {
		return nil, nil
	}
	ret := &common.EnvTerrains{}
	ret.DeleteAt = DeleteAt
	for _, terrain := range strings.Split(terrains, ",") {
		if len(terrain) == 0 {
			continue
		}
		split := strings.Split(terrain, ":")
		switch split[0] {
		case "1":
		case "2":
			//todo:
		case "3":
			//todo:

		}
	}
	return ret, nil
}
