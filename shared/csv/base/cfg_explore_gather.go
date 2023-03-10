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

type CfgExploreGather struct {
	Id         int32
	GatherTime int32
	GatherName string
	GatherDrop int32
}

type CfgExploreGatherConfig struct {
	data map[int32]*CfgExploreGather
}

func NewCfgExploreGatherConfig() *CfgExploreGatherConfig {
	return &CfgExploreGatherConfig{
		data: make(map[int32]*CfgExploreGather),
	}
}

func (c *CfgExploreGatherConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreGather)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreGather.Id field error,value:", vId)
			return false
		}

		/* parse GatherTime field */
		vGatherTime, _ := parse.GetFieldByName(uint32(i), "gatherTime")
		var GatherTimeRet bool
		data.GatherTime, GatherTimeRet = String2Int32(vGatherTime)
		if !GatherTimeRet {
			glog.Error("Parse CfgExploreGather.GatherTime field error,value:", vGatherTime)
			return false
		}

		/* parse GatherName field */
		data.GatherName, _ = parse.GetFieldByName(uint32(i), "gatherName")

		/* parse GatherDrop field */
		vGatherDrop, _ := parse.GetFieldByName(uint32(i), "gatherDrop")
		var GatherDropRet bool
		data.GatherDrop, GatherDropRet = String2Int32(vGatherDrop)
		if !GatherDropRet {
			glog.Error("Parse CfgExploreGather.GatherDrop field error,value:", vGatherDrop)
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

func (c *CfgExploreGatherConfig) Clear() {
}

func (c *CfgExploreGatherConfig) Find(id int32) (*CfgExploreGather, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreGatherConfig) GetAllData() map[int32]*CfgExploreGather {
	return c.data
}

func (c *CfgExploreGatherConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.GatherTime, ",", v.GatherName, ",", v.GatherDrop)
	}
}
