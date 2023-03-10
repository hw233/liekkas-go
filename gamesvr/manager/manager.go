package manager

import (
	"context"
	"database/sql"
	"time"

	"gamesvr/config"
	"shared/common"
	"shared/csv/base"
	"shared/csv/entry"
	"shared/global"
	utilityConfig "shared/utility/config"
	"shared/utility/dblog"
	"shared/utility/etcd"
	"shared/utility/event"
	uglobal "shared/utility/global"
	"shared/utility/glog"
	"shared/utility/lua_state"
	"shared/utility/mysql"
	"shared/utility/servertime"
	"shared/utility/session"

	"github.com/antlabs/timer"
	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	Conf        *config.Config
	CSV         *entry.Manager
	Timer       timer.Timer
	DB          *sql.DB
	MySQL       *mysql.Handler
	EtcdClient  *clientv3.Client
	RedisClient *redis.Client
	// RPCPortalClient *grpc.RPCPortalClient
	// RPCGameServer   *grpc.GameServer
	// RPCGuildClient  *grpc.RPCGuildClient
	// RPCMailClient   *grpc.RPCMailClient

	EventQueue  *event.EventQueue
	SessManager session.Manager
	Global      *global.Global
	LuaState    *lua_state.LuaMatcher
	// RPCLoginClient pb.RPCLoginClient
	Announcements *AnnouncementInfo
)

func Init(builder session.Builder) error {
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

	MySQL = mysql.NewHandler(DB)

	Timer = timer.NewTimer()
	go Timer.Run()

	BIInit()

	RedisClient = redis.NewClient(Conf.Redis)
	_, err = RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		glog.Errorf("connect to redis error: %v", err)
		return err
	}

	SessManager = session.NewManager(builder, &session.ManagerConfig{Expire: time.Minute * 15, Capacity: 1000})

	Global = global.NewGlobal(DB, RedisClient)

	Announcements = NewAnnouncementInfo()
	err = AnnouncementInit()
	if err != nil {
		glog.Error("announcement init failed!!!\n", err.Error())
	}

	EventQueue = event.NewEventQueue(RedisClient)

	LuaState = lua_state.NewLuaMatcher()
	err = LuaState.Init()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	servertime.SetDailyRefreshHour(CSV.GlobalEntry.DailyRefreshHour)
	timeoffset, err := Global.GetInt64(ctx, servertime.TimeOffsetRedisName)
	if err != nil {
		if err != uglobal.ErrNil {
			glog.Errorf("error, get timeoffset failed:%+v\n", err)
		}
		timeoffset = 0
	}
	servertime.SetTimeOffset(timeoffset)

	glog.Infof("now time: %+v", servertime.Now())

	err = initRPC()
	if err != nil {
		glog.Errorf("connect to redis error: %v", err)
		return err
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
