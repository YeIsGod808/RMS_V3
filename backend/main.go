package main

import (
	"fmt"
	"strconv"

	"github.com/RMS_V3/commonlib"
	"github.com/RMS_V3/config"
	"github.com/RMS_V3/log"
	"github.com/RMS_V3/logger"
	"github.com/RMS_V3/routes"
	"github.com/gin-contrib/cors"
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
	// 得到配置好路由的gin engine
	r := routes.SetRoute()

	// 使用 Swagger 中间件
	url := ginSwagger.URL("http://192.168.80.128:8005/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	// 启用日志和恢复中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	// 注意这里的addr就只能是host：port拼接，不能再format里面有任何其他东西
	addr := fmt.Sprintf("%s:%d", config.GetGlobalConfig().SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port)
	// 启动gin engine， 监听外部http请求并处理
	if err := r.Run(addr); err != nil {
		zap.L().Panic("Router.Run error: ", zap.Error(err))
	}
	zap.L().Sugar().Infof("listen on %s:%d", config.GetGlobalConfig().SvrConfig.Host, config.GetGlobalConfig().SvrConfig.Port)
}
