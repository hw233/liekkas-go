package config

import (
	"shared/utility/glog"
	"shared/utility/mysql"
	"time"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Service        string
	ServerName     string
	HTTPListenPort int32
	CSVPath        string
	ETCDEndpoints  []string
	MySQL          *mysql.Config
	Redis          *redis.Options
	Log            *glog.Options
}

func NewDefaultConfig() *Config {
	return &Config{
		Service:        "patch",
		ServerName:     "patch01",
		HTTPListenPort: 10080,
		CSVPath:        "./csv/data",
		ETCDEndpoints: []string{
			"localhost:2379",
			"localhost:2380",
			"localhost:2381",
		},
		MySQL: &mysql.Config{
			Addr:            "root:root@tcp(127.0.0.1:3306)/overlord_user_go?charset=utf8mb4&parseTime=true&loc=Local",
			MaxOpenConn:     5,
			MaxIdleConn:     5,
			ConnMaxLifetime: 4 * time.Hour,
		},
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
