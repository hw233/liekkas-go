package session

import (
	"context"

	"gamesvr/manager"
	"shared/protobuf/pb"
	"shared/utility/errors"
)

func (s *Session) RPCGuildIsDissolved(ctx context.Context) (*pb.GuildIsDissolvedResp, error) {
	resp, err := manager.RPCGuildClient.IsDissolved(ctx, &pb.GuildIsDissolvedReq{
		GuildID: s.Guild.GuildID,
	})

	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildRecommend(ctx context.Context) (*pb.GuildRecommendResp, error) {
	resp, err := manager.RPCGuildClient.Recommend(ctx, &pb.GuildRecommendReq{})
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildSync(ctx context.Context, status int32) (*pb.GuildSyncResp, error) {
	resp, err := manager.RPCGuildClient.Sync(ctx, &pb.GuildSyncReq{
		GuildID:       s.Guild.GuildID,
		UserID:        s.ID,
		LastLoginTime: s.Info.LastLoginTime,
		Status:        status,
	})

	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildShowInfo(ctx context.Context, id int64) (*pb.GuildShowInfoResp, error) {
	resp, err := manager.RPCGuildClient.ShowInfo(ctx, &pb.GuildShowInfoReq{
		GuildID: id,
	})
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildInfo(ctx context.Context) (*pb.GuildInfoResp, error) {
	req := &pb.GuildInfoReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
	}

	resp, err := manager.RPCGuildClient.Info(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildCreate(ctx context.Context, name string, icon *pb.VOGuildIcon, joinModel int32, title string) (*pb.GuildCreateResp, error) {
	req := &pb.GuildCreateReq{
		GuildID:   0,
		UserID:    s.ID,
		Name:      name,
		Icon:      icon,
		Title:     title,
		JoinModel: joinModel,
	}

	resp, err := manager.RPCGuildClient.Create(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildDissolve(ctx context.Context) (*pb.GuildDissolveResp, error) {
	req := &pb.GuildDissolveReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
	}

	resp, err := manager.RPCGuildClient.Dissolve(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildCancelDissolve(ctx context.Context) (*pb.GuildCancelDissolveResp, error) {
	req := &pb.GuildCancelDissolveReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
	}

	resp, err := manager.RPCGuildClient.CancelDissolve(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildApply(ctx context.Context, guildID int64) (*pb.GuildApplyResp, error) {
	// err := manager.RPCGuildClient.SetSinglecastID(ctx, guildID)
	// if err != nil {
	// 	return nil, errors.WrapTrace(err)
	// }

	resp, err := manager.RPCGuildClient.Apply(ctx, &pb.GuildApplyReq{
		GuildID: guildID,
		UserID:  s.ID,
	})

	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildCancelApply(ctx context.Context, guildID int64) (*pb.GuildCancelApplyResp, error) {
	// err := manager.RPCGuildClient.SetSinglecastID(ctx, guildID)
	// if err != nil {
	// 	return nil, errors.WrapTrace(err)
	// }

	resp, err := manager.RPCGuildClient.CancelApply(ctx, &pb.GuildCancelApplyReq{
		GuildID: guildID,
		UserID:  s.ID,
	})

	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildQuit(ctx context.Context) (*pb.GuildQuitResp, error) {
	req := &pb.GuildQuitReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
	}

	resp, err := manager.RPCGuildClient.Quit(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildHandleApplied(ctx context.Context, approve, refuse []int64) (*pb.GuildHandleAppliedResp, error) {
	req := &pb.GuildHandleAppliedReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
		Approve: approve,
		Refuse:  refuse,
	}

	resp, err := manager.RPCGuildClient.HandleApplied(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildKick(ctx context.Context, userID int64) (*pb.GuildKickResp, error) {
	req := &pb.GuildKickReq{
		GuildID:      s.Guild.GuildID,
		UserID:       s.ID,
		KickedUserID: userID,
	}

	resp, err := manager.RPCGuildClient.Kick(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildSendGroupMail(ctx context.Context, title, content string) (*pb.GuildSendGroupMailResp, error) {
	req := &pb.GuildSendGroupMailReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
		Title:   title,
		Content: content,
		Sender:  s.Name,
	}

	resp, err := manager.RPCGuildClient.SendGroupMail(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildChat(ctx context.Context, content string) (*pb.GuildChatResp, error) {
	req := &pb.GuildChatReq{
		GuildID:  s.Guild.GuildID,
		UserID:   s.ID,
		UserName: s.Name,
		Content:  content,
		Avatar:   s.Info.Avatar,
		Frame:    s.Info.Frame,
	}

	resp, err := manager.RPCGuildClient.Chat(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildPromotion(ctx context.Context, target int64) (*pb.GuildPromotionResp, error) {
	req := &pb.GuildPromotionReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
		Target:  target,
	}

	resp, err := manager.RPCGuildClient.Promotion(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildGetApplyList(ctx context.Context) (*pb.GuildGetApplyListResp, error) {
	req := &pb.GuildGetApplyListReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
	}

	resp, err := manager.RPCGuildClient.GetApplyList(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildDemotion(ctx context.Context, target int64) (*pb.GuildDemotionResp, error) {
	req := &pb.GuildDemotionReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
		Target:  target,
	}

	resp, err := manager.RPCGuildClient.Demotion(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildTransfer(ctx context.Context, target int64) (*pb.GuildTransferResp, error) {
	req := &pb.GuildTransferReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
		Target:  target,
	}

	resp, err := manager.RPCGuildClient.Transfer(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildMembers(ctx context.Context) (*pb.GuildMembersResp, error) {
	req := &pb.GuildMembersReq{
		GuildID: s.Guild.GuildID,
	}

	resp, err := manager.RPCGuildClient.Members(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

// ----------------联合建造------------------
func (s *Session) RPCGuildCoBuildGetInfo(ctx context.Context, buildID int32) (*pb.GuildCoBuildGetInfoResp, error) {
	req := &pb.GuildCoBuildGetInfoReq{
		GuildID: s.Guild.GuildID,
		BuildID: buildID,
	}

	resp, err := manager.RPCGuildClient.CoBuildGetInfo(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildCoBuildImprove(ctx context.Context, buildID int32) (*pb.GuildCoBuildImproveResp, error) {
	req := &pb.GuildCoBuildImproveReq{
		GuildID: s.Guild.GuildID,
		BuildID: buildID,
		UserID:  s.ID,
	}

	resp, err := manager.RPCGuildClient.CoBuildImprove(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildCoBuildUse(ctx context.Context, buildID int32) (*pb.GuildCoBuildUseResp, error) {

	req := &pb.GuildCoBuildUseReq{
		GuildID: s.Guild.GuildID,
		BuildID: buildID,
		UserID:  s.ID,
	}

	resp, err := manager.RPCGuildClient.CoBuildUse(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

// -------------------派遣---------------------

// todo 调用之前需要判断玩家是否有公会
// func (s *Session) RPCGuildGetDispatchCharac(ctx context.Context, characId int32) (*pb.GuildGetDispatchCharacResp, error) {
// 	req := &pb.GuildGetDispatchCharacReq{
// 		GuildID:     s.Guild.GuildID,
// 		CharacterId: characId,
// 	}

// 	// fmt.Println("=========guildId, ", s.Guild.GuildID)
// 	resp, err := manager.RPCGuildClient.GetDispatchCharac(balancer.WithContext(ctx, &balancer.Context{
// 		ID:     s.Guild.GuildID,
// 		Server: "",
// 	}), req)
// 	if err != nil {
// 		return nil, errors.WrapTrace(err)
// 	}
// 	return resp, nil
// }

// -----------------异界问候----------------------
func (s *Session) RPCGuildAddGreetings(ctx context.Context, greetings []*pb.VOGreetings) (*pb.GuildAddGreetingsResp, error) {
	resp, err := manager.RPCGuildClient.AddGreetings(ctx, &pb.GuildAddGreetingsReq{
		GuildID:   s.Guild.GuildID,
		UserId:    s.GetUserId(),
		Greetings: greetings,
	})
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildGetGreetings(ctx context.Context) (*pb.GuildGetGreetingsResp, error) {
	resp, err := manager.RPCGuildClient.GetGreetings(ctx, &pb.GuildGetGreetingsReq{
		GuildID: s.Guild.GuildID,
		UserId:  s.GetUserId(),
	})
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildUpdateGreetings(ctx context.Context, lastTimestamp int64) (*pb.GuildUpdateGreetingsResp, error) {
	resp, err := manager.RPCGuildClient.UpdateGreetingsRecord(ctx, &pb.GuildUpdateGreetingsReq{
		GuildID:       s.Guild.GuildID,
		UserId:        s.GetUserId(),
		LastTimestamp: lastTimestamp,
	})

	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

// -----------------公会互助----------------------

func (s *Session) RPCGuildAddGraveyardRequest(ctx context.Context, add *pb.VOAddGraveyardRequest) (*pb.GuildAddGraveyardRequestResp, error) {
	req := &pb.GuildAddGraveyardRequestReq{
		GuildID: s.Guild.GuildID,
		Add:     add,
	}

	resp, err := manager.RPCGuildClient.AddGraveyardRequest(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildGetGraveyardRequests(ctx context.Context) (*pb.GuildGetGraveyardRequestsResp, error) {
	req := &pb.GuildGetGraveyardRequestsReq{
		GuildID: s.Guild.GuildID,
		UserId:  s.ID,
	}

	resp, err := manager.RPCGuildClient.GetGraveyardRequests(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildHelpRequestsHandle(ctx context.Context) (*pb.GuildHelpRequestsHandleResp, error) {
	req := &pb.GuildHelpRequestsHandleReq{
		GuildID:            s.Guild.GuildID,
		UserId:             s.ID,
		DailyAddActivation: s.Graveyard.DailyAddActivationByHelp,
		DailyAddGold:       s.Graveyard.DailyGuildGoldByHelp,
	}

	resp, err := manager.RPCGuildClient.HelpRequestsHandle(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildGetElite(ctx context.Context) (*pb.GuildGetElitesResp, error) {
	req := &pb.GuildGetElitesReq{
		GuildId: s.Guild.GuildID,
	}

	resp, err := manager.RPCGuildClient.GetElites(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildRecommendInfo(ctx context.Context, guildId int64) (*pb.GuildRecommendInfoResp, error) {

	req := &pb.GuildRecommendInfoReq{
		GuildId: guildId,
	}

	resp, err := manager.RPCGuildClient.RecommendInfo(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return resp, nil
}

// --------------公会任务-------------------
func (s *Session) RPCGuildTaskList(ctx context.Context) (*pb.GuildTaskListResp, error) {
	req := &pb.GuildTaskListReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
	}

	resp, err := manager.RPCGuildClient.GetTaskList(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildTaskRewards(ctx context.Context, taskId int32) (*pb.GuildTaskRewardsResp, error) {
	req := &pb.GuildTaskRewardsReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
		TaskId:  taskId,
	}

	resp, err := manager.RPCGuildClient.TaskRewards(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildTaskAddProgress(ctx context.Context, tasks []*pb.VOGuildTaskItem) (*pb.GuildTaskAddProgressResp, error) {

	req := &pb.GuildTaskAddProgressReq{
		GuildID: s.Guild.GuildID,
		UserID:  s.ID,
		Items:   tasks,
	}

	resp, err := manager.RPCGuildClient.TaskAddProgress(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}

func (s *Session) RPCGuildModify(ctx context.Context, icon *pb.VOGuildIcon, joinModel int32, title string) (*pb.GuildModifyResp, error) {
	req := &pb.GuildModifyReq{
		GuildId:   s.Guild.GuildID,
		UserId:    s.ID,
		Icon:      icon,
		JoinModel: joinModel,
		Title:     title,
	}

	resp, err := manager.RPCGuildClient.Modify(ctx, req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return resp, nil
}
