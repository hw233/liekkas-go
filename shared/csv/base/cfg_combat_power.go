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
)

type CfgCombatPower struct {
	Id            int32
	HpMaxFinal    float64
	PhyAtkFinal   float64
	MagAtkFinal   float64
	PhyDfsFinal   float64
	MagDfsFinal   float64
	CritAtkRatio  float64
	CritDfsRatio  float64
	CritAtkValue  float64
	CritDfsValue  float64
	HitRateValue  float64
	EvadeValue    float64
	NormalAtkUp   float64
	NormalDfsUp   float64
	SkillAtkUp    float64
	SkillDfsUp    float64
	UltraAtkUp    float64
	UltraDfsUp    float64
	SkillPhyAtkUp float64
	SkillPhyDfsUp float64
	SkillMagAtkUp float64
	SkillMagDfsUp float64
	CureUp        float64
	HealUp        float64
	PhyPen        float64
	MagPen        float64
	FinalDmg      float64
}

type CfgCombatPowerConfig struct {
	data map[int32]*CfgCombatPower
}

func NewCfgCombatPowerConfig() *CfgCombatPowerConfig {
	return &CfgCombatPowerConfig{
		data: make(map[int32]*CfgCombatPower),
	}
}

func (c *CfgCombatPowerConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCombatPower)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCombatPower.Id field error,value:", vId)
			return false
		}

		/* parse HpMaxFinal field */
		vHpMaxFinal, _ := parse.GetFieldByName(uint32(i), "hp_max_Final")
		var HpMaxFinalRet bool
		data.HpMaxFinal, HpMaxFinalRet = String2Float(vHpMaxFinal)
		if !HpMaxFinalRet {
			glog.Error("Parse CfgCombatPower.HpMaxFinal field error,value:", vHpMaxFinal)
		}

		/* parse PhyAtkFinal field */
		vPhyAtkFinal, _ := parse.GetFieldByName(uint32(i), "phy_atk_Final")
		var PhyAtkFinalRet bool
		data.PhyAtkFinal, PhyAtkFinalRet = String2Float(vPhyAtkFinal)
		if !PhyAtkFinalRet {
			glog.Error("Parse CfgCombatPower.PhyAtkFinal field error,value:", vPhyAtkFinal)
		}

		/* parse MagAtkFinal field */
		vMagAtkFinal, _ := parse.GetFieldByName(uint32(i), "mag_atk_Final")
		var MagAtkFinalRet bool
		data.MagAtkFinal, MagAtkFinalRet = String2Float(vMagAtkFinal)
		if !MagAtkFinalRet {
			glog.Error("Parse CfgCombatPower.MagAtkFinal field error,value:", vMagAtkFinal)
		}

		/* parse PhyDfsFinal field */
		vPhyDfsFinal, _ := parse.GetFieldByName(uint32(i), "phy_dfs_Final")
		var PhyDfsFinalRet bool
		data.PhyDfsFinal, PhyDfsFinalRet = String2Float(vPhyDfsFinal)
		if !PhyDfsFinalRet {
			glog.Error("Parse CfgCombatPower.PhyDfsFinal field error,value:", vPhyDfsFinal)
		}

		/* parse MagDfsFinal field */
		vMagDfsFinal, _ := parse.GetFieldByName(uint32(i), "mag_dfs_Final")
		var MagDfsFinalRet bool
		data.MagDfsFinal, MagDfsFinalRet = String2Float(vMagDfsFinal)
		if !MagDfsFinalRet {
			glog.Error("Parse CfgCombatPower.MagDfsFinal field error,value:", vMagDfsFinal)
		}

		/* parse CritAtkRatio field */
		vCritAtkRatio, _ := parse.GetFieldByName(uint32(i), "crit_atk_ratio")
		var CritAtkRatioRet bool
		data.CritAtkRatio, CritAtkRatioRet = String2Float(vCritAtkRatio)
		if !CritAtkRatioRet {
			glog.Error("Parse CfgCombatPower.CritAtkRatio field error,value:", vCritAtkRatio)
		}

		/* parse CritDfsRatio field */
		vCritDfsRatio, _ := parse.GetFieldByName(uint32(i), "crit_dfs_ratio")
		var CritDfsRatioRet bool
		data.CritDfsRatio, CritDfsRatioRet = String2Float(vCritDfsRatio)
		if !CritDfsRatioRet {
			glog.Error("Parse CfgCombatPower.CritDfsRatio field error,value:", vCritDfsRatio)
		}

		/* parse CritAtkValue field */
		vCritAtkValue, _ := parse.GetFieldByName(uint32(i), "crit_atk_value")
		var CritAtkValueRet bool
		data.CritAtkValue, CritAtkValueRet = String2Float(vCritAtkValue)
		if !CritAtkValueRet {
			glog.Error("Parse CfgCombatPower.CritAtkValue field error,value:", vCritAtkValue)
		}

		/* parse CritDfsValue field */
		vCritDfsValue, _ := parse.GetFieldByName(uint32(i), "crit_dfs_value")
		var CritDfsValueRet bool
		data.CritDfsValue, CritDfsValueRet = String2Float(vCritDfsValue)
		if !CritDfsValueRet {
			glog.Error("Parse CfgCombatPower.CritDfsValue field error,value:", vCritDfsValue)
		}

		/* parse HitRateValue field */
		vHitRateValue, _ := parse.GetFieldByName(uint32(i), "hit_rate_value")
		var HitRateValueRet bool
		data.HitRateValue, HitRateValueRet = String2Float(vHitRateValue)
		if !HitRateValueRet {
			glog.Error("Parse CfgCombatPower.HitRateValue field error,value:", vHitRateValue)
		}

		/* parse EvadeValue field */
		vEvadeValue, _ := parse.GetFieldByName(uint32(i), "evade_value")
		var EvadeValueRet bool
		data.EvadeValue, EvadeValueRet = String2Float(vEvadeValue)
		if !EvadeValueRet {
			glog.Error("Parse CfgCombatPower.EvadeValue field error,value:", vEvadeValue)
		}

		/* parse NormalAtkUp field */
		vNormalAtkUp, _ := parse.GetFieldByName(uint32(i), "normal_atk_up")
		var NormalAtkUpRet bool
		data.NormalAtkUp, NormalAtkUpRet = String2Float(vNormalAtkUp)
		if !NormalAtkUpRet {
			glog.Error("Parse CfgCombatPower.NormalAtkUp field error,value:", vNormalAtkUp)
		}

		/* parse NormalDfsUp field */
		vNormalDfsUp, _ := parse.GetFieldByName(uint32(i), "normal_dfs_up")
		var NormalDfsUpRet bool
		data.NormalDfsUp, NormalDfsUpRet = String2Float(vNormalDfsUp)
		if !NormalDfsUpRet {
			glog.Error("Parse CfgCombatPower.NormalDfsUp field error,value:", vNormalDfsUp)
		}

		/* parse SkillAtkUp field */
		vSkillAtkUp, _ := parse.GetFieldByName(uint32(i), "skillAtkUp")
		var SkillAtkUpRet bool
		data.SkillAtkUp, SkillAtkUpRet = String2Float(vSkillAtkUp)
		if !SkillAtkUpRet {
			glog.Error("Parse CfgCombatPower.SkillAtkUp field error,value:", vSkillAtkUp)
		}

		/* parse SkillDfsUp field */
		vSkillDfsUp, _ := parse.GetFieldByName(uint32(i), "skillDfsUp")
		var SkillDfsUpRet bool
		data.SkillDfsUp, SkillDfsUpRet = String2Float(vSkillDfsUp)
		if !SkillDfsUpRet {
			glog.Error("Parse CfgCombatPower.SkillDfsUp field error,value:", vSkillDfsUp)
		}

		/* parse UltraAtkUp field */
		vUltraAtkUp, _ := parse.GetFieldByName(uint32(i), "ultraAtkUp")
		var UltraAtkUpRet bool
		data.UltraAtkUp, UltraAtkUpRet = String2Float(vUltraAtkUp)
		if !UltraAtkUpRet {
			glog.Error("Parse CfgCombatPower.UltraAtkUp field error,value:", vUltraAtkUp)
		}

		/* parse UltraDfsUp field */
		vUltraDfsUp, _ := parse.GetFieldByName(uint32(i), "ultraDfsUp")
		var UltraDfsUpRet bool
		data.UltraDfsUp, UltraDfsUpRet = String2Float(vUltraDfsUp)
		if !UltraDfsUpRet {
			glog.Error("Parse CfgCombatPower.UltraDfsUp field error,value:", vUltraDfsUp)
		}

		/* parse SkillPhyAtkUp field */
		vSkillPhyAtkUp, _ := parse.GetFieldByName(uint32(i), "skill_phy_atk_up")
		var SkillPhyAtkUpRet bool
		data.SkillPhyAtkUp, SkillPhyAtkUpRet = String2Float(vSkillPhyAtkUp)
		if !SkillPhyAtkUpRet {
			glog.Error("Parse CfgCombatPower.SkillPhyAtkUp field error,value:", vSkillPhyAtkUp)
		}

		/* parse SkillPhyDfsUp field */
		vSkillPhyDfsUp, _ := parse.GetFieldByName(uint32(i), "skill_phy_dfs_up")
		var SkillPhyDfsUpRet bool
		data.SkillPhyDfsUp, SkillPhyDfsUpRet = String2Float(vSkillPhyDfsUp)
		if !SkillPhyDfsUpRet {
			glog.Error("Parse CfgCombatPower.SkillPhyDfsUp field error,value:", vSkillPhyDfsUp)
		}

		/* parse SkillMagAtkUp field */
		vSkillMagAtkUp, _ := parse.GetFieldByName(uint32(i), "skill_mag_atk_up")
		var SkillMagAtkUpRet bool
		data.SkillMagAtkUp, SkillMagAtkUpRet = String2Float(vSkillMagAtkUp)
		if !SkillMagAtkUpRet {
			glog.Error("Parse CfgCombatPower.SkillMagAtkUp field error,value:", vSkillMagAtkUp)
		}

		/* parse SkillMagDfsUp field */
		vSkillMagDfsUp, _ := parse.GetFieldByName(uint32(i), "skill_mag_dfs_up")
		var SkillMagDfsUpRet bool
		data.SkillMagDfsUp, SkillMagDfsUpRet = String2Float(vSkillMagDfsUp)
		if !SkillMagDfsUpRet {
			glog.Error("Parse CfgCombatPower.SkillMagDfsUp field error,value:", vSkillMagDfsUp)
		}

		/* parse CureUp field */
		vCureUp, _ := parse.GetFieldByName(uint32(i), "cure_up")
		var CureUpRet bool
		data.CureUp, CureUpRet = String2Float(vCureUp)
		if !CureUpRet {
			glog.Error("Parse CfgCombatPower.CureUp field error,value:", vCureUp)
		}

		/* parse HealUp field */
		vHealUp, _ := parse.GetFieldByName(uint32(i), "heal_up")
		var HealUpRet bool
		data.HealUp, HealUpRet = String2Float(vHealUp)
		if !HealUpRet {
			glog.Error("Parse CfgCombatPower.HealUp field error,value:", vHealUp)
		}

		/* parse PhyPen field */
		vPhyPen, _ := parse.GetFieldByName(uint32(i), "phyPen")
		var PhyPenRet bool
		data.PhyPen, PhyPenRet = String2Float(vPhyPen)
		if !PhyPenRet {
			glog.Error("Parse CfgCombatPower.PhyPen field error,value:", vPhyPen)
		}

		/* parse MagPen field */
		vMagPen, _ := parse.GetFieldByName(uint32(i), "magPen")
		var MagPenRet bool
		data.MagPen, MagPenRet = String2Float(vMagPen)
		if !MagPenRet {
			glog.Error("Parse CfgCombatPower.MagPen field error,value:", vMagPen)
		}

		/* parse FinalDmg field */
		vFinalDmg, _ := parse.GetFieldByName(uint32(i), "final_dmg")
		var FinalDmgRet bool
		data.FinalDmg, FinalDmgRet = String2Float(vFinalDmg)
		if !FinalDmgRet {
			glog.Error("Parse CfgCombatPower.FinalDmg field error,value:", vFinalDmg)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgCombatPowerConfig) Clear() {
}

func (c *CfgCombatPowerConfig) Find(id int32) (*CfgCombatPower, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCombatPowerConfig) GetAllData() map[int32]*CfgCombatPower {
	return c.data
}

func (c *CfgCombatPowerConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.HpMaxFinal, ",", v.PhyAtkFinal, ",", v.MagAtkFinal, ",", v.PhyDfsFinal, ",", v.MagDfsFinal, ",", v.CritAtkRatio, ",", v.CritDfsRatio, ",", v.CritAtkValue, ",", v.CritDfsValue, ",", v.HitRateValue, ",", v.EvadeValue, ",", v.NormalAtkUp, ",", v.NormalDfsUp, ",", v.SkillAtkUp, ",", v.SkillDfsUp, ",", v.UltraAtkUp, ",", v.UltraDfsUp, ",", v.SkillPhyAtkUp, ",", v.SkillPhyDfsUp, ",", v.SkillMagAtkUp, ",", v.SkillMagDfsUp, ",", v.CureUp, ",", v.HealUp, ",", v.PhyPen, ",", v.MagPen, ",", v.FinalDmg)
	}
}
