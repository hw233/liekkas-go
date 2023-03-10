package balancer

import (
	balancer2 "google.golang.org/grpc/balancer"
	"google.golang.org/grpc/resolver"
)

type serverMap struct {
	m map[string]*serverMapEntry
}

func newServerMap() *serverMap {
	return &serverMap{
		m: map[string]*serverMapEntry{},
	}
}

func (s *serverMap) Set(addr resolver.Address, subConn balancer2.SubConn) {
	s.m[addr.ServerName] = &serverMapEntry{
		addr:    addr,
		subConn: subConn,
	}
}

func (s *serverMap) Get(addr resolver.Address) (balancer2.SubConn, bool) {
	entry, ok := s.m[addr.ServerName]
	if !ok {
		return nil, false
	}

	return entry.subConn, true
}

func (s *serverMap) Delete(addr resolver.Address) {
	delete(s.m, addr.ServerName)
}

func (s *serverMap) Len() int {
	return len(s.m)
}

func (s *serverMap) Keys() []resolver.Address {
	addrs := make([]resolver.Address, 0, s.Len())

	for _, v := range s.m {
		addrs = append(addrs, v.addr)
	}

	return addrs
}

type serverMapEntry struct {
	addr    resolver.Address
	subConn balancer2.SubConn
}
