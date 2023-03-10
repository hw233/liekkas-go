package global

import (
	"context"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/go-redis/redis/v8"
)

type Hash struct {
	client *redis.Client
}

func NewHash(client *redis.Client) *Hash {
	return &Hash{
		client: client,
	}
}

func (h *Hash) HSet(ctx context.Context, key interface{}, field string, val ...interface{}) error {
	var tmp []interface{}
	tmp = append(append(tmp, field), val...)
	return h.client.HSet(ctx, makeHashKey(key), tmp...).Err()
}

func (h *Hash) HSetAll(ctx context.Context, key, src interface{}) error {
	return h.client.HSet(ctx, makeHashKey(key), redigo.Args{}.AddFlat(src)...).Err()
}

func (h *Hash) HDel(ctx context.Context, key interface{}, field ...string) error {
	return h.client.HDel(ctx, makeHashKey(key), field...).Err()
}

func (h *Hash) HGet(ctx context.Context, key interface{}, field string) (string, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Result()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetInt(ctx context.Context, key interface{}, field string) (int, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Int()
	if err == redis.Nil {
		return 0, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetInt8(ctx context.Context, key interface{}, field string) (int8, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Int64()
	if err == redis.Nil {
		return int8(ret), ErrNil
	}

	return int8(ret), err
}

func (h *Hash) HGetUint8(ctx context.Context, key interface{}, field string) (uint8, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Uint64()
	if err == redis.Nil {
		return uint8(ret), ErrNil
	}

	return uint8(ret), err
}

func (h *Hash) HGetInt16(ctx context.Context, key interface{}, field string) (int16, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Int64()
	if err == redis.Nil {
		return int16(ret), ErrNil
	}

	return int16(ret), err
}

func (h *Hash) HGetUint16(ctx context.Context, key interface{}, field string) (uint16, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Uint64()
	if err == redis.Nil {
		return uint16(ret), ErrNil
	}

	return uint16(ret), err
}

func (h *Hash) HGetInt32(ctx context.Context, key interface{}, field string) (int32, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Int64()
	if err == redis.Nil {
		return int32(ret), ErrNil
	}

	return int32(ret), err
}

func (h *Hash) HGetUint32(ctx context.Context, key interface{}, field string) (uint32, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Uint64()
	if err == redis.Nil {
		return uint32(ret), ErrNil
	}

	return uint32(ret), err
}

func (h *Hash) HGetInt64(ctx context.Context, key interface{}, field string) (int64, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Int64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetUint64(ctx context.Context, key interface{}, field string) (uint64, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Uint64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetBool(ctx context.Context, key interface{}, field string) (bool, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Bool()
	if err == redis.Nil {
		return false, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetBytes(ctx context.Context, key interface{}, field string) ([]byte, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Bytes()
	if err == redis.Nil {
		return []byte{}, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetFloat32(ctx context.Context, key interface{}, field string) (float32, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Float32()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetFloat64(ctx context.Context, key interface{}, field string) (float64, error) {
	ret, err := h.client.HGet(ctx, makeHashKey(key), field).Float64()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetAll(ctx context.Context, key interface{}) (map[string]string, error) {
	ret, err := h.client.HGetAll(ctx, makeHashKey(key)).Result()
	if err == redis.Nil {
		return ret, ErrNil
	}

	return ret, err
}

func (h *Hash) HGetAllScan(ctx context.Context, key, dst interface{}) error {
	ret := h.client.HGetAll(ctx, makeHashKey(key))
	err := ret.Err()
	if err == redis.Nil {
		return ErrNil
	}
	if err != nil {
		return err
	}

	return ret.Scan(dst)
}

func (h *Hash) DelHash(ctx context.Context, key interface{}) (bool, error) {
	ret, err := h.client.Del(ctx, makeHashKey(key)).Result()
	return ret > 0, err
}
