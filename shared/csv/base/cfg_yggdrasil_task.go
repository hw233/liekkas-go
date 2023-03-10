// ===================================== //
// author:       gavingqf                //
// == Please don'g change me by hand ==  //
//====================================== //

/*you have defined the following interface:
type IConfig interface {
	// load interface
	Load(path string) bool

	// clear interface
	Clear()
}
*/

package base

import (
	"shared/utility/glog"
	"strings"
)

type CfgYggdrasilTask struct {
	Id                int32
	TaskType          int32
	FinishTaskCity    int32
	UnlockCondition   []string
	TaskGroup         int32
	NextSubTaskId     int32
	DropId            int32
	EnableMatchAreaId int32
	AbandonTask       bool
}

type CfgYggdrasilTaskConfig struct {
	data map[int32]*CfgYggdrasilTask
}

func NewCfgYggdrasilTaskConfig() *CfgYggdrasilTaskConfig {
	return &CfgYggdrasilTaskConfig{
		data: make(map[int32]*CfgYggdrasilTask),
	}
}

func (c *CfgYggdrasilTaskConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilTask)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilTask.Id field error,value:", vId)
			return false
		}

		/* parse TaskType field */
		vTaskType, _ := parse.GetFieldByName(uint32(i), "taskType")
		var TaskTypeRet bool
		data.TaskType, TaskTypeRet = String2Int32(vTaskType)
		if !TaskTypeRet {
			glog.Error("Parse CfgYggdrasilTask.TaskType field error,value:", vTaskType)
			return false
		}

		/* parse FinishTaskCity field */
		vFinishTaskCity, _ := parse.GetFieldByName(uint32(i), "finishTaskCity")
		var FinishTaskCityRet bool
		data.FinishTaskCity, FinishTaskCityRet = String2Int32(vFinishTaskCity)
		if !FinishTaskCityRet {
			glog.Error("Parse CfgYggdrasilTask.FinishTaskCity field error,value:", vFinishTaskCity)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse TaskGroup field */
		vTaskGroup, _ := parse.GetFieldByName(uint32(i), "taskGroup")
		var TaskGroupRet bool
		data.TaskGroup, TaskGroupRet = String2Int32(vTaskGroup)
		if !TaskGroupRet {
			glog.Error("Parse CfgYggdrasilTask.TaskGroup field error,value:", vTaskGroup)
			return false
		}

		/* parse NextSubTaskId field */
		vNextSubTaskId, _ := parse.GetFieldByName(uint32(i), "nextSubTaskId")
		var NextSubTaskIdRet bool
		data.NextSubTaskId, NextSubTaskIdRet = String2Int32(vNextSubTaskId)
		if !NextSubTaskIdRet {
			glog.Error("Parse CfgYggdrasilTask.NextSubTaskId field error,value:", vNextSubTaskId)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgYggdrasilTask.DropId field error,value:", vDropId)
			return false
		}

		/* parse EnableMatchAreaId field */
		vEnableMatchAreaId, _ := parse.GetFieldByName(uint32(i), "enableMatchAreaId")
		var EnableMatchAreaIdRet bool
		data.EnableMatchAreaId, EnableMatchAreaIdRet = String2Int32(vEnableMatchAreaId)
		if !EnableMatchAreaIdRet {
			glog.Error("Parse CfgYggdrasilTask.EnableMatchAreaId field error,value:", vEnableMatchAreaId)
			return false
		}

		/* parse AbandonTask field */
		vAbandonTask, _ := parse.GetFieldByName(uint32(i), "abandonTask")
		var AbandonTaskRet bool
		data.AbandonTask, AbandonTaskRet = String2Bool(vAbandonTask)
		if !AbandonTaskRet {
			glog.Error("Parse CfgYggdrasilTask.AbandonTask field error,value:", vAbandonTask)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgYggdrasilTaskConfig) Clear() {
}

func (c *CfgYggdrasilTaskConfig) Find(id int32) (*CfgYggdrasilTask, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilTaskConfig) GetAllData() map[int32]*CfgYggdrasilTask {
	return c.data
}

func (c *CfgYggdrasilTaskConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.TaskType, ",", v.FinishTaskCity, ",", v.UnlockCondition, ",", v.TaskGroup, ",", v.NextSubTaskId, ",", v.DropId, ",", v.EnableMatchAreaId, ",", v.AbandonTask)
	}
}
