package service

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"shared/common"
	"time"

	"login/internal/base"
	"login/model"
)

func Sha1Hash(data string) string {
	hash := sha1.New()
	_, _ = io.WriteString(hash, data)
	return fmt.Sprintf("%x", hash.Sum(nil))
}

type InnerRegisterResp struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

func InnerRegister(ctx context.Context, account, password string) (*InnerRegisterResp, error) {
	if !base.IsTestServerENV() {
		return nil, errors.New("not test ENV")
	}
	_, err := model.GetLoginAccountByAccount(ctx, account)
	if err != nil {
		if err != common.ErrAccountNotFound {
			return nil, err
		}
	}
	if err == nil {
		return nil, common.ErrAccountRepeat
	}

	password = Sha1Hash(account + password)
	token := Sha1Hash(time.Now().String() + password)

	newUserID, err := base.UserID.GenID(ctx)
	if err != nil {
		return nil, err
	}

	login := &model.LoginAccount{
		UserID:     newUserID,
		ThirdParty: "inner",
		OpenID:     fmt.Sprint(newUserID),
		Account:    account,
		Password:   password,
		Token:      token,
	}

	err = model.AddLoginAccount(ctx, login)
	if err != nil {
		return nil, err
	}
	err = model.IncreaseRegisterNum(ctx)
	if err != nil {
		return nil, err
	}
	return &InnerRegisterResp{
		UserID: newUserID,
		Token:  token,
	}, nil
}

type InnerLoginResp struct {
	UserID int64  `json:"user_id"`
	Token  string `json:"token"`
}

func InnerLogin(ctx context.Context, account, password string) (*InnerLoginResp, error) {
	if !base.IsTestServerENV() {
		return nil, errors.New("not test ENV")
	}

	login, err := model.GetLoginAccountByAccount(ctx, account)
	if err != nil {
		return nil, err
	}

	if login.Password != Sha1Hash(account+password) {
		return nil, errors.New("password wrong")
	}

	// 重新生成token
	token := Sha1Hash(time.Now().String() + password)

	// 刷新token
	err = model.RefreshToken(ctx, login.UserID, token)
	if err != nil {
		return nil, err
	}

	login.Token = token

	return &InnerLoginResp{
		UserID: login.UserID,
		Token:  login.Token,
	}, nil
}
