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

type CfgYggdrasilObject struct {
	Id          int32
	ObjectState int32
}

type CfgYggdrasilObjectConfig struct {
	data map[int32]*CfgYggdrasilObject
}

func NewCfgYggdrasilObjectConfig() *CfgYggdrasilObjectConfig {
	return &CfgYggdrasilObjectConfig{
		data: make(map[int32]*CfgYggdrasilObject),
	}
}

func (c *CfgYggdrasilObjectConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilObject)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilObject.Id field error,value:", vId)
			return false
		}

		/* parse ObjectState field */
		vObjectState, _ := parse.GetFieldByName(uint32(i), "objectState")
		var ObjectStateRet bool
		data.ObjectState, ObjectStateRet = String2Int32(vObjectState)
		if !ObjectStateRet {
			glog.Error("Parse CfgYggdrasilObject.ObjectState field error,value:", vObjectState)
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

func (c *CfgYggdrasilObjectConfig) Clear() {
}

func (c *CfgYggdrasilObjectConfig) Find(id int32) (*CfgYggdrasilObject, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilObjectConfig) GetAllData() map[int32]*CfgYggdrasilObject {
	return c.data
}

func (c *CfgYggdrasilObjectConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ObjectState)
	}
}
