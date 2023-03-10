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

type CfgHeroExtraData struct {
	Id             int32
	HeroID         int32
	HeroLv         int32
	RarityDiscount []float64
	CombatPower    int32
	HpMax          int32
	PhyAtk         int32
	MagAtk         int32
	PhyDfs         int32
	MagDfs         int32
	CritAtkRatio   int32
	CritAtkValue   int32
}

type CfgHeroExtraDataConfig struct {
	data map[int32]*CfgHeroExtraData
}

func NewCfgHeroExtraDataConfig() *CfgHeroExtraDataConfig {
	return &CfgHeroExtraDataConfig{
		data: make(map[int32]*CfgHeroExtraData),
	}
}

func (c *CfgHeroExtraDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgHeroExtraData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgHeroExtraData.Id field error,value:", vId)
			return false
		}

		/* parse HeroID field */
		vHeroID, _ := parse.GetFieldByName(uint32(i), "heroID")
		var HeroIDRet bool
		data.HeroID, HeroIDRet = String2Int32(vHeroID)
		if !HeroIDRet {
			glog.Error("Parse CfgHeroExtraData.HeroID field error,value:", vHeroID)
			return false
		}

		/* parse HeroLv field */
		vHeroLv, _ := parse.GetFieldByName(uint32(i), "heroLv")
		var HeroLvRet bool
		data.HeroLv, HeroLvRet = String2Int32(vHeroLv)
		if !HeroLvRet {
			glog.Error("Parse CfgHeroExtraData.HeroLv field error,value:", vHeroLv)
			return false
		}

		/* parse RarityDiscount field */
		vecRarityDiscount, _ := parse.GetFieldByName(uint32(i), "rarityDiscount")
		arrayRarityDiscount := strings.Split(vecRarityDiscount, ",")
		for j := 0; j < len(arrayRarityDiscount); j++ {
			v, ret := String2Float(arrayRarityDiscount[j])
			if !ret {
				glog.Error("Parse CfgHeroExtraData.RarityDiscount field error,value:", arrayRarityDiscount[j])
				return false
			}
			data.RarityDiscount = append(data.RarityDiscount, v)
		}

		/* parse CombatPower field */
		vCombatPower, _ := parse.GetFieldByName(uint32(i), "combatPower")
		var CombatPowerRet bool
		data.CombatPower, CombatPowerRet = String2Int32(vCombatPower)
		if !CombatPowerRet {
			glog.Error("Parse CfgHeroExtraData.CombatPower field error,value:", vCombatPower)
			return false
		}

		/* parse HpMax field */
		vHpMax, _ := parse.GetFieldByName(uint32(i), "hpMax")
		var HpMaxRet bool
		data.HpMax, HpMaxRet = String2Int32(vHpMax)
		if !HpMaxRet {
			glog.Error("Parse CfgHeroExtraData.HpMax field error,value:", vHpMax)
			return false
		}

		/* parse PhyAtk field */
		vPhyAtk, _ := parse.GetFieldByName(uint32(i), "phyAtk")
		var PhyAtkRet bool
		data.PhyAtk, PhyAtkRet = String2Int32(vPhyAtk)
		if !PhyAtkRet {
			glog.Error("Parse CfgHeroExtraData.PhyAtk field error,value:", vPhyAtk)
			return false
		}

		/* parse MagAtk field */
		vMagAtk, _ := parse.GetFieldByName(uint32(i), "magAtk")
		var MagAtkRet bool
		data.MagAtk, MagAtkRet = String2Int32(vMagAtk)
		if !MagAtkRet {
			glog.Error("Parse CfgHeroExtraData.MagAtk field error,value:", vMagAtk)
			return false
		}

		/* parse PhyDfs field */
		vPhyDfs, _ := parse.GetFieldByName(uint32(i), "phyDfs")
		var PhyDfsRet bool
		data.PhyDfs, PhyDfsRet = String2Int32(vPhyDfs)
		if !PhyDfsRet {
			glog.Error("Parse CfgHeroExtraData.PhyDfs field error,value:", vPhyDfs)
			return false
		}

		/* parse MagDfs field */
		vMagDfs, _ := parse.GetFieldByName(uint32(i), "magDfs")
		var MagDfsRet bool
		data.MagDfs, MagDfsRet = String2Int32(vMagDfs)
		if !MagDfsRet {
			glog.Error("Parse CfgHeroExtraData.MagDfs field error,value:", vMagDfs)
			return false
		}

		/* parse CritAtkRatio field */
		vCritAtkRatio, _ := parse.GetFieldByName(uint32(i), "critAtkRatio")
		var CritAtkRatioRet bool
		data.CritAtkRatio, CritAtkRatioRet = String2Int32(vCritAtkRatio)
		if !CritAtkRatioRet {
			glog.Error("Parse CfgHeroExtraData.CritAtkRatio field error,value:", vCritAtkRatio)
			return false
		}

		/* parse CritAtkValue field */
		vCritAtkValue, _ := parse.GetFieldByName(uint32(i), "critAtkValue")
		var CritAtkValueRet bool
		data.CritAtkValue, CritAtkValueRet = String2Int32(vCritAtkValue)
		if !CritAtkValueRet {
			glog.Error("Parse CfgHeroExtraData.CritAtkValue field error,value:", vCritAtkValue)
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

func (c *CfgHeroExtraDataConfig) Clear() {
}

func (c *CfgHeroExtraDataConfig) Find(id int32) (*CfgHeroExtraData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgHeroExtraDataConfig) GetAllData() map[int32]*CfgHeroExtraData {
	return c.data
}

func (c *CfgHeroExtraDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.HeroID, ",", v.HeroLv, ",", v.RarityDiscount, ",", v.CombatPower, ",", v.HpMax, ",", v.PhyAtk, ",", v.MagAtk, ",", v.PhyDfs, ",", v.MagDfs, ",", v.CritAtkRatio, ",", v.CritAtkValue)
	}
}
