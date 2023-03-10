package whitelist

import "sync"

type WhiteListType string

const (
	Id WhiteListType = "id"
)

type Options []*Option

func EmptyOps() *Options {
	return (*Options)(&[]*Option{})
}
func (i *Options) With(T WhiteListType, V interface{}) *Options {
	*i = append(*i, newOption(T, V))
	return i
}
func NewOpsWithUid(uid int64) *Options {
	return EmptyOps().With(Id, uid)
}

type Option struct {
	T WhiteListType
	V interface{}
}

func newOption(T WhiteListType, V interface{}) *Option {
	return &Option{
		T: T,
		V: V,
	}
}

type MultiWhiteList struct {
	sync.RWMutex
	Interceptors map[WhiteListType]*interceptor
}

func NewMultiWhiteList() *MultiWhiteList {
	return &MultiWhiteList{
		Interceptors: map[WhiteListType]*interceptor{},
	}
}

func (m *MultiWhiteList) Reload(Options *Options) {
	m.Lock()
	defer m.Unlock()
	m.Interceptors = map[WhiteListType]*interceptor{}
	for _, op := range *Options {
		i, ok := m.Interceptors[op.T]
		if !ok {
			i = NewInterceptor()
			m.Interceptors[op.T] = i
		}
		i.add(op.V)
	}

}

func (m *MultiWhiteList) Add(Options *Options) {
	m.Lock()
	defer m.Unlock()
	for _, op := range *Options {
		i, ok := m.Interceptors[op.T]
		if !ok {
			i = NewInterceptor()
			m.Interceptors[op.T] = i
		}
		i.add(op.V)
	}

}

func (m *MultiWhiteList) Del(Options *Options) {
	m.Lock()
	defer m.Unlock()
	for _, op := range *Options {
		i, ok := m.Interceptors[op.T]
		if !ok {
			continue
		}
		i.del(op.V)
	}

}

func (m *MultiWhiteList) Filter(Options *Options) bool {
	m.RLock()
	defer m.RUnlock()
	for _, op := range *Options {
		i, ok := m.Interceptors[op.T]
		if !ok {
			return false
		}
		if !i.try(op.V) {
			return false
		}
	}
	return true
}

type interceptor map[interface{}]struct{}

func NewInterceptor() *interceptor {
	return (*interceptor)(&map[interface{}]struct{}{})
}

func (w *interceptor) del(value interface{}) {
	delete(*w, value)
}

func (w *interceptor) add(value interface{}) {
	(*w)[value] = struct{}{}
}

func (w *interceptor) try(value interface{}) bool {
	_, ok := (*w)[value]
	return ok
}
