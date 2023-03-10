package session

import (
	"context"
	"shared/protobuf/pb"
)

func (s *Session) FormationInfo(ctx context.Context, req *pb.C2SFormationInfo) (*pb.S2CFormationInfo, error) {
	return &pb.S2CFormationInfo{Formations: s.User.FormationInfo.VOFormations()}, nil
}

func (s *Session) SetFormation(ctx context.Context, req *pb.C2SSetFormation) (*pb.S2CSetFormation, error) {
	err := s.User.SetFormation(req.Id, req.FormationType, req.BattleFormation)
	if err != nil {
		return nil, err
	}

	return &pb.S2CSetFormation{
		Formation: &pb.VOFormation{
			Id:              req.Id,
			FormationType:   req.FormationType,
			BattleFormation: req.BattleFormation,
		},
	}, nil
}
