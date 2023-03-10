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

type CfgActivity struct {
	Id               int32
	TimeType         int32
	StartTime        string
	EndTime          string
	UnlockConditions []string
	CloseConditions  []string
}

type CfgActivityConfig struct {
	data map[int32]*CfgActivity
}

func NewCfgActivityConfig() *CfgActivityConfig {
	return &CfgActivityConfig{
		data: make(map[int32]*CfgActivity),
	}
}

func (c *CfgActivityConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgActivity)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgActivity.Id field error,value:", vId)
			return false
		}

		/* parse TimeType field */
		vTimeType, _ := parse.GetFieldByName(uint32(i), "timeType")
		var TimeTypeRet bool
		data.TimeType, TimeTypeRet = String2Int32(vTimeType)
		if !TimeTypeRet {
			glog.Error("Parse CfgActivity.TimeType field error,value:", vTimeType)
			return false
		}

		/* parse StartTime field */
		data.StartTime, _ = parse.GetFieldByName(uint32(i), "startTime")

		/* parse EndTime field */
		data.EndTime, _ = parse.GetFieldByName(uint32(i), "endTime")

		/* parse UnlockConditions field */
		vecUnlockConditions, _ := parse.GetFieldByName(uint32(i), "unlockConditions")
		arrayUnlockConditions := strings.Split(vecUnlockConditions, ",")
		for j := 0; j < len(arrayUnlockConditions); j++ {
			v := arrayUnlockConditions[j]
			data.UnlockConditions = append(data.UnlockConditions, v)
		}

		/* parse CloseConditions field */
		vecCloseConditions, _ := parse.GetFieldByName(uint32(i), "closeConditions")
		arrayCloseConditions := strings.Split(vecCloseConditions, ",")
		for j := 0; j < len(arrayCloseConditions); j++ {
			v := arrayCloseConditions[j]
			data.CloseConditions = append(data.CloseConditions, v)
		}

		if _, ok := c.data[data.Id]; ok {
			glog.Errorf("Find %d repeated", data.Id)
			return false
		}
		c.data[data.Id] = data
	}

	return true
}

func (c *CfgActivityConfig) Clear() {
}

func (c *CfgActivityConfig) Find(id int32) (*CfgActivity, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgActivityConfig) GetAllData() map[int32]*CfgActivity {
	return c.data
}

func (c *CfgActivityConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.TimeType, ",", v.StartTime, ",", v.EndTime, ",", v.UnlockConditions, ",", v.CloseConditions)
	}
}
