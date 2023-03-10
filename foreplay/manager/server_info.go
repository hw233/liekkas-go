package manager

import (
	"context"
	"time"
)

type ServerInfo struct {
	Maintain bool
}

func NewServerInfo() *ServerInfo {
	return &ServerInfo{
		Maintain: true,
	}
}

func (si *ServerInfo) Init() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err := si.LoadServerInfo(ctx)
	if err != nil {
		return err
	}

	Timer.ScheduleFunc(time.Minute, si.OnTick)

	return nil
}

func (si *ServerInfo) LoadServerInfo(ctx context.Context) error {
	maintain, err := Global.FetchMaintainSwitch(ctx)
	if err != nil {
		return err
	}

	si.Maintain = maintain

	return nil
}

func (si *ServerInfo) IsMaintain() bool {
	return si.Maintain
}

func (si *ServerInfo) OnTick() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	si.LoadServerInfo(ctx)
}
