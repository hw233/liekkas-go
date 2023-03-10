package common

import (
	"time"

	"shared/protobuf/pb"
	"shared/utility/number"
	"shared/utility/servertime"
)

const (
	// 公会成员身份
	GuildPositionSystem       = 0 // 系统
	GuildPositionCommon       = 1 // 普通会员
	GuildPositionElite        = 2 // 精英
	GuildPositionViceChairman = 3 // 副会长
	GuildPositionChairman     = 4 // 话事人
)

// type Guild struct {
// 	GuildID             int32      `db:"id" where:"=" major:"true"`
// 	Name           string     `db:"name"`
// 	Chairman       int64      `db:"chairman"`
// 	ViceChairMen   []int64    `db:"vice_chairmen"`
// 	Icon           *GuildIcon `db:"icon"`
// 	Title          string     `db:"title"`
// 	Exp            int32      `db:"exp"`
// 	Level          int32      `db:"level"`
// 	JoinModel      int32      `db:"join_model"`
// 	CreateTime     int64      `db:"create_time"`
// 	Members        []int64    `db:"members"`
// 	Chats          []string   `db:"chats"`
// 	*mysql.Handler `db:"-"`
// }

// type Guild struct {
// 	GuildID           int32           `json:"id" redis:"id"`
// 	Name         string          `json:"name" redis:"name"`
// 	Chairman     *GuildChairman  `json:"chairman" redis:"chairman"`
// 	ViceChairMan []GuildChairman `json:"vice_chairman" redis:"vice_chairman"`
// 	Icon         *GuildIcon      `json:"icon" redis:"icon"`
// 	Title        string          `json:"chairman" redis:"chairman"`
// 	Level        int32           `json:"level" redis:"level"`
// 	JoinModel    int32           `json:"join_model" redis:"join_model"`
// 	CreateTime   int64           `json:"create_time" redis:"create_time"`
// 	MemberCount  int64           `json:"member_count" redis:"member_count"`
// 	Member       []int64         `json:"member" redis:"member"`
// }

// 公会成员
type GuildMember struct {
	UserID              int64             `json:"uid"`
	Position            int32             `json:"position"`        // 职位
	Activation          *number.CalNumber `json:"activation"`      // 活跃度，贡献值
	LastLoginTime       int64             `json:"last_login_time"` // 上次登录时间
	Status              int8              `json:"status"`          // 状态：是否在线
	TaskRewardsReceived map[int32]bool    `json:"task_rewards_received"`
}

func NewGuildMember(userID int64) *GuildMember {
	return &GuildMember{
		UserID:              userID,
		Position:            GuildPositionCommon,
		Activation:          number.NewCalNumber(0),
		LastLoginTime:       0,
		Status:              UserOffline,
		TaskRewardsReceived: map[int32]bool{},
	}
}

func (gm *GuildMember) VOGuildMember(intimacy int32) *pb.VOGuildMember {
	return &pb.VOGuildMember{
		UserID:        gm.UserID,
		Position:      gm.Position,
		Activation:    gm.Activation.Value(),
		LastLoginTime: gm.LastLoginTime,
		IsOnline:      gm.Status == UserOnline,
		Intimacy:      intimacy,
	}
}

// 公会任务
type GuildTask struct {
	ID    int32             `json:"id"`
	Count *number.CalNumber `json:"count"`
}

func NewGuildTask(id int32) *GuildTask {
	return &GuildTask{
		ID:    id,
		Count: number.NewCalNumber(0),
	}
}

func (t *GuildTask) VOGuildTask(isReceived bool) *pb.VOGuildTask {
	return &pb.VOGuildTask{
		Id:         t.ID,
		Count:      t.Count.Value(),
		IsReceived: isReceived,
	}
}

// 申请列表
type GuildApply struct {
	UserID    int64 `json:"user_id"`    // 申请者id
	ApplyTime int64 `json:"apply_time"` // 申请时间
}

func NewGuildApply(userID int64) *GuildApply {
	return &GuildApply{
		UserID:    userID,
		ApplyTime: servertime.Now().Unix(),
	}
}

// 公会聊天
type GuildChat struct {
	UserID     int64  `json:"user_id"`   // 发布者id
	UserName   string `json:"user_name"` // 发布者名称
	Type       int8   `json:"type"`      // 发布类型
	CreateTime int64  `json:"ctime"`     // 发布时间
	Position   int32  `json:"position"`  // 发布者公会职位
	Content    string `json:"content"`   // 发布内容
	Avatar     int32  `json:"avatar"`
	Frame      int32  `json:"frame"`
}

func NewGuildChat(typ int8, userID int64, userName string, position int32, content string, avatar, frame int32) *GuildChat {
	return &GuildChat{
		UserID:     userID,
		UserName:   userName,
		Type:       typ,
		CreateTime: time.Now().Unix(),
		Position:   position,
		Content:    content,
		Avatar:     avatar,
		Frame:      frame,
	}
}

func (gc *GuildChat) VOGuildChat() *pb.VOGuildChat {
	return &pb.VOGuildChat{
		UserID:     gc.UserID,
		UserName:   gc.UserName,
		Position:   gc.Position,
		Content:    gc.Content,
		CreateTime: gc.CreateTime,
		ChatType:   int32(gc.Type),
		Avatar:     gc.Avatar,
		Frame:      gc.Frame,
	}
}

// 公会图标
type GuildIcon struct {
	DecorateID      int32  `json:"decorate_id"`      // 装饰id
	DecorateColor   string `json:"decorate_color"`   // 装饰颜色
	BackgroundID    int32  `json:"background_id"`    // 背景id
	BackgroundColor string `json:"background_color"` // 背景颜色
	PictureID       int32  `json:"picture_id"`       // 图片id
	PictureColor    string `json:"picture_color"`    // 图片颜色
}

func (gi *GuildIcon) VOGuildIcon() *pb.VOGuildIcon {
	return &pb.VOGuildIcon{
		DecorateID:      gi.DecorateID,
		DecorateColor:   gi.DecorateColor,
		BackgroundID:    gi.BackgroundID,
		BackgroundColor: gi.BackgroundColor,
		PictureID:       gi.PictureID,
		PictureColor:    gi.PictureColor,
	}
}

func (gi *GuildIcon) LoadFromVOGuildIcon(vo *pb.VOGuildIcon) {
	gi.DecorateID = vo.DecorateID
	gi.DecorateColor = vo.DecorateColor
	gi.BackgroundID = vo.BackgroundID
	gi.BackgroundColor = vo.BackgroundColor
	gi.PictureID = vo.PictureID
	gi.PictureColor = vo.PictureColor

}

// type GuildChairman struct {
// 	UID   int32 `json:"uid" redis:"uid"`
// 	Name  int32 `json:"name" redis:"name"`
// 	Level int32 `json:"level" redis:"level"`
// 	Head  int32 `json:"head" redis:"head"`
// }

// 公会基本信息
type GuildSimpleInfo struct {
	GuildID     int64      `json:"guild_id"`
	GuildName   string     `json:"guild_name"`
	Icon        *GuildIcon `json:"icon"`
	Level       int32      `json:"level"`
	JoinModel   int32      `json:"join_model"`
	MemberCount int32      `json:"member_count"`
	CreateTime  int64      `json:"ctime"`
}

func (g *GuildSimpleInfo) LoadFromVOGuildSimpleInfo(info *pb.VOGuildSimpleInfo) {
	icon := &GuildIcon{}
	if info.Icon != nil {
		icon.LoadFromVOGuildIcon(info.Icon)
	}
	g.GuildID = info.GuildID
	g.GuildName = info.Name
	g.Icon = icon
	g.Level = info.Level
	g.JoinModel = info.JoinModel
	g.MemberCount = info.MemberCount
	g.CreateTime = info.GetCreateTime()
}

func (g *GuildSimpleInfo) VOGuildSimpleInfo() *pb.VOGuildSimpleInfo {
	icon := &pb.VOGuildIcon{}
	if g.Icon != nil {
		icon = g.Icon.VOGuildIcon()
	}
	return &pb.VOGuildSimpleInfo{
		GuildID:     g.GuildID,
		Name:        g.GuildName,
		Icon:        icon,
		Level:       g.Level,
		JoinModel:   g.JoinModel,
		MemberCount: g.MemberCount,
		CreateTime:  g.CreateTime,
	}
}
