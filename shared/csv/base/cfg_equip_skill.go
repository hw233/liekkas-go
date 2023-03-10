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

type CfgEquipSkill struct {
	Id          int32
	Comment     string
	SkillLevel  int32
	CombatPower []int32
}

type CfgEquipSkillConfig struct {
	data map[int32]*CfgEquipSkill
}

func NewCfgEquipSkillConfig() *CfgEquipSkillConfig {
	return &CfgEquipSkillConfig{
		data: make(map[int32]*CfgEquipSkill),
	}
}

func (c *CfgEquipSkillConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEquipSkill)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEquipSkill.Id field error,value:", vId)
			return false
		}

		/* parse Comment field */
		data.Comment, _ = parse.GetFieldByName(uint32(i), "comment")

		/* parse SkillLevel field */
		vSkillLevel, _ := parse.GetFieldByName(uint32(i), "skillLevel")
		var SkillLevelRet bool
		data.SkillLevel, SkillLevelRet = String2Int32(vSkillLevel)
		if !SkillLevelRet {
			glog.Error("Parse CfgEquipSkill.SkillLevel field error,value:", vSkillLevel)
			return false
		}

		/* parse CombatPower field */
		vecCombatPower, _ := parse.GetFieldByName(uint32(i), "combatPower")
		if vecCombatPower != "" {
			arrayCombatPower := strings.Split(vecCombatPower, ",")
			for j := 0; j < len(arrayCombatPower); j++ {
				v, ret := String2Int32(arrayCombatPower[j])
				if !ret {
					glog.Error("Parse CfgEquipSkill.CombatPower field error, value:", arrayCombatPower[j])
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

func (c *CfgEquipSkillConfig) Clear() {
}

func (c *CfgEquipSkillConfig) Find(id int32) (*CfgEquipSkill, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEquipSkillConfig) GetAllData() map[int32]*CfgEquipSkill {
	return c.data
}

func (c *CfgEquipSkillConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Comment, ",", v.SkillLevel, ",", v.CombatPower)
	}
}
