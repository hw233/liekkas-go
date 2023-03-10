package base

import (
	"database/sql"
	"shared/common"
	"shared/csv/base"
	"shared/csv/entry"

	"login/config"
	"shared/global"
	"shared/grpc"
	utilityConfig "shared/utility/config"
	"shared/utility/etcd"
	"shared/utility/glog"

	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	Config *config.Config

	MySQLClient *sql.DB
	RedisClient *redis.Client
	EtcdClient  *clientv3.Client

	RPCServer *grpc.Server
	UserID    *global.UserID
	CSV       *entry.Manager
)

func Init() {
	initConfig()
	initMySQLClient()
	initRedisClient()
	initEtcdClient()
	initRPC()
	initUserID()
	initCSV()
}

func initConfig() {
	Config = config.NewDefaultConfig()

	err := utilityConfig.Load(Config)
	if err != nil {
		glog.Fatalf("load config error: %v", err)
	}

	glog.Infof("config: %+v", Config)
}

func initMySQLClient() {
	mysqlConf := Config.MySQL
	var err error
	MySQLClient, err = sql.Open("mysql", mysqlConf.Addr)
	if err != nil {
		glog.Fatalf("mysql[%s] open error: %v", mysqlConf.Addr, err)
	}

	MySQLClient.SetMaxIdleConns(mysqlConf.MaxIdleConn)
	MySQLClient.SetMaxOpenConns(mysqlConf.MaxOpenConn)
	MySQLClient.SetConnMaxLifetime(mysqlConf.ConnMaxLifetime)

	// MySQLHandler = mysql.NewHandler(MySQLDB)
}

func initRedisClient() {
	RedisClient = redis.NewClient(Config.Redis)
	_, err := RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		glog.Fatalf("redis[%+v] ping error: %v", Config.Redis, err)
	}
}

func initEtcdClient() {
	var err error
	EtcdClient, err = etcd.Dial(Config.ETCDEndpoints)
	if err != nil {
		glog.Fatalf("etcd[%+v] dial error: %v", Config.ETCDEndpoints, err)
	}
}

func initRPC() {
	RPCServer = grpc.NewServer(Config.Service, Config.ServerName, EtcdClient, RedisClient)
}

func initUserID() {
	UserID = global.NewUserID(RedisClient)
}
func initCSV() {
	csvBase := &base.ConfigManager{}
	csvBase.Init()
	csvBase.LoadConfig(Config.CSVPath)
	CSV = entry.NewManager()
	err := CSV.Reload(csvBase)
	if err != nil {
		glog.Fatalf("load config error: %v", err)
	}

	// 设置奖励类型自动填充对
	common.SetRewardsType(CSV.Reward.RewardsType())
}
