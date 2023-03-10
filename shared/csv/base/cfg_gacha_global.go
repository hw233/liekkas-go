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

type CfgGachaGlobal struct {
	Id                       int32
	DropWeightN              int32
	DropWeightR              int32
	DropWeightSR             int32
	DropProbSSR              int32
	GuaranteeCountSR         int32
	GuaranteeTrigCountSSR    int32
	GuaranteeTrigProbIncrSSR int32
}

type CfgGachaGlobalConfig struct {
	data map[int32]*CfgGachaGlobal
}

func NewCfgGachaGlobalConfig() *CfgGachaGlobalConfig {
	return &CfgGachaGlobalConfig{
		data: make(map[int32]*CfgGachaGlobal),
	}
}

func (c *CfgGachaGlobalConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGachaGlobal)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGachaGlobal.Id field error,value:", vId)
			return false
		}

		/* parse DropWeightN field */
		vDropWeightN, _ := parse.GetFieldByName(uint32(i), "dropWeightN")
		var DropWeightNRet bool
		data.DropWeightN, DropWeightNRet = String2Int32(vDropWeightN)
		if !DropWeightNRet {
			glog.Error("Parse CfgGachaGlobal.DropWeightN field error,value:", vDropWeightN)
			return false
		}

		/* parse DropWeightR field */
		vDropWeightR, _ := parse.GetFieldByName(uint32(i), "dropWeightR")
		var DropWeightRRet bool
		data.DropWeightR, DropWeightRRet = String2Int32(vDropWeightR)
		if !DropWeightRRet {
			glog.Error("Parse CfgGachaGlobal.DropWeightR field error,value:", vDropWeightR)
			return false
		}

		/* parse DropWeightSR field */
		vDropWeightSR, _ := parse.GetFieldByName(uint32(i), "dropWeightSR")
		var DropWeightSRRet bool
		data.DropWeightSR, DropWeightSRRet = String2Int32(vDropWeightSR)
		if !DropWeightSRRet {
			glog.Error("Parse CfgGachaGlobal.DropWeightSR field error,value:", vDropWeightSR)
			return false
		}

		/* parse DropProbSSR field */
		vDropProbSSR, _ := parse.GetFieldByName(uint32(i), "dropProbSSR")
		var DropProbSSRRet bool
		data.DropProbSSR, DropProbSSRRet = String2Int32(vDropProbSSR)
		if !DropProbSSRRet {
			glog.Error("Parse CfgGachaGlobal.DropProbSSR field error,value:", vDropProbSSR)
			return false
		}

		/* parse GuaranteeCountSR field */
		vGuaranteeCountSR, _ := parse.GetFieldByName(uint32(i), "guaranteeCountSR")
		var GuaranteeCountSRRet bool
		data.GuaranteeCountSR, GuaranteeCountSRRet = String2Int32(vGuaranteeCountSR)
		if !GuaranteeCountSRRet {
			glog.Error("Parse CfgGachaGlobal.GuaranteeCountSR field error,value:", vGuaranteeCountSR)
			return false
		}

		/* parse GuaranteeTrigCountSSR field */
		vGuaranteeTrigCountSSR, _ := parse.GetFieldByName(uint32(i), "guaranteeTrigCountSSR")
		var GuaranteeTrigCountSSRRet bool
		data.GuaranteeTrigCountSSR, GuaranteeTrigCountSSRRet = String2Int32(vGuaranteeTrigCountSSR)
		if !GuaranteeTrigCountSSRRet {
			glog.Error("Parse CfgGachaGlobal.GuaranteeTrigCountSSR field error,value:", vGuaranteeTrigCountSSR)
			return false
		}

		/* parse GuaranteeTrigProbIncrSSR field */
		vGuaranteeTrigProbIncrSSR, _ := parse.GetFieldByName(uint32(i), "guaranteeTrigProbIncrSSR")
		var GuaranteeTrigProbIncrSSRRet bool
		data.GuaranteeTrigProbIncrSSR, GuaranteeTrigProbIncrSSRRet = String2Int32(vGuaranteeTrigProbIncrSSR)
		if !GuaranteeTrigProbIncrSSRRet {
			glog.Error("Parse CfgGachaGlobal.GuaranteeTrigProbIncrSSR field error,value:", vGuaranteeTrigProbIncrSSR)
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

func (c *CfgGachaGlobalConfig) Clear() {
}

func (c *CfgGachaGlobalConfig) Find(id int32) (*CfgGachaGlobal, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGachaGlobalConfig) GetAllData() map[int32]*CfgGachaGlobal {
	return c.data
}

func (c *CfgGachaGlobalConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.DropWeightN, ",", v.DropWeightR, ",", v.DropWeightSR, ",", v.DropProbSSR, ",", v.GuaranteeCountSR, ",", v.GuaranteeTrigCountSSR, ",", v.GuaranteeTrigProbIncrSSR)
	}
}
