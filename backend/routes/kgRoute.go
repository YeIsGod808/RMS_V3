package routes

import (
	"github.com/RMS_V3/internal/api"
	"github.com/gin-gonic/gin"
)

func KgRoutes(r *gin.RouterGroup) {
	// 知识图谱相关路由
	knowledge := r.Group("/")
	{
		knowledge.GET("/graph/chapters", api.GetChapterGraph)
		knowledge.GET("/graph/sections/:chapterId", api.GetSectionGraph)
		knowledge.GET("/graph/points/:sectionId", api.GetPointGraph)
		knowledge.GET("/knowledge/point/:pointId", api.GetKnowledgePoint)
		knowledge.GET("/knowledge/videos/:pointId", api.GetKnowledgeVideos)
		knowledge.GET("/knowledge/exercises/:pointId", api.GetKnowledgeExercises)
		knowledge.GET("/knowledge/search", api.SearchKnowledge)
	}
}
