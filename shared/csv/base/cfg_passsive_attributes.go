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

type CfgPasssiveAttributes struct {
	Id                  int32
	HpMax               int32
	HpMaxExtra          int32
	HpMaxPercent        int32
	GlobalHpMaxPercent  int32
	PhyAtk              int32
	PhyAtkExtra         int32
	PhyAtkPercent       int32
	GlobalPhyAtkPercent int32
	PhyDfs              int32
	PhyDfsExtra         int32
	PhyDfsPercent       int32
	GlobalPhyDfsPercent int32
	MagAtk              int32
	MagAtkExtra         int32
	MagAtkPercent       int32
	GlobalMagAtkPercent int32
	MagDfs              int32
	MagDfsExtra         int32
	MagDfsPercent       int32
	GlobalMagDfsPercent int32
	CritAtkRatio        int32
	CritDfsRatio        int32
	CritAtkValue        int32
	CritDfsValue        int32
	HitRateValue        int32
	EvadeValue          int32
	NormalAtkUp         int32
	NormalDfsUp         int32
	SkillAtkUp          int32
	SkillDfsUp          int32
	UltraAtkUp          int32
	UltraDfsUp          int32
	SkillPhyAtkUp       int32
	SkillPhyDfsUp       int32
	SkillMagAtkUp       int32
	SkillMagDfsUp       int32
	CureUp              int32
	CureDown            int32
	HealUp              int32
}

type CfgPasssiveAttributesConfig struct {
	data map[int32]*CfgPasssiveAttributes
}

func NewCfgPasssiveAttributesConfig() *CfgPasssiveAttributesConfig {
	return &CfgPasssiveAttributesConfig{
		data: make(map[int32]*CfgPasssiveAttributes),
	}
}

func (c *CfgPasssiveAttributesConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgPasssiveAttributes)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgPasssiveAttributes.Id field error,value:", vId)
			return false
		}

		/* parse HpMax field */
		vHpMax, _ := parse.GetFieldByName(uint32(i), "hpMax")
		var HpMaxRet bool
		data.HpMax, HpMaxRet = String2Int32(vHpMax)
		if !HpMaxRet {
			glog.Error("Parse CfgPasssiveAttributes.HpMax field error,value:", vHpMax)
			return false
		}

		/* parse HpMaxExtra field */
		vHpMaxExtra, _ := parse.GetFieldByName(uint32(i), "hpMaxExtra")
		var HpMaxExtraRet bool
		data.HpMaxExtra, HpMaxExtraRet = String2Int32(vHpMaxExtra)
		if !HpMaxExtraRet {
			glog.Error("Parse CfgPasssiveAttributes.HpMaxExtra field error,value:", vHpMaxExtra)
			return false
		}

		/* parse HpMaxPercent field */
		vHpMaxPercent, _ := parse.GetFieldByName(uint32(i), "hpMaxPercent")
		var HpMaxPercentRet bool
		data.HpMaxPercent, HpMaxPercentRet = String2Int32(vHpMaxPercent)
		if !HpMaxPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.HpMaxPercent field error,value:", vHpMaxPercent)
			return false
		}

		/* parse GlobalHpMaxPercent field */
		vGlobalHpMaxPercent, _ := parse.GetFieldByName(uint32(i), "globalHpMaxPercent")
		var GlobalHpMaxPercentRet bool
		data.GlobalHpMaxPercent, GlobalHpMaxPercentRet = String2Int32(vGlobalHpMaxPercent)
		if !GlobalHpMaxPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.GlobalHpMaxPercent field error,value:", vGlobalHpMaxPercent)
			return false
		}

		/* parse PhyAtk field */
		vPhyAtk, _ := parse.GetFieldByName(uint32(i), "phyAtk")
		var PhyAtkRet bool
		data.PhyAtk, PhyAtkRet = String2Int32(vPhyAtk)
		if !PhyAtkRet {
			glog.Error("Parse CfgPasssiveAttributes.PhyAtk field error,value:", vPhyAtk)
			return false
		}

		/* parse PhyAtkExtra field */
		vPhyAtkExtra, _ := parse.GetFieldByName(uint32(i), "phyAtkExtra")
		var PhyAtkExtraRet bool
		data.PhyAtkExtra, PhyAtkExtraRet = String2Int32(vPhyAtkExtra)
		if !PhyAtkExtraRet {
			glog.Error("Parse CfgPasssiveAttributes.PhyAtkExtra field error,value:", vPhyAtkExtra)
			return false
		}

		/* parse PhyAtkPercent field */
		vPhyAtkPercent, _ := parse.GetFieldByName(uint32(i), "phyAtkPercent")
		var PhyAtkPercentRet bool
		data.PhyAtkPercent, PhyAtkPercentRet = String2Int32(vPhyAtkPercent)
		if !PhyAtkPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.PhyAtkPercent field error,value:", vPhyAtkPercent)
			return false
		}

		/* parse GlobalPhyAtkPercent field */
		vGlobalPhyAtkPercent, _ := parse.GetFieldByName(uint32(i), "globalPhyAtkPercent")
		var GlobalPhyAtkPercentRet bool
		data.GlobalPhyAtkPercent, GlobalPhyAtkPercentRet = String2Int32(vGlobalPhyAtkPercent)
		if !GlobalPhyAtkPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.GlobalPhyAtkPercent field error,value:", vGlobalPhyAtkPercent)
			return false
		}

		/* parse PhyDfs field */
		vPhyDfs, _ := parse.GetFieldByName(uint32(i), "phyDfs")
		var PhyDfsRet bool
		data.PhyDfs, PhyDfsRet = String2Int32(vPhyDfs)
		if !PhyDfsRet {
			glog.Error("Parse CfgPasssiveAttributes.PhyDfs field error,value:", vPhyDfs)
			return false
		}

		/* parse PhyDfsExtra field */
		vPhyDfsExtra, _ := parse.GetFieldByName(uint32(i), "phyDfsExtra")
		var PhyDfsExtraRet bool
		data.PhyDfsExtra, PhyDfsExtraRet = String2Int32(vPhyDfsExtra)
		if !PhyDfsExtraRet {
			glog.Error("Parse CfgPasssiveAttributes.PhyDfsExtra field error,value:", vPhyDfsExtra)
			return false
		}

		/* parse PhyDfsPercent field */
		vPhyDfsPercent, _ := parse.GetFieldByName(uint32(i), "phyDfsPercent")
		var PhyDfsPercentRet bool
		data.PhyDfsPercent, PhyDfsPercentRet = String2Int32(vPhyDfsPercent)
		if !PhyDfsPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.PhyDfsPercent field error,value:", vPhyDfsPercent)
			return false
		}

		/* parse GlobalPhyDfsPercent field */
		vGlobalPhyDfsPercent, _ := parse.GetFieldByName(uint32(i), "globalPhyDfsPercent")
		var GlobalPhyDfsPercentRet bool
		data.GlobalPhyDfsPercent, GlobalPhyDfsPercentRet = String2Int32(vGlobalPhyDfsPercent)
		if !GlobalPhyDfsPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.GlobalPhyDfsPercent field error,value:", vGlobalPhyDfsPercent)
			return false
		}

		/* parse MagAtk field */
		vMagAtk, _ := parse.GetFieldByName(uint32(i), "magAtk")
		var MagAtkRet bool
		data.MagAtk, MagAtkRet = String2Int32(vMagAtk)
		if !MagAtkRet {
			glog.Error("Parse CfgPasssiveAttributes.MagAtk field error,value:", vMagAtk)
			return false
		}

		/* parse MagAtkExtra field */
		vMagAtkExtra, _ := parse.GetFieldByName(uint32(i), "magAtkExtra")
		var MagAtkExtraRet bool
		data.MagAtkExtra, MagAtkExtraRet = String2Int32(vMagAtkExtra)
		if !MagAtkExtraRet {
			glog.Error("Parse CfgPasssiveAttributes.MagAtkExtra field error,value:", vMagAtkExtra)
			return false
		}

		/* parse MagAtkPercent field */
		vMagAtkPercent, _ := parse.GetFieldByName(uint32(i), "magAtkPercent")
		var MagAtkPercentRet bool
		data.MagAtkPercent, MagAtkPercentRet = String2Int32(vMagAtkPercent)
		if !MagAtkPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.MagAtkPercent field error,value:", vMagAtkPercent)
			return false
		}

		/* parse GlobalMagAtkPercent field */
		vGlobalMagAtkPercent, _ := parse.GetFieldByName(uint32(i), "globalMagAtkPercent")
		var GlobalMagAtkPercentRet bool
		data.GlobalMagAtkPercent, GlobalMagAtkPercentRet = String2Int32(vGlobalMagAtkPercent)
		if !GlobalMagAtkPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.GlobalMagAtkPercent field error,value:", vGlobalMagAtkPercent)
			return false
		}

		/* parse MagDfs field */
		vMagDfs, _ := parse.GetFieldByName(uint32(i), "magDfs")
		var MagDfsRet bool
		data.MagDfs, MagDfsRet = String2Int32(vMagDfs)
		if !MagDfsRet {
			glog.Error("Parse CfgPasssiveAttributes.MagDfs field error,value:", vMagDfs)
			return false
		}

		/* parse MagDfsExtra field */
		vMagDfsExtra, _ := parse.GetFieldByName(uint32(i), "magDfsExtra")
		var MagDfsExtraRet bool
		data.MagDfsExtra, MagDfsExtraRet = String2Int32(vMagDfsExtra)
		if !MagDfsExtraRet {
			glog.Error("Parse CfgPasssiveAttributes.MagDfsExtra field error,value:", vMagDfsExtra)
			return false
		}

		/* parse MagDfsPercent field */
		vMagDfsPercent, _ := parse.GetFieldByName(uint32(i), "magDfsPercent")
		var MagDfsPercentRet bool
		data.MagDfsPercent, MagDfsPercentRet = String2Int32(vMagDfsPercent)
		if !MagDfsPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.MagDfsPercent field error,value:", vMagDfsPercent)
			return false
		}

		/* parse GlobalMagDfsPercent field */
		vGlobalMagDfsPercent, _ := parse.GetFieldByName(uint32(i), "globalMagDfsPercent")
		var GlobalMagDfsPercentRet bool
		data.GlobalMagDfsPercent, GlobalMagDfsPercentRet = String2Int32(vGlobalMagDfsPercent)
		if !GlobalMagDfsPercentRet {
			glog.Error("Parse CfgPasssiveAttributes.GlobalMagDfsPercent field error,value:", vGlobalMagDfsPercent)
			return false
		}

		/* parse CritAtkRatio field */
		vCritAtkRatio, _ := parse.GetFieldByName(uint32(i), "critAtkRatio")
		var CritAtkRatioRet bool
		data.CritAtkRatio, CritAtkRatioRet = String2Int32(vCritAtkRatio)
		if !CritAtkRatioRet {
			glog.Error("Parse CfgPasssiveAttributes.CritAtkRatio field error,value:", vCritAtkRatio)
			return false
		}

		/* parse CritDfsRatio field */
		vCritDfsRatio, _ := parse.GetFieldByName(uint32(i), "critDfsRatio")
		var CritDfsRatioRet bool
		data.CritDfsRatio, CritDfsRatioRet = String2Int32(vCritDfsRatio)
		if !CritDfsRatioRet {
			glog.Error("Parse CfgPasssiveAttributes.CritDfsRatio field error,value:", vCritDfsRatio)
			return false
		}

		/* parse CritAtkValue field */
		vCritAtkValue, _ := parse.GetFieldByName(uint32(i), "critAtkValue")
		var CritAtkValueRet bool
		data.CritAtkValue, CritAtkValueRet = String2Int32(vCritAtkValue)
		if !CritAtkValueRet {
			glog.Error("Parse CfgPasssiveAttributes.CritAtkValue field error,value:", vCritAtkValue)
			return false
		}

		/* parse CritDfsValue field */
		vCritDfsValue, _ := parse.GetFieldByName(uint32(i), "critDfsValue")
		var CritDfsValueRet bool
		data.CritDfsValue, CritDfsValueRet = String2Int32(vCritDfsValue)
		if !CritDfsValueRet {
			glog.Error("Parse CfgPasssiveAttributes.CritDfsValue field error,value:", vCritDfsValue)
			return false
		}

		/* parse HitRateValue field */
		vHitRateValue, _ := parse.GetFieldByName(uint32(i), "hitRateValue")
		var HitRateValueRet bool
		data.HitRateValue, HitRateValueRet = String2Int32(vHitRateValue)
		if !HitRateValueRet {
			glog.Error("Parse CfgPasssiveAttributes.HitRateValue field error,value:", vHitRateValue)
			return false
		}

		/* parse EvadeValue field */
		vEvadeValue, _ := parse.GetFieldByName(uint32(i), "evadeValue")
		var EvadeValueRet bool
		data.EvadeValue, EvadeValueRet = String2Int32(vEvadeValue)
		if !EvadeValueRet {
			glog.Error("Parse CfgPasssiveAttributes.EvadeValue field error,value:", vEvadeValue)
			return false
		}

		/* parse NormalAtkUp field */
		vNormalAtkUp, _ := parse.GetFieldByName(uint32(i), "normalAtkUp")
		var NormalAtkUpRet bool
		data.NormalAtkUp, NormalAtkUpRet = String2Int32(vNormalAtkUp)
		if !NormalAtkUpRet {
			glog.Error("Parse CfgPasssiveAttributes.NormalAtkUp field error,value:", vNormalAtkUp)
			return false
		}

		/* parse NormalDfsUp field */
		vNormalDfsUp, _ := parse.GetFieldByName(uint32(i), "normalDfsUp")
		var NormalDfsUpRet bool
		data.NormalDfsUp, NormalDfsUpRet = String2Int32(vNormalDfsUp)
		if !NormalDfsUpRet {
			glog.Error("Parse CfgPasssiveAttributes.NormalDfsUp field error,value:", vNormalDfsUp)
			return false
		}

		/* parse SkillAtkUp field */
		vSkillAtkUp, _ := parse.GetFieldByName(uint32(i), "skillAtkUp")
		var SkillAtkUpRet bool
		data.SkillAtkUp, SkillAtkUpRet = String2Int32(vSkillAtkUp)
		if !SkillAtkUpRet {
			glog.Error("Parse CfgPasssiveAttributes.SkillAtkUp field error,value:", vSkillAtkUp)
			return false
		}

		/* parse SkillDfsUp field */
		vSkillDfsUp, _ := parse.GetFieldByName(uint32(i), "skillDfsUp")
		var SkillDfsUpRet bool
		data.SkillDfsUp, SkillDfsUpRet = String2Int32(vSkillDfsUp)
		if !SkillDfsUpRet {
			glog.Error("Parse CfgPasssiveAttributes.SkillDfsUp field error,value:", vSkillDfsUp)
			return false
		}

		/* parse UltraAtkUp field */
		vUltraAtkUp, _ := parse.GetFieldByName(uint32(i), "ultraAtkUp")
		var UltraAtkUpRet bool
		data.UltraAtkUp, UltraAtkUpRet = String2Int32(vUltraAtkUp)
		if !UltraAtkUpRet {
			glog.Error("Parse CfgPasssiveAttributes.UltraAtkUp field error,value:", vUltraAtkUp)
			return false
		}

		/* parse UltraDfsUp field */
		vUltraDfsUp, _ := parse.GetFieldByName(uint32(i), "ultraDfsUp")
		var UltraDfsUpRet bool
		data.UltraDfsUp, UltraDfsUpRet = String2Int32(vUltraDfsUp)
		if !UltraDfsUpRet {
			glog.Error("Parse CfgPasssiveAttributes.UltraDfsUp field error,value:", vUltraDfsUp)
			return false
		}

		/* parse SkillPhyAtkUp field */
		vSkillPhyAtkUp, _ := parse.GetFieldByName(uint32(i), "skillPhyAtkUp")
		var SkillPhyAtkUpRet bool
		data.SkillPhyAtkUp, SkillPhyAtkUpRet = String2Int32(vSkillPhyAtkUp)
		if !SkillPhyAtkUpRet {
			glog.Error("Parse CfgPasssiveAttributes.SkillPhyAtkUp field error,value:", vSkillPhyAtkUp)
			return false
		}

		/* parse SkillPhyDfsUp field */
		vSkillPhyDfsUp, _ := parse.GetFieldByName(uint32(i), "skillPhyDfsUp")
		var SkillPhyDfsUpRet bool
		data.SkillPhyDfsUp, SkillPhyDfsUpRet = String2Int32(vSkillPhyDfsUp)
		if !SkillPhyDfsUpRet {
			glog.Error("Parse CfgPasssiveAttributes.SkillPhyDfsUp field error,value:", vSkillPhyDfsUp)
			return false
		}

		/* parse SkillMagAtkUp field */
		vSkillMagAtkUp, _ := parse.GetFieldByName(uint32(i), "skillMagAtkUp")
		var SkillMagAtkUpRet bool
		data.SkillMagAtkUp, SkillMagAtkUpRet = String2Int32(vSkillMagAtkUp)
		if !SkillMagAtkUpRet {
			glog.Error("Parse CfgPasssiveAttributes.SkillMagAtkUp field error,value:", vSkillMagAtkUp)
			return false
		}

		/* parse SkillMagDfsUp field */
		vSkillMagDfsUp, _ := parse.GetFieldByName(uint32(i), "skillMagDfsUp")
		var SkillMagDfsUpRet bool
		data.SkillMagDfsUp, SkillMagDfsUpRet = String2Int32(vSkillMagDfsUp)
		if !SkillMagDfsUpRet {
			glog.Error("Parse CfgPasssiveAttributes.SkillMagDfsUp field error,value:", vSkillMagDfsUp)
			return false
		}

		/* parse CureUp field */
		vCureUp, _ := parse.GetFieldByName(uint32(i), "cureUp")
		var CureUpRet bool
		data.CureUp, CureUpRet = String2Int32(vCureUp)
		if !CureUpRet {
			glog.Error("Parse CfgPasssiveAttributes.CureUp field error,value:", vCureUp)
			return false
		}

		/* parse CureDown field */
		vCureDown, _ := parse.GetFieldByName(uint32(i), "cureDown")
		var CureDownRet bool
		data.CureDown, CureDownRet = String2Int32(vCureDown)
		if !CureDownRet {
			glog.Error("Parse CfgPasssiveAttributes.CureDown field error,value:", vCureDown)
			return false
		}

		/* parse HealUp field */
		vHealUp, _ := parse.GetFieldByName(uint32(i), "healUp")
		var HealUpRet bool
		data.HealUp, HealUpRet = String2Int32(vHealUp)
		if !HealUpRet {
			glog.Error("Parse CfgPasssiveAttributes.HealUp field error,value:", vHealUp)
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

func (c *CfgPasssiveAttributesConfig) Clear() {
}

func (c *CfgPasssiveAttributesConfig) Find(id int32) (*CfgPasssiveAttributes, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgPasssiveAttributesConfig) GetAllData() map[int32]*CfgPasssiveAttributes {
	return c.data
}

func (c *CfgPasssiveAttributesConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.HpMax, ",", v.HpMaxExtra, ",", v.HpMaxPercent, ",", v.GlobalHpMaxPercent, ",", v.PhyAtk, ",", v.PhyAtkExtra, ",", v.PhyAtkPercent, ",", v.GlobalPhyAtkPercent, ",", v.PhyDfs, ",", v.PhyDfsExtra, ",", v.PhyDfsPercent, ",", v.GlobalPhyDfsPercent, ",", v.MagAtk, ",", v.MagAtkExtra, ",", v.MagAtkPercent, ",", v.GlobalMagAtkPercent, ",", v.MagDfs, ",", v.MagDfsExtra, ",", v.MagDfsPercent, ",", v.GlobalMagDfsPercent, ",", v.CritAtkRatio, ",", v.CritDfsRatio, ",", v.CritAtkValue, ",", v.CritDfsValue, ",", v.HitRateValue, ",", v.EvadeValue, ",", v.NormalAtkUp, ",", v.NormalDfsUp, ",", v.SkillAtkUp, ",", v.SkillDfsUp, ",", v.UltraAtkUp, ",", v.UltraDfsUp, ",", v.SkillPhyAtkUp, ",", v.SkillPhyDfsUp, ",", v.SkillMagAtkUp, ",", v.SkillMagDfsUp, ",", v.CureUp, ",", v.CureDown, ",", v.HealUp)
	}
}
