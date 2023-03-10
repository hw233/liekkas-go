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

type CfgHeroTaskData struct {
	Id              int32
	Module          int32
	UnlockCondition []string
	ConditionType   int32
	ConditionParams []int32
	TargetCount     int32
	Cname           string
	IsAutoReceive   int32
	AddPoint        int32
	DropId          int32
	Group           int32
}

type CfgHeroTaskDataConfig struct {
	data map[int32]*CfgHeroTaskData
}

func NewCfgHeroTaskDataConfig() *CfgHeroTaskDataConfig {
	return &CfgHeroTaskDataConfig{
		data: make(map[int32]*CfgHeroTaskData),
	}
}

func (c *CfgHeroTaskDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgHeroTaskData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgHeroTaskData.Id field error,value:", vId)
			return false
		}

		/* parse Module field */
		vModule, _ := parse.GetFieldByName(uint32(i), "module")
		var ModuleRet bool
		data.Module, ModuleRet = String2Int32(vModule)
		if !ModuleRet {
			glog.Error("Parse CfgHeroTaskData.Module field error,value:", vModule)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse ConditionType field */
		vConditionType, _ := parse.GetFieldByName(uint32(i), "conditionType")
		var ConditionTypeRet bool
		data.ConditionType, ConditionTypeRet = String2Int32(vConditionType)
		if !ConditionTypeRet {
			glog.Error("Parse CfgHeroTaskData.ConditionType field error,value:", vConditionType)
			return false
		}

		/* parse ConditionParams field */
		vecConditionParams, _ := parse.GetFieldByName(uint32(i), "conditionParams")
		if vecConditionParams != "" {
			arrayConditionParams := strings.Split(vecConditionParams, ",")
			for j := 0; j < len(arrayConditionParams); j++ {
				v, ret := String2Int32(arrayConditionParams[j])
				if !ret {
					glog.Error("Parse CfgHeroTaskData.ConditionParams field error, value:", arrayConditionParams[j])
					return false
				}
				data.ConditionParams = append(data.ConditionParams, v)
			}
		}

		/* parse TargetCount field */
		vTargetCount, _ := parse.GetFieldByName(uint32(i), "targetCount")
		var TargetCountRet bool
		data.TargetCount, TargetCountRet = String2Int32(vTargetCount)
		if !TargetCountRet {
			glog.Error("Parse CfgHeroTaskData.TargetCount field error,value:", vTargetCount)
			return false
		}

		/* parse Cname field */
		data.Cname, _ = parse.GetFieldByName(uint32(i), "cname")

		/* parse IsAutoReceive field */
		vIsAutoReceive, _ := parse.GetFieldByName(uint32(i), "isAutoReceive")
		var IsAutoReceiveRet bool
		data.IsAutoReceive, IsAutoReceiveRet = String2Int32(vIsAutoReceive)
		if !IsAutoReceiveRet {
			glog.Error("Parse CfgHeroTaskData.IsAutoReceive field error,value:", vIsAutoReceive)
			return false
		}

		/* parse AddPoint field */
		vAddPoint, _ := parse.GetFieldByName(uint32(i), "addPoint")
		var AddPointRet bool
		data.AddPoint, AddPointRet = String2Int32(vAddPoint)
		if !AddPointRet {
			glog.Error("Parse CfgHeroTaskData.AddPoint field error,value:", vAddPoint)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgHeroTaskData.DropId field error,value:", vDropId)
			return false
		}

		/* parse Group field */
		vGroup, _ := parse.GetFieldByName(uint32(i), "group")
		var GroupRet bool
		data.Group, GroupRet = String2Int32(vGroup)
		if !GroupRet {
			glog.Error("Parse CfgHeroTaskData.Group field error,value:", vGroup)
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

func (c *CfgHeroTaskDataConfig) Clear() {
}

func (c *CfgHeroTaskDataConfig) Find(id int32) (*CfgHeroTaskData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgHeroTaskDataConfig) GetAllData() map[int32]*CfgHeroTaskData {
	return c.data
}

func (c *CfgHeroTaskDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Module, ",", v.UnlockCondition, ",", v.ConditionType, ",", v.ConditionParams, ",", v.TargetCount, ",", v.Cname, ",", v.IsAutoReceive, ",", v.AddPoint, ",", v.DropId, ",", v.Group)
	}
}
