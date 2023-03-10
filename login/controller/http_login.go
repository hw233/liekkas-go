package controller

import (
	"net/http"
	"strconv"

	"login/service"
)

func InnerRegister(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {

	}

	account := req.Form.Get("account")
	if account == "" {
		_, _ = resp.Write(encodeHTTPResp("", nil))
		return
	}

	password := req.Form.Get("password")
	if password == "" {
		_, _ = resp.Write(encodeHTTPResp("", nil))
		return
	}

	ret, err := service.InnerRegister(req.Context(), account, password)
	_, _ = resp.Write(encodeHTTPResp(ret, err))
}

func InnerLogin(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {

	}

	account := req.Form.Get("account")
	if account == "" {
		_, _ = resp.Write(encodeHTTPResp("", nil))
		return
	}

	password := req.Form.Get("password")
	if password == "" {
		_, _ = resp.Write(encodeHTTPResp("", nil))
		return
	}

	ret, err := service.InnerLogin(req.Context(), account, password)
	_, _ = resp.Write(encodeHTTPResp(ret, err))
}

func ThirdPartyLogin(resp http.ResponseWriter, req *http.Request) {
	err := req.ParseForm()
	if err != nil {

	}
	thirdParty := req.Form.Get("third_party")
	if thirdParty == "" {
		_, _ = resp.Write(encodeHTTPResp("", nil))
		return
	}
	openID := req.Form.Get("open_id")
	if openID == "" {
		_, _ = resp.Write(encodeHTTPResp("", nil))
		return
	}
	token := req.Form.Get("token")
	if token == "" {
		_, _ = resp.Write(encodeHTTPResp("", nil))
		return
	}
	iosStr := req.Form.Get("ios")
	if iosStr == "" {
		_, _ = resp.Write(encodeHTTPResp("", nil))
		return
	}
	ios, err := strconv.ParseBool(iosStr)
	if err != nil {
		_, _ = resp.Write(encodeHTTPResp("", err))
		return
	}
	ret, err := service.ThirdPartyLogin(req.Context(), thirdParty, openID, token, ios)
	_, _ = resp.Write(encodeHTTPResp(ret, err))

}
