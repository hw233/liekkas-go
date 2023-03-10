package global

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type List struct {
	client *redis.Client
}

func NewList(client *redis.Client) *List {
	return &List{
		client: client,
	}
}

func (l *List) LPush(ctx context.Context, key interface{}, val ...interface{}) error {
	return l.client.LPush(ctx, makeListKey(key), val...).Err()
}

func (l *List) RPush(ctx context.Context, key interface{}, val ...interface{}) error {
	return l.client.RPush(ctx, makeListKey(key), val...).Err()
}

func (l *List) LPop(ctx context.Context, key interface{}) (string, error) {
	ret, err := l.client.LPop(ctx, makeListKey(key)).Result()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) LPopInt(ctx context.Context, key interface{}) (int, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Int()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) LPopInt8(ctx context.Context, key interface{}) (int8, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int8(ret), ErrNil
	}

	return int8(ret), err
}

func (l *List) LPopUint8(ctx context.Context, key interface{}) (uint8, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint8(ret), ErrNil
	}

	return uint8(ret), err
}

func (l *List) LPopInt16(ctx context.Context, key interface{}) (int16, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int16(ret), ErrNil
	}

	return int16(ret), err
}

func (l *List) LPopUint16(ctx context.Context, key interface{}) (uint16, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint16(ret), ErrNil
	}

	return uint16(ret), err
}

func (l *List) LPopInt32(ctx context.Context, key interface{}) (int32, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int32(ret), ErrNil
	}

	return int32(ret), err
}

func (l *List) LPopUint32(ctx context.Context, key interface{}) (uint32, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint32(ret), ErrNil
	}

	return uint32(ret), err
}

func (l *List) LPopInt64(ctx context.Context, key interface{}) (int64, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) LPopUint64(ctx context.Context, key interface{}) (uint64, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) LPopBool(ctx context.Context, key interface{}) (bool, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Bool()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) LPopFloat32(ctx context.Context, key interface{}) (float32, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Float32()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) LPopFloat64(ctx context.Context, key interface{}) (float64, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Float64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) LPopBytes(ctx context.Context, key interface{}) ([]byte, error) {
	ret, err := l.client.LPop(ctx, makeStringKey(key)).Bytes()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) RPop(ctx context.Context, key interface{}) (string, error) {
	return l.client.RPop(ctx, makeListKey(key)).Result()
}

func (l *List) RPopInt(ctx context.Context, key interface{}) (int, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Int()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) RPopInt8(ctx context.Context, key interface{}) (int8, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int8(ret), ErrNil
	}

	return int8(ret), err
}

func (l *List) RPopUint8(ctx context.Context, key interface{}) (uint8, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint8(ret), ErrNil
	}

	return uint8(ret), err
}

func (l *List) RPopInt16(ctx context.Context, key interface{}) (int16, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int16(ret), ErrNil
	}

	return int16(ret), err
}

func (l *List) RPopUint16(ctx context.Context, key interface{}) (uint16, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint16(ret), ErrNil
	}

	return uint16(ret), err
}

func (l *List) RPopInt32(ctx context.Context, key interface{}) (int32, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return int32(ret), ErrNil
	}

	return int32(ret), err
}

func (l *List) RPopUint32(ctx context.Context, key interface{}) (uint32, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return uint32(ret), ErrNil
	}

	return uint32(ret), err
}

func (l *List) RPopInt64(ctx context.Context, key interface{}) (int64, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Int64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) RPopUint64(ctx context.Context, key interface{}) (uint64, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Uint64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) RPopBool(ctx context.Context, key interface{}) (bool, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Bool()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) RPopFloat32(ctx context.Context, key interface{}) (float32, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Float32()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) RPopFloat64(ctx context.Context, key interface{}) (float64, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Float64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (l *List) RPopBytes(ctx context.Context, key interface{}) ([]byte, error) {
	ret, err := l.client.RPop(ctx, makeStringKey(key)).Bytes()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}
