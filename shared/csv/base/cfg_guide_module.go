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

type CfgGuideModule struct {
	Id                  int32
	Desc                string
	UnlockCondition     []string
	UnlockConditionOr   []string
	CompleteCondition   []string
	CompleteConditionOr []string
	GuideOrder          int32
	DropId              int32
}

type CfgGuideModuleConfig struct {
	data map[int32]*CfgGuideModule
}

func NewCfgGuideModuleConfig() *CfgGuideModuleConfig {
	return &CfgGuideModuleConfig{
		data: make(map[int32]*CfgGuideModule),
	}
}

func (c *CfgGuideModuleConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGuideModule)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGuideModule.Id field error,value:", vId)
			return false
		}

		/* parse Desc field */
		data.Desc, _ = parse.GetFieldByName(uint32(i), "desc")

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse UnlockConditionOr field */
		vecUnlockConditionOr, _ := parse.GetFieldByName(uint32(i), "unlockConditionOr")
		arrayUnlockConditionOr := strings.Split(vecUnlockConditionOr, ",")
		for j := 0; j < len(arrayUnlockConditionOr); j++ {
			v := arrayUnlockConditionOr[j]
			data.UnlockConditionOr = append(data.UnlockConditionOr, v)
		}

		/* parse CompleteCondition field */
		vecCompleteCondition, _ := parse.GetFieldByName(uint32(i), "completeCondition")
		arrayCompleteCondition := strings.Split(vecCompleteCondition, ",")
		for j := 0; j < len(arrayCompleteCondition); j++ {
			v := arrayCompleteCondition[j]
			data.CompleteCondition = append(data.CompleteCondition, v)
		}

		/* parse CompleteConditionOr field */
		vecCompleteConditionOr, _ := parse.GetFieldByName(uint32(i), "completeConditionOr")
		arrayCompleteConditionOr := strings.Split(vecCompleteConditionOr, ",")
		for j := 0; j < len(arrayCompleteConditionOr); j++ {
			v := arrayCompleteConditionOr[j]
			data.CompleteConditionOr = append(data.CompleteConditionOr, v)
		}

		/* parse GuideOrder field */
		vGuideOrder, _ := parse.GetFieldByName(uint32(i), "guide_order")
		var GuideOrderRet bool
		data.GuideOrder, GuideOrderRet = String2Int32(vGuideOrder)
		if !GuideOrderRet {
			glog.Error("Parse CfgGuideModule.GuideOrder field error,value:", vGuideOrder)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "drop_id")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgGuideModule.DropId field error,value:", vDropId)
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

func (c *CfgGuideModuleConfig) Clear() {
}

func (c *CfgGuideModuleConfig) Find(id int32) (*CfgGuideModule, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGuideModuleConfig) GetAllData() map[int32]*CfgGuideModule {
	return c.data
}

func (c *CfgGuideModuleConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Desc, ",", v.UnlockCondition, ",", v.UnlockConditionOr, ",", v.CompleteCondition, ",", v.CompleteConditionOr, ",", v.GuideOrder, ",", v.DropId)
	}
}
