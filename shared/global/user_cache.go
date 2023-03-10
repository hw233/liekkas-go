package global

import (
	"context"
	"shared/utility/errors"

	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/global"
	"shared/utility/mysql"
	"shared/utility/number"
	"shared/utility/uid"
)

type UserHelper struct {
	ID                 int64     `db:"id" major:"true"`
	Name               string    `db:"name"`
	Info               *UserInfo `db:"info"`
	*mysql.EmbedModule `db:"-"`
}

func NewUserHelper(userID int64) *UserHelper {
	return &UserHelper{
		ID:   userID,
		Info: &UserInfo{},

		EmbedModule: &mysql.EmbedModule{},
	}
}

type EquipmentPack struct {
	UID        *uid.UID                    `json:"uid"`        // 生成唯一ID
	Equipments map[int64]*common.Equipment `json:"equipments"` // 装备数据
}

type WorldItemPack struct {
	RecastID   int64                       `json:"recast_id"`   // 重铸随机的阵营
	RecastCamp int8                        `json:"recast_camp"` // 重铸随机的阵营
	UID        *uid.UID                    `json:"uid"`         // 生成唯一ID
	WorldItems map[int64]*common.WorldItem `json:"equipments"`  // 装备数据
}

type UserInfo struct {
	NameIndex              int32                    `json:"name_index"` // 昵称后缀
	Exp                    int32                    `json:"exp"`
	Level                  int32                    `json:"level"`
	Avatar                 int32                    `json:"avatar"`
	Frame                  int32                    `json:"frame"`
	Signature              string                   `json:"signature"`
	Power                  int32                    `json:"power"`       // 总战力
	RegisterAt             int64                    `json:"register_at"` // 注册日期
	CardShow               *common.VisitingCardShow `json:"card_show"`   // 名片展示
	LatestExploreLevel     int32                    `json:"latest_explore_level"`
	ExploredLevelVns       *number.NonRepeatableArr `json:"explored_level_vns"`
	AchievementFinishCount int32                    `json:"achievement_finish_count"`
	ManualCount            int32                    `json:"manual_count"`    // 总获得的图鉴数（角色和世界级道具）
	LastLoginTime          int64                    `json:"last_login_time"` // 上次登录时间
}

type UserGuild struct {
	GuildID   int64  `json:"guild_id" redis:"guild_id"`
	GuildName string `json:"guild_name" redis:"guild_name"`
}

type UserCache struct {
	*global.ImplementedCache

	ID        int64  `json:"id"`
	Name      string `json:"name"`
	NameIndex int32  `json:"name_index"`
	Exp       int32  `json:"exp"`
	Level     int32  `json:"level"`
	Head      int32  `json:"head"`
	Frame     int32  `json:"frame"`
	Avatar    int32  `json:"avatar"`
	Signature string `json:"signature"`

	Power                  int32                    `json:"power"`       // 总战力
	RegisterAt             int64                    `json:"register_at"` // 注册日期
	CardShow               *common.VisitingCardShow `json:"card_show"`   // 名片展示
	LatestExploreLevel     int32                    `json:"latest_explore_level"`
	ExploredLevelVns       *number.NonRepeatableArr `json:"explored_level_vns"`
	AchievementFinishCount int32                    `json:"achievement_finish_count"`
	ManualCount            int32                    `json:"manual_count"` // 总获得的图鉴数（角色和世界级道具）

	// extension
	// UserCacheWithGuild
	GuildID   int64  `json:"guild_id"`
	GuildName string `json:"guild_name"`

	// UserCacheWithOnline
	OnlineStatus  int32 `json:"online_status"`   // 在线状态, 0:离线 1:在线
	LastLoginTime int64 `json:"last_login_time"` // 上次登录时间
}

func (uc *UserCache) VOUserInfoSimple() *pb.VOUserInfoSimple {
	return &pb.VOUserInfoSimple{
		UserId:    uc.ID,
		NameIndex: uc.NameIndex,
		Nickname:  uc.Name,
		Avatar:    uc.Avatar,
		TeamLevel: uc.Level,
	}
}

type CharacterPack map[int32]*Character

type Character struct {
	ID               int32           `json:"id"`
	Exp              int32           `json:"exp"`
	Level            int32           `json:"level"`
	Star             int32           `json:"star"`
	Stage            int32           `json:"stage"`
	CTime            int64           `json:"c_time"`
	Skills           map[int32]int32 `json:"skills"`
	HeroId           int32           `json:"hero_id"`
	HeroSlot         int32           `json:"hero_slot"`
	Power            int32           `json:"power"`
	CanYggdrasilTime int64           `json:"can_yggdrasil_time"`
	Rarity           int32           `json:"rarity"`
	Equipments       []int64         `json:"equipments"`
	WorldItem        int64           `json:"world_item"`
}

func (c *Character) VOUserCharacter() *pb.VOUserCharacter {
	skills := make([]*pb.VOCharacterSkill, 0, len(c.Skills))
	for skillId, level := range c.Skills {
		skills = append(skills, &pb.VOCharacterSkill{
			SkillId: skillId,
			Level:   level,
		})
	}

	return &pb.VOUserCharacter{
		CharacterId:      c.ID,
		Exp:              c.Exp,
		Level:            c.Level,
		Star:             c.Star,
		Stage:            c.Stage,
		CreateAt:         c.CTime,
		Skills:           skills,
		HeroId:           c.HeroId,
		HeroPos:          c.HeroSlot,
		Power:            c.Power,
		CanYggdrasilTime: c.CanYggdrasilTime,
		Rarity:           c.Rarity,
		// Equipments:      c.Equipments,
		// WorldItem:        c.WorldItem,
	}
}

func (uc *UserCache) Load(ctx context.Context, id int64) error {
	helper := NewUserHelper(id)

	handler := mysql.NewHandler(uc.DB)
	err := handler.Init(helper)
	if err != nil {
		return errors.WrapTrace(err)
	}
	helper.SetTable("user")

	err = handler.Load(ctx, helper)
	if err != nil {
		return err
	}

	uc.ID = helper.ID
	uc.Name = helper.Name
	uc.NameIndex = helper.Info.NameIndex
	uc.Exp = helper.Info.Exp
	uc.Level = helper.Info.Level
	uc.Avatar = helper.Info.Avatar
	uc.Frame = helper.Info.Frame
	uc.Signature = helper.Info.Signature

	uc.Power = helper.Info.Power
	uc.RegisterAt = helper.Info.RegisterAt
	uc.CardShow = helper.Info.CardShow
	uc.LatestExploreLevel = helper.Info.LatestExploreLevel
	uc.ExploredLevelVns = helper.Info.ExploredLevelVns
	uc.AchievementFinishCount = helper.Info.AchievementFinishCount
	uc.ManualCount = helper.Info.ManualCount
	uc.LastLoginTime = helper.Info.LastLoginTime
	return nil
}

func (uc *UserCache) Key() int64 {
	return uc.ID
}
