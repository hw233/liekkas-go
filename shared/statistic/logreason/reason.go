package logreason

var emptyReason = &Reason{}

type Reason struct {
	id          int32
	relateLog   string
	relateOrder string
}

func NewReason(id int32, opts ...ReasonOption) *Reason {
	reason := &Reason{id: id}
	for _, opt := range opts {
		opt.apply(reason)
	}

	return reason
}

func (r *Reason) Id() int32 {
	return r.id
}

func (r *Reason) RelateLog() string {
	return r.relateLog
}

func (r *Reason) RelateOrder() string {
	return r.relateOrder
}

type ReasonOption interface {
	apply(reason *Reason)
}

type reasonOptionFunc func(reason *Reason)

func (f reasonOptionFunc) apply(reason *Reason) {
	f(reason)
}

func AddRelateLog(logId string) ReasonOption {
	return reasonOptionFunc(func(reason *Reason) {
		reason.relateLog = logId
	})
}

func AddRelateOrder(orderId string) ReasonOption {
	return reasonOptionFunc(func(reason *Reason) {
		reason.relateOrder = orderId
	})
}

func EmptyReason() *Reason {
	return emptyReason
}

type reasonEnum struct {
	value int32
	desc  string
}

var reasons []reasonEnum

func addReasonEnum(value int32, desc string) int32 {
	enum := reasonEnum{
		value: value,
		desc:  desc,
	}

	reasons = append(reasons, enum)

	return value
}

// func CheckReasonDuplicate() error {
// 	reasonMap := map[int32]int32{}
// 	for _, reason := range reasons {
// 		value := reason.value
// 		_, ok := reasonMap[value]
// 		if ok {
// 			return errors.NewErrorf("duplicate reason value: %d", value)
// 		}
// 		reasonMap[value] = value
// 	}

// 	return nil
// }

var (
	GMAddReward              = addReasonEnum(1, "gm获取物品")
	GMPassGuide              = addReasonEnum(2, "gm通过引导")
	CreateAccount            = addReasonEnum(1001, "创建账号")
	LevelUp                  = addReasonEnum(1002, "玩家升级")
	Guide                    = addReasonEnum(1003, "新手引导")
	Manual                   = addReasonEnum(1004, "图鉴")
	DailyReward              = addReasonEnum(1005, "每日奖励")
	Questionnaire            = addReasonEnum(1006, "问卷调查")
	VisualNovel              = addReasonEnum(1007, "视觉小说")
	ChangeNickName           = addReasonEnum(1008, "改名")
	ActiviyEnd               = addReasonEnum(1009, "活动及结束回收")
	ItemUse                  = addReasonEnum(2001, "使用道具")
	PassLevel                = addReasonEnum(3001, "通关副本")
	LevelTarget              = addReasonEnum(3002, "副本目标")
	LevelAchievement         = addReasonEnum(3003, "副本成就")
	SweepLevel               = addReasonEnum(3004, "副本扫荡")
	ChapterReward            = addReasonEnum(3101, "章节奖励")
	QuestComplete            = addReasonEnum(4001, "完成任务")
	QuestActivityReward      = addReasonEnum(4002, "任务活跃度奖励")
	ExploreNPC               = addReasonEnum(5001, "主线探索npc")
	ExploreRewardPoint       = addReasonEnum(5002, "主线探索采集点")
	ExploreResource          = addReasonEnum(5003, "主线探索资源点")
	GachaOnce                = addReasonEnum(6001, "单抽")
	GachaTen                 = addReasonEnum(6002, "十连")
	GraveyardCreateCurtain   = addReasonEnum(7001, "模拟经营建造揭幕")
	GraveyardLevelUpCurtain  = addReasonEnum(7002, "模拟经营升级揭幕")
	GraveyardStageUpCurtain  = addReasonEnum(7003, "模拟经营升阶揭幕")
	GraveyardProduction      = addReasonEnum(7004, "模拟经营生产")
	GraveyardHelp            = addReasonEnum(7005, "模拟经营帮助")
	GraveyardPlot            = addReasonEnum(7006, "模拟经营对话")
	GraveyardBuild           = addReasonEnum(7007, "模拟经营建造")
	GraveyardBuildingLevelUp = addReasonEnum(7008, "模拟经营建筑升级")
	GraveyardBuildingStageUp = addReasonEnum(7009, "模拟经营建筑升阶")
	GraveyardProduce         = addReasonEnum(7010, "模拟经营建筑生产")
	GraveyardAccelerate      = addReasonEnum(7011, "模拟经营建筑加速生产")
	GraveyardUseBuff         = addReasonEnum(7012, "模拟经营使用buff")
	ReceiveMail              = addReasonEnum(8001, "领取邮件")
	GuildGreating            = addReasonEnum(9001, "公会问候")
	GuildTask                = addReasonEnum(9002, "公会任务奖励")
	GuildCreate              = addReasonEnum(9003, "创建公会")
	Shop                     = addReasonEnum(10001, "商店")
	QuickBuyStamina          = addReasonEnum(10002, "快速购买体力")
	RefreshShop              = addReasonEnum(10003, "刷新商店")
	ScorePassReward          = addReasonEnum(11001, "积分任务活动奖励")
	SignIn                   = addReasonEnum(12001, "签到")
	Yggdrasil                = addReasonEnum(13001, "世界探索")
	YggDailyDispatch         = addReasonEnum(13002, "世界探索日常派遣")
	YggGuildDispatch         = addReasonEnum(13003, "世界探索工会派遣")
	YggTaskComplete          = addReasonEnum(13004, "世界探索任务")
	YggSubTask               = addReasonEnum(13005, "世界探索子任务")
	YggMail                  = addReasonEnum(13006, "世界探索邮件")
	YggReturnToCity          = addReasonEnum(13007, "世界探索返回城市")
	YggHandleChest           = addReasonEnum(13008, "世界探索宝箱")
	YggLevelPass             = addReasonEnum(13009, "世界探索通关副本")
	YggProgress              = addReasonEnum(13010, "世界探索进度")
	YggBuild                 = addReasonEnum(13011, "世界探索建造")
	YggUsingBuilding         = addReasonEnum(13012, "世界探索使用建筑")
	YggBuildTranportPortal   = addReasonEnum(13013, "世界探索建造传送门")
	YggUseTranportPortal     = addReasonEnum(13014, "世界探索使用传送门")
	YggStartExplore          = addReasonEnum(13015, "世界探索开始探索")
	YggDeliverTaskItem       = addReasonEnum(13015, "世界探索交付任务道具")
	NewCharacter             = addReasonEnum(14001, "获得角色")
	CharacterLevelUp         = addReasonEnum(14002, "角色升级")
	CharacterStarUp          = addReasonEnum(14003, "角色升星")
	CharacterStageUp         = addReasonEnum(14004, "角色升阶")
	CharacterSkillLevelUp    = addReasonEnum(14005, "角色技能升级")
	HeroUnlock               = addReasonEnum(15001, "至尊解锁")
	HeroSkillLevelUp         = addReasonEnum(15002, "至尊升级")
	WorldItemAdvance         = addReasonEnum(16001, "世界道具升阶")
	WorldItemStrengthen      = addReasonEnum(16002, "世界道具强化")
	EquipmentAdvance         = addReasonEnum(17001, "装备升阶")
	EquipmentStrengthen      = addReasonEnum(17002, "装备强化")
	EquipmentRecastCamp      = addReasonEnum(17003, "装备重铸阵营")
)
