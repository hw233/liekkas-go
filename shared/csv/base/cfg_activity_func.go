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

type CfgActivityFunc struct {
	Id              int32
	ActivityId      int32
	FuncType        int32
	FuncArgs        []int32
	StartTimeOffset []int32
	EndTimeOffset   []int32
	EndItemRemove   []int32
}

type CfgActivityFuncConfig struct {
	data map[int32]*CfgActivityFunc
}

func NewCfgActivityFuncConfig() *CfgActivityFuncConfig {
	return &CfgActivityFuncConfig{
		data: make(map[int32]*CfgActivityFunc),
	}
}

func (c *CfgActivityFuncConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgActivityFunc)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgActivityFunc.Id field error,value:", vId)
			return false
		}

		/* parse ActivityId field */
		vActivityId, _ := parse.GetFieldByName(uint32(i), "activityId")
		var ActivityIdRet bool
		data.ActivityId, ActivityIdRet = String2Int32(vActivityId)
		if !ActivityIdRet {
			glog.Error("Parse CfgActivityFunc.ActivityId field error,value:", vActivityId)
			return false
		}

		/* parse FuncType field */
		vFuncType, _ := parse.GetFieldByName(uint32(i), "funcType")
		var FuncTypeRet bool
		data.FuncType, FuncTypeRet = String2Int32(vFuncType)
		if !FuncTypeRet {
			glog.Error("Parse CfgActivityFunc.FuncType field error,value:", vFuncType)
			return false
		}

		/* parse FuncArgs field */
		vecFuncArgs, _ := parse.GetFieldByName(uint32(i), "funcArgs")
		if vecFuncArgs != "" {
			arrayFuncArgs := strings.Split(vecFuncArgs, ",")
			for j := 0; j < len(arrayFuncArgs); j++ {
				v, ret := String2Int32(arrayFuncArgs[j])
				if !ret {
					glog.Error("Parse CfgActivityFunc.FuncArgs field error, value:", arrayFuncArgs[j])
					return false
				}
				data.FuncArgs = append(data.FuncArgs, v)
			}
		}

		/* parse StartTimeOffset field */
		vecStartTimeOffset, _ := parse.GetFieldByName(uint32(i), "startTimeOffset")
		if vecStartTimeOffset != "" {
			arrayStartTimeOffset := strings.Split(vecStartTimeOffset, ",")
			for j := 0; j < len(arrayStartTimeOffset); j++ {
				v, ret := String2Int32(arrayStartTimeOffset[j])
				if !ret {
					glog.Error("Parse CfgActivityFunc.StartTimeOffset field error, value:", arrayStartTimeOffset[j])
					return false
				}
				data.StartTimeOffset = append(data.StartTimeOffset, v)
			}
		}

		/* parse EndTimeOffset field */
		vecEndTimeOffset, _ := parse.GetFieldByName(uint32(i), "endTimeOffset")
		if vecEndTimeOffset != "" {
			arrayEndTimeOffset := strings.Split(vecEndTimeOffset, ",")
			for j := 0; j < len(arrayEndTimeOffset); j++ {
				v, ret := String2Int32(arrayEndTimeOffset[j])
				if !ret {
					glog.Error("Parse CfgActivityFunc.EndTimeOffset field error, value:", arrayEndTimeOffset[j])
					return false
				}
				data.EndTimeOffset = append(data.EndTimeOffset, v)
			}
		}

		/* parse EndItemRemove field */
		vecEndItemRemove, _ := parse.GetFieldByName(uint32(i), "endItemRemove")
		if vecEndItemRemove != "" {
			arrayEndItemRemove := strings.Split(vecEndItemRemove, ",")
			for j := 0; j < len(arrayEndItemRemove); j++ {
				v, ret := String2Int32(arrayEndItemRemove[j])
				if !ret {
					glog.Error("Parse CfgActivityFunc.EndItemRemove field error, value:", arrayEndItemRemove[j])
					return false
				}
				data.EndItemRemove = append(data.EndItemRemove, v)
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

func (c *CfgActivityFuncConfig) Clear() {
}

func (c *CfgActivityFuncConfig) Find(id int32) (*CfgActivityFunc, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgActivityFuncConfig) GetAllData() map[int32]*CfgActivityFunc {
	return c.data
}

func (c *CfgActivityFuncConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ActivityId, ",", v.FuncType, ",", v.FuncArgs, ",", v.StartTimeOffset, ",", v.EndTimeOffset, ",", v.EndItemRemove)
	}
}
