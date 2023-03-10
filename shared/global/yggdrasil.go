package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"shared/utility/global"
	"shared/utility/key"
)

const (
	KeyYggEntities           = "YggEntities"
	KeyYggBuildTotalUseCount = "YggBuildTotalUseCount"
	KeyYggGoodsTakenBack     = "YggGoodsTakenBack"
)

type Yggdrasil struct {
	*YggdrasilEntityUid
	*YggBuildTotalUseCount
	*YggGoodsTakenBack
}

func NewYggdrasil(client *redis.Client) *Yggdrasil {
	return &Yggdrasil{
		YggdrasilEntityUid:    NewYggdrasilEntityUid(client),
		YggBuildTotalUseCount: NewYggBuildTotalUseCount(client),
		YggGoodsTakenBack:     NewYggGoodsTakenBack(client),
	}
}

type YggdrasilEntityUid struct {
	incr *global.IncrID
}

func NewYggdrasilEntityUid(client *redis.Client) *YggdrasilEntityUid {
	return &YggdrasilEntityUid{
		incr: global.NewIncrID(client),
	}

}
func (y *YggdrasilEntityUid) GenYggdrasilEntityUid(ctx context.Context) (int64, error) {
	return y.incr.GenID(ctx, KeyYggEntities)
}

type YggBuildTotalUseCount struct {
	client *redis.Client
}

func NewYggBuildTotalUseCount(client *redis.Client) *YggBuildTotalUseCount {
	return &YggBuildTotalUseCount{
		client: client,
	}
}

func (r *YggBuildTotalUseCount) FetchBuildUseNum(ctx context.Context, buildUid int64) (int32, error) {
	count, err := r.client.HGet(ctx, KeyYggBuildTotalUseCount, key.MakeRedisKey(buildUid)).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}
func (r *YggBuildTotalUseCount) IncrBuildUseNum(ctx context.Context, buildUid int64) (int32, error) {
	count, err := r.client.HIncrBy(ctx, KeyYggBuildTotalUseCount, key.MakeRedisKey(buildUid), 1).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}

	return int32(count), nil
}

// YggGoodsTakenBack 记录每个goods是否被带回
type YggGoodsTakenBack struct {
	client *redis.Client
}

func NewYggGoodsTakenBack(client *redis.Client) *YggGoodsTakenBack {
	return &YggGoodsTakenBack{
		client: client,
	}
}

func (r *YggGoodsTakenBack) FetchGoodsTakenBack(ctx context.Context, goodsUid int64) (bool, error) {
	takenBack, err := r.client.HGet(ctx, KeyYggGoodsTakenBack, key.MakeRedisKey(goodsUid)).Int64()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return takenBack > 0, nil
}
func (r *YggGoodsTakenBack) SetGoodsTakenBack(ctx context.Context, goodsUid int64) error {
	_, err := r.client.HIncrBy(ctx, KeyYggGoodsTakenBack, key.MakeRedisKey(goodsUid), 1).Result()
	if err == redis.Nil {
		return nil
	}
	if err != nil {
		return err
	}

	return nil
}
