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

type CfgScorePassReward struct {
	Id      int32
	PhaseId int32
	Score   int32
	DropId  int32
}

type CfgScorePassRewardConfig struct {
	data map[int32]*CfgScorePassReward
}

func NewCfgScorePassRewardConfig() *CfgScorePassRewardConfig {
	return &CfgScorePassRewardConfig{
		data: make(map[int32]*CfgScorePassReward),
	}
}

func (c *CfgScorePassRewardConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgScorePassReward)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgScorePassReward.Id field error,value:", vId)
			return false
		}

		/* parse PhaseId field */
		vPhaseId, _ := parse.GetFieldByName(uint32(i), "phaseId")
		var PhaseIdRet bool
		data.PhaseId, PhaseIdRet = String2Int32(vPhaseId)
		if !PhaseIdRet {
			glog.Error("Parse CfgScorePassReward.PhaseId field error,value:", vPhaseId)
			return false
		}

		/* parse Score field */
		vScore, _ := parse.GetFieldByName(uint32(i), "score")
		var ScoreRet bool
		data.Score, ScoreRet = String2Int32(vScore)
		if !ScoreRet {
			glog.Error("Parse CfgScorePassReward.Score field error,value:", vScore)
			return false
		}

		/* parse DropId field */
		vDropId, _ := parse.GetFieldByName(uint32(i), "dropId")
		var DropIdRet bool
		data.DropId, DropIdRet = String2Int32(vDropId)
		if !DropIdRet {
			glog.Error("Parse CfgScorePassReward.DropId field error,value:", vDropId)
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

func (c *CfgScorePassRewardConfig) Clear() {
}

func (c *CfgScorePassRewardConfig) Find(id int32) (*CfgScorePassReward, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgScorePassRewardConfig) GetAllData() map[int32]*CfgScorePassReward {
	return c.data
}

func (c *CfgScorePassRewardConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.PhaseId, ",", v.Score, ",", v.DropId)
	}
}
