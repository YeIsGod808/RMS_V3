package snowflake

import (
	"strconv"
	"sync"
	"time"

	"github.com/RMS_V3/config"
	"github.com/RMS_V3/log"
	"github.com/bwmarrin/snowflake"
)

var (
	node          *snowflake.Node
	snowflakeOnce sync.Once
)

// initSnowflake 初始化Snowflake分布式ID生成器。
// 参数:
//
//	startTime: Snowflake计时的起始时间，格式为"2006-01-02 15:04:05"。
//	machineID: 机器ID，用于分布式环境中区分不同的节点。
//
// 该函数通过解析传入的startTime来设置Snowflake的Epoch时间，
// 并根据machineID初始化Snowflake节点，以便生成唯一的ID。
func initSnowflake(startTime string, machineID int) {
	// 初始化起始时间变量
	var st time.Time

	// 解析传入的startTime，将其转换为时间对象
	st, err := time.Parse("2006-01-02 15:04:05", startTime)
	// 如果解析出错，抛出panic
	if err != nil {
		panic(err)
	}

	// 计算Snowflake的Epoch时间，即起始时间的Unix毫秒时间戳，
	// Snowflake 算法通过将当前时间戳减去一个自定义的 Epoch 时间点来计算出相对于该 Epoch 的时间差，
	// 这个时间差会被编码到生成的唯一 ID 中。
	// 在库中是一个全局常量，所以重新定义了之后就会在NewNode时被用到，即使不作为参数传入
	snowflake.Epoch = st.UnixNano() / 1000000

	// 根据机器ID初始化Snowflake节点
	node, err = snowflake.NewNode(int64(machineID))
	// 如果初始化失败，抛出panic
	if err != nil {
		panic(err)
	}

	// 日志记录Snowflake初始化成功
	log.Infof("snowflake init success")
}

// 获取一个 snowflake generator node
func GetNode() *snowflake.Node {
	// 保证只执行一次
	snowflakeOnce.Do(func() {
		initSnowflake(time.Now().Format("2006-01-02 00:00:00"), config.GetGlobalConfig().SvrConfig.MachineID)
	})
	return node
}

// 这个函数的目的是提供一个简单的方式来生成和返回一个唯一的分布式ID字符串。
func GenID() string {
	// 该函数通过调用 GetNode() 方法获取一个生成器节点，并调用其 Generate() 方法生成一个int64类型的ID。
	id := GetNode().Generate().Int64()
	// 最后，使用 strconv.FormatInt() 方法将生成的ID转换为字符串格式并返回。
	return strconv.FormatInt(id, 10)
}
