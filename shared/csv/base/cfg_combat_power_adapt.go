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

type CfgCombatPowerAdapt struct {
	Id                   int32
	CareerId             int32
	Symbol               int32
	HpAdaption           float64
	PhyAtkAdaption       float64
	MagAtkAdaption       float64
	PhyDfsAdaption       float64
	MagDfsAdaption       float64
	CritAtkRatioAdaption float64
}

type CfgCombatPowerAdaptConfig struct {
	data map[int32]*CfgCombatPowerAdapt
}

func NewCfgCombatPowerAdaptConfig() *CfgCombatPowerAdaptConfig {
	return &CfgCombatPowerAdaptConfig{
		data: make(map[int32]*CfgCombatPowerAdapt),
	}
}

func (c *CfgCombatPowerAdaptConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCombatPowerAdapt)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCombatPowerAdapt.Id field error,value:", vId)
			return false
		}

		/* parse CareerId field */
		vCareerId, _ := parse.GetFieldByName(uint32(i), "careerId")
		var CareerIdRet bool
		data.CareerId, CareerIdRet = String2Int32(vCareerId)
		if !CareerIdRet {
			glog.Error("Parse CfgCombatPowerAdapt.CareerId field error,value:", vCareerId)
			return false
		}

		/* parse Symbol field */
		vSymbol, _ := parse.GetFieldByName(uint32(i), "symbol")
		var SymbolRet bool
		data.Symbol, SymbolRet = String2Int32(vSymbol)
		if !SymbolRet {
			glog.Error("Parse CfgCombatPowerAdapt.Symbol field error,value:", vSymbol)
			return false
		}

		/* parse HpAdaption field */
		vHpAdaption, _ := parse.GetFieldByName(uint32(i), "hpAdaption")
		var HpAdaptionRet bool
		data.HpAdaption, HpAdaptionRet = String2Float(vHpAdaption)
		if !HpAdaptionRet {
			glog.Error("Parse CfgCombatPowerAdapt.HpAdaption field error,value:", vHpAdaption)
		}

		/* parse PhyAtkAdaption field */
		vPhyAtkAdaption, _ := parse.GetFieldByName(uint32(i), "phyAtkAdaption")
		var PhyAtkAdaptionRet bool
		data.PhyAtkAdaption, PhyAtkAdaptionRet = String2Float(vPhyAtkAdaption)
		if !PhyAtkAdaptionRet {
			glog.Error("Parse CfgCombatPowerAdapt.PhyAtkAdaption field error,value:", vPhyAtkAdaption)
		}

		/* parse MagAtkAdaption field */
		vMagAtkAdaption, _ := parse.GetFieldByName(uint32(i), "magAtkAdaption")
		var MagAtkAdaptionRet bool
		data.MagAtkAdaption, MagAtkAdaptionRet = String2Float(vMagAtkAdaption)
		if !MagAtkAdaptionRet {
			glog.Error("Parse CfgCombatPowerAdapt.MagAtkAdaption field error,value:", vMagAtkAdaption)
		}

		/* parse PhyDfsAdaption field */
		vPhyDfsAdaption, _ := parse.GetFieldByName(uint32(i), "phyDfsAdaption")
		var PhyDfsAdaptionRet bool
		data.PhyDfsAdaption, PhyDfsAdaptionRet = String2Float(vPhyDfsAdaption)
		if !PhyDfsAdaptionRet {
			glog.Error("Parse CfgCombatPowerAdapt.PhyDfsAdaption field error,value:", vPhyDfsAdaption)
		}

		/* parse MagDfsAdaption field */
		vMagDfsAdaption, _ := parse.GetFieldByName(uint32(i), "magDfsAdaption")
		var MagDfsAdaptionRet bool
		data.MagDfsAdaption, MagDfsAdaptionRet = String2Float(vMagDfsAdaption)
		if !MagDfsAdaptionRet {
			glog.Error("Parse CfgCombatPowerAdapt.MagDfsAdaption field error,value:", vMagDfsAdaption)
		}

		/* parse CritAtkRatioAdaption field */
		vCritAtkRatioAdaption, _ := parse.GetFieldByName(uint32(i), "critAtkRatioAdaption")
		var CritAtkRatioAdaptionRet bool
		data.CritAtkRatioAdaption, CritAtkRatioAdaptionRet = String2Float(vCritAtkRatioAdaption)
		if !CritAtkRatioAdaptionRet {
			glog.Error("Parse CfgCombatPowerAdapt.CritAtkRatioAdaption field error,value:", vCritAtkRatioAdaption)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgCombatPowerAdaptConfig) Clear() {
}

func (c *CfgCombatPowerAdaptConfig) Find(id int32) (*CfgCombatPowerAdapt, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCombatPowerAdaptConfig) GetAllData() map[int32]*CfgCombatPowerAdapt {
	return c.data
}

func (c *CfgCombatPowerAdaptConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CareerId, ",", v.Symbol, ",", v.HpAdaption, ",", v.PhyAtkAdaption, ",", v.MagAtkAdaption, ",", v.PhyDfsAdaption, ",", v.MagDfsAdaption, ",", v.CritAtkRatioAdaption)
	}
}
