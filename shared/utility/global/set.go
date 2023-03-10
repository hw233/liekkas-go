package global

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Set struct {
	client *redis.Client
}

func NewSet(client *redis.Client) *Set {
	return &Set{
		client: client,
	}
}

func (s *Set) SAdd(ctx context.Context, key, val interface{}) (bool, error) {
	ret, err := s.client.SAdd(ctx, makeSetKey(key), val).Result()
	if err != nil {
		return false, err
	}

	return ret == 1, nil
}

func (s *Set) SRem(ctx context.Context, key, val interface{}) error {
	return s.client.SRem(ctx, makeSetKey(key), val).Err()
}

func (s *Set) SMembers(ctx context.Context, key interface{}) ([]string, error) {
	ret, err := s.client.SMembers(ctx, makeSetKey(key)).Result()
	if err == redis.Nil {
		return ret, ErrNil
	}
	return ret, err

}

func (s *Set) SDel(ctx context.Context, key interface{}) error {
	_, err := s.client.Del(ctx, makeSetKey(key)).Result()
	return err
}
