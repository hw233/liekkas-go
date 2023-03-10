package service

import (
	"context"
	"login/internal/base"
	"login/model"
	"shared/common"
	"shared/utility/errors"
	"shared/utility/key"
	"strconv"
)

var koMoe = &KoMoe{
	host:      "https://line3-sdk-adapter.komoejoy.com",
	spareHost: "https://line1-sdk-adapter.komoejoy.com",
	GP: &KoMoeAppConfig{
		merchant_id: 1,
		game_id:     7495,
		secretKey:   "e8a583d6e42b4ccf9e388ba8763b4d9b",
	},
	IOS: &KoMoeAppConfig{
		merchant_id: 1,
		game_id:     7496,
		secretKey:   "0c56137d107543be808f933afeb5b620",
	},
}

func ThirdPartyLoginCheck(thirdParty, openID, token string, ios bool) error {
	switch thirdParty {
	case "koMoe":
		uid, err := strconv.ParseInt(openID, 10, 64)
		if err != nil {
			return errors.WrapTrace(err)
		}

		_, err = koMoe.DoSessionVerify(token, uid, ios)

		if err != nil {
			return errors.WrapTrace(err)
		}
	default:
		return errors.Swrapf(common.ErrUnknownSdk, thirdParty)
	}
	return nil
}

func ThirdPartyRegister(ctx context.Context, thirdParty, openID string) (*model.LoginAccount, error) {
	registerNum, err := model.GetRegisterNum(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}

	//  配置
	if registerNum >= base.CSV.GlobalEntry.CBT1RegisterLimit {
		return nil, common.ErrRegisterLimit
	}
	newUserID, err := base.UserID.GenID(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	account := key.MakeKey('_', thirdParty, openID)

	loginAccount := &model.LoginAccount{
		UserID:     newUserID,
		OpenID:     openID,
		ThirdParty: thirdParty,
		Account:    account,
	}

	err = model.AddLoginAccount(ctx, loginAccount)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	err = model.IncreaseRegisterNum(ctx)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	return loginAccount, nil
}
func ThirdPartyLogin(ctx context.Context, thirdParty, openID, token string, ios bool) (*InnerLoginResp, error) {
	err := ThirdPartyLoginCheck(thirdParty, openID, token, ios)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	loginAccount, err := model.GetLoginAccountByOpenIdAndTP(ctx, openID, thirdParty)
	if err != nil {
		if err != common.ErrAccountNotFound {
			return nil, errors.WrapTrace(err)
		}
		loginAccount, err = ThirdPartyRegister(ctx, thirdParty, openID)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}

	err = model.RefreshToken(ctx, loginAccount.UserID, token)
	if err != nil {
		return nil, err
	}
	return &InnerLoginResp{
		UserID: loginAccount.UserID,
		Token:  token,
	}, nil

}
