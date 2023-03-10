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

type CfgEquipRandAttributesValue struct {
	Id     int32
	Range1 []int32
	Prob1  int32
	Range2 []int32
	Prob2  int32
	Range3 []int32
	Prob3  int32
}

type CfgEquipRandAttributesValueConfig struct {
	data map[int32]*CfgEquipRandAttributesValue
}

func NewCfgEquipRandAttributesValueConfig() *CfgEquipRandAttributesValueConfig {
	return &CfgEquipRandAttributesValueConfig{
		data: make(map[int32]*CfgEquipRandAttributesValue),
	}
}

func (c *CfgEquipRandAttributesValueConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEquipRandAttributesValue)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEquipRandAttributesValue.Id field error,value:", vId)
			return false
		}

		/* parse Range1 field */
		vecRange1, _ := parse.GetFieldByName(uint32(i), "range1")
		if vecRange1 != "" {
			arrayRange1 := strings.Split(vecRange1, ",")
			for j := 0; j < len(arrayRange1); j++ {
				v, ret := String2Int32(arrayRange1[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributesValue.Range1 field error, value:", arrayRange1[j])
					return false
				}
				data.Range1 = append(data.Range1, v)
			}
		}

		/* parse Prob1 field */
		vProb1, _ := parse.GetFieldByName(uint32(i), "prob1")
		var Prob1Ret bool
		data.Prob1, Prob1Ret = String2Int32(vProb1)
		if !Prob1Ret {
			glog.Error("Parse CfgEquipRandAttributesValue.Prob1 field error,value:", vProb1)
			return false
		}

		/* parse Range2 field */
		vecRange2, _ := parse.GetFieldByName(uint32(i), "range2")
		if vecRange2 != "" {
			arrayRange2 := strings.Split(vecRange2, ",")
			for j := 0; j < len(arrayRange2); j++ {
				v, ret := String2Int32(arrayRange2[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributesValue.Range2 field error, value:", arrayRange2[j])
					return false
				}
				data.Range2 = append(data.Range2, v)
			}
		}

		/* parse Prob2 field */
		vProb2, _ := parse.GetFieldByName(uint32(i), "prob2")
		var Prob2Ret bool
		data.Prob2, Prob2Ret = String2Int32(vProb2)
		if !Prob2Ret {
			glog.Error("Parse CfgEquipRandAttributesValue.Prob2 field error,value:", vProb2)
			return false
		}

		/* parse Range3 field */
		vecRange3, _ := parse.GetFieldByName(uint32(i), "range3")
		if vecRange3 != "" {
			arrayRange3 := strings.Split(vecRange3, ",")
			for j := 0; j < len(arrayRange3); j++ {
				v, ret := String2Int32(arrayRange3[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributesValue.Range3 field error, value:", arrayRange3[j])
					return false
				}
				data.Range3 = append(data.Range3, v)
			}
		}

		/* parse Prob3 field */
		vProb3, _ := parse.GetFieldByName(uint32(i), "prob3")
		var Prob3Ret bool
		data.Prob3, Prob3Ret = String2Int32(vProb3)
		if !Prob3Ret {
			glog.Error("Parse CfgEquipRandAttributesValue.Prob3 field error,value:", vProb3)
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

func (c *CfgEquipRandAttributesValueConfig) Clear() {
}

func (c *CfgEquipRandAttributesValueConfig) Find(id int32) (*CfgEquipRandAttributesValue, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEquipRandAttributesValueConfig) GetAllData() map[int32]*CfgEquipRandAttributesValue {
	return c.data
}

func (c *CfgEquipRandAttributesValueConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Range1, ",", v.Prob1, ",", v.Range2, ",", v.Prob2, ",", v.Range3, ",", v.Prob3)
	}
}
