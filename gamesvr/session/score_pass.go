package session

import (
	"context"
	"gamesvr/manager"
	"shared/protobuf/pb"
)

func (s *Session) ScorePassInfo(ctx context.Context, req *pb.C2SScorePassInfo) (*pb.S2CScorePassInfo, error) {
	return &pb.S2CScorePassInfo{Info: s.User.ScorePassInfo.VOScorePassInfo()}, nil
}

func (s *Session) ReceiveScorePassReward(ctx context.Context, req *pb.C2SReceiveScorePassReward) (*pb.S2CReceiveScorePassReward, error) {
	err := s.User.ReceiveScorePassReward(req.RewardId)
	if err != nil {
		return nil, err
	}

	return &pb.S2CReceiveScorePassReward{
		RewardId:       req.RewardId,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

//----------------------------------------
//Notify
//----------------------------------------
func (s *Session) TryPushScorePassNotify() {
	notify := s.User.PopScorePassNotify()
	if notify == nil {
		return
	}

	s.push(manager.CSV.Protocol.Pushes.ScorePassNotify, notify)
}
