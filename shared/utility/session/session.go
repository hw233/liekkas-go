package session

import (
	"context"
	"time"

	"shared/utility/number"

	"github.com/antlabs/timer"
)

type CreateOpts struct {
	ID            int64
	Timer         timer.Timer
	Implemented   implementedSession
	WrapperClosed func()
	AllowNil      bool
}

type Session interface {
	OnCreated(context.Context, OnCreatedOpts) error
	OnTriggered(context.Context)
	OnClosed()

	AfterCall(expire time.Duration, callback func())
	ScheduleCall(expire time.Duration, callback func())
}

type OnCreatedOpts struct {
	ID       int64
	AllowNil bool
}

// for manager
type managedSession interface {
	Create(context.Context, CreateOpts) error
	Trigger(context.Context)
	Close()
	Expired()
}

// for rewrite
type implementedSession interface {
	OnCreated(context.Context, OnCreatedOpts) error
	OnTriggered(context.Context)
	OnClosed()
}

type unimplementedSession struct{}

func (s *unimplementedSession) OnCreated(context.Context, OnCreatedOpts) error { return nil }
func (s *unimplementedSession) OnTriggered(context.Context)                    {}
func (s *unimplementedSession) OnClosed()                                      {}

// 线程不安全，按需求自己加锁
type EmbedManagedSession struct {
	*unimplementedSession
	implemented  implementedSession
	id           int64
	timeNodes    []timer.TimeNoder
	alter        *number.CalNumber
	timer        timer.Timer
	wrapperClose func()
}

func (s *EmbedManagedSession) Create(ctx context.Context, opts CreateOpts) error {
	s.unimplementedSession = &unimplementedSession{}
	s.id = opts.ID
	s.timeNodes = make([]timer.TimeNoder, 0, 1)
	s.alter = number.NewCalNumber(0)
	s.timer = opts.Timer
	s.implemented = opts.Implemented
	s.wrapperClose = opts.WrapperClosed
	return s.implemented.OnCreated(ctx, OnCreatedOpts{ID: opts.ID, AllowNil: opts.AllowNil})
}

func (s *EmbedManagedSession) Trigger(ctx context.Context) {
	s.alter.Plus(1)
	s.implemented.OnTriggered(ctx)
}

func (s *EmbedManagedSession) Expired() {
	if s.alter.Equal(0) {
		s.Close()
		return
	}

	s.alter.SetValue(0)
}

func (s *EmbedManagedSession) ScheduleCall(expire time.Duration, f func()) {
	s.timeNodes = append(s.timeNodes, s.timer.ScheduleFunc(expire, f))
}

func (s *EmbedManagedSession) AfterCall(expire time.Duration, f func()) {
	s.timeNodes = append(s.timeNodes, s.timer.AfterFunc(expire, f))
}

func (s *EmbedManagedSession) Close() {
	s.wrapperClose()

	for i, _ := range s.timeNodes {
		s.timeNodes[i].Stop()
	}

	s.implemented.OnClosed()
}
