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

type CfgGuildAction struct {
	Id      int32
	EquipID int32
}

type CfgGuildActionConfig struct {
	data map[int32]*CfgGuildAction
}

func NewCfgGuildActionConfig() *CfgGuildActionConfig {
	return &CfgGuildActionConfig{
		data: make(map[int32]*CfgGuildAction),
	}
}

func (c *CfgGuildActionConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGuildAction)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGuildAction.Id field error,value:", vId)
			return false
		}

		/* parse EquipID field */
		vEquipID, _ := parse.GetFieldByName(uint32(i), "equipID")
		var EquipIDRet bool
		data.EquipID, EquipIDRet = String2Int32(vEquipID)
		if !EquipIDRet {
			glog.Error("Parse CfgGuildAction.EquipID field error,value:", vEquipID)
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

func (c *CfgGuildActionConfig) Clear() {
}

func (c *CfgGuildActionConfig) Find(id int32) (*CfgGuildAction, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGuildActionConfig) GetAllData() map[int32]*CfgGuildAction {
	return c.data
}

func (c *CfgGuildActionConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.EquipID)
	}
}
