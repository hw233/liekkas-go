package global

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"

	"shared/utility/errors"
	"shared/utility/key"
)

const keyIntimacy = "Intimacy:guild"

type Intimacy struct {
	*redis.Client
}

func NewIntimacy(client *redis.Client) *Intimacy {
	return &Intimacy{
		Client: client,
	}
}

func (i *Intimacy) genIntimacyHashField(userId1, userId2 int64) string {
	sortUserId1, sortUserId2 := userId1, userId2
	if userId1 > userId2 {
		sortUserId1, sortUserId2 = userId2, userId1
	}
	return key.MakeRedisKey(sortUserId1, sortUserId2)
}

func (i *Intimacy) genIntimacyHashKey(guildId int64) string {
	return key.MakeRedisKey(keyIntimacy, guildId)
}

// GetIntimacy 获得亲密度
func (i *Intimacy) GetIntimacy(ctx context.Context, guildId, userId1, userId2 int64) (int32, error) {
	hashKey := i.genIntimacyHashKey(guildId)
	v, err := i.HGet(ctx, hashKey, i.genIntimacyHashField(userId1, userId2)).Int64()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, errors.WrapTrace(err)
	}
	return int32(v), nil
}

// GetGuildIntimacyMap 我和公会其他成员默契值
func (i *Intimacy) GetGuildIntimacyMap(ctx context.Context, guildId int64, userId int64) (map[int64]int32, error) {
	hashKey := i.genIntimacyHashKey(guildId)
	all, err := i.HGetAll(ctx, hashKey).Result()
	ret := map[int64]int32{}
	if err != nil {
		if err == redis.Nil {
			return ret, nil
		}
		return nil, errors.WrapTrace(err)
	}

	for k, v := range all {
		//todo: 待优化
		split := strings.Split(k, ":")
		userId1, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		userId2, err := strconv.ParseInt(split[1], 10, 64)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		if userId1 == userId || userId2 == userId {
			intimacy, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, errors.WrapTrace(err)
			}
			var otherUserId int64
			if userId1 == userId {
				otherUserId = userId2
			} else {
				otherUserId = userId1

			}
			ret[otherUserId] = int32(intimacy)
		}
	}
	return ret, nil
}

// ChangeIntimacy 亲密度变化，返回变化后的值
func (i *Intimacy) ChangeIntimacy(ctx context.Context, guildId, userId1, userId2 int64, changeVal int32) (int32, error) {
	hashKey := i.genIntimacyHashKey(guildId)
	intimacy, err := i.HIncrBy(ctx, hashKey, i.genIntimacyHashField(userId1, userId2), int64(changeVal)).Result()
	if err != nil {
		return 0, errors.WrapTrace(err)
	}
	return int32(intimacy), nil
}

// ClearIntimacy 成员退出，清理与其相关的请密度
func (i *Intimacy) ClearIntimacy(ctx context.Context, guildId, userId int64, otherMembers ...int64) error {
	hashKey := i.genIntimacyHashKey(guildId)

	for _, memberUserId := range otherMembers {
		err := i.HDel(ctx, hashKey, i.genIntimacyHashField(userId, memberUserId)).Err()
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil
}

// ClearGuildIntimacy 公会解散 清理整个公会的亲密度
func (i *Intimacy) ClearGuildIntimacy(ctx context.Context, guildId int64) error {
	hashKey := i.genIntimacyHashKey(guildId)
	err := i.Del(ctx, hashKey).Err()
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}
