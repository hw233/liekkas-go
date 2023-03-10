package controller

import (
	"encoding/json"
	"gm/method"
	"io/ioutil"
	"net/http"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/httputil"
	"shared/utility/safe"
	"sync"
)

type Dispatcher struct {
	f       sync.Map
	filters []HttpRequestFilter
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer safe.Recover()
	glog.Debugf("req :%+v", r)

	//接入规范 响应头 Content-Type 必须 是： application/json
	w.Header().Set("content-type", "application/json")
	err := d.Filter(r)
	if err != nil {
		glog.Errorf("ServeHTTP Filter err:%v", err)
		w.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return

	}
	path := r.URL.Path[1:]
	moduleRouter, ok := d.Dispatch(r.Method)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	var content *httputil.HttpContent
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		glog.Debugf("req get path:%s ,param :%v", path, query)
		content = moduleRouter.HandleGet(r.Context(), path, query)
	case http.MethodPost:
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			glog.Errorf("ServeHTTP Parse body err:%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		glog.Debugf("req post path:%s,body :%s", path, bytes)
		content = moduleRouter.HandlePost(r.Context(), path, bytes)
	default:
		content = httputil.NewHttpContent(nil, errors.NewCode(404, "404"))
	}

	if content.ErrorCode != 0 {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	marshal, err := json.Marshal(content)
	if err != nil {
		glog.Errorf("ServeHTTP json Marshal error: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		glog.Errorf("ServeHTTP Write Err:%v", err)
	}

}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		f: sync.Map{},
	}
}

func (d *Dispatcher) RegisterRouter(router *PathRouter) {
	d.f.Store(router.GetMethod(), router)

}

func (d *Dispatcher) AddFilter(filter HttpRequestFilter) {
	d.filters = append(d.filters, filter)
}
func (d *Dispatcher) Filter(r *http.Request) error {
	for _, filter := range d.filters {
		err := filter.filter(r)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}
	return nil
}

func (d *Dispatcher) RegisterRouters() error {
	get, err := NewPathRouter(method.NewHttpGetHandler())
	if err != nil {
		glog.Fatalf("NewPathRouter Method GET error: %v", err)
		return err
	}
	d.RegisterRouter(get)
	post, err := NewPathRouter(method.NewHttpPostHandler())
	if err != nil {
		glog.Fatalf("NewPathRouter  Method POST error: %v", err)
		return err

	}
	d.RegisterRouter(post)
	return nil
}

func (d *Dispatcher) Dispatch(method string) (*PathRouter, bool) {

	v, ok := d.f.Load(method)
	if !ok {
		return nil, false
	}
	return v.(*PathRouter), true

}
