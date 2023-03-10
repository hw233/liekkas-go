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

type CfgGuildTask struct {
	Id         int32
	Count      int32
	Activation int32
	DropId     int32
	IsSeparate bool
}

type CfgGuildTaskConfig struct {
	data map[int32]*CfgGuildTask
}

func NewCfgGuildTaskConfig() *CfgGuildTaskConfig {
	return &CfgGuildTaskConfig{
		data: make(map[int32]*CfgGuildTask),
	}
}

func (c *CfgGuildTaskConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGuildTask)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGuildTask.Id field error,value:", vId)
			return false
		}

		/* parse Count field */
		vCount, _ := parse.GetFieldByName(uint32(i), "count")
		var CountRet bool
		data.Count, CountRet = String2Int32(vCount)
		if !CountRet {
			glog.Error("Parse CfgGuildTask.Count field error,value:", vCount)
			return false
		}

		/* parse Activation field */
		vActivation, _ := parse.GetFieldByName(uint32(i), "activation")
		var ActivationRet bool
		data.Activation, ActivationRet = String2Int32(vActivation)
		if !ActivationRet {
			glog.Error("Parse CfgGuildTask.Activation field error,value:", vActivation)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgGuildTask.DropId field error,value:", vDropId)
			return false
		}

		/* parse IsSeparate field */
		vIsSeparate, _ := parse.GetFieldByName(uint32(i), "isSeparate")
		var IsSeparateRet bool
		data.IsSeparate, IsSeparateRet = String2Bool(vIsSeparate)
		if !IsSeparateRet {
			glog.Error("Parse CfgGuildTask.IsSeparate field error,value:", vIsSeparate)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgGuildTaskConfig) Clear() {
}

func (c *CfgGuildTaskConfig) Find(id int32) (*CfgGuildTask, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGuildTaskConfig) GetAllData() map[int32]*CfgGuildTask {
	return c.data
}

func (c *CfgGuildTaskConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Count, ",", v.Activation, ",", v.DropId, ",", v.IsSeparate)
	}
}
