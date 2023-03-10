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

type CfgWorldItemAdvance struct {
	Id         int32
	Stage      int32
	Star       []int32
	LevelLimit []int32
	GoldCost   []int32
}

type CfgWorldItemAdvanceConfig struct {
	data map[int32]*CfgWorldItemAdvance
}

func NewCfgWorldItemAdvanceConfig() *CfgWorldItemAdvanceConfig {
	return &CfgWorldItemAdvanceConfig{
		data: make(map[int32]*CfgWorldItemAdvance),
	}
}

func (c *CfgWorldItemAdvanceConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgWorldItemAdvance)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgWorldItemAdvance.Id field error,value:", vId)
			return false
		}

		/* parse Stage field */
		vStage, _ := parse.GetFieldByName(uint32(i), "stage")
		var StageRet bool
		data.Stage, StageRet = String2Int32(vStage)
		if !StageRet {
			glog.Error("Parse CfgWorldItemAdvance.Stage field error,value:", vStage)
			return false
		}

		/* parse Star field */
		vecStar, _ := parse.GetFieldByName(uint32(i), "star")
		if vecStar != "" {
			arrayStar := strings.Split(vecStar, ",")
			for j := 0; j < len(arrayStar); j++ {
				v, ret := String2Int32(arrayStar[j])
				if !ret {
					glog.Error("Parse CfgWorldItemAdvance.Star field error, value:", arrayStar[j])
					return false
				}
				data.Star = append(data.Star, v)
			}
		}

		/* parse LevelLimit field */
		vecLevelLimit, _ := parse.GetFieldByName(uint32(i), "levelLimit")
		if vecLevelLimit != "" {
			arrayLevelLimit := strings.Split(vecLevelLimit, ",")
			for j := 0; j < len(arrayLevelLimit); j++ {
				v, ret := String2Int32(arrayLevelLimit[j])
				if !ret {
					glog.Error("Parse CfgWorldItemAdvance.LevelLimit field error, value:", arrayLevelLimit[j])
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
					glog.Error("Parse CfgWorldItemAdvance.GoldCost field error, value:", arrayGoldCost[j])
					return false
				}
				data.GoldCost = append(data.GoldCost, v)
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

func (c *CfgWorldItemAdvanceConfig) Clear() {
}

func (c *CfgWorldItemAdvanceConfig) Find(id int32) (*CfgWorldItemAdvance, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgWorldItemAdvanceConfig) GetAllData() map[int32]*CfgWorldItemAdvance {
	return c.data
}

func (c *CfgWorldItemAdvanceConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Stage, ",", v.Star, ",", v.LevelLimit, ",", v.GoldCost)
	}
}
