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

type CfgSkill struct {
	Id               int32
	Comment          string
	EventsID         []int32
	Duration         float64
	Cd               float64
	Icd              float64
	Anim             string
	AnimType         int32
	AnimAngle        []int32
	AnimCenter       []float64
	UltraAnim        int32
	IsAtkSpd         bool
	Next             int32
	Target           int32
	Cost             int32
	PreloadCharacter []int32
	PreloadFx        []int32
	PreloadBullet    []int32
}

type CfgSkillConfig struct {
	data map[int32]*CfgSkill
}

func NewCfgSkillConfig() *CfgSkillConfig {
	return &CfgSkillConfig{
		data: make(map[int32]*CfgSkill),
	}
}

func (c *CfgSkillConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgSkill)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgSkill.Id field error,value:", vId)
			return false
		}

		/* parse Comment field */
		data.Comment, _ = parse.GetFieldByName(uint32(i), "comment")

		/* parse EventsID field */
		vecEventsID, _ := parse.GetFieldByName(uint32(i), "eventsID")
		if vecEventsID != "" {
			arrayEventsID := strings.Split(vecEventsID, ",")
			for j := 0; j < len(arrayEventsID); j++ {
				v, ret := String2Int32(arrayEventsID[j])
				if !ret {
					glog.Error("Parse CfgSkill.EventsID field error, value:", arrayEventsID[j])
					return false
				}
				data.EventsID = append(data.EventsID, v)
			}
		}

		/* parse Duration field */
		vDuration, _ := parse.GetFieldByName(uint32(i), "duration")
		var DurationRet bool
		data.Duration, DurationRet = String2Float(vDuration)
		if !DurationRet {
			glog.Error("Parse CfgSkill.Duration field error,value:", vDuration)
		}

		/* parse Cd field */
		vCd, _ := parse.GetFieldByName(uint32(i), "cd")
		var CdRet bool
		data.Cd, CdRet = String2Float(vCd)
		if !CdRet {
			glog.Error("Parse CfgSkill.Cd field error,value:", vCd)
		}

		/* parse Icd field */
		vIcd, _ := parse.GetFieldByName(uint32(i), "icd")
		var IcdRet bool
		data.Icd, IcdRet = String2Float(vIcd)
		if !IcdRet {
			glog.Error("Parse CfgSkill.Icd field error,value:", vIcd)
		}

		/* parse Anim field */
		data.Anim, _ = parse.GetFieldByName(uint32(i), "anim")

		/* parse AnimType field */
		vAnimType, _ := parse.GetFieldByName(uint32(i), "animType")
		var AnimTypeRet bool
		data.AnimType, AnimTypeRet = String2Int32(vAnimType)
		if !AnimTypeRet {
			glog.Error("Parse CfgSkill.AnimType field error,value:", vAnimType)
			return false
		}

		/* parse AnimAngle field */
		vecAnimAngle, _ := parse.GetFieldByName(uint32(i), "animAngle")
		if vecAnimAngle != "" {
			arrayAnimAngle := strings.Split(vecAnimAngle, ",")
			for j := 0; j < len(arrayAnimAngle); j++ {
				v, ret := String2Int32(arrayAnimAngle[j])
				if !ret {
					glog.Error("Parse CfgSkill.AnimAngle field error, value:", arrayAnimAngle[j])
					return false
				}
				data.AnimAngle = append(data.AnimAngle, v)
			}
		}

		/* parse AnimCenter field */
		vecAnimCenter, _ := parse.GetFieldByName(uint32(i), "animCenter")
		arrayAnimCenter := strings.Split(vecAnimCenter, ",")
		for j := 0; j < len(arrayAnimCenter); j++ {
			v, ret := String2Float(arrayAnimCenter[j])
			if !ret {
				glog.Error("Parse CfgSkill.AnimCenter field error,value:", arrayAnimCenter[j])
				return false
			}
			data.AnimCenter = append(data.AnimCenter, v)
		}

		/* parse UltraAnim field */
		vUltraAnim, _ := parse.GetFieldByName(uint32(i), "ultraAnim")
		var UltraAnimRet bool
		data.UltraAnim, UltraAnimRet = String2Int32(vUltraAnim)
		if !UltraAnimRet {
			glog.Error("Parse CfgSkill.UltraAnim field error,value:", vUltraAnim)
			return false
		}

		/* parse IsAtkSpd field */
		vIsAtkSpd, _ := parse.GetFieldByName(uint32(i), "isAtkSpd")
		var IsAtkSpdRet bool
		data.IsAtkSpd, IsAtkSpdRet = String2Bool(vIsAtkSpd)
		if !IsAtkSpdRet {
			glog.Error("Parse CfgSkill.IsAtkSpd field error,value:", vIsAtkSpd)
		}

		/* parse Next field */
		vNext, _ := parse.GetFieldByName(uint32(i), "next")
		var NextRet bool
		data.Next, NextRet = String2Int32(vNext)
		if !NextRet {
			glog.Error("Parse CfgSkill.Next field error,value:", vNext)
			return false
		}

		/* parse Target field */
		vTarget, _ := parse.GetFieldByName(uint32(i), "target")
		var TargetRet bool
		data.Target, TargetRet = String2Int32(vTarget)
		if !TargetRet {
			glog.Error("Parse CfgSkill.Target field error,value:", vTarget)
			return false
		}

		/* parse Cost field */
		vCost, _ := parse.GetFieldByName(uint32(i), "cost")
		var CostRet bool
		data.Cost, CostRet = String2Int32(vCost)
		if !CostRet {
			glog.Error("Parse CfgSkill.Cost field error,value:", vCost)
			return false
		}

		/* parse PreloadCharacter field */
		vecPreloadCharacter, _ := parse.GetFieldByName(uint32(i), "preloadCharacter")
		if vecPreloadCharacter != "" {
			arrayPreloadCharacter := strings.Split(vecPreloadCharacter, ",")
			for j := 0; j < len(arrayPreloadCharacter); j++ {
				v, ret := String2Int32(arrayPreloadCharacter[j])
				if !ret {
					glog.Error("Parse CfgSkill.PreloadCharacter field error, value:", arrayPreloadCharacter[j])
					return false
				}
				data.PreloadCharacter = append(data.PreloadCharacter, v)
			}
		}

		/* parse PreloadFx field */
		vecPreloadFx, _ := parse.GetFieldByName(uint32(i), "preloadFx")
		if vecPreloadFx != "" {
			arrayPreloadFx := strings.Split(vecPreloadFx, ",")
			for j := 0; j < len(arrayPreloadFx); j++ {
				v, ret := String2Int32(arrayPreloadFx[j])
				if !ret {
					glog.Error("Parse CfgSkill.PreloadFx field error, value:", arrayPreloadFx[j])
					return false
				}
				data.PreloadFx = append(data.PreloadFx, v)
			}
		}

		/* parse PreloadBullet field */
		vecPreloadBullet, _ := parse.GetFieldByName(uint32(i), "preloadBullet")
		if vecPreloadBullet != "" {
			arrayPreloadBullet := strings.Split(vecPreloadBullet, ",")
			for j := 0; j < len(arrayPreloadBullet); j++ {
				v, ret := String2Int32(arrayPreloadBullet[j])
				if !ret {
					glog.Error("Parse CfgSkill.PreloadBullet field error, value:", arrayPreloadBullet[j])
					return false
				}
				data.PreloadBullet = append(data.PreloadBullet, v)
			}
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgSkillConfig) Clear() {
}

func (c *CfgSkillConfig) Find(id int32) (*CfgSkill, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgSkillConfig) GetAllData() map[int32]*CfgSkill {
	return c.data
}

func (c *CfgSkillConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Comment, ",", v.EventsID, ",", v.Duration, ",", v.Cd, ",", v.Icd, ",", v.Anim, ",", v.AnimType, ",", v.AnimAngle, ",", v.AnimCenter, ",", v.UltraAnim, ",", v.IsAtkSpd, ",", v.Next, ",", v.Target, ",", v.Cost, ",", v.PreloadCharacter, ",", v.PreloadFx, ",", v.PreloadBullet)
	}
}
