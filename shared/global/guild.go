package global

import (
	"context"
	"strconv"

	"github.com/go-redis/redis/v8"

	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/key"
)

const (
	// field
	FieldGuildID   = "guild_id"
	FieldGuildName = "guild_name"
)

type Guild struct {
	*GuildID
	*GuildName
	*GuildUser
	*GuildSet
}

func NewGuild(client *redis.Client) *Guild {
	return &Guild{
		GuildID:   NewGuildID(client),
		GuildName: NewGuildName(client),
		GuildUser: NewGuildUser(client),
		GuildSet:  NewGuildSet(client),
	}
}

// 生成公会ID
type GuildID struct {
	key  string
	incr *global.IncrID
}

func NewGuildID(client *redis.Client) *GuildID {
	return &GuildID{
		key:  KeyGuildID,
		incr: global.NewIncrID(client),
	}
}

func (g *GuildID) GenGuildID(ctx context.Context) (int64, error) {
	return g.incr.GenID(ctx, g.key)
}

// 公会名称去重和索引
type GuildName struct {
	key    string
	client *redis.Client
}

func NewGuildName(client *redis.Client) *GuildName {
	return &GuildName{
		key:    KeyGuildName,
		client: client,
	}
}

func (g *GuildName) GuildNameExist(ctx context.Context, guildName string) (bool, error) {
	_, err := g.client.HGet(ctx, g.key, guildName).Result()
	if err == redis.Nil {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func (g *GuildName) AddGuildNameIfNotExist(ctx context.Context, guildID int64, guildName string) (bool, error) {
	return g.client.HSetNX(ctx, g.key, guildName, guildID).Result()
}

func (g *GuildName) GetGuildID(ctx context.Context, guildName string) (int64, error) {
	return g.client.HGet(ctx, g.key, guildName).Int64()
}

func (g *GuildName) DelGuildName(ctx context.Context, guildName string) error {
	return g.client.HDel(ctx, g.key, guildName).Err()
}

// 玩家公会数据
type GuildUser struct {
	key    string
	client *redis.Client
}

func NewGuildUser(client *redis.Client) *GuildUser {
	return &GuildUser{
		key:    KeyGuildUser,
		client: client,
	}
}

func (g *GuildUser) SetUserGuildData(ctx context.Context, userID, guildID int64, guildName string) error {
	return g.client.HSet(ctx, key.MakeRedisKey(g.key, userID), FieldGuildID, guildID, FieldGuildName, guildName).Err()
}

func (g *GuildUser) GetUserGuildData(ctx context.Context, userID int64) (*UserGuild, error) {
	userGuild := &UserGuild{}

	err := g.client.HGetAll(ctx, key.MakeRedisKey(g.key, userID)).Scan(userGuild)
	if err != nil {
		return nil, err
	}

	return userGuild, nil
}

func (g *GuildUser) DelUserGuildData(ctx context.Context, userID int64) error {
	return g.client.Del(ctx, key.MakeRedisKey(g.key, userID)).Err()
}

// 玩家公会数据
type GuildSet struct {
	key    string
	client *redis.Client
	Lock   *global.Locker
}

func NewGuildSet(client *redis.Client) *GuildSet {
	return &GuildSet{
		key:    KeyGuildSet,
		client: client,
		Lock:   global.NewLocker(client),
	}
}

func (gs *GuildSet) GuildSetAdd(ctx context.Context, guildId int64) error {
	hashKey := key.MakeRedisKey(gs.key)

	err := gs.client.SAdd(ctx, hashKey, guildId).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (gs *GuildSet) GuildSetDelete(ctx context.Context, guildId int64) error {
	hashKey := key.MakeRedisKey(gs.key)

	err := gs.client.SRem(ctx, hashKey, guildId).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func (gs *GuildSet) GuildSetRandomGet(ctx context.Context, num int64) ([]int64, error) {
	hashKey := key.MakeRedisKey(gs.key)

	result := make([]int64, 0, num)

	ids, err := gs.client.SRandMemberN(ctx, hashKey, num).Result()
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	for _, id := range ids {
		idInt, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		result = append(result, idInt)
	}
	return result, nil
}
