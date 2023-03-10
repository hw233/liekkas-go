package manager

import (
	"database/sql"
	"mailserver/config"
	"shared/common"
	"shared/csv/base"
	"shared/csv/entry"
	"shared/grpc"
	"shared/protobuf/pb"
	utilityConfig "shared/utility/config"
	"shared/utility/dblog"
	"shared/utility/etcd"
	"shared/utility/global"
	"shared/utility/glog"
	"shared/utility/mysql"

	"github.com/go-redis/redis/v8"
	"github.com/openzipkin/zipkin-go/reporter"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	Conf        *config.Config
	CSV         *entry.Manager
	DB          *sql.DB
	EtcdClient  *clientv3.Client
	RedisClient *redis.Client
	Global      *global.Global

	RPCServer     *grpc.Server
	RPCClient     *grpc.Client
	RPCGameClient pb.GameClient

	ZipkinReporterClientGame reporter.Reporter
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
	mysql.SetLogger(dblog.NewDBLogger())

	csvBase := &base.ConfigManager{}
	csvBase.Init()
	csvBase.LoadConfig(Conf.CSVPath)
	CSV = entry.NewManager()
	err = CSV.Reload(csvBase)
	if err != nil {
		glog.Errorf("load config error: %v", err)
		return err
	}

	// 设置奖励类型自动填充对
	common.SetRewardsType(CSV.Reward.RewardsType())

	EtcdClient, err = etcd.Dial(Conf.ETCDEndpoints)
	if err != nil {
		glog.Errorf("etcd connect error: %v", err)
		return err
	}

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

	RedisClient = redis.NewClient(Conf.Redis)
	_, err = RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		glog.Errorf("connect to redis error: %v", err)
		return err
	}

	Global = global.NewGlobal(RedisClient)

	RPCClient = grpc.NewClient(EtcdClient, RedisClient)
	RPCServer = grpc.NewServer(Conf.Service, Conf.ServerName, EtcdClient, RedisClient)

	gameClientConn, gameReporter, err := RPCClient.Dial(common.ServiceGame, Conf.ServerName)
	ZipkinReporterClientGame = gameReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	gameClientConn.RegisterFetcherFunc(func(args interface{}) int64 {
		if f, ok := args.(interface {
			GetUid() int64
		}); ok {
			return f.GetUid()
		}

		return 0
	})

	gameClientConn.RegisterBroadcastMethod("/Game/NewGroupMailNotify")
	gameClientConn.IgnoreNullServer()

	RPCGameClient = pb.NewGameClient(gameClientConn)

	return nil
}

func ReporterClose() {
	ZipkinReporterClientGame.Close()
}
