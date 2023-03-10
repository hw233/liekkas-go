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

type CfgExploreNpc struct {
	Id            int32
	NpcName       string
	NpcOpt1       string
	NpcOpt2       string
	NpcOpt3       string
	NpcOptCost1   []string
	NpcOptCost2   []string
	NpcOptCost3   []string
	NpcOptDropId1 int32
	NpcOptDropId2 int32
	NpcOptDropId3 int32
}

type CfgExploreNpcConfig struct {
	data map[int32]*CfgExploreNpc
}

func NewCfgExploreNpcConfig() *CfgExploreNpcConfig {
	return &CfgExploreNpcConfig{
		data: make(map[int32]*CfgExploreNpc),
	}
}

func (c *CfgExploreNpcConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreNpc)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreNpc.Id field error,value:", vId)
			return false
		}

		/* parse NpcName field */
		data.NpcName, _ = parse.GetFieldByName(uint32(i), "npcName")

		/* parse NpcOpt1 field */
		data.NpcOpt1, _ = parse.GetFieldByName(uint32(i), "npcOpt1")

		/* parse NpcOpt2 field */
		data.NpcOpt2, _ = parse.GetFieldByName(uint32(i), "npcOpt2")

		/* parse NpcOpt3 field */
		data.NpcOpt3, _ = parse.GetFieldByName(uint32(i), "npcOpt3")

		/* parse NpcOptCost1 field */
		vecNpcOptCost1, _ := parse.GetFieldByName(uint32(i), "npcOptCost1")
		arrayNpcOptCost1 := strings.Split(vecNpcOptCost1, ",")
		for j := 0; j < len(arrayNpcOptCost1); j++ {
			v := arrayNpcOptCost1[j]
			data.NpcOptCost1 = append(data.NpcOptCost1, v)
		}

		/* parse NpcOptCost2 field */
		vecNpcOptCost2, _ := parse.GetFieldByName(uint32(i), "npcOptCost2")
		arrayNpcOptCost2 := strings.Split(vecNpcOptCost2, ",")
		for j := 0; j < len(arrayNpcOptCost2); j++ {
			v := arrayNpcOptCost2[j]
			data.NpcOptCost2 = append(data.NpcOptCost2, v)
		}

		/* parse NpcOptCost3 field */
		vecNpcOptCost3, _ := parse.GetFieldByName(uint32(i), "npcOptCost3")
		arrayNpcOptCost3 := strings.Split(vecNpcOptCost3, ",")
		for j := 0; j < len(arrayNpcOptCost3); j++ {
			v := arrayNpcOptCost3[j]
			data.NpcOptCost3 = append(data.NpcOptCost3, v)
		}

		/* parse NpcOptDropId1 field */
		vNpcOptDropId1, _ := parse.GetFieldByName(uint32(i), "npcOptDropId1")
		var NpcOptDropId1Ret bool
		data.NpcOptDropId1, NpcOptDropId1Ret = String2Int32(vNpcOptDropId1)
		if !NpcOptDropId1Ret {
			glog.Error("Parse CfgExploreNpc.NpcOptDropId1 field error,value:", vNpcOptDropId1)
			return false
		}

		/* parse NpcOptDropId2 field */
		vNpcOptDropId2, _ := parse.GetFieldByName(uint32(i), "npcOptDropId2")
		var NpcOptDropId2Ret bool
		data.NpcOptDropId2, NpcOptDropId2Ret = String2Int32(vNpcOptDropId2)
		if !NpcOptDropId2Ret {
			glog.Error("Parse CfgExploreNpc.NpcOptDropId2 field error,value:", vNpcOptDropId2)
			return false
		}

		/* parse NpcOptDropId3 field */
		vNpcOptDropId3, _ := parse.GetFieldByName(uint32(i), "npcOptDropId3")
		var NpcOptDropId3Ret bool
		data.NpcOptDropId3, NpcOptDropId3Ret = String2Int32(vNpcOptDropId3)
		if !NpcOptDropId3Ret {
			glog.Error("Parse CfgExploreNpc.NpcOptDropId3 field error,value:", vNpcOptDropId3)
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

func (c *CfgExploreNpcConfig) Clear() {
}

func (c *CfgExploreNpcConfig) Find(id int32) (*CfgExploreNpc, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreNpcConfig) GetAllData() map[int32]*CfgExploreNpc {
	return c.data
}

func (c *CfgExploreNpcConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.NpcName, ",", v.NpcOpt1, ",", v.NpcOpt2, ",", v.NpcOpt3, ",", v.NpcOptCost1, ",", v.NpcOptCost2, ",", v.NpcOptCost3, ",", v.NpcOptDropId1, ",", v.NpcOptDropId2, ",", v.NpcOptDropId3)
	}
}
