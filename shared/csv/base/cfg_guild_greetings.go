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

type CfgGuildGreetings struct {
	Id     int32
	DropId int32
}

type CfgGuildGreetingsConfig struct {
	data map[int32]*CfgGuildGreetings
}

func NewCfgGuildGreetingsConfig() *CfgGuildGreetingsConfig {
	return &CfgGuildGreetingsConfig{
		data: make(map[int32]*CfgGuildGreetings),
	}
}

func (c *CfgGuildGreetingsConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGuildGreetings)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGuildGreetings.Id field error,value:", vId)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgGuildGreetings.DropId field error,value:", vDropId)
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

func (c *CfgGuildGreetingsConfig) Clear() {
}

func (c *CfgGuildGreetingsConfig) Find(id int32) (*CfgGuildGreetings, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGuildGreetingsConfig) GetAllData() map[int32]*CfgGuildGreetings {
	return c.data
}

func (c *CfgGuildGreetingsConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.DropId)
	}
}
