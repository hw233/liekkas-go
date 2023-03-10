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

type CfgPassiveEffect struct {
	Id               int32
	Comment          string
	EventsID         []int32
	ObjectEffectsID  []int32
	TriggerEffectsID []int32
	TriggerObject    []int32
	TriggerType      int32
	TriggerParam     []int32
	TriggerCD        float64
	TriggerPR        int32
	Duration         float64
	Type             int32
}

type CfgPassiveEffectConfig struct {
	data map[int32]*CfgPassiveEffect
}

func NewCfgPassiveEffectConfig() *CfgPassiveEffectConfig {
	return &CfgPassiveEffectConfig{
		data: make(map[int32]*CfgPassiveEffect),
	}
}

func (c *CfgPassiveEffectConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgPassiveEffect)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgPassiveEffect.Id field error,value:", vId)
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
					glog.Error("Parse CfgPassiveEffect.EventsID field error, value:", arrayEventsID[j])
					return false
				}
				data.EventsID = append(data.EventsID, v)
			}
		}

		/* parse ObjectEffectsID field */
		vecObjectEffectsID, _ := parse.GetFieldByName(uint32(i), "objectEffectsID")
		if vecObjectEffectsID != "" {
			arrayObjectEffectsID := strings.Split(vecObjectEffectsID, ",")
			for j := 0; j < len(arrayObjectEffectsID); j++ {
				v, ret := String2Int32(arrayObjectEffectsID[j])
				if !ret {
					glog.Error("Parse CfgPassiveEffect.ObjectEffectsID field error, value:", arrayObjectEffectsID[j])
					return false
				}
				data.ObjectEffectsID = append(data.ObjectEffectsID, v)
			}
		}

		/* parse TriggerEffectsID field */
		vecTriggerEffectsID, _ := parse.GetFieldByName(uint32(i), "triggerEffectsID")
		if vecTriggerEffectsID != "" {
			arrayTriggerEffectsID := strings.Split(vecTriggerEffectsID, ",")
			for j := 0; j < len(arrayTriggerEffectsID); j++ {
				v, ret := String2Int32(arrayTriggerEffectsID[j])
				if !ret {
					glog.Error("Parse CfgPassiveEffect.TriggerEffectsID field error, value:", arrayTriggerEffectsID[j])
					return false
				}
				data.TriggerEffectsID = append(data.TriggerEffectsID, v)
			}
		}

		/* parse TriggerObject field */
		vecTriggerObject, _ := parse.GetFieldByName(uint32(i), "triggerObject")
		if vecTriggerObject != "" {
			arrayTriggerObject := strings.Split(vecTriggerObject, ",")
			for j := 0; j < len(arrayTriggerObject); j++ {
				v, ret := String2Int32(arrayTriggerObject[j])
				if !ret {
					glog.Error("Parse CfgPassiveEffect.TriggerObject field error, value:", arrayTriggerObject[j])
					return false
				}
				data.TriggerObject = append(data.TriggerObject, v)
			}
		}

		/* parse TriggerType field */
		vTriggerType, _ := parse.GetFieldByName(uint32(i), "triggerType")
		var TriggerTypeRet bool
		data.TriggerType, TriggerTypeRet = String2Int32(vTriggerType)
		if !TriggerTypeRet {
			glog.Error("Parse CfgPassiveEffect.TriggerType field error,value:", vTriggerType)
			return false
		}

		/* parse TriggerParam field */
		vecTriggerParam, _ := parse.GetFieldByName(uint32(i), "triggerParam")
		if vecTriggerParam != "" {
			arrayTriggerParam := strings.Split(vecTriggerParam, ",")
			for j := 0; j < len(arrayTriggerParam); j++ {
				v, ret := String2Int32(arrayTriggerParam[j])
				if !ret {
					glog.Error("Parse CfgPassiveEffect.TriggerParam field error, value:", arrayTriggerParam[j])
					return false
				}
				data.TriggerParam = append(data.TriggerParam, v)
			}
		}

		/* parse TriggerCD field */
		vTriggerCD, _ := parse.GetFieldByName(uint32(i), "triggerCD")
		var TriggerCDRet bool
		data.TriggerCD, TriggerCDRet = String2Float(vTriggerCD)
		if !TriggerCDRet {
			glog.Error("Parse CfgPassiveEffect.TriggerCD field error,value:", vTriggerCD)
		}

		/* parse TriggerPR field */
		vTriggerPR, _ := parse.GetFieldByName(uint32(i), "triggerPR")
		var TriggerPRRet bool
		data.TriggerPR, TriggerPRRet = String2Int32(vTriggerPR)
		if !TriggerPRRet {
			glog.Error("Parse CfgPassiveEffect.TriggerPR field error,value:", vTriggerPR)
			return false
		}

		/* parse Duration field */
		vDuration, _ := parse.GetFieldByName(uint32(i), "duration")
		var DurationRet bool
		data.Duration, DurationRet = String2Float(vDuration)
		if !DurationRet {
			glog.Error("Parse CfgPassiveEffect.Duration field error,value:", vDuration)
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgPassiveEffect.Type field error,value:", vType)
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

func (c *CfgPassiveEffectConfig) Clear() {
}

func (c *CfgPassiveEffectConfig) Find(id int32) (*CfgPassiveEffect, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgPassiveEffectConfig) GetAllData() map[int32]*CfgPassiveEffect {
	return c.data
}

func (c *CfgPassiveEffectConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Comment, ",", v.EventsID, ",", v.ObjectEffectsID, ",", v.TriggerEffectsID, ",", v.TriggerObject, ",", v.TriggerType, ",", v.TriggerParam, ",", v.TriggerCD, ",", v.TriggerPR, ",", v.Duration, ",", v.Type)
	}
}
