package method

import (
	"context"
	"net/http"
	"reflect"
	"shared/utility/glog"
	"shared/utility/httputil"
)

type HttpGetHandler struct {
	*httputil.HttpMethodAdapter
}

func NewHttpGetHandler() *HttpGetHandler {
	return &HttpGetHandler{
		HttpMethodAdapter: httputil.NewHttpMethodAdapter(http.MethodGet),
	}
}

func (g *HttpGetHandler) FilterMethod(f interface{}) bool {
	t := reflect.TypeOf(f)

	// check kind
	if t.Kind() != reflect.Func {
		return false
	}

	// check in
	numIn := t.NumIn()
	if numIn != 3 && numIn != 2 {
		return false
	}
	if t.In(1).String() != "context.Context" {
		return false
	}
	if numIn == 3 {
		_, ok := reflect.New(t.In(2)).Elem().Interface().(*httputil.Params)
		if !ok {
			return false
		}
	}

	numOut := t.NumOut()
	// check out
	if numOut != 2 && numOut != 1 {
		return false
	}

	if numOut == 2 {
		if t.Out(1).String() != "error" {
			return false
		}
	} else {
		if t.Out(0).String() != "error" {
			return false
		}
	}

	return true
}

func (g *HttpGetHandler) Test(ctx context.Context, param *httputil.Params) (string, error) {
	glog.Info(param.Get("A"))
	return param.Encode(), nil
}
