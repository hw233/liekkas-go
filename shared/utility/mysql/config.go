package mysql

import "time"

type Config struct {
	Addr            string
	MaxIdleConn     int
	MaxOpenConn     int
	ConnMaxLifetime time.Duration
}
