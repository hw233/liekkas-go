package httputil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"shared/utility/errors"
	"shared/utility/glog"
)

//DoPostUrl  参数放在url里面
func DoPostUrl(url string, params *Params, fields ...*HeaderField) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s?%s", url, params.Encode()), nil)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	for _, field := range fields {
		req.Header.Add(field.K, field.V)
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	glog.Debugf("DoPostUrl resp body:%s", bytes)
	return bytes, nil
}

//DoPostForm  参数放在form里面
func DoPostForm(url string, values url.Values) ([]byte, error) {
	resp, err := http.PostForm(http.MethodPost, values)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	return bytes, nil
}
