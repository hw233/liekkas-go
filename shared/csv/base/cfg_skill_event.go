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

type CfgSkillEvent struct {
	Id       int32
	Delay    float64
	Target1  int32
	Effects1 []int32
	Target2  int32
	Effects2 []int32
	Target3  int32
	Effects3 []int32
}

type CfgSkillEventConfig struct {
	data map[int32]*CfgSkillEvent
}

func NewCfgSkillEventConfig() *CfgSkillEventConfig {
	return &CfgSkillEventConfig{
		data: make(map[int32]*CfgSkillEvent),
	}
}

func (c *CfgSkillEventConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgSkillEvent)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgSkillEvent.Id field error,value:", vId)
			return false
		}

		/* parse Delay field */
		vDelay, _ := parse.GetFieldByName(uint32(i), "delay")
		var DelayRet bool
		data.Delay, DelayRet = String2Float(vDelay)
		if !DelayRet {
			glog.Error("Parse CfgSkillEvent.Delay field error,value:", vDelay)
		}

		/* parse Target1 field */
		vTarget1, _ := parse.GetFieldByName(uint32(i), "target1")
		var Target1Ret bool
		data.Target1, Target1Ret = String2Int32(vTarget1)
		if !Target1Ret {
			glog.Error("Parse CfgSkillEvent.Target1 field error,value:", vTarget1)
			return false
		}

		/* parse Effects1 field */
		vecEffects1, _ := parse.GetFieldByName(uint32(i), "effects1")
		if vecEffects1 != "" {
			arrayEffects1 := strings.Split(vecEffects1, ",")
			for j := 0; j < len(arrayEffects1); j++ {
				v, ret := String2Int32(arrayEffects1[j])
				if !ret {
					glog.Error("Parse CfgSkillEvent.Effects1 field error, value:", arrayEffects1[j])
					return false
				}
				data.Effects1 = append(data.Effects1, v)
			}
		}

		/* parse Target2 field */
		vTarget2, _ := parse.GetFieldByName(uint32(i), "target2")
		var Target2Ret bool
		data.Target2, Target2Ret = String2Int32(vTarget2)
		if !Target2Ret {
			glog.Error("Parse CfgSkillEvent.Target2 field error,value:", vTarget2)
			return false
		}

		/* parse Effects2 field */
		vecEffects2, _ := parse.GetFieldByName(uint32(i), "effects2")
		if vecEffects2 != "" {
			arrayEffects2 := strings.Split(vecEffects2, ",")
			for j := 0; j < len(arrayEffects2); j++ {
				v, ret := String2Int32(arrayEffects2[j])
				if !ret {
					glog.Error("Parse CfgSkillEvent.Effects2 field error, value:", arrayEffects2[j])
					return false
				}
				data.Effects2 = append(data.Effects2, v)
			}
		}

		/* parse Target3 field */
		vTarget3, _ := parse.GetFieldByName(uint32(i), "target3")
		var Target3Ret bool
		data.Target3, Target3Ret = String2Int32(vTarget3)
		if !Target3Ret {
			glog.Error("Parse CfgSkillEvent.Target3 field error,value:", vTarget3)
			return false
		}

		/* parse Effects3 field */
		vecEffects3, _ := parse.GetFieldByName(uint32(i), "effects3")
		if vecEffects3 != "" {
			arrayEffects3 := strings.Split(vecEffects3, ",")
			for j := 0; j < len(arrayEffects3); j++ {
				v, ret := String2Int32(arrayEffects3[j])
				if !ret {
					glog.Error("Parse CfgSkillEvent.Effects3 field error, value:", arrayEffects3[j])
					return false
				}
				data.Effects3 = append(data.Effects3, v)
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

func (c *CfgSkillEventConfig) Clear() {
}

func (c *CfgSkillEventConfig) Find(id int32) (*CfgSkillEvent, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgSkillEventConfig) GetAllData() map[int32]*CfgSkillEvent {
	return c.data
}

func (c *CfgSkillEventConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Delay, ",", v.Target1, ",", v.Effects1, ",", v.Target2, ",", v.Effects2, ",", v.Target3, ",", v.Effects3)
	}
}
