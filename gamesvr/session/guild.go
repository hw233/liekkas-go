package session

import (
	"context"

	"gamesvr/manager"
	"gamesvr/model"
	"shared/common"
	"shared/protobuf/pb"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/glog"
	"shared/utility/servertime"
)

func (s *Session) GuildInfo(ctx context.Context, req *pb.C2SGuildInfo) (*pb.S2CGuildInfo, error) {
	// err := s.User.CheckActionUnlock(static.ActionIdTypeGuildunlock)
	// if err != nil {
	// 	return nil, err
	// }

	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildInfo(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildInfo{
		GuildInfo: ret.GuildInfo,
	}, nil
}

func (s *Session) GuildSendGroupMail(ctx context.Context, req *pb.C2SGuildSendGroupMail) (*pb.S2CGuildSendGroupMail, error) {
	// err := s.User.CheckActionUnlock(static.ActionIdTypeGuildunlock)
	// if err != nil {
	// 	return nil, err
	// }

	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否加入公会
	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	_, err = s.RPCGuildSendGroupMail(ctx, req.Title, req.Content)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	mail := model.NewMail(0, 9, req.Title, []string{}, req.Content, []string{},
		nil, "", servertime.Now().Unix(), 0)

	s.User.BIMail(mail, s.User.GetUserId(), bilog.MailOpSend)

	return &pb.S2CGuildSendGroupMail{}, nil
}

func (s *Session) GuildCreate(ctx context.Context, req *pb.C2SGuildCreate) (*pb.S2CGuildCreate, error) {
	// err := s.User.CheckActionUnlock(static.ActionIdTypeGuildunlock)
	// if err != nil {
	// 	return nil, err
	// }

	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	err = s.Guild.CheckHasJoinedGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查待消耗的物品是否足够
	rewardId := manager.CSV.GlobalEntry.GuildCreateCost[0]
	rewardNum := manager.CSV.GlobalEntry.GuildCreateCost[1]
	r := common.NewReward(rewardId, rewardNum)
	costs := common.NewRewards()
	costs.AddReward(r)
	err = s.User.CheckRewardsEnough(costs)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildCreate(ctx, req.Name, req.Icon, req.JoinModel, req.Title)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	s.Guild.Join(ret.GuildInfo.GuildID, ret.GuildInfo.Name)

	// 消耗
	reason := logreason.NewReason(logreason.GuildCreate)
	err = s.User.CostRewards(costs, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildCreate{
		GuildInfo: ret.GuildInfo,
		Resource:  s.VOResourceResult(),
	}, nil
}

func (s *Session) GuildDissolve(ctx context.Context, req *pb.C2SGuildDissolve) (*pb.S2CGuildDissolve, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildDissolve(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildDissolve{
		GuildInfo: ret.GuildInfo,
	}, nil
}

func (s *Session) GuildCancelDissolve(ctx context.Context, req *pb.C2SGuildCancelDissolve) (*pb.S2CGuildCancelDissolve, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildCancelDissolve(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildCancelDissolve{
		GuildInfo: ret.GuildInfo,
	}, nil
}

func (s *Session) GuildApply(ctx context.Context, req *pb.C2SGuildApply) (*pb.S2CGuildApply, error) {
	// err := s.User.CheckActionUnlock(static.ActionIdTypeGuildunlock)
	// if err != nil {
	// 	return nil, err
	// }

	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查没加工会
	err = s.Guild.CheckHasJoinedGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查退会CD
	err = s.Guild.CheckQuitGuildCD()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = s.Guild.CheckApplyNumLimit()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildApply(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	resp, err := s.RPCGuildRecommendInfo(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	if ret.IsJoined {
		s.Guild.Join(ret.GuildInfo.GuildID, ret.GuildInfo.Name)
	} else {
		s.Guild.Apply(resp.Info)
	}

	return &pb.S2CGuildApply{
		IsJoined:  ret.IsJoined,
		UserGuild: s.Guild.VOGuildInfo(),
		GuildInfo: ret.GuildInfo,
	}, nil
}

func (s *Session) GuildCancelApply(ctx context.Context, req *pb.C2SGuildApply) (*pb.S2CGuildApply, error) {
	// err := s.User.CheckActionUnlock(static.ActionIdTypeGuildunlock)
	// if err != nil {
	// 	return nil, err
	// }
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查没加工会
	err = s.Guild.CheckHasJoinedGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	_, err = s.RPCGuildCancelApply(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	s.Guild.CancelApply(req.GuildID)

	return &pb.S2CGuildApply{
		UserGuild: s.Guild.VOGuildInfo(),
	}, nil
}

func (s *Session) GuildHandleApplied(ctx context.Context, req *pb.C2SGuildHandleApplied) (*pb.S2CGuildHandleApplied, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否加入公会
	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildHandleApplied(ctx, req.Approve, req.Refuse)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildHandleApplied{
		GuildInfo: ret.GuildInfo,
	}, nil
}

func (s *Session) GuildKick(ctx context.Context, req *pb.C2SGuildKick) (*pb.S2CGuildKick, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否加入公会
	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildKick(ctx, req.UserID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildKick{
		GuildInfo: ret.GuildInfo,
	}, nil
}

func (s *Session) GuildQuit(ctx context.Context, req *pb.C2SGuildQuit) (*pb.S2CGuildQuit, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否加入公会
	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	_, err = s.RPCGuildQuit(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 只有主动退出才有CD
	s.Guild.LastQuitTime = servertime.Now().Unix()
	s.Guild.Quit()

	return &pb.S2CGuildQuit{
		UserGuild: s.Guild.VOGuildInfo(),
	}, nil
}

func (s *Session) GuildChat(ctx context.Context, req *pb.C2SGuildChat) (*pb.S2CGuildChat, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否加入公会
	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildChat(ctx, req.Content)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildChat{
		GuildInfo: ret.GuildInfo,
	}, nil
}

func (s *Session) GuildPromotion(ctx context.Context, req *pb.C2SGuildPromotion) (*pb.S2CGuildPromotion, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否加入公会
	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	resp, err := s.RPCGuildPromotion(ctx, req.UserID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildPromotion{
		GuildInfo: resp.GuildInfo,
	}, nil
}

func (s *Session) GuildDemotion(ctx context.Context, req *pb.C2SGuildDemotion) (*pb.S2CGuildDemotion, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否加入公会
	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	resp, err := s.RPCGuildDemotion(ctx, req.UserID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildDemotion{
		GuildInfo: resp.GuildInfo,
	}, nil
}

func (s *Session) GuildTransfer(ctx context.Context, req *pb.C2SGuildTransfer) (*pb.S2CGuildTransfer, error) {
	req.GetUserID()
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// 检查是否加入公会
	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	resp, err := s.RPCGuildTransfer(ctx, req.UserID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildTransfer{
		GuildInfo: resp.GuildInfo,
	}, nil
}

func (s *Session) GuildSearch(ctx context.Context, req *pb.C2SGuildSearch) (*pb.S2CGuildSearch, error) {
	// err := s.User.CheckActionUnlock(static.ActionIdTypeGuildunlock)
	// if err != nil {
	// 	return nil, err
	// }

	guildID, err := manager.Global.GetGuildID(ctx, req.GuildName)
	if err == global.ErrNil {
		// 没查到结果
		return nil, common.ErrGuildNotFound
	} else if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildShowInfo(ctx, guildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildSearch{
		GuildShowInfo: ret.GuildShowInfo,
	}, nil
}

func (s *Session) GuildDataRefresh(ctx context.Context) error {
	err := s.Guild.RefreshApplied(ctx, s.ID)
	if err != nil {
		return errors.WrapTrace(err)
	}

	err = s.GuildRefreshDissolved(ctx)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (s *Session) GuildSync(ctx context.Context, status int32) error {
	err := s.Guild.RefreshApplied(ctx, s.ID)
	if err != nil {
		return errors.WrapTrace(err)
	}

	if s.Guild.HasJoinedGuild() {
		_, err = s.RPCGuildSync(ctx, status)
		if err != nil {
			glog.Errorf("RPCGuildSync(%d) error: %v", s.Guild.GuildID, err)
		}
	}

	return nil
}

func (s *Session) GuildRefreshDissolved(ctx context.Context) error {
	if s.Guild.NeedCheckDissolved() {
		ret, err := s.RPCGuildIsDissolved(ctx)
		if err != nil {
			return errors.WrapTrace(err)
		}

		if ret.IsDissolved {
			s.Guild.Quit()
			return common.ErrGuildDissolved
		}
		if !ret.IsInDissolving {
			s.Guild.LastCheckDissolve = servertime.Now().Unix()
		}
	}

	return nil
}

func (s *Session) GuildGetList(ctx context.Context, req *pb.C2SGuildGetList) (*pb.S2CGuildGetList, error) {

	result := make([]*pb.VOGuildSimpleInfo, 0, 50)

	// redis 里获取50个id
	ids, err := manager.Global.GuildSetRandomGet(ctx, 50)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	for _, guildId := range ids {
		if s.User.Guild.CheckIsAlreadyApply(guildId) {
			continue
		}
		resp, err := s.RPCGuildRecommendInfo(ctx, guildId)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		result = append(result, resp.Info)
	}

	return &pb.S2CGuildGetList{
		Info:         result,
		LastQuitTime: s.User.Guild.LastQuitTime,
		ApplyNum:     int32(len(s.User.Guild.ApplyList)),
	}, nil
}

func (s *Session) GuildGetApplyList(ctx context.Context, req *pb.C2SGuildGetApplyList) (*pb.S2CGuildGetApplyList, error) {
	resp, err := s.RPCGuildGetApplyList(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildGetApplyList{
		ApplyList: resp.ApplyList,
	}, nil
}

func (s *Session) GuildListInfo(ctx context.Context, req *pb.C2SGuildListInfo) (*pb.S2CGuildListInfo, error) {

	resp, err := s.RPCGuildShowInfo(ctx, req.GuildId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	// for _, member := range resp.GuildShowInfo.ChairmenInfo {
	// 	fmt.Println("-----------------------", member.UserID, member.UserName)
	// }

	return &pb.S2CGuildListInfo{
		GuildInfo: resp.GuildShowInfo,
	}, nil
}

func (s *Session) GuildApplyList(ctx context.Context, req *pb.C2SGuildApplyList) (*pb.S2CGuildApplyList, error) {
	voInfo := make([]*pb.VOGuildSimpleInfo, 0, len(s.User.Guild.ApplyList))

	for _, apply := range s.User.Guild.ApplyList {
		voInfo = append(voInfo, apply.VOGuildSimpleInfo())
	}

	return &pb.S2CGuildApplyList{
		Info: voInfo,
	}, nil
}

func (s *Session) GuildTaskRewards(ctx context.Context, req *pb.C2SGuildTaskRewards) (*pb.S2CGuildTaskRewards, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	ret, err := s.RPCGuildTaskRewards(ctx, req.Id)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	taskConfig, err := manager.CSV.Guild.GetTaskConfig(req.Id)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	reason := logreason.NewReason(logreason.GuildTask)
	_, err = s.User.AddRewardsByDropId(taskConfig.DropId, reason)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildTaskRewards{
		ResourceResult:   s.VOResourceResult(),
		ActivationChange: ret.Activation,
	}, nil
}

func (s *Session) GuildTaskList(ctx context.Context, req *pb.C2SGuildTaskList) (*pb.S2CGuildTaskList, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	ret, err := s.RPCGuildTaskList(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildTaskList{
		Tasks: ret.Tasks,
	}, nil
}

func (s *Session) PushGuildTasks(ctx context.Context) error {

	if len(s.User.Guild.Tasks) <= 0 {
		return nil
	}

	if s.Guild.GuildID == 0 {
		// 清空greetings
		s.User.Guild.ClearTaskItems()
		return nil
	}

	voTasks := make([]*pb.VOGuildTaskItem, 0, len(s.User.Guild.Tasks))

	for _, task := range s.User.Guild.Tasks {
		voTasks = append(voTasks, task.VOGuildTaskItem())
	}

	_, err := s.RPCGuildTaskAddProgress(ctx, voTasks)
	if err != nil {
		return errors.WrapTrace(err)
	}

	s.User.Guild.ClearTaskItems()

	return nil
}

// ----------------公会编辑---------------------
func (s *Session) GuildModify(ctx context.Context, req *pb.C2SGuildModify) (*pb.S2CGuildModify, error) {
	err := s.GuildDataRefresh(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	err = s.Guild.CheckHasNotJoinGuild()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if req.Icon == nil {
		return nil, errors.Swrapf(common.ErrGuildModifyIconIsNull, s.Guild.GuildID)
	}
	ret, err := s.RPCGuildModify(ctx, req.Icon, req.JoinModel, req.Title)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return &pb.S2CGuildModify{
		GuildId:   s.Guild.GuildID,
		Icon:      ret.GetInfo().Icon,
		JoinModel: ret.Info.JoinModel,
		Title:     ret.Info.Title,
	}, nil
}
