package session

import (
	"context"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// GetGachaList 获取卡池信息
func (s *Session) GetGachaList(ctx context.Context, req *pb.C2SGetGachaList) (*pb.S2CGetGachaList, error) {
	return s.User.GetGachaList(), nil
}

// UserGachaDrop 抽卡
func (s *Session) UserGachaDrop(ctx context.Context, req *pb.C2SUserGachaDrop) (*pb.S2CUserGachaDrop, error) {
	vos, err := s.User.UserGachaDrop(req.GachaId, req.IsSingle)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CUserGachaDrop{
		GachaInfo:      vos,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}
func (s *Session) UserGachaRecords(ctx context.Context, req *pb.C2SUserGachaRecords) (*pb.S2CUserGachaRecords, error) {
	if req.Num == 0 || req.Offset == 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}

	vos := s.User.UserGachaRecords(req.CharaOrWorldItem, req.Offset, req.Num)
	return &pb.S2CUserGachaRecords{
		Records: vos,
	}, nil
}
