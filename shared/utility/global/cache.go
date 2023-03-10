package global

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/go-redis/redis/v8"

	"shared/utility/key"
)

type CacheBuilder interface {
	NewCache() Cache
}

type Cache interface {
	Key() int64
	Load(ctx context.Context, key int64) error
}

type ImplementedCache struct {
	DB     *sql.DB       `json:"-"`
	Client *redis.Client `json:"-"`
}

func (c ImplementedCache) Key() int64 {
	return 0
}

func (c ImplementedCache) Load(ctx context.Context, key int64) error {
	return nil
}

type CacheHandler struct {
	key     string
	client  *redis.Client
	builder func() Cache
	expire  time.Duration
}

func NewCacheHandler(client *redis.Client, keyPrefix string, builder func() Cache, expire time.Duration) *CacheHandler {
	return &CacheHandler{
		key:     keyPrefix,
		client:  client,
		builder: builder,
		expire:  expire,
	}
}

func (h *CacheHandler) loadCaches(ctx context.Context, ids []int64) ([]Cache, error) {
	caches := make([]Cache, 0, len(ids))

	for _, id := range ids {
		cache := h.builder()

		err := cache.Load(ctx, id)
		if err != nil {
			log.Printf("load cache error: %v", err)
			continue
		}

		caches = append(caches, cache)
	}

	return caches, nil
}

func (h *CacheHandler) GetCaches(ctx context.Context, ids []int64) (map[int64]Cache, error) {
	pipeline := h.client.Pipeline()
	defer pipeline.Close()

	for _, id := range ids {
		pipeline.Get(ctx, key.MakeRedisKey(h.key, id))
	}

	cmdrs, err := pipeline.Exec(ctx)
	if err != nil && err != redis.Nil {
		return nil, err
	}

	ret := map[int64]Cache{}

	for _, cmdr := range cmdrs {
		if cmdr.Err() != nil {
			continue
		}

		cmd, ok := cmdr.(*redis.StringCmd)
		if !ok {
			continue
		}

		val, err := cmd.Result()
		if err != nil {
			continue
		}

		cache := h.builder()

		err = json.Unmarshal([]byte(val), cache)
		if err != nil {
			continue
		}

		ret[cache.Key()] = cache
	}

	// 有未命中缓存
	if len(ids) > len(ret) {
		unHit := make([]int64, 0, len(ids)-len(ret))

		for _, id := range ids {
			if _, ok := ret[id]; !ok {
				unHit = append(unHit, id)
			}
		}

		caches, err := h.loadCaches(ctx, unHit)
		if err != nil {
			return nil, err
		}

		for _, cache := range caches {
			ret[cache.Key()] = cache

			val, err := json.Marshal(cache)
			if err != nil {
				return nil, err
			}

			err = h.client.Set(ctx, key.MakeRedisKey(h.key, cache.Key()), string(val), h.expire).Err()
			if err != nil {
				return nil, err
			}
		}
	}

	return ret, nil
}
