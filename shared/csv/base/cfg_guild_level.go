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

type CfgGuildLevel struct {
	Id                              int32
	Exp                             int32
	SingleHelpAddIntimacy           int32
	SingleHelpAddGold               int32
	SingleHelpAddActivation         int32
	DailyGoldUpperLimitByHelp       int32
	DailyActivationUpperLimitByHelp int32
}

type CfgGuildLevelConfig struct {
	data map[int32]*CfgGuildLevel
}

func NewCfgGuildLevelConfig() *CfgGuildLevelConfig {
	return &CfgGuildLevelConfig{
		data: make(map[int32]*CfgGuildLevel),
	}
}

func (c *CfgGuildLevelConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGuildLevel)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGuildLevel.Id field error,value:", vId)
			return false
		}

		/* parse Exp field */
		vExp, _ := parse.GetFieldByName(uint32(i), "exp")
		var ExpRet bool
		data.Exp, ExpRet = String2Int32(vExp)
		if !ExpRet {
			glog.Error("Parse CfgGuildLevel.Exp field error,value:", vExp)
			return false
		}

		/* parse SingleHelpAddIntimacy field */
		vSingleHelpAddIntimacy, _ := parse.GetFieldByName(uint32(i), "singleHelpAddIntimacy")
		var SingleHelpAddIntimacyRet bool
		data.SingleHelpAddIntimacy, SingleHelpAddIntimacyRet = String2Int32(vSingleHelpAddIntimacy)
		if !SingleHelpAddIntimacyRet {
			glog.Error("Parse CfgGuildLevel.SingleHelpAddIntimacy field error,value:", vSingleHelpAddIntimacy)
			return false
		}

		/* parse SingleHelpAddGold field */
		vSingleHelpAddGold, _ := parse.GetFieldByName(uint32(i), "singleHelpAddGold")
		var SingleHelpAddGoldRet bool
		data.SingleHelpAddGold, SingleHelpAddGoldRet = String2Int32(vSingleHelpAddGold)
		if !SingleHelpAddGoldRet {
			glog.Error("Parse CfgGuildLevel.SingleHelpAddGold field error,value:", vSingleHelpAddGold)
			return false
		}

		/* parse SingleHelpAddActivation field */
		vSingleHelpAddActivation, _ := parse.GetFieldByName(uint32(i), "singleHelpAddActivation")
		var SingleHelpAddActivationRet bool
		data.SingleHelpAddActivation, SingleHelpAddActivationRet = String2Int32(vSingleHelpAddActivation)
		if !SingleHelpAddActivationRet {
			glog.Error("Parse CfgGuildLevel.SingleHelpAddActivation field error,value:", vSingleHelpAddActivation)
			return false
		}

		/* parse DailyGoldUpperLimitByHelp field */
		vDailyGoldUpperLimitByHelp, _ := parse.GetFieldByName(uint32(i), "dailyGoldUpperLimitByHelp")
		var DailyGoldUpperLimitByHelpRet bool
		data.DailyGoldUpperLimitByHelp, DailyGoldUpperLimitByHelpRet = String2Int32(vDailyGoldUpperLimitByHelp)
		if !DailyGoldUpperLimitByHelpRet {
			glog.Error("Parse CfgGuildLevel.DailyGoldUpperLimitByHelp field error,value:", vDailyGoldUpperLimitByHelp)
			return false
		}

		/* parse DailyActivationUpperLimitByHelp field */
		vDailyActivationUpperLimitByHelp, _ := parse.GetFieldByName(uint32(i), "dailyActivationUpperLimitByHelp")
		var DailyActivationUpperLimitByHelpRet bool
		data.DailyActivationUpperLimitByHelp, DailyActivationUpperLimitByHelpRet = String2Int32(vDailyActivationUpperLimitByHelp)
		if !DailyActivationUpperLimitByHelpRet {
			glog.Error("Parse CfgGuildLevel.DailyActivationUpperLimitByHelp field error,value:", vDailyActivationUpperLimitByHelp)
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

func (c *CfgGuildLevelConfig) Clear() {
}

func (c *CfgGuildLevelConfig) Find(id int32) (*CfgGuildLevel, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGuildLevelConfig) GetAllData() map[int32]*CfgGuildLevel {
	return c.data
}

func (c *CfgGuildLevelConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Exp, ",", v.SingleHelpAddIntimacy, ",", v.SingleHelpAddGold, ",", v.SingleHelpAddActivation, ",", v.DailyGoldUpperLimitByHelp, ",", v.DailyActivationUpperLimitByHelp)
	}
}
