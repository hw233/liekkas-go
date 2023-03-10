package global

import (
	"context"
	"errors"

	"github.com/go-redis/redis/v8"
)

var ErrNil = errors.New("nil")

type Global struct {
	*Hash
	*Set
	*String
	*Locker
	*IncrID
	*List
}

func NewGlobal(client *redis.Client) *Global {
	global := &Global{
		Hash:   NewHash(client),
		Set:    NewSet(client),
		String: NewString(client),
		Locker: NewLocker(client),
		IncrID: NewIncrID(client),
	}

	return global
}

func (s *String) Del(ctx context.Context, key interface{}) (bool, error) {
	ret, err := s.client.Del(ctx, makeStringKey(key)).Result()
	return ret > 0, err
}
