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

type CfgEverydayEnergyReceive struct {
	Id              int32
	StarReceiveTime string
	EndReceiveTime  string
	DropID          int32
}

type CfgEverydayEnergyReceiveConfig struct {
	data map[int32]*CfgEverydayEnergyReceive
}

func NewCfgEverydayEnergyReceiveConfig() *CfgEverydayEnergyReceiveConfig {
	return &CfgEverydayEnergyReceiveConfig{
		data: make(map[int32]*CfgEverydayEnergyReceive),
	}
}

func (c *CfgEverydayEnergyReceiveConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEverydayEnergyReceive)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEverydayEnergyReceive.Id field error,value:", vId)
			return false
		}

		/* parse StarReceiveTime field */
		data.StarReceiveTime, _ = parse.GetFieldByName(uint32(i), "starReceiveTime")

		/* parse EndReceiveTime field */
		data.EndReceiveTime, _ = parse.GetFieldByName(uint32(i), "endReceiveTime")

		/* parse DropID field */
		vDropID, _ := parse.GetFieldByName(uint32(i), "dropID")
		var DropIDRet bool
		data.DropID, DropIDRet = String2Int32(vDropID)
		if !DropIDRet {
			glog.Error("Parse CfgEverydayEnergyReceive.DropID field error,value:", vDropID)
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

func (c *CfgEverydayEnergyReceiveConfig) Clear() {
}

func (c *CfgEverydayEnergyReceiveConfig) Find(id int32) (*CfgEverydayEnergyReceive, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEverydayEnergyReceiveConfig) GetAllData() map[int32]*CfgEverydayEnergyReceive {
	return c.data
}

func (c *CfgEverydayEnergyReceiveConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.StarReceiveTime, ",", v.EndReceiveTime, ",", v.DropID)
	}
}
