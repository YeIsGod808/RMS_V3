package neo4jUtlis

import (
	"fmt"
	"sync"

	"github.com/RMS_V3/config"
	"github.com/RMS_V3/log"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	// 替换为你的日志包路径
)

var (
	// db 是一个包级别的变量，用于保存唯一的单例对象, 单例对象为neo4j.Driver
	driver neo4j.Driver
	// dbOnce 是一个用于确保初始化操作只执行一次的 sync.Once 对象
	dbOnce sync.Once
)

// openDB 连接Neo4j数据库
// 该函数根据全局配置获取数据库连接信息，并尝试建立连接
// 此函数在初始化数据库连接时会进行panic处理，如果连接失败，将终止程序运行
func openNeo4jDB() {
	// 获取数据库配置
	dbConfig := config.GetGlobalConfig().Neo4jConfig

	// 根据配置信息构造Neo4j连接字符串
	connUri := fmt.Sprintf("neo4j://%s:%s@%s:%s", dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port)
	log.Info("neo4j uri:" + connUri)

	// 尝试连接Neo4j数据库
	var err error
	driver, err = neo4j.NewDriver(connUri, neo4j.BasicAuth(dbConfig.Username, dbConfig.Password, ""))
	if err != nil {
		// 如果连接失败，抛出panic
		panic("failed to connect Neo4j database: " + err.Error())
	}

	// 确保在程序退出前关闭连接
	// go func() {
	// 	<-time.After(time.Duration(dbConfig.MaxIdleTime) * time.Second)
	// 	driver.Close()
	// }()
}

// GetSession 获取Neo4j数据库会话
// 该函数确保在并发环境下只打开一个数据库连接，并返回该连接实例
func GetSession() neo4j.Session {
	// 利用sync.Once确保openDB只被执行一次
	dbOnce.Do(openNeo4jDB)

	// 返回一个新的会话
	return driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
}

// CloseDB 关闭Neo4j数据库连接
// 在程序结束或者不需要数据库连接时，调用该函数释放数据库资源
func CloseNeo4jDB() {
	if driver != nil {
		// 关闭数据库连接
		driver.Close()
	}
}
