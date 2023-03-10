package model

type ActivityState int

const (
	ActivityStateWaiting ActivityState = iota
	ActivityStateProgressing
	ActivityStateEnded
)

type ActivityFunc struct {
	Id            int32         `json:"id"`
	State         ActivityState `json:"state"`
	StartTime     int64         `json:"start_time"`
	EndTime       int64         `json:"end_time"`
	RealStartTime int64         `json:"real_start_time"`
	RealEndTime   int64         `json:"real_end_time"`
}

type ActivityInfo struct {
	ActivityFuncs map[int32]*ActivityFunc `json:"activity_funcs"`

	WaitingFuncs     map[int32]int32 `json:"-"`
	ProgressingFuncs map[int32]int32 `json:"-"`
}

func NewActivityInfo() *ActivityInfo {
	return &ActivityInfo{
		ActivityFuncs:    map[int32]*ActivityFunc{},
		WaitingFuncs:     map[int32]int32{},
		ProgressingFuncs: map[int32]int32{},
	}
}

// ActivityInfo
func (ai *ActivityInfo) GetActivityFunc(id int32) (*ActivityFunc, bool) {
	activityFunc, ok := ai.ActivityFuncs[id]
	return activityFunc, ok
}

func (ai *ActivityInfo) AddActivityFunc(id int32, startTime, endTime int64) *ActivityFunc {
	_, ok := ai.ActivityFuncs[id]
	if ok {
		return nil
	}

	af := &ActivityFunc{
		Id:        id,
		State:     ActivityStateWaiting,
		StartTime: startTime,
		EndTime:   endTime,
	}

	ai.ActivityFuncs[id] = af
	ai.WaitingFuncs[id] = id

	return af
}

func (ai *ActivityInfo) StartActivityFunc(id int32, startTime int64) {
	activityFunc, ok := ai.GetActivityFunc(id)
	if !ok {
		return
	}

	activityFunc.start(startTime)

	delete(ai.WaitingFuncs, id)
	ai.ProgressingFuncs[id] = id
}

func (ai *ActivityInfo) EndActivityFunc(id int32, endTime int64) {
	activityFunc, ok := ai.GetActivityFunc(id)
	if !ok {
		return
	}

	activityFunc.end(endTime)

	delete(ai.WaitingFuncs, id)
	delete(ai.ProgressingFuncs, id)
}

func (ai *ActivityInfo) GetWaitingActivityFuncs() []int32 {
	list := []int32{}
	for id := range ai.WaitingFuncs {
		list = append(list, id)
	}

	return list
}

func (ai *ActivityInfo) GetProgressingActivityFuncs() []int32 {
	list := []int32{}
	for id := range ai.ProgressingFuncs {
		list = append(list, id)
	}

	return list
}

// ActivityFunc
func (af *ActivityFunc) ResetTime(startTime, endTime int64) {
	af.StartTime = startTime
	af.EndTime = endTime
}

func (af *ActivityFunc) IsArrivalStartTime(now int64) bool {
	return now >= af.StartTime && now < af.EndTime
}

func (af *ActivityFunc) IsArrivalEndTime(now int64) bool {
	return now >= af.EndTime
}

func (af *ActivityFunc) IsWaiting() bool {
	return af.State == ActivityStateWaiting
}

func (af *ActivityFunc) IsProgressing() bool {
	return af.State == ActivityStateProgressing
}

func (af *ActivityFunc) IsEnded() bool {
	return af.State >= ActivityStateEnded
}

func (af *ActivityFunc) start(startTime int64) {
	af.State = ActivityStateProgressing
	af.RealStartTime = startTime
}

func (af *ActivityFunc) end(endTime int64) {
	af.State = ActivityStateEnded
	af.RealEndTime = endTime
}
