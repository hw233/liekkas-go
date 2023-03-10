package model

import (
	"context"

	"gamesvr/manager"
	"shared/utility/mysql"
)

// var manager *mysql.Manager

type SurveyAnswer struct {
	QID    int32  `json:"qid"`
	Answer string `json:"answer"`
}

type SurveyAnswers []SurveyAnswer

type Survey struct {
	UID                int64          `db:"user_id" where:"=" major:"true"`
	SurveyID           int32          `db:"survey_id"`
	Answers            *SurveyAnswers `db:"survey_answers"`
	*mysql.EmbedModule `db:"-"`
}

func NewSurvey() *Survey {
	return &Survey{
		UID:         0,
		SurveyID:    0,
		Answers:     NewSurveyAnswers(),
		EmbedModule: &mysql.EmbedModule{},
	}
}

func NewSurveyAnswers() *SurveyAnswers {
	return &SurveyAnswers{}
}

func (s *Survey) SetSurvey(uid int64, sid int32) {
	s.UID = uid
	s.SurveyID = sid
}

func (s *Survey) AddAnswer(qid int32, answer string) {
	cur := &SurveyAnswer{}
	cur.QID = qid
	cur.Answer = answer
	answers := *s.Answers
	answers = append(answers, *cur)
	s.Answers = &answers
}

func (s *Survey) NewSurveyData(ctx context.Context) error {
	err := manager.MySQL.Create(ctx, s)
	if err != nil {
		return err
	}
	return nil
}
