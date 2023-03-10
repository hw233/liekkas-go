package model

import (
	"context"
	"database/sql"
	"gamesvr/manager"
	"log"
	"shared/common"
	"shared/csv/base"
	"shared/csv/entry"
	"shared/global"
	"shared/utility/dblog"
	"shared/utility/glog"
	"shared/utility/lua_state"
	"shared/utility/mysql"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

var TestUser *User

func TestMain(m *testing.M) {
	glog.InitLog(&glog.Options{
		LogLevel:      "DEBUG",
		LogDir:        "./log",
		FileSizeLimit: 4294967296,
	})
	mysql.SetLogger(dblog.NewDBLogger())
	csvBase := &base.ConfigManager{}
	csvBase.Init()
	csvBase.LoadConfig("../../shared/csv/data")
	CSV := entry.NewManager()
	err := CSV.Reload(csvBase)
	if err != nil {
		log.Printf("CSV.Reload err:%v", err)
		return
	}

	common.SetRewardsType(CSV.Reward.RewardsType())

	manager.CSV = CSV
	manager.RedisClient = redis.NewClient(&redis.Options{
		Username: "root",
		Password: "",
		Addr:     ":6379",
	})
	_, err = manager.RedisClient.Ping(manager.RedisClient.Context()).Result()
	if err != nil {
		log.Printf("connect to redis error: %v", err)
		return
	}

	conf := &mysql.Config{
		Addr:            "root:xuecm01,@tcp(localhost:3306)/overlord?charset=utf8mb4&parseTime=true&loc=Local",
		MaxOpenConn:     5,
		MaxIdleConn:     5,
		ConnMaxLifetime: 4 * time.Hour,
	}

	db, err := sql.Open("mysql", conf.Addr)
	if err != nil {
		return
	}

	db.SetMaxIdleConns(conf.MaxIdleConn)
	db.SetMaxOpenConns(conf.MaxOpenConn)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)

	MySQL := mysql.NewHandler(db)

	manager.MySQL = MySQL

	manager.Global = global.NewGlobal(db, manager.RedisClient)
	manager.LuaState = lua_state.NewLuaMatcher()

	TestUser = NewUser(1001)
	TestUser.InitForCreate(context.Background())
	TestUser.Init(context.Background())
	TestUser.RewardsResult.Clear()

	m.Run()
}
