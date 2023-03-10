package manager

import (
	"database/sql"
	"foreplay/config"
	"foreplay/csv/base"
	"foreplay/csv/entry"
	"shared/global"
	"shared/grpc"
	utilityConfig "shared/utility/config"
	"shared/utility/etcd"
	"shared/utility/glog"

	"github.com/antlabs/timer"
	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	Conf          *config.Config
	CSV           *entry.Manager
	DB            *sql.DB
	EtcdClient    *clientv3.Client
	RedisClient   *redis.Client
	Timer         timer.Timer
	Announcements *AnnouncementInfo
	Server        *ServerInfo
	RPCClient     *grpc.Client
	RPCServer     *grpc.Server
	Global        *global.Global
)

func Init() error {
	Conf = config.NewDefaultConfig()
	err := utilityConfig.Load(Conf)
	if err != nil {
		glog.Errorf("load config error: %v", err)
		return err
	}

	err = glog.InitLog(Conf.Log)
	if err != nil {
		glog.Fatalf("SetLogLevel error: %v", err)
		return err
	}

	baseCfgMgr := &base.ConfigManager{}
	baseCfgMgr.Init()
	baseCfgMgr.LoadConfig(Conf.CSVPath)
	CSV = entry.NewManager()
	err = CSV.Reload(baseCfgMgr)
	if err != nil {
		glog.Errorf("load config error: %v", err)
		return err
	}

	Timer = timer.NewTimer()
	go Timer.Run()

	DB, err = sql.Open("mysql", Conf.MySQL.Addr)
	if err != nil {
		return err
	}
	DB.SetMaxIdleConns(Conf.MySQL.MaxIdleConn)
	DB.SetMaxOpenConns(Conf.MySQL.MaxOpenConn)
	DB.SetConnMaxLifetime(Conf.MySQL.ConnMaxLifetime)
	err = DB.Ping()
	if err != nil {
		glog.Errorf("init mysql error: %v", err)
		return err
	}

	EtcdClient, err = etcd.Dial(Conf.ETCDEndpoints)
	if err != nil {
		glog.Errorf("etcd connect error: %v", err)
		return err
	}

	RedisClient = redis.NewClient(Conf.Redis)
	_, err = RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		glog.Errorf("connect to redis error: %v", err)
		return err
	}

	Global = global.NewGlobal(DB, RedisClient)

	RPCClient = grpc.NewClient(EtcdClient, RedisClient)
	RPCServer = grpc.NewServer(Conf.Service, Conf.ServerName, EtcdClient, RedisClient)

	Announcements = NewAnnouncementInfo()
	err = Announcements.Init()
	if err != nil {
		glog.Error("announcement init failed!!!\n", err.Error())
	}

	Server = NewServerInfo()
	err = Server.Init()
	if err != nil {
		glog.Error("announcement init failed!!!\n", err.Error())
	}

	return nil
}

func Reload() error {
	newConf := config.NewDefaultConfig()

	err := utilityConfig.Reload(newConf)
	if err != nil {
		return err
	}

	Conf = newConf

	return nil
}
