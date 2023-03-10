package model

import (
	"context"
	"database/sql"
	"shared/common"
	"time"

	"github.com/go-redis/redis/v8"

	"login/internal/base"
	"shared/utility/key"
)

const (
	insertSQL = "INSERT INTO `login` (`id`,`open_id`,`third_party`,`account`,`password`,`token`) VALUES (?,?,?,?,?,?)"

	selectByIDSQL      = "SELECT `open_id`,`account`,`password`,`token` FROM `login` WHERE `id`=?"
	selectByAccountSQL = "SELECT `id`,`open_id`,`password`,`token` FROM `login` WHERE `account`=?"
	selectByOpenIdSQL  = "SELECT `id`,`account`,`password`,`token` FROM `login` WHERE `open_id`=? and `third_party`=?"

	// 	UPDATE table_name
	// SET column1=value1,column2=value2,...
	// WHERE some_column=some_value;
	updateTokenSQL = "UPDATE `login` SET `token`=? WHERE `id`=?"

	// redis key
	prefixKeyLogin = "login"

	fieldKeyToken = "token"

	tokenExpire = time.Hour * 6

	keyRegisterNum = "registerNum"
)

type LoginAccount struct {
	UserID     int64
	OpenID     string
	ThirdParty string
	Account    string
	Password   string
	Token      string
}

func AddLoginAccount(ctx context.Context, login *LoginAccount) error {
	_, err := base.MySQLClient.ExecContext(ctx, insertSQL, login.UserID, login.OpenID, login.ThirdParty, login.Account, login.Password, login.Token)
	return err
}

func GetLoginAccountByUserID(ctx context.Context, userID int64) (*LoginAccount, error) {
	row := base.MySQLClient.QueryRowContext(ctx, selectByIDSQL, userID)

	login := &LoginAccount{
		UserID: userID,
	}

	err := row.Scan(&login.OpenID, &login.Account, &login.Password, &login.Token)
	if err != nil {
		return nil, err
	}

	return login, nil
}

func GetLoginAccountByAccount(ctx context.Context, account string) (*LoginAccount, error) {
	row := base.MySQLClient.QueryRowContext(ctx, selectByAccountSQL, account)

	login := &LoginAccount{
		Account: account,
	}

	err := row.Scan(&login.UserID, &login.OpenID, &login.Password, &login.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrAccountNotFound
		}
		return nil, err
	}

	return login, nil
}

func GetLoginAccountByOpenIdAndTP(ctx context.Context, openId string, thirdParty string) (*LoginAccount, error) {
	row := base.MySQLClient.QueryRowContext(ctx, selectByOpenIdSQL, openId, thirdParty)

	login := &LoginAccount{
		OpenID:     openId,
		ThirdParty: thirdParty,
	}

	err := row.Scan(&login.UserID, &login.Account, &login.Password, &login.Token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, common.ErrAccountNotFound
		}
		return nil, err
	}

	return login, nil
}

func RefreshToken(ctx context.Context, userID int64, token string) error {
	return base.RedisClient.SetNX(ctx, makeTokenKey(userID), token, tokenExpire).Err()
}

func GetToken(ctx context.Context, userID int64) (string, error) {
	ret, err := base.RedisClient.Get(ctx, makeTokenKey(userID)).Result()
	if err == redis.Nil {
		return "", nil
	}
	return ret, err
}

func makeTokenKey(userID int64) string {
	return key.MakeRedisKey(prefixKeyLogin, fieldKeyToken, userID)
}

func GetRegisterNum(ctx context.Context) (int64, error) {
	ret, err := base.RedisClient.Get(ctx, keyRegisterNum).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return ret, err
}
func IncreaseRegisterNum(ctx context.Context) error {
	return base.RedisClient.Incr(ctx, keyRegisterNum).Err()
}
