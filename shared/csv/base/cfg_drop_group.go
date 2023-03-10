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

type CfgDropGroup struct {
	Id         int32
	DropGroup  int32
	DropItem   int32
	DropPr     int32
	DropNumber []int32
}

type CfgDropGroupConfig struct {
	data map[int32]*CfgDropGroup
}

func NewCfgDropGroupConfig() *CfgDropGroupConfig {
	return &CfgDropGroupConfig{
		data: make(map[int32]*CfgDropGroup),
	}
}

func (c *CfgDropGroupConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgDropGroup)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgDropGroup.Id field error,value:", vId)
			return false
		}

		/* parse DropGroup field */
		vDropGroup, _ := parse.GetFieldByName(uint32(i), "dropGroup")
		var DropGroupRet bool
		data.DropGroup, DropGroupRet = String2Int32(vDropGroup)
		if !DropGroupRet {
			glog.Error("Parse CfgDropGroup.DropGroup field error,value:", vDropGroup)
			return false
		}

		/* parse DropItem field */
		vDropItem, _ := parse.GetFieldByName(uint32(i), "dropItem")
		var DropItemRet bool
		data.DropItem, DropItemRet = String2Int32(vDropItem)
		if !DropItemRet {
			glog.Error("Parse CfgDropGroup.DropItem field error,value:", vDropItem)
			return false
		}

		/* parse DropPr field */
		vDropPr, _ := parse.GetFieldByName(uint32(i), "dropPr")
		var DropPrRet bool
		data.DropPr, DropPrRet = String2Int32(vDropPr)
		if !DropPrRet {
			glog.Error("Parse CfgDropGroup.DropPr field error,value:", vDropPr)
			return false
		}

		/* parse DropNumber field */
		vecDropNumber, _ := parse.GetFieldByName(uint32(i), "dropNumber")
		if vecDropNumber != "" {
			arrayDropNumber := strings.Split(vecDropNumber, ",")
			for j := 0; j < len(arrayDropNumber); j++ {
				v, ret := String2Int32(arrayDropNumber[j])
				if !ret {
					glog.Error("Parse CfgDropGroup.DropNumber field error, value:", arrayDropNumber[j])
					return false
				}
				data.DropNumber = append(data.DropNumber, v)
			}
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgDropGroupConfig) Clear() {
}

func (c *CfgDropGroupConfig) Find(id int32) (*CfgDropGroup, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgDropGroupConfig) GetAllData() map[int32]*CfgDropGroup {
	return c.data
}

func (c *CfgDropGroupConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.DropGroup, ",", v.DropItem, ",", v.DropPr, ",", v.DropNumber)
	}
}
