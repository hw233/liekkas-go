package session

import (
	"context"
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/servertime"
)

func (s *Session) GetMercenaryList(ctx context.Context, req *pb.C2SGetMercenaryList) (*pb.S2CGetMercenaryList, error) {

	var elites []int64

	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	if s.Guild.GuildID != 0 {
		// 向公会请求所有精英成员的id和名字
		resp, err := s.RPCGuildMembers(ctx)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		elites = resp.GetUserIds()
	}

	voResults, err := s.User.GetMercenaryList(ctx, []int64{}, elites)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return voResults, nil
}

func (s *Session) MercenarySendApply(ctx context.Context, req *pb.C2SMercenarySendApply) (*pb.S2CMercenarySendApply, error) {

	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	var elites []int64

	if s.Guild.GuildID != 0 {
		resp, err := s.RPCGuildMembers(ctx)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		elites = resp.GetUserIds()
	}

	// 读取redis，查看是否已经被借出,之后发送申请
	status, err := s.User.SendApply(ctx, req.UserId, req.CharacterId, []int64{}, elites)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CMercenarySendApply{
		UserId:      req.UserId,
		CharacterId: req.CharacterId,
		Status:      status,
	}, nil
}
func (s *Session) MercenaryCancel(ctx context.Context, req *pb.C2SMercenaryCancel) (*pb.S2CMercenaryCancel, error) {
	timestamp, status, err := s.User.CancelApply(ctx, req.UserId, req.CharacterId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CMercenaryCancel{
		UserId:       req.UserId,
		CharacterId:  req.CharacterId,
		CoolDownTime: timestamp - servertime.Now().Unix(),
		Status:       status,
	}, nil
}

func (s *Session) MercenaryManagement(ctx context.Context, req *pb.C2SMercenaryManagement) (*pb.S2CMercenaryManagement, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	var members []int64

	if s.Guild.GuildID != 0 {
		resp, err := s.RPCGuildMembers(ctx)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		members = resp.GetUserIds()
	}

	voApply, err := s.User.GetMercenaryManagement(ctx, []int64{}, members)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CMercenaryManagement{
		Apply: voApply,
	}, nil
}

func (s *Session) MercenaryHandleApply(ctx context.Context, req *pb.C2SMercenaryHandleApply) (*pb.S2CMercenaryHandleApply, error) {
	// var realation int32 // 1 好友， 2 公会，3 都是

	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if s.Guild.GuildID == 0 {
		// todo 增加好友功能之后需要修改
		return nil, errors.Swrapf(common.ErrGuildNotFound)
	}
	members, err := s.RPCGuildMembers(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	status, err := s.User.HandleApply(ctx, req.IfIgnore, req.UserId, req.CharacterId, []int64{}, members.GetUserIds())
	if err != nil {
		return nil, nil
	}

	return &pb.S2CMercenaryHandleApply{
		CharacterId: req.CharacterId,
		UserId:      req.UserId,
		Status:      status, // 1 已经非好友或者公会 2此人已经从其他地方借到了, 0则借成功了
	}, nil
}

func (s *Session) MercenaryRecord(ctx context.Context, req *pb.C2SMercenaryRecord) (*pb.S2CMercenaryRecord, error) {
	voRecords := s.User.GetMercenaryRecord(ctx)

	return &pb.S2CMercenaryRecord{
		Records: voRecords,
	}, nil
}

// 每半分钟检查一次
func (s *Session) MercenaryCharacterList(ctx context.Context) {

	s.User.UpdateMercenaryBorrow(ctx)
}

// 请求佣兵数据
func (s *Session) GetMercenaryCharacter(ctx context.Context, req *pb.C2SGetMercenaryCharacter) (*pb.S2CGetMercenaryCharacter, error) {
	if req.SystemType != static.BattleTypeTower && req.SystemType != static.BattleTypeChallengeAltar {
		return nil, errors.Swrapf(common.ErrWrongSystemType, req.SystemType)
	}

	voDetail, err := s.User.GetMercenaryCharacter(ctx, req.SystemType)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.S2CGetMercenaryCharacter{
		Mercenaries: voDetail,
	}, nil
}

// 更新玩家的角色数据
func (s *Session) UserCharacterSimpleUpdate(ctx context.Context) error {
	var characters []*common.MercenaryCharacter
	for _, character := range *s.User.CharacterPack {
		if character.Rarity >= static.RaritySr {
			c := common.NewMercenaryCharacter(character.ID, character.Level, character.Star, character.Power)
			characters = append(characters, c)
		}
	}
	err := manager.Global.UserMercenary.SetMercenaryCharacterData(ctx, s.GetUserId(), characters)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}
