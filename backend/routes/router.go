package routes

import (
	"github.com/RMS_V3/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetRoute() *gin.Engine {
	if config.GetGlobalConfig().Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	r := gin.New()
	// 配置CORS中间件
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000","http://127.0.0.1:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	// 使用group方法为gin engine创建一个路由组
	rms := r.Group("/api")
	{
		UserRoutes(rms)
		KgRoutes(rms)
	}
	return r
}
