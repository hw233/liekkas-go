package method

import (
	"net/http"
	"reflect"
	"shared/utility/httputil"
)

type HttpPostHandler struct {
	*httputil.HttpMethodAdapter
}

func NewHttpPostHandler() *HttpPostHandler {
	return &HttpPostHandler{
		HttpMethodAdapter: httputil.NewHttpMethodAdapter(http.MethodPost),
	}
}
func (p *HttpPostHandler) FilterMethod(f interface{}) bool {
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
