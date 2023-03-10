package model

import (
	"context"
	"time"

	"guild/manager"
	"shared/common"
	"shared/csv/static"
	"shared/global"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/mysql"
)

const (
	GuildJoinModelAuto   = 0 // 公会自由加入
	GuildJoinModelHandle = 1 // 公会加入需要审批

	GuildChatJoin = 1 // 公会加入
	GuildChatQuit = 2 // 退出公会

	GuildNormalTaskNum   = 2 // 每次普通任务池随机出公会任务数量
	GuildSeparateTaskNum = 1 // 每次特殊任务池随机出公会任务数量
	GuildTaskTotalNum    = 8 // 公会任务总数量
)

type Guild struct {
	ID                     int64                 `db:"id" major:"true"`           // 公会ID
	Name                   string                `db:"name"`                      // 名称
	Chairman               int64                 `db:"chairman"`                  // 会长
	ViceChairmen           []int64               `db:"vice_chairmen"`             // 副会长
	Icon                   *common.GuildIcon     `db:"icon"`                      // 图标
	Title                  string                `db:"title"`                     // 说明
	Exp                    int32                 `db:"exp"`                       // 经验
	Level                  int32                 `db:"level"`                     // 等级
	JoinModel              int32                 `db:"join_model"`                // 加入模式：自由加入、手动审批
	Members                []*common.GuildMember `db:"members"`                   // 成员
	Chats                  []common.GuildChat    `db:"chats"`                     // 聊天框
	AppliedList            []common.GuildApply   `db:"applied_list"`              // 申请列表
	DissolveTime           int64                 `db:"dissolve_time"`             // 解散时间戳，可以24h反悔
	LastTimeCancelDissolve int64                 `db:"last_time_cancel_dissolve"` // 上次取消解散的时间
	HelpList               []int64               `db:"help_list"`                 // 帮助列表
	CreateTime             int64                 `db:"ctime"`                     // 创建时间
	ClearFlag              int32                 `db:"clear_flag"`                // 是否清除数据，解散后记录清除数据标记
	ChairmanLastLoginTime  int64                 `db:"chairman_last_login_time"`  // 会长上次登录时间，过久不登录自动转让会长
	MemberLastLoginTime    int64                 `db:"member_last_login_time"`    // 成员上次登录时间，过久没人不登录自动解散公会
	Tasks                  []*common.GuildTask   `db:"tasks"`                     // 任务
	TaskLastRefreshTime    int64                 `db:"task_last_refresh_time"`    // 上次任务刷新时间
	LastSendGroupMailTime  int64                 `db:"last_send_group_mail_time"` // 上次集体邮件时间
	// Area                  int32                 `db:"area"`                     // 区
	// BossDamage            *number.CalNumber     `db:"boss_damage"`              // BOSS伤害
	HelpRequests       *common.GraveyardHelpRequests `db:"help_requests"`   // 公会互助请求
	YggCoBuilds        map[int32]*YggCoBuild         `db:"-"`               // 世界探索公会联合建筑
	GuildGreetings     *GuildGreetings               `db:"guild_greetings"` // 角色异界问候
	*mysql.EmbedModule `db:"-"`
}

func NewGuild(id int64) *Guild {
	return &Guild{
		ID:             id,
		Name:           "",
		Chairman:       0,
		ViceChairmen:   []int64{},
		Icon:           &common.GuildIcon{},
		Title:          "",
		Exp:            0,
		Level:          0,
		JoinModel:      0,
		Members:        []*common.GuildMember{},
		Chats:          []common.GuildChat{},
		AppliedList:    []common.GuildApply{},
		DissolveTime:   0,
		HelpList:       []int64{},
		CreateTime:     time.Now().Unix(),
		ClearFlag:      0,
		HelpRequests:   common.NewGraveyardRequests(),
		EmbedModule:    &mysql.EmbedModule{},
		GuildGreetings: NewGreetings(),
	}
}

func (g *Guild) InitNewGuild(userID int64) bool {
	g.Chairman = userID

	member := common.NewGuildMember(userID)

	// 初始化在线状态
	member.Status = common.UserOnline
	member.LastLoginTime = time.Now().Unix()
	member.Position = common.GuildPositionChairman

	g.Members = append(g.Members, member)
	g.Level = 1

	return false
}

func (g *Guild) CheckInDissolving() error {
	if g.DissolveTime > 0 {
		return common.ErrGuildInDissolving
	}

	return nil
}

// 检查公会是否解散，人死了什么都没了，公会也是
func (g *Guild) CheckDissolved() error {
	if g.IsDissolved() {
		return common.ErrGuildDissolved
	}

	return nil
}

// 检查公会解散CD，返回值为true，代表还在CD中，不可解散
func (g *Guild) CheckDissolvedCD() bool {
	return g.LastTimeCancelDissolve > 0 && time.Unix(g.LastTimeCancelDissolve, 0).Add(time.Duration(manager.CSV.GlobalEntry.GuildDissolveCD)*time.Second).After(time.Now())
}

func (g *Guild) IsCleared() bool {
	return g.ClearFlag == 0
}

func (g *Guild) Clear() {
	g.ClearFlag = 1
}

func (g *Guild) IsDissolved() bool {
	return g.DissolveTime > 0 && time.Now().Add(-time.Duration(manager.CSV.GlobalEntry.GuildDissolveDuration)*time.Second).After(time.Unix(g.DissolveTime, 0))
}

func (g *Guild) CheckIsMember(userID int64) error {
	var isInGuild bool
	for _, member := range g.Members {
		if member.UserID == userID {
			isInGuild = true
		}
	}
	if !isInGuild {
		return errors.Swrapf(common.ErrGuildNotMember, g.ID, userID)
	}

	return nil
}

func (g *Guild) CheckPrivilegeChairMan(userID int64) error {
	return g.CheckPrivilege(userID, common.GuildPositionChairman)
}

func (g *Guild) CheckPrivilegeViceChairMan(userID int64) error {
	return g.CheckPrivilege(userID, common.GuildPositionViceChairman)
}

func (g *Guild) CheckSendGroupMailCD() error {
	// TODO：read config
	if !time.Now().Add(-24 * time.Hour).After(time.Unix(g.LastSendGroupMailTime, 0)) {
		return common.ErrGuildSendGroupMailInCD
	}

	return nil
}

func (g *Guild) IsChairMan(userID int64) bool {
	return g.Chairman == userID
}

func (g *Guild) IsViceChairMan(userID int64) bool {
	for _, v := range g.ViceChairmen {
		if v == userID {
			return true
		}
	}

	return false
}

func (g *Guild) IsElite(userID int64) bool {
	for _, v := range g.Members {
		if v.UserID == userID {
			return v.Position == common.GuildPositionElite
		}
	}

	return false
}

func (g *Guild) IsCommon(userID int64) bool {
	for _, v := range g.Members {
		if v.UserID == userID {
			return v.Position == common.GuildPositionCommon
		}
	}

	return false
}

func (g *Guild) EliteNum() int32 {
	var result int32
	for _, v := range g.Members {
		if v.Position == common.GuildPositionElite {
			result++
		}
	}
	return result
}

func (g *Guild) CheckPrivilege(userID int64, position int32) error {
	switch position {
	case common.GuildPositionChairman:
		if g.IsChairMan(userID) {
			return nil
		}

		return common.ErrGuildNoPrivilege
	case common.GuildPositionViceChairman:
		if g.IsChairMan(userID) {
			return nil
		}

		if g.IsViceChairMan(userID) {
			return nil
		}

		return common.ErrGuildNoPrivilege
	case common.GuildPositionElite:
		if g.IsChairMan(userID) {
			return nil
		}

		if g.IsViceChairMan(userID) {
			return nil
		}

		if g.IsElite(userID) {
			return nil
		}

		return common.ErrGuildNoPrivilege
	case common.GuildPositionCommon:
		return nil
	}
	return nil
}

func (g *Guild) IsAutoHandle() bool {
	return g.JoinModel == GuildJoinModelAuto
}

func (g *Guild) Apply(userID int64) {
	// 添加到申请列表
	g.AppliedList = append(g.AppliedList, *common.NewGuildApply(userID))
}

func (g *Guild) CancelApply(userID int64) {
	// 从申请列表移除
	for i, v := range g.AppliedList {
		if v.UserID == userID {
			g.AppliedList = append(g.AppliedList[:i], g.AppliedList[i+1:]...)
			break
		}
	}
}

func (g *Guild) Quit(ctx context.Context, userID int64) {
	// 从成员列表移除
	for i, v := range g.Members {
		if v.UserID == userID {
			g.Members = append(g.Members[:i], g.Members[i+1:]...)
			break
		}
	}

	// 从副会长列表移除
	for i, v := range g.ViceChairmen {
		if v == userID {
			g.ViceChairmen = append(g.ViceChairmen[:i], g.ViceChairmen[i+1:]...)
			break
		}
	}

	userCaches, err := manager.Global.GetUserCaches(ctx, []int64{userID})
	if err != nil {
		glog.Errorf("GetUserCaches(%v) error", userID, err)
	}
	if userCache, ok := userCaches[userID]; ok {
		g.ChatEvent(userID, userCache.Name, static.GuildActionGuildChatQuit)
	}
	g.SendGreetingsByMails(ctx, []int64{userID})
}

func (g *Guild) CheckMemberCount(approveCount int) error {
	if len(g.Members)+approveCount >= int(manager.CSV.GlobalEntry.GuildMemberMaxNum) {
		return common.ErrGuildIsFull
	}

	return nil
}

func (g *Guild) CheckAppliedListCount(maxCount int) error {
	if len(g.AppliedList) >= maxCount {
		return common.ErrGuildIsFull
	}

	return nil
}

func (g *Guild) CheckInAppliedList(userIDs []int64) error {
	appliedMap := map[int64]bool{}
	for _, v := range g.AppliedList {
		appliedMap[v.UserID] = true
	}

	for _, userID := range userIDs {
		if !appliedMap[userID] {
			return common.ErrGuildNotInAppliedList
		}
	}

	return nil
}

func (g *Guild) Approve(ctx context.Context, userID int64) {
	member := common.NewGuildMember(userID)

	// 初始化在线状态
	caches, err := manager.Global.GetUserCachesExtension(ctx, []int64{userID}, global.UserCacheWithOnline)
	if err != nil {
		glog.Errorf("GetUserCachesExtension(%v) error", userID, err)
	}

	cache, ok := caches[userID]
	if ok {
		member.LastLoginTime = cache.LastLoginTime
		member.Status = int8(cache.OnlineStatus)
	}

	// 添加成员列表
	g.Members = append(g.Members, member)
	g.ChatEvent(userID, cache.Name, static.GuildActionGuildChatJoin)
}

func (g *Guild) RefreshLoginTime(userID, lastLoginTime int64, status int8) {
	if g.IsChairMan(userID) {
		g.ChairmanLastLoginTime = lastLoginTime
	}

	for i, member := range g.Members {
		if member.UserID == userID {
			g.Members[i].LastLoginTime = lastLoginTime
			g.Members[i].Status = status
		}
	}

	g.MemberLastLoginTime = lastLoginTime
}

func (g *Guild) RefreshOnlineStatus() {
	// 玩家上线和下线都会主动同步状态，当玩家在线的时候需要持续自己去获取状态同步
	for i, member := range g.Members {

		if member.Status == common.UserOnline {
			// 如果玩家在线超过5分钟没有同步，就同步一次，防止服务器意外崩溃没有同步状态
			now := time.Now()
			if (now.Add(-5 * time.Minute)).After(time.Unix(member.LastLoginTime, 0)) {
				lastLoginTime, err := manager.Global.UserLastLoginTime(context.Background(), member.UserID)
				if err != nil {
					glog.Errorf("UserLastLoginTime(%d) error: %v", member.UserID, err)
				} else {
					// 不在线，更新状态
					if lastLoginTime == 0 {
						g.Members[i].LastLoginTime = now.Unix()
						g.Members[i].Status = common.UserOffline
					}
				}
			}
		}
	}
}

func (g *Guild) MemberPosition(userID int64) int32 {
	for _, v := range g.Members {
		if v.UserID == userID {
			return v.Position
		}
	}

	return 0
}

func (g *Guild) Chat(maxChatLen int, userID int64, userName, content string, avatar, frame int32) {
	// 添加成员列表
	g.Chats = append(g.Chats, *common.NewGuildChat(0, userID, userName, g.MemberPosition(userID), content, avatar, frame))

	// 会话长度超了，移除最旧的会话
	if len(g.Chats) > maxChatLen {
		g.Chats = g.Chats[1:]
	}
}
func (g *Guild) ChatEvent(userId int64, userName string, eventType int8) {
	g.Chats = append(g.Chats, *common.NewGuildChat(eventType, userId, userName, 0, "", 0, 0))

	if len(g.Chats) > int(manager.CSV.GlobalEntry.GuildChatLimit) {
		g.Chats = g.Chats[1:]
	}
}

// 升职，不做检查
func (g *Guild) Promotion(userID int64, position int32) {
	for i, member := range g.Members {
		if member.UserID == userID {
			g.Members[i].Position = position
		}
	}

	// 升职成为副会长
	if position == common.GuildPositionViceChairman {
		g.ViceChairmen = append(g.ViceChairmen, userID)
	}
}

// 降职，不做检查
func (g *Guild) Demotion(userID int64, position int32) {
	for i, member := range g.Members {
		if member.UserID == userID {
			g.Members[i].Position = position
		}
	}

	// 从副会长降职
	if position == common.GuildPositionElite {
		for i, member := range g.ViceChairmen {
			if member == userID {
				g.ViceChairmen = append(g.ViceChairmen[:i], g.ViceChairmen[i+1:]...)
			}
		}
	}
}

func (g *Guild) Transfer(ctx context.Context, userID int64) {
	oldChairman := g.Chairman

	g.Chairman = userID

	// 从副会长删除
	for i, member := range g.ViceChairmen {
		if member == userID {
			g.ViceChairmen = append(g.ViceChairmen[:i], g.ViceChairmen[i+1:]...)
		}
	}

	// 会长变成副会长
	g.ViceChairmen = append(g.ViceChairmen, oldChairman)

	// 职位变化
	for i, member := range g.Members {
		if member.UserID == userID {
			g.Members[i].Position = common.GuildPositionChairman
		} else if member.UserID == oldChairman {
			g.Members[i].Position = common.GuildPositionViceChairman
		}
	}

	userCaches, err := manager.Global.GetUserCaches(ctx, []int64{oldChairman, g.Chairman})
	if err != nil {
		glog.Errorf("GetUserCaches(%v) error", userID, err)
	}
	if userCache, ok := userCaches[oldChairman]; ok {
		g.ChatEvent(userID, userCache.Name, static.GuildActionGuildChatUnchairman)
	}
	if userCache, ok := userCaches[g.Chairman]; ok {
		g.ChatEvent(userID, userCache.Name, static.GuildActionGuildChatChairman)
	}
}

func (g *Guild) HandleApplied(users []int64) {
	userM := map[int64]bool{}

	for _, userID := range users {
		userM[userID] = true
	}

	// 添加成员列表
	for i, v := range g.AppliedList {
		if userM[v.UserID] {
			g.AppliedList = append(g.AppliedList[:i], g.AppliedList[i+1:]...)
		}
	}
}

func (g *Guild) VOGuildInfo(ctx context.Context, userId int64) (*pb.VOGuildInfo, error) {
	intimacyMap, err := manager.Global.GetGuildIntimacyMap(ctx, g.ID, userId)
	if err != nil {
		return nil, err
	}
	uids := make([]int64, 0, len(g.Members))
	voMembers := make([]*pb.VOGuildMember, 0, len(g.Members))
	for i, member := range g.Members {
		voMembers = append(voMembers, g.Members[i].VOGuildMember(intimacyMap[member.UserID]))
		uids = append(uids, member.UserID)
	}

	// 添加成员其他数据
	userCaches, err := manager.Global.GetUserCaches(ctx, uids)
	if err != nil {
		return nil, err
	}

	for i, member := range voMembers {
		if userCache, ok := userCaches[member.UserID]; ok {
			voMembers[i].UserName = userCache.Name
			voMembers[i].Avatar = userCache.Avatar
			voMembers[i].Frame = userCache.Frame
			voMembers[i].Power = userCache.Power
		}
	}

	voChats := make([]*pb.VOGuildChat, 0, len(g.Chats))
	for i, _ := range g.Chats {
		voChats = append(voChats, g.Chats[i].VOGuildChat())
	}

	return &pb.VOGuildInfo{
		GuildID:                g.ID,
		Name:                   g.Name,
		Chairman:               g.Chairman,
		ViceChairmen:           g.ViceChairmen,
		Icon:                   g.Icon.VOGuildIcon(),
		Title:                  g.Title,
		Exp:                    g.Exp,
		Level:                  g.Level,
		JoinModel:              g.JoinModel,
		CreateTime:             g.CreateTime,
		Members:                voMembers,
		Chats:                  voChats,
		DissolveTime:           g.DissolveTime,
		LastCancelDissolveTime: g.LastTimeCancelDissolve,
	}, nil
}

// 返回公会列表详情
func (g *Guild) VOGuildShowInfo(ctx context.Context) (*pb.VOGuildShowInfo, error) {
	voChairmen := make([]*pb.VOGuildMember, 0, 1+manager.CSV.GlobalEntry.GuildViceChairmenNum)
	uids := make([]int64, 0, 1+manager.CSV.GlobalEntry.GuildViceChairmenNum)
	for i, member := range g.Members {
		if member.Position >= common.GuildPositionViceChairman {
			voChairmen = append(voChairmen, g.Members[i].VOGuildMember(0))
			uids = append(uids, member.UserID)
		}
	}

	// 添加成员其他数据
	userCaches, err := manager.Global.GetUserCaches(ctx, uids)
	if err != nil {
		return nil, err
	}

	for i, member := range voChairmen {
		if userCache, ok := userCaches[member.UserID]; ok {
			// fmt.Println("===============================", userCache.ID, userCache.Name)
			voChairmen[i].UserName = userCache.Name
			voChairmen[i].Avatar = userCache.Avatar
			voChairmen[i].Frame = userCache.Frame
			voChairmen[i].Power = userCache.Power
		}
	}
	icon := &pb.VOGuildIcon{}
	if g.Icon != nil {
		icon = g.Icon.VOGuildIcon()
	}
	return &pb.VOGuildShowInfo{
		GuildID:      g.ID,
		Name:         g.Name,
		Chairman:     g.Chairman,
		ViceChairmen: g.ViceChairmen,
		Icon:         icon,
		Title:        g.Title,
		Exp:          g.Exp,
		Level:        g.Level,
		JoinModel:    g.JoinModel,
		CreateTime:   g.CreateTime,
		ChairmenInfo: voChairmen,
	}, nil
}

func (g *Guild) AddGuildExp(addExp int32) {
	if addExp <= 0 {
		return
	}
	expArr := manager.CSV.Guild.GetExpArr()
	nowLevel := g.Level
	maxLevel := int32(len(expArr))
	if nowLevel >= maxLevel {
		return
	}
	nowExp := g.Exp + addExp

	for i := nowLevel - 1; i < maxLevel; i++ {
		exp := expArr[i]
		if exp == nowExp {
			nowLevel = i + 1
			break
		} else if exp > nowExp {
			break
		}
		nowLevel = i + 1
	}

	if nowLevel >= maxLevel {
		// 满级后经验就不会再加了,停在满级0经验
		g.Exp = expArr[maxLevel-1]
		g.Level = maxLevel

	} else {
		g.Exp = nowExp
		g.Level = nowLevel

	}

}

func (g *Guild) AddGuildMemberActivation(member *common.GuildMember, addExp int32) {
	// 增加 贡献度（公会等级经验）
	g.AddGuildExp(addExp * manager.CSV.Guild.GetGuildContributionExpRatio())
	// 增加个人贡献度
	member.Activation.Plus(addExp)
}

func (g *Guild) GetMemberUserIDs() []int64 {
	userIDs := make([]int64, 0, len(g.Members))

	for _, v := range g.Members {
		userIDs = append(userIDs, v.UserID)
	}

	return userIDs
}
