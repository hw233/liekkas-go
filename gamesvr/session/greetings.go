package session

import (
	"context"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

// 获取公会异界问候
func (s *Session) GuildGetMemberGreetings(ctx context.Context, req *pb.C2SGuildGetMemberGreetings) (*pb.S2CGuildGetMemberGreetings, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	if s.Guild.GuildID == 0 {
		return nil, errors.Swrapf(common.ErrGuildNotFound)
	}

	// 从公会获取问候
	resp, err := s.RPCGuildGetGreetings(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	voGreetings, err := s.User.GetGreetings(ctx, resp.Greetings)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildGetMemberGreetings{
		Greetings: voGreetings,
	}, nil
}

func (s *Session) GuildGreetingsRewards(ctx context.Context, req *pb.C2SGuildGreetingsRewards) (*pb.S2CGuildGreetingsRewards, error) {

	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	if s.Guild.GuildID == 0 {
		return nil, errors.Swrapf(common.ErrGuildNotFound)
	}

	lastTimestamp, err := s.User.GreetingsRewards(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	// 为0,代表并没有领取任何奖励
	if lastTimestamp == 0 {
		return nil, nil
	}

	_, err = s.RPCGuildUpdateGreetings(ctx, lastTimestamp)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildGreetingsRewards{
		ResourceResult: s.VOResourceResult(),
	}, nil
}

func (s *Session) PushGreetings(ctx context.Context) error {

	if len(s.User.Guild.GuildGreetings.SendQueue) <= 0 {
		return nil
	}

	if s.Guild.GuildID == 0 {
		// 清空greetings
		s.User.ClearGreetings()
		return nil
	}

	voGreetings := make([]*pb.VOGreetings, 0, len(s.User.Guild.GuildGreetings.SendQueue))

	for _, greeting := range s.User.Guild.GuildGreetings.SendQueue {
		voGreetings = append(voGreetings, greeting.VOGreetings())
	}

	_, err := s.RPCGuildAddGreetings(ctx, voGreetings)
	if err != nil {
		return errors.WrapTrace(err)
	}

	s.User.ClearGreetings()

	return nil
}
