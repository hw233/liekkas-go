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

type CfgQuickPurchaseStamina struct {
	Id    int32
	Count int32
	Cost  int32
}

type CfgQuickPurchaseStaminaConfig struct {
	data map[int32]*CfgQuickPurchaseStamina
}

func NewCfgQuickPurchaseStaminaConfig() *CfgQuickPurchaseStaminaConfig {
	return &CfgQuickPurchaseStaminaConfig{
		data: make(map[int32]*CfgQuickPurchaseStamina),
	}
}

func (c *CfgQuickPurchaseStaminaConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgQuickPurchaseStamina)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgQuickPurchaseStamina.Id field error,value:", vId)
			return false
		}

		/* parse Count field */
		vCount, _ := parse.GetFieldByName(uint32(i), "count")
		var CountRet bool
		data.Count, CountRet = String2Int32(vCount)
		if !CountRet {
			glog.Error("Parse CfgQuickPurchaseStamina.Count field error,value:", vCount)
			return false
		}

		/* parse Cost field */
		vCost, _ := parse.GetFieldByName(uint32(i), "cost")
		var CostRet bool
		data.Cost, CostRet = String2Int32(vCost)
		if !CostRet {
			glog.Error("Parse CfgQuickPurchaseStamina.Cost field error,value:", vCost)
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

func (c *CfgQuickPurchaseStaminaConfig) Clear() {
}

func (c *CfgQuickPurchaseStaminaConfig) Find(id int32) (*CfgQuickPurchaseStamina, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgQuickPurchaseStaminaConfig) GetAllData() map[int32]*CfgQuickPurchaseStamina {
	return c.data
}

func (c *CfgQuickPurchaseStaminaConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Count, ",", v.Cost)
	}
}
