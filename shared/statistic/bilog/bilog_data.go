package bilog

import (
	"bytes"
	"math"
	"shared/common"
	"shared/statistic/logreason"
	"shared/utility/rand"
	"strconv"
	"time"

	"go.uber.org/zap/zapcore"
)

const (
	EvenNameCreate          = "create_role"
	EventNameLogin          = "player_login"
	EventNameLogout         = "player_logout"
	EventNameGuide          = "guide_flow"
	EventNameExp            = "player_exp"
	EventNameResourceChange = "gold_flow"
	EventNameItemChange     = "item_flow"
	EventNameGraveyard      = "home_flow"
	EventNameQuest          = "mission_flow"
	EventNameGacha          = "gacha"
	EventNameLevel          = "stage_flow"
	EventNameCharaGet       = "card_flow"
	EventNameCharaOp        = "card_op_flow"
	EventNameHeroOp         = "master_op_flow"
	EventNameEquipmentOp    = "equip_op_flow"
	EventNameWorldItemOp    = "equip2_op_flow"
	EventNameMail           = "mail_flow"
	EventNameSnapshot       = "playersnapshot"
	EventNameOnlineCount    = "online_num"
)

func FormatLogId(appName, eventName string, uid int64, timeStamp int64) string {
	buf := bytes.NewBufferString(appName)
	buf.WriteString(eventName)
	buf.WriteString(strconv.Itoa(int(uid)))
	buf.WriteString(strconv.Itoa(int(timeStamp)))
	buf.WriteString(strconv.Itoa(int(rand.RangeInt32(1, math.MaxInt32))))

	return buf.String()
}

type UserCommon struct {
	Uid       int64
	EventName string
	LogId     string
	LogTime   time.Time
}

func (uc *UserCommon) GetLogId() string {
	return uc.LogId
}

func (uc *UserCommon) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("b_log_id", uc.LogId)
	encoder.AddString("b_account_id", strconv.Itoa(int(uc.Uid)))
	encoder.AddString("b_role_id", strconv.Itoa(int(uc.Uid)))
	encoder.AddInt64("b_utc_timestamp", uc.LogTime.Unix())
	encoder.AddString("b_datetime", uc.LogTime.Format("2006-01-02 15:04:05"))
	encoder.AddString("b_eventname", uc.EventName)

	return nil
}

type Reward common.Reward

func (r *Reward) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("id", r.ID)
	encoder.AddInt32("type", r.Type)
	encoder.AddInt32("count", r.Num)

	return nil
}

type Rewards common.Rewards

func (r *Rewards) MarshalLogArray(encoder zapcore.ArrayEncoder) error {
	for _, rewards := range *r {
		for _, reward := range rewards {
			encoder.AppendObject(((*Reward)(&reward)))
		}
	}

	return nil
}

//bi data

type UserCreate struct {
	RoleName  string
	Model     string
	OSVersion string
	Network   string
	Mac       string
	Ip        string
}

func (uc *UserCreate) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("role_name", uc.RoleName)
	encoder.AddString("model", uc.Model)
	encoder.AddString("os_version", uc.OSVersion)
	encoder.AddString("network", uc.Network)
	encoder.AddString("mac", uc.Mac)
	encoder.AddString("ip", uc.Ip)

	return nil
}

type UserLoginLog struct {
	RoleName      string
	Model         string
	OSVersion     string
	Network       string
	Mac           string
	Ip            string
	Udid          string
	SDKUdid       string
	SDKUid        string
	GameBaseId    string
	GameId        int64
	Platform      string
	ZoneId        int64
	ChannelId     int64
	ClientVersion string
}

func (ull *UserLoginLog) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("role_name", ull.RoleName)
	encoder.AddString("model", ull.Model)
	encoder.AddString("os_version", ull.OSVersion)
	encoder.AddString("network", ull.Network)
	encoder.AddString("mac", ull.Mac)
	encoder.AddString("ip", ull.Ip)
	encoder.AddString("b_udid", ull.Udid)
	encoder.AddString("b_sdk_udid", ull.SDKUdid)
	encoder.AddString("b_sdk_uid", ull.SDKUid)
	encoder.AddString("b_game_base_id", ull.GameBaseId)
	encoder.AddInt64("b_game_id", ull.GameId)
	encoder.AddString("b_platform", ull.Platform)
	encoder.AddInt64("b_zone_id", ull.ZoneId)
	encoder.AddInt64("b_channel_id", ull.ChannelId)
	encoder.AddString("b_version", ull.ClientVersion)

	return nil
}

type UserLogoutLog struct {
	OnlineTime int64
	Gender     int8
	Model      string
	OSVersion  string
	Network    string
	Mac        string
	Ip         string
}

func (ull *UserLogoutLog) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt64("online_time", ull.OnlineTime)
	encoder.AddInt8("gender", ull.Gender)
	encoder.AddString("model", ull.Model)
	encoder.AddString("os_version", ull.OSVersion)
	encoder.AddString("network", ull.Network)
	encoder.AddString("mac", ull.Mac)
	encoder.AddString("ip", ull.Ip)

	return nil
}

type UserGuide struct {
	GuideId int32
}

func (ug *UserGuide) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("guide_id", ug.GuideId)

	return nil
}

type UserResourceChange struct {
	ResourceId int32
	ChangeNum  int64
	BeforeNum  int64
	AfterNum   int64
	Reason     logreason.Reason
}

func (urc *UserResourceChange) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	absChangeNum := urc.ChangeNum
	if urc.ChangeNum > 0 {
		encoder.AddString("act_type", "add")
	} else {
		absChangeNum = -absChangeNum
		encoder.AddString("act_type", "reduce")
	}
	encoder.AddString("related_event_logid", urc.Reason.RelateLog())
	encoder.AddInt32("gold_id", urc.ResourceId)
	encoder.AddInt64("gold_num", absChangeNum)
	encoder.AddInt64("before_count", urc.BeforeNum)
	encoder.AddInt64("after_count", urc.AfterNum)
	encoder.AddInt32("reason", urc.Reason.Id())
	encoder.AddString("related_order_id", urc.Reason.RelateOrder())

	return nil
}

type UserItemChange struct {
	ItemId    int32
	ChangeNum int32
	BeforeNum int32
	AfterNum  int32
	Reason    logreason.Reason
}

func (urc *UserItemChange) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	absChangeNum := urc.ChangeNum
	if urc.ChangeNum > 0 {
		encoder.AddString("act_type", "add")
	} else {
		absChangeNum = -absChangeNum
		encoder.AddString("act_type", "reduce")
	}
	encoder.AddString("related_event_logid", urc.Reason.RelateLog())
	encoder.AddInt32("item_id", urc.ItemId)
	encoder.AddInt32("item_num", absChangeNum)
	encoder.AddInt32("before_count", urc.BeforeNum)
	encoder.AddInt32("after_count", urc.AfterNum)
	encoder.AddInt32("reason", urc.Reason.Id())
	encoder.AddString("subreason", "")
	encoder.AddString("related_order_id", urc.Reason.RelateOrder())

	return nil
}

type UserExpLevel struct {
	Optype string
	Exp    int32
	Level  int32
	Reason *logreason.Reason
}

func (uel *UserExpLevel) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("change_type", uel.Optype)
	encoder.AddInt32("exp", uel.Exp)
	encoder.AddInt32("level", uel.Level)
	encoder.AddInt32("reason", uel.Reason.Id())

	return nil
}

type UserGetCharacter struct {
	CharaId    int32
	CharaLevel int32
	CharaStar  int32
	CharaStage int32
	CharaExp   int32
	Reason     *logreason.Reason
}

func (ugc *UserGetCharacter) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("act_type", "add")
	encoder.AddInt32("reason", ugc.Reason.Id())
	encoder.AddString("related_event_logid", ugc.Reason.RelateLog())
	encoder.AddInt32("card_id", ugc.CharaId)
	encoder.AddInt32("card_level", ugc.CharaLevel)
	encoder.AddInt32("card_rank", ugc.CharaStage)
	encoder.AddInt32("card_star", ugc.CharaStar)
	encoder.AddInt32("card_exp", ugc.CharaExp)

	return nil
}

const (
	CharaOpExp = iota
	CharaOpLevel
	CharaOpStage
	CharaOpStar
)

type UserCharaOperationP struct {
	OpType  int8
	CharaId int32
	Level   int32
	Star    int32
	Stage   int32
	Exp     int32
}

func (uco *UserCharaOperationP) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt8("act_type", uco.OpType)
	encoder.AddInt32("card_id", uco.CharaId)
	encoder.AddInt32("card_level", uco.Level)
	encoder.AddInt32("card_rank", uco.Star)
	encoder.AddInt32("card_exp", uco.Exp)

	return nil
}

type UserCharaOperationCP struct {
	SkillLevels map[int32]int32
}

func charaSkillLevelEncoder(skillNum, level int32) zapcore.ObjectMarshalerFunc {
	return zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
		encoder.AddInt32("skill_num", skillNum)
		encoder.AddInt32("skill_level", level)
		return nil
	})
}

func charaSkillsEncoder(skillLevels map[int32]int32) zapcore.ArrayMarshalerFunc {
	return zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
		for skillNum, level := range skillLevels {
			encoder.AppendObject(charaSkillLevelEncoder(skillNum, level))
		}

		return nil
	})
}

func (uco *UserCharaOperationCP) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddArray("card_skill", charaSkillsEncoder(uco.SkillLevels))

	return nil
}

const (
	GraveyardOpBuild = iota
	GraveyardOpLevelUp
	GraveyardOpStageUp
	GraveyardOpStartProduce
	GraveyardOpGetProduction
	GraveyardOpAccelerate
	GraveyardOpCharaDispatch
)

type UserGraveyard struct {
	OpType         int
	MainTowerLevel int32
	BuildingId     int32
	BuildingLevel  int32
	BuildingStage  int32
	Charas         []int32
	Duration       int64
}

func (ug *UserGraveyard) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt("act_type", ug.OpType)
	encoder.AddInt32("facility_id", ug.BuildingId)
	encoder.AddInt32("city_lv", ug.MainTowerLevel)
	encoder.AddInt32("facility_lv", ug.BuildingLevel)
	encoder.AddInt32("facility_rank", ug.BuildingId)
	encoder.AddArray("addcard_id", int32SliceEncoder(ug.Charas))

	return nil
}

const (
	LevelsResultFail = iota
	LevelsResultWin
)

type UserLevelsProperties struct {
	LevelId   int32
	LevelType int32
	IsWin     bool
	Score     int32
	IsFirst   bool
}

func (ulp *UserLevelsProperties) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("stage_id", ulp.LevelId)
	encoder.AddInt32("stage_type", ulp.LevelType)
	encoder.AddInt8("stage_result", boolToInt8(ulp.IsWin))
	encoder.AddInt32("stage_score", ulp.Score)
	encoder.AddInt8("is_first", boolToInt8(ulp.IsFirst))

	return nil
}

type UserQuestData struct {
	QuestId    int32
	QuestName  string
	QuestType  int32
	QuestState int32
	UserExp    int32
	UserPower  int32
}

func (uqd *UserQuestData) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("mission_id", uqd.QuestId)
	encoder.AddString("mission_name", uqd.QuestName)
	encoder.AddInt32("mission_type", uqd.QuestType)
	encoder.AddInt32("mission_status", uqd.QuestState)
	encoder.AddInt32("exp", uqd.UserExp)
	encoder.AddInt32("power", uqd.UserPower)

	return nil
}

type UserGachaData struct {
	PoolId     int32
	PoolType   int32
	IsTen      bool
	Reward     *Rewards
	Cost       *Rewards
	TotalTimes int32
}

const (
	gachaModeSingle = 1
	gachaModeTen    = 2
)

func (ugd *UserGachaData) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("gacha_id", ugd.PoolId)
	encoder.AddInt32("gacha_type", ugd.PoolType)
	if ugd.IsTen {
		encoder.AddInt32("gacha_mode", gachaModeTen)
	} else {
		encoder.AddInt32("gacha_mode", gachaModeSingle)
	}
	encoder.AddArray("award", ugd.Reward)
	encoder.AddArray("cost", ugd.Cost)
	encoder.AddInt32("total_gacha", ugd.TotalTimes)

	return nil
}

const (
	HeroOpLevelUp = iota
	HeroOpSkillLevelUp
	HeroOpAddChara
	HeroOpUnlock
)

type UserHeroCp struct {
	OpType     int8
	HeroId     int32
	HeroLevel  int32
	HeroSkills map[int32]int32
	Charas     []int32
}

func (uhc *UserHeroCp) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt8("act_type", uhc.OpType)
	encoder.AddInt32("master_id", uhc.HeroId)
	encoder.AddInt32("master_level", uhc.HeroLevel)
	encoder.AddArray("master_skill", charaSkillsEncoder(uhc.HeroSkills))
	encoder.AddArray("pos_detail", int32SliceEncoder(uhc.Charas))

	return nil
}

const (
	EquipmentOpGet = iota
	EquipmentOpLevelUp
	EquipmentOpStageUp
	EquipmentOpWear
	EquipmentOpTakeOff
)

type UserEquipmentP struct {
	OpType  int8
	Id      int64
	Eid     int32
	CharaId int32
}

func (uep *UserEquipmentP) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("card_id", uep.CharaId)
	encoder.AddInt8("act_type", uep.OpType)
	encoder.AddInt64("equip_uid", uep.Id)
	encoder.AddInt32("equip_id", uep.Eid)

	return nil
}

type UserEquipmentCP struct {
	Stage  int32
	Attrs  map[int8]int32
	Reason *logreason.Reason
}

func EquipmentAttrEncoder(attrId int8, value int32) zapcore.ObjectMarshalerFunc {
	return zapcore.ObjectMarshalerFunc(func(encoder zapcore.ObjectEncoder) error {
		encoder.AddInt8("attr_id", attrId)
		encoder.AddInt32("value", value)
		return nil
	})
}

func EquipmentAttrsEncoder(equipmentAttrs map[int8]int32) zapcore.ArrayMarshalerFunc {
	return zapcore.ArrayMarshalerFunc(func(encoder zapcore.ArrayEncoder) error {
		for attrId, value := range equipmentAttrs {
			encoder.AppendObject(EquipmentAttrEncoder(attrId, value))
		}

		return nil
	})
}

func (uep *UserEquipmentCP) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("reason", uep.Reason.Id())
	encoder.AddInt32("equip_rank", uep.Stage)
	encoder.AddArray("equip_extra", EquipmentAttrsEncoder(uep.Attrs))

	return nil
}

const (
	WorldItemOpGet = iota
	WorldItemOpLevelUp
	WorldItemOpStageUp
	WorldItemOpWear
	WorldItemOpTakeOff
)

type UserWorldItemCP struct {
	OpType  int8
	Id      int64
	Wid     int32
	CharaId int32
	Stage   int32
	Reason  *logreason.Reason
}

func (uwi *UserWorldItemCP) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("card_id", uwi.CharaId)
	encoder.AddInt8("act_type", uwi.OpType)
	encoder.AddInt64("equip2_uid", uwi.Id)
	encoder.AddInt32("equip2_id", uwi.Wid)
	encoder.AddInt32("reason", uwi.Reason.Id())
	encoder.AddInt32("equip2_rank", uwi.Stage)

	return nil
}

const (
	MailOpSend = iota
	MailOpReceive
	MailOpRead
	MailOpRemove
	MailOpReceiveAttachment
)

type UserMail struct {
	OpType               int8
	TemplateId           int32
	SenderId             int64
	MailId               int64
	Title                string
	Content              string
	SendTime             int64
	ExpireTime           int64
	IsAttachmentReceived bool
	Attachment           *Rewards
}

func (um *UserMail) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt8("act_type", um.OpType)
	encoder.AddInt32("mail_type", um.TemplateId)
	encoder.AddInt64("from_id", um.SenderId)
	encoder.AddInt64("mail_id", um.MailId)
	encoder.AddString("mail_title", um.Title)
	encoder.AddString("mail_content", um.Content)
	encoder.AddInt64("mail_createtime", um.SendTime)
	encoder.AddInt64("mail_deltime", um.ExpireTime)
	encoder.AddInt8("item_got", boolToInt8(um.IsAttachmentReceived))
	if um.Attachment != nil {
		encoder.AddArray("bind_item_num", um.Attachment)
	}

	return nil
}

//----------------------------------------
//snapshot
//----------------------------------------
type UserSnapshot struct {
	Time            time.Time
	UserCreateTime  time.Time
	LastLoginTime   time.Time
	GuildId         int64
	Level           int32
	Exp             int32
	Power           int32
	FirstPayTime    time.Time
	FirstOrderId    string
	FirstOrderLevel int32
	PayDiamond      int32
	FreeDiamond     int32
	TotalPay        int64
	TotalPayTimes   int32
	TotalLoginDay   int32
	TotalOnlineTime int64
}

func (us *UserSnapshot) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("snapshot_name", "user_snapshot")
	encoder.AddString("snapshot_date", us.Time.Format("20060102"))
	encoder.AddInt64("snapshot_time", us.Time.Unix())
	encoder.AddString("role_ctime", us.UserCreateTime.Format("2006-01-02 15:04:05"))
	encoder.AddString("last_login_time", us.LastLoginTime.Format("2006-01-02 15:04:05"))
	encoder.AddInt64("guild_id", us.GuildId)
	encoder.AddInt32("level", us.Level)
	encoder.AddInt32("exp", us.Exp)
	encoder.AddInt32("fight", us.Power)
	encoder.AddString("first_order_time", us.FirstPayTime.Format("2006-01-02 15:04:05"))
	encoder.AddString("first_order_id", us.FirstOrderId)
	encoder.AddInt32("first_order_level", us.FirstOrderLevel)
	encoder.AddInt32("balance_gold", us.PayDiamond)
	encoder.AddInt32("balance_silver", us.FreeDiamond)
	encoder.AddInt64("cumulative_payment", us.TotalPay)
	encoder.AddInt32("cumulative_order", us.TotalPayTimes)
	encoder.AddInt32("cumulative_days", us.TotalLoginDay)
	encoder.AddInt64("cumulative_onlinetime", us.TotalOnlineTime)

	return nil
}

type UserSnapshotCp struct {
	LatestExploreId int32
}

func (us *UserSnapshotCp) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddInt32("maintask_stop_id", us.LatestExploreId)

	return nil
}

func FormatGameLogId(appName, eventName string, timeStamp int64) string {
	buf := bytes.NewBufferString(appName)
	buf.WriteString(eventName)
	buf.WriteString(strconv.Itoa(int(timeStamp)))
	buf.WriteString(strconv.Itoa(int(rand.RangeInt32(1, math.MaxInt32))))

	return buf.String()
}

type GameOnlineUsers struct {
	Time        time.Time
	LogId       string
	GameBaseId  string
	GameId      int64
	Platform    string
	ZoneId      int64
	ChannelId   int64
	OnlineCount int32
}

func (gou *GameOnlineUsers) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	encoder.AddString("b_log_id", gou.LogId)
	encoder.AddInt64("b_utc_timestamp", gou.Time.Unix())
	encoder.AddString("b_datetime", gou.Time.Format("2006-01-02 15:04:05"))
	encoder.AddString("b_game_base_id", gou.GameBaseId)
	encoder.AddInt64("b_game_id", gou.GameId)
	encoder.AddString("b_platform", gou.Platform)
	encoder.AddInt64("b_zone_id", gou.ZoneId)
	encoder.AddInt64("b_channel_id", gou.ChannelId)
	encoder.AddInt32("online_num", gou.OnlineCount)

	return nil
}
