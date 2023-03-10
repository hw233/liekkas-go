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

type CfgStoreGeneral struct {
	Id              int32
	UpdateRule      int32
	UnlockCondition []string
	Currencies      []int32
	SubStores       []int32
	SubStoreType    int32
}

type CfgStoreGeneralConfig struct {
	data map[int32]*CfgStoreGeneral
}

func NewCfgStoreGeneralConfig() *CfgStoreGeneralConfig {
	return &CfgStoreGeneralConfig{
		data: make(map[int32]*CfgStoreGeneral),
	}
}

func (c *CfgStoreGeneralConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgStoreGeneral)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgStoreGeneral.Id field error,value:", vId)
			return false
		}

		/* parse UpdateRule field */
		vUpdateRule, _ := parse.GetFieldByName(uint32(i), "updateRule")
		var UpdateRuleRet bool
		data.UpdateRule, UpdateRuleRet = String2Int32(vUpdateRule)
		if !UpdateRuleRet {
			glog.Error("Parse CfgStoreGeneral.UpdateRule field error,value:", vUpdateRule)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse Currencies field */
		vecCurrencies, _ := parse.GetFieldByName(uint32(i), "currencies")
		if vecCurrencies != "" {
			arrayCurrencies := strings.Split(vecCurrencies, ",")
			for j := 0; j < len(arrayCurrencies); j++ {
				v, ret := String2Int32(arrayCurrencies[j])
				if !ret {
					glog.Error("Parse CfgStoreGeneral.Currencies field error, value:", arrayCurrencies[j])
					return false
				}
				data.Currencies = append(data.Currencies, v)
			}
		}

		/* parse SubStores field */
		vecSubStores, _ := parse.GetFieldByName(uint32(i), "subStores")
		if vecSubStores != "" {
			arraySubStores := strings.Split(vecSubStores, ",")
			for j := 0; j < len(arraySubStores); j++ {
				v, ret := String2Int32(arraySubStores[j])
				if !ret {
					glog.Error("Parse CfgStoreGeneral.SubStores field error, value:", arraySubStores[j])
					return false
				}
				data.SubStores = append(data.SubStores, v)
			}
		}

		/* parse SubStoreType field */
		vSubStoreType, _ := parse.GetFieldByName(uint32(i), "subStoreType")
		var SubStoreTypeRet bool
		data.SubStoreType, SubStoreTypeRet = String2Int32(vSubStoreType)
		if !SubStoreTypeRet {
			glog.Error("Parse CfgStoreGeneral.SubStoreType field error,value:", vSubStoreType)
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

func (c *CfgStoreGeneralConfig) Clear() {
}

func (c *CfgStoreGeneralConfig) Find(id int32) (*CfgStoreGeneral, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgStoreGeneralConfig) GetAllData() map[int32]*CfgStoreGeneral {
	return c.data
}

func (c *CfgStoreGeneralConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.UpdateRule, ",", v.UnlockCondition, ",", v.Currencies, ",", v.SubStores, ",", v.SubStoreType)
	}
}
