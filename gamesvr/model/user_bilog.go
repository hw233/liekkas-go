package model

import (
	"gamesvr/manager"
	"shared/common"
	"shared/csv/static"
	"shared/protobuf/pb"
	"shared/statistic/bilog"
	"shared/statistic/logreason"
	"shared/utility/servertime"
	"strings"
	"time"
)

type biLogData struct {
	common     *bilog.UserCommon
	data       bilog.LogObjMarshaler
	properties bilog.LogObjMarshaler
	cpParams   bilog.LogObjMarshaler
}

func (bld *biLogData) Serialize() ([]bilog.LogObjMarshaler, []*bilog.EventData) {
	datas := []bilog.LogObjMarshaler{}

	if bld.common != nil {
		datas = append(datas, bld.common)
	}

	if bld.data != nil {
		datas = append(datas, bld.data)
	}

	eventDatas := []*bilog.EventData{}
	if bld.properties != nil {
		eventDatas = append(eventDatas, bilog.NewEventData("properties", bld.properties))
	}

	if bld.cpParams != nil {
		eventDatas = append(eventDatas, bilog.NewEventData("cp_params", bld.cpParams))
	}

	return datas, eventDatas
}

func (u *User) WriteBILog(logData *biLogData) {
	bilog.EventLog(logData.Serialize())
}

func (u *User) WriteSnapshot(logData *biLogData) {
	bilog.SnapshotLog(logData.Serialize())
}

func (u *User) CreateBICommon(eventName string) *bilog.UserCommon {
	data := &bilog.UserCommon{
		Uid:       u.GetUserId(),
		EventName: eventName,
		LogTime:   servertime.Now(),
	}

	data.LogId = bilog.FormatLogId("overlord", eventName, data.Uid, data.LogTime.Unix())

	return data
}

func (u *User) BICreateUser() {
	logData := &biLogData{
		common: u.CreateBICommon(bilog.EvenNameCreate),
		properties: &bilog.UserCreate{
			RoleName:  u.GetUserName(),
			Model:     u.OutGameInfo.GetDevice(),
			OSVersion: u.OutGameInfo.GetOsVersion(),
			Network:   u.OutGameInfo.GetNetwork(),
			Mac:       u.OutGameInfo.GetMac(),
			Ip:        u.OutGameInfo.GetIp(),
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BILogin() {
	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameLogin),
		properties: &bilog.UserLoginLog{
			RoleName:      u.GetUserName(),
			Model:         u.OutGameInfo.GetDevice(),
			OSVersion:     u.OutGameInfo.GetOsVersion(),
			Network:       u.OutGameInfo.GetNetwork(),
			Mac:           u.OutGameInfo.GetMac(),
			Ip:            u.OutGameInfo.GetIp(),
			Udid:          u.OutGameInfo.GetUdid(),
			SDKUdid:       u.OutGameInfo.GetSdkUdid(),
			SDKUid:        u.OutGameInfo.GetUdid(),
			GameBaseId:    "",
			GameId:        1,
			Platform:      "",
			ZoneId:        0,
			ChannelId:     1,
			ClientVersion: "1.1.1",
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BILogout() {
	now := servertime.Now()
	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameLogout),

		properties: &bilog.UserLogoutLog{
			OnlineTime: now.Unix() - u.Info.LastLoginTime,
			Gender:     1,
			Model:      u.OutGameInfo.GetDevice(),
			OSVersion:  u.OutGameInfo.GetOsVersion(),
			Network:    u.OutGameInfo.GetNetwork(),
			Mac:        u.OutGameInfo.GetMac(),
			Ip:         u.OutGameInfo.GetIp(),
		},
	}
	u.WriteBILog(logData)
}

func (u *User) BIGuide(guideId int32, commonData *bilog.UserCommon) {
	logData := &biLogData{
		common: commonData,
		properties: &bilog.UserGuide{
			GuideId: guideId,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIExpLevel(opType string, level, exp int32, reason *logreason.Reason) {
	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameExp),
		properties: &bilog.UserExpLevel{
			Exp:    exp,
			Level:  level,
			Reason: reason,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIResourceChange(resourceId int32, change, after int64,
	reason *logreason.Reason) {

	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameResourceChange),
		properties: &bilog.UserResourceChange{
			ResourceId: resourceId,
			ChangeNum:  change,
			BeforeNum:  after - change,
			AfterNum:   after,
			Reason:     *reason,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIItemChange(itemId int32, change, before, after int32,
	reason *logreason.Reason) {

	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameItemChange),
		properties: &bilog.UserItemChange{
			ItemId:    itemId,
			ChangeNum: change,
			BeforeNum: before,
			AfterNum:  after,
			Reason:    *reason,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIGraveyard(opType int, mainTowerLevel int32, building *common.UserGraveyardBuild,
	duration int64, commonData *bilog.UserCommon) {

	var charas []int32
	if building.UserGraveyardProduce != nil && building.Characters != nil {
		charas = building.Characters.GetCharacters()
	} else {
		charas = []int32{}
	}

	logData := &biLogData{
		common: commonData,
		cpParams: &bilog.UserGraveyard{
			OpType:         opType,
			MainTowerLevel: mainTowerLevel,
			BuildingId:     building.BuildId,
			BuildingLevel:  building.Lv,
			BuildingStage:  building.Stage,
			Charas:         charas,
			Duration:       duration,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIQuest(quest *Quest, commonData *bilog.UserCommon) {
	questCfg, _ := manager.CSV.Quest.Get(quest.Id)

	logData := &biLogData{
		common: commonData,
		data: &bilog.UserQuestData{
			QuestId:    quest.Id,
			QuestName:  questCfg.Cname,
			QuestType:  questCfg.Module,
			QuestState: int32(quest.State),
			UserExp:    u.Info.Exp.Value(),
			UserPower:  u.GetUserPower(),
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIGacha(poolId, poolType int32, isTen bool, cost, rewards *common.Rewards,
	totalTimes int32, commonData *bilog.UserCommon) {

	if cost == nil {
		cost = common.NewRewards()
	}

	if rewards == nil {
		rewards = common.NewRewards()
	}

	logData := &biLogData{
		common: commonData,
		data: &bilog.UserGachaData{
			PoolId:     poolId,
			IsTen:      isTen,
			Cost:       (*bilog.Rewards)(cost),
			Reward:     (*bilog.Rewards)(rewards),
			TotalTimes: totalTimes,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BILevels(levelId, systemType int32, isWin, isFirst bool, score int32, formation *pb.VOBattleFormation,
	cost, rewards *common.Rewards, battleData *BattleVerifyData, commonData *bilog.UserCommon) {

	logData := &biLogData{
		common: commonData,
		properties: &bilog.UserLevelsProperties{
			LevelId:   levelId,
			LevelType: systemType,
			IsWin:     isWin,
			Score:     score,
			IsFirst:   isFirst,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BICharaGet(chara *Character, reason *logreason.Reason) {
	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameCharaGet),
		properties: &bilog.UserGetCharacter{
			CharaId:    chara.GetId(),
			CharaLevel: chara.GetLevel(),
			CharaStar:  chara.GetStar(),
			CharaStage: chara.GetStage(),
			CharaExp:   chara.GetExp(),
			Reason:     reason,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BICharaOp(chara *Character, opType int8) {
	skills := map[int32]int32{}

	for skillNum, level := range chara.Skills {
		skills[skillNum] = level
	}

	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameCharaOp),
		properties: &bilog.UserCharaOperationP{
			OpType:  opType,
			CharaId: chara.GetId(),
			Level:   chara.GetLevel(),
			Star:    chara.GetStar(),
			Stage:   chara.GetStage(),
			Exp:     chara.GetExp(),
		},
		cpParams: &bilog.UserCharaOperationCP{
			SkillLevels: skills,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIHeroOp(hero *Hero, heroLevel int32, opType int8) {
	skills := map[int32]int32{}
	for skillId, level := range hero.Skills {
		skills[skillId] = level
	}

	charas := make([]int32, 0, len(hero.Attendants))
	for _, attendant := range hero.Attendants {
		charas = append(charas, attendant.CharaId)
	}

	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameHeroOp),

		cpParams: &bilog.UserHeroCp{
			OpType:     opType,
			HeroId:     hero.ID,
			HeroLevel:  heroLevel,
			HeroSkills: skills,
			Charas:     charas,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIEquipmentOp(equipment *common.Equipment, opType int8, reason *logreason.Reason) {
	attrs := map[int8]int32{}
	for _, attr := range equipment.Attrs {
		if attr.Unlock {
			attrs[attr.Attr] = attr.Value
		}
	}

	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameEquipmentOp),
		properties: &bilog.UserEquipmentP{
			OpType:  opType,
			Id:      equipment.ID,
			Eid:     equipment.EID,
			CharaId: equipment.CID,
		},
		cpParams: &bilog.UserEquipmentCP{
			Stage:  equipment.Stage.Value(),
			Attrs:  attrs,
			Reason: reason,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIWorldItemOp(worldItem *common.WorldItem, opType int8, reason *logreason.Reason) {

	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameWorldItemOp),
		cpParams: &bilog.UserWorldItemCP{
			OpType:  opType,
			Id:      worldItem.ID,
			Wid:     worldItem.WID,
			CharaId: worldItem.CID,
			Stage:   worldItem.Stage.Value(),
			Reason:  reason,
		},
	}

	u.WriteBILog(logData)
}

func (u *User) BIMail(mail *Mail, senderId int64, opType int8) {
	var title, content string = "", ""
	if opType == bilog.MailOpSend || opType == bilog.MailOpReceive {
		if len(mail.TitleArgs) > 0 {
			title = strings.Join(mail.TitleArgs, ",")
		} else {
			title = mail.Title
		}

		if len(mail.ContentArgs) > 0 {
			content = strings.Join(mail.ContentArgs, ",")
		} else {
			content = mail.Content
		}
	}

	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameMail),
		properties: &bilog.UserMail{
			OpType:               opType,
			TemplateId:           mail.TemplateId,
			SenderId:             senderId,
			MailId:               mail.Id,
			Title:                title,
			Content:              content,
			SendTime:             mail.SendTime,
			ExpireTime:           mail.ExpireTime,
			IsAttachmentReceived: mail.State == static.MailStateReceived,
			Attachment:           (*bilog.Rewards)(mail.Attachment),
		},
	}

	u.WriteBILog(logData)
}

//----------------------------------------
//snapshot
//----------------------------------------
func (u *User) BISnapshot() {
	logData := &biLogData{
		common: u.CreateBICommon(bilog.EventNameSnapshot),
		data: &bilog.UserSnapshot{
			Time:            time.Now(),
			UserCreateTime:  time.Unix(u.Info.RegisterAt, 0),
			LastLoginTime:   time.Unix(u.Info.LastLoginTime, 0),
			GuildId:         u.Guild.GuildID,
			Level:           u.Info.GetLevel(),
			Exp:             u.Info.GetExp(),
			Power:           u.GetUserPower(),
			FirstPayTime:    time.Unix(0, 0),
			FirstOrderId:    "",
			FirstOrderLevel: 0,
			PayDiamond:      u.Info.DiamondCash.Value(),
			FreeDiamond:     u.Info.DiamondGift.Value(),
			TotalPay:        0,
			TotalPayTimes:   0,
			TotalLoginDay:   u.Info.TotalLoginDay,
			TotalOnlineTime: u.Info.TotalLoginTime,
		},
		cpParams: &bilog.UserSnapshotCp{
			LatestExploreId: u.Info.LatestExploreLevel,
		},
	}

	u.WriteSnapshot(logData)
}
