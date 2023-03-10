package config

import (
	"shared/utility/glog"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Service          string
	ServerName       string
	TCPListenPort    string
	GRPCListenPort   string
	ETCDEndpoints    []string
	TCPConnKeepAlive time.Duration // tcp连接保持时间
	MaxConn          int
	RequestRateCap   int
	RequestRateSec   float64
	Redis            *redis.Options
	Log              *glog.Options
}

func NewDefaultConfig() *Config {
	return &Config{
		Service:        "portal",
		ServerName:     "portal01",
		TCPListenPort:  "9080",
		GRPCListenPort: "9090",
		ETCDEndpoints: []string{
			"localhost:2379",
			"localhost:2380",
			"localhost:2381",
		},
		TCPConnKeepAlive: 5 * time.Minute,
		MaxConn:          10000,
		RequestRateCap:   600,
		RequestRateSec:   10,
		Redis: &redis.Options{
			Username: "root",
			Password: "",
			Addr:     ":6379",
		},
		Log: &glog.Options{
			LogLevel:      "DEBUG",
			LogDir:        "./log",
			FileSizeLimit: 4294967296,
		},
	}
}
