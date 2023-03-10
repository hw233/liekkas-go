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

type CfgEquipBreak struct {
	Id             int32
	EquipID        int32
	Stage          int32
	BreakItem      []int32
	BreakLv        int32
	CostGoldCount  int32
	CostEquipCount int32
}

type CfgEquipBreakConfig struct {
	data map[int32]*CfgEquipBreak
}

func NewCfgEquipBreakConfig() *CfgEquipBreakConfig {
	return &CfgEquipBreakConfig{
		data: make(map[int32]*CfgEquipBreak),
	}
}

func (c *CfgEquipBreakConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEquipBreak)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEquipBreak.Id field error,value:", vId)
			return false
		}

		/* parse EquipID field */
		vEquipID, _ := parse.GetFieldByName(uint32(i), "equipID")
		var EquipIDRet bool
		data.EquipID, EquipIDRet = String2Int32(vEquipID)
		if !EquipIDRet {
			glog.Error("Parse CfgEquipBreak.EquipID field error,value:", vEquipID)
			return false
		}

		/* parse Stage field */
		vStage, _ := parse.GetFieldByName(uint32(i), "stage")
		var StageRet bool
		data.Stage, StageRet = String2Int32(vStage)
		if !StageRet {
			glog.Error("Parse CfgEquipBreak.Stage field error,value:", vStage)
			return false
		}

		/* parse BreakItem field */
		vecBreakItem, _ := parse.GetFieldByName(uint32(i), "breakItem")
		if vecBreakItem != "" {
			arrayBreakItem := strings.Split(vecBreakItem, ",")
			for j := 0; j < len(arrayBreakItem); j++ {
				v, ret := String2Int32(arrayBreakItem[j])
				if !ret {
					glog.Error("Parse CfgEquipBreak.BreakItem field error, value:", arrayBreakItem[j])
					return false
				}
				data.BreakItem = append(data.BreakItem, v)
			}
		}

		/* parse BreakLv field */
		vBreakLv, _ := parse.GetFieldByName(uint32(i), "breakLv")
		var BreakLvRet bool
		data.BreakLv, BreakLvRet = String2Int32(vBreakLv)
		if !BreakLvRet {
			glog.Error("Parse CfgEquipBreak.BreakLv field error,value:", vBreakLv)
			return false
		}

		/* parse CostGoldCount field */
		vCostGoldCount, _ := parse.GetFieldByName(uint32(i), "costGoldCount")
		var CostGoldCountRet bool
		data.CostGoldCount, CostGoldCountRet = String2Int32(vCostGoldCount)
		if !CostGoldCountRet {
			glog.Error("Parse CfgEquipBreak.CostGoldCount field error,value:", vCostGoldCount)
			return false
		}

		/* parse CostEquipCount field */
		vCostEquipCount, _ := parse.GetFieldByName(uint32(i), "costEquipCount")
		var CostEquipCountRet bool
		data.CostEquipCount, CostEquipCountRet = String2Int32(vCostEquipCount)
		if !CostEquipCountRet {
			glog.Error("Parse CfgEquipBreak.CostEquipCount field error,value:", vCostEquipCount)
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

func (c *CfgEquipBreakConfig) Clear() {
}

func (c *CfgEquipBreakConfig) Find(id int32) (*CfgEquipBreak, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEquipBreakConfig) GetAllData() map[int32]*CfgEquipBreak {
	return c.data
}

func (c *CfgEquipBreakConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.EquipID, ",", v.Stage, ",", v.BreakItem, ",", v.BreakLv, ",", v.CostGoldCount, ",", v.CostEquipCount)
	}
}
