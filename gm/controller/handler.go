package controller

import (
	"context"
	"encoding/json"
	"net/url"
	"reflect"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/httputil"
)

type PathRouter struct {
	*Router
	httputil.HttpHandler
}

func NewPathRouter(handler httputil.HttpHandler) (*PathRouter, error) {
	t := reflect.TypeOf(handler)
	if t.Elem().Kind() != reflect.Struct {
		return nil, errors.New("kind error ")
	}
	newRouter := NewRouter()
	err := newRouter.RegisterHandler(handler, WithFilter(handler.FilterMethod))
	if err != nil {
		return nil, err
	}
	return &PathRouter{
		Router:      newRouter,
		HttpHandler: handler,
	}, nil
}

func (h *PathRouter) HandleGet(ctx context.Context, path string, values url.Values) *httputil.HttpContent {

	reflectFunc, ok := h.Route(path)
	if !ok {
		glog.Errorf("Handle Get Route not found,path: %s", path)
		return httputil.NewHttpContent(nil, errors.New("Route not found"))
	}
	var ret []interface{}
	if reflectFunc.NumIn() == 3 {
		ret = reflectFunc.Call(h.HttpHandler, ctx, httputil.NewParams(values))
	} else {
		ret = reflectFunc.Call(h.HttpHandler, ctx)
	}

	if len(ret) == 2 {
		err, ok := ret[1].(error)
		if ok && err != nil {
			glog.Errorf("Handle Call return error: %+v", errors.Format(err))
		}

		content := httputil.NewHttpContent(ret[0], err)
		glog.Debugf("resp :%v", content)
		return content

	} else {
		err, ok := ret[0].(error)
		if ok && err != nil {
			glog.Errorf("Handle Call return error: %+v", errors.Format(err))
		}

		content := httputil.NewHttpContent(nil, err)
		glog.Debugf("resp :%v", content)
		return content
	}
}

func (h *PathRouter) HandlePost(ctx context.Context, path string, bytes []byte) *httputil.HttpContent {

	reflectFunc, ok := h.Route(path)
	if !ok {
		glog.Errorf("Handle Post Route not found,path: %s", path)
		return httputil.NewHttpContent(nil, errors.New("Route not found"))
	}
	in2 := reflectFunc.In(2)
	err := json.Unmarshal(bytes, in2)
	if err != nil {
		glog.Errorf("Handle Call return error: %+v", errors.Format(err))
		content := httputil.NewHttpContent(nil, err)
		glog.Debugf("resp :%v", content)
		return content
	}

	ret := reflectFunc.Call(h.HttpHandler, ctx, in2)

	if len(ret) == 2 {
		err, ok := ret[1].(error)
		if ok && err != nil {
			glog.Errorf("Handle Call return error: %+v", errors.Format(err))
		}

		content := httputil.NewHttpContent(ret[0], err)
		glog.Debugf("resp :%v", content)
		return content

	} else {
		err, ok := ret[0].(error)
		if ok && err != nil {
			glog.Errorf("Handle Call return error: %+v", errors.Format(err))
		}

		content := httputil.NewHttpContent(nil, err)
		glog.Debugf("resp :%v", content)
		return content
	}
}
