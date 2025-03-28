package analysis

import (
	"strconv"

	"github.com/RMS_V3/log"

	"github.com/RMS_V3/middleware/neo4jUtils"
	"github.com/RMS_V3/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

// 分析知识点关联度
func AnalyzeKnowledgeConnections(c *gin.Context) {
	session := neo4jUtils.GetSession()
	defer session.Close()

	// 使用Neo4j的中心度算法
	query := `
    CALL gds.degree.stream('knowledge-graph')
    YIELD nodeId, score
    RETURN gds.util.asNode(nodeId).name AS name, score
    ORDER BY score DESC
    `

	result, err := session.Run(query, nil)
	if err != nil {
		c.JSON(500, response.Error(500, "分析失败"))
		return
	}

	// 处理分析结果
	analysis := processAnalysisResult(result)
	c.JSON(200, response.Success(analysis))
}

// 处理分析结果
func processAnalysisResult(result neo4j.Result) []map[string]interface{} {
	var analysis []map[string]interface{}

	for result.Next() {
		record := result.Record()
		name, _ := record.Get("name")
		score, _ := record.Get("score")

		analysis = append(analysis, map[string]interface{}{
			"name":  name.(string),
			"score": score.(float64),
		})
	}
	return analysis
}

// 评估学习难度
func AssessLearningDifficulty(c *gin.Context) {
	// 获取请求参数
	pointIdStr := c.Query("point_id")
	if pointIdStr == "" {
		c.JSON(400, response.Error(400, "缺少参数: point_id"))
		return
	}
	pointId, err := strconv.Atoi(pointIdStr)
	// Neo4j 查询语句
	query := `
	MATCH (p:point)
	WHERE id(p) = $pointId
	CALL {
		WITH p
		MATCH (p)<-[:前置*]-(prereq)
		RETURN count(prereq) AS prereqCount,
			max(COUNT { (prereq)-[:前置*]->(p) }) AS maxDepth
	}
	RETURN prereqCount, maxDepth
	`

	// 执行查询
	session := neo4jUtils.GetSession()
	defer session.Close()

	result, err := session.Run(query, map[string]interface{}{"pointId": pointId})
	if err != nil {
		c.JSON(500, response.Error(500, "查询失败"))
		return
	}
	log.Infof("查询结果: %v", result)

	// 处理查询结果
	var prereqCount int64
	var maxDepth int64
	if result.Next() {
		record := result.Record()
		prereqCountValue, _ := record.Get("prereqCount")
		maxDepthValue, _ := record.Get("maxDepth")

		prereqCount = prereqCountValue.(int64)
		if maxDepthValue != nil {
			maxDepth = maxDepthValue.(int64)
		} else {
			maxDepth = 0 // 如果没有前置知识点，深度为 0
		}
	} else {
		c.JSON(404, response.Error(404, "未找到对应的知识点"))
		return
	}

	// 计算难度分数
	difficultyScore := calculateDifficultyScore(prereqCount, maxDepth)

	// 返回结果
	c.JSON(200, response.Success(map[string]interface{}{
		"point_id":         pointId,
		"prereq_count":     prereqCount,
		"max_depth":        maxDepth,
		"difficulty_score": difficultyScore,
	}))
}

// 计算难度分数
func calculateDifficultyScore(prereqCount, maxDepth int64) float64 {
	// 示例公式：难度分数 = 前置知识点数量 + 最深深度 * 权重
	depthWeight := 1.5 // 深度权重，可以根据实际情况调整
	return float64(prereqCount) + float64(maxDepth)*depthWeight
}
