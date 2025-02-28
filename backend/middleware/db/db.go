package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/RMS_V3/config"
	"github.com/RMS_V3/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	// db 是一个包级别的变量，用于保存唯一的单例对象,单例对象为gorm.DB
	db *gorm.DB
	// dbonce 是一个用于确保初始化操作只执行一次的 sync. Once 对象
	dbOnce sync.Once
)

// openDB 连接数据库
// 该函数根据全局配置获取数据库连接信息，并尝试建立连接
// 此函数在初始化数据库连接时会进行panic处理，如果连接失败，将终止程序运行
func openDB() {
	// 获取数据库配置
	dbConfig := config.GetGlobalConfig().DbConfig
	// 根据配置信息构造数据库连接字符串
	connArgs := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.Username,
		dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)
	// 记录数据库连接信息，便于调试
	log.Info("mdb addr:" + connArgs)

	// 尝试连接数据库
	var err error
	db, err = gorm.Open(mysql.Open(connArgs), &gorm.Config{})
	if err != nil {
		// 如果连接失败，抛出panic
		panic("failed to connect database")
	}

	// 获取数据库连接的原始接口，用于设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		// 如果获取数据库连接出错，抛出panic
		panic("fetch db connection err:" + err.Error())
	}
	// _ = db.AutoMigrate(&repository.User{})
	// 设置数据库连接池参数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConn)                                        // 设置最大空闲连接
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConn)                                        // 设置最大打开的连接
	sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxIdleTime * int64(time.Second))) // 设置空闲时间为(s)
}

// GetDB 获取数据库连接
// 该函数确保在并发环境下只打开一个数据库连接，并返回该连接实例
func GetDB() *gorm.DB {
	// 利用sync.Once确保openDB只被执行一次
	dbOnce.Do(openDB)
	// 返回数据库连接实例
	return db
}

// CloseDB 关闭数据库连接
// 在程序结束或者不需要数据库连接时，调用该函数释放数据库资源
func CloseDB() {
	if db != nil {
		// 获取数据库连接的原始接口
		sqlDB, err := db.DB()
		if err != nil {
			// 如果获取数据库连接出错，抛出panic
			panic("fetch db connection err:" + err.Error())
		}
		// 关闭数据库连接
		sqlDB.Close()
	}
}
