package base

import (
	"sync"
)

type ConnPool struct {
	sync.RWMutex
	pool map[int64]*Conn
}

func NewConnPool() *ConnPool {
	return &ConnPool{
		pool: map[int64]*Conn{},
	}
}

func (p *ConnPool) Len() int {
	p.RLock()
	defer p.RUnlock()

	return len(p.pool)
}

func (p *ConnPool) GetConn(id int64) (*Conn, bool) {
	p.RLock()
	defer p.RUnlock()

	conn, ok := p.pool[id]
	if !ok {
		return nil, false
	}

	return conn, true
}

func (p *ConnPool) PutConn(conn *Conn) {
	p.Lock()
	defer p.Unlock()

	p.pool[conn.UID] = conn
}

func (p *ConnPool) DelConn(conn *Conn) {
	p.Lock()
	defer p.Unlock()

	delete(p.pool, conn.UID)
}
