package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

const (
	httpPort = "8091"

	apiInnerRegister = "/v1/inner-register"

	account  = "xuezhe"
	password = "123456"
)

type InnerRegisterResp struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

type InnerLoginResp struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

func TestLoginRegister(t *testing.T) {
	client := http.Client{}

	vals := url.Values{}
	vals.Add("account", account)
	vals.Add("password", password)

	ret, err := client.PostForm("http://127.0.0.1:"+httpPort+apiInnerRegister, vals)
	if err != nil {
		t.Errorf("client.PostForm(%s) error: %v", apiInnerRegister, err)
		return
	}

	content, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		t.Errorf("ioutil.ReadAll(%s) error: %v", apiInnerRegister, err)
		return
	}

	resp := &InnerRegisterResp{}

	err = json.Unmarshal(content, resp)
	if err != nil {
		t.Errorf("json.Unmarshal(%s) error: %v", string(content), err)
		return
	}

	t.Logf("register content: %s resp: %+v", string(content), resp)
}
