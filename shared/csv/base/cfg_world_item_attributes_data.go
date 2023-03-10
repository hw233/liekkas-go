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

type CfgWorldItemAttributesData struct {
	Id                 int32
	WorldItemId        int32
	Star               int32
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
	SkillPhyAtkUp      int32
	SkillPhyDfsUp      int32
	SkillMagAtkUp      int32
	SkillMagDfsUp      int32
	CureUp             int32
	CureDown           int32
	HealUp             int32
	NormalLeechValue   int32
	SkillPhyLeechValue int32
	SkillMagLeechValue int32
	NormalAtkNoDfs     int32
	PhyPen             int32
	MagPen             int32
}

type CfgWorldItemAttributesDataConfig struct {
	data map[int32]*CfgWorldItemAttributesData
}

func NewCfgWorldItemAttributesDataConfig() *CfgWorldItemAttributesDataConfig {
	return &CfgWorldItemAttributesDataConfig{
		data: make(map[int32]*CfgWorldItemAttributesData),
	}
}

func (c *CfgWorldItemAttributesDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgWorldItemAttributesData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgWorldItemAttributesData.Id field error,value:", vId)
			return false
		}

		/* parse WorldItemId field */
		vWorldItemId, _ := parse.GetFieldByName(uint32(i), "worldItemId")
		var WorldItemIdRet bool
		data.WorldItemId, WorldItemIdRet = String2Int32(vWorldItemId)
		if !WorldItemIdRet {
			glog.Error("Parse CfgWorldItemAttributesData.WorldItemId field error,value:", vWorldItemId)
			return false
		}

		/* parse Star field */
		vStar, _ := parse.GetFieldByName(uint32(i), "star")
		var StarRet bool
		data.Star, StarRet = String2Int32(vStar)
		if !StarRet {
			glog.Error("Parse CfgWorldItemAttributesData.Star field error,value:", vStar)
			return false
		}

		/* parse HpExtraRatio field */
		vHpExtraRatio, _ := parse.GetFieldByName(uint32(i), "hp_extra_ratio")
		var HpExtraRatioRet bool
		data.HpExtraRatio, HpExtraRatioRet = String2Float(vHpExtraRatio)
		if !HpExtraRatioRet {
			glog.Error("Parse CfgWorldItemAttributesData.HpExtraRatio field error,value:", vHpExtraRatio)
		}

		/* parse PhyAtkExtraRatio field */
		vPhyAtkExtraRatio, _ := parse.GetFieldByName(uint32(i), "phy_atk_extra_ratio")
		var PhyAtkExtraRatioRet bool
		data.PhyAtkExtraRatio, PhyAtkExtraRatioRet = String2Float(vPhyAtkExtraRatio)
		if !PhyAtkExtraRatioRet {
			glog.Error("Parse CfgWorldItemAttributesData.PhyAtkExtraRatio field error,value:", vPhyAtkExtraRatio)
		}

		/* parse MagAtkExtraRatio field */
		vMagAtkExtraRatio, _ := parse.GetFieldByName(uint32(i), "mag_atk_extra_ratio")
		var MagAtkExtraRatioRet bool
		data.MagAtkExtraRatio, MagAtkExtraRatioRet = String2Float(vMagAtkExtraRatio)
		if !MagAtkExtraRatioRet {
			glog.Error("Parse CfgWorldItemAttributesData.MagAtkExtraRatio field error,value:", vMagAtkExtraRatio)
		}

		/* parse PhyDfsExtraRatio field */
		vPhyDfsExtraRatio, _ := parse.GetFieldByName(uint32(i), "phy_dfs_extra_ratio")
		var PhyDfsExtraRatioRet bool
		data.PhyDfsExtraRatio, PhyDfsExtraRatioRet = String2Float(vPhyDfsExtraRatio)
		if !PhyDfsExtraRatioRet {
			glog.Error("Parse CfgWorldItemAttributesData.PhyDfsExtraRatio field error,value:", vPhyDfsExtraRatio)
		}

		/* parse MagDfsExtraRatio field */
		vMagDfsExtraRatio, _ := parse.GetFieldByName(uint32(i), "mag_dfs_extra_ratio")
		var MagDfsExtraRatioRet bool
		data.MagDfsExtraRatio, MagDfsExtraRatioRet = String2Float(vMagDfsExtraRatio)
		if !MagDfsExtraRatioRet {
			glog.Error("Parse CfgWorldItemAttributesData.MagDfsExtraRatio field error,value:", vMagDfsExtraRatio)
		}

		/* parse HpMaxExtra field */
		vHpMaxExtra, _ := parse.GetFieldByName(uint32(i), "hp_max_extra")
		var HpMaxExtraRet bool
		data.HpMaxExtra, HpMaxExtraRet = String2Int32(vHpMaxExtra)
		if !HpMaxExtraRet {
			glog.Error("Parse CfgWorldItemAttributesData.HpMaxExtra field error,value:", vHpMaxExtra)
			return false
		}

		/* parse PhyAtkExtra field */
		vPhyAtkExtra, _ := parse.GetFieldByName(uint32(i), "phy_atk_extra")
		var PhyAtkExtraRet bool
		data.PhyAtkExtra, PhyAtkExtraRet = String2Int32(vPhyAtkExtra)
		if !PhyAtkExtraRet {
			glog.Error("Parse CfgWorldItemAttributesData.PhyAtkExtra field error,value:", vPhyAtkExtra)
			return false
		}

		/* parse MagAtkExtra field */
		vMagAtkExtra, _ := parse.GetFieldByName(uint32(i), "mag_atk_extra")
		var MagAtkExtraRet bool
		data.MagAtkExtra, MagAtkExtraRet = String2Int32(vMagAtkExtra)
		if !MagAtkExtraRet {
			glog.Error("Parse CfgWorldItemAttributesData.MagAtkExtra field error,value:", vMagAtkExtra)
			return false
		}

		/* parse PhyDfsExtra field */
		vPhyDfsExtra, _ := parse.GetFieldByName(uint32(i), "phy_dfs_extra")
		var PhyDfsExtraRet bool
		data.PhyDfsExtra, PhyDfsExtraRet = String2Int32(vPhyDfsExtra)
		if !PhyDfsExtraRet {
			glog.Error("Parse CfgWorldItemAttributesData.PhyDfsExtra field error,value:", vPhyDfsExtra)
			return false
		}

		/* parse MagDfsExtra field */
		vMagDfsExtra, _ := parse.GetFieldByName(uint32(i), "mag_dfs_extra")
		var MagDfsExtraRet bool
		data.MagDfsExtra, MagDfsExtraRet = String2Int32(vMagDfsExtra)
		if !MagDfsExtraRet {
			glog.Error("Parse CfgWorldItemAttributesData.MagDfsExtra field error,value:", vMagDfsExtra)
			return false
		}

		/* parse HpMaxPercent field */
		vHpMaxPercent, _ := parse.GetFieldByName(uint32(i), "hp_max_percent")
		var HpMaxPercentRet bool
		data.HpMaxPercent, HpMaxPercentRet = String2Int32(vHpMaxPercent)
		if !HpMaxPercentRet {
			glog.Error("Parse CfgWorldItemAttributesData.HpMaxPercent field error,value:", vHpMaxPercent)
			return false
		}

		/* parse PhyAtkPercent field */
		vPhyAtkPercent, _ := parse.GetFieldByName(uint32(i), "phy_atk_percent")
		var PhyAtkPercentRet bool
		data.PhyAtkPercent, PhyAtkPercentRet = String2Int32(vPhyAtkPercent)
		if !PhyAtkPercentRet {
			glog.Error("Parse CfgWorldItemAttributesData.PhyAtkPercent field error,value:", vPhyAtkPercent)
			return false
		}

		/* parse MagAtkPercent field */
		vMagAtkPercent, _ := parse.GetFieldByName(uint32(i), "mag_atk_percent")
		var MagAtkPercentRet bool
		data.MagAtkPercent, MagAtkPercentRet = String2Int32(vMagAtkPercent)
		if !MagAtkPercentRet {
			glog.Error("Parse CfgWorldItemAttributesData.MagAtkPercent field error,value:", vMagAtkPercent)
			return false
		}

		/* parse PhyDfsPercent field */
		vPhyDfsPercent, _ := parse.GetFieldByName(uint32(i), "phy_dfs_percent")
		var PhyDfsPercentRet bool
		data.PhyDfsPercent, PhyDfsPercentRet = String2Int32(vPhyDfsPercent)
		if !PhyDfsPercentRet {
			glog.Error("Parse CfgWorldItemAttributesData.PhyDfsPercent field error,value:", vPhyDfsPercent)
			return false
		}

		/* parse MagDfsPercent field */
		vMagDfsPercent, _ := parse.GetFieldByName(uint32(i), "mag_dfs_percent")
		var MagDfsPercentRet bool
		data.MagDfsPercent, MagDfsPercentRet = String2Int32(vMagDfsPercent)
		if !MagDfsPercentRet {
			glog.Error("Parse CfgWorldItemAttributesData.MagDfsPercent field error,value:", vMagDfsPercent)
			return false
		}

		/* parse CritAtkRatio field */
		vCritAtkRatio, _ := parse.GetFieldByName(uint32(i), "crit_atk_ratio")
		var CritAtkRatioRet bool
		data.CritAtkRatio, CritAtkRatioRet = String2Int32(vCritAtkRatio)
		if !CritAtkRatioRet {
			glog.Error("Parse CfgWorldItemAttributesData.CritAtkRatio field error,value:", vCritAtkRatio)
			return false
		}

		/* parse CritDfsRatio field */
		vCritDfsRatio, _ := parse.GetFieldByName(uint32(i), "crit_dfs_ratio")
		var CritDfsRatioRet bool
		data.CritDfsRatio, CritDfsRatioRet = String2Int32(vCritDfsRatio)
		if !CritDfsRatioRet {
			glog.Error("Parse CfgWorldItemAttributesData.CritDfsRatio field error,value:", vCritDfsRatio)
			return false
		}

		/* parse CritAtkValue field */
		vCritAtkValue, _ := parse.GetFieldByName(uint32(i), "crit_atk_value")
		var CritAtkValueRet bool
		data.CritAtkValue, CritAtkValueRet = String2Int32(vCritAtkValue)
		if !CritAtkValueRet {
			glog.Error("Parse CfgWorldItemAttributesData.CritAtkValue field error,value:", vCritAtkValue)
			return false
		}

		/* parse CritDfsValue field */
		vCritDfsValue, _ := parse.GetFieldByName(uint32(i), "crit_dfs_value")
		var CritDfsValueRet bool
		data.CritDfsValue, CritDfsValueRet = String2Int32(vCritDfsValue)
		if !CritDfsValueRet {
			glog.Error("Parse CfgWorldItemAttributesData.CritDfsValue field error,value:", vCritDfsValue)
			return false
		}

		/* parse HitRateValue field */
		vHitRateValue, _ := parse.GetFieldByName(uint32(i), "hit_rate_value")
		var HitRateValueRet bool
		data.HitRateValue, HitRateValueRet = String2Int32(vHitRateValue)
		if !HitRateValueRet {
			glog.Error("Parse CfgWorldItemAttributesData.HitRateValue field error,value:", vHitRateValue)
			return false
		}

		/* parse EvadeValue field */
		vEvadeValue, _ := parse.GetFieldByName(uint32(i), "evade_value")
		var EvadeValueRet bool
		data.EvadeValue, EvadeValueRet = String2Int32(vEvadeValue)
		if !EvadeValueRet {
			glog.Error("Parse CfgWorldItemAttributesData.EvadeValue field error,value:", vEvadeValue)
			return false
		}

		/* parse NormalAtkUp field */
		vNormalAtkUp, _ := parse.GetFieldByName(uint32(i), "normal_atk_up")
		var NormalAtkUpRet bool
		data.NormalAtkUp, NormalAtkUpRet = String2Int32(vNormalAtkUp)
		if !NormalAtkUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.NormalAtkUp field error,value:", vNormalAtkUp)
			return false
		}

		/* parse NormalDfsUp field */
		vNormalDfsUp, _ := parse.GetFieldByName(uint32(i), "normal_dfs_up")
		var NormalDfsUpRet bool
		data.NormalDfsUp, NormalDfsUpRet = String2Int32(vNormalDfsUp)
		if !NormalDfsUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.NormalDfsUp field error,value:", vNormalDfsUp)
			return false
		}

		/* parse SkillAtkUp field */
		vSkillAtkUp, _ := parse.GetFieldByName(uint32(i), "skillAtkUp")
		var SkillAtkUpRet bool
		data.SkillAtkUp, SkillAtkUpRet = String2Int32(vSkillAtkUp)
		if !SkillAtkUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.SkillAtkUp field error,value:", vSkillAtkUp)
			return false
		}

		/* parse SkillDfsUp field */
		vSkillDfsUp, _ := parse.GetFieldByName(uint32(i), "skillDfsUp")
		var SkillDfsUpRet bool
		data.SkillDfsUp, SkillDfsUpRet = String2Int32(vSkillDfsUp)
		if !SkillDfsUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.SkillDfsUp field error,value:", vSkillDfsUp)
			return false
		}

		/* parse UltraAtkUp field */
		vUltraAtkUp, _ := parse.GetFieldByName(uint32(i), "ultraAtkUp")
		var UltraAtkUpRet bool
		data.UltraAtkUp, UltraAtkUpRet = String2Int32(vUltraAtkUp)
		if !UltraAtkUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.UltraAtkUp field error,value:", vUltraAtkUp)
			return false
		}

		/* parse UltraDfsUp field */
		vUltraDfsUp, _ := parse.GetFieldByName(uint32(i), "ultraDfsUp")
		var UltraDfsUpRet bool
		data.UltraDfsUp, UltraDfsUpRet = String2Int32(vUltraDfsUp)
		if !UltraDfsUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.UltraDfsUp field error,value:", vUltraDfsUp)
			return false
		}

		/* parse SkillPhyAtkUp field */
		vSkillPhyAtkUp, _ := parse.GetFieldByName(uint32(i), "skill_phy_atk_up")
		var SkillPhyAtkUpRet bool
		data.SkillPhyAtkUp, SkillPhyAtkUpRet = String2Int32(vSkillPhyAtkUp)
		if !SkillPhyAtkUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.SkillPhyAtkUp field error,value:", vSkillPhyAtkUp)
			return false
		}

		/* parse SkillPhyDfsUp field */
		vSkillPhyDfsUp, _ := parse.GetFieldByName(uint32(i), "skill_phy_dfs_up")
		var SkillPhyDfsUpRet bool
		data.SkillPhyDfsUp, SkillPhyDfsUpRet = String2Int32(vSkillPhyDfsUp)
		if !SkillPhyDfsUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.SkillPhyDfsUp field error,value:", vSkillPhyDfsUp)
			return false
		}

		/* parse SkillMagAtkUp field */
		vSkillMagAtkUp, _ := parse.GetFieldByName(uint32(i), "skill_mag_atk_up")
		var SkillMagAtkUpRet bool
		data.SkillMagAtkUp, SkillMagAtkUpRet = String2Int32(vSkillMagAtkUp)
		if !SkillMagAtkUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.SkillMagAtkUp field error,value:", vSkillMagAtkUp)
			return false
		}

		/* parse SkillMagDfsUp field */
		vSkillMagDfsUp, _ := parse.GetFieldByName(uint32(i), "skill_mag_dfs_up")
		var SkillMagDfsUpRet bool
		data.SkillMagDfsUp, SkillMagDfsUpRet = String2Int32(vSkillMagDfsUp)
		if !SkillMagDfsUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.SkillMagDfsUp field error,value:", vSkillMagDfsUp)
			return false
		}

		/* parse CureUp field */
		vCureUp, _ := parse.GetFieldByName(uint32(i), "cure_up")
		var CureUpRet bool
		data.CureUp, CureUpRet = String2Int32(vCureUp)
		if !CureUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.CureUp field error,value:", vCureUp)
			return false
		}

		/* parse CureDown field */
		vCureDown, _ := parse.GetFieldByName(uint32(i), "cure_down")
		var CureDownRet bool
		data.CureDown, CureDownRet = String2Int32(vCureDown)
		if !CureDownRet {
			glog.Error("Parse CfgWorldItemAttributesData.CureDown field error,value:", vCureDown)
			return false
		}

		/* parse HealUp field */
		vHealUp, _ := parse.GetFieldByName(uint32(i), "heal_up")
		var HealUpRet bool
		data.HealUp, HealUpRet = String2Int32(vHealUp)
		if !HealUpRet {
			glog.Error("Parse CfgWorldItemAttributesData.HealUp field error,value:", vHealUp)
			return false
		}

		/* parse NormalLeechValue field */
		vNormalLeechValue, _ := parse.GetFieldByName(uint32(i), "normal_leech_value")
		var NormalLeechValueRet bool
		data.NormalLeechValue, NormalLeechValueRet = String2Int32(vNormalLeechValue)
		if !NormalLeechValueRet {
			glog.Error("Parse CfgWorldItemAttributesData.NormalLeechValue field error,value:", vNormalLeechValue)
			return false
		}

		/* parse SkillPhyLeechValue field */
		vSkillPhyLeechValue, _ := parse.GetFieldByName(uint32(i), "skill_phy_leech_value")
		var SkillPhyLeechValueRet bool
		data.SkillPhyLeechValue, SkillPhyLeechValueRet = String2Int32(vSkillPhyLeechValue)
		if !SkillPhyLeechValueRet {
			glog.Error("Parse CfgWorldItemAttributesData.SkillPhyLeechValue field error,value:", vSkillPhyLeechValue)
			return false
		}

		/* parse SkillMagLeechValue field */
		vSkillMagLeechValue, _ := parse.GetFieldByName(uint32(i), "skill_mag_leech_value")
		var SkillMagLeechValueRet bool
		data.SkillMagLeechValue, SkillMagLeechValueRet = String2Int32(vSkillMagLeechValue)
		if !SkillMagLeechValueRet {
			glog.Error("Parse CfgWorldItemAttributesData.SkillMagLeechValue field error,value:", vSkillMagLeechValue)
			return false
		}

		/* parse NormalAtkNoDfs field */
		vNormalAtkNoDfs, _ := parse.GetFieldByName(uint32(i), "normal_atk_no_dfs")
		var NormalAtkNoDfsRet bool
		data.NormalAtkNoDfs, NormalAtkNoDfsRet = String2Int32(vNormalAtkNoDfs)
		if !NormalAtkNoDfsRet {
			glog.Error("Parse CfgWorldItemAttributesData.NormalAtkNoDfs field error,value:", vNormalAtkNoDfs)
			return false
		}

		/* parse PhyPen field */
		vPhyPen, _ := parse.GetFieldByName(uint32(i), "phyPen")
		var PhyPenRet bool
		data.PhyPen, PhyPenRet = String2Int32(vPhyPen)
		if !PhyPenRet {
			glog.Error("Parse CfgWorldItemAttributesData.PhyPen field error,value:", vPhyPen)
			return false
		}

		/* parse MagPen field */
		vMagPen, _ := parse.GetFieldByName(uint32(i), "magPen")
		var MagPenRet bool
		data.MagPen, MagPenRet = String2Int32(vMagPen)
		if !MagPenRet {
			glog.Error("Parse CfgWorldItemAttributesData.MagPen field error,value:", vMagPen)
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

func (c *CfgWorldItemAttributesDataConfig) Clear() {
}

func (c *CfgWorldItemAttributesDataConfig) Find(id int32) (*CfgWorldItemAttributesData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgWorldItemAttributesDataConfig) GetAllData() map[int32]*CfgWorldItemAttributesData {
	return c.data
}

func (c *CfgWorldItemAttributesDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.WorldItemId, ",", v.Star, ",", v.HpExtraRatio, ",", v.PhyAtkExtraRatio, ",", v.MagAtkExtraRatio, ",", v.PhyDfsExtraRatio, ",", v.MagDfsExtraRatio, ",", v.HpMaxExtra, ",", v.PhyAtkExtra, ",", v.MagAtkExtra, ",", v.PhyDfsExtra, ",", v.MagDfsExtra, ",", v.HpMaxPercent, ",", v.PhyAtkPercent, ",", v.MagAtkPercent, ",", v.PhyDfsPercent, ",", v.MagDfsPercent, ",", v.CritAtkRatio, ",", v.CritDfsRatio, ",", v.CritAtkValue, ",", v.CritDfsValue, ",", v.HitRateValue, ",", v.EvadeValue, ",", v.NormalAtkUp, ",", v.NormalDfsUp, ",", v.SkillAtkUp, ",", v.SkillDfsUp, ",", v.UltraAtkUp, ",", v.UltraDfsUp, ",", v.SkillPhyAtkUp, ",", v.SkillPhyDfsUp, ",", v.SkillMagAtkUp, ",", v.SkillMagDfsUp, ",", v.CureUp, ",", v.CureDown, ",", v.HealUp, ",", v.NormalLeechValue, ",", v.SkillPhyLeechValue, ",", v.SkillMagLeechValue, ",", v.NormalAtkNoDfs, ",", v.PhyPen, ",", v.MagPen)
	}
}
