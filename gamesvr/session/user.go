package session

import (
	"context"
	"gamesvr/manager"
	"gamesvr/model"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

func (s *Session) SuitUser(ctx context.Context, req *pb.C2SSuitUser) (*pb.S2CSuitUser, error) {
	// log.Printf("session: name: %s uid: %d \n", req.Name, req.Uid)

	return &pb.S2CSuitUser{
		UserInfo:    s.User.VOUserInfo(),
		Resource:    s.User.VOUserResource(),
		Items:       s.User.VOItemInfo(),
		Characters:  s.User.VOUserCharacter(),
		Equipments:  s.User.VOUserEquipment(),
		WorldItems:  s.User.VOUserWorldItem(), // []*VOUserWorldItem
		ManualInfos: s.User.VOManualInfo(),    // []*VOManualInfo
	}, nil
}

// HeartBeatReq 心跳
func (s *Session) HeartBeatReq(ctx context.Context, req *pb.C2SHeartBeatReq) (*pb.S2CHeartBeatReq, error) {
	return &pb.S2CHeartBeatReq{}, nil
}

func (s *Session) ExchangeCDKey(ctx context.Context, req *pb.C2SExchangeCDKey) (*pb.S2CExchangeCDKey, error) {
	if s.isGMCode(req.Code) {
		err := s.handelGMCode(ctx, req.Code)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}

	return &pb.S2CExchangeCDKey{
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

func (s *Session) RefreshRedPoint(ctx context.Context, req *pb.C2SRefreshRedPoint) (*pb.S2CRefreshRedPoint, error) {
	var retRedPoints []*pb.VORedPoint

	for _, redPoint := range req.NeedRefreshRedPoints {
		switch redPoint {
		case static.ReddotTypeElitereward:
			count := s.User.CheckChapterRewardNoticeByType(model.ChapterTypeElite)
			if count > 0 {
				retRedPoints = append(retRedPoints, &pb.VORedPoint{
					Id:    redPoint,
					Value: count,
				})
			}
		}
	}

	return &pb.S2CRefreshRedPoint{
		RedPoints: retRedPoints,
	}, nil
}

func (s *Session) GM(ctx context.Context, req *pb.C2SGM) (*pb.S2CGM, error) {
	ret, err := s.User.ExecuteGM(ctx, req.GmType, req.Args)
	if err != nil {
		return nil, err
	}
	ret.ResourceResult = s.User.VOResourceResult()

	return ret, nil
}

func (s *Session) FetchAnnouncement(ctx context.Context, req *pb.C2SFetchAnnouncement) (*pb.S2CFetchAnnouncement, error) {
	announcements, banners := manager.Announcements.GetAnnouncements()

	voAnncs := make([]*pb.VOAnnouncement, 0, len(announcements))
	voBanners := make([]*pb.VOBanner, 0, len(banners))

	for _, annc := range announcements {
		voAnncs = append(voAnncs, annc.VOAnnouncement())
	}

	for _, banner := range banners {
		voBanners = append(voBanners, banner.VOBanner())
	}

	return &pb.S2CFetchAnnouncement{
		Announcements: voAnncs,
		Banners:       voBanners,
	}, nil
}

//----------------------------------------
//Notify
//----------------------------------------
func (s *Session) DailyRefreshNotify() {
	notify := s.User.PopDailyRefreshNotify()
	if notify != nil {
		s.push(manager.CSV.Protocol.Pushes.DailyRefreshNotify, notify)
	}
}

//----------------------------------------
//Power
//----------------------------------------
func (s *Session) UpdateUserPower() {
	s.User.GetUserPower()
}
