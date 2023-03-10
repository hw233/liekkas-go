package controller

import (
	"encoding/json"
	"time"

	"shared/utility/errors"
	"shared/utility/glog"
)

type httpResp struct {
	Code int         `json:"code"`
	Time int64       `json:"time"`
	Data interface{} `json:"data"`
}

func (hr *httpResp) Encode() []byte {
	ret, _ := json.Marshal(hr)
	return ret
}

func encodeHTTPResp(data interface{}, err error) []byte {
	hr := &httpResp{
		Time: time.Now().Unix(),
	}

	if err != nil {
		err = errors.Format(err)
		glog.Errorf("error: %+v", err)
		hr.Code = errors.Code(err)
		return hr.Encode()
	}

	hr.Data = data

	ret := hr.Encode()

	glog.Debugf("ret: %s", string(ret))

	return ret
}
