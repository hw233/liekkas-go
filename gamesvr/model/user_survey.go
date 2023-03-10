package model

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/servertime"
	"time"
)

// 对回答的长度和数量进行检查，没问题的话就调用函数保存survey数据
func (u *User) SaveSurvey(ctx context.Context, id int32, answers []*pb.VOSurveyAnswer) error {

	surveyData, err := manager.CSV.Survey.GetInfo(id)
	if err != nil {
		return err
	}
	if id > 31 {
		return errors.Swrapf(common.ErrSurveyCntOutOfBitNumber, id)
	}
	if u.Info.SurveyRecord.IsMarked(int(id)) {
		return errors.Swrapf(common.ErrSurveyDuplicateSurveyData, id)
	}

	numOfQuestion := surveyData.QuestionCnt
	maxLen := surveyData.MaxLen
	if len(answers) > int(numOfQuestion) {
		return errors.Swrapf(common.ErrSurveyQuestionCntBeyondLimit, id)
	}

	survey := NewSurvey()

	survey.SetSurvey(u.ID, surveyData.Id)
	for _, answer := range answers {
		if len(answer.Answer) > int(maxLen) {
			return errors.Swrapf(common.ErrSurveyAnswerTextBeyondLimit, id)
		}
		survey.AddAnswer(answer.Qid, answer.Answer)
	}
	err = survey.NewSurveyData(ctx) // 将survey写入数据库
	if err != nil {
		return err
	}

	// 发放奖励
	reason := logreason.NewReason(logreason.Questionnaire)
	_, err = u.AddRewardsByDropId(surveyData.DropID, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}

	u.Info.SurveyRecord.Mark(int(id))
	return nil
}

// 对所给id对应的survey进行检查，如果未达到解锁条件或者问卷已过期则返回错误,否则返回nil
func (u *User) CheckForSurveyType(id int32) error {

	surveyData, err := manager.CSV.Survey.GetInfo(id)
	if err != nil {
		return err
	}
	surveyType := surveyData.SurveyType
	if surveyType == 1 {
		err = u.checkSurveyUnlockCondition(id)
		if err != nil {
			return err
		}
	} else if surveyType == 2 {
		err = checkSurveyTimeLimit(id, surveyData.StartDate, surveyData.DayCnt)
		if err != nil {
			return err
		}
	} else {

		err = u.checkSurveyUnlockCondition(id)
		if err != nil {
			return err
		}

		err = checkSurveyTimeLimit(id, surveyData.StartDate, surveyData.DayCnt)
		if err != nil {
			return err
		}
	}

	return nil

}

//----------------------------------------
//SurveyUnlockCondition
//----------------------------------------
func (u *User) checkSurveyUnlockCondition(surveyID int32) error {
	surveyCfg, err := manager.CSV.Survey.GetInfo(surveyID)
	if err != nil {
		return err
	}

	if u.CheckUserConditions(surveyCfg.UnlockCondition) != nil {
		return errors.Swrapf(common.ErrSurveyUnlockConditionNotSatisfied, surveyID)
	}

	return nil
}

//----------------------------------------
//SurveyTimeLimit
//----------------------------------------
func checkSurveyTimeLimit(id int32, startDate int32, dayCnt int32) error {
	now := servertime.Now()
	start := DailyRefreshTime(int2Time(startDate))
	end := start.AddDate(0, 0, int(dayCnt))

	if now.Before(start) || now.After(end) {
		return errors.Swrapf(common.ErrSurveyValidityPeriodBeyondLimit, id)
	}

	return nil
}

// 例：20211026 -> 时间2021年10月26日0点0分0秒 ->时间戳
func int2Time(tm int32) time.Time {
	day := tm % 100
	month := tm / 100 % 100
	year := tm / 10000
	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)
	return date
}
