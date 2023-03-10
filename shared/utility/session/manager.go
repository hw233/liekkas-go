package session

import (
	"context"
	"sync"
	"time"

	"github.com/antlabs/timer"

	"shared/utility/errors"
)

var (
	ErrManagerIsFull   = errors.New("manager is full")
	ErrSessionNotExist = errors.New("session not exist")
	ErrNotImplManaged  = errors.New("session not impl managed")
)

type Builder interface {
	NewSession() Session
}

type Manager interface {
	GetSession(ctx context.Context, id int64) (Session, error)
	NewSessionIfNotExist(ctx context.Context, id int64) (Session, error)
	GetAllSessionId() []int64
	GetSessionCount() int32
	HasSession(id int64) bool
	Close()
}

type ManagerConfig struct {
	Expire   time.Duration
	Capacity int
}

type manager struct {
	sync.RWMutex
	builder  Builder
	timer    timer.Timer
	sessions map[int64]managedSession
	closed   chan int64 // 异步加锁，和manager容量一致，防止卡死
	close    chan struct{}
	expire   time.Duration // 过期时间
	cap      int           // 容量
}

func NewManager(builder Builder, config *ManagerConfig) *manager {
	m := &manager{
		builder:  builder,
		timer:    timer.NewTimer(),
		sessions: map[int64]managedSession{},
		closed:   make(chan int64, config.Capacity),
		close:    make(chan struct{}, 1),
		expire:   config.Expire,
		cap:      config.Capacity,
	}

	go m.timer.Run()
	go m.watchClosed()

	return m
}

func (m *manager) GetSession(ctx context.Context, id int64) (Session, error) {
	m.Lock()
	defer m.Unlock()

	return m.getSession(ctx, id, false)
}

func (m *manager) NewSessionIfNotExist(ctx context.Context, id int64) (Session, error) {
	m.Lock()
	defer m.Unlock()

	return m.getSession(ctx, id, true)
}

func (m *manager) getSession(ctx context.Context, id int64, allowNil bool) (Session, error) {
	sess, ok := m.sessions[id]
	if !ok {
		if len(m.sessions) >= m.cap {
			return nil, ErrManagerIsFull
		}

		newSess := m.builder.NewSession()

		managedSess, ok := newSess.(managedSession)
		if !ok {
			return nil, ErrNotImplManaged
		}

		err := managedSess.Create(ctx, CreateOpts{
			ID:          id,
			Timer:       m.timer,
			Implemented: newSess,
			WrapperClosed: func() {
				m.closed <- id
			},
			AllowNil: allowNil,
		})
		if err != nil {
			return nil, err
		}

		newSess.ScheduleCall(m.expire, managedSess.Expired)

		m.sessions[id] = managedSess

		return m.sessions[id].(Session), nil
	}

	sess.Trigger(ctx)

	return sess.(Session), nil
}

func (m *manager) GetAllSessionId() []int64 {
	m.RLock()
	defer m.RUnlock()

	ret := make([]int64, 0, len(m.sessions))
	for sessionId := range m.sessions {
		ret = append(ret, sessionId)
	}

	return ret
}

func (m *manager) GetSessionCount() int32 {
	m.RLock()
	defer m.RUnlock()

	return int32(len(m.sessions))
}

func (m *manager) HasSession(id int64) bool {
	m.RLock()
	defer m.RUnlock()

	_, ok := m.sessions[id]
	return ok
}

func (m *manager) watchClosed() {
	for {
		select {
		case id := <-m.closed:
			m.Lock()
			delete(m.sessions, id)
			m.Unlock()
		case <-m.close:
			return
		}
	}
}

func (m *manager) Close() {
	m.Lock()
	defer m.Unlock()

	for _, sess := range m.sessions {
		sess.Close()
	}

	m.timer.Stop()
	m.close <- struct{}{}
}
