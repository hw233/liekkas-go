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

type CfgEquipData struct {
	Id           int32
	Name         string
	Rarity       int32
	UseType      int32
	UseParam     int32
	SellPrice    string
	Part         int32
	EquipLmt     []int32
	RepickCost   []string
	CampAddition int32
	CampWeight   []int32
	GoldperExp   float64
	IsLocked     bool
	PassiveID    string
	AttributeID  string
	CanProduced  bool
	BreakEquipId []int32
	RandEntries  []int32
}

type CfgEquipDataConfig struct {
	data map[int32]*CfgEquipData
}

func NewCfgEquipDataConfig() *CfgEquipDataConfig {
	return &CfgEquipDataConfig{
		data: make(map[int32]*CfgEquipData),
	}
}

func (c *CfgEquipDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEquipData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEquipData.Id field error,value:", vId)
			return false
		}

		/* parse Name field */
		data.Name, _ = parse.GetFieldByName(uint32(i), "name")

		/* parse Rarity field */
		vRarity, _ := parse.GetFieldByName(uint32(i), "rarity")
		var RarityRet bool
		data.Rarity, RarityRet = String2Int32(vRarity)
		if !RarityRet {
			glog.Error("Parse CfgEquipData.Rarity field error,value:", vRarity)
			return false
		}

		/* parse UseType field */
		vUseType, _ := parse.GetFieldByName(uint32(i), "useType")
		var UseTypeRet bool
		data.UseType, UseTypeRet = String2Int32(vUseType)
		if !UseTypeRet {
			glog.Error("Parse CfgEquipData.UseType field error,value:", vUseType)
			return false
		}

		/* parse UseParam field */
		vUseParam, _ := parse.GetFieldByName(uint32(i), "useParam")
		var UseParamRet bool
		data.UseParam, UseParamRet = String2Int32(vUseParam)
		if !UseParamRet {
			glog.Error("Parse CfgEquipData.UseParam field error,value:", vUseParam)
			return false
		}

		/* parse SellPrice field */
		data.SellPrice, _ = parse.GetFieldByName(uint32(i), "sellPrice")

		/* parse Part field */
		vPart, _ := parse.GetFieldByName(uint32(i), "part")
		var PartRet bool
		data.Part, PartRet = String2Int32(vPart)
		if !PartRet {
			glog.Error("Parse CfgEquipData.Part field error,value:", vPart)
			return false
		}

		/* parse EquipLmt field */
		vecEquipLmt, _ := parse.GetFieldByName(uint32(i), "equipLmt")
		if vecEquipLmt != "" {
			arrayEquipLmt := strings.Split(vecEquipLmt, ",")
			for j := 0; j < len(arrayEquipLmt); j++ {
				v, ret := String2Int32(arrayEquipLmt[j])
				if !ret {
					glog.Error("Parse CfgEquipData.EquipLmt field error, value:", arrayEquipLmt[j])
					return false
				}
				data.EquipLmt = append(data.EquipLmt, v)
			}
		}

		/* parse RepickCost field */
		vecRepickCost, _ := parse.GetFieldByName(uint32(i), "repickCost")
		arrayRepickCost := strings.Split(vecRepickCost, ",")
		for j := 0; j < len(arrayRepickCost); j++ {
			v := arrayRepickCost[j]
			data.RepickCost = append(data.RepickCost, v)
		}

		/* parse CampAddition field */
		vCampAddition, _ := parse.GetFieldByName(uint32(i), "campAddition")
		var CampAdditionRet bool
		data.CampAddition, CampAdditionRet = String2Int32(vCampAddition)
		if !CampAdditionRet {
			glog.Error("Parse CfgEquipData.CampAddition field error,value:", vCampAddition)
			return false
		}

		/* parse CampWeight field */
		vecCampWeight, _ := parse.GetFieldByName(uint32(i), "campWeight")
		if vecCampWeight != "" {
			arrayCampWeight := strings.Split(vecCampWeight, ",")
			for j := 0; j < len(arrayCampWeight); j++ {
				v, ret := String2Int32(arrayCampWeight[j])
				if !ret {
					glog.Error("Parse CfgEquipData.CampWeight field error, value:", arrayCampWeight[j])
					return false
				}
				data.CampWeight = append(data.CampWeight, v)
			}
		}

		/* parse GoldperExp field */
		vGoldperExp, _ := parse.GetFieldByName(uint32(i), "goldperExp")
		var GoldperExpRet bool
		data.GoldperExp, GoldperExpRet = String2Float(vGoldperExp)
		if !GoldperExpRet {
			glog.Error("Parse CfgEquipData.GoldperExp field error,value:", vGoldperExp)
		}

		/* parse IsLocked field */
		vIsLocked, _ := parse.GetFieldByName(uint32(i), "isLocked")
		var IsLockedRet bool
		data.IsLocked, IsLockedRet = String2Bool(vIsLocked)
		if !IsLockedRet {
			glog.Error("Parse CfgEquipData.IsLocked field error,value:", vIsLocked)
		}

		/* parse PassiveID field */
		data.PassiveID, _ = parse.GetFieldByName(uint32(i), "passiveID")

		/* parse AttributeID field */
		data.AttributeID, _ = parse.GetFieldByName(uint32(i), "attributeID")

		/* parse CanProduced field */
		vCanProduced, _ := parse.GetFieldByName(uint32(i), "canProduced")
		var CanProducedRet bool
		data.CanProduced, CanProducedRet = String2Bool(vCanProduced)
		if !CanProducedRet {
			glog.Error("Parse CfgEquipData.CanProduced field error,value:", vCanProduced)
		}

		/* parse BreakEquipId field */
		vecBreakEquipId, _ := parse.GetFieldByName(uint32(i), "breakEquipId")
		if vecBreakEquipId != "" {
			arrayBreakEquipId := strings.Split(vecBreakEquipId, ",")
			for j := 0; j < len(arrayBreakEquipId); j++ {
				v, ret := String2Int32(arrayBreakEquipId[j])
				if !ret {
					glog.Error("Parse CfgEquipData.BreakEquipId field error, value:", arrayBreakEquipId[j])
					return false
				}
				data.BreakEquipId = append(data.BreakEquipId, v)
			}
		}

		/* parse RandEntries field */
		vecRandEntries, _ := parse.GetFieldByName(uint32(i), "randEntries")
		if vecRandEntries != "" {
			arrayRandEntries := strings.Split(vecRandEntries, ",")
			for j := 0; j < len(arrayRandEntries); j++ {
				v, ret := String2Int32(arrayRandEntries[j])
				if !ret {
					glog.Error("Parse CfgEquipData.RandEntries field error, value:", arrayRandEntries[j])
					return false
				}
				data.RandEntries = append(data.RandEntries, v)
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

func (c *CfgEquipDataConfig) Clear() {
}

func (c *CfgEquipDataConfig) Find(id int32) (*CfgEquipData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEquipDataConfig) GetAllData() map[int32]*CfgEquipData {
	return c.data
}

func (c *CfgEquipDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Name, ",", v.Rarity, ",", v.UseType, ",", v.UseParam, ",", v.SellPrice, ",", v.Part, ",", v.EquipLmt, ",", v.RepickCost, ",", v.CampAddition, ",", v.CampWeight, ",", v.GoldperExp, ",", v.IsLocked, ",", v.PassiveID, ",", v.AttributeID, ",", v.CanProduced, ",", v.BreakEquipId, ",", v.RandEntries)
	}
}
