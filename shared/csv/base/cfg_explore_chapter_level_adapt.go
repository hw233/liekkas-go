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

type CfgExploreChapterLevelAdapt struct {
	Id            int32
	AdaptionId    int32
	AdaptionLevel int32
	MonsterLevel  int32
	CombatPower   float64
	HpMax         float64
	PhyAtk        float64
	MagAtk        float64
	PhyDfs        float64
	MagDfs        float64
}

type CfgExploreChapterLevelAdaptConfig struct {
	data map[int32]*CfgExploreChapterLevelAdapt
}

func NewCfgExploreChapterLevelAdaptConfig() *CfgExploreChapterLevelAdaptConfig {
	return &CfgExploreChapterLevelAdaptConfig{
		data: make(map[int32]*CfgExploreChapterLevelAdapt),
	}
}

func (c *CfgExploreChapterLevelAdaptConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreChapterLevelAdapt)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.Id field error,value:", vId)
			return false
		}

		/* parse AdaptionId field */
		vAdaptionId, _ := parse.GetFieldByName(uint32(i), "adaptionId")
		var AdaptionIdRet bool
		data.AdaptionId, AdaptionIdRet = String2Int32(vAdaptionId)
		if !AdaptionIdRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.AdaptionId field error,value:", vAdaptionId)
			return false
		}

		/* parse AdaptionLevel field */
		vAdaptionLevel, _ := parse.GetFieldByName(uint32(i), "adaptionLevel")
		var AdaptionLevelRet bool
		data.AdaptionLevel, AdaptionLevelRet = String2Int32(vAdaptionLevel)
		if !AdaptionLevelRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.AdaptionLevel field error,value:", vAdaptionLevel)
			return false
		}

		/* parse MonsterLevel field */
		vMonsterLevel, _ := parse.GetFieldByName(uint32(i), "monsterLevel")
		var MonsterLevelRet bool
		data.MonsterLevel, MonsterLevelRet = String2Int32(vMonsterLevel)
		if !MonsterLevelRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.MonsterLevel field error,value:", vMonsterLevel)
			return false
		}

		/* parse CombatPower field */
		vCombatPower, _ := parse.GetFieldByName(uint32(i), "combatPower")
		var CombatPowerRet bool
		data.CombatPower, CombatPowerRet = String2Float(vCombatPower)
		if !CombatPowerRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.CombatPower field error,value:", vCombatPower)
		}

		/* parse HpMax field */
		vHpMax, _ := parse.GetFieldByName(uint32(i), "hpMax")
		var HpMaxRet bool
		data.HpMax, HpMaxRet = String2Float(vHpMax)
		if !HpMaxRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.HpMax field error,value:", vHpMax)
		}

		/* parse PhyAtk field */
		vPhyAtk, _ := parse.GetFieldByName(uint32(i), "phyAtk")
		var PhyAtkRet bool
		data.PhyAtk, PhyAtkRet = String2Float(vPhyAtk)
		if !PhyAtkRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.PhyAtk field error,value:", vPhyAtk)
		}

		/* parse MagAtk field */
		vMagAtk, _ := parse.GetFieldByName(uint32(i), "magAtk")
		var MagAtkRet bool
		data.MagAtk, MagAtkRet = String2Float(vMagAtk)
		if !MagAtkRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.MagAtk field error,value:", vMagAtk)
		}

		/* parse PhyDfs field */
		vPhyDfs, _ := parse.GetFieldByName(uint32(i), "phyDfs")
		var PhyDfsRet bool
		data.PhyDfs, PhyDfsRet = String2Float(vPhyDfs)
		if !PhyDfsRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.PhyDfs field error,value:", vPhyDfs)
		}

		/* parse MagDfs field */
		vMagDfs, _ := parse.GetFieldByName(uint32(i), "magDfs")
		var MagDfsRet bool
		data.MagDfs, MagDfsRet = String2Float(vMagDfs)
		if !MagDfsRet {
			glog.Error("Parse CfgExploreChapterLevelAdapt.MagDfs field error,value:", vMagDfs)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgExploreChapterLevelAdaptConfig) Clear() {
}

func (c *CfgExploreChapterLevelAdaptConfig) Find(id int32) (*CfgExploreChapterLevelAdapt, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreChapterLevelAdaptConfig) GetAllData() map[int32]*CfgExploreChapterLevelAdapt {
	return c.data
}

func (c *CfgExploreChapterLevelAdaptConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.AdaptionId, ",", v.AdaptionLevel, ",", v.MonsterLevel, ",", v.CombatPower, ",", v.HpMax, ",", v.PhyAtk, ",", v.MagAtk, ",", v.PhyDfs, ",", v.MagDfs)
	}
}
