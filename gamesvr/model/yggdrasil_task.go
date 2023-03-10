package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/coordinate"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/number"
)

type YggdrasilTask struct {
	CompleteTaskIds *number.NonRepeatableArr     `json:"complete_task_ids"`
	TaskInProgress  map[int32]*YggdrasilTaskInfo `json:"task_in_progress"`
	TrackTaskId     int32                        `json:"track_task_id"`
}

func NewYggdrasilTask() *YggdrasilTask {
	return &YggdrasilTask{
		CompleteTaskIds: number.NewNonRepeatableArr(),
		TaskInProgress:  map[int32]*YggdrasilTaskInfo{},
		TrackTaskId:     0,
	}
}

func (y *YggdrasilTask) VOYggdrasilTaskTotalInfo() *pb.VOYggdrasilTaskTotalInfo {
	taskInfos := make([]*pb.VOYggdrasilTaskInfo, 0, len(y.TaskInProgress))
	for _, info := range y.TaskInProgress {
		taskInfos = append(taskInfos, info.VOYggdrasilTaskInfo())
	}

	return &pb.VOYggdrasilTaskTotalInfo{
		CompletedTasks: y.CompleteTaskIds.Values(),
		TaskInfos:      taskInfos,
		TrackTask:      y.FetchTrackTaskId(),
	}
}

type YggdrasilTaskInfo struct {
	TaskId             int32                    `json:"task_id"`
	TemporalObjects    map[int64]*YggObject     `json:"temporal_objects"` // 刚接取任务的时候 obj的状态和位置，任务失败或者放弃的时候重置
	CompleteSubTaskIds *number.NonRepeatableArr `json:"complete_sub_task_ids"`

	*YggSubTaskProgressInfoWrap
}
type YggSubTaskProgressInfoWrap struct {
	Base  *YggSubTaskProgressInfo   `json:"base"`
	Multi []*YggSubTaskProgressInfo `json:"multi_sub_task_info"` // 一个子任务有多个子任务组成的情况(保证有序所以用slice)
}

func (y *YggSubTaskProgressInfoWrap) Vo() (*pb.VOYggdrasilSubTaskInfo, []*pb.VOYggdrasilSubTaskInfo) {
	multi := make([]*pb.VOYggdrasilSubTaskInfo, 0, len(y.Multi))

	for _, info := range y.Multi {
		multi = append(multi, info.VOYggdrasilSubTaskInfo())
	}
	return y.Base.VOYggdrasilSubTaskInfo(), multi
}

func NewYggSubTaskProgressInfoWrap(base *YggSubTaskProgressInfo, multi []*YggSubTaskProgressInfo) *YggSubTaskProgressInfoWrap {
	return &YggSubTaskProgressInfoWrap{
		Base:  base,
		Multi: multi,
	}
}

type YggSubTaskProgressInfo struct {
	SubTaskId          int32                             `json:"sub_task_id"`
	IsSubTaskCompleted bool                              `json:"is_sub_task_completed"`
	Processes          []*common.YggdrasilSubTaskProcess `json:"processes"`
	EnvId              int32                             `json:"env_id"`
}

func (y *YggSubTaskProgressInfo) VOYggdrasilSubTaskInfo() *pb.VOYggdrasilSubTaskInfo {
	vo := make([]*pb.VOYggdrasilSubTaskProcess, 0, len(y.Processes))

	for _, process := range y.Processes {
		vo = append(vo, process.VOYggdrasilSubTaskProcess())
	}

	return &pb.VOYggdrasilSubTaskInfo{
		SubTaskId:          y.SubTaskId,
		IsSubTaskCompleted: y.IsSubTaskCompleted,
		TaskProcesses:      vo,
		EnvId:              y.EnvId,
	}
}

func NewYggSubTaskProgressInfo() *YggSubTaskProgressInfo {
	return &YggSubTaskProgressInfo{
		IsSubTaskCompleted: false,
	}
}

func NewYggdrasilTaskInfo(taskId int32) *YggdrasilTaskInfo {
	return &YggdrasilTaskInfo{
		TaskId:                     taskId,
		TemporalObjects:            map[int64]*YggObject{},
		CompleteSubTaskIds:         number.NewNonRepeatableArr(),
		YggSubTaskProgressInfoWrap: nil,
	}

}

func (y *YggdrasilTaskInfo) VOYggdrasilTaskInfo() *pb.VOYggdrasilTaskInfo {

	base, multi := y.Vo()
	return &pb.VOYggdrasilTaskInfo{
		TaskId:             y.TaskId,
		SubTaskId:          y.Base.SubTaskId,
		TaskMain:           base,
		Multi:              multi,
		CompleteSubTaskIds: y.CompleteSubTaskIds.Values(),
	}

}

// TaskAbandon 任务放弃
func (y *YggdrasilTask) TaskAbandon(ctx context.Context, u *User, taskId int32) error {
	// 判断任务是否可以放弃
	config, err := manager.CSV.Yggdrasil.GetTaskConfig(taskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if !config.AbandonTask {
		return errors.WrapTrace(common.ErrYggdrasilTaskCannotAbandon)
	}
	info, ok := y.TaskInProgress[taskId]
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilTaskNotInProgress)
	}
	return y.TaskFail(ctx, u.ID, u.Yggdrasil, info)

}

// TaskFail  任务失败时候 要做的事情   1.任务物品删除  2.env覆盖的obj恢复 3.创建的env删除 4.还原当前子任务期间操作过的任务obj的状态和位置
func (y *YggdrasilTask) TaskFail(ctx context.Context, userId int64, yggdrasil *Yggdrasil, info *YggdrasilTaskInfo) error {

	allSubTaskIds := number.NewNonRepeatableArr()
	allSubTaskIds.Append(info.CompleteSubTaskIds.Values()...)
	allSubTaskIds.Append(info.Base.SubTaskId)
	for _, multi := range info.Multi {
		allSubTaskIds.Append(multi.SubTaskId)
	}

	for _, subTaskId := range allSubTaskIds.Values() {
		//1. 任务物品删除
		yggdrasil.EntityChange.TaskItemChanges(yggdrasil.TaskPack.Remove(subTaskId))

		// 2.env覆盖的obj恢复
		err := yggdrasil.RecoverEnv(ctx, userId, subTaskId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// 3.创建的env删除
		for _, object := range yggdrasil.Entities.EnvObjectCreateAt[subTaskId] {
			yggdrasil.RemoveObject(object)
		}
		yggdrasil.Entities.RemoveEnvTerrainCreateAt(subTaskId, yggdrasil)
	}

	// 4.还原期间操作过的任务obj的状态和位置
	for uid, before := range info.TemporalObjects {
		now, ok := yggdrasil.Entities.FindObjectByUid(uid)
		if ok {
			// 如果位置或者状态改变了 则还原
			yggdrasil.RecoverFromBefore(now, before)
		} else {
			// 如果被删除了则新增
			yggdrasil.AppendObject(ctx, before)

		}
	}
	// 整个任务变成未领取状态
	delete(y.TaskInProgress, info.TaskId)
	return nil
}

// InitSubTask 初始化subtask
func (y *YggdrasilTask) InitSubTask(ctx context.Context, user *User, nextSubTaskId int32) (*YggdrasilTaskInfo, error) {
	subTaskConfig, err := manager.CSV.Yggdrasil.GetSubTaskConfig(nextSubTaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	taskInfo, ok := y.TaskInProgress[subTaskConfig.TaskId]
	if !ok {
		taskInfo = NewYggdrasilTaskInfo(subTaskConfig.TaskId)
	}
	// 获得任务初始进度
	process, err := y.GetTaskInitProcess(ctx, user, nextSubTaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 赋值任务初始进度
	taskInfo.YggSubTaskProgressInfoWrap = process
	y.TaskInProgress[taskInfo.TaskId] = taskInfo

	return taskInfo, nil
}

func (y *YggdrasilTask) GetTaskInitProcess(ctx context.Context, user *User, nextSubTaskId int32) (*YggSubTaskProgressInfoWrap, error) {
	subTaskConfig, err := manager.CSV.Yggdrasil.GetSubTaskConfig(nextSubTaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	env, err := subTaskConfig.Envs.RandomEnv()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 先初始化env 有时handler.init会生成env生成的object的任务进度
	err = y.InitEnv(ctx, user, nextSubTaskId, env)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	handler, err := TaskHandlers.getHandler(env.SubTaskType)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 获得任务初始进度
	process, err := handler.init(ctx, env, user)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return process, nil
}

//WhenSubTaskComplete 子任务完成时 1.添加到已完成子任务列表 2.env暂时删除的entity恢复 3.env创造的obj删除 4. 子任务发奖
func (y *YggdrasilTask) WhenSubTaskComplete(ctx context.Context, u *User, info *YggdrasilTaskInfo, progressInfo *YggSubTaskProgressInfo) error {
	completeSubTaskId := progressInfo.SubTaskId
	config, err := manager.CSV.Yggdrasil.GetSubTaskConfig(completeSubTaskId)
	if err != nil {
		return errors.WrapTrace(err)
	}

	// 添加到已完成子任务列表
	info.CompleteSubTaskIds.Append(completeSubTaskId)

	err = u.Yggdrasil.RecoverEnv(ctx, u.ID, completeSubTaskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// env创造的obj删除
	for _, object := range u.Yggdrasil.Entities.EnvObjectDeleteAt[completeSubTaskId] {
		u.Yggdrasil.RemoveObject(object)
	}
	// 子任务发奖
	if config.YggDropId > 0 {
		reason := logreason.NewReason(logreason.YggSubTask)
		err = u.Yggdrasil.AddRewardsByDropId(ctx, u, config.YggDropId, completeSubTaskId, reason)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	if config.DropId > 0 {
		reason := logreason.NewReason(logreason.YggSubTask)
		_, err = u.AddRewardsByDropId(config.DropId, reason)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// 隐藏式给玩家
		u.RewardsResult.Rewards = common.NewRewards()

	}
	// 移除env产生的地形
	u.Yggdrasil.Entities.RemoveEnvTerrainDeleteAt(completeSubTaskId, u.Yggdrasil)

	return nil
}

// InitEnv 初始化env
func (y *YggdrasilTask) InitEnv(ctx context.Context, user *User, subTaskId int32, env *entry.Env) error {

	// 先清理地块 再刷怪
	for _, position := range env.ClearPosGroup.Points() {
		// 暂时移除
		err := user.Yggdrasil.EnvRemove(ctx, user.ID, position, subTaskId, false)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	// 填对应subtaskID就销毁，填0or不填自动销毁。填-1永久不删除。除非走clearPosGroup流程删除
	for _, obj := range *env.CreateObjects {
		var pos coordinate.Position
		if obj.Relative {
			pos = *coordinate.NewPosition(obj.Position.X+user.Yggdrasil.TravelPos.X, obj.Position.Y+user.Yggdrasil.TravelPos.Y)
		} else {
			pos = obj.Position
		}
		// 暂时移除
		err := user.Yggdrasil.EnvRemove(ctx, user.ID, pos, subTaskId, true)
		if err != nil {
			return errors.WrapTrace(err)
		}
		objectConfig, err := manager.CSV.Yggdrasil.GetYggdrasilObjectConfig(obj.ObjectId)
		if err != nil {
			return errors.WrapTrace(err)
		}
		newObject, err := NewYggObject(ctx, pos, obj.ObjectId, objectConfig.DefaultState)
		if err != nil {
			return errors.WrapTrace(err)
		}
		// 填0or不填自动销毁。填-1永久不删除。
		if obj.DeleteAt == 0 {
			newObject.SetCreateAndDeleteAt(subTaskId, subTaskId)
		} else {
			newObject.SetCreateAndDeleteAt(subTaskId, obj.DeleteAt)

		}
		// 在初始化obj的时候可能会初始化area
		user.Yggdrasil.AppendObject(ctx, newObject)

	}

	// 添加任务道具
	user.Yggdrasil.AddRewardsInTaskPack(ctx, user, TaskItemToRewards(env.AddTaskItem), subTaskId)
	// 删除任务道具
	user.Yggdrasil.CostYggdrasilTaskPack(TaskItemToRewards(env.DeleteTaskItem))

	//  改变obj状态
	for _, objectId := range env.ChangeObjectState {
		objs, err := user.Yggdrasil.FindObjectById(objectId)
		if err != nil {
			glog.Errorf("InitEnv ChangeObjectState FindObjectById,objectId:%d,err:%d", objectId, err)
			continue
		}
		for _, obj := range objs {
			user.Yggdrasil.ObjectChangeToNextState(ctx, user, obj)
		}
	}
	// 创建env产生的地形
	if env.Terrains != nil {
		// 不填or填0表示随subtask
		//自动删除，填-1表示永不删除
		if env.Terrains.DeleteAt == 0 {

			user.Yggdrasil.Entities.CreateEnvTerrain(subTaskId, subTaskId, env.Id, user.Yggdrasil)
		} else {
			user.Yggdrasil.Entities.CreateEnvTerrain(subTaskId, env.Terrains.DeleteAt, env.Id, user.Yggdrasil)

		}
	}

	return nil
}

func TaskItemToRewards(items *common.TaskItems) *common.Rewards {
	rewards := common.NewRewards()
	for _, item := range *items {
		rewards.AddReward(common.NewReward(item.ID, item.Num))
	}
	return rewards
}

func (y *YggdrasilTask) ProcessNum(ctx context.Context, u *User, subTaskType, process int32, filterCondition ...int32) {
	y.process(ctx, u, subTaskType, []int32{process}, filterCondition...)
}

func (y *YggdrasilTask) ProcessPos(ctx context.Context, u *User, subTaskType int32, position coordinate.Position, filterCondition ...int32) {
	y.process(ctx, u, subTaskType, []int32{position.X, position.Y}, filterCondition...)

}

// 流转任务进度 注：只单纯做任务进度处理和推送
func (y *YggdrasilTask) process(ctx context.Context, u *User, subTaskType int32, process []int32, filterCondition ...int32) {
	for _, info := range y.TaskInProgress {
		if info.YggSubTaskProgressInfoWrap == nil {
			// 任务进度还没初始化完成
			glog.Debugf("process pass,cause YggSubTaskProgressInfoWrap not inited")
			continue
		}
		// 先流转子任务的子任务
		multiChanged := y.processMulti(ctx, u, info, subTaskType, process, filterCondition...)
		// 流转子任务
		changed := y.process_(ctx, u, info, info.Base, subTaskType, process, filterCondition...)

		if !multiChanged && !changed {
			continue
		}
		// 推送任务进度
		u.AddYggPush(&pb.S2CYggdrasilTaskInfoChange{
			TaskInfo:       info.VOYggdrasilTaskInfo(),
			EntityChange:   u.VOYggdrasilEntityChange(ctx),
			ResourceResult: u.VOResourceResult(),
		})

	}
}

// 进入任务监听，返回进度是否有改变
func (y *YggdrasilTask) process_(ctx context.Context, u *User, info *YggdrasilTaskInfo, progressInfo *YggSubTaskProgressInfo, subTaskType int32, process []int32, filterCondition ...int32) bool {
	subTaskConfig, err := manager.CSV.Yggdrasil.GetSubTaskConfig(progressInfo.SubTaskId)
	if err != nil {
		glog.Errorf("Task process err:%+v", err)
		return false
	}
	env := subTaskConfig.Envs.Env(progressInfo.EnvId)
	if env.SubTaskType != subTaskType {
		return false
	}

	// 如果已经完成 continue
	if progressInfo.IsSubTaskCompleted {
		return false
	}

	handler, err := TaskHandlers.getHandler(subTaskType)
	if err != nil {
		glog.Errorf("Task process err:%+v", err)
		return false

	}
	// 过滤条件
	index, ok := progressInfo.filter(env.SubTaskTargets, filterCondition)
	if !ok {
		return false
	}
	handler.process(progressInfo, index, env, process, filterCondition)
	if progressInfo.IsSubTaskCompleted {
		//  子任务完成时 1.添加到已完成子任务列表 2.env暂时删除的entity恢复 3.env创造的obj删除 4. 子任务发奖
		err := y.WhenSubTaskComplete(ctx, u, info, progressInfo)
		if err != nil {
			glog.Errorf("Task process err:%+v", err)
		}
	}
	return true

}

func (y *YggdrasilTask) processMulti(ctx context.Context, u *User, info *YggdrasilTaskInfo, subTaskType int32, process []int32, filterCondition ...int32) bool {

	var changed bool
	// 先流转子任务的子任务
	for _, multi := range info.Multi {
		changed = y.process_(ctx, u, info, multi, subTaskType, process, filterCondition...) || changed

		// 子任务的子任务完成后，流转子任务,写这里处理，不放在process_中防止嵌套递归
		if multi.IsSubTaskCompleted {
			y.process_(ctx, u, info, info.Base, static.YggdrasilSubTaskTypeTypeMultiSubTask, []int32{}, multi.SubTaskId)
		}

	}
	return changed
}

// 流转任务失败
func (y *YggdrasilTask) processFail(ctx context.Context, u *User, subTaskType int32, filterCondition ...int32) {
	//todo: 子任务包含子任务的情况 任务失败
	for _, info := range y.TaskInProgress {
		subTaskConfig, err := manager.CSV.Yggdrasil.GetSubTaskConfig(info.Base.SubTaskId)
		if err != nil {
			glog.Errorf("Task processFail err:%+v", err)
			continue
		}
		env := subTaskConfig.Envs.Env(info.Base.EnvId)
		if env.SubTaskType != subTaskType {
			continue
		}
		// 过滤条件
		index, ok := info.Base.filter(env.SubTaskTargets, filterCondition)
		if !ok {
			continue
		}
		// 如果已经完成 continue
		if info.Base.IsSubTaskCompleted {
			continue
		}

		handler, err := TaskFailHandler.getHandler(subTaskType)
		if err != nil {
			glog.Errorf("Task processFail err:%+v", err)
			continue
		}
		// 任务失败
		if handler.processFail(info.Base.Processes[index].Attach) {
			// yggTask Fail
			err := y.TaskFail(ctx, u.ID, u.Yggdrasil, info)
			if err != nil {
				glog.Errorf("Task processFail  err:%+v", err)

			} else {
				// 推送
				u.AddYggPush(&pb.S2CYggdrasilTaskFail{
					TaskId:       info.TaskId,
					EntityChange: u.VOYggdrasilEntityChange(ctx),
				})
			}

		}

	}
}

func (y *YggdrasilTask) AcceptTask(ctx context.Context, u *User, taskConfig *entry.CfgYggdrasilTask) (*YggdrasilTaskInfo, error) {
	// 任务已完成
	if y.CompleteTaskIds.Contains(taskConfig.Id) {
		return nil, errors.WrapTrace(common.ErrYggdrasilCompleteTaskBefore)
	}
	// 任务已接取
	if _, ok := y.TaskInProgress[taskConfig.Id]; ok {
		return nil, errors.WrapTrace(common.ErrYggdrasilAcceptTaskBefore)
	}
	var otherCount int32 = 0
	// 互斥任务组
	if taskConfig.TaskGroup != 0 {
		for i := range y.TaskInProgress {
			temp, err := manager.CSV.Yggdrasil.GetTaskConfig(i)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			if temp.TaskGroup == taskConfig.TaskGroup {
				return nil, errors.WrapTrace(common.ErrYggdrasilSameTaskGroup)
			}
			if temp.TaskType != static.YggdrasilTaskTypeArea {
				otherCount++
			}
		}
	}
	if otherCount >= manager.CSV.Yggdrasil.GetYggOtherMaxTaskCount() {
		return nil, errors.WrapTrace(common.ErrYggdrasilOtherMaxTaskCount)
	}

	// 奇遇和支线任务是否开启
	if taskConfig.TaskType == static.YggdrasilTaskTypeHostel {
		err := u.CheckActionUnlock(static.ActionIdTypeHosteltask)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	} else if taskConfig.TaskType == static.YggdrasilTaskTypeAdventure {
		err := u.CheckActionUnlock(static.ActionIdTypeAdventure)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}

	return y.InitSubTask(ctx, u, taskConfig.NextSubTaskId)

}

func (y *YggdrasilTask) SetTrackTask(taskId int32) error {
	if taskId != 0 {
		// 任务已接取
		if _, ok := y.TaskInProgress[taskId]; !ok {
			return errors.WrapTrace(common.ErrYggdrasilNotDoingTaskNow)
		}
	}
	y.TrackTaskId = taskId
	return nil
}

func (y *YggdrasilTask) FetchTrackTaskId() int32 {
	if y.TrackTaskId == 0 {
		return 0
	}
	if _, ok := y.TaskInProgress[y.TrackTaskId]; !ok {
		return 0
	}
	return y.TrackTaskId
}

func (y *YggdrasilTask) CompleteTask(ctx context.Context, u *User, taskId int32) error {

	info, ok := y.TaskInProgress[taskId]
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilTaskNotInProgress)
	}
	subTaskConfig, err := manager.CSV.Yggdrasil.GetSubTaskConfig(info.Base.SubTaskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	env := subTaskConfig.Envs.Env(info.Base.EnvId)

	// 子任务未完成
	handler, err := TaskHandlers.getHandler(env.SubTaskType)
	if err != nil {
		return errors.WrapTrace(err)
	}
	if !handler.IsComplete(info.Base) {
		return errors.WrapTrace(common.ErrYggdrasilTaskCannotComplete)
	}
	// 不是最后一个子任务
	if !subTaskConfig.NextSubTaskIds.IsEmpty() {
		return errors.WrapTrace(common.ErrYggdrasilTaskCannotComplete)
	}

	return y.ForceCompleteTask(ctx, u, taskId)
}

func (y *YggdrasilTask) ForceCompleteTask(ctx context.Context, u *User, taskId int32) error {
	config, err := manager.CSV.Yggdrasil.GetTaskConfig(taskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 设置成已完成任务
	y.CompleteTaskIds.Append(taskId)
	delete(y.TaskInProgress, taskId)
	// 任务监听
	u.Yggdrasil.Task.ProcessNum(ctx, u, static.YggdrasilSubTaskTypeTypeMultiTask, 1, taskId)
	if config.EnableMatchAreaId != 0 {
		area := u.Yggdrasil.Areas.getByCreate(ctx, u.Yggdrasil, config.EnableMatchAreaId)
		area.IsTaskDone = true
	}
	return nil

}
func (y *YggdrasilTask) ChooseNext(ctx context.Context, u *User, nextSubTaskId int32) (*YggdrasilTaskInfo, error) {
	nextSubTaskConfig, err := manager.CSV.Yggdrasil.GetSubTaskConfig(nextSubTaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	info, ok := y.TaskInProgress[nextSubTaskConfig.TaskId]
	if !ok {
		return nil, errors.WrapTrace(common.ErrYggdrasilTaskNotInProgress)
	}
	nowSubTaskConfig, err := manager.CSV.Yggdrasil.GetSubTaskConfig(info.Base.SubTaskId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	env := nowSubTaskConfig.Envs.Env(info.Base.EnvId)
	handler, err := TaskHandlers.getHandler(env.SubTaskType)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 子任务未完成
	if !handler.IsComplete(info.Base) {
		return nil, errors.WrapTrace(common.ErrYggdrasilTaskCannotComplete)
	}
	if !nowSubTaskConfig.NextSubTaskIds.Contains(nextSubTaskId) {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	return y.InitSubTask(ctx, u, nextSubTaskId)
}

// AppendTemporalObject 记录obj原始state和坐标，任务放弃或者失败的时候重置成操作前的状态
func (y *YggdrasilTask) AppendTemporalObject(obj *YggObject) {
	// env 创造的obj不记录
	if obj.CreateAtSubTaskId > 0 {
		return
	}

	state, err := manager.CSV.Yggdrasil.GetObjectState(obj.ObjectId, obj.State)
	if err != nil {
		glog.Errorf("AppendTemporalObject GetObjectState err:%+v", err)
		return
	}
	if state.SubTaskID > 0 {

		_, info, ok := y.YggSubTaskProgressInfoInProcess(state.SubTaskID)
		if !ok {
			return
		}
		info.AppendTemporalObject(obj)
	}
}

func (y *YggdrasilTask) DeliverTaskGoods(ctx context.Context, u *User, subTaskId int32, resources []*pb.VOResource) error {
	progress, info, ok := y.YggSubTaskProgressInfoInProcess(subTaskId)
	if !ok {
		return errors.WrapTrace(common.ErrYggdrasilTaskNotInProgress)
	}
	config, err := manager.CSV.Yggdrasil.GetSubTaskConfig(subTaskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	env := config.Envs.Env(progress.EnvId)

	subTaskType := env.SubTaskType
	if subTaskType != static.YggdrasilSubTaskTypeTypeDeliverItem && subTaskType != static.YggdrasilSubTaskTypeTypeDeliverItemSelectOne {
		return errors.WrapTrace(common.ErrParamError)
	}
	var rewards []*common.Reward
	for _, data := range *env.SubTaskTargets {
		rewards = append(rewards, common.NewReward(data.FilterConditions[0], data.FilterConditions[1]))
	}

	if subTaskType == static.YggdrasilSubTaskTypeTypeDeliverItem {
		err = y.DeliverAll(u, rewards...)
		if err != nil {
			return errors.WrapTrace(err)
		}
	} else if subTaskType == static.YggdrasilSubTaskTypeTypeDeliverItemSelectOne {
		if len(resources) != 1 {
			return errors.WrapTrace(common.ErrParamError)
		}
		err = y.DeliverSelectOne(u, rewards, resources[0])
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	err = y.ForceCompleteSubTask(ctx, u, progress, info)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

// YggSubTaskProgressInfoInProcess 根据subTaskId 获得任务进度
func (y *YggdrasilTask) YggSubTaskProgressInfoInProcess(subTaskId int32) (*YggSubTaskProgressInfo, *YggdrasilTaskInfo, bool) {
	for _, info := range y.TaskInProgress {
		if info.Base.SubTaskId == subTaskId {
			return info.Base, info, true
		}
		for _, multi := range info.Multi {
			if multi.SubTaskId == subTaskId {
				return multi, info, true
			}
		}
	}
	return nil, nil, false
}

// DeliverSelectOne 交付n选1任务道具
func (y *YggdrasilTask) DeliverSelectOne(u *User, rewards []*common.Reward, resource *pb.VOResource) error {
	var success bool
	for _, reward := range rewards {
		if reward.ID == resource.ItemId && reward.Num == resource.Count {
			success = true
		}
	}
	if !success {
		return errors.WrapTrace(common.ErrParamError)
	}
	return y.DeliverAll(u, common.NewReward(resource.ItemId, resource.Count))
}

// DeliverAll 交付所有任务道具
func (y *YggdrasilTask) DeliverAll(u *User, rewards ...*common.Reward) error {
	newRewards := common.NewRewards()
	for _, reward := range rewards {
		newRewards.AddReward(reward)
	}
	err := u.Yggdrasil.CheckRewardsEnough(u, newRewards)
	if err != nil {
		return errors.WrapTrace(err)
	}

	reason := logreason.NewReason(logreason.YggDeliverTaskItem)
	u.Yggdrasil.CostRewards(u, newRewards, reason)

	return nil
}

// ForceCompleteSubTask 不走任务监听完成子任务
func (y *YggdrasilTask) ForceCompleteSubTask(ctx context.Context, u *User, progressInfo *YggSubTaskProgressInfo, taskInfo *YggdrasilTaskInfo) error {

	subTaskConfig, err := manager.CSV.Yggdrasil.GetSubTaskConfig(progressInfo.SubTaskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	env := subTaskConfig.Envs.Env(progressInfo.EnvId)
	handler, err := TaskHandlers.getHandler(env.SubTaskType)
	if err != nil {
		return errors.WrapTrace(err)
	}
	err = handler.completeProcess(progressInfo)
	if err != nil {
		return errors.WrapTrace(err)
	}

	progressInfo.IsSubTaskCompleted = true
	//  子任务完成时 1.添加到已完成子任务列表 2.env暂时删除的entity恢复 3.env创造的obj删除 4. 子任务发奖
	err = y.WhenSubTaskComplete(ctx, u, taskInfo, progressInfo)
	if err != nil {
		return errors.WrapTrace(err)
	}

	// 子任务的子任务完成后，流转子任务
	y.process_(ctx, u, taskInfo, taskInfo.Base, static.YggdrasilSubTaskTypeTypeMultiSubTask, []int32{}, progressInfo.SubTaskId)
	u.AddYggPush(&pb.S2CYggdrasilTaskInfoChange{
		TaskInfo:       taskInfo.VOYggdrasilTaskInfo(),
		EntityChange:   u.VOYggdrasilEntityChange(ctx),
		ResourceResult: u.VOResourceResult(),
	})
	return nil
}

func (y *YggSubTaskProgressInfo) filter(subTaskTargets *common.YggdrasilSubTaskTargets, filterCondition []int32) (int, bool) {
	index, ok := subTaskTargets.CanProcess(filterCondition)
	return index, ok

}

func (y *YggdrasilTaskInfo) AppendTemporalObject(obj *YggObject) {
	_, ok := y.TemporalObjects[obj.Uid]
	//如果已经有备份了 return
	if ok {
		return
	}
	// 备份
	y.TemporalObjects[obj.Uid] = obj.Clone()
}
