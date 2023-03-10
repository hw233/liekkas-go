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

type CfgTeamLevel struct {
	Id                  int32
	LvExp               int32
	CharLv              int32
	MaxStamina          int32
	StaminaRecover      int32
	MaxAp               int32
	ExploreAp           int32
	HeroMaxLevel        int32
	EquipMaxLevel       int32
	WorldItemLevelUp    int32
	BattleAdaptionLevel int32
}

type CfgTeamLevelConfig struct {
	data map[int32]*CfgTeamLevel
}

func NewCfgTeamLevelConfig() *CfgTeamLevelConfig {
	return &CfgTeamLevelConfig{
		data: make(map[int32]*CfgTeamLevel),
	}
}

func (c *CfgTeamLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgTeamLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgTeamLevel.Id field error,value:", vId)
			return false
		}

		/* parse LvExp field */
		vLvExp, _ := parse.GetFieldByName(uint32(i), "lv_exp")
		var LvExpRet bool
		data.LvExp, LvExpRet = String2Int32(vLvExp)
		if !LvExpRet {
			glog.Error("Parse CfgTeamLevel.LvExp field error,value:", vLvExp)
			return false
		}

		/* parse CharLv field */
		vCharLv, _ := parse.GetFieldByName(uint32(i), "char_lv")
		var CharLvRet bool
		data.CharLv, CharLvRet = String2Int32(vCharLv)
		if !CharLvRet {
			glog.Error("Parse CfgTeamLevel.CharLv field error,value:", vCharLv)
			return false
		}

		/* parse MaxStamina field */
		vMaxStamina, _ := parse.GetFieldByName(uint32(i), "max_stamina")
		var MaxStaminaRet bool
		data.MaxStamina, MaxStaminaRet = String2Int32(vMaxStamina)
		if !MaxStaminaRet {
			glog.Error("Parse CfgTeamLevel.MaxStamina field error,value:", vMaxStamina)
			return false
		}

		/* parse StaminaRecover field */
		vStaminaRecover, _ := parse.GetFieldByName(uint32(i), "stamina_recover")
		var StaminaRecoverRet bool
		data.StaminaRecover, StaminaRecoverRet = String2Int32(vStaminaRecover)
		if !StaminaRecoverRet {
			glog.Error("Parse CfgTeamLevel.StaminaRecover field error,value:", vStaminaRecover)
			return false
		}

		/* parse MaxAp field */
		vMaxAp, _ := parse.GetFieldByName(uint32(i), "max_ap")
		var MaxApRet bool
		data.MaxAp, MaxApRet = String2Int32(vMaxAp)
		if !MaxApRet {
			glog.Error("Parse CfgTeamLevel.MaxAp field error,value:", vMaxAp)
			return false
		}

		/* parse ExploreAp field */
		vExploreAp, _ := parse.GetFieldByName(uint32(i), "explore_ap")
		var ExploreApRet bool
		data.ExploreAp, ExploreApRet = String2Int32(vExploreAp)
		if !ExploreApRet {
			glog.Error("Parse CfgTeamLevel.ExploreAp field error,value:", vExploreAp)
			return false
		}

		/* parse HeroMaxLevel field */
		vHeroMaxLevel, _ := parse.GetFieldByName(uint32(i), "hero_max_level")
		var HeroMaxLevelRet bool
		data.HeroMaxLevel, HeroMaxLevelRet = String2Int32(vHeroMaxLevel)
		if !HeroMaxLevelRet {
			glog.Error("Parse CfgTeamLevel.HeroMaxLevel field error,value:", vHeroMaxLevel)
			return false
		}

		/* parse EquipMaxLevel field */
		vEquipMaxLevel, _ := parse.GetFieldByName(uint32(i), "equip_max_level")
		var EquipMaxLevelRet bool
		data.EquipMaxLevel, EquipMaxLevelRet = String2Int32(vEquipMaxLevel)
		if !EquipMaxLevelRet {
			glog.Error("Parse CfgTeamLevel.EquipMaxLevel field error,value:", vEquipMaxLevel)
			return false
		}

		/* parse WorldItemLevelUp field */
		vWorldItemLevelUp, _ := parse.GetFieldByName(uint32(i), "world_item_level_up")
		var WorldItemLevelUpRet bool
		data.WorldItemLevelUp, WorldItemLevelUpRet = String2Int32(vWorldItemLevelUp)
		if !WorldItemLevelUpRet {
			glog.Error("Parse CfgTeamLevel.WorldItemLevelUp field error,value:", vWorldItemLevelUp)
			return false
		}

		/* parse BattleAdaptionLevel field */
		vBattleAdaptionLevel, _ := parse.GetFieldByName(uint32(i), "battle_adaption_level")
		var BattleAdaptionLevelRet bool
		data.BattleAdaptionLevel, BattleAdaptionLevelRet = String2Int32(vBattleAdaptionLevel)
		if !BattleAdaptionLevelRet {
			glog.Error("Parse CfgTeamLevel.BattleAdaptionLevel field error,value:", vBattleAdaptionLevel)
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

func (c *CfgTeamLevelConfig) Clear() {
}

func (c *CfgTeamLevelConfig) Find(id int32) (*CfgTeamLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgTeamLevelConfig) GetAllData() map[int32]*CfgTeamLevel {
	return c.data
}

func (c *CfgTeamLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.LvExp, ",", v.CharLv, ",", v.MaxStamina, ",", v.StaminaRecover, ",", v.MaxAp, ",", v.ExploreAp, ",", v.HeroMaxLevel, ",", v.EquipMaxLevel, ",", v.WorldItemLevelUp, ",", v.BattleAdaptionLevel)
	}
}
