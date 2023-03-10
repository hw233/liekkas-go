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

type CfgGachaDrop struct {
	Id     int32
	DropId int32
	Item   int32
	Count  int32
	Weight int32
}

type CfgGachaDropConfig struct {
	data map[int32]*CfgGachaDrop
}

func NewCfgGachaDropConfig() *CfgGachaDropConfig {
	return &CfgGachaDropConfig{
		data: make(map[int32]*CfgGachaDrop),
	}
}

func (c *CfgGachaDropConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGachaDrop)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGachaDrop.Id field error,value:", vId)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgGachaDrop.DropId field error,value:", vDropId)
			return false
		}

		/* parse Item field */
		vItem, _ := parse.GetFieldByName(uint32(i), "item")
		var ItemRet bool
		data.Item, ItemRet = String2Int32(vItem)
		if !ItemRet {
			glog.Error("Parse CfgGachaDrop.Item field error,value:", vItem)
			return false
		}

		/* parse Count field */
		vCount, _ := parse.GetFieldByName(uint32(i), "count")
		var CountRet bool
		data.Count, CountRet = String2Int32(vCount)
		if !CountRet {
			glog.Error("Parse CfgGachaDrop.Count field error,value:", vCount)
			return false
		}

		/* parse Weight field */
		vWeight, _ := parse.GetFieldByName(uint32(i), "weight")
		var WeightRet bool
		data.Weight, WeightRet = String2Int32(vWeight)
		if !WeightRet {
			glog.Error("Parse CfgGachaDrop.Weight field error,value:", vWeight)
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

func (c *CfgGachaDropConfig) Clear() {
}

func (c *CfgGachaDropConfig) Find(id int32) (*CfgGachaDrop, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGachaDropConfig) GetAllData() map[int32]*CfgGachaDrop {
	return c.data
}

func (c *CfgGachaDropConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.DropId, ",", v.Item, ",", v.Count, ",", v.Weight)
	}
}
