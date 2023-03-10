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

type CfgTpGate struct {
	Id                  int32
	StartMapType        int32
	StartMapId          int32
	UnlockConditions    []string
	UnlockUseConditions []string
	Destroy             int32
	PassLevelId         int32
	DestroyShow         bool
	CompleteLimit       int32
}

type CfgTpGateConfig struct {
	data map[int32]*CfgTpGate
}

func NewCfgTpGateConfig() *CfgTpGateConfig {
	return &CfgTpGateConfig{
		data: make(map[int32]*CfgTpGate),
	}
}

func (c *CfgTpGateConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgTpGate)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgTpGate.Id field error,value:", vId)
			return false
		}

		/* parse StartMapType field */
		vStartMapType, _ := parse.GetFieldByName(uint32(i), "startMapType")
		var StartMapTypeRet bool
		data.StartMapType, StartMapTypeRet = String2Int32(vStartMapType)
		if !StartMapTypeRet {
			glog.Error("Parse CfgTpGate.StartMapType field error,value:", vStartMapType)
			return false
		}

		/* parse StartMapId field */
		vStartMapId, _ := parse.GetFieldByName(uint32(i), "startMapId")
		var StartMapIdRet bool
		data.StartMapId, StartMapIdRet = String2Int32(vStartMapId)
		if !StartMapIdRet {
			glog.Error("Parse CfgTpGate.StartMapId field error,value:", vStartMapId)
			return false
		}

		/* parse UnlockConditions field */
		vecUnlockConditions, _ := parse.GetFieldByName(uint32(i), "unlockConditions")
		arrayUnlockConditions := strings.Split(vecUnlockConditions, ",")
		for j := 0; j < len(arrayUnlockConditions); j++ {
			v := arrayUnlockConditions[j]
			data.UnlockConditions = append(data.UnlockConditions, v)
		}

		/* parse UnlockUseConditions field */
		vecUnlockUseConditions, _ := parse.GetFieldByName(uint32(i), "unlockUseConditions")
		arrayUnlockUseConditions := strings.Split(vecUnlockUseConditions, ",")
		for j := 0; j < len(arrayUnlockUseConditions); j++ {
			v := arrayUnlockUseConditions[j]
			data.UnlockUseConditions = append(data.UnlockUseConditions, v)
		}

		/* parse Destroy field */
		vDestroy, _ := parse.GetFieldByName(uint32(i), "destroy")
		var DestroyRet bool
		data.Destroy, DestroyRet = String2Int32(vDestroy)
		if !DestroyRet {
			glog.Error("Parse CfgTpGate.Destroy field error,value:", vDestroy)
			return false
		}

		/* parse PassLevelId field */
		vPassLevelId, _ := parse.GetFieldByName(uint32(i), "passLevelId")
		var PassLevelIdRet bool
		data.PassLevelId, PassLevelIdRet = String2Int32(vPassLevelId)
		if !PassLevelIdRet {
			glog.Error("Parse CfgTpGate.PassLevelId field error,value:", vPassLevelId)
			return false
		}

		/* parse DestroyShow field */
		vDestroyShow, _ := parse.GetFieldByName(uint32(i), "destroyShow")
		var DestroyShowRet bool
		data.DestroyShow, DestroyShowRet = String2Bool(vDestroyShow)
		if !DestroyShowRet {
			glog.Error("Parse CfgTpGate.DestroyShow field error,value:", vDestroyShow)
		}

		/* parse CompleteLimit field */
		vCompleteLimit, _ := parse.GetFieldByName(uint32(i), "completeLimit")
		var CompleteLimitRet bool
		data.CompleteLimit, CompleteLimitRet = String2Int32(vCompleteLimit)
		if !CompleteLimitRet {
			glog.Error("Parse CfgTpGate.CompleteLimit field error,value:", vCompleteLimit)
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

func (c *CfgTpGateConfig) Clear() {
}

func (c *CfgTpGateConfig) Find(id int32) (*CfgTpGate, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgTpGateConfig) GetAllData() map[int32]*CfgTpGate {
	return c.data
}

func (c *CfgTpGateConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.StartMapType, ",", v.StartMapId, ",", v.UnlockConditions, ",", v.UnlockUseConditions, ",", v.Destroy, ",", v.PassLevelId, ",", v.DestroyShow, ",", v.CompleteLimit)
	}
}
