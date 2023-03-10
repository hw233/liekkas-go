package session

import (
	"context"
	"shared/protobuf/pb"
)

func (s *Session) SurveyData(ctx context.Context, req *pb.C2SSurveyData) (*pb.S2CSurveyData, error) {
	id := req.SurveyId
	answers := req.SurveyAnswers
	err := s.User.CheckForSurveyType(id)
	if err != nil {
		return nil, err
	}

	err = s.User.SaveSurvey(ctx, id, answers)
	if err != nil {
		return nil, err
	}

	return &pb.S2CSurveyData{
		SurveyId:       id,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) SurveyCheckMark(ctx context.Context, req *pb.C2SSurveyCheckMark) (*pb.S2CSurveyCheckMark, error) {

	result := s.User.Info.SurveyRecord.IsMarked(int(req.SurveyId))

	return &pb.S2CSurveyCheckMark{
		SurveyId: req.SurveyId,
		HasDone:  result,
	}, nil
}
