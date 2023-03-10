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

type CfgCharacterData struct {
	Id                   int32
	CharaName            string
	HpMax                int32
	PhyAtk               int32
	MagAtk               int32
	PhyDfs               int32
	MagDfs               int32
	CritAtkRatio         int32
	CritDfsValue         int32
	HitRateBasic         int32
	HitRateValue         int32
	EvadeValue           int32
	EnergyMax            int32
	EnergyRecover        int32
	HeroEnergyRecover    int32
	HpRecover            int32
	Energy               int32
	EnergyDamagedRecover int32
	MoveSpd              int32
	CombatPower          int32
}

type CfgCharacterDataConfig struct {
	data map[int32]*CfgCharacterData
}

func NewCfgCharacterDataConfig() *CfgCharacterDataConfig {
	return &CfgCharacterDataConfig{
		data: make(map[int32]*CfgCharacterData),
	}
}

func (c *CfgCharacterDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterData.Id field error,value:", vId)
			return false
		}

		/* parse CharaName field */
		data.CharaName, _ = parse.GetFieldByName(uint32(i), "charaName")

		/* parse HpMax field */
		vHpMax, _ := parse.GetFieldByName(uint32(i), "hpMax")
		var HpMaxRet bool
		data.HpMax, HpMaxRet = String2Int32(vHpMax)
		if !HpMaxRet {
			glog.Error("Parse CfgCharacterData.HpMax field error,value:", vHpMax)
			return false
		}

		/* parse PhyAtk field */
		vPhyAtk, _ := parse.GetFieldByName(uint32(i), "phyAtk")
		var PhyAtkRet bool
		data.PhyAtk, PhyAtkRet = String2Int32(vPhyAtk)
		if !PhyAtkRet {
			glog.Error("Parse CfgCharacterData.PhyAtk field error,value:", vPhyAtk)
			return false
		}

		/* parse MagAtk field */
		vMagAtk, _ := parse.GetFieldByName(uint32(i), "magAtk")
		var MagAtkRet bool
		data.MagAtk, MagAtkRet = String2Int32(vMagAtk)
		if !MagAtkRet {
			glog.Error("Parse CfgCharacterData.MagAtk field error,value:", vMagAtk)
			return false
		}

		/* parse PhyDfs field */
		vPhyDfs, _ := parse.GetFieldByName(uint32(i), "phyDfs")
		var PhyDfsRet bool
		data.PhyDfs, PhyDfsRet = String2Int32(vPhyDfs)
		if !PhyDfsRet {
			glog.Error("Parse CfgCharacterData.PhyDfs field error,value:", vPhyDfs)
			return false
		}

		/* parse MagDfs field */
		vMagDfs, _ := parse.GetFieldByName(uint32(i), "magDfs")
		var MagDfsRet bool
		data.MagDfs, MagDfsRet = String2Int32(vMagDfs)
		if !MagDfsRet {
			glog.Error("Parse CfgCharacterData.MagDfs field error,value:", vMagDfs)
			return false
		}

		/* parse CritAtkRatio field */
		vCritAtkRatio, _ := parse.GetFieldByName(uint32(i), "critAtkRatio")
		var CritAtkRatioRet bool
		data.CritAtkRatio, CritAtkRatioRet = String2Int32(vCritAtkRatio)
		if !CritAtkRatioRet {
			glog.Error("Parse CfgCharacterData.CritAtkRatio field error,value:", vCritAtkRatio)
			return false
		}

		/* parse CritDfsValue field */
		vCritDfsValue, _ := parse.GetFieldByName(uint32(i), "critDfsValue")
		var CritDfsValueRet bool
		data.CritDfsValue, CritDfsValueRet = String2Int32(vCritDfsValue)
		if !CritDfsValueRet {
			glog.Error("Parse CfgCharacterData.CritDfsValue field error,value:", vCritDfsValue)
			return false
		}

		/* parse HitRateBasic field */
		vHitRateBasic, _ := parse.GetFieldByName(uint32(i), "hitRateBasic")
		var HitRateBasicRet bool
		data.HitRateBasic, HitRateBasicRet = String2Int32(vHitRateBasic)
		if !HitRateBasicRet {
			glog.Error("Parse CfgCharacterData.HitRateBasic field error,value:", vHitRateBasic)
			return false
		}

		/* parse HitRateValue field */
		vHitRateValue, _ := parse.GetFieldByName(uint32(i), "hitRateValue")
		var HitRateValueRet bool
		data.HitRateValue, HitRateValueRet = String2Int32(vHitRateValue)
		if !HitRateValueRet {
			glog.Error("Parse CfgCharacterData.HitRateValue field error,value:", vHitRateValue)
			return false
		}

		/* parse EvadeValue field */
		vEvadeValue, _ := parse.GetFieldByName(uint32(i), "evadeValue")
		var EvadeValueRet bool
		data.EvadeValue, EvadeValueRet = String2Int32(vEvadeValue)
		if !EvadeValueRet {
			glog.Error("Parse CfgCharacterData.EvadeValue field error,value:", vEvadeValue)
			return false
		}

		/* parse EnergyMax field */
		vEnergyMax, _ := parse.GetFieldByName(uint32(i), "energyMax")
		var EnergyMaxRet bool
		data.EnergyMax, EnergyMaxRet = String2Int32(vEnergyMax)
		if !EnergyMaxRet {
			glog.Error("Parse CfgCharacterData.EnergyMax field error,value:", vEnergyMax)
			return false
		}

		/* parse EnergyRecover field */
		vEnergyRecover, _ := parse.GetFieldByName(uint32(i), "energyRecover")
		var EnergyRecoverRet bool
		data.EnergyRecover, EnergyRecoverRet = String2Int32(vEnergyRecover)
		if !EnergyRecoverRet {
			glog.Error("Parse CfgCharacterData.EnergyRecover field error,value:", vEnergyRecover)
			return false
		}

		/* parse HeroEnergyRecover field */
		vHeroEnergyRecover, _ := parse.GetFieldByName(uint32(i), "heroEnergyRecover")
		var HeroEnergyRecoverRet bool
		data.HeroEnergyRecover, HeroEnergyRecoverRet = String2Int32(vHeroEnergyRecover)
		if !HeroEnergyRecoverRet {
			glog.Error("Parse CfgCharacterData.HeroEnergyRecover field error,value:", vHeroEnergyRecover)
			return false
		}

		/* parse HpRecover field */
		vHpRecover, _ := parse.GetFieldByName(uint32(i), "hpRecover")
		var HpRecoverRet bool
		data.HpRecover, HpRecoverRet = String2Int32(vHpRecover)
		if !HpRecoverRet {
			glog.Error("Parse CfgCharacterData.HpRecover field error,value:", vHpRecover)
			return false
		}

		/* parse Energy field */
		vEnergy, _ := parse.GetFieldByName(uint32(i), "energy")
		var EnergyRet bool
		data.Energy, EnergyRet = String2Int32(vEnergy)
		if !EnergyRet {
			glog.Error("Parse CfgCharacterData.Energy field error,value:", vEnergy)
			return false
		}

		/* parse EnergyDamagedRecover field */
		vEnergyDamagedRecover, _ := parse.GetFieldByName(uint32(i), "energyDamagedRecover")
		var EnergyDamagedRecoverRet bool
		data.EnergyDamagedRecover, EnergyDamagedRecoverRet = String2Int32(vEnergyDamagedRecover)
		if !EnergyDamagedRecoverRet {
			glog.Error("Parse CfgCharacterData.EnergyDamagedRecover field error,value:", vEnergyDamagedRecover)
			return false
		}

		/* parse MoveSpd field */
		vMoveSpd, _ := parse.GetFieldByName(uint32(i), "moveSpd")
		var MoveSpdRet bool
		data.MoveSpd, MoveSpdRet = String2Int32(vMoveSpd)
		if !MoveSpdRet {
			glog.Error("Parse CfgCharacterData.MoveSpd field error,value:", vMoveSpd)
			return false
		}

		/* parse CombatPower field */
		vCombatPower, _ := parse.GetFieldByName(uint32(i), "combatPower")
		var CombatPowerRet bool
		data.CombatPower, CombatPowerRet = String2Int32(vCombatPower)
		if !CombatPowerRet {
			glog.Error("Parse CfgCharacterData.CombatPower field error,value:", vCombatPower)
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

func (c *CfgCharacterDataConfig) Clear() {
}

func (c *CfgCharacterDataConfig) Find(id int32) (*CfgCharacterData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterDataConfig) GetAllData() map[int32]*CfgCharacterData {
	return c.data
}

func (c *CfgCharacterDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CharaName, ",", v.HpMax, ",", v.PhyAtk, ",", v.MagAtk, ",", v.PhyDfs, ",", v.MagDfs, ",", v.CritAtkRatio, ",", v.CritDfsValue, ",", v.HitRateBasic, ",", v.HitRateValue, ",", v.EvadeValue, ",", v.EnergyMax, ",", v.EnergyRecover, ",", v.HeroEnergyRecover, ",", v.HpRecover, ",", v.Energy, ",", v.EnergyDamagedRecover, ",", v.MoveSpd, ",", v.CombatPower)
	}
}
