package model

import (
	"context"
	"gamesvr/manager"
	"shared/utility/errors"
	"shared/utility/global"
	"shared/utility/key"
)

const (
	RedisKeyBuild        = "build"
	RedisKeyMessage      = "message"
	RedisKeyDiscardGoods = "discardGoods"
	RedisKeyMailNum      = "mailNum"
)

func AppendYggBuildInRedis(ctx context.Context, userId int64, build *YggBuild) error {
	if !build.IsOwn() {
		return nil
	}
	err := manager.Global.Hash.HSet(ctx, key.MakeRedisKey(RedisKeyBuild, userId), key.MakeRedisKey(build.Uid), build)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func AppendYggMessageInRedis(ctx context.Context, userId int64, message *YggMessage) error {
	if !message.IsOwn() {
		return nil
	}

	err := manager.Global.Hash.HSet(ctx, key.MakeRedisKey(RedisKeyMessage, userId), key.MakeRedisKey(message.Uid), message)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func AppendYggDiscardGoodsInRedis(ctx context.Context, userId int64, good *YggDiscardGoods) error {
	if !good.IsOwn() {
		return nil
	}
	err := manager.Global.Hash.HSet(ctx, key.MakeRedisKey(RedisKeyDiscardGoods, userId), key.MakeRedisKey(good.Uid), good)
	if err != nil {
		return errors.WrapTrace(err)
	}
	return nil
}

func DelYggBuildInRedis(ctx context.Context, userId int64, build *YggBuild) error {
	if !build.IsOwn() {
		return nil
	}
	return manager.Global.Hash.HDel(ctx, key.MakeRedisKey(RedisKeyBuild, userId), key.MakeRedisKey(build.Uid))
}

func DelYggMessageInRedis(ctx context.Context, userId int64, message *YggMessage) error {
	if !message.IsOwn() {
		return nil
	}
	return manager.Global.Hash.HDel(ctx, key.MakeRedisKey(RedisKeyMessage, userId), key.MakeRedisKey(message.Uid))
}

func DelYggDiscardGoodsInRedis(ctx context.Context, userId int64, good *YggDiscardGoods) error {
	if !good.IsOwn() {
		return nil
	}
	return manager.Global.Hash.HDel(ctx, key.MakeRedisKey(RedisKeyDiscardGoods, userId), key.MakeRedisKey(good.Uid))
}

func SetYggMailNumInRedis(ctx context.Context, userId int64, mailNum int) error {
	return manager.Global.Hash.HSet(ctx, RedisKeyMailNum, key.MakeRedisKey(userId), mailNum)
}
func LoadUserMatchEntities(ctx context.Context, userId int64) (*UserMatchEntities, error) {

	allBuild, err := manager.Global.Hash.HGetAll(ctx, key.MakeRedisKey(RedisKeyBuild, userId))
	if err != nil && err != global.ErrNil {
		return nil, errors.WrapTrace(err)
	}
	allMessage, err := manager.Global.Hash.HGetAll(ctx, key.MakeRedisKey(RedisKeyMessage, userId))
	if err != nil && err != global.ErrNil {
		return nil, errors.WrapTrace(err)
	}
	allDiscardGoods, err := manager.Global.Hash.HGetAll(ctx, key.MakeRedisKey(RedisKeyDiscardGoods, userId))
	if err != nil && err != global.ErrNil {
		return nil, errors.WrapTrace(err)
	}
	yggMailNum, err := manager.Global.String.GetInt(ctx, key.MakeRedisKey(RedisKeyMailNum, userId))
	if err != nil && err != global.ErrNil {
		return nil, errors.WrapTrace(err)
	}

	return NewUserMatchEntities(allBuild, allMessage, allDiscardGoods, yggMailNum)

}
