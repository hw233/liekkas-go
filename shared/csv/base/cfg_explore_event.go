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

type CfgExploreEvent struct {
	Id               int32
	EventType        int32
	EventParam       int32
	UnlockConditions []string
	RefreshTime      int32
	RefreshLimit     int32
}

type CfgExploreEventConfig struct {
	data map[int32]*CfgExploreEvent
}

func NewCfgExploreEventConfig() *CfgExploreEventConfig {
	return &CfgExploreEventConfig{
		data: make(map[int32]*CfgExploreEvent),
	}
}

func (c *CfgExploreEventConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreEvent)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreEvent.Id field error,value:", vId)
			return false
		}

		/* parse EventType field */
		vEventType, _ := parse.GetFieldByName(uint32(i), "eventType")
		var EventTypeRet bool
		data.EventType, EventTypeRet = String2Int32(vEventType)
		if !EventTypeRet {
			glog.Error("Parse CfgExploreEvent.EventType field error,value:", vEventType)
			return false
		}

		/* parse EventParam field */
		vEventParam, _ := parse.GetFieldByName(uint32(i), "eventParam")
		var EventParamRet bool
		data.EventParam, EventParamRet = String2Int32(vEventParam)
		if !EventParamRet {
			glog.Error("Parse CfgExploreEvent.EventParam field error,value:", vEventParam)
			return false
		}

		/* parse UnlockConditions field */
		vecUnlockConditions, _ := parse.GetFieldByName(uint32(i), "unlockConditions")
		arrayUnlockConditions := strings.Split(vecUnlockConditions, ",")
		for j := 0; j < len(arrayUnlockConditions); j++ {
			v := arrayUnlockConditions[j]
			data.UnlockConditions = append(data.UnlockConditions, v)
		}

		/* parse RefreshTime field */
		vRefreshTime, _ := parse.GetFieldByName(uint32(i), "refreshTime")
		var RefreshTimeRet bool
		data.RefreshTime, RefreshTimeRet = String2Int32(vRefreshTime)
		if !RefreshTimeRet {
			glog.Error("Parse CfgExploreEvent.RefreshTime field error,value:", vRefreshTime)
			return false
		}

		/* parse RefreshLimit field */
		vRefreshLimit, _ := parse.GetFieldByName(uint32(i), "refreshLimit")
		var RefreshLimitRet bool
		data.RefreshLimit, RefreshLimitRet = String2Int32(vRefreshLimit)
		if !RefreshLimitRet {
			glog.Error("Parse CfgExploreEvent.RefreshLimit field error,value:", vRefreshLimit)
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

func (c *CfgExploreEventConfig) Clear() {
}

func (c *CfgExploreEventConfig) Find(id int32) (*CfgExploreEvent, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreEventConfig) GetAllData() map[int32]*CfgExploreEvent {
	return c.data
}

func (c *CfgExploreEventConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.EventType, ",", v.EventParam, ",", v.UnlockConditions, ",", v.RefreshTime, ",", v.RefreshLimit)
	}
}
