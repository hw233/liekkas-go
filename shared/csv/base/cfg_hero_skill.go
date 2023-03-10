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

type CfgHeroSkill struct {
	Id           int32
	SkillID      int32
	SkillLevel   int32
	CostItems    []string
	Unlock       []string
	InvokeType   int32
	CampFilterId int32
	CombatPower  []int32
	AttrRates    []int32
	AutoTarget   int32
	CrossHairObj string
	CrossHairEnd string
	SoundName    string
	SoundSheet   string
}

type CfgHeroSkillConfig struct {
	data map[int32]*CfgHeroSkill
}

func NewCfgHeroSkillConfig() *CfgHeroSkillConfig {
	return &CfgHeroSkillConfig{
		data: make(map[int32]*CfgHeroSkill),
	}
}

func (c *CfgHeroSkillConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgHeroSkill)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgHeroSkill.Id field error,value:", vId)
			return false
		}

		/* parse SkillID field */
		vSkillID, _ := parse.GetFieldByName(uint32(i), "skillID")
		var SkillIDRet bool
		data.SkillID, SkillIDRet = String2Int32(vSkillID)
		if !SkillIDRet {
			glog.Error("Parse CfgHeroSkill.SkillID field error,value:", vSkillID)
			return false
		}

		/* parse SkillLevel field */
		vSkillLevel, _ := parse.GetFieldByName(uint32(i), "skillLevel")
		var SkillLevelRet bool
		data.SkillLevel, SkillLevelRet = String2Int32(vSkillLevel)
		if !SkillLevelRet {
			glog.Error("Parse CfgHeroSkill.SkillLevel field error,value:", vSkillLevel)
			return false
		}

		/* parse CostItems field */
		vecCostItems, _ := parse.GetFieldByName(uint32(i), "costItems")
		arrayCostItems := strings.Split(vecCostItems, ",")
		for j := 0; j < len(arrayCostItems); j++ {
			v := arrayCostItems[j]
			data.CostItems = append(data.CostItems, v)
		}

		/* parse Unlock field */
		vecUnlock, _ := parse.GetFieldByName(uint32(i), "unlock")
		arrayUnlock := strings.Split(vecUnlock, ",")
		for j := 0; j < len(arrayUnlock); j++ {
			v := arrayUnlock[j]
			data.Unlock = append(data.Unlock, v)
		}

		/* parse InvokeType field */
		vInvokeType, _ := parse.GetFieldByName(uint32(i), "invokeType")
		var InvokeTypeRet bool
		data.InvokeType, InvokeTypeRet = String2Int32(vInvokeType)
		if !InvokeTypeRet {
			glog.Error("Parse CfgHeroSkill.InvokeType field error,value:", vInvokeType)
			return false
		}

		/* parse CampFilterId field */
		vCampFilterId, _ := parse.GetFieldByName(uint32(i), "campFilterId")
		var CampFilterIdRet bool
		data.CampFilterId, CampFilterIdRet = String2Int32(vCampFilterId)
		if !CampFilterIdRet {
			glog.Error("Parse CfgHeroSkill.CampFilterId field error,value:", vCampFilterId)
			return false
		}

		/* parse CombatPower field */
		vecCombatPower, _ := parse.GetFieldByName(uint32(i), "combatPower")
		if vecCombatPower != "" {
			arrayCombatPower := strings.Split(vecCombatPower, ",")
			for j := 0; j < len(arrayCombatPower); j++ {
				v, ret := String2Int32(arrayCombatPower[j])
				if !ret {
					glog.Error("Parse CfgHeroSkill.CombatPower field error, value:", arrayCombatPower[j])
					return false
				}
				data.CombatPower = append(data.CombatPower, v)
			}
		}

		/* parse AttrRates field */
		vecAttrRates, _ := parse.GetFieldByName(uint32(i), "attrRates")
		if vecAttrRates != "" {
			arrayAttrRates := strings.Split(vecAttrRates, ",")
			for j := 0; j < len(arrayAttrRates); j++ {
				v, ret := String2Int32(arrayAttrRates[j])
				if !ret {
					glog.Error("Parse CfgHeroSkill.AttrRates field error, value:", arrayAttrRates[j])
					return false
				}
				data.AttrRates = append(data.AttrRates, v)
			}
		}

		/* parse AutoTarget field */
		vAutoTarget, _ := parse.GetFieldByName(uint32(i), "autoTarget")
		var AutoTargetRet bool
		data.AutoTarget, AutoTargetRet = String2Int32(vAutoTarget)
		if !AutoTargetRet {
			glog.Error("Parse CfgHeroSkill.AutoTarget field error,value:", vAutoTarget)
			return false
		}

		/* parse CrossHairObj field */
		data.CrossHairObj, _ = parse.GetFieldByName(uint32(i), "crossHairObj")

		/* parse CrossHairEnd field */
		data.CrossHairEnd, _ = parse.GetFieldByName(uint32(i), "crossHairEnd")

		/* parse SoundName field */
		data.SoundName, _ = parse.GetFieldByName(uint32(i), "soundName")

		/* parse SoundSheet field */
		data.SoundSheet, _ = parse.GetFieldByName(uint32(i), "soundSheet")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgHeroSkillConfig) Clear() {
}

func (c *CfgHeroSkillConfig) Find(id int32) (*CfgHeroSkill, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgHeroSkillConfig) GetAllData() map[int32]*CfgHeroSkill {
	return c.data
}

func (c *CfgHeroSkillConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.SkillID, ",", v.SkillLevel, ",", v.CostItems, ",", v.Unlock, ",", v.InvokeType, ",", v.CampFilterId, ",", v.CombatPower, ",", v.AttrRates, ",", v.AutoTarget, ",", v.CrossHairObj, ",", v.CrossHairEnd, ",", v.SoundName, ",", v.SoundSheet)
	}
}
