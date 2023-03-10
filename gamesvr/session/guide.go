package session

import (
	"context"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// UserGuideGetAll  新手引导-获取所有已通过引导
func (s *Session) UserGuideGetAll(ctx context.Context, req *pb.C2SUserGuideGetAll) (*pb.S2CUserGuideGetAll, error) {

	return &pb.S2CUserGuideGetAll{
		PassedGuideIds: s.User.Info.PassedGuideIds.Values(),
	}, nil
}

// UserGuide  新手引导-通过某一引导
func (s *Session) UserGuide(ctx context.Context, req *pb.C2SUserGuide) (*pb.S2CUserGuide, error) {
	err := s.User.UserGuide(req.GuideId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CUserGuide{
		GuideId:        req.GuideId,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// DestinyChild 获得新手池子本命角色(引导用)
func (s *Session) DestinyChild(ctx context.Context, req *pb.C2SDestinyChild) (*pb.S2CDestinyChild, error) {
	return &pb.S2CDestinyChild{
		DestinyChild: s.User.GachaRecords.DestinyChild,
	}, nil
}
