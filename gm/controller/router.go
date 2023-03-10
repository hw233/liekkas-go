package controller

import (
	"reflect"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/httputil"
	"shared/utility/naming"
	"shared/utility/router"
	"sync"
)

type Router struct {
	f sync.Map
}

func NewRouter() *Router {
	return &Router{
		f: sync.Map{},
	}
}

type routerOption struct {
	filter func(interface{}) bool
	// prefix interface{}
}

func WithFilter(filter func(interface{}) bool) func(*routerOption) {
	return func(option *routerOption) {
		option.filter = filter
	}
}

func (h *Router) Route(key interface{}) (*router.ReflectFunc, bool) {
	v, ok := h.f.Load(key)
	if !ok {
		return nil, false
	}

	return v.(*router.ReflectFunc), true
}

func (h *Router) RegisterHandler(handler httputil.HttpHandler, options ...func(option *routerOption)) error {
	t := reflect.TypeOf(handler)
	if t.Elem().Kind() != reflect.Struct {
		return errors.New("kind error ")
	}
	// load options
	ro := &routerOption{}
	for _, option := range options {
		option(ro)
	}

	// by traverse handler
	for i := 0; i < t.NumMethod(); i++ {
		k := naming.FirstLower(t.Method(i).Name)
		f := t.Method(i).Func.Interface()
		if ro.filter != nil && !ro.filter(f) {
			// method was filtered
			continue
		}

		h.f.Store(k, router.NewReflectFunc(f))
		glog.Infof("RegisterHandler Method:%s,Path:%s,Func:%v", handler.GetMethod(), k, t.Method(i).Func.String())
	}

	return nil
}
