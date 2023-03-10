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

type CfgManualGroup struct {
	Id int32
}

type CfgManualGroupConfig struct {
	data map[int32]*CfgManualGroup
}

func NewCfgManualGroupConfig() *CfgManualGroupConfig {
	return &CfgManualGroupConfig{
		data: make(map[int32]*CfgManualGroup),
	}
}

func (c *CfgManualGroupConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgManualGroup)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgManualGroup.Id field error,value:", vId)
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

func (c *CfgManualGroupConfig) Clear() {
}

func (c *CfgManualGroupConfig) Find(id int32) (*CfgManualGroup, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgManualGroupConfig) GetAllData() map[int32]*CfgManualGroup {
	return c.data
}

func (c *CfgManualGroupConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id)
	}
}
