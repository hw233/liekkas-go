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

type CfgCharacterLevel struct {
	Id       int32
	RiseExp  int32
	CostItem []int32
}

type CfgCharacterLevelConfig struct {
	data map[int32]*CfgCharacterLevel
}

func NewCfgCharacterLevelConfig() *CfgCharacterLevelConfig {
	return &CfgCharacterLevelConfig{
		data: make(map[int32]*CfgCharacterLevel),
	}
}

func (c *CfgCharacterLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterLevel.Id field error,value:", vId)
			return false
		}

		/* parse RiseExp field */
		vRiseExp, _ := parse.GetFieldByName(uint32(i), "riseExp")
		var RiseExpRet bool
		data.RiseExp, RiseExpRet = String2Int32(vRiseExp)
		if !RiseExpRet {
			glog.Error("Parse CfgCharacterLevel.RiseExp field error,value:", vRiseExp)
			return false
		}

		/* parse CostItem field */
		vecCostItem, _ := parse.GetFieldByName(uint32(i), "costItem")
		if vecCostItem != "" {
			arrayCostItem := strings.Split(vecCostItem, ",")
			for j := 0; j < len(arrayCostItem); j++ {
				v, ret := String2Int32(arrayCostItem[j])
				if !ret {
					glog.Error("Parse CfgCharacterLevel.CostItem field error, value:", arrayCostItem[j])
					return false
				}
				data.CostItem = append(data.CostItem, v)
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

func (c *CfgCharacterLevelConfig) Clear() {
}

func (c *CfgCharacterLevelConfig) Find(id int32) (*CfgCharacterLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterLevelConfig) GetAllData() map[int32]*CfgCharacterLevel {
	return c.data
}

func (c *CfgCharacterLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.RiseExp, ",", v.CostItem)
	}
}
