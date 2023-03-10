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

type CfgGachaPool struct {
	Id               int32
	Type             int32
	UnlockCondition  []string
	OpenTime         string
	CloseTime        string
	DailyLimit       int32
	SingleConsume    string
	DropN            int32
	DropR            int32
	DropSR           int32
	DropSSR          int32
	DropSafeSSR      int32
	DropUpSSR        int32
	UpGuaranteeCount int32
}

type CfgGachaPoolConfig struct {
	data map[int32]*CfgGachaPool
}

func NewCfgGachaPoolConfig() *CfgGachaPoolConfig {
	return &CfgGachaPoolConfig{
		data: make(map[int32]*CfgGachaPool),
	}
}

func (c *CfgGachaPoolConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgGachaPool)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgGachaPool.Id field error,value:", vId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgGachaPool.Type field error,value:", vType)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse OpenTime field */
		data.OpenTime, _ = parse.GetFieldByName(uint32(i), "openTime")

		/* parse CloseTime field */
		data.CloseTime, _ = parse.GetFieldByName(uint32(i), "closeTime")

		/* parse DailyLimit field */
		vDailyLimit, _ := parse.GetFieldByName(uint32(i), "dailyLimit")
		var DailyLimitRet bool
		data.DailyLimit, DailyLimitRet = String2Int32(vDailyLimit)
		if !DailyLimitRet {
			glog.Error("Parse CfgGachaPool.DailyLimit field error,value:", vDailyLimit)
			return false
		}

		/* parse SingleConsume field */
		data.SingleConsume, _ = parse.GetFieldByName(uint32(i), "singleConsume")

		/* parse DropN field */
		vDropN, _ := parse.GetFieldByName(uint32(i), "dropN")
		var DropNRet bool
		data.DropN, DropNRet = String2Int32(vDropN)
		if !DropNRet {
			glog.Error("Parse CfgGachaPool.DropN field error,value:", vDropN)
			return false
		}

		/* parse DropR field */
		vDropR, _ := parse.GetFieldByName(uint32(i), "dropR")
		var DropRRet bool
		data.DropR, DropRRet = String2Int32(vDropR)
		if !DropRRet {
			glog.Error("Parse CfgGachaPool.DropR field error,value:", vDropR)
			return false
		}

		/* parse DropSR field */
		vDropSR, _ := parse.GetFieldByName(uint32(i), "dropSR")
		var DropSRRet bool
		data.DropSR, DropSRRet = String2Int32(vDropSR)
		if !DropSRRet {
			glog.Error("Parse CfgGachaPool.DropSR field error,value:", vDropSR)
			return false
		}

		/* parse DropSSR field */
		vDropSSR, _ := parse.GetFieldByName(uint32(i), "dropSSR")
		var DropSSRRet bool
		data.DropSSR, DropSSRRet = String2Int32(vDropSSR)
		if !DropSSRRet {
			glog.Error("Parse CfgGachaPool.DropSSR field error,value:", vDropSSR)
			return false
		}

		/* parse DropSafeSSR field */
		vDropSafeSSR, _ := parse.GetFieldByName(uint32(i), "dropSafeSSR")
		var DropSafeSSRRet bool
		data.DropSafeSSR, DropSafeSSRRet = String2Int32(vDropSafeSSR)
		if !DropSafeSSRRet {
			glog.Error("Parse CfgGachaPool.DropSafeSSR field error,value:", vDropSafeSSR)
			return false
		}

		/* parse DropUpSSR field */
		vDropUpSSR, _ := parse.GetFieldByName(uint32(i), "dropUpSSR")
		var DropUpSSRRet bool
		data.DropUpSSR, DropUpSSRRet = String2Int32(vDropUpSSR)
		if !DropUpSSRRet {
			glog.Error("Parse CfgGachaPool.DropUpSSR field error,value:", vDropUpSSR)
			return false
		}

		/* parse UpGuaranteeCount field */
		vUpGuaranteeCount, _ := parse.GetFieldByName(uint32(i), "upGuaranteeCount")
		var UpGuaranteeCountRet bool
		data.UpGuaranteeCount, UpGuaranteeCountRet = String2Int32(vUpGuaranteeCount)
		if !UpGuaranteeCountRet {
			glog.Error("Parse CfgGachaPool.UpGuaranteeCount field error,value:", vUpGuaranteeCount)
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

func (c *CfgGachaPoolConfig) Clear() {
}

func (c *CfgGachaPoolConfig) Find(id int32) (*CfgGachaPool, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgGachaPoolConfig) GetAllData() map[int32]*CfgGachaPool {
	return c.data
}

func (c *CfgGachaPoolConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Type, ",", v.UnlockCondition, ",", v.OpenTime, ",", v.CloseTime, ",", v.DailyLimit, ",", v.SingleConsume, ",", v.DropN, ",", v.DropR, ",", v.DropSR, ",", v.DropSSR, ",", v.DropSafeSSR, ",", v.DropUpSSR, ",", v.UpGuaranteeCount)
	}
}
