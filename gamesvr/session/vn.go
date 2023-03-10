package session

import (
	"context"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// VnGetInfo 视觉小说-获取信息
func (s *Session) VnGetInfo(ctx context.Context, req *pb.C2SVnGetInfo) (*pb.S2CVnGetInfo, error) {
	return &pb.S2CVnGetInfo{
		VnIdList: s.User.Info.RewardedVnIds.Values(),
	}, nil
}

// VnGetReward 视觉小说-领取剧情奖励
func (s *Session) VnGetReward(ctx context.Context, req *pb.C2SVnGetReward) (*pb.S2CVnGetReward, error) {
	err := s.User.VnGetReward(req.VnId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CVnGetReward{
		VnId:           req.VnId,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// VnRead 视觉小说-阅读
func (s *Session) VnRead(ctx context.Context, req *pb.C2SVnRead) (*pb.S2CVnRead, error) {
	s.User.VnRead(ctx, req.VnId)
	return &pb.S2CVnRead{
		VnId: req.VnId,
	}, nil
}
