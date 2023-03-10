package config

import (
	"time"

	"shared/utility/mysql"

	"github.com/go-redis/redis/v8"
)

type Config struct {
	Service        string
	ServerName     string
	HTTPListenPort string
	GRPCListenPort string
	CSVPath        string
	ETCDEndpoints  []string
	MySQL          *mysql.Config
	Redis          *redis.Options
	LogLevel       string
	ServerENV      string
	// TCPConnKeepAlive time.Duration // tcp连接保持时间
}

func NewDefaultConfig() *Config {
	return &Config{
		Service:        "login",
		ServerName:     "login01",
		HTTPListenPort: "8090",
		GRPCListenPort: "",
		CSVPath:        "../shared/csv/data",
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
		LogLevel:  "DEBUG",
		ServerENV: "",
		// TCPConnKeepAlive: 5 * time.Minute,
	}
}
