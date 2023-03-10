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

type CfgMonsterData struct {
	Id                   int32
	Lv                   int32
	MonsterId            int32
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
	BreakMax             int32
	BreakValue           int32
	BreakTime            int32
	CombatPower          int32
}

type CfgMonsterDataConfig struct {
	data map[int32]*CfgMonsterData
}

func NewCfgMonsterDataConfig() *CfgMonsterDataConfig {
	return &CfgMonsterDataConfig{
		data: make(map[int32]*CfgMonsterData),
	}
}

func (c *CfgMonsterDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgMonsterData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgMonsterData.Id field error,value:", vId)
			return false
		}

		/* parse Lv field */
		vLv, _ := parse.GetFieldByName(uint32(i), "lv")
		var LvRet bool
		data.Lv, LvRet = String2Int32(vLv)
		if !LvRet {
			glog.Error("Parse CfgMonsterData.Lv field error,value:", vLv)
			return false
		}

		/* parse MonsterId field */
		vMonsterId, _ := parse.GetFieldByName(uint32(i), "monsterId")
		var MonsterIdRet bool
		data.MonsterId, MonsterIdRet = String2Int32(vMonsterId)
		if !MonsterIdRet {
			glog.Error("Parse CfgMonsterData.MonsterId field error,value:", vMonsterId)
			return false
		}

		/* parse CharaName field */
		data.CharaName, _ = parse.GetFieldByName(uint32(i), "charaName")

		/* parse HpMax field */
		vHpMax, _ := parse.GetFieldByName(uint32(i), "hpMax")
		var HpMaxRet bool
		data.HpMax, HpMaxRet = String2Int32(vHpMax)
		if !HpMaxRet {
			glog.Error("Parse CfgMonsterData.HpMax field error,value:", vHpMax)
			return false
		}

		/* parse PhyAtk field */
		vPhyAtk, _ := parse.GetFieldByName(uint32(i), "phyAtk")
		var PhyAtkRet bool
		data.PhyAtk, PhyAtkRet = String2Int32(vPhyAtk)
		if !PhyAtkRet {
			glog.Error("Parse CfgMonsterData.PhyAtk field error,value:", vPhyAtk)
			return false
		}

		/* parse MagAtk field */
		vMagAtk, _ := parse.GetFieldByName(uint32(i), "magAtk")
		var MagAtkRet bool
		data.MagAtk, MagAtkRet = String2Int32(vMagAtk)
		if !MagAtkRet {
			glog.Error("Parse CfgMonsterData.MagAtk field error,value:", vMagAtk)
			return false
		}

		/* parse PhyDfs field */
		vPhyDfs, _ := parse.GetFieldByName(uint32(i), "phyDfs")
		var PhyDfsRet bool
		data.PhyDfs, PhyDfsRet = String2Int32(vPhyDfs)
		if !PhyDfsRet {
			glog.Error("Parse CfgMonsterData.PhyDfs field error,value:", vPhyDfs)
			return false
		}

		/* parse MagDfs field */
		vMagDfs, _ := parse.GetFieldByName(uint32(i), "magDfs")
		var MagDfsRet bool
		data.MagDfs, MagDfsRet = String2Int32(vMagDfs)
		if !MagDfsRet {
			glog.Error("Parse CfgMonsterData.MagDfs field error,value:", vMagDfs)
			return false
		}

		/* parse CritAtkRatio field */
		vCritAtkRatio, _ := parse.GetFieldByName(uint32(i), "critAtkRatio")
		var CritAtkRatioRet bool
		data.CritAtkRatio, CritAtkRatioRet = String2Int32(vCritAtkRatio)
		if !CritAtkRatioRet {
			glog.Error("Parse CfgMonsterData.CritAtkRatio field error,value:", vCritAtkRatio)
			return false
		}

		/* parse CritDfsValue field */
		vCritDfsValue, _ := parse.GetFieldByName(uint32(i), "critDfsValue")
		var CritDfsValueRet bool
		data.CritDfsValue, CritDfsValueRet = String2Int32(vCritDfsValue)
		if !CritDfsValueRet {
			glog.Error("Parse CfgMonsterData.CritDfsValue field error,value:", vCritDfsValue)
			return false
		}

		/* parse HitRateBasic field */
		vHitRateBasic, _ := parse.GetFieldByName(uint32(i), "hitRateBasic")
		var HitRateBasicRet bool
		data.HitRateBasic, HitRateBasicRet = String2Int32(vHitRateBasic)
		if !HitRateBasicRet {
			glog.Error("Parse CfgMonsterData.HitRateBasic field error,value:", vHitRateBasic)
			return false
		}

		/* parse HitRateValue field */
		vHitRateValue, _ := parse.GetFieldByName(uint32(i), "hitRateValue")
		var HitRateValueRet bool
		data.HitRateValue, HitRateValueRet = String2Int32(vHitRateValue)
		if !HitRateValueRet {
			glog.Error("Parse CfgMonsterData.HitRateValue field error,value:", vHitRateValue)
			return false
		}

		/* parse EvadeValue field */
		vEvadeValue, _ := parse.GetFieldByName(uint32(i), "evadeValue")
		var EvadeValueRet bool
		data.EvadeValue, EvadeValueRet = String2Int32(vEvadeValue)
		if !EvadeValueRet {
			glog.Error("Parse CfgMonsterData.EvadeValue field error,value:", vEvadeValue)
			return false
		}

		/* parse EnergyMax field */
		vEnergyMax, _ := parse.GetFieldByName(uint32(i), "energyMax")
		var EnergyMaxRet bool
		data.EnergyMax, EnergyMaxRet = String2Int32(vEnergyMax)
		if !EnergyMaxRet {
			glog.Error("Parse CfgMonsterData.EnergyMax field error,value:", vEnergyMax)
			return false
		}

		/* parse EnergyRecover field */
		vEnergyRecover, _ := parse.GetFieldByName(uint32(i), "energyRecover")
		var EnergyRecoverRet bool
		data.EnergyRecover, EnergyRecoverRet = String2Int32(vEnergyRecover)
		if !EnergyRecoverRet {
			glog.Error("Parse CfgMonsterData.EnergyRecover field error,value:", vEnergyRecover)
			return false
		}

		/* parse HeroEnergyRecover field */
		vHeroEnergyRecover, _ := parse.GetFieldByName(uint32(i), "heroEnergyRecover")
		var HeroEnergyRecoverRet bool
		data.HeroEnergyRecover, HeroEnergyRecoverRet = String2Int32(vHeroEnergyRecover)
		if !HeroEnergyRecoverRet {
			glog.Error("Parse CfgMonsterData.HeroEnergyRecover field error,value:", vHeroEnergyRecover)
			return false
		}

		/* parse HpRecover field */
		vHpRecover, _ := parse.GetFieldByName(uint32(i), "hpRecover")
		var HpRecoverRet bool
		data.HpRecover, HpRecoverRet = String2Int32(vHpRecover)
		if !HpRecoverRet {
			glog.Error("Parse CfgMonsterData.HpRecover field error,value:", vHpRecover)
			return false
		}

		/* parse Energy field */
		vEnergy, _ := parse.GetFieldByName(uint32(i), "energy")
		var EnergyRet bool
		data.Energy, EnergyRet = String2Int32(vEnergy)
		if !EnergyRet {
			glog.Error("Parse CfgMonsterData.Energy field error,value:", vEnergy)
			return false
		}

		/* parse EnergyDamagedRecover field */
		vEnergyDamagedRecover, _ := parse.GetFieldByName(uint32(i), "energyDamagedRecover")
		var EnergyDamagedRecoverRet bool
		data.EnergyDamagedRecover, EnergyDamagedRecoverRet = String2Int32(vEnergyDamagedRecover)
		if !EnergyDamagedRecoverRet {
			glog.Error("Parse CfgMonsterData.EnergyDamagedRecover field error,value:", vEnergyDamagedRecover)
			return false
		}

		/* parse MoveSpd field */
		vMoveSpd, _ := parse.GetFieldByName(uint32(i), "moveSpd")
		var MoveSpdRet bool
		data.MoveSpd, MoveSpdRet = String2Int32(vMoveSpd)
		if !MoveSpdRet {
			glog.Error("Parse CfgMonsterData.MoveSpd field error,value:", vMoveSpd)
			return false
		}

		/* parse BreakMax field */
		vBreakMax, _ := parse.GetFieldByName(uint32(i), "breakMax")
		var BreakMaxRet bool
		data.BreakMax, BreakMaxRet = String2Int32(vBreakMax)
		if !BreakMaxRet {
			glog.Error("Parse CfgMonsterData.BreakMax field error,value:", vBreakMax)
			return false
		}

		/* parse BreakValue field */
		vBreakValue, _ := parse.GetFieldByName(uint32(i), "breakValue")
		var BreakValueRet bool
		data.BreakValue, BreakValueRet = String2Int32(vBreakValue)
		if !BreakValueRet {
			glog.Error("Parse CfgMonsterData.BreakValue field error,value:", vBreakValue)
			return false
		}

		/* parse BreakTime field */
		vBreakTime, _ := parse.GetFieldByName(uint32(i), "breakTime")
		var BreakTimeRet bool
		data.BreakTime, BreakTimeRet = String2Int32(vBreakTime)
		if !BreakTimeRet {
			glog.Error("Parse CfgMonsterData.BreakTime field error,value:", vBreakTime)
			return false
		}

		/* parse CombatPower field */
		vCombatPower, _ := parse.GetFieldByName(uint32(i), "combatPower")
		var CombatPowerRet bool
		data.CombatPower, CombatPowerRet = String2Int32(vCombatPower)
		if !CombatPowerRet {
			glog.Error("Parse CfgMonsterData.CombatPower field error,value:", vCombatPower)
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

func (c *CfgMonsterDataConfig) Clear() {
}

func (c *CfgMonsterDataConfig) Find(id int32) (*CfgMonsterData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgMonsterDataConfig) GetAllData() map[int32]*CfgMonsterData {
	return c.data
}

func (c *CfgMonsterDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Lv, ",", v.MonsterId, ",", v.CharaName, ",", v.HpMax, ",", v.PhyAtk, ",", v.MagAtk, ",", v.PhyDfs, ",", v.MagDfs, ",", v.CritAtkRatio, ",", v.CritDfsValue, ",", v.HitRateBasic, ",", v.HitRateValue, ",", v.EvadeValue, ",", v.EnergyMax, ",", v.EnergyRecover, ",", v.HeroEnergyRecover, ",", v.HpRecover, ",", v.Energy, ",", v.EnergyDamagedRecover, ",", v.MoveSpd, ",", v.BreakMax, ",", v.BreakValue, ",", v.BreakTime, ",", v.CombatPower)
	}
}
