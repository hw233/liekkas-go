package model

import (
	"gamesvr/manager"
	"shared/global"
	"shared/protobuf/pb"
)

type UserVisitingCard struct {
	UserId    int64
	Name      string
	NameIndex int32
	TeamLevel int32
	TeamExp   int32
	Avatar    int32
	Frame     int32
	GuildId   int64
	GuildName string

	RegisterAt          int64
	Signature           string
	ExploreLevel        int32
	ManualProgress      int32
	AchievementProgress int32
	WorldItemId         int32
	WorldItemUid        int64
	Characters          map[int32]*pb.VOUserCharacter
}

func NewOwnUserVisitingCard(userId int64, nickname string, info *UserInfo, guild *UserGuild) *UserVisitingCard {

	manualCount := int32(len(info.ExploredLevelVns.Values())) + info.ManualCount

	return &UserVisitingCard{
		UserId:              userId,
		Name:                nickname,
		NameIndex:           info.NameIndex,
		TeamLevel:           info.Level.Value(),
		TeamExp:             info.Exp.Value(),
		Avatar:              info.Avatar,
		Frame:               info.Frame,
		GuildId:             guild.GuildID,
		GuildName:           guild.GuildName,
		RegisterAt:          info.RegisterAt,
		Signature:           info.Signature,
		ExploreLevel:        info.GetLatestExploreLevel(),
		ManualProgress:      manualCount * 100 / manager.CSV.Manual.GetManualTotalCount(),
		AchievementProgress: info.GetAchievementFinishCount() * 100 / manager.CSV.Quest.GetAchievementTotalCount(),
		Characters:          info.CardShow.CharactersCache,
		WorldItemId:         info.CardShow.WorldItemId,
		WorldItemUid:        info.CardShow.WorldItemUId,
	}
}

func NewOthersUserVisitingCard(cache *global.UserCache) *UserVisitingCard {

	manualCount := int32(len(cache.ExploredLevelVns.Values())) + cache.ManualCount

	return &UserVisitingCard{
		UserId:              cache.ID,
		Name:                cache.Name,
		NameIndex:           cache.NameIndex,
		TeamLevel:           cache.Level,
		TeamExp:             cache.Exp,
		Avatar:              cache.Avatar,
		Frame:               cache.Frame,
		GuildId:             cache.GuildID,
		GuildName:           cache.GuildName,
		RegisterAt:          cache.RegisterAt,
		Signature:           cache.Signature,
		ExploreLevel:        cache.LatestExploreLevel,
		ManualProgress:      manualCount * 100 / manager.CSV.Manual.GetManualTotalCount(),
		AchievementProgress: cache.AchievementFinishCount * 100 / manager.CSV.Quest.GetAchievementTotalCount(),
		Characters:          cache.CardShow.CharactersCache,
		WorldItemId:         cache.CardShow.WorldItemId,
		WorldItemUid:        cache.CardShow.WorldItemUId,
	}
}

func (u *UserVisitingCard) VOVisitingCard() *pb.VOVisitingCard {
	characters := make([]*pb.VOCharacterVisitingCard, 0, len(u.Characters))
	for pos, characterVo := range u.Characters {
		characters = append(characters, &pb.VOCharacterVisitingCard{
			Pos:       pos,
			Character: characterVo,
		})
	}

	return &pb.VOVisitingCard{
		UserInfo: &pb.VOUserInfo{
			UserId:    u.UserId,
			Nickname:  u.Name,
			NameIndex: u.NameIndex,
			TeamLevel: u.TeamLevel,
			TeamExp:   u.TeamExp,
			Avatar:    u.Avatar,
			Frame:     u.Frame,
			GuildId:   u.GuildId,
			GuildName: u.GuildName,
		},
		RegisterAt:          u.RegisterAt,
		Signature:           u.Signature,
		ExploreLevel:        u.ExploreLevel,
		ManualProgress:      u.ManualProgress,
		AchievementProgress: u.AchievementProgress,
		Characters:          characters,
		WorldItemId:         u.WorldItemId,
		WorldItemUId:        u.WorldItemUid,
	}
}
