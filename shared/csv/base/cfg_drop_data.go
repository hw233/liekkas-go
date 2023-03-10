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

type CfgDropData struct {
	Id        int32
	DropGroup int32
	DropMode  int32
	DropType  int32
	DropCount int32
}

type CfgDropDataConfig struct {
	data map[int32]*CfgDropData
}

func NewCfgDropDataConfig() *CfgDropDataConfig {
	return &CfgDropDataConfig{
		data: make(map[int32]*CfgDropData),
	}
}

func (c *CfgDropDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgDropData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgDropData.Id field error,value:", vId)
			return false
		}

		/* parse DropGroup field */
		vDropGroup, _ := parse.GetFieldByName(uint32(i), "dropGroup")
		var DropGroupRet bool
		data.DropGroup, DropGroupRet = String2Int32(vDropGroup)
		if !DropGroupRet {
			glog.Error("Parse CfgDropData.DropGroup field error,value:", vDropGroup)
			return false
		}

		/* parse DropMode field */
		vDropMode, _ := parse.GetFieldByName(uint32(i), "dropMode")
		var DropModeRet bool
		data.DropMode, DropModeRet = String2Int32(vDropMode)
		if !DropModeRet {
			glog.Error("Parse CfgDropData.DropMode field error,value:", vDropMode)
			return false
		}

		/* parse DropType field */
		vDropType, _ := parse.GetFieldByName(uint32(i), "dropType")
		var DropTypeRet bool
		data.DropType, DropTypeRet = String2Int32(vDropType)
		if !DropTypeRet {
			glog.Error("Parse CfgDropData.DropType field error,value:", vDropType)
			return false
		}

		/* parse DropCount field */
		vDropCount, _ := parse.GetFieldByName(uint32(i), "dropCount")
		var DropCountRet bool
		data.DropCount, DropCountRet = String2Int32(vDropCount)
		if !DropCountRet {
			glog.Error("Parse CfgDropData.DropCount field error,value:", vDropCount)
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

func (c *CfgDropDataConfig) Clear() {
}

func (c *CfgDropDataConfig) Find(id int32) (*CfgDropData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgDropDataConfig) GetAllData() map[int32]*CfgDropData {
	return c.data
}

func (c *CfgDropDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.DropGroup, ",", v.DropMode, ",", v.DropType, ",", v.DropCount)
	}
}
