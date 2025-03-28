package main

import (
	"fmt"
	"strconv"

	"github.com/RMS_V3/config"
	"github.com/RMS_V3/log"
	"github.com/RMS_V3/log/logger"
	"github.com/RMS_V3/pkg/commonlib"
	"github.com/RMS_V3/routes"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go.uber.org/zap"
)

func Init() {
	if err := config.Init(); err != nil {
		log.Fatalf("init config failed, err:%v\n", err)
	}
	log.InitLog()
	loggerBuf, _ := strconv.Atoi(config.GetGlobalConfig().LogConfig.LoggerBuf)
	logger.InitLogger(loggerBuf)
	commonlib.InitDBConn()
	log.Info("log init success...")
}
func main() {
	Init()
	defer log.Sync()

	// 设置Gin模式
	if config.GetGlobalConfig().Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	r := routes.SetRoute()

	// 启用日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 使用 Swagger 中间件
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 注意这里的addr就只能是host：port拼接，不能再format里面有任何其他东西
	addr := fmt.Sprintf("%s:%d", config.GetGlobalConfig().SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port)
	// 启动gin engine， 监听外部http请求并处理
	if err := r.Run(addr); err != nil {
		zap.L().Panic("Router.Run error: ", zap.Error(err))
	}
	zap.L().Sugar().Infof("listen on %s:%d", config.GetGlobalConfig().SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port)
}
