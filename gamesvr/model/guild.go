package model

import (
	"context"
	"encoding/json"
	"time"

	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/global"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/servertime"
)

type UserGuild struct {
	GuildID           int64                     `json:"guild_id"`            // 公会ID
	GuildName         string                    `json:"guild_name"`          // 公会名称
	ApplyList         []*common.GuildSimpleInfo `json:"apply_list"`          // 申请列表
	LastQuitTime      int64                     `json:"last_quit_time"`      // 上次退会时间
	LastCheckDissolve int64                     `json:"last_check_dissolve"` // 上次确认公会是否解散的时间，公会解散有24h，所以只需要24h查询一次即可
	GuildGold         *common.BIEventNumber     `json:"guild_gold"`          // 公会代币
	Tasks             []*TaskItem               `json:"-"`
	GuildGreetings    Greetings                 `json:"guild_greetings"`
}

type TaskItem struct {
	TaskId int32
	Value  int32
}

func NewUserGuild() *UserGuild {
	guildGoldLimit := manager.CSV.Item.GetLimit(static.CommonResourceTypeGuildMoney)
	return &UserGuild{
		GuildID:           0,
		GuildName:         "",
		ApplyList:         []*common.GuildSimpleInfo{},
		LastQuitTime:      0,
		LastCheckDissolve: 0,
		GuildGold:         common.NewBIEventNumber(0, 0, guildGoldLimit),
		Tasks:             []*TaskItem{},
		GuildGreetings:    Greetings{},
	}
}

func NewTaskItem(id, value int32) *TaskItem {
	return &TaskItem{
		TaskId: id,
		Value:  value,
	}
}
func (t *TaskItem) VOGuildTaskItem() *pb.VOGuildTaskItem {
	return &pb.VOGuildTaskItem{
		TaskId: t.TaskId,
		Value:  t.Value,
	}
}

func (g *UserGuild) AddTaskItem(id, value int32) {
	item := NewTaskItem(id, value)
	g.Tasks = append(g.Tasks, item)
}

func (g *UserGuild) ClearTaskItems() {
	g.Tasks = []*TaskItem{}
}

func (g *UserGuild) RefreshApplied(ctx context.Context, userID int64) error {
	if g.GuildID == 0 {
		// 如果没有公会，并且有在申请中的公会，检查有无公会通过审批
		if len(g.ApplyList) > 0 {
			userGuild := &global.UserGuild{}
			err := manager.Global.HGetAllScan(ctx, userID, userGuild)
			if err != nil {
				return errors.WrapTrace(err)
			}

			// 有工会通过审批
			if userGuild.GuildID != 0 {
				g.Join(userGuild.GuildID, userGuild.GuildName)
				return errors.WrapTrace(err)
			}
		}
	}

	return nil
}

func (g *UserGuild) HasJoinedGuild() bool {
	return g.GuildID != 0
}

func (g *UserGuild) CheckHasJoinedGuild() error {
	if g.HasJoinedGuild() {
		return common.ErrGuildHasJoined
	}

	return nil
}

func (g *UserGuild) CheckHasNotJoinGuild() error {
	if !g.HasJoinedGuild() {
		return common.ErrGuildHasNotJoin
	}

	return nil
}
func (g *UserGuild) CheckApplyNumLimit() error {
	if len(g.ApplyList) >= 10 {
		return common.ErrGuildExceedApplyNumLimit
	}
	return nil
}

func (g *UserGuild) CheckQuitGuildCD() error {
	if time.Unix(g.LastQuitTime, 0).Add(time.Duration(manager.CSV.GlobalEntry.GuildQuitCD) * time.Second).After(time.Now()) {
		return common.ErrGuildQuitCD
	}

	return nil
}

func (g *UserGuild) NeedCheckDissolved() bool {
	return g.HasJoinedGuild() && time.Unix(g.LastCheckDissolve, 0).Add(time.Duration(manager.CSV.GlobalEntry.GuildDissolveDuration)*time.Second).Before(time.Now())
}

func (g *UserGuild) Apply(info *pb.VOGuildSimpleInfo) {
	simple := &common.GuildSimpleInfo{}
	simple.LoadFromVOGuildSimpleInfo(info)
	g.ApplyList = append(g.ApplyList, simple)
}

func (g *UserGuild) CancelApply(guildID int64) {
	for i, applyGuild := range g.ApplyList {
		if applyGuild.GuildID == guildID {
			g.ApplyList = append(g.ApplyList[:i], g.ApplyList[i+1:]...)
		}
	}
}

func (g *UserGuild) CheckIsAlreadyApply(guildID int64) bool {
	for _, apply := range g.ApplyList {
		if guildID == apply.GuildID {
			return true
		}
	}
	return false
}

func (g *UserGuild) Join(guildID int64, guildName string) {
	g.ApplyList = []*common.GuildSimpleInfo{}
	g.GuildID = guildID
	g.GuildName = guildName

	// 可加入的公会一定不是处于解散状态的
	g.LastCheckDissolve = servertime.Now().Unix()
}

func (g *UserGuild) Quit() {
	g.GuildID = 0
	g.GuildName = ""
}

func (g *UserGuild) Marshal() ([]byte, error) {
	return json.Marshal(g)
}

func (g *UserGuild) Unmarshal(bs []byte) error {
	if len(bs) == 0 {
		return nil
	}

	return json.Unmarshal(bs, g)
}

func (g *UserGuild) VOGuildInfo() *pb.VOUserGuild {
	list := make([]int64, 0, len(g.ApplyList))
	for _, guild := range g.ApplyList {
		list = append(list, guild.GuildID)
	}

	return &pb.VOUserGuild{
		GuildID:      g.GuildID,
		GuildName:    g.GuildName,
		ApplyList:    list,
		LastQuitTime: g.LastQuitTime,
	}
}
