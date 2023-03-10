package config

import (
	"time"

	"shared/utility/glog"
	"shared/utility/mysql"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Service          string
	ServerName       string
	GRPCListenPort   string
	ETCDEndpoints    []string
	TCPConnKeepAlive time.Duration // tcp连接保持时间
	MaxConn          int
	Redis            *redis.Options
	MySQL            *mysql.Config
	CSVPath          string
	Log              *glog.Options
}

func NewDefaultConfig() *Config {
	return &Config{
		Service:        "guild",
		ServerName:     "guild01",
		GRPCListenPort: "9060",
		ETCDEndpoints: []string{
			"localhost:2379",
			"localhost:2380",
			"localhost:2381",
		},
		TCPConnKeepAlive: 5 * time.Minute,
		MaxConn:          10000,
		Redis: &redis.Options{
			Username: "root",
			Password: "",
			Addr:     ":6379",
		},
		MySQL: &mysql.Config{
			Addr:            "root:root@tcp(127.0.0.1:3306)/overlord_user_go?charset=utf8mb4&parseTime=true&loc=Local",
			MaxOpenConn:     5,
			MaxIdleConn:     5,
			ConnMaxLifetime: 4 * time.Hour,
		},
		CSVPath: "../shared/csv/data",
		Log: &glog.Options{
			LogLevel:      "DEBUG",
			LogDir:        "./log",
			FileSizeLimit: 4294967296,
		},
	}
}
