package controller

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestThirdPartyLogin(t *testing.T) {
	resp, err := http.PostForm("http://127.0.0.1:8091/v1/third-party-login",
		url.Values{"third_party": {"koMoe"}, "open_id": {"1234"}, "token": {"123"}, "ios": {"true"}})
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))

}

func TestInnerRegister(t *testing.T) {
	resp, err := http.PostForm("http://127.0.0.1:8091/v1/inner-register",
		url.Values{"account": {"user1"}, "password": {"password1"}})
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}

func TestInnerLogin(t *testing.T) {
	resp, err := http.PostForm("http://127.0.0.1:8091/v1/inner-login",
		url.Values{"account": {"user1"}, "password": {"password1"}})
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
	}

	fmt.Println(string(body))
}
