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

type CfgStoreGoods struct {
	Id              int32
	RewardsIdAndCnt []string
	UnlockCondition []string
	Probability     int32
	Times           int32
	Currencies      []int32
	RealPrice       []int32
}

type CfgStoreGoodsConfig struct {
	data map[int32]*CfgStoreGoods
}

func NewCfgStoreGoodsConfig() *CfgStoreGoodsConfig {
	return &CfgStoreGoodsConfig{
		data: make(map[int32]*CfgStoreGoods),
	}
}

func (c *CfgStoreGoodsConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgStoreGoods)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgStoreGoods.Id field error,value:", vId)
			return false
		}

		/* parse RewardsIdAndCnt field */
		vecRewardsIdAndCnt, _ := parse.GetFieldByName(uint32(i), "rewardsIdAndCnt")
		arrayRewardsIdAndCnt := strings.Split(vecRewardsIdAndCnt, ",")
		for j := 0; j < len(arrayRewardsIdAndCnt); j++ {
			v := arrayRewardsIdAndCnt[j]
			data.RewardsIdAndCnt = append(data.RewardsIdAndCnt, v)
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse Probability field */
		vProbability, _ := parse.GetFieldByName(uint32(i), "probability")
		var ProbabilityRet bool
		data.Probability, ProbabilityRet = String2Int32(vProbability)
		if !ProbabilityRet {
			glog.Error("Parse CfgStoreGoods.Probability field error,value:", vProbability)
			return false
		}

		/* parse Times field */
		vTimes, _ := parse.GetFieldByName(uint32(i), "times")
		var TimesRet bool
		data.Times, TimesRet = String2Int32(vTimes)
		if !TimesRet {
			glog.Error("Parse CfgStoreGoods.Times field error,value:", vTimes)
			return false
		}

		/* parse Currencies field */
		vecCurrencies, _ := parse.GetFieldByName(uint32(i), "currencies")
		if vecCurrencies != "" {
			arrayCurrencies := strings.Split(vecCurrencies, ",")
			for j := 0; j < len(arrayCurrencies); j++ {
				v, ret := String2Int32(arrayCurrencies[j])
				if !ret {
					glog.Error("Parse CfgStoreGoods.Currencies field error, value:", arrayCurrencies[j])
					return false
				}
				data.Currencies = append(data.Currencies, v)
			}
		}

		/* parse RealPrice field */
		vecRealPrice, _ := parse.GetFieldByName(uint32(i), "realPrice")
		if vecRealPrice != "" {
			arrayRealPrice := strings.Split(vecRealPrice, ",")
			for j := 0; j < len(arrayRealPrice); j++ {
				v, ret := String2Int32(arrayRealPrice[j])
				if !ret {
					glog.Error("Parse CfgStoreGoods.RealPrice field error, value:", arrayRealPrice[j])
					return false
				}
				data.RealPrice = append(data.RealPrice, v)
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

func (c *CfgStoreGoodsConfig) Clear() {
}

func (c *CfgStoreGoodsConfig) Find(id int32) (*CfgStoreGoods, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgStoreGoodsConfig) GetAllData() map[int32]*CfgStoreGoods {
	return c.data
}

func (c *CfgStoreGoodsConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.RewardsIdAndCnt, ",", v.UnlockCondition, ",", v.Probability, ",", v.Times, ",", v.Currencies, ",", v.RealPrice)
	}
}
