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

type CfgYggdrasilDailyMonster struct {
	Id          int32
	ObjectId    int32
	PosX        int32
	PosY        int32
	Radius      int32
	IntervalDay int32
}

type CfgYggdrasilDailyMonsterConfig struct {
	data map[int32]*CfgYggdrasilDailyMonster
}

func NewCfgYggdrasilDailyMonsterConfig() *CfgYggdrasilDailyMonsterConfig {
	return &CfgYggdrasilDailyMonsterConfig{
		data: make(map[int32]*CfgYggdrasilDailyMonster),
	}
}

func (c *CfgYggdrasilDailyMonsterConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgYggdrasilDailyMonster)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgYggdrasilDailyMonster.Id field error,value:", vId)
			return false
		}

		/* parse ObjectId field */
		vObjectId, _ := parse.GetFieldByName(uint32(i), "objectId")
		var ObjectIdRet bool
		data.ObjectId, ObjectIdRet = String2Int32(vObjectId)
		if !ObjectIdRet {
			glog.Error("Parse CfgYggdrasilDailyMonster.ObjectId field error,value:", vObjectId)
			return false
		}

		/* parse PosX field */
		vPosX, _ := parse.GetFieldByName(uint32(i), "posX")
		var PosXRet bool
		data.PosX, PosXRet = String2Int32(vPosX)
		if !PosXRet {
			glog.Error("Parse CfgYggdrasilDailyMonster.PosX field error,value:", vPosX)
			return false
		}

		/* parse PosY field */
		vPosY, _ := parse.GetFieldByName(uint32(i), "posY")
		var PosYRet bool
		data.PosY, PosYRet = String2Int32(vPosY)
		if !PosYRet {
			glog.Error("Parse CfgYggdrasilDailyMonster.PosY field error,value:", vPosY)
			return false
		}

		/* parse Radius field */
		vRadius, _ := parse.GetFieldByName(uint32(i), "radius")
		var RadiusRet bool
		data.Radius, RadiusRet = String2Int32(vRadius)
		if !RadiusRet {
			glog.Error("Parse CfgYggdrasilDailyMonster.Radius field error,value:", vRadius)
			return false
		}

		/* parse IntervalDay field */
		vIntervalDay, _ := parse.GetFieldByName(uint32(i), "intervalDay")
		var IntervalDayRet bool
		data.IntervalDay, IntervalDayRet = String2Int32(vIntervalDay)
		if !IntervalDayRet {
			glog.Error("Parse CfgYggdrasilDailyMonster.IntervalDay field error,value:", vIntervalDay)
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

func (c *CfgYggdrasilDailyMonsterConfig) Clear() {
}

func (c *CfgYggdrasilDailyMonsterConfig) Find(id int32) (*CfgYggdrasilDailyMonster, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgYggdrasilDailyMonsterConfig) GetAllData() map[int32]*CfgYggdrasilDailyMonster {
	return c.data
}

func (c *CfgYggdrasilDailyMonsterConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ObjectId, ",", v.PosX, ",", v.PosY, ",", v.Radius, ",", v.IntervalDay)
	}
}
