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

type CfgScorePassPhase struct {
	Id             int32
	SeasonId       int32
	StartDayOffset int32
	ScoreItemId    int32
}

type CfgScorePassPhaseConfig struct {
	data map[int32]*CfgScorePassPhase
}

func NewCfgScorePassPhaseConfig() *CfgScorePassPhaseConfig {
	return &CfgScorePassPhaseConfig{
		data: make(map[int32]*CfgScorePassPhase),
	}
}

func (c *CfgScorePassPhaseConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgScorePassPhase)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgScorePassPhase.Id field error,value:", vId)
			return false
		}

		/* parse SeasonId field */
		vSeasonId, _ := parse.GetFieldByName(uint32(i), "seasonId")
		var SeasonIdRet bool
		data.SeasonId, SeasonIdRet = String2Int32(vSeasonId)
		if !SeasonIdRet {
			glog.Error("Parse CfgScorePassPhase.SeasonId field error,value:", vSeasonId)
			return false
		}

		/* parse StartDayOffset field */
		vStartDayOffset, _ := parse.GetFieldByName(uint32(i), "startDayOffset")
		var StartDayOffsetRet bool
		data.StartDayOffset, StartDayOffsetRet = String2Int32(vStartDayOffset)
		if !StartDayOffsetRet {
			glog.Error("Parse CfgScorePassPhase.StartDayOffset field error,value:", vStartDayOffset)
			return false
		}

		/* parse ScoreItemId field */
		vScoreItemId, _ := parse.GetFieldByName(uint32(i), "scoreItemId")
		var ScoreItemIdRet bool
		data.ScoreItemId, ScoreItemIdRet = String2Int32(vScoreItemId)
		if !ScoreItemIdRet {
			glog.Error("Parse CfgScorePassPhase.ScoreItemId field error,value:", vScoreItemId)
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

func (c *CfgScorePassPhaseConfig) Clear() {
}

func (c *CfgScorePassPhaseConfig) Find(id int32) (*CfgScorePassPhase, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgScorePassPhaseConfig) GetAllData() map[int32]*CfgScorePassPhase {
	return c.data
}

func (c *CfgScorePassPhaseConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.SeasonId, ",", v.StartDayOffset, ",", v.ScoreItemId)
	}
}
