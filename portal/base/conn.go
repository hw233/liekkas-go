package base

import (
	"sync"
	"time"

	"shared/utility/servertime"
	"shared/utility/uuid"

	"github.com/panjf2000/gnet"
	"golang.org/x/time/rate"
)

// 连接上下文
// TODO：做流量控制和其他验证
type Conn struct {
	sync.RWMutex
	*rate.Limiter
	gnet.Conn

	UID         int64  // uid
	Token       string // 连接的token，id+token是连接的唯一标识
	Server      string // 唯一服务器名称，查找addr
	IsConnected bool   // 是否连接
	CTime       int64
}

func NewConn(c gnet.Conn, timeout time.Duration, limiter *rate.Limiter) *Conn {
	return &Conn{
		Limiter:     limiter,
		Conn:        c,
		UID:         0,
		Token:       uuid.GenUUID(),
		Server:      "",
		IsConnected: false,
		CTime:       servertime.Now().Unix(),
	}
}
