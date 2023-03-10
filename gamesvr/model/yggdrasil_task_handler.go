package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/utility/errors"
	"shared/utility/glog"
)

const Done = 1
const UnDone = 0

var TaskHandlers *YggdrasilTaskHandlers
var TaskFailHandler *YggdrasilTaskFailHandlers

func init() {
	TaskHandlers = NewYggdrasilTaskHandlers()
	TaskHandlers.setHandler(NewTaskHandler(&TaskHandlerSimple{}), static.YggdrasilSubTaskTypeTypeVn, static.YggdrasilSubTaskTypeTypeMultiTask, static.YggdrasilSubTaskTypeTypeMultiSubTask)
	TaskHandlers.setHandler(NewTaskHandler(&TaskHandlerCoordinate{}), static.YggdrasilSubTaskTypeTypeLeadWay, static.YggdrasilSubTaskTypeTypeMove)
	TaskHandlers.setHandler(NewTaskHandler(&TaskHandlerConvey{}), static.YggdrasilSubTaskTypeTypeConvoy)
	TaskHandlers.setHandler(NewTaskHandler(&TaskHandlerCumulative{}), static.YggdrasilSubTaskTypeTypeChapter, static.YggdrasilSubTaskTypeTypeMonster, static.YggdrasilSubTaskTypeTypeHelpBuild, static.YggdrasilSubTaskTypeTypeBuild)
	TaskHandlers.setHandler(NewTaskHandler(&TaskHandlerReplaced{}), static.YggdrasilSubTaskTypeTypeObjectStateChange, static.YggdrasilSubTaskTypeTypeCity)
	TaskHandlers.setHandler(NewTaskHandler(&TaskHandlerOwn{}), static.YggdrasilSubTaskTypeTypeOwn)
	TaskHandlers.setHandler(NewTaskHandler(&TaskHandlerDeliver{}), static.YggdrasilSubTaskTypeTypeDeliverItem, static.YggdrasilSubTaskTypeTypeDeliverItemSelectOne)

	TaskFailHandler = NewYggdrasilTaskFailHandler()
	TaskFailHandler.setHandler(&TaskHandlerConvey{}, static.YggdrasilSubTaskTypeTypeConvoy)

}

type YggdrasilTaskHandlers map[int32]*YggdrasilTaskHandler

func NewYggdrasilTaskHandlers() *YggdrasilTaskHandlers {
	return (*YggdrasilTaskHandlers)(&map[int32]*YggdrasilTaskHandler{})
}
func (y *YggdrasilTaskHandlers) setHandler(handler *YggdrasilTaskHandler, subTaskTypes ...int32) {
	for _, subTaskType := range subTaskTypes {
		(*y)[subTaskType] = handler
	}
}

func (y *YggdrasilTaskHandlers) getHandler(subTaskType int32) (*YggdrasilTaskHandler, error) {
	handler, ok := (*y)[subTaskType]
	if !ok {
		return nil, errors.WrapTrace(common.ErrYggdrasilTaskProcessError)
	}
	return handler, nil
}

type YggdrasilTaskBase interface {
	initProcess(env *entry.Env, yggdrasil *Yggdrasil) ([]*common.YggdrasilSubTaskProcess, error)
	process(org *common.YggdrasilSubTaskProcess, changeProcess []int32)
	completeProcess(env *entry.Env, info *YggSubTaskProgressInfo) error
	isComplete(env *entry.Env, info *YggSubTaskProgressInfo) bool
}

type YggdrasilTaskHandler struct {
	YggdrasilTaskBase
}

func NewTaskHandler(h YggdrasilTaskBase) *YggdrasilTaskHandler {
	return &YggdrasilTaskHandler{
		YggdrasilTaskBase: h,
	}
}
func (t *YggdrasilTaskHandler) process(info *YggSubTaskProgressInfo, index int, env *entry.Env, process []int32, filterCondition []int32) {
	// 过滤条件
	index, ok := info.filter(env.SubTaskTargets, filterCondition)
	if !ok {
		return
	}
	t.YggdrasilTaskBase.process(info.Processes[index], process)
	if t.isComplete(env, info) {
		info.IsSubTaskCompleted = true
	}
}

func (t *YggdrasilTaskHandler) init(ctx context.Context, env *entry.Env, u *User) (*YggSubTaskProgressInfoWrap, error) {
	process, err := t.initProcess(env, u.Yggdrasil)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	info := NewYggSubTaskProgressInfo()
	info.SubTaskId = env.SubTaskId
	info.Processes = process
	info.IsSubTaskCompleted = false
	info.EnvId = env.Id

	var multi []*YggSubTaskProgressInfo
	if env.SubTaskType == static.YggdrasilSubTaskTypeTypeMultiSubTask {
		multi = make([]*YggSubTaskProgressInfo, 0, len(*env.SubTaskTargets))
		for _, target := range *env.SubTaskTargets {
			subtaskId := target.FilterConditions[0]
			initProcess, err := u.Yggdrasil.Task.GetTaskInitProcess(ctx, u, subtaskId)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			multi = append(multi, initProcess.Base)
		}
	}
	return NewYggSubTaskProgressInfoWrap(info, multi), nil

}

func (t *YggdrasilTaskHandler) completeProcess(info *YggSubTaskProgressInfo) error {
	config, err := manager.CSV.Yggdrasil.GetSubTaskConfig(info.SubTaskId)
	if err != nil {
		return errors.WrapTrace(err)
	}
	env := config.Envs.Env(info.EnvId)
	return t.YggdrasilTaskBase.completeProcess(env, info)
}

func (t *YggdrasilTaskHandler) IsComplete(info *YggSubTaskProgressInfo) bool {
	config, err := manager.CSV.Yggdrasil.GetSubTaskConfig(info.SubTaskId)
	if err != nil {
		glog.Errorf("IsComplete err:%v", err)
		return false
	}
	env := config.Envs.Env(info.EnvId)
	return t.YggdrasilTaskBase.isComplete(env, info)
}

// TaskHandlerSimple done或者Undone类
type TaskHandlerSimple struct {
}

func (t *TaskHandlerSimple) initProcess(env *entry.Env, yggdrasil *Yggdrasil) ([]*common.YggdrasilSubTaskProcess, error) {
	var processes []*common.YggdrasilSubTaskProcess

	switch env.SubTaskType {
	case static.YggdrasilSubTaskTypeTypeMultiTask:
		for _, target := range *env.SubTaskTargets {
			taskId := target.FilterConditions[0]
			if yggdrasil.Task.CompleteTaskIds.Contains(taskId) {
				processes = append(processes, common.NewYggdrasilSubTaskProcess(Done))
			} else {
				processes = append(processes, common.NewYggdrasilSubTaskProcess(UnDone))
			}
		}

	default:
		for range *env.SubTaskTargets {
			processes = append(processes, common.NewYggdrasilSubTaskProcess(UnDone))
		}
	}
	return processes, nil
}

func (t *TaskHandlerSimple) process(org *common.YggdrasilSubTaskProcess, changeProcess []int32) {
	org.SetProcess(Done)

}
func (t *TaskHandlerSimple) completeProcess(env *entry.Env, info *YggSubTaskProgressInfo) error {
	for _, process := range info.Processes {
		process.SetProcess(Done)
	}
	return nil
}
func (t *TaskHandlerSimple) isComplete(env *entry.Env, info *YggSubTaskProgressInfo) bool {
	for _, process := range info.Processes {
		if process.Process[0] == UnDone {
			return false
		}
	}
	return true
}

// TaskHandlerCoordinate 坐标类
type TaskHandlerCoordinate struct {
}

func (t *TaskHandlerCoordinate) initProcess(env *entry.Env, yggdrasil *Yggdrasil) ([]*common.YggdrasilSubTaskProcess, error) {
	var processes []*common.YggdrasilSubTaskProcess
	for _, data := range *env.SubTaskTargets {

		switch env.SubTaskType {
		case static.YggdrasilSubTaskTypeTypeMove:
			processes = append(processes, common.NewYggdrasilSubTaskProcess(yggdrasil.TravelPos.X, yggdrasil.TravelPos.Y))

		case static.YggdrasilSubTaskTypeTypeLeadWay:
			objects, err := yggdrasil.FindObjectById(data.FilterConditions[0])
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			if len(objects) != 1 {
				return nil, errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilObjectRepeated, data.FilterConditions[0]))

			}
			processes = append(processes, common.NewYggdrasilSubTaskProcess(objects[0].X, objects[0].Y))
		}
	}
	return processes, nil
}
func (t *TaskHandlerCoordinate) process(org *common.YggdrasilSubTaskProcess, changeProcess []int32) {
	if len(changeProcess) != 2 {
		glog.Errorf("Task process len error :%+v", errors.WrapTrace(common.ErrYggdrasilTaskProcessError))
		return
	}
	org.SetPosProcess(changeProcess[0], changeProcess[1])
}
func (t *TaskHandlerCoordinate) completeProcess(env *entry.Env, info *YggSubTaskProgressInfo) error {
	for i, process := range info.Processes {
		process.SetPosProcess((*env.SubTaskTargets)[i].Process[0], (*env.SubTaskTargets)[i].Process[1])
	}
	return nil
}
func (t *TaskHandlerCoordinate) isComplete(env *entry.Env, info *YggSubTaskProgressInfo) bool {
	for i, data := range *env.SubTaskTargets {
		for j, process := range data.Process {
			if process != info.Processes[i].Process[j] {
				return false
			}
		}
	}
	return true
}

// TaskHandlerConvey 护送npc //todo: 血量怎么扣？
type TaskHandlerConvey struct {
	*TaskHandlerCoordinate
}

func (t *TaskHandlerConvey) initProcess(env *entry.Env, yggdrasil *Yggdrasil) ([]*common.YggdrasilSubTaskProcess, error) {
	var processes []*common.YggdrasilSubTaskProcess
	for _, data := range *env.SubTaskTargets {
		objects, err := yggdrasil.FindObjectById(data.FilterConditions[0])
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		if len(objects) != 1 {
			return nil, errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilObjectRepeated, data.FilterConditions[0]))

		}
		process := common.NewYggdrasilSubTaskProcess(objects[0].X, objects[0].Y)
		process.AppendNpc()
		processes = append(processes, process)

	}
	return processes, nil
}

func (t *TaskHandlerConvey) processFail(attach *common.YggdrasilSubTaskAttach) bool {
	return attach.NpcHp.IsDead()
}
func (t *TaskHandlerConvey) isComplete(env *entry.Env, info *YggSubTaskProgressInfo) bool {
	for i, data := range *env.SubTaskTargets {
		for j, process := range data.Process {
			if process != info.Processes[i].Process[j] {
				return false
			}
		}
	}
	return true
}

// TaskHandlerCumulative 累加类
type TaskHandlerCumulative struct {
}

func (t *TaskHandlerCumulative) initProcess(env *entry.Env, yggdrasil *Yggdrasil) ([]*common.YggdrasilSubTaskProcess, error) {
	var processes []*common.YggdrasilSubTaskProcess
	for range *env.SubTaskTargets {
		processes = append(processes, common.NewYggdrasilSubTaskProcess(UnDone))
	}
	return processes, nil
}

func (t *TaskHandlerCumulative) process(org *common.YggdrasilSubTaskProcess, changeProcess []int32) {
	if len(changeProcess) != 1 {
		glog.Errorf("Task process len error :%+v", errors.WrapTrace(common.ErrYggdrasilTaskProcessError))
		return
	}
	org.AddProcess(changeProcess[0])
}

func (t *TaskHandlerCumulative) completeProcess(env *entry.Env, info *YggSubTaskProgressInfo) error {

	for i, process := range info.Processes {
		process.SetProcess((*env.SubTaskTargets)[i].Process[0])
	}
	return nil
}
func (t *TaskHandlerCumulative) isComplete(env *entry.Env, info *YggSubTaskProgressInfo) bool {
	for i, data := range *env.SubTaskTargets {
		for j, process := range data.Process {
			if process > info.Processes[i].Process[j] {
				return false
			}
		}
	}
	return true
}

// TaskHandlerReplaced 替换类
type TaskHandlerReplaced struct {
}

func (t *TaskHandlerReplaced) initProcess(env *entry.Env, yggdrasil *Yggdrasil) ([]*common.YggdrasilSubTaskProcess, error) {
	var processes []*common.YggdrasilSubTaskProcess
	for _, data := range *env.SubTaskTargets {
		switch env.SubTaskType {
		case static.YggdrasilSubTaskTypeTypeObjectStateChange:
			objects, err := yggdrasil.FindObjectById(data.FilterConditions[0])
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			if len(objects) != 1 {
				return nil, errors.WrapTrace(errors.Swrapf(common.ErrYggdrasilObjectRepeated, data.FilterConditions[0]))

			}
			processes = append(processes, common.NewYggdrasilSubTaskProcess(objects[0].State))

		case static.YggdrasilSubTaskTypeTypeCity:

			processes = append(processes, common.NewYggdrasilSubTaskProcess(yggdrasil.CityId))

		}
	}
	return processes, nil
}
func (t *TaskHandlerReplaced) process(org *common.YggdrasilSubTaskProcess, changeProcess []int32) {
	if len(changeProcess) != 1 {
		glog.Errorf("Task process len error :%+v", errors.WrapTrace(common.ErrYggdrasilTaskProcessError))
		return
	}
	org.SetProcess(changeProcess[0])
}
func (t *TaskHandlerReplaced) completeProcess(env *entry.Env, info *YggSubTaskProgressInfo) error {

	for i, process := range info.Processes {
		process.SetProcess((*env.SubTaskTargets)[i].Process[0])
	}
	return nil
}
func (t *TaskHandlerReplaced) isComplete(env *entry.Env, info *YggSubTaskProgressInfo) bool {
	for i, data := range *env.SubTaskTargets {
		for j, process := range data.Process {
			if process != info.Processes[i].Process[j] {
				return false
			}
		}
	}
	return true
}

type TaskHandlerOwn struct {
}

func (t *TaskHandlerOwn) initProcess(env *entry.Env, yggdrasil *Yggdrasil) ([]*common.YggdrasilSubTaskProcess, error) {
	var processes []*common.YggdrasilSubTaskProcess
	for range *env.SubTaskTargets {
		// todo: 初始化拥有的物品
		processes = append(processes, common.NewYggdrasilSubTaskProcess(0))
	}
	return processes, nil
}
func (t *TaskHandlerOwn) process(org *common.YggdrasilSubTaskProcess, changeProcess []int32) {
	if len(changeProcess) != 1 {

		glog.Errorf("Task process len error :%+v", errors.WrapTrace(common.ErrYggdrasilTaskProcessError))
		return
	}

	org.SetProcess(changeProcess[0])
}

func (t *TaskHandlerOwn) completeProcess(env *entry.Env, info *YggSubTaskProgressInfo) error {

	for i, process := range info.Processes {
		process.SetProcess((*env.SubTaskTargets)[i].Process[0])
	}
	return nil
}
func (t *TaskHandlerOwn) isComplete(env *entry.Env, info *YggSubTaskProgressInfo) bool {
	for i, data := range *env.SubTaskTargets {
		for j, process := range data.Process {
			if process > info.Processes[i].Process[j] {
				return false
			}
		}
	}
	return true
}

type TaskHandlerDeliver struct {
}

func (t *TaskHandlerDeliver) initProcess(env *entry.Env, yggdrasil *Yggdrasil) ([]*common.YggdrasilSubTaskProcess, error) {
	var processes []*common.YggdrasilSubTaskProcess
	for range *env.SubTaskTargets {

		processes = append(processes, common.NewYggdrasilSubTaskProcess(0))
	}
	return processes, nil
}
func (t *TaskHandlerDeliver) process(org *common.YggdrasilSubTaskProcess, changeProcess []int32) {
	if len(changeProcess) != 1 {

		glog.Errorf("Task process len error :%+v", errors.WrapTrace(common.ErrYggdrasilTaskProcessError))
		return
	}

	org.SetProcess(changeProcess[0])
}
func (t *TaskHandlerDeliver) completeProcess(env *entry.Env, info *YggSubTaskProgressInfo) error {
	for _, process := range info.Processes {
		process.SetProcess(Done)
	}
	return nil
}
func (t *TaskHandlerDeliver) isComplete(env *entry.Env, info *YggSubTaskProgressInfo) bool {
	for i, data := range *env.SubTaskTargets {
		for j, process := range data.Process {
			if process != info.Processes[i].Process[j] {
				return false
			}
		}
	}
	return true
}

type YggdrasilTaskFailHandlers map[int32]YggdrasilTaskFailHandler

func NewYggdrasilTaskFailHandler() *YggdrasilTaskFailHandlers {
	return (*YggdrasilTaskFailHandlers)(&map[int32]YggdrasilTaskFailHandler{})
}
func (y *YggdrasilTaskFailHandlers) setHandler(handler YggdrasilTaskFailHandler, subTaskTypes ...int32) {
	for _, subTaskType := range subTaskTypes {
		(*y)[subTaskType] = handler
	}
}

func (y *YggdrasilTaskFailHandlers) getHandler(subTaskType int32) (YggdrasilTaskFailHandler, error) {
	handler, ok := (*y)[subTaskType]
	if !ok {
		return nil, errors.WrapTrace(common.ErrYggdrasilTaskProcessError)
	}
	return handler, nil
}

type YggdrasilTaskFailHandler interface {
	processFail(attach *common.YggdrasilSubTaskAttach) bool
}
