package global

import (
	"context"
	"shared/utility/errors"
	"shared/utility/key"

	"github.com/go-redis/redis/v8"
)

const (
	KeyGreeting = "Greeting"
)

type Greetings struct {
	client *redis.Client
}

func NewGreetings(client *redis.Client) *Greetings {
	return &Greetings{
		client: client,
	}
}

func (g *Greetings) makePersonalHashKey(userId int64) string {
	return key.MakeRedisKey(KeyGreeting, userId)
}

func (g *Greetings) makeGreetingField(gid, gType int32) string {
	return key.MakeRedisKey(gid, gType)
}

func (g *Greetings) GetUserGreetingCount(ctx context.Context, userId int64, gid int32, gType int32) (int32, error) {
	hashKey := g.makePersonalHashKey(userId)

	count, err := g.client.HGet(ctx, hashKey, g.makeGreetingField(gid, gType)).Int()
	if err != nil {
		if err == redis.Nil {
			return 0, nil
		}
		return 0, errors.WrapTrace(err)
	}

	return int32(count), nil
}

func (g *Greetings) SetUserGreetingCount(ctx context.Context, userId int64, gid int32, gType int32, add int32) error {
	hashKey := g.makePersonalHashKey(userId)

	count, err := g.GetUserGreetingCount(ctx, userId, gid, gType)
	if err != nil {
		return errors.WrapTrace(err)
	}

	count += add

	return g.client.HSet(ctx, hashKey, g.makeGreetingField(gid, gType), count).Err()
}
