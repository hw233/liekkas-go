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

type CfgTowerStage struct {
	Id      int32
	TowerId int32
	Stage   int32
	LevelId int32
}

type CfgTowerStageConfig struct {
	data map[int32]*CfgTowerStage
}

func NewCfgTowerStageConfig() *CfgTowerStageConfig {
	return &CfgTowerStageConfig{
		data: make(map[int32]*CfgTowerStage),
	}
}

func (c *CfgTowerStageConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgTowerStage)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgTowerStage.Id field error,value:", vId)
			return false
		}

		/* parse TowerId field */
		vTowerId, _ := parse.GetFieldByName(uint32(i), "towerId")
		var TowerIdRet bool
		data.TowerId, TowerIdRet = String2Int32(vTowerId)
		if !TowerIdRet {
			glog.Error("Parse CfgTowerStage.TowerId field error,value:", vTowerId)
			return false
		}

		/* parse Stage field */
		vStage, _ := parse.GetFieldByName(uint32(i), "stage")
		var StageRet bool
		data.Stage, StageRet = String2Int32(vStage)
		if !StageRet {
			glog.Error("Parse CfgTowerStage.Stage field error,value:", vStage)
			return false
		}

		/* parse LevelId field */
		vLevelId, _ := parse.GetFieldByName(uint32(i), "levelId")
		var LevelIdRet bool
		data.LevelId, LevelIdRet = String2Int32(vLevelId)
		if !LevelIdRet {
			glog.Error("Parse CfgTowerStage.LevelId field error,value:", vLevelId)
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

func (c *CfgTowerStageConfig) Clear() {
}

func (c *CfgTowerStageConfig) Find(id int32) (*CfgTowerStage, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgTowerStageConfig) GetAllData() map[int32]*CfgTowerStage {
	return c.data
}

func (c *CfgTowerStageConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.TowerId, ",", v.Stage, ",", v.LevelId)
	}
}
