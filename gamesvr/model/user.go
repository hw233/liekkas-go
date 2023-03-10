package model

import (
	"context"
	"fmt"
	"gamesvr/manager"
	"math"
	"shared/common"
	"shared/csv/entry"
	"shared/csv/static"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/errors"
	"shared/utility/mysql"
	"shared/utility/number"
)

type User struct {
	ID                 int64             `db:"id" where:"=" major:"true"`
	Name               string            `db:"name"`
	Info               *UserInfo         `db:"info"`
	ItemPack           *ItemPack         `db:"item_pack"`
	CharacterPack      *CharacterPack    `db:"character_pack"`
	EquipmentPack      *EquipmentPack    `db:"equipment_pack"`
	WorldItemPack      *WorldItemPack    `db:"world_item_pack"`
	HeroPack           *HeroPack         `db:"hero_pack"`
	ManualInfo         *ManualInfo       `db:"manual_info"`
	QuestPack          *QuestPack        `db:"quest_pack"`
	QuestCache         *QuestCache       `db:"-"`
	Graveyard          *Graveyard        `db:"graveyard"`
	LevelsInfo         *LevelsInfo       `db:"levels_info"`
	ChapterInfo        *ChapterInfo      `db:"chapter_info"`
	ExploreInfo        *ExploreInfo      `db:"explore_info"`
	TowerInfo          *TowerInfo        `db:"tower_info"`
	StoreInfo          *StoreInfo        `db:"store_info"`
	Guild              *UserGuild        `db:"guild"`
	MailInfo           *MailInfo         `db:"mail_info"`
	Yggdrasil          *Yggdrasil        `db:"yggdrasil"`
	ActivityInfo       *ActivityInfo     `db:"activity_info"`
	ScorePassInfo      *ScorePassInfo    `db:"score_pass_info"`
	GachaRecords       *GachaRecords     `db:"gacha_record"`
	FormationInfo      *FormationInfo    `db:"formation_info"`
	RewardsResult      *RewardsResult    `db:"-"`
	Notifies           *UserNotifies     `db:"-"`
	Online             bool              `db:"-"`
	DailyRefreshers    []*DailyRefresher `db:"-"`
	OutGameInfo        *UserOutGameInfo  `db:"-"`
	Mercenary          *Mercenaries      `db:"mercenary"`
	*mysql.EmbedModule `db:"-"`
}

func NewUser(uid int64) *User {
	user := &User{
		ID:            uid,
		Name:          "",
		Info:          NewUserInfo(),
		ItemPack:      NewItemPack(),
		CharacterPack: NewCharacterPack(),
		WorldItemPack: NewWorldItemPack(),
		HeroPack:      NewHeroPack(),
		ManualInfo:    NewManualInfo(),
		QuestPack:     NewQuestPack(),
		QuestCache:    NewQuestCache(),
		Graveyard:     NewGraveyard(),
		LevelsInfo:    NewLevelsInfo(),
		ChapterInfo:   NewChapterInfo(),
		ExploreInfo:   NewExploreInfo(),
		TowerInfo:     NewTowerInfo(),
		StoreInfo:     NewStoreInfo(),
		MailInfo:      NewMailInfo(),
		Yggdrasil:     NewYggdrasil(),
		ActivityInfo:  NewActivityInfo(),
		ScorePassInfo: NewScorePassInfo(),
		GachaRecords:  NewGachaRecords(),
		Mercenary:     NewMercenaries(),
		FormationInfo: NewFormationInfo(),

		// ItemDroppedInfo: NewItemDroppedInfo(),
		Guild:       NewUserGuild(),
		Notifies:    NewUserNotifies(),
		Online:      false,
		OutGameInfo: NewUserOutGameInfo(),
		EmbedModule: &mysql.EmbedModule{},
	}

	user.EquipmentPack = NewEquipmentPack()
	user.RewardsResult = NewRewardsResult(user)

	return user
}

func (u *User) GetUserId() int64 {
	return u.ID
}

func (u *User) GetUserName() string {
	return u.Name
}

// 建号时初始设置
func (u *User) InitForCreate(ctx context.Context) error {
	suffix, err := manager.Global.GenNickNameSuffix(ctx)
	if err != nil {
		return errors.WrapTrace(err)
	}
	// 初始昵称
	err = u.changeNickname(ctx, fmt.Sprintf("飞鼠_%d", suffix))
	if err != nil {
		return errors.WrapTrace(err)
	}

	reason := logreason.NewReason(logreason.CreateAccount)

	// 初始等级
	u.Info.Level.SetValue(1, reason)
	levelConfig, ok := manager.CSV.TeamLevelCache.GetByLv(u.Info.Level.Value())
	if !ok {
		return errors.WrapTrace(common.ErrNotFoundInCSV)
	}
	// 初始精力
	u.Info.Energy.SetValue(levelConfig.MaxStamina)
	// 初始ap
	u.Info.Ap.SetValue(levelConfig.MaxAp)

	// 发建号资源
	_, err = u.AddRewardsByDropId(manager.CSV.GlobalEntry.InitUserDrop, reason)
	if err != nil {
		return errors.WrapTrace(err)
	}

	heroId := manager.CSV.GlobalEntry.InitUserHero
	u.HeroPack.AddHero(heroId)
	hero, ok := u.HeroPack.GetHero(heroId)
	if !ok {
		return errors.Swrapf(common.ErrHeroNotFound, heroId)
	}

	power, err := u.CalHeroCombatPower(hero)
	if err != nil {
		return errors.WrapTrace(err)
	}
	hero.Power = power

	// 初始化模拟经营结构
	u.Graveyard.initForCreate()

	// 更新名片中角色展示的数据
	u.RefreshCardShowCharacterCache()

	//
	u.RewardsResult.Clear()

	u.BICreateUser()

	return nil
}

// Init 登陆时初始化
func (u *User) Init(ctx context.Context) {
	u.RegisterDailyRefreshers()

	// 精力回复上限和速度
	u.Info.Energy.RegisterInterval(manager.CSV.TeamLevelCache.GetStaminaRecoverSeconds())
	// ap恢复时间间隔
	u.Info.Ap.RegisterInterval(manager.CSV.Yggdrasil.GetYggApRestoreTime())

	levelConfig, ok := manager.CSV.TeamLevelCache.GetByLv(u.Info.Level.Value())
	if ok {
		u.Info.Energy.SetTimerUpper(levelConfig.MaxStamina)
		u.Info.Ap.SetTimerUpper(levelConfig.MaxAp)

	}

	u.initItemPack()

	u.RegisterNumberEvent()

	u.Yggdrasil.init(ctx)

	u.InitQuest()
}

func (u *User) RegisterNumberEvent() {
	u.Info.Exp.RegisterChangedEvent(u.expChangeEvent)
	u.Info.Level.RegisterChangedEvent(u.levelChangeEvent)

	u.Info.Gold.RegisterChangedEvent(u.genResourceEvent(static.CommonResourceTypeMoney))
	u.Info.DiamondGift.RegisterChangedEvent(u.genResourceEvent(static.CommonResourceTypeDiamondGift))
	u.Info.DiamondCash.RegisterChangedEvent(u.genResourceEvent(static.CommonResourceTypeDiamondCash))
	u.Guild.GuildGold.RegisterChangedEvent(u.genResourceEvent(static.CommonResourceTypeGuildMoney))
}

// AddRewardsByDropId 通过dropId掉落
func (u *User) AddRewardsByDropId(dropId int32, reason *logreason.Reason) (*common.Rewards, error) {
	rewards, err := manager.CSV.Drop.DropRewards(dropId)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return u.addRewards(rewards, reason)
}

func (u *User) AddRewardsByDropIds(dropIds []int32, reason *logreason.Reason) (*common.Rewards, error) {
	realRewards := common.NewRewards()
	for _, dropId := range dropIds {
		rr, err := u.AddRewardsByDropId(dropId, reason)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}

		realRewards.Append(rr)
	}

	return realRewards, nil
}

func (u *User) addRewardsForMerge(rewards *common.Rewards, merge bool, reason *logreason.Reason) (*common.Rewards, error) {
	// todo: 写一个邮件列表，然后发奖之后触发一次扫描列表，写在这里耦合度过高
	// var fullToMail []*common.Reward

	value := rewards.Value()
	if merge {
		value = rewards.MergeValue()
	}

	realRewards := common.NewRewards()

	for i := range value {
		reward := &value[i]

		rr, addToRewardsResult, err := u.addReward(reward, reason)
		if err != nil {
			return nil, err
		}
		if addToRewardsResult {
			u.RewardsResult.AddReward(reward)
		}

		realRewards.Append(rr)
	}

	return realRewards, nil
}

// addRewards
func (u *User) AddRewards(rewards *common.Rewards, reason *logreason.Reason) (*common.Rewards, error) {
	return u.addRewardsForMerge(rewards, true, reason)
}

func (u *User) addRewards(rewards *common.Rewards, reason *logreason.Reason) (*common.Rewards, error) {
	return u.addRewardsForMerge(rewards, true, reason)
}

func (u *User) addReward(reward *common.Reward, reason *logreason.Reason) (*common.Rewards, bool, error) {
	realRewards := common.NewRewards()
	realRewards.AddReward(common.NewReward(reward.ID, reward.Num))

	switch reward.Type {
	case static.ItemTypeCurrency: // 一般道具（金币，钻石，体力等）
		switch reward.ID {
		case static.CommonResourceTypeTeamExp: // 账号经验
			u.AddExp(reward.Num, reason)
		case static.CommonResourceTypeMoney: // 金币
			u.Info.Gold.Plus(reward.Num, reason)
		case static.CommonResourceTypeDiamondGift: // 免费钻石
			u.Info.DiamondGift.Plus(reward.Num, reason)
		case static.CommonResourceTypeEnergy: // 体力
			u.Info.Energy.Plus(reward.Num)
		case static.CommonResourceTypeDiamondCash: // 付费钻石
			u.Info.DiamondCash.Plus(reward.Num, reason)
		case static.CommonResourceTypeAp: // 世界探索精力
			u.Info.Ap.Plus(reward.Num)
		case static.CommonResourceTypeHeroExp:
			before := u.HeroPack.Exp
			u.HeroPack.AddExp(reward.Num)
			after := u.HeroPack.Exp
			u.BIResourceChange(reward.ID, int64(after-before), int64(after), reason)
		case static.CommonResourceTypeArenaMoney: // 竞技场代币
		case static.CommonResourceTypeGuildMoney: // 公会代币
			u.Guild.GuildGold.Plus(reward.Num, reason)
		}

	case static.ItemTypeGacha, // 抽卡道具
		static.ItemTypeCharExp,              // 角色经验道具
		static.ItemTypeCharSkill,            // 角色进阶/技能道具
		static.ItemTypeCharPiece,            // 角色碎片
		static.ItemTypeWorldItemExp,         // 神器经验道具
		static.ItemTypeWorldItemBreak,       // 神器突破道具
		static.ItemTypeEquipExp,             // 装备经验道具
		static.ItemTypeEquipBreak,           // 装备突破道具
		static.ItemTypeGraveyardAccelerate,  // 模拟经营加速道具
		static.ItemTypeGraveyardGetProduct,  // 模拟经营使用道具直接获得产出
		static.ItemTypeGraveyardProduceBuff, // 模拟经营生产buff道具
		static.ItemTypeEnergyItem,           // 体力药水
		static.ItemTypeCommon,               // 一般道具
		static.ItemTypeGiftSelectOne,        // N选1礼包
		static.ItemTypeGiftRandomDrop,       // 随机普通礼包
		static.ItemTypeCharacterAvatar,      // 头像
		static.ItemTypeCharacterFrame,       // 头像框
		static.ItemTypeHeroUnlockItem,       //至尊道具
		static.ItemTypeYggPrestige:          // 声望
		u.addItem(reward.ID, reward.Num, reason)
	case static.ItemTypeCharacterAvatarAuto: // 角色Q版头像
		u.addItem(reward.ID, reward.Num, reason)
		// 角色Q版头像不弹奖励弹框
		return realRewards, false, nil

	case static.ItemTypeAutoRandom: // 获得后不进背包，自动使用 且 使用后执行一般掉落逻辑的道具-即按照原掉落规则纯随机
		dropIDs, err := manager.CSV.Reward.AutoRandomParam(reward.ID)
		if err != nil {
			return nil, false, errors.WrapTrace(err)
		}
		for i := 0; i < int(reward.Num); i++ {
			index, err := u.Info.Drop.RandIndex(reward.ID, len(dropIDs))
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
			rewards, err := manager.CSV.Drop.DropRewards(dropIDs[index])
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
			if rewards.ContainsType(static.ItemTypeAutoRandom) {
				return nil, false, errors.WrapTrace(common.ErrRecursionReward)
			}
			realRewards, err = u.addRewards(rewards, reason)
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
		}
		return realRewards, false, nil

	case static.ItemTypeAutoRandomRemove: // 获得后不进背包，自动使用 且 使用后执行奖池掉落逻辑的道具
		dropIDs, err := manager.CSV.Reward.AutoRandomRemoveParam(reward.ID)
		if err != nil {
			return nil, false, errors.WrapTrace(err)
		}

		realRewards = common.NewRewards()

		for i := 0; i < int(reward.Num); i++ {
			index, err := u.Info.Drop.RandIndex(reward.ID, len(dropIDs))
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
			rewards, err := manager.CSV.Drop.DropRewards(dropIDs[index])
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
			u.Info.Drop.MarkIndex(reward.ID, index)

			if rewards.ContainsType(static.ItemTypeAutoRandomRemove) {
				return nil, false, errors.WrapTrace(common.ErrRecursionReward)
			}
			rr, err := u.addRewards(rewards, reason)
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
			realRewards.Append(rr)
		}

		return realRewards, false, nil

	case static.ItemTypeYggItemSpecial: // 世界探索特殊道具（魔石，不占用世界探索格子，不能存到仓库）
	case static.ItemTypeCharacterAdaptiveGift: // 根据新手池子抽出的角色自适应的道具礼包
		adaptiveGift, ok := manager.CSV.Item.GetCharacterAdaptiveGift(reward.ID)
		if !ok {
			return nil, false, errors.WrapTrace(errors.Swrapf(common.ErrNotFoundInCSV, entry.CfgItemData))
		}
		dropId, ok := adaptiveGift.DropMap[u.GachaRecords.DestinyChild]
		if !ok {
			// 还没抽新手池子
			return realRewards, false, nil
		}
		rewards, err := manager.CSV.Drop.DropRewards(dropId)
		if err != nil {
			return nil, false, errors.WrapTrace(err)
		}

		realRewards, err = u.addRewards(rewards, reason)
		if err != nil {
			return nil, false, errors.WrapTrace(err)
		}
		return realRewards, false, nil

	case static.ItemTypeEquipment: // 装备
		for i := int32(0); i < reward.Num; i++ {
			if u.EquipmentPack.Count() >= manager.CSV.GlobalEntry.EquipmentCountLimit {
				break
			}

			e, err := u.EquipmentPack.NewEquipment(reward.ID)
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}

			u.TriggerQuestUpdate(static.TaskTypeEquipmentLevelCount, e.Rarity, int32(0), e.Level.Value())

			u.RewardsResult.AddEquipment(e)
			u.BIEquipmentOp(e, bilog.EquipmentOpGet, reason)
		}
	case static.ItemTypeWorldItem: // 世界级道具
		for i := int32(0); i < reward.Num; i++ {
			if u.WorldItemPack.Count() >= manager.CSV.GlobalEntry.WorlItemCountLimit {
				break
			}

			w, err := u.WorldItemPack.NewWorldItem(reward.ID)
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
			u.ActiveManual(static.ManualTypeWorldItem, w.WID)
			u.RewardsResult.AddWorldItem(w)

			u.BIWorldItemOp(w, bilog.WorldItemOpGet, reason)

			if int32(w.Rarity) == static.RaritySsr {
				u.AddGreetings(reward.ID, 2, 1)
			}
		}
	case static.ItemTypeCharacter: // 角色
		repeatNum, err := u.addCharacter(reward.ID, reward.Num, u, reason)
		if err != nil {
			return nil, false, errors.WrapTrace(err)
		}
		if repeatNum > 0 {
			charaShard, err := manager.CSV.Gacha.GetCharaShards(reward.ID)
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
			rewards := common.NewRewards()
			rewards.AddReward(charaShard.CharaShard)

			realRewards, err = u.addRewards(rewards.Multiple(repeatNum), reason)
			if err != nil {
				return nil, false, errors.WrapTrace(err)
			}
		}
		if repeatNum == reward.Num {
			return realRewards, false, nil

		}

		rare, err := manager.CSV.Character.GetRare(reward.ID)
		if err != nil {
			return nil, false, errors.WrapTrace(err)
		}
		if rare == static.RaritySsr {
			u.AddGreetings(reward.ID, 1, 1)
		}

	}

	return realRewards, true, nil
}

// 角色
func (u *User) addCharacter(id int32, num int32, user *User, reason *logreason.Reason) (int32, error) {
	if num <= 0 {
		return 0, errors.WrapTrace(common.ErrParamError)
	}
	if !u.CharacterPack.Contains(id) {
		// 新角色
		rare, err := manager.CSV.Character.GetRare(id)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}

		avatar, err := manager.CSV.Character.Avatar(id)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}
		character := u.CharacterPack.NewCharacter(id)
		// 附带q版头像
		//todo:策划暂时没空配
		if avatar > 0 {
			reason := logreason.NewReason(logreason.NewCharacter)
			u.addItem(avatar, 1, reason)
		}

		// 初始化角色技能
		user.ActiveCharacterSkill(id)
		character.SetRare(rare)

		u.TriggerQuestUpdate(static.TaskTypeCharacterLevelCount, character, int32(0))

		camp, err := manager.CSV.Character.Camp(id)
		if err == nil {
			u.TriggerQuestUpdate(static.TaskTypeCharacterCampCount, id, camp)
		}

		u.TriggerQuestUpdate(static.TaskTypeHasCharacters, id)

		// 激活图鉴
		u.ActiveManual(static.ManualTypeCharacter, id)

		u.RewardsResult.AddCharacter(character)

		power, err := u.CalCharacterCombatPower(character)
		if err != nil {
			return 0, errors.WrapTrace(err)
		}
		character.Power = power

		u.BICharaGet(character, reason)

		return num - 1, nil
	} else {
		return num, nil
	}

}

func (u *User) initItemPack() {
	for id, item := range *u.ItemPack {
		limit := manager.CSV.Item.GetLimit(id)
		item.SetLimit(0, limit)
	}
}

func (u *User) addItem(id, num int32, reason *logreason.Reason) {
	newItem := false
	if u.ItemPack.Count(id) <= 0 {
		newItem = true
	}
	before, after := u.ItemPack.Add(id, num)

	if newItem {
		limit := manager.CSV.Item.GetLimit(id)
		u.ItemPack.SetLimit(id, limit)
		after = u.ItemPack.Count(id)
	}

	u.RewardsResult.AddItem(id, after)
	u.BIItemChange(id, num, before, after, reason)
}
func (u *User) minusItem(id, num int32, reason *logreason.Reason) {
	before, after := u.ItemPack.Minus(id, num)
	u.RewardsResult.AddItem(id, after)
	u.BIItemChange(id, -num, before, after, reason)
}

func (u *User) AddExp(value int32, reason *logreason.Reason) {
	calcuator := &UserExpLevelCalculator{
		UserInfo: u.Info,
		Reason:   reason,
	}

	u.AddExpCommon(value, calcuator, manager.CSV.TeamLevelCache)
}

func (u *User) CheckRewardsEnough(rewards *common.Rewards) error {
	for _, reward := range rewards.MergeValue() {
		// 不需要数量的直接返回
		if reward.Num == 0 {
			return nil
		}

		switch reward.Type {
		case static.ItemTypeCurrency:
			switch reward.ID {
			case static.CommonResourceTypeMoney: // 金币
				if !u.Info.Gold.Enough(reward.Num) {
					return errors.Swrapf(common.ErrItemNotEnough, reward.ID, reward.Num)
				}
			case static.CommonResourceTypeDiamondGift: // 免费钻石
				sum := number.NewEventNumber(u.Info.DiamondGift.Value()+u.Info.DiamondCash.Value(), 0, math.MaxInt32)
				if !sum.Enough(reward.Num) {
					return errors.Swrapf(common.ErrItemNotEnough, reward.ID, reward.Num)
				}
			case static.CommonResourceTypeEnergy: // 体力
				if !u.Info.Energy.Enough(reward.Num) {
					return errors.Swrapf(common.ErrItemNotEnough, reward.ID, reward.Num)
				}
			case static.CommonResourceTypeDiamondCash: // 付费钻石
				if !u.Info.DiamondCash.Enough(reward.Num) {
					return errors.Swrapf(common.ErrItemNotEnough, reward.ID, reward.Num)
				}
			case static.CommonResourceTypeAp: // 世界探索精力
				if !u.Info.Ap.Enough(reward.Num) {
					return errors.Swrapf(common.ErrItemNotEnough, reward.ID, reward.Num)
				}
			case static.CommonResourceTypeArenaMoney: // 竞技场代币
			case static.CommonResourceTypeGuildMoney: // 公会代币
				u.Guild.GuildGold.Enough(reward.Num)
			default:
				return errors.WrapTrace(common.ErrItemTypeCannotConsume)
			}
		case static.ItemTypeGacha, // 抽卡道具
			static.ItemTypeCharExp,              // 角色经验道具
			static.ItemTypeCharSkill,            // 角色进阶/技能道具
			static.ItemTypeCharPiece,            // 角色碎片
			static.ItemTypeWorldItemExp,         // 神器经验道具
			static.ItemTypeWorldItemBreak,       // 神器突破道具
			static.ItemTypeEquipExp,             // 装备经验道具
			static.ItemTypeEquipBreak,           // 装备突破道具
			static.ItemTypeGraveyardAccelerate,  // 模拟经营加速道具
			static.ItemTypeGraveyardGetProduct,  // 模拟经营使用道具直接获得产出
			static.ItemTypeGraveyardProduceBuff, // 模拟经营生产buff道具
			static.ItemTypeEnergyItem,           // 体力药水
			static.ItemTypeCommon,               // 一般道具
			static.ItemTypeGiftSelectOne,        // N选1礼包
			static.ItemTypeGiftRandomDrop,       // 随机普通礼包
			static.ItemTypeHeroUnlockItem,       //至尊解锁道具
			static.ItemTypeYggPrestige:          // 声望
			if !u.ItemPack.Enough(reward.ID, reward.Num) {
				return errors.Swrapf(common.ErrItemNotEnough, reward.ID, reward.Num)
			}
		default:
			return errors.WrapTrace(common.ErrItemTypeCannotConsume)
		}
	}

	return nil
}

func (u *User) CostRewards(rewards *common.Rewards, reason *logreason.Reason) error {
	mergeRewards := rewards.MergeValue()

	for _, reward := range mergeRewards {
		// 不需要数量的直接返回
		if reward.Num == 0 {
			return nil
		}

		switch reward.Type {
		case static.ItemTypeCurrency:
			switch reward.ID {
			case static.CommonResourceTypeMoney: // 金币
				u.Info.Gold.Minus(reward.Num, reason)
			case static.CommonResourceTypeDiamondGift: // 免费钻石
				u.Info.DiamondGift.Minus(reward.Num, reason)
				if !u.Info.DiamondGift.Enough(0) {
					u.Info.DiamondCash.Plus(u.Info.DiamondGift.Value(), reason)
					u.Info.DiamondGift.SetValue(0, reason)
				}
			case static.CommonResourceTypeEnergy: // 体力
				u.Info.Energy.Minus(reward.Num)
				u.AddExp(reward.Num*manager.CSV.GlobalEntry.StaminaExpRatio, reason)
				u.Guild.AddTaskItem(static.GuildTaskStamina, reward.Num)
			case static.CommonResourceTypeDiamondCash: // 付费钻石
				u.Info.DiamondCash.Minus(reward.Num, reason)
			case static.CommonResourceTypeAp: // 世界探索精力
				u.Info.Ap.Minus(reward.Num)
			case static.CommonResourceTypeArenaMoney: // 竞技场代币
			case static.CommonResourceTypeGuildMoney: // 公会代币
				u.Guild.GuildGold.Minus(reward.Num, reason)
			default:
				return errors.WrapTrace(common.ErrItemTypeCannotConsume)
			}

		case static.ItemTypeGacha, // 抽卡道具
			static.ItemTypeCharExp,              // 角色经验道具
			static.ItemTypeCharSkill,            // 角色进阶/技能道具
			static.ItemTypeCharPiece,            // 角色碎片
			static.ItemTypeWorldItemExp,         // 神器经验道具
			static.ItemTypeWorldItemBreak,       // 神器突破道具
			static.ItemTypeEquipExp,             // 装备经验道具
			static.ItemTypeEquipBreak,           // 装备突破道具
			static.ItemTypeGraveyardAccelerate,  // 模拟经营加速道具
			static.ItemTypeGraveyardGetProduct,  // 模拟经营使用道具直接获得产出
			static.ItemTypeGraveyardProduceBuff, // 模拟经营生产buff道具
			static.ItemTypeEnergyItem,           // 体力药水
			static.ItemTypeCommon,               // 一般道具
			static.ItemTypeGiftSelectOne,        // N选1礼包
			static.ItemTypeGiftRandomDrop:       // 随机普通礼包
			u.minusItem(reward.ID, reward.Num, reason)
		default:
			return errors.WrapTrace(common.ErrItemTypeCannotConsume)
		}

		u.TriggerQuestUpdate(static.TaskTypeCostItem, reward.ID, reward.Num)
	}

	return nil
}

func (u *User) IsOnline() bool {
	return u.Online
}

func (u *User) OnOnline(now int64) {
	u.TryDailyRefresh(now)
	u.TryAcceptNewQuests()
	u.RemoveExpireMails()
	u.ActivityPreSetting()
	u.MailOnLine()
	u.Info.LastLoginTime = now

	u.Online = true

	u.BILogin()
}

func (u *User) OnOffline() {
	u.Online = false
	u.RefreshCardShowCharacterCache()
	u.Info.OnOffline()
	u.BILogout()
	u.BISnapshot()
}

func (u *User) OnSecond(now int64) {
	u.TryDailyRefresh(now)
}

func (u *User) On30Second(now int64) {
	u.UpdateActivities(now)
}

func (u *User) OnHour() {
	u.RefreshCardShowCharacterCache()
}

func (u *User) DailyRefresh(refreshTime int64) {
	u.CheckScorePassUpdate(refreshTime)
	u.DailyGiftDailyRefresh(refreshTime)
	u.Info.RecordLoginDay()
}
