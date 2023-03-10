package session

import (
	"context"
	"gamesvr/manager"
	"gamesvr/model"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// UserKanban 设置看板娘
func (s *Session) UserKanban(ctx context.Context, req *pb.C2SUserKanban) (*pb.S2CUserKanban, error) {
	err := s.User.SetKanban(req.CharacterId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CUserKanban{
		CharacterId: req.CharacterId,
	}, nil
}

// UserAvatar 设置头像
func (s *Session) UserAvatar(ctx context.Context, req *pb.C2SUserAvatar) (*pb.S2CUserAvatar, error) {
	err := s.User.SetAvatar(req.Avatar)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CUserAvatar{
		Avatar: req.Avatar,
	}, nil
}

// UserNickname 设置昵称
func (s *Session) UserNickname(ctx context.Context, req *pb.C2SUserNickname) (*pb.S2CUserNickname, error) {
	err := s.User.ChangeNickname(ctx, req.Nickname)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CUserNickname{
		Nickname:       s.User.Name,
		NameIndex:      s.User.Info.NameIndex,
		ResourceResult: s.VOResourceResult(),
	}, nil
}

// UserSignature 设置签名
func (s *Session) UserSignature(ctx context.Context, req *pb.C2SUserSignature) (*pb.S2CUserSignature, error) {
	err := s.User.ChangeSignature(req.Signature)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CUserSignature{
		Signature: req.Signature,
	}, nil
}

// VisitingCardSetCharacter 设置名片展示角色
func (s *Session) VisitingCardSetCharacter(ctx context.Context, req *pb.C2SVisitingCardSetCharacter) (*pb.S2CVisitingCardSetCharacter, error) {
	err := s.User.VisitingCardSetCharacter(req.PosList)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CVisitingCardSetCharacter{}, nil
}

// VisitingCardSetWorldItem 设置名片展示世界级道具
func (s *Session) VisitingCardSetWorldItem(ctx context.Context, req *pb.C2SVisitingCardSetWorldItem) (*pb.S2CVisitingCardSetWorldItem, error) {
	err := s.User.VisitingCardSetWorldItem(req.WorldItemUId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CVisitingCardSetWorldItem{}, nil
}

// VisitingCard 查看名片信息
func (s *Session) VisitingCard(ctx context.Context, req *pb.C2SVisitingCard) (*pb.S2CVisitingCard, error) {
	if req.UserId == s.User.ID {
		err := s.GuildDataRefresh(ctx)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		return &pb.S2CVisitingCard{
			Card: s.User.VOVisitingCard(),
		}, nil
	} else {
		card, err := OtherUserVisitingCard(ctx, req.UserId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		return &pb.S2CVisitingCard{
			Card: card,
		}, nil
	}

}

func OtherUserVisitingCard(ctx context.Context, id int64) (*pb.VOVisitingCard, error) {
	userCache, err := manager.Global.GetUserCache(ctx, id)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return model.NewOthersUserVisitingCard(userCache).VOVisitingCard(), nil

}

// DailyGiftInfo 整点奖励领取情况
func (s *Session) DailyGiftInfo(ctx context.Context, req *pb.C2SDailyGiftInfo) (*pb.S2CDailyGiftInfo, error) {
	return &pb.S2CDailyGiftInfo{
		Rewarded: s.User.Info.DailyGiftInfoRewarded,
	}, nil
}

// GetDailyGiftReward 领取整点奖励
func (s *Session) GetDailyGiftReward(ctx context.Context, req *pb.C2SGetDailyGiftReward) (*pb.S2CGetDailyGiftReward, error) {
	err := s.User.GetDailyGiftReward(int(req.Index))
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGetDailyGiftReward{
		Rewarded:       s.User.Info.DailyGiftInfoRewarded,
		ResourceResult: s.User.VOResourceResult(),
	}, nil
}

// UserFrame 设置头像框
func (s *Session) UserFrame(ctx context.Context, req *pb.C2SUserFrame) (*pb.S2CUserFrame, error) {
	err := s.User.SetFrame(req.Frame)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CUserFrame{
		Frame: req.Frame,
	}, nil
}

// FetchIsFirstRename 查询是否是第一次改名
func (s *Session) FetchIsFirstRename(ctx context.Context, req *pb.C2SFetchIsFirstRename) (*pb.S2CFetchIsFirstRename, error) {
	return &pb.S2CFetchIsFirstRename{
		IsFirstRename: s.Info.IsFirstRename,
	}, nil
}
