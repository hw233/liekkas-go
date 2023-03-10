package session

import (
	"context"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// 使用道具
func (s *Session) ItemUse(ctx context.Context, req *pb.C2SItemUse) (*pb.S2CItemUse, error) {
	if req.GetAmount() <= 0 {
		return nil, errors.WrapTrace(common.ErrParamError)
	}
	reward := common.NewReward(req.GetItemId(), req.GetAmount())
	err := s.User.ItemUse(reward, req.GetParam())
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CItemUse{ResourceResult: s.User.VOResourceResult()}, nil
}
