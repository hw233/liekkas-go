package grpc

import "time"

type serverCache struct {
	server string
	expire time.Time
}

func newServerCache(server string, expire time.Duration) *serverCache {
	return &serverCache{
		server: server,
		expire: time.Now().Add(expire),
	}
}

func (sc *serverCache) IsExpired() bool {
	return time.Now().After(sc.expire)
}
