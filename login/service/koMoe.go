package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"shared/common"
	"shared/utility/errors"
	"shared/utility/httputil"
	"shared/utility/servertime"
	"sort"
	"strings"
)

var koMoeSessionVerifyUrl = "%s/api/server/session.verify"

// KoMoe 小萌sdk/*
type KoMoe struct {
	host      string          // 请求host
	spareHost string          // 备用host
	GP        *KoMoeAppConfig // google play
	IOS       *KoMoeAppConfig // ios

}
type KoMoeAppConfig struct {
	game_id     int //游戏id，[参数表]中为CP分配的game_id
	merchant_id int //商户id，[参数表]中为CP分配的merchant_id
	secretKey   string
}

func (ko *KoMoe) getSessionVerifyUrl() string {
	return ko.getSessionVerifyUrl_(false)
}

func (ko *KoMoe) getSessionVerifyUrl_(spare bool) string {
	host := ko.host
	if spare {
		host = ko.spareHost
	}
	return fmt.Sprintf(koMoeSessionVerifyUrl, host)
}

type PublicResponse struct {
	Timestamp int64  `json:"timestamp"` // 时间戳(毫秒)，对应request的
	Code      int32  `json:"code"`      // 状态码
	Message   string `json:"message"`   // 错误信息，code不为0的时候出现
}
type SessionVerifyResponse struct {
	*PublicResponse
	OpenId int64  `json:"open_id"` // Game⽤户的⽤户id，也就是客户端登录接⼝返回的uid
	Uname  string `json:"uname"`   // Game⽤户的昵称
}

func (ko *KoMoe) GenSessionVerifyParams(access_key string, uid int64, ios bool) *httputil.Params {
	params := ko.GenPublicParams(ios)
	params.Put("access_key", access_key) // ⽤户登录后身份令牌
	params.Put("uid", fmt.Sprint(uid))   // ⽤户唯⼀ID
	params.Put("region", "7")            // 游戏账号域，region=7
	params.Put("version", "2")           // 接⼝版本,version=2

	return params
}

// GenPublicParams 请求消息公共字段(公共请求参数
func (ko *KoMoe) GenPublicParams(ios bool) *httputil.Params {
	params := httputil.EmptyParams()
	var game_Id int
	var merchant_id int

	if ios {
		game_Id = ko.IOS.game_id
		merchant_id = ko.IOS.merchant_id

	} else {
		game_Id = ko.GP.game_id
		merchant_id = ko.GP.merchant_id

	}
	params.Put("game_id", fmt.Sprint(game_Id))                           // 游戏id，[参数表]中为CP分配的game_id
	params.Put("merchant_id", fmt.Sprint(merchant_id))                   // 商户id，[参数表]中为CP分配的merchant_id
	params.Put("timestamp", fmt.Sprint(servertime.Now().UnixNano()/1e6)) // 当前时间戳(毫秒)
	params.Put("sign", "")                                               // 签名 ,后续生成，先放个位置
	return params
}

// 对所有字段进行签名
func (ko *KoMoe) getSign(params *httputil.Params, ios bool) string {
	var buf strings.Builder
	values := params.Values()
	keys := make([]string, 0, len(values))
	for k := range values {
		keys = append(keys, k)
	}
	// 对keys 排序
	sort.Strings(keys)
	for _, k := range keys {
		//跳过不需要签名的字段

		if "sign" == k {
			continue
		}
		buf.WriteString(params.Get(k).String())
	}
	var secretKey string
	if ios {
		secretKey = ko.IOS.secretKey
	} else {
		secretKey = ko.GP.secretKey

	}
	buf.WriteString(secretKey)
	s := buf.String()
	// md5hex&lower
	md5ctx := md5.New()
	md5ctx.Write([]byte(s))
	return strings.ToLower(hex.EncodeToString(md5ctx.Sum(nil)))
}

func (ko *KoMoe) DoSessionVerify(access_key string, uid int64, ios bool) (*SessionVerifyResponse, error) {
	params := ko.GenSessionVerifyParams(access_key, uid, ios)

	bytes, err := ko.DoPost(ko.getSessionVerifyUrl(), params, ios)

	if err != nil {
		// 尝试备用路线
		bytes, err = ko.DoPost(ko.getSessionVerifyUrl_(true), params, ios)
		if err != nil {
			return nil, errors.WrapTrace(err)
		}
	}
	ret := &SessionVerifyResponse{}
	err = json.Unmarshal(bytes, ret)
	if err != nil {
		return nil, errors.WrapTrace(err)
	}
	if ret.Code != 0 || ret.Message != "" {
		return nil, errors.Swrapf(common.ErrLoginSdkError, "koMoe", ret.Code, ret.Message)
	}
	return ret, nil
}
func (ko *KoMoe) DoPost(url string, params *httputil.Params, ios bool) ([]byte, error) {
	params.Put("sign", ko.getSign(params, ios))
	// header User-Agent Y Mozilla/5.0 GameServer 常量，所有API都是这个值
	return httputil.DoPostUrl(url, params, httputil.NewHeaderField("User-Agent", "Mozilla/5.0 GameServer"))

}
