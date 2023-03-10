package module

import (
	"context"

	"github.com/go-redis/redis/v8"

	"shared/utility/key"
)

// record user in which server
type Recorder struct {
	prefix string
	client *redis.Client
}

func NewRecorder(service string, client *redis.Client) *Recorder {
	return &Recorder{
		prefix: key.MakeRedisKey(RDKey, service),
		client: client,
	}
}

func (r *Recorder) SetRecord(ctx context.Context, server string, id int64) error {
	return r.client.SetNX(ctx, r.makeKey(id), server, 0).Err()
}

func (r *Recorder) DelRecord(ctx context.Context, id int64) error {
	return r.client.Del(ctx, r.makeKey(id)).Err()
}

func (r *Recorder) GetRecord(ctx context.Context, id int64) (string, error) {
	ret, err := r.client.Get(ctx, r.makeKey(id)).Result()
	if err == redis.Nil {
		return ret, nil
	}

	return ret, err
}

func (r *Recorder) makeKey(id int64) string {
	return key.MakeRedisKey(r.prefix, id)
}
