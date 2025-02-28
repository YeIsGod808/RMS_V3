package routes

import (
	"github.com/RMS_V3/config"
	"github.com/gin-gonic/gin"
)

func SetRoute() *gin.Engine {
	if config.GetGlobalConfig().Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	// 使用group方法为gin engine创建一个路由组
	rms := r.Group("/api")
	{
		UserRoutes(rms)
		KgRoutes(rms)
	}
	return r
}
