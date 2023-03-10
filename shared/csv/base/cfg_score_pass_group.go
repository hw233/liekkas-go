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
)

type CfgScorePassGroup struct {
	Id             int32
	PhaseId        int32
	StartDayOffset int32
	TaskGroupId    int32
}

type CfgScorePassGroupConfig struct {
	data map[int32]*CfgScorePassGroup
}

func NewCfgScorePassGroupConfig() *CfgScorePassGroupConfig {
	return &CfgScorePassGroupConfig{
		data: make(map[int32]*CfgScorePassGroup),
	}
}

func (c *CfgScorePassGroupConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgScorePassGroup)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgScorePassGroup.Id field error,value:", vId)
			return false
		}

		/* parse PhaseId field */
		vPhaseId, _ := parse.GetFieldByName(uint32(i), "phaseId")
		var PhaseIdRet bool
		data.PhaseId, PhaseIdRet = String2Int32(vPhaseId)
		if !PhaseIdRet {
			glog.Error("Parse CfgScorePassGroup.PhaseId field error,value:", vPhaseId)
			return false
		}

		/* parse StartDayOffset field */
		vStartDayOffset, _ := parse.GetFieldByName(uint32(i), "startDayOffset")
		var StartDayOffsetRet bool
		data.StartDayOffset, StartDayOffsetRet = String2Int32(vStartDayOffset)
		if !StartDayOffsetRet {
			glog.Error("Parse CfgScorePassGroup.StartDayOffset field error,value:", vStartDayOffset)
			return false
		}

		/* parse TaskGroupId field */
		vTaskGroupId, _ := parse.GetFieldByName(uint32(i), "taskGroupId")
		var TaskGroupIdRet bool
		data.TaskGroupId, TaskGroupIdRet = String2Int32(vTaskGroupId)
		if !TaskGroupIdRet {
			glog.Error("Parse CfgScorePassGroup.TaskGroupId field error,value:", vTaskGroupId)
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

func (c *CfgScorePassGroupConfig) Clear() {
}

func (c *CfgScorePassGroupConfig) Find(id int32) (*CfgScorePassGroup, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgScorePassGroupConfig) GetAllData() map[int32]*CfgScorePassGroup {
	return c.data
}

func (c *CfgScorePassGroupConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.PhaseId, ",", v.StartDayOffset, ",", v.TaskGroupId)
	}
}
