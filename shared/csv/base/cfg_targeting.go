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

type CfgTargeting struct {
	Id                    int32
	Comment               string
	TargetGroup           int32
	Shape                 []float64
	HitPointFilter        string
	HitPointPercentFilter string
	Camp                  int32
	Career                int32
	Sex                   int32
	SortType              int32
	SortOrder             int32
	Amount                int32
	NextTarget            int32
}

type CfgTargetingConfig struct {
	data map[int32]*CfgTargeting
}

func NewCfgTargetingConfig() *CfgTargetingConfig {
	return &CfgTargetingConfig{
		data: make(map[int32]*CfgTargeting),
	}
}

func (c *CfgTargetingConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgTargeting)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgTargeting.Id field error,value:", vId)
			return false
		}

		/* parse Comment field */
		data.Comment, _ = parse.GetFieldByName(uint32(i), "comment")

		/* parse TargetGroup field */
		vTargetGroup, _ := parse.GetFieldByName(uint32(i), "targetGroup")
		var TargetGroupRet bool
		data.TargetGroup, TargetGroupRet = String2Int32(vTargetGroup)
		if !TargetGroupRet {
			glog.Error("Parse CfgTargeting.TargetGroup field error,value:", vTargetGroup)
			return false
		}

		/* parse Shape field */
		vecShape, _ := parse.GetFieldByName(uint32(i), "shape")
		arrayShape := strings.Split(vecShape, ",")
		for j := 0; j < len(arrayShape); j++ {
			v, ret := String2Float(arrayShape[j])
			if !ret {
				glog.Error("Parse CfgTargeting.Shape field error,value:", arrayShape[j])
				return false
			}
			data.Shape = append(data.Shape, v)
		}

		/* parse HitPointFilter field */
		data.HitPointFilter, _ = parse.GetFieldByName(uint32(i), "hitPointFilter")

		/* parse HitPointPercentFilter field */
		data.HitPointPercentFilter, _ = parse.GetFieldByName(uint32(i), "hitPointPercentFilter")

		/* parse Camp field */
		vCamp, _ := parse.GetFieldByName(uint32(i), "camp")
		var CampRet bool
		data.Camp, CampRet = String2Int32(vCamp)
		if !CampRet {
			glog.Error("Parse CfgTargeting.Camp field error,value:", vCamp)
			return false
		}

		/* parse Career field */
		vCareer, _ := parse.GetFieldByName(uint32(i), "career")
		var CareerRet bool
		data.Career, CareerRet = String2Int32(vCareer)
		if !CareerRet {
			glog.Error("Parse CfgTargeting.Career field error,value:", vCareer)
			return false
		}

		/* parse Sex field */
		vSex, _ := parse.GetFieldByName(uint32(i), "sex")
		var SexRet bool
		data.Sex, SexRet = String2Int32(vSex)
		if !SexRet {
			glog.Error("Parse CfgTargeting.Sex field error,value:", vSex)
			return false
		}

		/* parse SortType field */
		vSortType, _ := parse.GetFieldByName(uint32(i), "sortType")
		var SortTypeRet bool
		data.SortType, SortTypeRet = String2Int32(vSortType)
		if !SortTypeRet {
			glog.Error("Parse CfgTargeting.SortType field error,value:", vSortType)
			return false
		}

		/* parse SortOrder field */
		vSortOrder, _ := parse.GetFieldByName(uint32(i), "sortOrder")
		var SortOrderRet bool
		data.SortOrder, SortOrderRet = String2Int32(vSortOrder)
		if !SortOrderRet {
			glog.Error("Parse CfgTargeting.SortOrder field error,value:", vSortOrder)
			return false
		}

		/* parse Amount field */
		vAmount, _ := parse.GetFieldByName(uint32(i), "amount")
		var AmountRet bool
		data.Amount, AmountRet = String2Int32(vAmount)
		if !AmountRet {
			glog.Error("Parse CfgTargeting.Amount field error,value:", vAmount)
			return false
		}

		/* parse NextTarget field */
		vNextTarget, _ := parse.GetFieldByName(uint32(i), "nextTarget")
		var NextTargetRet bool
		data.NextTarget, NextTargetRet = String2Int32(vNextTarget)
		if !NextTargetRet {
			glog.Error("Parse CfgTargeting.NextTarget field error,value:", vNextTarget)
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

func (c *CfgTargetingConfig) Clear() {
}

func (c *CfgTargetingConfig) Find(id int32) (*CfgTargeting, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgTargetingConfig) GetAllData() map[int32]*CfgTargeting {
	return c.data
}

func (c *CfgTargetingConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Comment, ",", v.TargetGroup, ",", v.Shape, ",", v.HitPointFilter, ",", v.HitPointPercentFilter, ",", v.Camp, ",", v.Career, ",", v.Sex, ",", v.SortType, ",", v.SortOrder, ",", v.Amount, ",", v.NextTarget)
	}
}
