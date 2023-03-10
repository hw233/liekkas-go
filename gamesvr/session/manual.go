package session

import (
	"context"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/number"
)

func (s *Session) ManualGetReward(ctx context.Context, req *pb.C2SManualGetReward) (*pb.S2CManualGetReward, error) {
	if len(req.ManualIds) == 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	repeatableArr := number.NewNonRepeatableArr()
	repeatableArr.Append(req.ManualIds...)
	vos, err := s.User.ManualGetReward(repeatableArr)
	if err != nil {
		return nil, err
	}
	return &pb.S2CManualGetReward{
		Infos:          vos,
		ResourceResult: s.VOResourceResult(),
	}, nil
}
func (s *Session) ManualGetAll(ctx context.Context, req *pb.C2SManualGetAll) (*pb.S2CManualGetAll, error) {
	return &pb.S2CManualGetAll{
		Infos: s.User.VOManualInfo(),
	}, nil
}
