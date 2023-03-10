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

type CfgSurveyQuestion struct {
	Id           int32
	SurveyID     int32
	QuestionType int32
}

type CfgSurveyQuestionConfig struct {
	data map[int32]*CfgSurveyQuestion
}

func NewCfgSurveyQuestionConfig() *CfgSurveyQuestionConfig {
	return &CfgSurveyQuestionConfig{
		data: make(map[int32]*CfgSurveyQuestion),
	}
}

func (c *CfgSurveyQuestionConfig) Load(filePath string) bool {
	parse := NewParser()
	if err := parse.Load(filePath, true); err != nil {
		glog.Info("Load", filePath, "err: ", err)
		return false
	}

	// iterator all lines' content
	for i := 2; i < parse.GetAllCount(); i++ {
		data := new(CfgSurveyQuestion)

		/* parse Id field */
		vId, _ := parse.GetFieldByName(uint32(i), "id")
		var IdRet bool
		data.Id, IdRet = String2Int32(vId)
		if !IdRet {
			glog.Error("Parse CfgSurveyQuestion.Id field error,value:", vId)
			return false
		}

		/* parse SurveyID field */
		vSurveyID, _ := parse.GetFieldByName(uint32(i), "surveyID")
		var SurveyIDRet bool
		data.SurveyID, SurveyIDRet = String2Int32(vSurveyID)
		if !SurveyIDRet {
			glog.Error("Parse CfgSurveyQuestion.SurveyID field error,value:", vSurveyID)
			return false
		}

		/* parse QuestionType field */
		vQuestionType, _ := parse.GetFieldByName(uint32(i), "questionType")
		var QuestionTypeRet bool
		data.QuestionType, QuestionTypeRet = String2Int32(vQuestionType)
		if !QuestionTypeRet {
			glog.Error("Parse CfgSurveyQuestion.QuestionType field error,value:", vQuestionType)
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

func (c *CfgSurveyQuestionConfig) Clear() {
}

func (c *CfgSurveyQuestionConfig) Find(id int32) (*CfgSurveyQuestion, bool) {
	v, ok := c.data[id]
	return v, ok
}

func (c *CfgSurveyQuestionConfig) GetAllData() map[int32]*CfgSurveyQuestion {
	return c.data
}

func (c *CfgSurveyQuestionConfig) Traverse() {
	for _, v := range c.data {
		glog.Info(v.Id, ",", v.SurveyID, ",", v.QuestionType)
	}
}
