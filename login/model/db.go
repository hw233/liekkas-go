package model

// import (
// 	"database/sql"
//
// 	"login/config"
// 	"shared/utility/glog"
// )
//
// var MySQLDB *sql.DB
//
// // var MySQLHandler *mysql.Handler
//
// func init() {
// 	mysqlConf := config.Config.MySQL
// 	var err error
// 	MySQLDB, err = sql.Open("mysql", mysqlConf.Addr)
// 	if err != nil {
// 		glog.Fatalf("sql.Open(%s) error: %v", mysqlConf.Addr, err)
// 	}
//
// 	MySQLDB.SetMaxIdleConns(mysqlConf.MaxIdleConn)
// 	MySQLDB.SetMaxOpenConns(mysqlConf.MaxOpenConn)
// 	MySQLDB.SetConnMaxLifetime(mysqlConf.ConnMaxLifetime)
//
// 	// MySQLHandler = mysql.NewHandler(MySQLDB)
// }
