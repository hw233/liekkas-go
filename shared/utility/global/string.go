package global

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type String struct {
	client *redis.Client
}

func NewString(client *redis.Client) *String {
	return &String{
		client: client,
	}
}

func (s *String) Get(ctx context.Context, key interface{}) (string, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Result()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (s *String) GetInt(ctx context.Context, key interface{}) (int, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Int()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (s *String) GetInt8(ctx context.Context, key interface{}) (int8, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int8(ret), ErrNil
	}

	return int8(ret), err
}

func (s *String) GetUint8(ctx context.Context, key interface{}) (uint8, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint8(ret), ErrNil
	}

	return uint8(ret), err
}

func (s *String) GetInt16(ctx context.Context, key interface{}) (int16, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int16(ret), ErrNil
	}

	return int16(ret), err
}

func (s *String) GetUint16(ctx context.Context, key interface{}) (uint16, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint16(ret), ErrNil
	}

	return uint16(ret), err
}

func (s *String) GetInt32(ctx context.Context, key interface{}) (int32, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int32(ret), ErrNil
	}

	return int32(ret), err
}

func (s *String) GetUint32(ctx context.Context, key interface{}) (uint32, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint32(ret), ErrNil
	}

	return uint32(ret), err
}

func (s *String) GetInt64(ctx context.Context, key interface{}) (int64, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (s *String) GetUint64(ctx context.Context, key interface{}) (uint64, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (s *String) GetBool(ctx context.Context, key interface{}) (bool, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Bool()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (s *String) GetFloat32(ctx context.Context, key interface{}) (float32, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Float32()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (s *String) GetFloat64(ctx context.Context, key interface{}) (float64, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Float64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (s *String) GetBytes(ctx context.Context, key interface{}) ([]byte, error) {
	ret, err := s.client.Get(ctx, makeStringKey(key)).Bytes()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (s *String) Set(ctx context.Context, key, val interface{}) error {
	return s.client.Set(ctx, makeStringKey(key), val, 0).Err()
}

func (s *String) SetEX(ctx context.Context, key, val interface{}, ex time.Duration) error {
	return s.client.SetEX(ctx, makeStringKey(key), val, ex).Err()
}

func (s *String) Incr(ctx context.Context, key interface{}) (int64, error) {
	return s.client.Incr(ctx, makeStringKey(key)).Result()
}
