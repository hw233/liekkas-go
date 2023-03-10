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

type CfgUltraSkill struct {
	Id         int32
	StartTime  float64
	PlayTime   float64
	EndTime    float64
	PauseTime  float64
	ResumeTime float64
}

type CfgUltraSkillConfig struct {
	data map[int32]*CfgUltraSkill
}

func NewCfgUltraSkillConfig() *CfgUltraSkillConfig {
	return &CfgUltraSkillConfig{
		data: make(map[int32]*CfgUltraSkill),
	}
}

func (c *CfgUltraSkillConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgUltraSkill)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgUltraSkill.Id field error,value:", vId)
			return false
		}

		/* parse StartTime field */
		vStartTime, _ := parse.GetFieldByName(uint32(i), "startTime")
		var StartTimeRet bool
		data.StartTime, StartTimeRet = String2Float(vStartTime)
		if !StartTimeRet {
			glog.Error("Parse CfgUltraSkill.StartTime field error,value:", vStartTime)
		}

		/* parse PlayTime field */
		vPlayTime, _ := parse.GetFieldByName(uint32(i), "playTime")
		var PlayTimeRet bool
		data.PlayTime, PlayTimeRet = String2Float(vPlayTime)
		if !PlayTimeRet {
			glog.Error("Parse CfgUltraSkill.PlayTime field error,value:", vPlayTime)
		}

		/* parse EndTime field */
		vEndTime, _ := parse.GetFieldByName(uint32(i), "endTime")
		var EndTimeRet bool
		data.EndTime, EndTimeRet = String2Float(vEndTime)
		if !EndTimeRet {
			glog.Error("Parse CfgUltraSkill.EndTime field error,value:", vEndTime)
		}

		/* parse PauseTime field */
		vPauseTime, _ := parse.GetFieldByName(uint32(i), "pauseTime")
		var PauseTimeRet bool
		data.PauseTime, PauseTimeRet = String2Float(vPauseTime)
		if !PauseTimeRet {
			glog.Error("Parse CfgUltraSkill.PauseTime field error,value:", vPauseTime)
		}

		/* parse ResumeTime field */
		vResumeTime, _ := parse.GetFieldByName(uint32(i), "resumeTime")
		var ResumeTimeRet bool
		data.ResumeTime, ResumeTimeRet = String2Float(vResumeTime)
		if !ResumeTimeRet {
			glog.Error("Parse CfgUltraSkill.ResumeTime field error,value:", vResumeTime)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgUltraSkillConfig) Clear() {
}

func (c *CfgUltraSkillConfig) Find(id int32) (*CfgUltraSkill, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgUltraSkillConfig) GetAllData() map[int32]*CfgUltraSkill {
	return c.data
}

func (c *CfgUltraSkillConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.StartTime, ",", v.PlayTime, ",", v.EndTime, ",", v.PauseTime, ",", v.ResumeTime)
	}
}
