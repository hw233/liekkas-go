package config

import (
	"shared/utility/glog"
	"shared/utility/mysql"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Service        string
	ServerName     string
	GRPCListenPort string
	ETCDEndpoints  []string
	Redis          *redis.Options
	MySQL          *mysql.Config
	CSVPath        string
	Log            *glog.Options
}

func NewDefaultConfig() *Config {
	return &Config{
		Service:        "mail",
		ServerName:     "mail01",
		GRPCListenPort: "9101",
		ETCDEndpoints: []string{
			"localhost:2379",
			"localhost:2380",
			"localhost:2381",
		},
		Redis: &redis.Options{
			Username: "root",
			Password: "",
			Addr:     ":6379",
		},
		MySQL: &mysql.Config{
			Addr:        "root:root@tcp(127.0.0.1:3306)/overlord_user_go?charset=utf8mb4&parseTime=true&loc=Local",
			MaxOpenConn: 5,
			MaxIdleConn: 5,
		},
		CSVPath: "../shared/csv/data",
		Log: &glog.Options{
			LogLevel:      "DEBUG",
			LogDir:        "./log",
			FileSizeLimit: 4294967296,
		},
	}
}
