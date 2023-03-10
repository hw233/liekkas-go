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

type CfgSurveyInfo struct {
	Id              int32
	SurveyType      int32
	StartDate       int32
	DayCnt          int32
	UnlockCondition []string
	QuestionCnt     int32
	MaxLen          int32
	DropID          int32
}

type CfgSurveyInfoConfig struct {
	data map[int32]*CfgSurveyInfo
}

func NewCfgSurveyInfoConfig() *CfgSurveyInfoConfig {
	return &CfgSurveyInfoConfig{
		data: make(map[int32]*CfgSurveyInfo),
	}
}

func (c *CfgSurveyInfoConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgSurveyInfo)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgSurveyInfo.Id field error,value:", vId)
			return false
		}

		/* parse SurveyType field */
		vSurveyType, _ := parse.GetFieldByName(uint32(i), "surveyType")
		var SurveyTypeRet bool
		data.SurveyType, SurveyTypeRet = String2Int32(vSurveyType)
		if !SurveyTypeRet {
			glog.Error("Parse CfgSurveyInfo.SurveyType field error,value:", vSurveyType)
			return false
		}

		/* parse StartDate field */
		vStartDate, _ := parse.GetFieldByName(uint32(i), "startDate")
		var StartDateRet bool
		data.StartDate, StartDateRet = String2Int32(vStartDate)
		if !StartDateRet {
			glog.Error("Parse CfgSurveyInfo.StartDate field error,value:", vStartDate)
			return false
		}

		/* parse DayCnt field */
		vDayCnt, _ := parse.GetFieldByName(uint32(i), "dayCnt")
		var DayCntRet bool
		data.DayCnt, DayCntRet = String2Int32(vDayCnt)
		if !DayCntRet {
			glog.Error("Parse CfgSurveyInfo.DayCnt field error,value:", vDayCnt)
			return false
		}

		/* parse UnlockCondition field */
		vecUnlockCondition, _ := parse.GetFieldByName(uint32(i), "unlockCondition")
		arrayUnlockCondition := strings.Split(vecUnlockCondition, ",")
		for j := 0; j < len(arrayUnlockCondition); j++ {
			v := arrayUnlockCondition[j]
			data.UnlockCondition = append(data.UnlockCondition, v)
		}

		/* parse QuestionCnt field */
		vQuestionCnt, _ := parse.GetFieldByName(uint32(i), "questionCnt")
		var QuestionCntRet bool
		data.QuestionCnt, QuestionCntRet = String2Int32(vQuestionCnt)
		if !QuestionCntRet {
			glog.Error("Parse CfgSurveyInfo.QuestionCnt field error,value:", vQuestionCnt)
			return false
		}

		/* parse MaxLen field */
		vMaxLen, _ := parse.GetFieldByName(uint32(i), "maxLen")
		var MaxLenRet bool
		data.MaxLen, MaxLenRet = String2Int32(vMaxLen)
		if !MaxLenRet {
			glog.Error("Parse CfgSurveyInfo.MaxLen field error,value:", vMaxLen)
			return false
		}

		/* parse DropID field */
		vDropID, _ := parse.GetFieldByName(uint32(i), "dropID")
		var DropIDRet bool
		data.DropID, DropIDRet = String2Int32(vDropID)
		if !DropIDRet {
			glog.Error("Parse CfgSurveyInfo.DropID field error,value:", vDropID)
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

func (c *CfgSurveyInfoConfig) Clear() {
}

func (c *CfgSurveyInfoConfig) Find(id int32) (*CfgSurveyInfo, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgSurveyInfoConfig) GetAllData() map[int32]*CfgSurveyInfo {
	return c.data
}

func (c *CfgSurveyInfoConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.SurveyType, ",", v.StartDate, ",", v.DayCnt, ",", v.UnlockCondition, ",", v.QuestionCnt, ",", v.MaxLen, ",", v.DropID)
	}
}
