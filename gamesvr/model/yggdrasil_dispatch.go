package model

import (
	"math"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/servertime"
)

const (
	YggdrasilDispatchStateReadyForMission     = 0
	YggdrasilDispatchStateOnMission           = 1
	YggdrasilDispatchStateMissionComplete     = 2
	YggdrasilDispatchStateMissionAfterRewards = 3
)

type YggdrasilDispatches struct {
	*DailyRefreshChecker
	DailyTasks    map[int32]*YggdrasilDailyDispatch `json:"daily_tasks"`
	GuildTasks    map[int32]*YggdrasilGuildDispatch `json:"guild_tasks"`
	IfUpdateGuild bool                              `json:"if_update_guild"`
	// DelayTasks          map[int32]*YggDispatchAreaStar // key: 区域Id
	// CharacStateForDaily map[int32]int32 // 为true代表处于派遣中，为false则是未派遣
	// CharacStateForGuild map[int32]int32
}

type YggdrasilDailyDispatch struct {
	Id           int32   `json:"id"`
	EndTime      int64   `json:"end_time"`
	TaskState    int32   `json:"task_state"`
	ExtraRewards []int32 `json:"extra_rewards"`
	CharacterId  []int32 `json:"character_ids"`
}

type YggdrasilGuildDispatch struct {
	Id           int32   `json:"id"`
	EndTime      int64   `json:"end_time"`
	TaskState    int32   `json:"task_state"`
	ExtraRewards []int32 `json:"extra_rewards"`
	CloseTime    int64   `json:"close_time"`
	CharacterId  []int32 `json:"character_ids"`
}

func NewYggDispatches() *YggdrasilDispatches {
	return &YggdrasilDispatches{
		DailyRefreshChecker: NewDailyRefreshChecker(),
		DailyTasks:          map[int32]*YggdrasilDailyDispatch{},
		GuildTasks:          map[int32]*YggdrasilGuildDispatch{},
		IfUpdateGuild:       false,
	}
}

func NewYggdrasilDispatch(taskId int32) *YggdrasilDailyDispatch {
	return &YggdrasilDailyDispatch{
		Id:           taskId,
		EndTime:      -1,
		TaskState:    YggdrasilDispatchStateReadyForMission,
		ExtraRewards: []int32{},
		CharacterId:  []int32{},
	}
}

func NewYggdrasilGuildDispatch(taskId int32, closeTime int64) *YggdrasilGuildDispatch {
	return &YggdrasilGuildDispatch{
		Id:           taskId,
		EndTime:      -1,
		TaskState:    YggdrasilDispatchStateReadyForMission,
		CloseTime:    closeTime + servertime.Now().Unix(),
		ExtraRewards: []int32{},
		CharacterId:  []int32{},
	}
}

// -----------------YggdrasilDispatches------------

func (yd *YggdrasilDispatches) CancelDispatch(taskId int32, dispatchType int32) (*pb.VOYggDispatchTaskState, error) {
	if dispatchType == 1 {
		task, ok := yd.DailyTasks[taskId]
		if !ok {
			return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskNotFound, taskId)
		}
		if !task.CheckStateOnMission() {
			return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskStateNotOnMission, taskId, task.TaskState)
		}
		yd.DailyTasks[taskId] = NewYggdrasilDispatch(taskId)
		return task.VOYggDispatchTaskState(), nil
	} else {
		task, ok := yd.GuildTasks[taskId]
		if !ok {
			return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskNotFound, taskId)
		}
		if !task.CheckStateOnMission() {
			return nil, errors.Swrapf(common.ErrYggdrasilDispatchTaskStateNotOnMission, taskId, task.TaskState)
		}
		// todo closeTime 如何处理
		yd.GuildTasks[taskId] = NewYggdrasilGuildDispatch(taskId, task.CloseTime)
		return task.VOYggDispatchTaskState(), nil
	}
}

// 新增任务
func (yd *YggdrasilDispatches) AddNewDailyTask(taskId int32) {
	yd.DailyTasks[taskId] = NewYggdrasilDispatch(taskId)
}
func (yd *YggdrasilDispatches) AddNewGuildTask(taskId int32, closeTime int64) {
	yd.GuildTasks[taskId] = NewYggdrasilGuildDispatch(taskId, closeTime)
}

// 删除任务
func (yd *YggdrasilDispatches) DeleteDailyTask(taskId int32) {
	delete(yd.DailyTasks, taskId)
}
func (yd *YggdrasilDispatches) DeleteGuildTask(taskId int32) {
	delete(yd.GuildTasks, taskId)
}

func (yd *YggdrasilDispatches) DailyDispatchComplete(taskId int32) {
	// 将角色对应的任务给置0
	// for characId, taskStateId := range yd.CharacStateForDaily {
	// 	if taskStateId == taskId {
	// 		yd.CharacStateForDaily[characId] = 0
	// 	}
	// }
	yd.DeleteDailyTask(taskId)
}

func (yd *YggdrasilDispatches) GuildDispatchComplete(taskId int32) {

}

// 刷新日常任务状态
func (yd *YggdrasilDispatches) YggDispatchUpdateDailyStates() {
	for _, task := range yd.DailyTasks {
		if task.TaskState == YggdrasilDispatchStateOnMission && task.EndTime <= servertime.Now().Unix() {
			task.NextState()
			yd.DailyTasks[task.Id] = task
		}
	}
}

// 刷新公会任务状态
func (yd *YggdrasilDispatches) YggDispatchUpdateGuildStates() {
	for _, task := range yd.GuildTasks {
		// 关闭过期的任务
		if task.TaskState == YggdrasilDispatchStateReadyForMission && task.CloseTime <= servertime.Now().Unix() {
			delete(yd.GuildTasks, task.Id)
		}
		// 把派遣完成的状态改为待领取奖励
		if task.TaskState == YggdrasilDispatchStateOnMission && task.EndTime <= servertime.Now().Unix() {
			task.NextState()
			yd.GuildTasks[task.Id] = task
		}
	}

}

// 派遣刷新，对于日常派遣，删除所有待派遣的任务；对于公会派遣，将是否更新置true
func (yd *YggdrasilDispatches) YggDispatchDailyRefresh() {
	//
	for _, task := range yd.DailyTasks {
		if task.TaskState == YggdrasilDispatchStateReadyForMission || task.TaskState == YggdrasilDispatchStateMissionAfterRewards {
			delete(yd.DailyTasks, task.Id)
		}
	}
	yd.IfUpdateGuild = true
}

// 返回值为真，代表这个角色在派遣中
func (yd *YggdrasilDispatches) FindCharacterDailyOnMission(characId int32) bool {
	for _, task := range yd.DailyTasks {
		for _, id := range task.CharacterId {
			if characId == id {
				return true
			}
		}
	}
	return false
}

// ----------------YggdrasilDispatch 日常派遣单个任务
func (ydd *YggdrasilDailyDispatch) YggDispatchUpdateTask() {
	if ydd.TaskState == YggdrasilDispatchStateOnMission && ydd.EndTime <= servertime.Now().Unix() {
		ydd.NextState()
	}
}

// 返回true，代表处于待派遣状态
func (ydd *YggdrasilDailyDispatch) CheckStateReady() bool {
	return ydd.TaskState == YggdrasilDispatchStateReadyForMission
}

// 返回true，代表处于派遣中状态
func (ydd *YggdrasilDailyDispatch) CheckStateOnMission() bool {
	return ydd.TaskState == YggdrasilDispatchStateOnMission
}

// 返回true，代表处于待领取奖励状态
func (ydd *YggdrasilDailyDispatch) CheckStateMissionComplete() bool {
	ydd.YggDispatchUpdateTask()
	return ydd.TaskState == YggdrasilDispatchStateMissionComplete
}

// 返回true，代表返回afterRewards状态
func (ydd *YggdrasilDailyDispatch) CheckStateAfterRewards() bool {
	return ydd.TaskState == YggdrasilDispatchStateMissionAfterRewards
}

func (ydd *YggdrasilDailyDispatch) StateReadyToGo(timeCost int64, userCharacters []int32) {
	ydd.EndTime = servertime.Now().Unix() + timeCost
	ydd.CharacterId = userCharacters
	ydd.NextState()
}

// 修改状态
func (ydd *YggdrasilDailyDispatch) NextState() {
	ydd.TaskState = (ydd.TaskState + 1) % 4
	if ydd.TaskState != YggdrasilDispatchStateOnMission {
		ydd.EndTime = -1
	}
	if ydd.TaskState == YggdrasilDispatchStateMissionAfterRewards {
		ydd.CharacterId = []int32{}
	}
}

// 获得剩余的秒数
func getCountDownSec(EndAt int64) int64 {
	return int64(math.Max(0, float64(EndAt-servertime.Now().Unix())))
}

func (ydd *YggdrasilDailyDispatch) VOYggDispatchTaskState() *pb.VOYggDispatchTaskState {
	restTime := int64(-1)
	if ydd.TaskState == YggdrasilDispatchStateOnMission {
		restTime = getCountDownSec(ydd.EndTime)
	}
	return &pb.VOYggDispatchTaskState{
		TaskId:    ydd.Id,
		TaskState: ydd.TaskState,
		RestTime:  restTime,
		CloseTime: -1,
	}
}

//----------------YggdrasilGuildDispatch---------------

func (ygd *YggdrasilGuildDispatch) VOYggDispatchTaskState() *pb.VOYggDispatchTaskState {
	restTime := int64(-1)
	if ygd.TaskState == YggdrasilDispatchStateOnMission {
		restTime = getCountDownSec(ygd.EndTime)
	}
	closeTime := int64(-1)
	if ygd.TaskState == YggdrasilDispatchStateReadyForMission {
		closeTime = getCountDownSec(ygd.CloseTime)
	}
	return &pb.VOYggDispatchTaskState{
		TaskId:    ygd.Id,
		TaskState: ygd.TaskState,
		RestTime:  restTime,
		CloseTime: closeTime,
	}
}

func (ygd *YggdrasilGuildDispatch) StateReadyToGo(timeCost int64, userCharacters []int32) {
	ygd.EndTime = servertime.Now().Unix() + timeCost
	ygd.CharacterId = userCharacters
	ygd.NextState()
}

func (ygd *YggdrasilGuildDispatch) NextState() {
	ygd.TaskState = (ygd.TaskState + 1) % 3
	// if ygd.TaskState != YggdrasilDispatchStateReadyForMission {
	// 	ygd.CloseTime = -1
	// }
	if ygd.TaskState != YggdrasilDispatchStateOnMission {
		ygd.EndTime = -1
	}
}

// 返回true，代表处于待派遣状态
func (ygd *YggdrasilGuildDispatch) CheckStateReady() bool {
	if ygd.CloseTime <= servertime.Now().Unix() {
		return false
	}
	return ygd.TaskState == YggdrasilDispatchStateReadyForMission
}

// 返回true，代表处于派遣中状态
func (ygd *YggdrasilGuildDispatch) CheckStateOnMission() bool {
	if ygd.EndTime <= servertime.Now().Unix() {
		return false
	}
	return ygd.TaskState == YggdrasilDispatchStateOnMission
}

// 返回true，代表处于待领取奖励状态
func (ygd *YggdrasilGuildDispatch) CheckStateMissionComplete() bool {
	if ygd.TaskState == YggdrasilDispatchStateOnMission && ygd.EndTime <= servertime.Now().Unix() {
		ygd.NextState()
	}
	return ygd.TaskState == YggdrasilDispatchStateMissionComplete
}
