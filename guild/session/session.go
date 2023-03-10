package session

import (
	"context"
	"database/sql"
	"log"
	"math"
	"sync"
	"time"

	"guild/manager"
	"guild/model"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/safe"
	"shared/utility/servertime"

	"shared/utility/session"
)

type Builder struct{}

func (b *Builder) NewSession() session.Session {
	return &Session{
		EmbedManagedSession: &session.EmbedManagedSession{},
		Guild:               &model.Guild{},
	}
}

type Session struct {
	sync.RWMutex
	*session.EmbedManagedSession
	*model.Guild
}

func (s *Session) OnCreated(ctx context.Context, opts session.OnCreatedOpts) error {
	s.Guild = model.NewGuild(opts.ID)
	s.Guild.RefreshTask()
	err := manager.MySQL.Load(ctx, s.Guild)
	if err != nil {
		if err == sql.ErrNoRows {
			if !opts.AllowNil {
				return common.ErrGuildNotFound
			}
		} else {
			return err
		}
	}

	s.ScheduleCall(5*time.Second, func() {
		safe.Exec(5, func(i int) error {
			err := manager.MySQL.Save(context.Background(), s.Guild)
			if err != nil {
				log.Printf("save user error: %v, times: %d ,\n", err, i)
			}

			return err
		})
	})

	return nil
}

func (s *Session) OnClosed() {
	s.ScheduleCall(5*time.Second, func() {
		safe.Exec(5, func(i int) error {
			err := manager.MySQL.Save(context.Background(), s.Guild)
			if err != nil {
				log.Printf("save user error: %v, times: %d ,\n", err, i)
			}

			return err
		})
	})

	log.Printf("guild session closed； %d", s.Guild.ID)
}

// func (s *Session) GuildInfo(ctx context.Context, userID int64) (*model.Guild, error) {
// 	// if !s.IsMember(userID) {
// 	// 	return nil, errors.New("not member")
// 	// }
//
// 	return s.Guild, nil
// }

func (s *Session) GuildCreate(ctx context.Context, userID int64, name string, icon *pb.VOGuildIcon, joinModel int32, title string) error {
	// if s.IsLoaded {
	// 	// 公会已经创建
	// 	return errors.New("guild exist")
	// }

	// 初始化数据
	s.InitNewGuild(userID)
	s.Name = name
	s.Icon.LoadFromVOGuildIcon(icon)
	s.JoinModel = joinModel
	s.Title = title

	// 入库
	err := manager.MySQL.Create(ctx, s.Guild)
	if err != nil {
		log.Printf("new user error: %v", err)
		return err
	}

	err = manager.Global.SetUserGuildData(ctx, userID, s.ID, s.Name)
	if err != nil {
		return err
	}

	return nil
}

func (s *Session) GuildDissolve(ctx context.Context, userID int64) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 判断权限，会长才有权限解散公会
	err = s.CheckPrivilegeChairMan(userID)
	if err != nil {
		return err
	}

	if s.CheckDissolvedCD() {
		return errors.Swrapf(common.ErrGuildInDissolvedCD, s.ID)
	}

	// 一个人直接解散公会，不进行倒计时
	if len(s.Members) == 1 {
		s.DissolveTime = time.Now().Add(-24 * time.Hour).Unix()
	}
	for _, member := range s.Members {
		err := s.SendGreetingsToMember(ctx, member.UserID)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	// 设定解散时间
	s.DissolveTime = time.Now().Unix()

	return nil
}

func (s *Session) GuildCancelDissolve(ctx context.Context, userID int64) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 判断权限，会长才有权限解散公会
	err = s.CheckPrivilegeChairMan(userID)
	if err != nil {
		return err
	}

	// 设定解散时间
	s.DissolveTime = 0

	s.LastTimeCancelDissolve = servertime.Now().Unix()

	return nil
}

func (s *Session) GuildApply(ctx context.Context, userID int64) (bool, error) {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return false, err
	}

	// 检查公会是否已经在解散中
	err = s.CheckInDissolving()
	if err != nil {
		return false, err
	}

	lock, err := manager.Global.ObtainLock(ctx, userID)
	if err != nil {
		return false, err
	}
	defer lock.Release()

	guildID, err := manager.Global.HGetInt64(ctx, userID, "guild_id")
	if err != nil {
		if err == global.ErrNil {
			guildID = 0
		} else {
			return false, err
		}
	}

	// 检查玩家是否有公会，可能是申请的时候同时被其他公会通过，看谁抢到锁
	if guildID != 0 {
		return false, common.ErrGuildHasJoined
	}

	// 是否是自动加入
	if s.IsAutoHandle() {
		// 检查公会成员达到上限
		err := s.CheckMemberCount(1)
		if err != nil {
			return false, err
		}

		// 先设置公会ID和名称，这一步可能会因为网络异常失败，要保证加入公会后这一步必须成功
		// err = manager.Global.HSetAll(ctx, userID, &common.UserGuild{GetUserGuildID: s.ID, GetUserGuildName: s.Name})
		// if err != nil {
		// 	return false, err
		// }

		err = manager.Global.SetUserGuildData(ctx, userID, s.ID, s.Name)
		if err != nil {
			return false, err
		}

		// 加入公会
		s.Approve(ctx, userID)

		return true, nil
	}

	// 检查申请列表是否达到上限
	err = s.CheckAppliedListCount(100)
	if err != nil {
		return false, err
	}

	// 需要审批，加入审批列表
	s.Apply(userID)

	// 非自动加入返回空的公会信息
	return false, nil
}

func (s *Session) GuildCancelApply(ctx context.Context, userID int64) error {
	// 不报错，只是删除申请列表中的对应玩家

	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return nil
	}

	// 检查玩家还在不在申请列表
	err = s.CheckInAppliedList([]int64{userID})
	if err != nil {
		return nil
	}

	// 检查公会是否已经在解散中
	s.CancelApply(userID)

	return nil
}

func (s *Session) GuildQuit(ctx context.Context, userID int64) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 检查是否是公会成员，公会成员才能退出公会
	err = s.CheckIsMember(userID)
	if err != nil {
		return err
	}

	// 检查是否不是会长，会长不能退出公会
	if s.IsChairMan(userID) {
		return common.ErrGuildChairmanCantQuit
	}

	// 清除玩家公会信息
	// err = manager.Global.HSetAll(ctx, userID, &common.UserGuild{GetUserGuildID: 0, GetUserGuildName: ""})
	// if err != nil {
	// 	return err
	// }

	err = manager.Global.SetUserGuildData(ctx, userID, 0, "")
	if err != nil {
		return err
	}

	err = manager.Global.ClearIntimacy(ctx, s.ID, userID, s.GetMembers()...)
	if err != nil {
		return err
	}
	s.Quit(ctx, userID)

	return nil
}

func (s *Session) GuildHandleApplied(ctx context.Context, userID int64, approve, refuse []int64) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 判断权限，会长和副会长才有权限审批
	err = s.CheckPrivilegeViceChairMan(userID)
	if err != nil {
		return err
	}

	// 检查公会是否已经在解散中，不管通不通过都要消失
	err = s.CheckInDissolving()
	if err != nil {
		s.HandleApplied(append(approve, refuse...))
		return err
	}

	// 检查玩家还在不在申请列表
	err = s.CheckInAppliedList(append(approve, refuse...))
	if err != nil {
		return err
	}

	// 检查公会成员达到上限
	err = s.CheckMemberCount(len(approve))
	if err != nil {
		return err
	}

	f := func(userID int64) error {
		lock, err := manager.Global.ObtainLock(ctx, userID)
		if err != nil {
			return err
		}
		defer lock.Release()

		// guildID, err := manager.Global.HGetInt64(ctx, userID, "guild_id")
		// if err != nil {
		// 	if err == global.ErrNil {
		// 		guildID = 0
		// 	} else {
		// 		return err
		// 	}
		// }

		guild, err := manager.Global.GetUserGuildData(ctx, userID)
		if err != nil {
			return err
		}

		// 检查玩家是否有公会，可能通过的时候同时被其他公会通过，看谁抢到锁
		if guild.GuildID != 0 {
			return common.ErrGuildHasJoined
		}

		// 先设置公会ID和名称，这一步可能会因为网络异常失败，要保证加入公会后这一步必须成功
		// err = manager.Global.HSetAll(ctx, userID, &common.UserGuild{GetUserGuildID: s.ID, GetUserGuildName: s.Name})
		// if err != nil {
		// 	return err
		// }

		err = manager.Global.SetUserGuildData(ctx, userID, s.ID, s.Name)
		if err != nil {
			return err
		}

		// 加入公会
		s.Approve(ctx, userID)

		return nil
	}

	for _, approveUserID := range approve {
		err := f(approveUserID)
		if err != nil {
			return err
		}
	}

	// 从申请列表删除
	s.HandleApplied(append(approve, refuse...))

	return nil
}

func (s *Session) GuildGetApplyList(ctx context.Context, userID int64) ([]*pb.VOGuildApplyUser, error) {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return nil, err
	}

	// 判断权限，会长和副会长才能看
	err = s.CheckPrivilegeViceChairMan(userID)
	if err != nil {
		return nil, err
	}

	ret := make([]*pb.VOGuildApplyUser, 0, len(s.AppliedList))
	uids := make([]int64, 0, len(s.AppliedList))
	for _, v := range s.AppliedList {
		uids = append(uids, v.UserID)
	}

	userCaches, err := manager.Global.GetUserCaches(ctx, uids)
	if err != nil {
		return nil, err
	}

	for _, v := range s.AppliedList {
		userCache, ok := userCaches[v.UserID]
		if !ok {
			continue
		}

		ret = append(ret, &pb.VOGuildApplyUser{
			UserID:    userCache.ID,
			UserName:  userCache.Name,
			Power:     userCache.Power,
			Avatar:    userCache.Avatar,
			Frame:     userCache.Frame,
			Level:     userCache.Level,
			ApplyTime: v.ApplyTime,
		})
	}

	return ret, nil
}

func (s *Session) GuildKick(ctx context.Context, userID, kickedUserID int64) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 判断权限，会长和副会长才有权限审批
	err = s.CheckPrivilegeViceChairMan(userID)
	if err != nil {
		return err
	}

	// 判断职位，普通成员才可以被提
	if !s.IsCommon(kickedUserID) {
		return common.ErrGuildDissolved
	}

	err = manager.Global.SetUserGuildData(ctx, kickedUserID, 0, "")
	if err != nil {
		return err
	}

	err = manager.Global.ClearIntimacy(ctx, s.ID, kickedUserID, s.GetMembers()...)
	if err != nil {
		return err
	}

	s.Quit(ctx, kickedUserID)

	return nil
}

func (s *Session) GuildChat(ctx context.Context, userID int64, userName, content string, avatar, frame int32) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 检查是否是公会成员，公会成员才能聊天
	err = s.CheckIsMember(userID)
	if err != nil {
		return err
	}

	s.Chat(int(manager.CSV.GlobalEntry.GuildChatLimit), userID, userName, content, avatar, frame)

	return nil
}

func (s *Session) GuildPromotion(ctx context.Context, userID, target int64) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	position := s.MemberPosition(target)

	switch position {
	case common.GuildPositionChairman,
		common.GuildPositionViceChairman: // 会长和副会长不能升职

		return common.ErrGuildNoPrivilege
	case common.GuildPositionElite: // 升职成为副会长
		// 会长才有权限提升副会长
		err := s.CheckPrivilegeChairMan(userID)
		if err != nil {
			return err
		}
		if len(s.ViceChairmen) >= int(manager.CSV.GlobalEntry.GuildViceChairmenNum) {
			return errors.Swrapf(common.ErrGuildViceCharimanNumOutOfLimit, s.ID)
		}

		s.Promotion(target, common.GuildPositionViceChairman)
	case common.GuildPositionCommon: // 升职成为精英
		// 会长和副会长才有权限提升精英
		err := s.CheckPrivilegeViceChairMan(userID)
		if err != nil {
			return err
		}
		if s.EliteNum() >= manager.CSV.GlobalEntry.GuildElitenNum {
			return errors.Swrapf(common.ErrGuildEliteNumOutOfLimit, s.ID)
		}

		s.Promotion(target, common.GuildPositionElite)
	}

	return nil
}

func (s *Session) GuildDemotion(ctx context.Context, userID, target int64) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	position := s.MemberPosition(target)

	switch position {
	case common.GuildPositionChairman,
		common.GuildPositionCommon: // 会长和普通成员不能降职

		return common.ErrGuildNoPrivilege
	case common.GuildPositionElite: // 精英降为普通成员
		// 会长和副会长才有权限降职精英
		err := s.CheckPrivilegeViceChairMan(userID)
		if err != nil {
			return err
		}

		s.Demotion(target, common.GuildPositionCommon)
	case common.GuildPositionViceChairman: // 副会长降为精英
		// 会长才有权限降职副会长
		err := s.CheckPrivilegeChairMan(userID)
		if err != nil {
			return err
		}

		s.Demotion(target, common.GuildPositionElite)
	}

	return nil
}

func (s *Session) GuildTransfer(ctx context.Context, userID, target int64) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 会长才有权限转让会长
	err = s.CheckPrivilegeChairMan(userID)
	if err != nil {
		return err
	}

	s.Transfer(ctx, target)

	return nil
}

func (s *Session) GuildIsDissolved(ctx context.Context) (bool, bool, error) {
	if s.IsDissolved() {
		if !s.IsCleared() {
			// 删除 公会名称-公会ID 索引
			err := manager.Global.DelGuildName(ctx, s.Name)
			if err != nil {
				return false, false, err
			}
			// 删除 公会set
			err = manager.Global.GuildSetDelete(ctx, s.ID)
			if err != nil {
				return false, false, err
			}

			s.Clear()
		}

		return true, false, nil
	}

	return false, s.Guild.DissolveTime > 0, nil
}

func (s *Session) GuildSync(ctx context.Context, userID int64, lastLoginTime int64, status int32) error {
	s.RefreshLoginTime(userID, lastLoginTime, int8(status))
	return nil
}

// GuildMembers 获取公会所有成员
func (s *Session) GuildMembers(ctx context.Context) []int64 {
	userIds := make([]int64, 0, len(s.Guild.Members))
	for _, member := range s.Guild.Members {
		userIds = append(userIds, member.UserID)
	}
	return userIds
}

func (s *Session) GuildElite(ctx context.Context) []int64 {
	userIds := make([]int64, 0, len(s.Guild.Members))
	for _, member := range s.Guild.Members {
		if member.Position >= common.GuildPositionElite {
			userIds = append(userIds, member.UserID)
		}
	}
	return userIds
}

// -----------联合建筑------------
func (s *Session) GuildCoBuildImprove(ctx context.Context, buildId int32, uid int64) (*pb.VOGuildYggdrasilCoBuild, bool, error) {
	coBuild, buildNotExecuted, err := s.YggCoBuildImprove(buildId, uid)
	if err != nil {
		return nil, false, errors.WrapTrace(err)
	}
	return coBuild.VOGuildYggdrailCoBuild(), buildNotExecuted, nil
}

// todo 目前来看过程中没有error可返回
func (s *Session) GuildCoBuildGetInfo(ctx context.Context, buildID int32) (*pb.VOGuildYggdrasilCoBuild, error) {
	coBuild := s.Guild.YggCoBuildGetInfo(buildID)
	return coBuild.VOGuildYggdrailCoBuild(), nil
}

func (s *Session) GuildCoBuildUse(ctx context.Context, buildId int32, uid int64) (*pb.VOGuildYggdrasilCoBuild, error) {
	coBuild, err := s.YggCoBuildUse(buildId, uid)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return coBuild.VOGuildYggdrailCoBuild(), nil
}

//-------------公会派遣----------------
// func (s *Session) GetGuildCharacters(ctx context.Context, characId int32) ([]*pb.VOYggDispatchGuildCharacter, error) {
// 	guildCharacs, err := s.Guild.GetGuildCharacters(ctx, characId)
// 	if err != nil {
// 		return nil, errors.WrapTrace(err)
// 	}

// 	return guildCharacs, nil
// }

// 公会异界问候
func (s *Session) AddGreetings(ctx context.Context, userId int64, greetings []*pb.VOGreetings) error {
	err := s.CheckDissolved()
	if err != nil {
		return nil
	}

	err = s.Guild.AddGreetings(ctx, userId, greetings)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (s *Session) GetMemberGreetings(ctx context.Context, userId int64) ([]*pb.VOGreetings, error) {
	voGreetings, err := s.Guild.MemberGreetings(ctx, userId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return voGreetings, err
}

func (s *Session) UpdateGreetingsRecord(ctx context.Context, userId int64, lastTimestamp int64) error {
	err := s.Guild.UpdateMember(ctx, userId, lastTimestamp)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (s *Session) AddGraveyardRequest(ctx context.Context, add *pb.VOAddGraveyardRequest) (int64, error) {

	err := s.CheckIsMember(add.UserId)
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	sec := int32(math.Round(float64(add.TotalSec) * manager.CSV.GraveyardEntry.GetGraveyardMinHelpPercent()))
	minSec := manager.CSV.GraveyardEntry.GetGraveyardMinHelpSec()
	if sec < minSec {
		sec = minSec
	}

	request := common.NewGraveyardHelpRequest(add.UserId, add.BuildUid, int(add.HelpType), sec, 0, add.BuildId, add.BuildLv, add.ExpireAt)

	requestUid := s.Guild.HelpRequests.Add(request)
	return requestUid, nil
}

func (s *Session) GraveyardRequest(ctx context.Context, userId int64) ([]*pb.VOGraveyardHelpRequest, error) {

	m := map[int64]struct{}{}
	for _, member := range s.Guild.Members {
		m[member.UserID] = struct{}{}
	}
	requests := s.Guild.HelpRequests.GetRequestsExcept(userId, m, manager.CSV.GraveyardEntry.GetGraveyardSingleHelpCount())
	info, err := s.VOGuildInfo(ctx, userId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	members := map[int64]*pb.VOGuildMember{}
	for _, member := range info.Members {
		members[member.UserID] = member
	}
	vos := make([]*pb.VOGraveyardHelpRequest, 0, len(requests))
	for _, request := range requests {
		vos = append(vos, request.VOGraveyardHelpRequest(members))
	}
	return vos, nil
}

func (s *Session) HelpRequestsHandle(ctx context.Context, userId int64, dailyAddActivation int32, dailyAddGold int32) (*pb.GuildHelpRequestsHandleResp, error) {
	m := map[int64]struct{}{}
	for _, member := range s.Guild.Members {
		m[member.UserID] = struct{}{}
	}
	requests := s.Guild.HelpRequests.GetRequestsExcept(userId, m, manager.CSV.GraveyardEntry.GetGraveyardSingleHelpCount())

	var member *common.GuildMember
	for _, tmp := range s.Members {
		if tmp.UserID == userId {
			member = tmp
		}
	}
	if member == nil {
		return nil, errors.Swrapf(common.ErrGuildNotMember, s.ID, userId)
	}

	var helpCount int32
	for _, request := range requests {
		guildLevel, err := manager.CSV.Guild.GetGuildLevel(s.Guild.Level)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		// 互助事件
		manager.EventQueue.Push(ctx, request.UserId, common.NewGraveyardHelpEvent(request.HelpType, request.BuildUid, request.Sec))
		// 已帮助次数1
		request.HelpedCount++
		// 设置本条已帮助
		request.HelpedUserIds[userId] = struct{}{}
		// 增加的友好度
		_, err = manager.Global.ChangeIntimacy(ctx, s.ID, userId, request.UserId, guildLevel.SingleHelpAddIntimacy)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		// 达到每日上限不加经验
		if dailyAddActivation < guildLevel.DailyActivationUpperLimitByHelp {
			max := guildLevel.DailyActivationUpperLimitByHelp - dailyAddActivation
			addActivation := guildLevel.SingleHelpAddActivation
			if max > addActivation {
				addActivation = max
			}
			s.Guild.AddGuildMemberActivation(member, addActivation)
			dailyAddActivation += addActivation
		}
		// 达到每日上限不加代币
		if dailyAddGold < guildLevel.DailyGoldUpperLimitByHelp {
			max := guildLevel.DailyGoldUpperLimitByHelp - dailyAddGold
			addGold := guildLevel.SingleHelpAddGold
			if max > addGold {
				addGold = max
			}
			// 增加 公会代币
			dailyAddGold += addGold
		}
		helpCount++
	}
	info, err := s.Guild.VOGuildInfo(ctx, userId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return &pb.GuildHelpRequestsHandleResp{
		DailyAddGold:       dailyAddGold,
		DailyAddActivation: dailyAddActivation,
		HelpCount:          helpCount,
		GuildInfo:          info,
	}, nil
}

// ----------公会推荐列表------------
func (s *Session) RecommendInfo(ctx context.Context) (*pb.VOGuildSimpleInfo, error) {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return nil, err
	}
	icon := &pb.VOGuildIcon{}
	if s.Icon != nil {
		icon = s.Icon.VOGuildIcon()
	}
	// 返回数据
	return &pb.VOGuildSimpleInfo{
		GuildID:     s.ID,
		Name:        s.Name,
		Icon:        icon,
		Level:       s.Level,
		JoinModel:   s.JoinModel,
		MemberCount: int32(len(s.Members)),
		CreateTime:  s.CreateTime,
	}, nil
}

func (s *Session) GuildSendGroupMail(ctx context.Context, userID int64, title, content, sender string) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 会长和副会长才有权限
	err = s.CheckPrivilegeViceChairMan(userID)
	if err != nil {
		return err
	}

	// 检查发送cd
	err = s.CheckSendGroupMailCD()
	if err != nil {
		return err
	}

	userIDs := s.GetMemberUserIDs()

	// 不包括自己
	for i, v := range userIDs {
		if v == userID {
			userIDs = append(userIDs[:i], userIDs[i+1:]...)
			break
		}
	}

	// TODO: read config ,not hardcode
	var guildGroupMailID int32 = 8
	mailTemplate, err := manager.CSV.Mail.GetTemplate(guildGroupMailID)
	if err != nil {
		return err
	}

	now := servertime.Now()
	expire := now.Add(time.Duration(mailTemplate.ExpireDays) * 24 * time.Hour).Unix()

	_, err = manager.RPCMailClient.SendGroupMail(ctx, &pb.SendGroupMailReq{
		TemplateId: guildGroupMailID,
		Title:      title,
		Content:    content,
		Sender:     sender,
		StartTime:  now.Unix(),
		ExpireTime: expire,
		EndTime:    expire,
		Users:      userIDs,
	})

	return err
}

func (s *Session) GuildModify(ctx context.Context, userID int64, icon *pb.VOGuildIcon, joinModel int32, title string) error {
	// 检查公会是否已经解散
	err := s.CheckDissolved()
	if err != nil {
		return err
	}

	// 判断权限，会长才有权限修改信息
	err = s.CheckPrivilegeChairMan(userID)
	if err != nil {
		return err
	}
	s.Icon.LoadFromVOGuildIcon(icon)
	s.JoinModel = joinModel
	s.Title = title

	return nil
}
