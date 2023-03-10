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

type CfgSigninData struct {
	Id     int32
	Type   int32
	Year   int32
	Month  int32
	Start  int32
	DayCnt int32
	DropID []int32
}

type CfgSigninDataConfig struct {
	data map[int32]*CfgSigninData
}

func NewCfgSigninDataConfig() *CfgSigninDataConfig {
	return &CfgSigninDataConfig{
		data: make(map[int32]*CfgSigninData),
	}
}

func (c *CfgSigninDataConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgSigninData)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgSigninData.Id field error,value:", vId)
			return false
		}

		/* parse Type field */
		vType, _ := parse.GetFieldByName(uint32(i), "type")
		var TypeRet bool
		data.Type, TypeRet = String2Int32(vType)
		if !TypeRet {
			glog.Error("Parse CfgSigninData.Type field error,value:", vType)
			return false
		}

		/* parse Year field */
		vYear, _ := parse.GetFieldByName(uint32(i), "year")
		var YearRet bool
		data.Year, YearRet = String2Int32(vYear)
		if !YearRet {
			glog.Error("Parse CfgSigninData.Year field error,value:", vYear)
			return false
		}

		/* parse Month field */
		vMonth, _ := parse.GetFieldByName(uint32(i), "month")
		var MonthRet bool
		data.Month, MonthRet = String2Int32(vMonth)
		if !MonthRet {
			glog.Error("Parse CfgSigninData.Month field error,value:", vMonth)
			return false
		}

		/* parse Start field */
		vStart, _ := parse.GetFieldByName(uint32(i), "start")
		var StartRet bool
		data.Start, StartRet = String2Int32(vStart)
		if !StartRet {
			glog.Error("Parse CfgSigninData.Start field error,value:", vStart)
			return false
		}

		/* parse DayCnt field */
		vDayCnt, _ := parse.GetFieldByName(uint32(i), "dayCnt")
		var DayCntRet bool
		data.DayCnt, DayCntRet = String2Int32(vDayCnt)
		if !DayCntRet {
			glog.Error("Parse CfgSigninData.DayCnt field error,value:", vDayCnt)
			return false
		}

		/* parse DropID field */
		vecDropID, _ := parse.GetFieldByName(uint32(i), "dropID")
		if vecDropID != "" {
			arrayDropID := strings.Split(vecDropID, ",")
			for j := 0; j < len(arrayDropID); j++ {
				v, ret := String2Int32(arrayDropID[j])
				if !ret {
					glog.Error("Parse CfgSigninData.DropID field error, value:", arrayDropID[j])
					return false
				}
				data.DropID = append(data.DropID, v)
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

func (c *CfgSigninDataConfig) Clear() {
}

func (c *CfgSigninDataConfig) Find(id int32) (*CfgSigninData, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgSigninDataConfig) GetAllData() map[int32]*CfgSigninData {
	return c.data
}

func (c *CfgSigninDataConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Type, ",", v.Year, ",", v.Month, ",", v.Start, ",", v.DayCnt, ",", v.DropID)
	}
}
