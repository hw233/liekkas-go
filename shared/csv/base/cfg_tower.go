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

type CfgTower struct {
	Id             int32
	ShowTowerFloor int32
	Camp           []int32
	GoUpLimit      int32
	ActiveDate     []int32
}

type CfgTowerConfig struct {
	data map[int32]*CfgTower
}

func NewCfgTowerConfig() *CfgTowerConfig {
	return &CfgTowerConfig{
		data: make(map[int32]*CfgTower),
	}
}

func (c *CfgTowerConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgTower)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgTower.Id field error,value:", vId)
			return false
		}

		/* parse ShowTowerFloor field */
		vShowTowerFloor, _ := parse.GetFieldByName(uint32(i), "showTowerFloor")
		var ShowTowerFloorRet bool
		data.ShowTowerFloor, ShowTowerFloorRet = String2Int32(vShowTowerFloor)
		if !ShowTowerFloorRet {
			glog.Error("Parse CfgTower.ShowTowerFloor field error,value:", vShowTowerFloor)
			return false
		}

		/* parse Camp field */
		vecCamp, _ := parse.GetFieldByName(uint32(i), "camp")
		if vecCamp != "" {
			arrayCamp := strings.Split(vecCamp, ",")
			for j := 0; j < len(arrayCamp); j++ {
				v, ret := String2Int32(arrayCamp[j])
				if !ret {
					glog.Error("Parse CfgTower.Camp field error, value:", arrayCamp[j])
					return false
				}
				data.Camp = append(data.Camp, v)
			}
		}

		/* parse GoUpLimit field */
		vGoUpLimit, _ := parse.GetFieldByName(uint32(i), "goUpLimit")
		var GoUpLimitRet bool
		data.GoUpLimit, GoUpLimitRet = String2Int32(vGoUpLimit)
		if !GoUpLimitRet {
			glog.Error("Parse CfgTower.GoUpLimit field error,value:", vGoUpLimit)
			return false
		}

		/* parse ActiveDate field */
		vecActiveDate, _ := parse.GetFieldByName(uint32(i), "activeDate")
		if vecActiveDate != "" {
			arrayActiveDate := strings.Split(vecActiveDate, ",")
			for j := 0; j < len(arrayActiveDate); j++ {
				v, ret := String2Int32(arrayActiveDate[j])
				if !ret {
					glog.Error("Parse CfgTower.ActiveDate field error, value:", arrayActiveDate[j])
					return false
				}
				data.ActiveDate = append(data.ActiveDate, v)
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

func (c *CfgTowerConfig) Clear() {
}

func (c *CfgTowerConfig) Find(id int32) (*CfgTower, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgTowerConfig) GetAllData() map[int32]*CfgTower {
	return c.data
}

func (c *CfgTowerConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ShowTowerFloor, ",", v.Camp, ",", v.GoUpLimit, ",", v.ActiveDate)
	}
}
