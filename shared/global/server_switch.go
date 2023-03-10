package global

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const maintainSwitchKey = "maintainSwitch"
const whiteListSwitchKey = "whiteListSwitch"

type ServerSwitch struct {
	*redis.Client
}

func NewServerSwitch(client *redis.Client) *ServerSwitch {
	return &ServerSwitch{
		Client: client,
	}
}

func (m *ServerSwitch) FetchMaintainSwitch(ctx context.Context) (bool, error) {
	ret, err := m.Get(ctx, maintainSwitchKey).Int64()
	if err == redis.Nil {
		// 默认关
		return false, nil
	}
	return ret > 0, err
}

func (m *ServerSwitch) SetMaintainSwitch(ctx context.Context, open bool) error {
	before, err := m.FetchMaintainSwitch(ctx)
	if err != nil {
		return err
	}
	if before == open {
		return nil
	}
	var changeVal int64
	if open {
		changeVal = 1

	} else {
		changeVal = -1
	}
	_, err = m.Set(ctx, maintainSwitchKey, changeVal, 0).Result()
	return err

}

func (m *ServerSwitch) FetchWLSwitch(ctx context.Context) (bool, error) {
	ret, err := m.Get(ctx, whiteListSwitchKey).Int64()
	if err == redis.Nil {
		// 默认开
		return true, nil
	}
	return ret > 0, err
}

func (m *ServerSwitch) SetWLSwitch(ctx context.Context, open bool) error {
	before, err := m.FetchWLSwitch(ctx)
	if err != nil {
		return err
	}
	if before == open {
		return nil
	}
	var changeVal int64
	if open {
		changeVal = 1

	} else {
		changeVal = -1
	}
	_, err = m.Set(ctx, whiteListSwitchKey, changeVal, 0).Result()
	return err

}
