package module

import (
	"context"

	"github.com/go-redis/redis/v8"

	"shared/utility/errors"
	"shared/utility/key"
)

// for load balance
type Balancer struct {
	key    string
	client *redis.Client
}

func NewBalancer(service string, client *redis.Client) *Balancer {
	return &Balancer{
		key:    key.MakeRedisKey(GRPCKey, BLKey, service),
		client: client,
	}
}

func (b *Balancer) IncrBalance(ctx context.Context, server string) error {
	return b.client.ZAddArgsIncr(ctx, b.key, redis.ZAddArgs{
		Members: []redis.Z{
			{
				Score:  1,
				Member: server,
			},
		},
	}).Err()
}

func (b *Balancer) DecrBalance(ctx context.Context, server string) error {
	return b.client.ZAddArgsIncr(ctx, b.key, redis.ZAddArgs{
		Members: []redis.Z{
			{
				Score:  -1,
				Member: server,
			},
		},
	}).Err()
}

func (b *Balancer) SetBalance(ctx context.Context, server string, count int) error {
	return b.client.ZAdd(ctx, b.key, &redis.Z{
		Score:  float64(count),
		Member: server,
	}).Err()
}

func (b *Balancer) GetBalance(ctx context.Context) (string, error) {
	ret, err := b.client.ZRangeWithScores(ctx, b.key, 0, 0).Result()
	switch {
	case err != nil:
		return "", err
	case len(ret) == 0:
		return "", errors.New("no server")
	}

	// TODO: 这里可以限制服务器玩家数量
	// if values[1].(int) >= 1000000 {
	// 	return "", errors.New("server exhausted")
	// }

	return ret[0].Member.(string), nil
}

func (b *Balancer) DelBalance(ctx context.Context, server string) error {
	return b.client.ZRem(ctx, b.key, server).Err()
}
