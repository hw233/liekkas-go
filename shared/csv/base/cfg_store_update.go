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

type CfgStoreUpdate struct {
	Id         int32
	Period     int32
	TimesLimit int32
	Currency   int32
	Cnt        []int32
}

type CfgStoreUpdateConfig struct {
	data map[int32]*CfgStoreUpdate
}

func NewCfgStoreUpdateConfig() *CfgStoreUpdateConfig {
	return &CfgStoreUpdateConfig{
		data: make(map[int32]*CfgStoreUpdate),
	}
}

func (c *CfgStoreUpdateConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgStoreUpdate)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgStoreUpdate.Id field error,value:", vId)
			return false
		}

		/* parse Period field */
		vPeriod, _ := parse.GetFieldByName(uint32(i), "period")
		var PeriodRet bool
		data.Period, PeriodRet = String2Int32(vPeriod)
		if !PeriodRet {
			glog.Error("Parse CfgStoreUpdate.Period field error,value:", vPeriod)
			return false
		}

		/* parse TimesLimit field */
		vTimesLimit, _ := parse.GetFieldByName(uint32(i), "timesLimit")
		var TimesLimitRet bool
		data.TimesLimit, TimesLimitRet = String2Int32(vTimesLimit)
		if !TimesLimitRet {
			glog.Error("Parse CfgStoreUpdate.TimesLimit field error,value:", vTimesLimit)
			return false
		}

		/* parse Currency field */
		vCurrency, _ := parse.GetFieldByName(uint32(i), "currency")
		var CurrencyRet bool
		data.Currency, CurrencyRet = String2Int32(vCurrency)
		if !CurrencyRet {
			glog.Error("Parse CfgStoreUpdate.Currency field error,value:", vCurrency)
			return false
		}

		/* parse Cnt field */
		vecCnt, _ := parse.GetFieldByName(uint32(i), "cnt")
		if vecCnt != "" {
			arrayCnt := strings.Split(vecCnt, ",")
			for j := 0; j < len(arrayCnt); j++ {
				v, ret := String2Int32(arrayCnt[j])
				if !ret {
					glog.Error("Parse CfgStoreUpdate.Cnt field error, value:", arrayCnt[j])
					return false
				}
				data.Cnt = append(data.Cnt, v)
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

func (c *CfgStoreUpdateConfig) Clear() {
}

func (c *CfgStoreUpdateConfig) Find(id int32) (*CfgStoreUpdate, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgStoreUpdateConfig) GetAllData() map[int32]*CfgStoreUpdate {
	return c.data
}

func (c *CfgStoreUpdateConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Period, ",", v.TimesLimit, ",", v.Currency, ",", v.Cnt)
	}
}
