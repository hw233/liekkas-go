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

type CfgExploreResource struct {
	Id              int32
	ResourceName    string
	ResourceMonster []string
	Monstermap      int32
	Time            int32
	Drop            int32
	RefreshLimit    int32
}

type CfgExploreResourceConfig struct {
	data map[int32]*CfgExploreResource
}

func NewCfgExploreResourceConfig() *CfgExploreResourceConfig {
	return &CfgExploreResourceConfig{
		data: make(map[int32]*CfgExploreResource),
	}
}

func (c *CfgExploreResourceConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreResource)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreResource.Id field error,value:", vId)
			return false
		}

		/* parse ResourceName field */
		data.ResourceName, _ = parse.GetFieldByName(uint32(i), "resourceName")

		/* parse ResourceMonster field */
		vecResourceMonster, _ := parse.GetFieldByName(uint32(i), "resourceMonster")
		arrayResourceMonster := strings.Split(vecResourceMonster, ",")
		for j := 0; j < len(arrayResourceMonster); j++ {
			v := arrayResourceMonster[j]
			data.ResourceMonster = append(data.ResourceMonster, v)
		}

		/* parse Monstermap field */
		vMonstermap, _ := parse.GetFieldByName(uint32(i), "monstermap")
		var MonstermapRet bool
		data.Monstermap, MonstermapRet = String2Int32(vMonstermap)
		if !MonstermapRet {
			glog.Error("Parse CfgExploreResource.Monstermap field error,value:", vMonstermap)
			return false
		}

		/* parse Time field */
		vTime, _ := parse.GetFieldByName(uint32(i), "time")
		var TimeRet bool
		data.Time, TimeRet = String2Int32(vTime)
		if !TimeRet {
			glog.Error("Parse CfgExploreResource.Time field error,value:", vTime)
			return false
		}

		/* parse Drop field */
		vDrop, _ := parse.GetFieldByName(uint32(i), "drop")
		var DropRet bool
		data.Drop, DropRet = String2Int32(vDrop)
		if !DropRet {
			glog.Error("Parse CfgExploreResource.Drop field error,value:", vDrop)
			return false
		}

		/* parse RefreshLimit field */
		vRefreshLimit, _ := parse.GetFieldByName(uint32(i), "refreshLimit")
		var RefreshLimitRet bool
		data.RefreshLimit, RefreshLimitRet = String2Int32(vRefreshLimit)
		if !RefreshLimitRet {
			glog.Error("Parse CfgExploreResource.RefreshLimit field error,value:", vRefreshLimit)
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

func (c *CfgExploreResourceConfig) Clear() {
}

func (c *CfgExploreResourceConfig) Find(id int32) (*CfgExploreResource, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreResourceConfig) GetAllData() map[int32]*CfgExploreResource {
	return c.data
}

func (c *CfgExploreResourceConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.ResourceName, ",", v.ResourceMonster, ",", v.Monstermap, ",", v.Time, ",", v.Drop, ",", v.RefreshLimit)
	}
}
