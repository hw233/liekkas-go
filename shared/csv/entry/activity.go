package entry

import (
	"shared/common"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/servertime"
	"shared/utility/transfer"
	"strconv"
	"sync"
	"time"
)

type Activity struct {
	Id               int32
	TimeType         int32
	StartTime        string
	EndTime          string
	UnlockConditions *common.Conditions `rule:"conditions"`
	CloseConditions  *common.Conditions `rule:"conditions"`
	StartTimeArg     int64              `ignore:"true"`
	EndTimeArg       int64              `ignore:"true"`
}

type ActivityFunc struct {
	Id                 int32
	ActivityId         int32
	FuncType           int32
	FuncArgs           []int32
	StartTimeOffset    []int32
	EndTimeOffset      []int32
	StartTimeSecOffset int64 `ignore:"true"`
	EndTimeSecOffset   int64 `ignore:"true"`
	EndItemRemove      []int32
}

type Activities struct {
	sync.RWMutex

	Activities    map[int32]*Activity
	ActivityFuncs map[int32]*ActivityFunc

	ActivityToFuncs map[int32]map[int32]int32
}

func NewActivities() *Activities {
	return &Activities{}
}

func arrayOffsetToSecOffset(arr []int32) int64 {
	offsets := make([]int64, 3)

	for i, v := range arr {
		offsets[i] = int64(v)
	}

	return offsets[0]*servertime.SecondPerDay + offsets[1]*servertime.SecondPerHour + offsets[2]*servertime.SecondPerMinute
}

func strToInt64(str string) (int64, error) {
	if str == "" {
		return 0, nil
	}

	val, err := strconv.ParseInt(str, 10, 64)
	return val, err
}

func (a *Activities) Check(config *Config) error {
	return nil
}

func (a *Activities) Reload(config *Config) error {
	a.Lock()
	defer a.Unlock()

	activities := map[int32]*Activity{}
	activityFuncs := map[int32]*ActivityFunc{}
	activityToFuncs := map[int32]map[int32]int32{}

	for _, csv := range config.CfgActivityConfig.GetAllData() {
		activityCfg := &Activity{}
		err := transfer.Transfer(csv, activityCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		switch activityCfg.TimeType {
		case static.ActivityTimeTypeDate:
			startTime, err := servertime.ParseTime(activityCfg.StartTime)
			if err != nil {
				return err
			}

			activityCfg.StartTimeArg = startTime

			endTime, err := servertime.ParseTime(activityCfg.EndTime)
			if err != nil {
				return err
			}

			activityCfg.EndTimeArg = endTime

		case static.ActivityTimeTypeUserCreate:
			start, err := strToInt64(activityCfg.StartTime)
			if err != nil {
				return err
			}
			activityCfg.StartTimeArg = int64((time.Hour * 24).Seconds()) * start

			end, err := strToInt64(activityCfg.EndTime)
			if err != nil {
				return err
			}
			activityCfg.EndTimeArg = int64((time.Hour * 24).Seconds()) * end
		}

		activities[activityCfg.Id] = activityCfg
	}

	for _, csv := range config.CfgActivityFuncConfig.GetAllData() {
		activityFuncCfg := &ActivityFunc{}
		err := transfer.Transfer(csv, activityFuncCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		activityFuncCfg.StartTimeSecOffset = arrayOffsetToSecOffset(activityFuncCfg.StartTimeOffset)
		activityFuncCfg.EndTimeSecOffset = arrayOffsetToSecOffset(activityFuncCfg.EndTimeOffset)

		activityId := activityFuncCfg.ActivityId
		funcIds, ok := activityToFuncs[activityId]
		if !ok {
			funcIds = map[int32]int32{}
			activityToFuncs[activityId] = funcIds
		}

		funcIds[activityFuncCfg.Id] = activityFuncCfg.Id
		activityFuncs[activityFuncCfg.Id] = activityFuncCfg
	}

	a.Activities = activities
	a.ActivityFuncs = activityFuncs
	a.ActivityToFuncs = activityToFuncs

	return nil
}

func (a *Activities) GetActivityFunc(id int32) (*ActivityFunc, error) {
	cfg, ok := a.ActivityFuncs[id]
	if !ok {
		return nil, errors.Swrapf(common.ErrActivityFuncCfgNotFound, id)
	}

	return cfg, nil
}

func (a *Activities) GetActivityByFuncId(funcId int32) (*Activity, error) {
	funcCfg, err := a.GetActivityFunc(funcId)
	if err != nil {
		return nil, err
	}

	activityId := funcCfg.ActivityId
	cfg, ok := a.Activities[activityId]
	if !ok {
		return nil, errors.Swrapf(common.ErrActivityCfgNotFound, activityId)
	}

	return cfg, nil
}

func (a *Activities) GetActivityFuncIdsByActivity(activityId int32) (map[int32]int32, error) {
	funcIds, ok := a.ActivityToFuncs[activityId]
	if !ok {
		return nil, errors.Swrapf(common.ErrActivityCfgNotFound, activityId)
	}

	return funcIds, nil
}

func (a *Activities) GetAllActivities() map[int32]*Activity {
	return a.Activities
}
