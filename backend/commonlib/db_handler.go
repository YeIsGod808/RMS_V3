package commonlib

import (
	// "SQL-OJ-BACKEND/conf"
	"database/sql"
	"errors"
	"fmt"

	"github.com/RMS_V3/config"
	_ "github.com/go-sql-driver/mysql"
)

var DB_problem *sql.DB
var DB_user *sql.DB

// Before start server, create db connection pools here,
// any invalid connection should cause panic
func InitDBConn() {
	// DB_problem = createDBConn("mysql", conf.GetConf("db_problem"))
	// 获取数据库配置
	dbConfig := config.GetGlobalConfig().DbConfig
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.Username,
		dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	DB_user = createDBConn("mysql", connArgs)
}

func createDBConn(driverName string, dataSourceName string) (db *sql.DB) {
	db, _ = sql.Open(driverName, dataSourceName)
	db.SetMaxOpenConns(500)
	db.SetMaxIdleConns(50)
	if err := db.Ping(); err != nil {
		panic(errors.New(fmt.Sprint("init", driverName, dataSourceName, "failed")))
	}
	return
}
