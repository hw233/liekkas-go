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

type CfgCharacterSkill struct {
	Id                    int32
	RoleID                int32
	SkillNumber           int32
	AssociatedSkillNumber int32
	SkillType             int32
	IsPassive             bool
	SkillLevel            int32
	SkillUnlock           int32
	UnlockParam           int32
	StarAutoLevelUp       int32
	Cost                  []string
	CombatPower           []int32
}

type CfgCharacterSkillConfig struct {
	data map[int32]*CfgCharacterSkill
}

func NewCfgCharacterSkillConfig() *CfgCharacterSkillConfig {
	return &CfgCharacterSkillConfig{
		data: make(map[int32]*CfgCharacterSkill),
	}
}

func (c *CfgCharacterSkillConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterSkill)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterSkill.Id field error,value:", vId)
			return false
		}

		/* parse RoleID field */
		vRoleID, _ := parse.GetFieldByName(uint32(i), "roleID")
		var RoleIDRet bool
		data.RoleID, RoleIDRet = String2Int32(vRoleID)
		if !RoleIDRet {
			glog.Error("Parse CfgCharacterSkill.RoleID field error,value:", vRoleID)
			return false
		}

		/* parse SkillNumber field */
		vSkillNumber, _ := parse.GetFieldByName(uint32(i), "skillNumber")
		var SkillNumberRet bool
		data.SkillNumber, SkillNumberRet = String2Int32(vSkillNumber)
		if !SkillNumberRet {
			glog.Error("Parse CfgCharacterSkill.SkillNumber field error,value:", vSkillNumber)
			return false
		}

		/* parse AssociatedSkillNumber field */
		vAssociatedSkillNumber, _ := parse.GetFieldByName(uint32(i), "associatedSkillNumber")
		var AssociatedSkillNumberRet bool
		data.AssociatedSkillNumber, AssociatedSkillNumberRet = String2Int32(vAssociatedSkillNumber)
		if !AssociatedSkillNumberRet {
			glog.Error("Parse CfgCharacterSkill.AssociatedSkillNumber field error,value:", vAssociatedSkillNumber)
			return false
		}

		/* parse SkillType field */
		vSkillType, _ := parse.GetFieldByName(uint32(i), "skillType")
		var SkillTypeRet bool
		data.SkillType, SkillTypeRet = String2Int32(vSkillType)
		if !SkillTypeRet {
			glog.Error("Parse CfgCharacterSkill.SkillType field error,value:", vSkillType)
			return false
		}

		/* parse IsPassive field */
		vIsPassive, _ := parse.GetFieldByName(uint32(i), "isPassive")
		var IsPassiveRet bool
		data.IsPassive, IsPassiveRet = String2Bool(vIsPassive)
		if !IsPassiveRet {
			glog.Error("Parse CfgCharacterSkill.IsPassive field error,value:", vIsPassive)
		}

		/* parse SkillLevel field */
		vSkillLevel, _ := parse.GetFieldByName(uint32(i), "skillLevel")
		var SkillLevelRet bool
		data.SkillLevel, SkillLevelRet = String2Int32(vSkillLevel)
		if !SkillLevelRet {
			glog.Error("Parse CfgCharacterSkill.SkillLevel field error,value:", vSkillLevel)
			return false
		}

		/* parse SkillUnlock field */
		vSkillUnlock, _ := parse.GetFieldByName(uint32(i), "skillUnlock")
		var SkillUnlockRet bool
		data.SkillUnlock, SkillUnlockRet = String2Int32(vSkillUnlock)
		if !SkillUnlockRet {
			glog.Error("Parse CfgCharacterSkill.SkillUnlock field error,value:", vSkillUnlock)
			return false
		}

		/* parse UnlockParam field */
		vUnlockParam, _ := parse.GetFieldByName(uint32(i), "unlockParam")
		var UnlockParamRet bool
		data.UnlockParam, UnlockParamRet = String2Int32(vUnlockParam)
		if !UnlockParamRet {
			glog.Error("Parse CfgCharacterSkill.UnlockParam field error,value:", vUnlockParam)
			return false
		}

		/* parse StarAutoLevelUp field */
		vStarAutoLevelUp, _ := parse.GetFieldByName(uint32(i), "starAutoLevelUp")
		var StarAutoLevelUpRet bool
		data.StarAutoLevelUp, StarAutoLevelUpRet = String2Int32(vStarAutoLevelUp)
		if !StarAutoLevelUpRet {
			glog.Error("Parse CfgCharacterSkill.StarAutoLevelUp field error,value:", vStarAutoLevelUp)
			return false
		}

		/* parse Cost field */
		vecCost, _ := parse.GetFieldByName(uint32(i), "cost")
		arrayCost := strings.Split(vecCost, ",")
		for j := 0; j < len(arrayCost); j++ {
			v := arrayCost[j]
			data.Cost = append(data.Cost, v)
		}

		/* parse CombatPower field */
		vecCombatPower, _ := parse.GetFieldByName(uint32(i), "combatPower")
		if vecCombatPower != "" {
			arrayCombatPower := strings.Split(vecCombatPower, ",")
			for j := 0; j < len(arrayCombatPower); j++ {
				v, ret := String2Int32(arrayCombatPower[j])
				if !ret {
					glog.Error("Parse CfgCharacterSkill.CombatPower field error, value:", arrayCombatPower[j])
					return false
				}
				data.CombatPower = append(data.CombatPower, v)
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

func (c *CfgCharacterSkillConfig) Clear() {
}

func (c *CfgCharacterSkillConfig) Find(id int32) (*CfgCharacterSkill, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterSkillConfig) GetAllData() map[int32]*CfgCharacterSkill {
	return c.data
}

func (c *CfgCharacterSkillConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.RoleID, ",", v.SkillNumber, ",", v.AssociatedSkillNumber, ",", v.SkillType, ",", v.IsPassive, ",", v.SkillLevel, ",", v.SkillUnlock, ",", v.UnlockParam, ",", v.StarAutoLevelUp, ",", v.Cost, ",", v.CombatPower)
	}
}
