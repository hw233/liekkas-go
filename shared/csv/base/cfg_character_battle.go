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

type CfgCharacterBattle struct {
	Id                  int32
	Name                string
	AiName              string
	Attackinterval      int32
	NormalAtkID         int32
	NormalMeleeAtkID    int32
	SkillID1            int32
	SkillID2            int32
	UltraSkillID        int32
	PassiveSkillID      int32
	EnterSkillID        int32
	MoveAtkID           int32
	AcitvePassiveSkills []int32
	AttackRange         []float64
	AttackVerticalRange []float64
	HitRect             []float64
	CollideRect         []float64
	MeleeHitRange       float64
	RemoteHitRange      float64
	KnockBackDuration   float64
	GetUpDuration       float64
	DamagedDuration     float64
	EnterDuration       float64
	EnterDelay          float64
	TranslateData       string
}

type CfgCharacterBattleConfig struct {
	data map[int32]*CfgCharacterBattle
}

func NewCfgCharacterBattleConfig() *CfgCharacterBattleConfig {
	return &CfgCharacterBattleConfig{
		data: make(map[int32]*CfgCharacterBattle),
	}
}

func (c *CfgCharacterBattleConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgCharacterBattle)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgCharacterBattle.Id field error,value:", vId)
			return false
		}

		/* parse Name field */
		data.Name, _ = parse.GetFieldByName(uint32(i), "name")

		/* parse AiName field */
		data.AiName, _ = parse.GetFieldByName(uint32(i), "aiName")

		/* parse Attackinterval field */
		vAttackinterval, _ := parse.GetFieldByName(uint32(i), "attackinterval")
		var AttackintervalRet bool
		data.Attackinterval, AttackintervalRet = String2Int32(vAttackinterval)
		if !AttackintervalRet {
			glog.Error("Parse CfgCharacterBattle.Attackinterval field error,value:", vAttackinterval)
			return false
		}

		/* parse NormalAtkID field */
		vNormalAtkID, _ := parse.GetFieldByName(uint32(i), "normalAtkID")
		var NormalAtkIDRet bool
		data.NormalAtkID, NormalAtkIDRet = String2Int32(vNormalAtkID)
		if !NormalAtkIDRet {
			glog.Error("Parse CfgCharacterBattle.NormalAtkID field error,value:", vNormalAtkID)
			return false
		}

		/* parse NormalMeleeAtkID field */
		vNormalMeleeAtkID, _ := parse.GetFieldByName(uint32(i), "normalMeleeAtkID")
		var NormalMeleeAtkIDRet bool
		data.NormalMeleeAtkID, NormalMeleeAtkIDRet = String2Int32(vNormalMeleeAtkID)
		if !NormalMeleeAtkIDRet {
			glog.Error("Parse CfgCharacterBattle.NormalMeleeAtkID field error,value:", vNormalMeleeAtkID)
			return false
		}

		/* parse SkillID1 field */
		vSkillID1, _ := parse.GetFieldByName(uint32(i), "skillID1")
		var SkillID1Ret bool
		data.SkillID1, SkillID1Ret = String2Int32(vSkillID1)
		if !SkillID1Ret {
			glog.Error("Parse CfgCharacterBattle.SkillID1 field error,value:", vSkillID1)
			return false
		}

		/* parse SkillID2 field */
		vSkillID2, _ := parse.GetFieldByName(uint32(i), "skillID2")
		var SkillID2Ret bool
		data.SkillID2, SkillID2Ret = String2Int32(vSkillID2)
		if !SkillID2Ret {
			glog.Error("Parse CfgCharacterBattle.SkillID2 field error,value:", vSkillID2)
			return false
		}

		/* parse UltraSkillID field */
		vUltraSkillID, _ := parse.GetFieldByName(uint32(i), "ultraSkillID")
		var UltraSkillIDRet bool
		data.UltraSkillID, UltraSkillIDRet = String2Int32(vUltraSkillID)
		if !UltraSkillIDRet {
			glog.Error("Parse CfgCharacterBattle.UltraSkillID field error,value:", vUltraSkillID)
			return false
		}

		/* parse PassiveSkillID field */
		vPassiveSkillID, _ := parse.GetFieldByName(uint32(i), "passiveSkillID")
		var PassiveSkillIDRet bool
		data.PassiveSkillID, PassiveSkillIDRet = String2Int32(vPassiveSkillID)
		if !PassiveSkillIDRet {
			glog.Error("Parse CfgCharacterBattle.PassiveSkillID field error,value:", vPassiveSkillID)
			return false
		}

		/* parse EnterSkillID field */
		vEnterSkillID, _ := parse.GetFieldByName(uint32(i), "enterSkillID")
		var EnterSkillIDRet bool
		data.EnterSkillID, EnterSkillIDRet = String2Int32(vEnterSkillID)
		if !EnterSkillIDRet {
			glog.Error("Parse CfgCharacterBattle.EnterSkillID field error,value:", vEnterSkillID)
			return false
		}

		/* parse MoveAtkID field */
		vMoveAtkID, _ := parse.GetFieldByName(uint32(i), "moveAtkID")
		var MoveAtkIDRet bool
		data.MoveAtkID, MoveAtkIDRet = String2Int32(vMoveAtkID)
		if !MoveAtkIDRet {
			glog.Error("Parse CfgCharacterBattle.MoveAtkID field error,value:", vMoveAtkID)
			return false
		}

		/* parse AcitvePassiveSkills field */
		vecAcitvePassiveSkills, _ := parse.GetFieldByName(uint32(i), "acitvePassiveSkills")
		if vecAcitvePassiveSkills != "" {
			arrayAcitvePassiveSkills := strings.Split(vecAcitvePassiveSkills, ",")
			for j := 0; j < len(arrayAcitvePassiveSkills); j++ {
				v, ret := String2Int32(arrayAcitvePassiveSkills[j])
				if !ret {
					glog.Error("Parse CfgCharacterBattle.AcitvePassiveSkills field error, value:", arrayAcitvePassiveSkills[j])
					return false
				}
				data.AcitvePassiveSkills = append(data.AcitvePassiveSkills, v)
			}
		}

		/* parse AttackRange field */
		vecAttackRange, _ := parse.GetFieldByName(uint32(i), "attackRange")
		arrayAttackRange := strings.Split(vecAttackRange, ",")
		for j := 0; j < len(arrayAttackRange); j++ {
			v, ret := String2Float(arrayAttackRange[j])
			if !ret {
				glog.Error("Parse CfgCharacterBattle.AttackRange field error,value:", arrayAttackRange[j])
				return false
			}
			data.AttackRange = append(data.AttackRange, v)
		}

		/* parse AttackVerticalRange field */
		vecAttackVerticalRange, _ := parse.GetFieldByName(uint32(i), "attackVerticalRange")
		arrayAttackVerticalRange := strings.Split(vecAttackVerticalRange, ",")
		for j := 0; j < len(arrayAttackVerticalRange); j++ {
			v, ret := String2Float(arrayAttackVerticalRange[j])
			if !ret {
				glog.Error("Parse CfgCharacterBattle.AttackVerticalRange field error,value:", arrayAttackVerticalRange[j])
				return false
			}
			data.AttackVerticalRange = append(data.AttackVerticalRange, v)
		}

		/* parse HitRect field */
		vecHitRect, _ := parse.GetFieldByName(uint32(i), "hitRect")
		arrayHitRect := strings.Split(vecHitRect, ",")
		for j := 0; j < len(arrayHitRect); j++ {
			v, ret := String2Float(arrayHitRect[j])
			if !ret {
				glog.Error("Parse CfgCharacterBattle.HitRect field error,value:", arrayHitRect[j])
				return false
			}
			data.HitRect = append(data.HitRect, v)
		}

		/* parse CollideRect field */
		vecCollideRect, _ := parse.GetFieldByName(uint32(i), "collideRect")
		arrayCollideRect := strings.Split(vecCollideRect, ",")
		for j := 0; j < len(arrayCollideRect); j++ {
			v, ret := String2Float(arrayCollideRect[j])
			if !ret {
				glog.Error("Parse CfgCharacterBattle.CollideRect field error,value:", arrayCollideRect[j])
				return false
			}
			data.CollideRect = append(data.CollideRect, v)
		}

		/* parse MeleeHitRange field */
		vMeleeHitRange, _ := parse.GetFieldByName(uint32(i), "meleeHitRange")
		var MeleeHitRangeRet bool
		data.MeleeHitRange, MeleeHitRangeRet = String2Float(vMeleeHitRange)
		if !MeleeHitRangeRet {
			glog.Error("Parse CfgCharacterBattle.MeleeHitRange field error,value:", vMeleeHitRange)
		}

		/* parse RemoteHitRange field */
		vRemoteHitRange, _ := parse.GetFieldByName(uint32(i), "remoteHitRange")
		var RemoteHitRangeRet bool
		data.RemoteHitRange, RemoteHitRangeRet = String2Float(vRemoteHitRange)
		if !RemoteHitRangeRet {
			glog.Error("Parse CfgCharacterBattle.RemoteHitRange field error,value:", vRemoteHitRange)
		}

		/* parse KnockBackDuration field */
		vKnockBackDuration, _ := parse.GetFieldByName(uint32(i), "knockBackDuration")
		var KnockBackDurationRet bool
		data.KnockBackDuration, KnockBackDurationRet = String2Float(vKnockBackDuration)
		if !KnockBackDurationRet {
			glog.Error("Parse CfgCharacterBattle.KnockBackDuration field error,value:", vKnockBackDuration)
		}

		/* parse GetUpDuration field */
		vGetUpDuration, _ := parse.GetFieldByName(uint32(i), "getUpDuration")
		var GetUpDurationRet bool
		data.GetUpDuration, GetUpDurationRet = String2Float(vGetUpDuration)
		if !GetUpDurationRet {
			glog.Error("Parse CfgCharacterBattle.GetUpDuration field error,value:", vGetUpDuration)
		}

		/* parse DamagedDuration field */
		vDamagedDuration, _ := parse.GetFieldByName(uint32(i), "damagedDuration")
		var DamagedDurationRet bool
		data.DamagedDuration, DamagedDurationRet = String2Float(vDamagedDuration)
		if !DamagedDurationRet {
			glog.Error("Parse CfgCharacterBattle.DamagedDuration field error,value:", vDamagedDuration)
		}

		/* parse EnterDuration field */
		vEnterDuration, _ := parse.GetFieldByName(uint32(i), "enterDuration")
		var EnterDurationRet bool
		data.EnterDuration, EnterDurationRet = String2Float(vEnterDuration)
		if !EnterDurationRet {
			glog.Error("Parse CfgCharacterBattle.EnterDuration field error,value:", vEnterDuration)
		}

		/* parse EnterDelay field */
		vEnterDelay, _ := parse.GetFieldByName(uint32(i), "enterDelay")
		var EnterDelayRet bool
		data.EnterDelay, EnterDelayRet = String2Float(vEnterDelay)
		if !EnterDelayRet {
			glog.Error("Parse CfgCharacterBattle.EnterDelay field error,value:", vEnterDelay)
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

func (c *CfgCharacterBattleConfig) Clear() {
}

func (c *CfgCharacterBattleConfig) Find(id int32) (*CfgCharacterBattle, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgCharacterBattleConfig) GetAllData() map[int32]*CfgCharacterBattle {
	return c.data
}

func (c *CfgCharacterBattleConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Name, ",", v.AiName, ",", v.Attackinterval, ",", v.NormalAtkID, ",", v.NormalMeleeAtkID, ",", v.SkillID1, ",", v.SkillID2, ",", v.UltraSkillID, ",", v.PassiveSkillID, ",", v.EnterSkillID, ",", v.MoveAtkID, ",", v.AcitvePassiveSkills, ",", v.AttackRange, ",", v.AttackVerticalRange, ",", v.HitRect, ",", v.CollideRect, ",", v.MeleeHitRange, ",", v.RemoteHitRange, ",", v.KnockBackDuration, ",", v.GetUpDuration, ",", v.DamagedDuration, ",", v.EnterDuration, ",", v.EnterDelay, ",", v.TranslateData)
	}
}
