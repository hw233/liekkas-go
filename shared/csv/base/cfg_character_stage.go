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

type CfgCharacterStage struct {
	Id            int32
	CharId        int32
	Stage         int32
	LevelLimit    int32
	Frame         int32
	HpMax         int32
	PhyAtk        int32
	MagAtk        int32
	PhyDfs        int32
	MagDfs        int32
	CritAtkRatio  int32
	CritDfsRatio  int32
	CritAtkValue  int32
	CritDfsValue  int32
	HitRateValue  int32
	EvadeValue    int32
	NormalAtkUp   int32
	NormalDfsUp   int32
	SkillAtkUp    int32
	SkillDfsUp    int32
	UltraAtkUp    int32
	UltraDfsUp    int32
	SkillPhyAtkUp int32
	SkillPhyDfsUp int32
	SkillMagAtkUp int32
	SkillMagDfsUp int32
	Cost          []string
}

type CfgCharacterStageConfig struct {
	data map[int32]*CfgCharacterStage
}

func NewCfgCharacterStageConfig() *CfgCharacterStageConfig {
	return &CfgCharacterStageConfig{
		data: make(map[int32]*CfgCharacterStage),
	}
}

func (c *CfgCharacterStageConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterStage)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterStage.Id field error,value:", vId)
			return false
		}

		/* parse CharId field */
		vCharId, _ := parse.GetFieldByName(uint32(i), "charId")
		var CharIdRet bool
		data.CharId, CharIdRet = String2Int32(vCharId)
		if !CharIdRet {
			glog.Error("Parse CfgCharacterStage.CharId field error,value:", vCharId)
			return false
		}

		/* parse Stage field */
		vStage, _ := parse.GetFieldByName(uint32(i), "stage")
		var StageRet bool
		data.Stage, StageRet = String2Int32(vStage)
		if !StageRet {
			glog.Error("Parse CfgCharacterStage.Stage field error,value:", vStage)
			return false
		}

		/* parse LevelLimit field */
		vLevelLimit, _ := parse.GetFieldByName(uint32(i), "levelLimit")
		var LevelLimitRet bool
		data.LevelLimit, LevelLimitRet = String2Int32(vLevelLimit)
		if !LevelLimitRet {
			glog.Error("Parse CfgCharacterStage.LevelLimit field error,value:", vLevelLimit)
			return false
		}

		/* parse Frame field */
		vFrame, _ := parse.GetFieldByName(uint32(i), "frame")
		var FrameRet bool
		data.Frame, FrameRet = String2Int32(vFrame)
		if !FrameRet {
			glog.Error("Parse CfgCharacterStage.Frame field error,value:", vFrame)
			return false
		}

		/* parse HpMax field */
		vHpMax, _ := parse.GetFieldByName(uint32(i), "hpMax")
		var HpMaxRet bool
		data.HpMax, HpMaxRet = String2Int32(vHpMax)
		if !HpMaxRet {
			glog.Error("Parse CfgCharacterStage.HpMax field error,value:", vHpMax)
			return false
		}

		/* parse PhyAtk field */
		vPhyAtk, _ := parse.GetFieldByName(uint32(i), "phyAtk")
		var PhyAtkRet bool
		data.PhyAtk, PhyAtkRet = String2Int32(vPhyAtk)
		if !PhyAtkRet {
			glog.Error("Parse CfgCharacterStage.PhyAtk field error,value:", vPhyAtk)
			return false
		}

		/* parse MagAtk field */
		vMagAtk, _ := parse.GetFieldByName(uint32(i), "magAtk")
		var MagAtkRet bool
		data.MagAtk, MagAtkRet = String2Int32(vMagAtk)
		if !MagAtkRet {
			glog.Error("Parse CfgCharacterStage.MagAtk field error,value:", vMagAtk)
			return false
		}

		/* parse PhyDfs field */
		vPhyDfs, _ := parse.GetFieldByName(uint32(i), "phyDfs")
		var PhyDfsRet bool
		data.PhyDfs, PhyDfsRet = String2Int32(vPhyDfs)
		if !PhyDfsRet {
			glog.Error("Parse CfgCharacterStage.PhyDfs field error,value:", vPhyDfs)
			return false
		}

		/* parse MagDfs field */
		vMagDfs, _ := parse.GetFieldByName(uint32(i), "magDfs")
		var MagDfsRet bool
		data.MagDfs, MagDfsRet = String2Int32(vMagDfs)
		if !MagDfsRet {
			glog.Error("Parse CfgCharacterStage.MagDfs field error,value:", vMagDfs)
			return false
		}

		/* parse CritAtkRatio field */
		vCritAtkRatio, _ := parse.GetFieldByName(uint32(i), "critAtkRatio")
		var CritAtkRatioRet bool
		data.CritAtkRatio, CritAtkRatioRet = String2Int32(vCritAtkRatio)
		if !CritAtkRatioRet {
			glog.Error("Parse CfgCharacterStage.CritAtkRatio field error,value:", vCritAtkRatio)
			return false
		}

		/* parse CritDfsRatio field */
		vCritDfsRatio, _ := parse.GetFieldByName(uint32(i), "critDfsRatio")
		var CritDfsRatioRet bool
		data.CritDfsRatio, CritDfsRatioRet = String2Int32(vCritDfsRatio)
		if !CritDfsRatioRet {
			glog.Error("Parse CfgCharacterStage.CritDfsRatio field error,value:", vCritDfsRatio)
			return false
		}

		/* parse CritAtkValue field */
		vCritAtkValue, _ := parse.GetFieldByName(uint32(i), "critAtkValue")
		var CritAtkValueRet bool
		data.CritAtkValue, CritAtkValueRet = String2Int32(vCritAtkValue)
		if !CritAtkValueRet {
			glog.Error("Parse CfgCharacterStage.CritAtkValue field error,value:", vCritAtkValue)
			return false
		}

		/* parse CritDfsValue field */
		vCritDfsValue, _ := parse.GetFieldByName(uint32(i), "critDfsValue")
		var CritDfsValueRet bool
		data.CritDfsValue, CritDfsValueRet = String2Int32(vCritDfsValue)
		if !CritDfsValueRet {
			glog.Error("Parse CfgCharacterStage.CritDfsValue field error,value:", vCritDfsValue)
			return false
		}

		/* parse HitRateValue field */
		vHitRateValue, _ := parse.GetFieldByName(uint32(i), "hitRateValue")
		var HitRateValueRet bool
		data.HitRateValue, HitRateValueRet = String2Int32(vHitRateValue)
		if !HitRateValueRet {
			glog.Error("Parse CfgCharacterStage.HitRateValue field error,value:", vHitRateValue)
			return false
		}

		/* parse EvadeValue field */
		vEvadeValue, _ := parse.GetFieldByName(uint32(i), "evadeValue")
		var EvadeValueRet bool
		data.EvadeValue, EvadeValueRet = String2Int32(vEvadeValue)
		if !EvadeValueRet {
			glog.Error("Parse CfgCharacterStage.EvadeValue field error,value:", vEvadeValue)
			return false
		}

		/* parse NormalAtkUp field */
		vNormalAtkUp, _ := parse.GetFieldByName(uint32(i), "normalAtkUp")
		var NormalAtkUpRet bool
		data.NormalAtkUp, NormalAtkUpRet = String2Int32(vNormalAtkUp)
		if !NormalAtkUpRet {
			glog.Error("Parse CfgCharacterStage.NormalAtkUp field error,value:", vNormalAtkUp)
			return false
		}

		/* parse NormalDfsUp field */
		vNormalDfsUp, _ := parse.GetFieldByName(uint32(i), "normalDfsUp")
		var NormalDfsUpRet bool
		data.NormalDfsUp, NormalDfsUpRet = String2Int32(vNormalDfsUp)
		if !NormalDfsUpRet {
			glog.Error("Parse CfgCharacterStage.NormalDfsUp field error,value:", vNormalDfsUp)
			return false
		}

		/* parse SkillAtkUp field */
		vSkillAtkUp, _ := parse.GetFieldByName(uint32(i), "skillAtkUp")
		var SkillAtkUpRet bool
		data.SkillAtkUp, SkillAtkUpRet = String2Int32(vSkillAtkUp)
		if !SkillAtkUpRet {
			glog.Error("Parse CfgCharacterStage.SkillAtkUp field error,value:", vSkillAtkUp)
			return false
		}

		/* parse SkillDfsUp field */
		vSkillDfsUp, _ := parse.GetFieldByName(uint32(i), "skillDfsUp")
		var SkillDfsUpRet bool
		data.SkillDfsUp, SkillDfsUpRet = String2Int32(vSkillDfsUp)
		if !SkillDfsUpRet {
			glog.Error("Parse CfgCharacterStage.SkillDfsUp field error,value:", vSkillDfsUp)
			return false
		}

		/* parse UltraAtkUp field */
		vUltraAtkUp, _ := parse.GetFieldByName(uint32(i), "ultraAtkUp")
		var UltraAtkUpRet bool
		data.UltraAtkUp, UltraAtkUpRet = String2Int32(vUltraAtkUp)
		if !UltraAtkUpRet {
			glog.Error("Parse CfgCharacterStage.UltraAtkUp field error,value:", vUltraAtkUp)
			return false
		}

		/* parse UltraDfsUp field */
		vUltraDfsUp, _ := parse.GetFieldByName(uint32(i), "ultraDfsUp")
		var UltraDfsUpRet bool
		data.UltraDfsUp, UltraDfsUpRet = String2Int32(vUltraDfsUp)
		if !UltraDfsUpRet {
			glog.Error("Parse CfgCharacterStage.UltraDfsUp field error,value:", vUltraDfsUp)
			return false
		}

		/* parse SkillPhyAtkUp field */
		vSkillPhyAtkUp, _ := parse.GetFieldByName(uint32(i), "skillPhyAtkUp")
		var SkillPhyAtkUpRet bool
		data.SkillPhyAtkUp, SkillPhyAtkUpRet = String2Int32(vSkillPhyAtkUp)
		if !SkillPhyAtkUpRet {
			glog.Error("Parse CfgCharacterStage.SkillPhyAtkUp field error,value:", vSkillPhyAtkUp)
			return false
		}

		/* parse SkillPhyDfsUp field */
		vSkillPhyDfsUp, _ := parse.GetFieldByName(uint32(i), "skillPhyDfsUp")
		var SkillPhyDfsUpRet bool
		data.SkillPhyDfsUp, SkillPhyDfsUpRet = String2Int32(vSkillPhyDfsUp)
		if !SkillPhyDfsUpRet {
			glog.Error("Parse CfgCharacterStage.SkillPhyDfsUp field error,value:", vSkillPhyDfsUp)
			return false
		}

		/* parse SkillMagAtkUp field */
		vSkillMagAtkUp, _ := parse.GetFieldByName(uint32(i), "skillMagAtkUp")
		var SkillMagAtkUpRet bool
		data.SkillMagAtkUp, SkillMagAtkUpRet = String2Int32(vSkillMagAtkUp)
		if !SkillMagAtkUpRet {
			glog.Error("Parse CfgCharacterStage.SkillMagAtkUp field error,value:", vSkillMagAtkUp)
			return false
		}

		/* parse SkillMagDfsUp field */
		vSkillMagDfsUp, _ := parse.GetFieldByName(uint32(i), "skillMagDfsUp")
		var SkillMagDfsUpRet bool
		data.SkillMagDfsUp, SkillMagDfsUpRet = String2Int32(vSkillMagDfsUp)
		if !SkillMagDfsUpRet {
			glog.Error("Parse CfgCharacterStage.SkillMagDfsUp field error,value:", vSkillMagDfsUp)
			return false
		}

		/* parse Cost field */
		vecCost, _ := parse.GetFieldByName(uint32(i), "cost")
		arrayCost := strings.Split(vecCost, ",")
		for j := 0; j < len(arrayCost); j++ {
			v := arrayCost[j]
			data.Cost = append(data.Cost, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgCharacterStageConfig) Clear() {
}

func (c *CfgCharacterStageConfig) Find(id int32) (*CfgCharacterStage, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterStageConfig) GetAllData() map[int32]*CfgCharacterStage {
	return c.data
}

func (c *CfgCharacterStageConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CharId, ",", v.Stage, ",", v.LevelLimit, ",", v.Frame, ",", v.HpMax, ",", v.PhyAtk, ",", v.MagAtk, ",", v.PhyDfs, ",", v.MagDfs, ",", v.CritAtkRatio, ",", v.CritDfsRatio, ",", v.CritAtkValue, ",", v.CritDfsValue, ",", v.HitRateValue, ",", v.EvadeValue, ",", v.NormalAtkUp, ",", v.NormalDfsUp, ",", v.SkillAtkUp, ",", v.SkillDfsUp, ",", v.UltraAtkUp, ",", v.UltraDfsUp, ",", v.SkillPhyAtkUp, ",", v.SkillPhyDfsUp, ",", v.SkillMagAtkUp, ",", v.SkillMagDfsUp, ",", v.Cost)
	}
}
