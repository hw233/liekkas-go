package service

import "testing"

func TestKoMoe_GetSign(t *testing.T) {
	moe := &KoMoe{}
	sign := moe.getSign(moe.GenSessionVerifyParams("1", 1, true), true)
	t.Log(sign)
}
