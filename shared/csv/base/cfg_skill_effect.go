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

type CfgSkillEffect struct {
	Id            int32
	BfPriority    int32
	BfGroup       int32
	BfOverlayMax  int32
	BfType        int32
	BfSubType     int32
	BfOverlayType int32
	EffectName    string
	Interval      float64
	Duration      float64
	Value1        string
	Value2        string
	Value3        string
	Value4        string
}

type CfgSkillEffectConfig struct {
	data map[int32]*CfgSkillEffect
}

func NewCfgSkillEffectConfig() *CfgSkillEffectConfig {
	return &CfgSkillEffectConfig{
		data: make(map[int32]*CfgSkillEffect),
	}
}

func (c *CfgSkillEffectConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgSkillEffect)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgSkillEffect.Id field error,value:", vId)
			return false
		}

		/* parse BfPriority field */
		vBfPriority, _ := parse.GetFieldByName(uint32(i), "bfPriority")
		var BfPriorityRet bool
		data.BfPriority, BfPriorityRet = String2Int32(vBfPriority)
		if !BfPriorityRet {
			glog.Error("Parse CfgSkillEffect.BfPriority field error,value:", vBfPriority)
			return false
		}

		/* parse BfGroup field */
		vBfGroup, _ := parse.GetFieldByName(uint32(i), "bfGroup")
		var BfGroupRet bool
		data.BfGroup, BfGroupRet = String2Int32(vBfGroup)
		if !BfGroupRet {
			glog.Error("Parse CfgSkillEffect.BfGroup field error,value:", vBfGroup)
			return false
		}

		/* parse BfOverlayMax field */
		vBfOverlayMax, _ := parse.GetFieldByName(uint32(i), "bfOverlayMax")
		var BfOverlayMaxRet bool
		data.BfOverlayMax, BfOverlayMaxRet = String2Int32(vBfOverlayMax)
		if !BfOverlayMaxRet {
			glog.Error("Parse CfgSkillEffect.BfOverlayMax field error,value:", vBfOverlayMax)
			return false
		}

		/* parse BfType field */
		vBfType, _ := parse.GetFieldByName(uint32(i), "bfType")
		var BfTypeRet bool
		data.BfType, BfTypeRet = String2Int32(vBfType)
		if !BfTypeRet {
			glog.Error("Parse CfgSkillEffect.BfType field error,value:", vBfType)
			return false
		}

		/* parse BfSubType field */
		vBfSubType, _ := parse.GetFieldByName(uint32(i), "bfSubType")
		var BfSubTypeRet bool
		data.BfSubType, BfSubTypeRet = String2Int32(vBfSubType)
		if !BfSubTypeRet {
			glog.Error("Parse CfgSkillEffect.BfSubType field error,value:", vBfSubType)
			return false
		}

		/* parse BfOverlayType field */
		vBfOverlayType, _ := parse.GetFieldByName(uint32(i), "bfOverlayType")
		var BfOverlayTypeRet bool
		data.BfOverlayType, BfOverlayTypeRet = String2Int32(vBfOverlayType)
		if !BfOverlayTypeRet {
			glog.Error("Parse CfgSkillEffect.BfOverlayType field error,value:", vBfOverlayType)
			return false
		}

		/* parse EffectName field */
		data.EffectName, _ = parse.GetFieldByName(uint32(i), "effectName")

		/* parse Interval field */
		vInterval, _ := parse.GetFieldByName(uint32(i), "interval")
		var IntervalRet bool
		data.Interval, IntervalRet = String2Float(vInterval)
		if !IntervalRet {
			glog.Error("Parse CfgSkillEffect.Interval field error,value:", vInterval)
		}

		/* parse Duration field */
		vDuration, _ := parse.GetFieldByName(uint32(i), "duration")
		var DurationRet bool
		data.Duration, DurationRet = String2Float(vDuration)
		if !DurationRet {
			glog.Error("Parse CfgSkillEffect.Duration field error,value:", vDuration)
		}

		/* parse Value1 field */
		data.Value1, _ = parse.GetFieldByName(uint32(i), "value1")

		/* parse Value2 field */
		data.Value2, _ = parse.GetFieldByName(uint32(i), "value2")

		/* parse Value3 field */
		data.Value3, _ = parse.GetFieldByName(uint32(i), "value3")

		/* parse Value4 field */
		data.Value4, _ = parse.GetFieldByName(uint32(i), "value4")

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgSkillEffectConfig) Clear() {
}

func (c *CfgSkillEffectConfig) Find(id int32) (*CfgSkillEffect, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgSkillEffectConfig) GetAllData() map[int32]*CfgSkillEffect {
	return c.data
}

func (c *CfgSkillEffectConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.BfPriority, ",", v.BfGroup, ",", v.BfOverlayMax, ",", v.BfType, ",", v.BfSubType, ",", v.BfOverlayType, ",", v.EffectName, ",", v.Interval, ",", v.Duration, ",", v.Value1, ",", v.Value2, ",", v.Value3, ",", v.Value4)
	}
}
