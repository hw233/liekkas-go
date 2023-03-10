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

type CfgEquipRandAttributes struct {
	Id                   int32
	HpAddPercent         []int32
	AtkAddPercent        []int32
	MAtkAddPercent       []int32
	DefAddPercent        []int32
	MDefAddPercent       []int32
	CritAddPercent       []int32
	CritDamPercent       []int32
	CritDamReducePercent []int32
	CureAddPercent       []int32
	PhyDamAdd            []int32
	PhyDamReduce         []int32
	MagDamAdd            []int32
	MagDamReduce         []int32
	PhyPen               []int32
	MagPen               []int32
}

type CfgEquipRandAttributesConfig struct {
	data map[int32]*CfgEquipRandAttributes
}

func NewCfgEquipRandAttributesConfig() *CfgEquipRandAttributesConfig {
	return &CfgEquipRandAttributesConfig{
		data: make(map[int32]*CfgEquipRandAttributes),
	}
}

func (c *CfgEquipRandAttributesConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEquipRandAttributes)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEquipRandAttributes.Id field error,value:", vId)
			return false
		}

		/* parse HpAddPercent field */
		vecHpAddPercent, _ := parse.GetFieldByName(uint32(i), "hpAddPercent")
		if vecHpAddPercent != "" {
			arrayHpAddPercent := strings.Split(vecHpAddPercent, ",")
			for j := 0; j < len(arrayHpAddPercent); j++ {
				v, ret := String2Int32(arrayHpAddPercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.HpAddPercent field error, value:", arrayHpAddPercent[j])
					return false
				}
				data.HpAddPercent = append(data.HpAddPercent, v)
			}
		}

		/* parse AtkAddPercent field */
		vecAtkAddPercent, _ := parse.GetFieldByName(uint32(i), "atkAddPercent")
		if vecAtkAddPercent != "" {
			arrayAtkAddPercent := strings.Split(vecAtkAddPercent, ",")
			for j := 0; j < len(arrayAtkAddPercent); j++ {
				v, ret := String2Int32(arrayAtkAddPercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.AtkAddPercent field error, value:", arrayAtkAddPercent[j])
					return false
				}
				data.AtkAddPercent = append(data.AtkAddPercent, v)
			}
		}

		/* parse MAtkAddPercent field */
		vecMAtkAddPercent, _ := parse.GetFieldByName(uint32(i), "mAtkAddPercent")
		if vecMAtkAddPercent != "" {
			arrayMAtkAddPercent := strings.Split(vecMAtkAddPercent, ",")
			for j := 0; j < len(arrayMAtkAddPercent); j++ {
				v, ret := String2Int32(arrayMAtkAddPercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.MAtkAddPercent field error, value:", arrayMAtkAddPercent[j])
					return false
				}
				data.MAtkAddPercent = append(data.MAtkAddPercent, v)
			}
		}

		/* parse DefAddPercent field */
		vecDefAddPercent, _ := parse.GetFieldByName(uint32(i), "defAddPercent")
		if vecDefAddPercent != "" {
			arrayDefAddPercent := strings.Split(vecDefAddPercent, ",")
			for j := 0; j < len(arrayDefAddPercent); j++ {
				v, ret := String2Int32(arrayDefAddPercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.DefAddPercent field error, value:", arrayDefAddPercent[j])
					return false
				}
				data.DefAddPercent = append(data.DefAddPercent, v)
			}
		}

		/* parse MDefAddPercent field */
		vecMDefAddPercent, _ := parse.GetFieldByName(uint32(i), "mDefAddPercent")
		if vecMDefAddPercent != "" {
			arrayMDefAddPercent := strings.Split(vecMDefAddPercent, ",")
			for j := 0; j < len(arrayMDefAddPercent); j++ {
				v, ret := String2Int32(arrayMDefAddPercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.MDefAddPercent field error, value:", arrayMDefAddPercent[j])
					return false
				}
				data.MDefAddPercent = append(data.MDefAddPercent, v)
			}
		}

		/* parse CritAddPercent field */
		vecCritAddPercent, _ := parse.GetFieldByName(uint32(i), "critAddPercent")
		if vecCritAddPercent != "" {
			arrayCritAddPercent := strings.Split(vecCritAddPercent, ",")
			for j := 0; j < len(arrayCritAddPercent); j++ {
				v, ret := String2Int32(arrayCritAddPercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.CritAddPercent field error, value:", arrayCritAddPercent[j])
					return false
				}
				data.CritAddPercent = append(data.CritAddPercent, v)
			}
		}

		/* parse CritDamPercent field */
		vecCritDamPercent, _ := parse.GetFieldByName(uint32(i), "critDamPercent")
		if vecCritDamPercent != "" {
			arrayCritDamPercent := strings.Split(vecCritDamPercent, ",")
			for j := 0; j < len(arrayCritDamPercent); j++ {
				v, ret := String2Int32(arrayCritDamPercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.CritDamPercent field error, value:", arrayCritDamPercent[j])
					return false
				}
				data.CritDamPercent = append(data.CritDamPercent, v)
			}
		}

		/* parse CritDamReducePercent field */
		vecCritDamReducePercent, _ := parse.GetFieldByName(uint32(i), "critDamReducePercent")
		if vecCritDamReducePercent != "" {
			arrayCritDamReducePercent := strings.Split(vecCritDamReducePercent, ",")
			for j := 0; j < len(arrayCritDamReducePercent); j++ {
				v, ret := String2Int32(arrayCritDamReducePercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.CritDamReducePercent field error, value:", arrayCritDamReducePercent[j])
					return false
				}
				data.CritDamReducePercent = append(data.CritDamReducePercent, v)
			}
		}

		/* parse CureAddPercent field */
		vecCureAddPercent, _ := parse.GetFieldByName(uint32(i), "cureAddPercent")
		if vecCureAddPercent != "" {
			arrayCureAddPercent := strings.Split(vecCureAddPercent, ",")
			for j := 0; j < len(arrayCureAddPercent); j++ {
				v, ret := String2Int32(arrayCureAddPercent[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.CureAddPercent field error, value:", arrayCureAddPercent[j])
					return false
				}
				data.CureAddPercent = append(data.CureAddPercent, v)
			}
		}

		/* parse PhyDamAdd field */
		vecPhyDamAdd, _ := parse.GetFieldByName(uint32(i), "phyDamAdd")
		if vecPhyDamAdd != "" {
			arrayPhyDamAdd := strings.Split(vecPhyDamAdd, ",")
			for j := 0; j < len(arrayPhyDamAdd); j++ {
				v, ret := String2Int32(arrayPhyDamAdd[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.PhyDamAdd field error, value:", arrayPhyDamAdd[j])
					return false
				}
				data.PhyDamAdd = append(data.PhyDamAdd, v)
			}
		}

		/* parse PhyDamReduce field */
		vecPhyDamReduce, _ := parse.GetFieldByName(uint32(i), "phyDamReduce")
		if vecPhyDamReduce != "" {
			arrayPhyDamReduce := strings.Split(vecPhyDamReduce, ",")
			for j := 0; j < len(arrayPhyDamReduce); j++ {
				v, ret := String2Int32(arrayPhyDamReduce[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.PhyDamReduce field error, value:", arrayPhyDamReduce[j])
					return false
				}
				data.PhyDamReduce = append(data.PhyDamReduce, v)
			}
		}

		/* parse MagDamAdd field */
		vecMagDamAdd, _ := parse.GetFieldByName(uint32(i), "magDamAdd")
		if vecMagDamAdd != "" {
			arrayMagDamAdd := strings.Split(vecMagDamAdd, ",")
			for j := 0; j < len(arrayMagDamAdd); j++ {
				v, ret := String2Int32(arrayMagDamAdd[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.MagDamAdd field error, value:", arrayMagDamAdd[j])
					return false
				}
				data.MagDamAdd = append(data.MagDamAdd, v)
			}
		}

		/* parse MagDamReduce field */
		vecMagDamReduce, _ := parse.GetFieldByName(uint32(i), "magDamReduce")
		if vecMagDamReduce != "" {
			arrayMagDamReduce := strings.Split(vecMagDamReduce, ",")
			for j := 0; j < len(arrayMagDamReduce); j++ {
				v, ret := String2Int32(arrayMagDamReduce[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.MagDamReduce field error, value:", arrayMagDamReduce[j])
					return false
				}
				data.MagDamReduce = append(data.MagDamReduce, v)
			}
		}

		/* parse PhyPen field */
		vecPhyPen, _ := parse.GetFieldByName(uint32(i), "phyPen")
		if vecPhyPen != "" {
			arrayPhyPen := strings.Split(vecPhyPen, ",")
			for j := 0; j < len(arrayPhyPen); j++ {
				v, ret := String2Int32(arrayPhyPen[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.PhyPen field error, value:", arrayPhyPen[j])
					return false
				}
				data.PhyPen = append(data.PhyPen, v)
			}
		}

		/* parse MagPen field */
		vecMagPen, _ := parse.GetFieldByName(uint32(i), "magPen")
		if vecMagPen != "" {
			arrayMagPen := strings.Split(vecMagPen, ",")
			for j := 0; j < len(arrayMagPen); j++ {
				v, ret := String2Int32(arrayMagPen[j])
				if !ret {
					glog.Error("Parse CfgEquipRandAttributes.MagPen field error, value:", arrayMagPen[j])
					return false
				}
				data.MagPen = append(data.MagPen, v)
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

func (c *CfgEquipRandAttributesConfig) Clear() {
}

func (c *CfgEquipRandAttributesConfig) Find(id int32) (*CfgEquipRandAttributes, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEquipRandAttributesConfig) GetAllData() map[int32]*CfgEquipRandAttributes {
	return c.data
}

func (c *CfgEquipRandAttributesConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.HpAddPercent, ",", v.AtkAddPercent, ",", v.MAtkAddPercent, ",", v.DefAddPercent, ",", v.MDefAddPercent, ",", v.CritAddPercent, ",", v.CritDamPercent, ",", v.CritDamReducePercent, ",", v.CureAddPercent, ",", v.PhyDamAdd, ",", v.PhyDamReduce, ",", v.MagDamAdd, ",", v.MagDamReduce, ",", v.PhyPen, ",", v.MagPen)
	}
}
