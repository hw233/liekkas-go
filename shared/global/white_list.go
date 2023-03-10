package global

import (
	"context"
	"github.com/go-redis/redis/v8"
	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/key"
	"shared/utility/servertime"
	"shared/utility/whitelist"
	"strconv"
)

const whiteListKey = "white:list"

type WhiteList struct {
	*global.Global
}

func NewWhiteList(client *redis.Client) *WhiteList {
	return &WhiteList{
		Global: global.NewGlobal(client),
	}
}

func (i *WhiteList) genWLKey(t whitelist.WhiteListType) string {
	return key.MakeRedisKey(whiteListKey, t)
}

func (i *WhiteList) GetIdWhiteList(ctx context.Context) (*whitelist.Options, error) {
	hashKey := i.genWLKey(whitelist.Id)
	// 分布式锁
	gLock, err := i.ObtainLock(ctx, hashKey)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	defer gLock.Release()
	ret := whitelist.EmptyOps()
	m, err := i.HGetAll(ctx, hashKey)
	if err != nil {
		if err == global.ErrNil {
			return ret, nil
		}
		return nil, errors.WrapTrace(err)
	}
	for k := range m {
		uid, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
		ret.With(whitelist.Id, uid)
	}
	return ret, nil
}

func (i *WhiteList) DelIdWhiteList(ctx context.Context, uid int64) error {
	hashKey := i.genWLKey(whitelist.Id)
	// 分布式锁
	gLock, err := i.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	err = i.HDel(ctx, hashKey, key.MakeRedisKey(uid))
	if err != nil {
		return errors.WrapTrace(err)
	}
	return err
}

func (i *WhiteList) AddIdWhiteList(ctx context.Context, uid int64) error {
	hashKey := i.genWLKey(whitelist.Id)
	// 分布式锁
	gLock, err := i.ObtainLock(ctx, hashKey)
	if err != nil {
		return errors.WrapTrace(err)
	}
	defer gLock.Release()

	err = i.HSet(ctx, hashKey, key.MakeRedisKey(uid), servertime.Now().Unix())
	if err != nil {
		return errors.WrapTrace(err)
	}
	return err
}
