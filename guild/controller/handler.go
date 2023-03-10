package controller

import (
	"context"

	"guild/manager"
	"guild/session"
	"shared/common"
	"shared/utility/glog"

	"shared/protobuf/pb"
	"shared/utility/errors"
)

type GuildHandler struct {
	*pb.UnimplementedGuildServer
}

func NewGuildHandler() *GuildHandler {
	return &GuildHandler{
		UnimplementedGuildServer: &pb.UnimplementedGuildServer{},
	}
}

func (h *GuildHandler) GetSession(ctx context.Context, id int64) (*session.Session, error) {
	sess, err := manager.SessManager.GetSession(ctx, id)
	if err != nil {
		return nil, err
	}

	return sess.(*session.Session), nil
}

func (h *GuildHandler) ShowInfo(ctx context.Context, req *pb.GuildShowInfoReq) (*pb.GuildShowInfoResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	voInfo, err := sess.VOGuildShowInfo(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GuildShowInfoResp{
		GuildShowInfo: voInfo,
	}, nil
}

func (h *GuildHandler) Recommend(ctx context.Context, req *pb.GuildRecommendReq) (*pb.GuildRecommendResp, error) {
	sess, err := h.GetSession(ctx, 1)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	return &pb.GuildRecommendResp{
		// GuildShowInfoList: sess.VOGuildShowInfo(),
	}, nil
}

func (h *GuildHandler) Info(ctx context.Context, req *pb.GuildInfoReq) (*pb.GuildInfoResp, error) {
	glog.Info("Info: req[%+v]", req)
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		glog.Errorf("Info.GetSession(%d) error: %v", req.GuildID, err)
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	// 成员才可以获取详细的公会信息
	err = sess.CheckIsMember(req.UserID)
	if err != nil {
		glog.Errorf("Info.CheckIsMember(%d) error: %v", req.UserID, err)
		return nil, err
	}

	// 刷新在线状态
	sess.RefreshOnlineStatus()

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		glog.Errorf("Info.VOGuildInfo(%d) error: %v", req.UserID, err)
		return nil, err
	}

	return &pb.GuildInfoResp{
		GuildInfo: voGuildInfo,
	}, nil
}

func (h *GuildHandler) Create(ctx context.Context, req *pb.GuildCreateReq) (*pb.GuildCreateResp, error) {
	// 检查重名
	exist, err := manager.Global.GuildNameExist(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	// 重名了
	if exist {
		glog.Errorf("exist guild name: %s", req.Name)
		return nil, common.ErrGuildRepeatName
	}

	// 生成公会ID
	guildID, err := manager.Global.GenGuildID(ctx)
	if err != nil {
		return nil, err
	}

	success, err := manager.Global.AddGuildNameIfNotExist(ctx, guildID, req.Name)
	if err != nil {
		return nil, err
	}

	// 重名了
	if !success {
		glog.Errorf("exist guild name: %s", req.Name)
		return nil, common.ErrGuildRepeatName
	}

	managedSess, err := manager.SessManager.NewSessionIfNotExist(ctx, guildID)
	if err != nil {
		return nil, err
	}

	sess := managedSess.(*session.Session)

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildCreate(ctx, req.UserID, req.Name, req.Icon, req.JoinModel, req.Title)
	if err != nil {
		return nil, err
	}

	// 加入set，方便后续的随机推荐
	err = manager.Global.GuildSetAdd(ctx, guildID)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, 0)
	if err != nil {
		return nil, err
	}

	return &pb.GuildCreateResp{
		GuildInfo: voGuildInfo,
	}, nil
}

func (h *GuildHandler) Dissolve(ctx context.Context, req *pb.GuildDissolveReq) (*pb.GuildDissolveResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildDissolve(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// 公会解散，则从set中删除
	err = manager.Global.GuildSetDelete(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, 0)
	if err != nil {
		return nil, err
	}

	return &pb.GuildDissolveResp{
		GuildInfo: voGuildInfo,
	}, nil
}

func (h *GuildHandler) GetApplyList(ctx context.Context, req *pb.GuildGetApplyListReq) (*pb.GuildGetApplyListResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	ret, err := sess.GuildGetApplyList(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildGetApplyListResp{
		ApplyList: ret,
	}, nil
}

func (h *GuildHandler) CancelDissolve(ctx context.Context, req *pb.GuildCancelDissolveReq) (*pb.GuildCancelDissolveResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildCancelDissolve(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// 加入set，方便后续的随机推荐
	err = manager.Global.GuildSetAdd(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildCancelDissolveResp{
		GuildInfo: voGuildInfo,
	}, nil
}

func (h *GuildHandler) Apply(ctx context.Context, req *pb.GuildApplyReq) (*pb.GuildApplyResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	isJoined, err := sess.GuildApply(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	var guildInfo *pb.VOGuildInfo
	if isJoined {
		guildInfo = voGuildInfo
	}

	return &pb.GuildApplyResp{
		IsJoined:  isJoined,
		GuildInfo: guildInfo,
	}, nil
}

func (h *GuildHandler) CancelApply(ctx context.Context, req *pb.GuildCancelApplyReq) (*pb.GuildCancelApplyResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildCancelApply(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildCancelApplyResp{}, nil
}

func (h *GuildHandler) Quit(ctx context.Context, req *pb.GuildQuitReq) (*pb.GuildQuitResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildQuit(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildQuitResp{
		// GuildInfo: guild.VOGuildInfo(),
	}, nil
}

func (h *GuildHandler) HandleApplied(ctx context.Context, req *pb.GuildHandleAppliedReq) (*pb.GuildHandleAppliedResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildHandleApplied(ctx, req.UserID, req.Approve, req.Refuse)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildHandleAppliedResp{
		GuildInfo: voGuildInfo,
	}, nil
}

func (h *GuildHandler) Kick(ctx context.Context, req *pb.GuildKickReq) (*pb.GuildKickResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildKick(ctx, req.UserID, req.KickedUserID)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildKickResp{
		GuildInfo: voGuildInfo,
	}, nil
}

// 升职
func (h *GuildHandler) Promotion(ctx context.Context, req *pb.GuildPromotionReq) (*pb.GuildPromotionResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildPromotion(ctx, req.UserID, req.Target)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildPromotionResp{
		GuildInfo: voGuildInfo,
	}, nil
}

// 降职
func (h *GuildHandler) Demotion(ctx context.Context, req *pb.GuildDemotionReq) (*pb.GuildDemotionResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildDemotion(ctx, req.UserID, req.Target)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildDemotionResp{
		GuildInfo: voGuildInfo,
	}, nil
}

// 降职
func (h *GuildHandler) Transfer(ctx context.Context, req *pb.GuildTransferReq) (*pb.GuildTransferResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildTransfer(ctx, req.UserID, req.Target)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildTransferResp{
		GuildInfo: voGuildInfo,
	}, nil
}

// 是否解散
func (h *GuildHandler) IsDissolved(ctx context.Context, req *pb.GuildIsDissolvedReq) (*pb.GuildIsDissolvedResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	isDissolved, isInDissolving, err := sess.GuildIsDissolved(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.GuildIsDissolvedResp{
		IsDissolved:    isDissolved,
		IsInDissolving: isInDissolving,
	}, nil
}

func (h *GuildHandler) Sync(ctx context.Context, req *pb.GuildSyncReq) (*pb.GuildSyncResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildSync(ctx, req.UserID, req.LastLoginTime, req.Status)
	if err != nil {
		return nil, err
	}

	return &pb.GuildSyncResp{}, nil
}

func (h *GuildHandler) Members(ctx context.Context, req *pb.GuildMembersReq) (*pb.GuildMembersResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	return &pb.GuildMembersResp{
		UserIds: sess.GuildMembers(ctx),
	}, nil
}

// --------------联合建筑---------------
func (h *GuildHandler) CoBuildGetInfo(ctx context.Context, req *pb.GuildCoBuildGetInfoReq) (*pb.GuildCoBuildGetInfoResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	voBuild, err := sess.GuildCoBuildGetInfo(ctx, req.BuildID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildCoBuildGetInfoResp{
		CoBuild: voBuild,
	}, nil
}

func (h *GuildHandler) CoBuildImprove(ctx context.Context, req *pb.GuildCoBuildImproveReq) (*pb.GuildCoBuildImproveResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	coBuild, buildNotExecuted, err := sess.GuildCoBuildImprove(ctx, req.BuildID, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildCoBuildImproveResp{
		CoBuild:         coBuild,
		BuildNotExcuted: buildNotExecuted,
	}, nil
}

func (h *GuildHandler) CoBuildUse(ctx context.Context, req *pb.GuildCoBuildUseReq) (*pb.GuildCoBuildUseResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	voBuild, err := sess.GuildCoBuildUse(ctx, req.BuildID, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pb.GuildCoBuildUseResp{
		CoBuild: voBuild,
	}, nil

}

// func (h *GuildHandler) GetDispatchCharac(ctx context.Context, req *pb.GuildGetDispatchCharacReq) (*pb.GuildGetDispatchCharacResp, error) {
// 	sess, err := h.GetSession(ctx,req.GuildID)
// 	if err != nil {
// 		return nil, errors.WrapTrace(err)
// 	}
// 	sess.Lock()
// 	defer sess.Unlock()

// 	voGuildCharacs, err := sess.GetGuildCharacters(ctx, req.CharacterId)
// 	if err != nil {
// 		return nil, errors.WrapTrace(err)
// 	}
// 	return &pb.GuildGetDispatchCharacResp{
// 		GuildCharacters: voGuildCharacs,
// 	}, nil
// }

// todo 处理请求
func (h *GuildHandler) AddGreetings(ctx context.Context, req *pb.GuildAddGreetingsReq) (*pb.GuildAddGreetingsResp, error) {

	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	sess.Lock()
	defer sess.Unlock()

	err = sess.AddGreetings(ctx, req.UserId, req.Greetings)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.GuildAddGreetingsResp{}, nil
}

func (h *GuildHandler) GetGreetings(ctx context.Context, req *pb.GuildGetGreetingsReq) (*pb.GuildGetGreetingsResp, error) {

	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	sess.Lock()
	defer sess.Unlock()

	voGreetings, err := sess.MemberGreetings(ctx, req.UserId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.GuildGetGreetingsResp{
		Greetings: voGreetings,
	}, nil
}

func (h *GuildHandler) UpdateGreetingsRecord(ctx context.Context, req *pb.GuildUpdateGreetingsReq) (*pb.GuildUpdateGreetingsResp, error) {

	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	sess.Lock()
	defer sess.Unlock()

	err = sess.UpdateMember(ctx, req.UserId, req.LastTimestamp)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.GuildUpdateGreetingsResp{}, nil
}

func (h *GuildHandler) GetGraveyardRequests(ctx context.Context, req *pb.GuildGetGraveyardRequestsReq) (*pb.GuildGetGraveyardRequestsResp, error) {

	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	sess.Lock()
	defer sess.Unlock()

	vos, err := sess.GraveyardRequest(ctx, req.UserId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.GuildGetGraveyardRequestsResp{
		Requests: vos,
	}, nil
}

func (h *GuildHandler) AddGraveyardRequest(ctx context.Context, req *pb.GuildAddGraveyardRequestReq) (*pb.GuildAddGraveyardRequestResp, error) {

	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	sess.Lock()
	defer sess.Unlock()

	requestUid, err := sess.AddGraveyardRequest(ctx, req.Add)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.GuildAddGraveyardRequestResp{
		RequestUid: requestUid,
	}, nil
}

func (h *GuildHandler) HelpRequestsHandle(ctx context.Context, req *pb.GuildHelpRequestsHandleReq) (*pb.GuildHelpRequestsHandleResp, error) {

	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	sess.Lock()
	defer sess.Unlock()

	return sess.HelpRequestsHandle(ctx, req.UserId, req.DailyAddActivation, req.DailyAddGold)

}

func (h *GuildHandler) GetElites(ctx context.Context, req *pb.GuildGetElitesReq) (*pb.GuildGetElitesResp, error) {
	sess, err := h.GetSession(ctx, req.GuildId)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	return &pb.GuildGetElitesResp{
		UserIds: sess.GuildElite(ctx),
	}, nil
}

// -------公会推荐列表------------
func (h *GuildHandler) RecommendInfo(ctx context.Context, req *pb.GuildRecommendInfoReq) (*pb.GuildRecommendInfoResp, error) {
	sess, err := h.GetSession(ctx, req.GuildId)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	info, err := sess.RecommendInfo(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.GuildRecommendInfoResp{
		Info: info,
	}, nil
}

// ------------公会任务------------------

func (h *GuildHandler) GetTaskList(ctx context.Context, req *pb.GuildTaskListReq) (*pb.GuildTaskListResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	// 成员才可以获取公会任务信息
	err = sess.CheckIsMember(req.UserID)
	if err != nil {
		return nil, err
	}

	voGuildTask := sess.VOGuildTaskList(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pb.GuildTaskListResp{
		Tasks: voGuildTask,
	}, nil
}

func (h *GuildHandler) TaskRewards(ctx context.Context, req *pb.GuildTaskRewardsReq) (*pb.GuildTaskRewardsResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	// 检查是否为成员
	err = sess.CheckIsMember(req.UserID)
	if err != nil {
		return nil, err
	}
	// 检查进度是否达到
	activation, err := sess.ReceiveTaskRewards(req.TaskId, req.UserID)
	if err != nil {
		return nil, err
	}
	return &pb.GuildTaskRewardsResp{
		Activation: activation,
	}, nil
}

func (h *GuildHandler) TaskAddProgress(ctx context.Context, req *pb.GuildTaskAddProgressReq) (*pb.GuildTaskAddProgressResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	// 检查是否为成员
	err = sess.CheckIsMember(req.UserID)
	if err != nil {
		return nil, err
	}
	for _, item := range req.Items {
		err = sess.GuildTaskAddProgress(item.TaskId, item.Value)
		if err != nil {
			return nil, err
		}
	}
	return &pb.GuildTaskAddProgressResp{}, nil
}

// 全体邮件
func (h *GuildHandler) SendGroupMail(ctx context.Context, req *pb.GuildSendGroupMailReq) (*pb.GuildSendGroupMailResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildSendGroupMail(ctx, req.UserID, req.Title, req.Content, req.Sender)
	if err != nil {
		return nil, err
	}

	return &pb.GuildSendGroupMailResp{}, nil
}

func (h *GuildHandler) Modify(ctx context.Context, req *pb.GuildModifyReq) (*pb.GuildModifyResp, error) {
	sess, err := h.GetSession(ctx, req.GuildId)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()
	err = sess.GuildModify(ctx, req.UserId, req.Icon, req.JoinModel, req.Title)
	if err != nil {
		return nil, err
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GuildModifyResp{
		Info: voGuildInfo,
	}, nil
}

func (h *GuildHandler) Chat(ctx context.Context, req *pb.GuildChatReq) (*pb.GuildChatResp, error) {
	sess, err := h.GetSession(ctx, req.GuildID)
	if err != nil {
		return nil, err
	}

	sess.Lock()
	defer sess.Unlock()

	err = sess.GuildChat(ctx, req.UserID, req.UserName, req.Content, req.Avatar, req.Frame)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	voGuildInfo, err := sess.VOGuildInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	return &pb.GuildChatResp{
		GuildInfo: voGuildInfo,
	}, nil
}
