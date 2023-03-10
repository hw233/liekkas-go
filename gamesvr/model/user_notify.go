package model

import (
	"shared/protobuf/pb"
	"shared/utility/servertime"

	"google.golang.org/protobuf/proto"
)

const (
	MailCountPerPage = 20
)

type UserNotifies struct {
	DailyRefreshNotify      *pb.S2CDailyRefreshNotify
	QuestUpdateNotify       *pb.S2CQuestUpdateNotify
	ExploreEventNotify      *pb.S2CExploreEventNotify
	ExploreResourceNotify   *pb.S2CExploreResourceNotify
	MailNotifies            []*pb.S2CMailNotify
	TowerNotify             *pb.S2CTowerUpdateNotify
	LevelNotify             *pb.S2CLevelNotify
	YggPush                 []proto.Message
	GraveyardPush           []proto.Message
	ScorePassNotify         *pb.S2CScorePassNotify
	StoreDailyRefreshNotify *pb.S2CStoreDailyRefresh
}

func NewUserNotifies() *UserNotifies {
	return &UserNotifies{}
}

//----------------------------------------
//DailyRefresh
//----------------------------------------
func (u *User) AddDailyRefreshNotify() {
	if !u.IsOnline() {
		return
	}

	if u.Notifies.DailyRefreshNotify == nil {
		u.Notifies.DailyRefreshNotify = &pb.S2CDailyRefreshNotify{}
	}
}

func (u *User) PopDailyRefreshNotify() *pb.S2CDailyRefreshNotify {
	notify := u.Notifies.DailyRefreshNotify
	u.Notifies.DailyRefreshNotify = nil
	return notify
}

//----------------------------------------
//Quest
//----------------------------------------
func (u *User) AddQuestUpdateNotify(questIds []int32) {
	if !u.IsOnline() {
		return
	}

	if u.Notifies.QuestUpdateNotify == nil {
		u.Notifies.QuestUpdateNotify = &pb.S2CQuestUpdateNotify{
			QuestInfo: &pb.VOQuestInfo{
				Quests: make([]*pb.VOQuest, 0, len(questIds)),
			},
		}
	}

	questInfo := u.Notifies.QuestUpdateNotify.QuestInfo
	for _, questId := range questIds {
		quest, _ := u.QuestPack.GetQuest(questId)
		questInfo.Quests = append(questInfo.Quests, quest.VOQuest())
	}

}

func (u *User) PopQuestUpdateNotify() *pb.S2CQuestUpdateNotify {
	notify := u.Notifies.QuestUpdateNotify
	u.Notifies.QuestUpdateNotify = nil

	if notify != nil && notify.QuestInfo != nil {
		notify.QuestInfo.NextDailyRefreshTime = u.QuestPack.GetLastDailyRefreshTime() + servertime.SecondPerDay
		notify.QuestInfo.NextWeeklyRefreshTime = u.QuestPack.GetLastWeeklyRefreshTime() + servertime.SecondPerWeek
	}

	return notify
}

//----------------------------------------
//Explore
//----------------------------------------
func (u *User) AddExploreEventNotify(eventId int32) {
	if !u.IsOnline() {
		return
	}

	if u.Notifies.ExploreEventNotify == nil {
		u.Notifies.ExploreEventNotify = &pb.S2CExploreEventNotify{
			ExploreEvent: make([]*pb.VOExploreEvent, 0, 1),
		}
	}

	event, _ := u.ExploreInfo.GetEventPoint(eventId)

	notify := u.Notifies.ExploreEventNotify
	notify.ExploreEvent = append(notify.ExploreEvent, event.VOExploreEvent())
}

func (u *User) PopExploreEventNotify() *pb.S2CExploreEventNotify {
	notify := u.Notifies.ExploreEventNotify

	u.Notifies.ExploreEventNotify = nil

	return notify
}

func (u *User) AddExploreResourceNotify(resourceId int32) {
	if !u.IsOnline() {
		return
	}

	if u.Notifies.ExploreResourceNotify == nil {
		u.Notifies.ExploreResourceNotify = &pb.S2CExploreResourceNotify{
			ResourcePoint: make([]*pb.VOExploreResource, 0, 1),
		}
	}

	resourcePoint, _ := u.ExploreInfo.GetResourcePoint(resourceId)

	notify := u.Notifies.ExploreResourceNotify
	notify.ResourcePoint = append(notify.ResourcePoint, resourcePoint.VOExploreResource())
}

func (u *User) PopExploreResourceNotify() *pb.S2CExploreResourceNotify {
	notify := u.Notifies.ExploreResourceNotify

	u.Notifies.ExploreResourceNotify = nil

	return notify
}

//----------------------------------------
//Mail
//----------------------------------------
func (u *User) AddMailNotify(mailId int64) {
	if !u.IsOnline() {
		return
	}

	mail, _ := u.MailInfo.GetMail(mailId)
	notify := u.getLastMailNotify()

	notify.Mails = append(notify.Mails, mail.VOMail())
}

func (u *User) AddMainRemoveNotify(mailId int64) {
	if !u.IsOnline() {
		return
	}

	notify := u.getLastMailNotify()
	notify.RemoveMails = append(notify.RemoveMails, mailId)
}

func (u *User) PopMailNotifies() []*pb.S2CMailNotify {
	notifies := u.Notifies.MailNotifies
	u.Notifies.MailNotifies = []*pb.S2CMailNotify{}

	return notifies
}

func (u *User) addMailNotifyBatch() *pb.S2CMailNotify {
	notify := &pb.S2CMailNotify{
		Mails:       []*pb.VOMail{},
		RemoveMails: []int64{},
	}
	u.Notifies.MailNotifies = append(u.Notifies.MailNotifies, notify)
	return notify
}

func (u *User) getLastMailNotify() *pb.S2CMailNotify {
	lastIdx := len(u.Notifies.MailNotifies) - 1

	if lastIdx < 0 {
		return u.addMailNotifyBatch()
	}

	notify := u.Notifies.MailNotifies[lastIdx]
	if len(notify.Mails) >= MailCountPerPage {
		return u.addMailNotifyBatch()
	}

	return notify
}

//----------------------------------------
//Tower
//----------------------------------------
func (u *User) AddTowerUpdateNotify(towerId int32) {
	if !u.IsOnline() {
		return
	}

	if u.Notifies.TowerNotify == nil {
		u.Notifies.TowerNotify = &pb.S2CTowerUpdateNotify{
			Towers: make([]*pb.VOTower, 0),
		}
	}

	tower, _ := u.TowerInfo.GetTower(towerId)

	notify := u.Notifies.TowerNotify
	notify.Towers = append(notify.Towers, tower.VOTower())
}

func (u *User) PopTowerUpdateNotify() *pb.S2CTowerUpdateNotify {
	notify := u.Notifies.TowerNotify
	u.Notifies.TowerNotify = nil

	return notify
}

//----------------------------------------
//Level
//----------------------------------------
func (u *User) AddLevelNotify(levelIds []int32) {
	if !u.IsOnline() {
		return
	}

	if u.Notifies.LevelNotify == nil {
		u.Notifies.LevelNotify = &pb.S2CLevelNotify{
			Levels: make([]*pb.VOLevel, 0, len(levelIds)),
		}
	}

	notify := u.Notifies.LevelNotify

	for _, levelId := range levelIds {
		level, _ := u.LevelsInfo.GetLevel(levelId)
		notify.Levels = append(notify.Levels, level.VOLevel())
	}
}

func (u *User) PopLevelNofity() *pb.S2CLevelNotify {
	notify := u.Notifies.LevelNotify
	u.Notifies.LevelNotify = nil

	return notify
}

func (u *User) AddYggPush(msg proto.Message) {
	u.Notifies.YggPush = append(u.Notifies.YggPush, msg)
}

func (u *User) AddGraveyardPush(msg proto.Message) {
	u.Notifies.GraveyardPush = append(u.Notifies.GraveyardPush, msg)
}

//----------------------------------------
//ScorePass
//----------------------------------------
func (u *User) AddScorePassNotify(seasonId int32) {
	if !u.IsOnline() {
		return
	}

	if u.Notifies.ScorePassNotify == nil {
		u.Notifies.ScorePassNotify = &pb.S2CScorePassNotify{
			Seasons: []*pb.VOScorePassSeason{},
		}
	}

	notify := u.Notifies.ScorePassNotify

	season, ok := u.ScorePassInfo.GetSeason(seasonId)
	if !ok {
		return
	}
	notify.Seasons = append(notify.Seasons, season.VOScorePassSeason())
}

func (u *User) PopScorePassNotify() *pb.S2CScorePassNotify {
	notify := u.Notifies.ScorePassNotify
	u.Notifies.ScorePassNotify = nil

	return notify
}

// ----------------------------------------
// store
// ----------------------------------------
func (u *User) AddStoreNotify() {
	if !u.IsOnline() {
		return
	}
	if u.Notifies.StoreDailyRefreshNotify == nil {
		u.Notifies.StoreDailyRefreshNotify = &pb.S2CStoreDailyRefresh{}
	}
}

func (u *User) PopStoreNotify() *pb.S2CStoreDailyRefresh {
	notify := u.Notifies.StoreDailyRefreshNotify
	u.Notifies.StoreDailyRefreshNotify = nil

	return notify
}
