package global

import (
	"context"
	"shared/utility/glog"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

type LockerOpts struct {
	TTL           time.Duration
	RetryDuration time.Duration
	RetryTimes    int
}

var defaultLockerOpts = &LockerOpts{
	TTL:           10 * time.Second,
	RetryDuration: 50 * time.Millisecond,
	RetryTimes:    30,
}

type Locker struct {
	client *redislock.Client
	opts   *LockerOpts
}

type Lock struct {
	ctx  context.Context
	lock *redislock.Lock
}

func (l *Lock) Release() error {
	err := l.lock.Release(l.ctx)
	if err != nil {
		glog.Error(err)
	}
	return err
}

func NewLocker(client *redis.Client) *Locker {
	return &Locker{
		client: redislock.New(client),
		opts:   defaultLockerOpts,
	}
}

func (l *Locker) ObtainLock(ctx context.Context, key interface{}) (*Lock, error) {
	lock, err := l.client.Obtain(ctx, makeLockKey(key), l.opts.TTL, &redislock.Options{
		RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(l.opts.RetryDuration), l.opts.RetryTimes),
	})
	if err != nil {
		return nil, err
	}

	return &Lock{
		ctx:  ctx,
		lock: lock,
	}, nil
}

func (l *Locker) SetLockerOpts(opts *LockerOpts) {
	l.opts = opts
}
