package manager

import (
	"database/sql"
	"log"
	"shared/protobuf/pb"
	"shared/utility/event"
	"time"

	"guild/config"
	"shared/common"
	"shared/csv/base"
	"shared/csv/entry"
	"shared/global"
	"shared/grpc"
	utilityConfig "shared/utility/config"
	"shared/utility/etcd"
	"shared/utility/glog"
	"shared/utility/mysql"
	utilitySession "shared/utility/session"

	"github.com/antlabs/timer"
	"github.com/go-redis/redis/v8"
	"github.com/openzipkin/zipkin-go/reporter"
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
	SessManager utilitySession.Manager
	Global      *global.Global
	EventQueue  *event.EventQueue

	RPCClient     *grpc.Client
	RPCServer     *grpc.Server
	RPCMailClient pb.MailClient

	ZipkinReporterClientMail reporter.Reporter
)

func Init(builder utilitySession.Builder) error {
	Conf = config.NewDefaultConfig()
	err := utilityConfig.Load(Conf)
	if err != nil {
		log.Printf("load config error: %v", err)
		return err
	}

	err = glog.InitLog(Conf.Log)
	if err != nil {
		glog.Fatalf("SetLogLevel error: %v", err)
		return err
	}

	csvBase := &base.ConfigManager{}
	csvBase.Init()
	csvBase.LoadConfig(Conf.CSVPath)
	CSV = entry.NewManager()
	err = CSV.Reload(csvBase)
	if err != nil {
		log.Printf("load config error: %v", err)
		return err
	}

	// 设置奖励类型自动填充对
	common.SetRewardsType(CSV.Reward.RewardsType())

	EtcdClient, err = etcd.Dial(Conf.ETCDEndpoints)
	if err != nil {
		log.Printf("etcd connect error: %v", err)
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

	RedisClient = redis.NewClient(Conf.Redis)
	_, err = RedisClient.Ping(RedisClient.Context()).Result()
	if err != nil {
		log.Printf("connect to redis error: %v", err)
		return err
	}

	SessManager = utilitySession.NewManager(builder, &utilitySession.ManagerConfig{Expire: time.Minute * 3, Capacity: 1000})

	Global = global.NewGlobal(DB, RedisClient)
	EventQueue = event.NewEventQueue(RedisClient)

	RPCClient = grpc.NewClient(EtcdClient, RedisClient)
	RPCServer = grpc.NewServer(Conf.Service, Conf.ServerName, EtcdClient, RedisClient)

	mailClientConn, mailReporter, err := RPCClient.Dial(common.ServiceMail, Conf.ServerName)
	ZipkinReporterClientMail = mailReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	RPCMailClient = pb.NewMailClient(mailClientConn)

	return nil
}

func ReporterClose() {
	ZipkinReporterClientMail.Close()
}
