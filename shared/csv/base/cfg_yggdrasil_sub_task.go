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

type CfgYggdrasilSubTask struct {
	Id             int32
	TaskId         int32
	NextSubTaskIds []int32
	YggDropId      int32
	DropId         int32
}

type CfgYggdrasilSubTaskConfig struct {
	data map[int32]*CfgYggdrasilSubTask
}

func NewCfgYggdrasilSubTaskConfig() *CfgYggdrasilSubTaskConfig {
	return &CfgYggdrasilSubTaskConfig{
		data: make(map[int32]*CfgYggdrasilSubTask),
	}
}

func (c *CfgYggdrasilSubTaskConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilSubTask)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilSubTask.Id field error,value:", vId)
			return false
		}

		/* parse TaskId field */
		vTaskId, _ := parse.GetFieldByName(uint32(i), "taskId")
		var TaskIdRet bool
		data.TaskId, TaskIdRet = String2Int32(vTaskId)
		if !TaskIdRet {
			glog.Error("Parse CfgYggdrasilSubTask.TaskId field error,value:", vTaskId)
			return false
		}

		/* parse NextSubTaskIds field */
		vecNextSubTaskIds, _ := parse.GetFieldByName(uint32(i), "nextSubTaskIds")
		if vecNextSubTaskIds != "" {
			arrayNextSubTaskIds := strings.Split(vecNextSubTaskIds, ",")
			for j := 0; j < len(arrayNextSubTaskIds); j++ {
				v, ret := String2Int32(arrayNextSubTaskIds[j])
				if !ret {
					glog.Error("Parse CfgYggdrasilSubTask.NextSubTaskIds field error, value:", arrayNextSubTaskIds[j])
					return false
				}
				data.NextSubTaskIds = append(data.NextSubTaskIds, v)
			}
		}

		/* parse YggDropId field */
		vYggDropId, _ := parse.GetFieldByName(uint32(i), "yggDropId")
		var YggDropIdRet bool
		data.YggDropId, YggDropIdRet = String2Int32(vYggDropId)
		if !YggDropIdRet {
			glog.Error("Parse CfgYggdrasilSubTask.YggDropId field error,value:", vYggDropId)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgYggdrasilSubTask.DropId field error,value:", vDropId)
			return false
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgYggdrasilSubTaskConfig) Clear() {
}

func (c *CfgYggdrasilSubTaskConfig) Find(id int32) (*CfgYggdrasilSubTask, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilSubTaskConfig) GetAllData() map[int32]*CfgYggdrasilSubTask {
	return c.data
}

func (c *CfgYggdrasilSubTaskConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.TaskId, ",", v.NextSubTaskIds, ",", v.YggDropId, ",", v.DropId)
	}
}
