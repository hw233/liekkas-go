package manager

import (
	"database/sql"
	"gm/config"
	"shared/common"
	"shared/csv/base"
	"shared/csv/entry"
	"shared/global"
	utilityConfig "shared/utility/config"
	"shared/utility/dblog"
	"shared/utility/etcd"
	"shared/utility/glog"
	"shared/utility/mysql"

	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	Conf        *config.Config
	CSV         *entry.Manager
	EtcdClient  *clientv3.Client
	RedisClient *redis.Client
	DB          *sql.DB
	MySQL       *mysql.Handler
	Global      *global.Global
)

func Init() error {
	Conf = config.NewDefaultConfig()
	err := utilityConfig.Load(Conf)
	if err != nil {
		glog.Fatalf("load config error: %v", err)
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
		glog.Fatalf("load config error: %v", err)
		return err
	}

	// 设置奖励类型自动填充对
	common.SetRewardsType(CSV.Reward.RewardsType())

	EtcdClient, err = etcd.Dial(Conf.ETCDEndpoints)
	if err != nil {
		glog.Fatalf("etcd connect error: %v", err)
		return err
	}
	RedisClient = redis.NewClient(Conf.Redis)
	_, err = RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		glog.Fatalf("connect to redis error: %v", err)
		return err
	}

	DB, err = sql.Open("mysql", Conf.MySQL.Addr)
	if err != nil {
		glog.Fatalf("mysql sql Open error: %v", err)

		return err
	}

	DB.SetMaxIdleConns(Conf.MySQL.MaxIdleConn)
	DB.SetMaxOpenConns(Conf.MySQL.MaxOpenConn)
	DB.SetConnMaxLifetime(Conf.MySQL.ConnMaxLifetime)
	err = DB.Ping()
	if err != nil {
		glog.Fatalf("init mysql error: %v", err)
		return err
	}

	MySQL = mysql.NewHandler(DB)
	Global = global.NewGlobal(DB, RedisClient)

	err = initRPC()
	if err != nil {
		glog.Fatalf("initRPC error: %v", err)
		return err
	}
	return nil
}
