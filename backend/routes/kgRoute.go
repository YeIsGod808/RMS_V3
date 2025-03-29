package routes

import (
	"github.com/RMS_V3/internal/kg/application"
	analysis "github.com/RMS_V3/internal/kg/application/Analysis"
	recommend "github.com/RMS_V3/internal/kg/application/Recommend"
	auto "github.com/RMS_V3/internal/kg/application/autoConstuct"
	"github.com/gin-gonic/gin"
)

func KgRoutes(r *gin.RouterGroup) {
	// 知识图谱相关路由
	knowledge := r.Group("/")
	{
		// 节点相关路由
		knowledge.POST("/knowledge/addNode", application.AddNode)
		knowledge.POST("/knowledge/deleteNode", application.DeleteNode)
		knowledge.POST("/knowledge/updateNode", application.UpdateNode)
		knowledge.GET("knowledge/searchByKeyword", application.SearchNodesByKeyword)
		// 关系相关路由
		knowledge.POST("/knowledge/addLink", application.AddRelationBetweenNodes)
		knowledge.POST("/knowledge/deleteLink", application.DeleteRelationBetweenNodes)
		knowledge.POST("/knowledge/updateLink", application.UpdateRelationBetweenNodes)
		knowledge.GET("knowledge/relation", application.QueryRelationsBetweenTypes)
		// 图相关路由
		knowledge.GET("/knowledge/chapter", application.QueryChapterNodesAndRelations)
		knowledge.GET("/knowledge/section", application.QuerySectionNodesAndRelations)
		knowledge.GET("/knowledge/point", application.QueryPointNodesAndRelations)
		knowledge.GET("/knowledge/sectionByID", application.QuerySectionsByChapterId)
		knowledge.GET("/knowledge/pointByID", application.QueryPointsBySectionId)
		// 资源相关路由
		knowledge.POST("/knowledge/uploadResource", application.UploadResource)
		knowledge.GET("/knowledge/videosByPointId", application.GetPointVideo)
		knowledge.GET("/knowledge/exercisesByPointId", application.GetPointExercise)
		knowledge.GET("/knowledge/coursewaresByPointId", application.GetPointCourseware)
		knowledge.POST("/knowledge/deleteVideo", application.DeletePointVideo)
		knowledge.POST("/knowledge/deleteExercise", application.DeletePointExercise)
		knowledge.POST("/knowledge/deleteCourseware", application.DeletePointCourseware)
		// 高级特性
		knowledge.POST("/knowledge/autoConstruct", auto.ExtractKnowledgeFromFile)
		knowledge.GET("/knowledge/learningDifficulty", analysis.AssessLearningDifficulty)
		knowledge.POST("/knowledge/pathRecommend", recommend.GenerateLearningPath)
		knowledge.GET("/knowledge/analyzeConnections", analysis.AnalyzeKnowledgeConnections)
	}
}
