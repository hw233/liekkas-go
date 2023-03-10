package manager

import (
	"context"
	"testing"

	"google.golang.org/grpc/status"

	"gamesvr/config"
	"shared/protobuf/pb"
	utilityConfig "shared/utility/config"
	"shared/utility/errors"
	"shared/utility/etcd"

	"github.com/go-redis/redis/v8"
)

func TestMain(m *testing.M) {
	utilityConfig.SetConfigPath("../config.toml")
	Conf = config.NewDefaultConfig()
	err := utilityConfig.Load(Conf)
	if err != nil {
		return
	}

	EtcdClient, err = etcd.Dial(Conf.ETCDEndpoints)
	if err != nil {
		return
	}

	RedisClient = redis.NewClient(Conf.Redis)
	_, err = RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		return
	}

	err = initRPC()
	if err != nil {
		return
	}

	m.Run()
}

func TestErrorHandle(t *testing.T) {
	_, err := RPCLoginClient.Test(context.Background(), &pb.LoginTestReq{
		Id:  1,
		Msg: "123",
	})

	s, ok := status.FromError(err)
	if !ok {
		t.Errorf("FromError(%+v) !ok", err)
	}

	t.Logf("grpc: code: %d, message: %s", s.Code(), s.Message())

	// err = errors.WrapTrace(err)
	err = errors.WrapText(err, "wwwwwwww")
	err = errors.WrapText(err, "xxxxxx")
	err = errors.WrapTrace(err)
	err = errors.WrapText(err, "zzzzz")
	err = errors.WrapText(err, "2222")
	// err = errors.WrapTrace(err)
	// err = errors.WrapTrace(err)

	t.Logf("errors: code: %d, message: %v", errors.Code(err), errors.Trace(err))
}
