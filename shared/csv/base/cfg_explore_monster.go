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

type CfgExploreMonster struct {
	Id         int32
	Monstermap int32
}

type CfgExploreMonsterConfig struct {
	data map[int32]*CfgExploreMonster
}

func NewCfgExploreMonsterConfig() *CfgExploreMonsterConfig {
	return &CfgExploreMonsterConfig{
		data: make(map[int32]*CfgExploreMonster),
	}
}

func (c *CfgExploreMonsterConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgExploreMonster)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgExploreMonster.Id field error,value:", vId)
			return false
		}

		/* parse Monstermap field */
		vMonstermap, _ := parse.GetFieldByName(uint32(i), "monstermap")
		var MonstermapRet bool
		data.Monstermap, MonstermapRet = String2Int32(vMonstermap)
		if !MonstermapRet {
			glog.Error("Parse CfgExploreMonster.Monstermap field error,value:", vMonstermap)
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

func (c *CfgExploreMonsterConfig) Clear() {
}

func (c *CfgExploreMonsterConfig) Find(id int32) (*CfgExploreMonster, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgExploreMonsterConfig) GetAllData() map[int32]*CfgExploreMonster {
	return c.data
}

func (c *CfgExploreMonsterConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.Monstermap)
	}
}
