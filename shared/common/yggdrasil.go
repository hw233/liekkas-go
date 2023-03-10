package common

import (
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/coordinate"
	"shared/utility/number"
)

type YggdrasilNpcHp struct {
	*number.CalNumber
}

func NewYggdrasilNpcHp(init int32) *YggdrasilNpcHp {
	return &YggdrasilNpcHp{
		CalNumber: number.NewCalNumber(init),
	}
}
func (y *YggdrasilNpcHp) IsDead() bool {
	return y.Value() <= 0
}

// YggdrasilSubTaskAttach 任务进度附带的其他进度
type YggdrasilSubTaskAttach struct {
	NpcHp *YggdrasilNpcHp `json:"npc_hp"`
}

func NewYggdrasilSubTaskProcessAttach() *YggdrasilSubTaskAttach {
	return &YggdrasilSubTaskAttach{}
}
func (y *YggdrasilSubTaskAttach) TaskFailed() bool {
	if y.NpcHp != nil {
		if y.NpcHp.IsDead() {
			return true
		}
	}
	return false
}
func (y *YggdrasilSubTaskAttach) VOYggdrasilSubTaskAttach() *pb.VOYggdrasilSubTaskAttach {
	if y == nil {
		return nil
	}
	var NpcHp int32
	if y.NpcHp != nil {
		NpcHp = y.NpcHp.Value()
	}

	return &pb.VOYggdrasilSubTaskAttach{
		NpcHp: NpcHp,
	}
}

func (y *YggdrasilSubTaskAttach) AppendNpc() {
	y.NpcHp = NewYggdrasilNpcHp(100)
}

type YggdrasilSubTaskProcess struct {
	Process []int32                 `json:"process"`
	Attach  *YggdrasilSubTaskAttach `json:"attach"`
}

func NewYggdrasilSubTaskProcess(Process ...int32) *YggdrasilSubTaskProcess {
	return &YggdrasilSubTaskProcess{
		Process: Process,
		Attach:  nil,
	}
}

func (y *YggdrasilSubTaskProcess) SetProcess(val int32) {
	y.Process[0] = val
}

func (y *YggdrasilSubTaskProcess) AddProcess(addVal int32) {
	y.Process[0] = y.Process[0] + addVal
}

func (y *YggdrasilSubTaskProcess) SetPosProcess(X, Y int32) {
	y.Process[0] = X
	y.Process[1] = Y
}

func (y *YggdrasilSubTaskProcess) AppendNpc() {
	if y.Attach == nil {
		y.Attach = NewYggdrasilSubTaskProcessAttach()
	}
	y.Attach.AppendNpc()

}

func (y *YggdrasilSubTaskProcess) MinusNpcHp(val int32) {
	if y.Attach != nil && y.Attach.NpcHp != nil {
		y.Attach.NpcHp.Minus(val)
	}

}

func (y *YggdrasilSubTaskProcess) VOYggdrasilSubTaskProcess() *pb.VOYggdrasilSubTaskProcess {
	l := len(y.Process)
	if l == 0 {
		return &pb.VOYggdrasilSubTaskProcess{
			Attach: y.Attach.VOYggdrasilSubTaskAttach(),
		}

	} else if l == 1 {
		return &pb.VOYggdrasilSubTaskProcess{
			Process1: y.Process[0],
			Attach:   y.Attach.VOYggdrasilSubTaskAttach(),
		}

	} else if l == 2 {
		return &pb.VOYggdrasilSubTaskProcess{
			Process1: y.Process[0],
			Process2: y.Process[1],
			Attach:   y.Attach.VOYggdrasilSubTaskAttach(),
		}
	}

	return nil

}

type YggdrasilSubTaskTarget struct {
	FilterConditions []int32
	Process          []int32
}

func NewYggdrasilSubTaskTarget(subTaskType int32, arr []int32) *YggdrasilSubTaskTarget {

	var FilterConditions []int32
	var Process []int32

	switch subTaskType {

	case static.YggdrasilSubTaskTypeTypeVn,
		static.YggdrasilSubTaskTypeTypeMultiTask,
		static.YggdrasilSubTaskTypeTypeMultiSubTask:
		FilterConditions = append(FilterConditions, arr[0])
	case static.YggdrasilSubTaskTypeTypeCity:
		Process = append(Process, arr[0])

	case static.YggdrasilSubTaskTypeTypeConvoy,
		static.YggdrasilSubTaskTypeTypeLeadWay:
		FilterConditions = append(FilterConditions, arr[0])
		Process = append(Process, arr[1], arr[2])
	case static.YggdrasilSubTaskTypeTypeMove:
		Process = append(Process, arr[0], arr[1])
	case static.YggdrasilSubTaskTypeTypeObjectStateChange,
		static.YggdrasilSubTaskTypeTypeMonster,
		static.YggdrasilSubTaskTypeTypeHelpBuild,
		static.YggdrasilSubTaskTypeTypeBuild,
		static.YggdrasilSubTaskTypeTypeOwn,
		static.YggdrasilSubTaskTypeTypeChapter:
		FilterConditions = append(FilterConditions, arr[0])
		Process = append(Process, arr[1])

	case static.YggdrasilSubTaskTypeTypeDeliverItem,
		static.YggdrasilSubTaskTypeTypeDeliverItemSelectOne:
		FilterConditions = append(FilterConditions, arr[0])
		FilterConditions = append(FilterConditions, arr[1])

	}

	return &YggdrasilSubTaskTarget{
		FilterConditions: FilterConditions,
		Process:          Process,
	}
}

func (y *YggdrasilSubTaskTarget) CanProcess(filterCondition []int32) bool {
	if len(filterCondition) != len(y.FilterConditions) {
		return false
	}
	for i, v := range filterCondition {
		if y.FilterConditions[i] != v {
			return false
		}
	}
	return true
}

type YggdrasilSubTaskTargets []*YggdrasilSubTaskTarget

func (y *YggdrasilSubTaskTargets) CanProcess(filterCondition []int32) (int, bool) {
	for i, data := range *y {
		if data.CanProcess(filterCondition) {
			return i, true
		}
	}
	return -1, false
}

type TaskItems []*TaskItem

func NewTaskItems() *TaskItems {
	return (*TaskItems)(&[]*TaskItem{})
}

type TaskItem struct {
	ID  int32
	Num int32
}

type EnvObjs []*EnvObj

func NewEnvObjs() *EnvObjs {
	return (*EnvObjs)(&[]*EnvObj{})
}

type EnvObj struct {
	Relative bool
	Position coordinate.Position
	ObjectId int32
	DeleteAt int32
}

type EnvTerrains struct {
	HeightTerrains []*EnvTerrainHeight
	TypeTerrains   []*EnvTerrainType
	DeleteAt       int32
}

type EnvTerrainHeight struct {
	*Area
	PosHeight int32
}
type EnvTerrainType struct {
	*Area
	PosType int32
}
