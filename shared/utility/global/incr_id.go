package global

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type IncrID struct {
	client *redis.Client
}

func NewIncrID(client *redis.Client) *IncrID {
	return &IncrID{
		client: client,
	}
}

func (i *IncrID) GenID(ctx context.Context, key interface{}) (int64, error) {
	return i.client.Incr(ctx, makeIncrIDKey(key)).Result()
}

func (i *IncrID) SetInitID(ctx context.Context, key interface{}, init int64) error {
	ret, err := i.client.Get(ctx, makeIncrIDKey(key)).Int64()
	if err != nil {
		return err
	}

	if ret < init {
		err := i.client.Set(ctx, makeIncrIDKey(key), init, 0).Err()
		if err != nil {
			return err
		}
	}

	return nil
}
