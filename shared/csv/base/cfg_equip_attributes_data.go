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

type CfgEquipAttributesData struct {
	Id                 int32
	EquipID            int32
	Stage              int32
	HpExtraRatio       float64
	PhyAtkExtraRatio   float64
	MagAtkExtraRatio   float64
	PhyDfsExtraRatio   float64
	MagDfsExtraRatio   float64
	HpMaxExtra         int32
	PhyAtkExtra        int32
	MagAtkExtra        int32
	PhyDfsExtra        int32
	MagDfsExtra        int32
	HpMaxPercent       int32
	PhyAtkPercent      int32
	MagAtkPercent      int32
	PhyDfsPercent      int32
	MagDfsPercent      int32
	CritAtkRatio       int32
	CritDfsRatio       int32
	CritAtkValue       int32
	CritDfsValue       int32
	HitRateValue       int32
	EvadeValue         int32
	NormalAtkUp        int32
	NormalDfsUp        int32
	SkillAtkUp         int32
	SkillDfsUp         int32
	UltraAtkUp         int32
	UltraDfsUp         int32
	PhyDamAdd          int32
	PhyDamReduce       int32
	MagDamAdd          int32
	MagDamReduce       int32
	CureUp             int32
	CureDown           int32
	HealUp             int32
	NormalLeechValue   int32
	SkillPhyLeechValue int32
	SkillMagLeechValue int32
	NormalAtkNoDfs     int32
	SkillPhyNoDfs      int32
	SkillMagNoDfs      int32
	PhyPen             int32
	MagPen             int32
}

type CfgEquipAttributesDataConfig struct {
	data map[int32]*CfgEquipAttributesData
}

func NewCfgEquipAttributesDataConfig() *CfgEquipAttributesDataConfig {
	return &CfgEquipAttributesDataConfig{
		data: make(map[int32]*CfgEquipAttributesData),
	}
}

func (c *CfgEquipAttributesDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgEquipAttributesData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgEquipAttributesData.Id field error,value:", vId)
			return false
		}

		/* parse EquipID field */
		vEquipID, _ := parse.GetFieldByName(uint32(i), "equipID")
		var EquipIDRet bool
		data.EquipID, EquipIDRet = String2Int32(vEquipID)
		if !EquipIDRet {
			glog.Error("Parse CfgEquipAttributesData.EquipID field error,value:", vEquipID)
			return false
		}

		/* parse Stage field */
		vStage, _ := parse.GetFieldByName(uint32(i), "stage")
		var StageRet bool
		data.Stage, StageRet = String2Int32(vStage)
		if !StageRet {
			glog.Error("Parse CfgEquipAttributesData.Stage field error,value:", vStage)
			return false
		}

		/* parse HpExtraRatio field */
		vHpExtraRatio, _ := parse.GetFieldByName(uint32(i), "hp_extra_ratio")
		var HpExtraRatioRet bool
		data.HpExtraRatio, HpExtraRatioRet = String2Float(vHpExtraRatio)
		if !HpExtraRatioRet {
			glog.Error("Parse CfgEquipAttributesData.HpExtraRatio field error,value:", vHpExtraRatio)
		}

		/* parse PhyAtkExtraRatio field */
		vPhyAtkExtraRatio, _ := parse.GetFieldByName(uint32(i), "phy_atk_extra_ratio")
		var PhyAtkExtraRatioRet bool
		data.PhyAtkExtraRatio, PhyAtkExtraRatioRet = String2Float(vPhyAtkExtraRatio)
		if !PhyAtkExtraRatioRet {
			glog.Error("Parse CfgEquipAttributesData.PhyAtkExtraRatio field error,value:", vPhyAtkExtraRatio)
		}

		/* parse MagAtkExtraRatio field */
		vMagAtkExtraRatio, _ := parse.GetFieldByName(uint32(i), "mag_atk_extra_ratio")
		var MagAtkExtraRatioRet bool
		data.MagAtkExtraRatio, MagAtkExtraRatioRet = String2Float(vMagAtkExtraRatio)
		if !MagAtkExtraRatioRet {
			glog.Error("Parse CfgEquipAttributesData.MagAtkExtraRatio field error,value:", vMagAtkExtraRatio)
		}

		/* parse PhyDfsExtraRatio field */
		vPhyDfsExtraRatio, _ := parse.GetFieldByName(uint32(i), "phy_dfs_extra_ratio")
		var PhyDfsExtraRatioRet bool
		data.PhyDfsExtraRatio, PhyDfsExtraRatioRet = String2Float(vPhyDfsExtraRatio)
		if !PhyDfsExtraRatioRet {
			glog.Error("Parse CfgEquipAttributesData.PhyDfsExtraRatio field error,value:", vPhyDfsExtraRatio)
		}

		/* parse MagDfsExtraRatio field */
		vMagDfsExtraRatio, _ := parse.GetFieldByName(uint32(i), "mag_dfs_extra_ratio")
		var MagDfsExtraRatioRet bool
		data.MagDfsExtraRatio, MagDfsExtraRatioRet = String2Float(vMagDfsExtraRatio)
		if !MagDfsExtraRatioRet {
			glog.Error("Parse CfgEquipAttributesData.MagDfsExtraRatio field error,value:", vMagDfsExtraRatio)
		}

		/* parse HpMaxExtra field */
		vHpMaxExtra, _ := parse.GetFieldByName(uint32(i), "hp_max_extra")
		var HpMaxExtraRet bool
		data.HpMaxExtra, HpMaxExtraRet = String2Int32(vHpMaxExtra)
		if !HpMaxExtraRet {
			glog.Error("Parse CfgEquipAttributesData.HpMaxExtra field error,value:", vHpMaxExtra)
			return false
		}

		/* parse PhyAtkExtra field */
		vPhyAtkExtra, _ := parse.GetFieldByName(uint32(i), "phy_atk_extra")
		var PhyAtkExtraRet bool
		data.PhyAtkExtra, PhyAtkExtraRet = String2Int32(vPhyAtkExtra)
		if !PhyAtkExtraRet {
			glog.Error("Parse CfgEquipAttributesData.PhyAtkExtra field error,value:", vPhyAtkExtra)
			return false
		}

		/* parse MagAtkExtra field */
		vMagAtkExtra, _ := parse.GetFieldByName(uint32(i), "mag_atk_extra")
		var MagAtkExtraRet bool
		data.MagAtkExtra, MagAtkExtraRet = String2Int32(vMagAtkExtra)
		if !MagAtkExtraRet {
			glog.Error("Parse CfgEquipAttributesData.MagAtkExtra field error,value:", vMagAtkExtra)
			return false
		}

		/* parse PhyDfsExtra field */
		vPhyDfsExtra, _ := parse.GetFieldByName(uint32(i), "phy_dfs_extra")
		var PhyDfsExtraRet bool
		data.PhyDfsExtra, PhyDfsExtraRet = String2Int32(vPhyDfsExtra)
		if !PhyDfsExtraRet {
			glog.Error("Parse CfgEquipAttributesData.PhyDfsExtra field error,value:", vPhyDfsExtra)
			return false
		}

		/* parse MagDfsExtra field */
		vMagDfsExtra, _ := parse.GetFieldByName(uint32(i), "mag_dfs_extra")
		var MagDfsExtraRet bool
		data.MagDfsExtra, MagDfsExtraRet = String2Int32(vMagDfsExtra)
		if !MagDfsExtraRet {
			glog.Error("Parse CfgEquipAttributesData.MagDfsExtra field error,value:", vMagDfsExtra)
			return false
		}

		/* parse HpMaxPercent field */
		vHpMaxPercent, _ := parse.GetFieldByName(uint32(i), "hp_max_percent")
		var HpMaxPercentRet bool
		data.HpMaxPercent, HpMaxPercentRet = String2Int32(vHpMaxPercent)
		if !HpMaxPercentRet {
			glog.Error("Parse CfgEquipAttributesData.HpMaxPercent field error,value:", vHpMaxPercent)
			return false
		}

		/* parse PhyAtkPercent field */
		vPhyAtkPercent, _ := parse.GetFieldByName(uint32(i), "phy_atk_percent")
		var PhyAtkPercentRet bool
		data.PhyAtkPercent, PhyAtkPercentRet = String2Int32(vPhyAtkPercent)
		if !PhyAtkPercentRet {
			glog.Error("Parse CfgEquipAttributesData.PhyAtkPercent field error,value:", vPhyAtkPercent)
			return false
		}

		/* parse MagAtkPercent field */
		vMagAtkPercent, _ := parse.GetFieldByName(uint32(i), "mag_atk_percent")
		var MagAtkPercentRet bool
		data.MagAtkPercent, MagAtkPercentRet = String2Int32(vMagAtkPercent)
		if !MagAtkPercentRet {
			glog.Error("Parse CfgEquipAttributesData.MagAtkPercent field error,value:", vMagAtkPercent)
			return false
		}

		/* parse PhyDfsPercent field */
		vPhyDfsPercent, _ := parse.GetFieldByName(uint32(i), "phy_dfs_percent")
		var PhyDfsPercentRet bool
		data.PhyDfsPercent, PhyDfsPercentRet = String2Int32(vPhyDfsPercent)
		if !PhyDfsPercentRet {
			glog.Error("Parse CfgEquipAttributesData.PhyDfsPercent field error,value:", vPhyDfsPercent)
			return false
		}

		/* parse MagDfsPercent field */
		vMagDfsPercent, _ := parse.GetFieldByName(uint32(i), "mag_dfs_percent")
		var MagDfsPercentRet bool
		data.MagDfsPercent, MagDfsPercentRet = String2Int32(vMagDfsPercent)
		if !MagDfsPercentRet {
			glog.Error("Parse CfgEquipAttributesData.MagDfsPercent field error,value:", vMagDfsPercent)
			return false
		}

		/* parse CritAtkRatio field */
		vCritAtkRatio, _ := parse.GetFieldByName(uint32(i), "crit_atk_ratio")
		var CritAtkRatioRet bool
		data.CritAtkRatio, CritAtkRatioRet = String2Int32(vCritAtkRatio)
		if !CritAtkRatioRet {
			glog.Error("Parse CfgEquipAttributesData.CritAtkRatio field error,value:", vCritAtkRatio)
			return false
		}

		/* parse CritDfsRatio field */
		vCritDfsRatio, _ := parse.GetFieldByName(uint32(i), "crit_dfs_ratio")
		var CritDfsRatioRet bool
		data.CritDfsRatio, CritDfsRatioRet = String2Int32(vCritDfsRatio)
		if !CritDfsRatioRet {
			glog.Error("Parse CfgEquipAttributesData.CritDfsRatio field error,value:", vCritDfsRatio)
			return false
		}

		/* parse CritAtkValue field */
		vCritAtkValue, _ := parse.GetFieldByName(uint32(i), "crit_atk_value")
		var CritAtkValueRet bool
		data.CritAtkValue, CritAtkValueRet = String2Int32(vCritAtkValue)
		if !CritAtkValueRet {
			glog.Error("Parse CfgEquipAttributesData.CritAtkValue field error,value:", vCritAtkValue)
			return false
		}

		/* parse CritDfsValue field */
		vCritDfsValue, _ := parse.GetFieldByName(uint32(i), "crit_dfs_value")
		var CritDfsValueRet bool
		data.CritDfsValue, CritDfsValueRet = String2Int32(vCritDfsValue)
		if !CritDfsValueRet {
			glog.Error("Parse CfgEquipAttributesData.CritDfsValue field error,value:", vCritDfsValue)
			return false
		}

		/* parse HitRateValue field */
		vHitRateValue, _ := parse.GetFieldByName(uint32(i), "hit_rate_value")
		var HitRateValueRet bool
		data.HitRateValue, HitRateValueRet = String2Int32(vHitRateValue)
		if !HitRateValueRet {
			glog.Error("Parse CfgEquipAttributesData.HitRateValue field error,value:", vHitRateValue)
			return false
		}

		/* parse EvadeValue field */
		vEvadeValue, _ := parse.GetFieldByName(uint32(i), "evade_value")
		var EvadeValueRet bool
		data.EvadeValue, EvadeValueRet = String2Int32(vEvadeValue)
		if !EvadeValueRet {
			glog.Error("Parse CfgEquipAttributesData.EvadeValue field error,value:", vEvadeValue)
			return false
		}

		/* parse NormalAtkUp field */
		vNormalAtkUp, _ := parse.GetFieldByName(uint32(i), "normal_atk_up")
		var NormalAtkUpRet bool
		data.NormalAtkUp, NormalAtkUpRet = String2Int32(vNormalAtkUp)
		if !NormalAtkUpRet {
			glog.Error("Parse CfgEquipAttributesData.NormalAtkUp field error,value:", vNormalAtkUp)
			return false
		}

		/* parse NormalDfsUp field */
		vNormalDfsUp, _ := parse.GetFieldByName(uint32(i), "normal_dfs_up")
		var NormalDfsUpRet bool
		data.NormalDfsUp, NormalDfsUpRet = String2Int32(vNormalDfsUp)
		if !NormalDfsUpRet {
			glog.Error("Parse CfgEquipAttributesData.NormalDfsUp field error,value:", vNormalDfsUp)
			return false
		}

		/* parse SkillAtkUp field */
		vSkillAtkUp, _ := parse.GetFieldByName(uint32(i), "skillAtkUp")
		var SkillAtkUpRet bool
		data.SkillAtkUp, SkillAtkUpRet = String2Int32(vSkillAtkUp)
		if !SkillAtkUpRet {
			glog.Error("Parse CfgEquipAttributesData.SkillAtkUp field error,value:", vSkillAtkUp)
			return false
		}

		/* parse SkillDfsUp field */
		vSkillDfsUp, _ := parse.GetFieldByName(uint32(i), "skillDfsUp")
		var SkillDfsUpRet bool
		data.SkillDfsUp, SkillDfsUpRet = String2Int32(vSkillDfsUp)
		if !SkillDfsUpRet {
			glog.Error("Parse CfgEquipAttributesData.SkillDfsUp field error,value:", vSkillDfsUp)
			return false
		}

		/* parse UltraAtkUp field */
		vUltraAtkUp, _ := parse.GetFieldByName(uint32(i), "ultraAtkUp")
		var UltraAtkUpRet bool
		data.UltraAtkUp, UltraAtkUpRet = String2Int32(vUltraAtkUp)
		if !UltraAtkUpRet {
			glog.Error("Parse CfgEquipAttributesData.UltraAtkUp field error,value:", vUltraAtkUp)
			return false
		}

		/* parse UltraDfsUp field */
		vUltraDfsUp, _ := parse.GetFieldByName(uint32(i), "ultraDfsUp")
		var UltraDfsUpRet bool
		data.UltraDfsUp, UltraDfsUpRet = String2Int32(vUltraDfsUp)
		if !UltraDfsUpRet {
			glog.Error("Parse CfgEquipAttributesData.UltraDfsUp field error,value:", vUltraDfsUp)
			return false
		}

		/* parse PhyDamAdd field */
		vPhyDamAdd, _ := parse.GetFieldByName(uint32(i), "phyDamAdd")
		var PhyDamAddRet bool
		data.PhyDamAdd, PhyDamAddRet = String2Int32(vPhyDamAdd)
		if !PhyDamAddRet {
			glog.Error("Parse CfgEquipAttributesData.PhyDamAdd field error,value:", vPhyDamAdd)
			return false
		}

		/* parse PhyDamReduce field */
		vPhyDamReduce, _ := parse.GetFieldByName(uint32(i), "phyDamReduce")
		var PhyDamReduceRet bool
		data.PhyDamReduce, PhyDamReduceRet = String2Int32(vPhyDamReduce)
		if !PhyDamReduceRet {
			glog.Error("Parse CfgEquipAttributesData.PhyDamReduce field error,value:", vPhyDamReduce)
			return false
		}

		/* parse MagDamAdd field */
		vMagDamAdd, _ := parse.GetFieldByName(uint32(i), "magDamAdd")
		var MagDamAddRet bool
		data.MagDamAdd, MagDamAddRet = String2Int32(vMagDamAdd)
		if !MagDamAddRet {
			glog.Error("Parse CfgEquipAttributesData.MagDamAdd field error,value:", vMagDamAdd)
			return false
		}

		/* parse MagDamReduce field */
		vMagDamReduce, _ := parse.GetFieldByName(uint32(i), "magDamReduce")
		var MagDamReduceRet bool
		data.MagDamReduce, MagDamReduceRet = String2Int32(vMagDamReduce)
		if !MagDamReduceRet {
			glog.Error("Parse CfgEquipAttributesData.MagDamReduce field error,value:", vMagDamReduce)
			return false
		}

		/* parse CureUp field */
		vCureUp, _ := parse.GetFieldByName(uint32(i), "cure_up")
		var CureUpRet bool
		data.CureUp, CureUpRet = String2Int32(vCureUp)
		if !CureUpRet {
			glog.Error("Parse CfgEquipAttributesData.CureUp field error,value:", vCureUp)
			return false
		}

		/* parse CureDown field */
		vCureDown, _ := parse.GetFieldByName(uint32(i), "cure_down")
		var CureDownRet bool
		data.CureDown, CureDownRet = String2Int32(vCureDown)
		if !CureDownRet {
			glog.Error("Parse CfgEquipAttributesData.CureDown field error,value:", vCureDown)
			return false
		}

		/* parse HealUp field */
		vHealUp, _ := parse.GetFieldByName(uint32(i), "heal_up")
		var HealUpRet bool
		data.HealUp, HealUpRet = String2Int32(vHealUp)
		if !HealUpRet {
			glog.Error("Parse CfgEquipAttributesData.HealUp field error,value:", vHealUp)
			return false
		}

		/* parse NormalLeechValue field */
		vNormalLeechValue, _ := parse.GetFieldByName(uint32(i), "normal_leech_value")
		var NormalLeechValueRet bool
		data.NormalLeechValue, NormalLeechValueRet = String2Int32(vNormalLeechValue)
		if !NormalLeechValueRet {
			glog.Error("Parse CfgEquipAttributesData.NormalLeechValue field error,value:", vNormalLeechValue)
			return false
		}

		/* parse SkillPhyLeechValue field */
		vSkillPhyLeechValue, _ := parse.GetFieldByName(uint32(i), "skill_phy_leech_value")
		var SkillPhyLeechValueRet bool
		data.SkillPhyLeechValue, SkillPhyLeechValueRet = String2Int32(vSkillPhyLeechValue)
		if !SkillPhyLeechValueRet {
			glog.Error("Parse CfgEquipAttributesData.SkillPhyLeechValue field error,value:", vSkillPhyLeechValue)
			return false
		}

		/* parse SkillMagLeechValue field */
		vSkillMagLeechValue, _ := parse.GetFieldByName(uint32(i), "skill_mag_leech_value")
		var SkillMagLeechValueRet bool
		data.SkillMagLeechValue, SkillMagLeechValueRet = String2Int32(vSkillMagLeechValue)
		if !SkillMagLeechValueRet {
			glog.Error("Parse CfgEquipAttributesData.SkillMagLeechValue field error,value:", vSkillMagLeechValue)
			return false
		}

		/* parse NormalAtkNoDfs field */
		vNormalAtkNoDfs, _ := parse.GetFieldByName(uint32(i), "normal_atk_no_dfs")
		var NormalAtkNoDfsRet bool
		data.NormalAtkNoDfs, NormalAtkNoDfsRet = String2Int32(vNormalAtkNoDfs)
		if !NormalAtkNoDfsRet {
			glog.Error("Parse CfgEquipAttributesData.NormalAtkNoDfs field error,value:", vNormalAtkNoDfs)
			return false
		}

		/* parse SkillPhyNoDfs field */
		vSkillPhyNoDfs, _ := parse.GetFieldByName(uint32(i), "skill_phy_no_dfs")
		var SkillPhyNoDfsRet bool
		data.SkillPhyNoDfs, SkillPhyNoDfsRet = String2Int32(vSkillPhyNoDfs)
		if !SkillPhyNoDfsRet {
			glog.Error("Parse CfgEquipAttributesData.SkillPhyNoDfs field error,value:", vSkillPhyNoDfs)
			return false
		}

		/* parse SkillMagNoDfs field */
		vSkillMagNoDfs, _ := parse.GetFieldByName(uint32(i), "skill_mag_no_dfs")
		var SkillMagNoDfsRet bool
		data.SkillMagNoDfs, SkillMagNoDfsRet = String2Int32(vSkillMagNoDfs)
		if !SkillMagNoDfsRet {
			glog.Error("Parse CfgEquipAttributesData.SkillMagNoDfs field error,value:", vSkillMagNoDfs)
			return false
		}

		/* parse PhyPen field */
		vPhyPen, _ := parse.GetFieldByName(uint32(i), "phyPen")
		var PhyPenRet bool
		data.PhyPen, PhyPenRet = String2Int32(vPhyPen)
		if !PhyPenRet {
			glog.Error("Parse CfgEquipAttributesData.PhyPen field error,value:", vPhyPen)
			return false
		}

		/* parse MagPen field */
		vMagPen, _ := parse.GetFieldByName(uint32(i), "magPen")
		var MagPenRet bool
		data.MagPen, MagPenRet = String2Int32(vMagPen)
		if !MagPenRet {
			glog.Error("Parse CfgEquipAttributesData.MagPen field error,value:", vMagPen)
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

func (c *CfgEquipAttributesDataConfig) Clear() {
}

func (c *CfgEquipAttributesDataConfig) Find(id int32) (*CfgEquipAttributesData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgEquipAttributesDataConfig) GetAllData() map[int32]*CfgEquipAttributesData {
	return c.data
}

func (c *CfgEquipAttributesDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.EquipID, ",", v.Stage, ",", v.HpExtraRatio, ",", v.PhyAtkExtraRatio, ",", v.MagAtkExtraRatio, ",", v.PhyDfsExtraRatio, ",", v.MagDfsExtraRatio, ",", v.HpMaxExtra, ",", v.PhyAtkExtra, ",", v.MagAtkExtra, ",", v.PhyDfsExtra, ",", v.MagDfsExtra, ",", v.HpMaxPercent, ",", v.PhyAtkPercent, ",", v.MagAtkPercent, ",", v.PhyDfsPercent, ",", v.MagDfsPercent, ",", v.CritAtkRatio, ",", v.CritDfsRatio, ",", v.CritAtkValue, ",", v.CritDfsValue, ",", v.HitRateValue, ",", v.EvadeValue, ",", v.NormalAtkUp, ",", v.NormalDfsUp, ",", v.SkillAtkUp, ",", v.SkillDfsUp, ",", v.UltraAtkUp, ",", v.UltraDfsUp, ",", v.PhyDamAdd, ",", v.PhyDamReduce, ",", v.MagDamAdd, ",", v.MagDamReduce, ",", v.CureUp, ",", v.CureDown, ",", v.HealUp, ",", v.NormalLeechValue, ",", v.SkillPhyLeechValue, ",", v.SkillMagLeechValue, ",", v.NormalAtkNoDfs, ",", v.SkillPhyNoDfs, ",", v.SkillMagNoDfs, ",", v.PhyPen, ",", v.MagPen)
	}
}
