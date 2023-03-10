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

type CfgEquipAdvance struct {
	Id         int32
	Stage      int32
	LevelLimit []int32
	GoldCost   []int32
	ItemCost   []int32
}

type CfgEquipAdvanceConfig struct {
	data map[int32]*CfgEquipAdvance
}

func NewCfgEquipAdvanceConfig() *CfgEquipAdvanceConfig {
	return &CfgEquipAdvanceConfig{
		data: make(map[int32]*CfgEquipAdvance),
	}
}

func (c *CfgEquipAdvanceConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEquipAdvance)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEquipAdvance.Id field error,value:", vId)
			return false
		}

		/* parse Stage field */
		vStage, _ := parse.GetFieldByName(uint32(i), "stage")
		var StageRet bool
		data.Stage, StageRet = String2Int32(vStage)
		if !StageRet {
			glog.Error("Parse CfgEquipAdvance.Stage field error,value:", vStage)
			return false
		}

		/* parse LevelLimit field */
		vecLevelLimit, _ := parse.GetFieldByName(uint32(i), "levelLimit")
		if vecLevelLimit != "" {
			arrayLevelLimit := strings.Split(vecLevelLimit, ",")
			for j := 0; j < len(arrayLevelLimit); j++ {
				v, ret := String2Int32(arrayLevelLimit[j])
				if !ret {
					glog.Error("Parse CfgEquipAdvance.LevelLimit field error, value:", arrayLevelLimit[j])
					return false
				}
				data.LevelLimit = append(data.LevelLimit, v)
			}
		}

		/* parse GoldCost field */
		vecGoldCost, _ := parse.GetFieldByName(uint32(i), "goldCost")
		if vecGoldCost != "" {
			arrayGoldCost := strings.Split(vecGoldCost, ",")
			for j := 0; j < len(arrayGoldCost); j++ {
				v, ret := String2Int32(arrayGoldCost[j])
				if !ret {
					glog.Error("Parse CfgEquipAdvance.GoldCost field error, value:", arrayGoldCost[j])
					return false
				}
				data.GoldCost = append(data.GoldCost, v)
			}
		}

		/* parse ItemCost field */
		vecItemCost, _ := parse.GetFieldByName(uint32(i), "itemCost")
		if vecItemCost != "" {
			arrayItemCost := strings.Split(vecItemCost, ",")
			for j := 0; j < len(arrayItemCost); j++ {
				v, ret := String2Int32(arrayItemCost[j])
				if !ret {
					glog.Error("Parse CfgEquipAdvance.ItemCost field error, value:", arrayItemCost[j])
					return false
				}
				data.ItemCost = append(data.ItemCost, v)
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

func (c *CfgEquipAdvanceConfig) Clear() {
}

func (c *CfgEquipAdvanceConfig) Find(id int32) (*CfgEquipAdvance, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEquipAdvanceConfig) GetAllData() map[int32]*CfgEquipAdvance {
	return c.data
}

func (c *CfgEquipAdvanceConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Stage, ",", v.LevelLimit, ",", v.GoldCost, ",", v.ItemCost)
	}
}
