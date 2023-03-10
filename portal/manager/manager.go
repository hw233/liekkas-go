package manager

import (
	"github.com/go-redis/redis/v8"

	"portal/base"
	"portal/config"
	sconfig "shared/utility/config"
	"shared/utility/etcd"
	"shared/utility/glog"

	"github.com/panjf2000/ants/v2"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	Conf        *config.Config
	GoPool      *ants.Pool
	ConnPool    *base.ConnPool
	EtcdClient  *clientv3.Client
	RedisClient *redis.Client
	// RPCGameClient   *grpc.GameClient
	// RPCPortalServer *grpc.PortalServer
)

func Init() error {
	Conf = config.NewDefaultConfig()
	err := sconfig.Load(Conf)
	if err != nil {
		glog.Errorf("load config error: %v", err)
		return err
	}

	err = glog.InitLog(Conf.Log)
	if err != nil {
		glog.Fatalf("SetLogLevel error: %v", err)
		return err
	}

	EtcdClient, err = etcd.Dial(Conf.ETCDEndpoints)
	if err != nil {
		glog.Errorf("etcd dial error: %v", err)
		return err
	}

	GoPool, err = ants.NewPool(100)
	if err != nil {
		glog.Errorf("new pool error: %v", err)
		return err
	}

	GoPool.Tune(10000)

	ConnPool = base.NewConnPool()

	RedisClient = redis.NewClient(Conf.Redis)
	_, err = RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		glog.Errorf("connect to redis error: %v", err)
		return err
	}

	err = initRPC()
	if err != nil {
		return err
	}
	// RPCGameClient = grpc.NewGameClient(EtcdClient, RedisClient, common.ServiceGame)
	// err = RPCGameClient.Dial()
	// if err != nil {
	// 	return err
	// }

	// RPCPortalServer = grpc.NewPortalServer(EtcdClient, RedisClient, common.ServicePortal, Conf.ServerName)

	return nil
}

func Reload() error {
	newConf := config.NewDefaultConfig()

	err := sconfig.Reload(newConf)
	if err != nil {
		return err
	}

	Conf = newConf

	return nil
}
