package config

import (
	"testing"
	"time"
)

type Config struct {
	ServerID         uint32
	GRPCAddr         string
	CSVPath          string `viper:"csv_path"`
	ETCDEndpoints    []string
	MySQL            *MySQLConfig
	TCPConnKeepAlive time.Duration
}

type MySQLConfig struct {
	Addr        string
	MaxIdleConn int
	MaxOpenConn int
}

func TestLoad(t *testing.T) {
	config := &Config{}
	err := Load(config)
	if err != nil {
		t.Errorf("load config error: %v", err)
	}

	t.Logf("config: %+v", *config)
}

func TestReload(t *testing.T) {
	config := &Config{}
	err := Load(config)
	if err != nil {
		t.Errorf("load config error: %v", err)
	}

	t.Logf("config: %+v", *config)

	config.ServerID = 100
	err = Reload(config)
	if err != nil {
		t.Errorf("reload config error: %v", err)
	}

	t.Logf("config: %+v", *config)
}

func BenchmarkReload(b *testing.B) {
	config := &Config{}
	err := Load(config)
	if err != nil {
		b.Errorf("load config error: %v", err)
	}

	b.ReportAllocs()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			err = Reload(config)
			if err != nil {
				b.Errorf("reload config error: %v", err)
			}

			_ = config.MySQL.Addr
			_ = config.MySQL.MaxIdleConn
			_ = config.MySQL.MaxOpenConn
			_ = config.ETCDEndpoints[0]
		}
	})
}
