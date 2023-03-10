package router

import (
	"errors"
	"reflect"
	"sync"
)

var typeInvalidErr = errors.New("invalid type for register function")

// Router: routing from command to function
type Router struct {
	f sync.Map
}

func NewRouter() *Router {
	return &Router{
		f: sync.Map{},
	}
}

type routerOption struct {
	config map[int32]string // key: command, val: function name
	// filter interface{}            // function type

	filter func(interface{}) bool
	// prefix interface{}
}

func defaultRouterOption() *routerOption {
	return &routerOption{
		config: nil,
		filter: nil,
	}
}

func WithConfig(config map[int32]string) func(*routerOption) {
	return func(option *routerOption) {
		option.config = config
	}
}

func WithFilter(filter func(interface{}) bool) func(*routerOption) {
	return func(option *routerOption) {
		option.filter = filter
	}
}

func (h *Router) Route(key interface{}) (*ReflectFunc, bool) {
	v, ok := h.f.Load(key)
	if !ok {
		return nil, false
	}

	return v.(*ReflectFunc), true
}

func (h *Router) RouteCall(key interface{}, in ...interface{}) ([]interface{}, bool) {
	f, ok := h.Route(key)
	if !ok {
		return nil, false
	}

	return f.Call(in...), ok
}

func (h *Router) Register(k int32, f interface{}) error {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return typeInvalidErr
	}

	h.f.Store(k, NewReflectFunc(f))
	return nil
}

func (h *Router) RegisterHandler(handler interface{}, options ...func(option *routerOption)) error {
	t := reflect.TypeOf(handler)
	if t.Elem().Kind() != reflect.Struct {
		return typeInvalidErr
	}

	// load options
	ro := &routerOption{}
	for _, option := range options {
		option(ro)
	}

	if ro.config != nil {
		// by config
		for k, v := range ro.config {
			m, ok := t.MethodByName(v)
			if !ok {
				// not found method
				continue
			}

			f := m.Func.Interface()

			if ro.filter != nil && !ro.filter(f) {
				// method was filtered
				continue
			}

			h.f.Store(k, NewReflectFunc(f))
		}
	} else {
		// by traverse handler
		for i := 0; i < t.NumMethod(); i++ {
			f := t.Method(i).Func.Interface()

			if ro.filter != nil && !ro.filter(f) {
				// method was filtered
				continue
			}

			h.f.Store(i, NewReflectFunc(f))
		}
	}

	return nil
}
