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

type CfgMonsterBattle struct {
	Id                  int32
	CharaName           string
	AiName              string
	Attackinterval      int32
	NormalAtkID         int32
	NormalMeleeAtkID    int32
	MoveAtkID           int32
	EnterSkillID        int32
	WarningSkill        int32
	UiSkillIDs          []int32
	SkillID             []int32
	UltraSkillID        int32
	AcitvePassiveSkills []int32
	CollideRect         []float64
	MeleeHitRange       float64
	RemoteHitRange      float64
	KnockBackDuration   float64
	GetUpDuration       float64
	DamagedDuration     float64
	EnterDuration       float64
	BkStartDuration     float64
	BkEndDuration       float64
	EnterDelay          float64
	TranslateData       string
}

type CfgMonsterBattleConfig struct {
	data map[int32]*CfgMonsterBattle
}

func NewCfgMonsterBattleConfig() *CfgMonsterBattleConfig {
	return &CfgMonsterBattleConfig{
		data: make(map[int32]*CfgMonsterBattle),
	}
}

func (c *CfgMonsterBattleConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgMonsterBattle)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgMonsterBattle.Id field error,value:", vId)
			return false
		}

		/* parse CharaName field */
		data.CharaName, _ = parse.GetFieldByName(uint32(i), "charaName")

		/* parse AiName field */
		data.AiName, _ = parse.GetFieldByName(uint32(i), "aiName")

		/* parse Attackinterval field */
		vAttackinterval, _ := parse.GetFieldByName(uint32(i), "attackinterval")
		var AttackintervalRet bool
		data.Attackinterval, AttackintervalRet = String2Int32(vAttackinterval)
		if !AttackintervalRet {
			glog.Error("Parse CfgMonsterBattle.Attackinterval field error,value:", vAttackinterval)
			return false
		}

		/* parse NormalAtkID field */
		vNormalAtkID, _ := parse.GetFieldByName(uint32(i), "normalAtkID")
		var NormalAtkIDRet bool
		data.NormalAtkID, NormalAtkIDRet = String2Int32(vNormalAtkID)
		if !NormalAtkIDRet {
			glog.Error("Parse CfgMonsterBattle.NormalAtkID field error,value:", vNormalAtkID)
			return false
		}

		/* parse NormalMeleeAtkID field */
		vNormalMeleeAtkID, _ := parse.GetFieldByName(uint32(i), "normalMeleeAtkID")
		var NormalMeleeAtkIDRet bool
		data.NormalMeleeAtkID, NormalMeleeAtkIDRet = String2Int32(vNormalMeleeAtkID)
		if !NormalMeleeAtkIDRet {
			glog.Error("Parse CfgMonsterBattle.NormalMeleeAtkID field error,value:", vNormalMeleeAtkID)
			return false
		}

		/* parse MoveAtkID field */
		vMoveAtkID, _ := parse.GetFieldByName(uint32(i), "moveAtkID")
		var MoveAtkIDRet bool
		data.MoveAtkID, MoveAtkIDRet = String2Int32(vMoveAtkID)
		if !MoveAtkIDRet {
			glog.Error("Parse CfgMonsterBattle.MoveAtkID field error,value:", vMoveAtkID)
			return false
		}

		/* parse EnterSkillID field */
		vEnterSkillID, _ := parse.GetFieldByName(uint32(i), "enterSkillID")
		var EnterSkillIDRet bool
		data.EnterSkillID, EnterSkillIDRet = String2Int32(vEnterSkillID)
		if !EnterSkillIDRet {
			glog.Error("Parse CfgMonsterBattle.EnterSkillID field error,value:", vEnterSkillID)
			return false
		}

		/* parse WarningSkill field */
		vWarningSkill, _ := parse.GetFieldByName(uint32(i), "warningSkill")
		var WarningSkillRet bool
		data.WarningSkill, WarningSkillRet = String2Int32(vWarningSkill)
		if !WarningSkillRet {
			glog.Error("Parse CfgMonsterBattle.WarningSkill field error,value:", vWarningSkill)
			return false
		}

		/* parse UiSkillIDs field */
		vecUiSkillIDs, _ := parse.GetFieldByName(uint32(i), "uiSkillIDs")
		if vecUiSkillIDs != "" {
			arrayUiSkillIDs := strings.Split(vecUiSkillIDs, ",")
			for j := 0; j < len(arrayUiSkillIDs); j++ {
				v, ret := String2Int32(arrayUiSkillIDs[j])
				if !ret {
					glog.Error("Parse CfgMonsterBattle.UiSkillIDs field error, value:", arrayUiSkillIDs[j])
					return false
				}
				data.UiSkillIDs = append(data.UiSkillIDs, v)
			}
		}

		/* parse SkillID field */
		vecSkillID, _ := parse.GetFieldByName(uint32(i), "skillID")
		if vecSkillID != "" {
			arraySkillID := strings.Split(vecSkillID, ",")
			for j := 0; j < len(arraySkillID); j++ {
				v, ret := String2Int32(arraySkillID[j])
				if !ret {
					glog.Error("Parse CfgMonsterBattle.SkillID field error, value:", arraySkillID[j])
					return false
				}
				data.SkillID = append(data.SkillID, v)
			}
		}

		/* parse UltraSkillID field */
		vUltraSkillID, _ := parse.GetFieldByName(uint32(i), "ultraSkillID")
		var UltraSkillIDRet bool
		data.UltraSkillID, UltraSkillIDRet = String2Int32(vUltraSkillID)
		if !UltraSkillIDRet {
			glog.Error("Parse CfgMonsterBattle.UltraSkillID field error,value:", vUltraSkillID)
			return false
		}

		/* parse AcitvePassiveSkills field */
		vecAcitvePassiveSkills, _ := parse.GetFieldByName(uint32(i), "acitvePassiveSkills")
		if vecAcitvePassiveSkills != "" {
			arrayAcitvePassiveSkills := strings.Split(vecAcitvePassiveSkills, ",")
			for j := 0; j < len(arrayAcitvePassiveSkills); j++ {
				v, ret := String2Int32(arrayAcitvePassiveSkills[j])
				if !ret {
					glog.Error("Parse CfgMonsterBattle.AcitvePassiveSkills field error, value:", arrayAcitvePassiveSkills[j])
					return false
				}
				data.AcitvePassiveSkills = append(data.AcitvePassiveSkills, v)
			}
		}

		/* parse CollideRect field */
		vecCollideRect, _ := parse.GetFieldByName(uint32(i), "collideRect")
		arrayCollideRect := strings.Split(vecCollideRect, ",")
		for j := 0; j < len(arrayCollideRect); j++ {
			v, ret := String2Float(arrayCollideRect[j])
			if !ret {
				glog.Error("Parse CfgMonsterBattle.CollideRect field error,value:", arrayCollideRect[j])
				return false
			}
			data.CollideRect = append(data.CollideRect, v)
		}

		/* parse MeleeHitRange field */
		vMeleeHitRange, _ := parse.GetFieldByName(uint32(i), "meleeHitRange")
		var MeleeHitRangeRet bool
		data.MeleeHitRange, MeleeHitRangeRet = String2Float(vMeleeHitRange)
		if !MeleeHitRangeRet {
			glog.Error("Parse CfgMonsterBattle.MeleeHitRange field error,value:", vMeleeHitRange)
		}

		/* parse RemoteHitRange field */
		vRemoteHitRange, _ := parse.GetFieldByName(uint32(i), "remoteHitRange")
		var RemoteHitRangeRet bool
		data.RemoteHitRange, RemoteHitRangeRet = String2Float(vRemoteHitRange)
		if !RemoteHitRangeRet {
			glog.Error("Parse CfgMonsterBattle.RemoteHitRange field error,value:", vRemoteHitRange)
		}

		/* parse KnockBackDuration field */
		vKnockBackDuration, _ := parse.GetFieldByName(uint32(i), "knockBackDuration")
		var KnockBackDurationRet bool
		data.KnockBackDuration, KnockBackDurationRet = String2Float(vKnockBackDuration)
		if !KnockBackDurationRet {
			glog.Error("Parse CfgMonsterBattle.KnockBackDuration field error,value:", vKnockBackDuration)
		}

		/* parse GetUpDuration field */
		vGetUpDuration, _ := parse.GetFieldByName(uint32(i), "getUpDuration")
		var GetUpDurationRet bool
		data.GetUpDuration, GetUpDurationRet = String2Float(vGetUpDuration)
		if !GetUpDurationRet {
			glog.Error("Parse CfgMonsterBattle.GetUpDuration field error,value:", vGetUpDuration)
		}

		/* parse DamagedDuration field */
		vDamagedDuration, _ := parse.GetFieldByName(uint32(i), "damagedDuration")
		var DamagedDurationRet bool
		data.DamagedDuration, DamagedDurationRet = String2Float(vDamagedDuration)
		if !DamagedDurationRet {
			glog.Error("Parse CfgMonsterBattle.DamagedDuration field error,value:", vDamagedDuration)
		}

		/* parse EnterDuration field */
		vEnterDuration, _ := parse.GetFieldByName(uint32(i), "enterDuration")
		var EnterDurationRet bool
		data.EnterDuration, EnterDurationRet = String2Float(vEnterDuration)
		if !EnterDurationRet {
			glog.Error("Parse CfgMonsterBattle.EnterDuration field error,value:", vEnterDuration)
		}

		/* parse BkStartDuration field */
		vBkStartDuration, _ := parse.GetFieldByName(uint32(i), "bkStartDuration")
		var BkStartDurationRet bool
		data.BkStartDuration, BkStartDurationRet = String2Float(vBkStartDuration)
		if !BkStartDurationRet {
			glog.Error("Parse CfgMonsterBattle.BkStartDuration field error,value:", vBkStartDuration)
		}

		/* parse BkEndDuration field */
		vBkEndDuration, _ := parse.GetFieldByName(uint32(i), "bkEndDuration")
		var BkEndDurationRet bool
		data.BkEndDuration, BkEndDurationRet = String2Float(vBkEndDuration)
		if !BkEndDurationRet {
			glog.Error("Parse CfgMonsterBattle.BkEndDuration field error,value:", vBkEndDuration)
		}

		/* parse EnterDelay field */
		vEnterDelay, _ := parse.GetFieldByName(uint32(i), "enterDelay")
		var EnterDelayRet bool
		data.EnterDelay, EnterDelayRet = String2Float(vEnterDelay)
		if !EnterDelayRet {
			glog.Error("Parse CfgMonsterBattle.EnterDelay field error,value:", vEnterDelay)
		}

		/* parse TranslateData field */
		data.TranslateData, _ = parse.GetFieldByName(uint32(i), "translateData")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgMonsterBattleConfig) Clear() {
}

func (c *CfgMonsterBattleConfig) Find(id int32) (*CfgMonsterBattle, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgMonsterBattleConfig) GetAllData() map[int32]*CfgMonsterBattle {
	return c.data
}

func (c *CfgMonsterBattleConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.CharaName, ",", v.AiName, ",", v.Attackinterval, ",", v.NormalAtkID, ",", v.NormalMeleeAtkID, ",", v.MoveAtkID, ",", v.EnterSkillID, ",", v.WarningSkill, ",", v.UiSkillIDs, ",", v.SkillID, ",", v.UltraSkillID, ",", v.AcitvePassiveSkills, ",", v.CollideRect, ",", v.MeleeHitRange, ",", v.RemoteHitRange, ",", v.KnockBackDuration, ",", v.GetUpDuration, ",", v.DamagedDuration, ",", v.EnterDuration, ",", v.BkStartDuration, ",", v.BkEndDuration, ",", v.EnterDelay, ",", v.TranslateData)
	}
}
