package recommend

import (
	"fmt"

	"github.com/RMS_V3/log"

	"github.com/RMS_V3/middleware/neo4jUtils"
	"github.com/RMS_V3/pkg/response"
	"github.com/gin-gonic/gin"
)

func GenerateLearningPath(c *gin.Context) {
	// 获取请求参数
	startPointName := c.Query("start_point_name")
	startPointType := c.Query("start_point_type")
	endPointName := c.Query("end_point_name")
	endPointType := c.Query("end_point_type")

	if startPointName == "" || endPointName == "" || startPointType == "" || endPointType == "" {
		c.JSON(400, response.Error(400, "缺少参数: start_point 或 end_point"))
		return
	}

	// 动态拼接标签
	query := fmt.Sprintf(`
	MATCH path = shortestPath((start:%s {name: $startPointName})-[*]->(end:%s {name: $endPointName}))
	RETURN [node in nodes(path) | node.name] AS pathNodes,
	       [rel in relationships(path) | type(rel)] AS pathRelationships
	`, startPointType, endPointType)

	// 执行查询
	session := neo4jUtils.GetSession()
	defer session.Close()

	result, err := session.Run(query, map[string]interface{}{
		"startPointName": startPointName,
		"endPointName":   endPointName,
	})
	log.Infof("查询结果: %v", result)
	if err != nil {
		c.JSON(500, response.Error(500, "路径生成失败"))
		return
	}

	// 处理查询结果
	if !result.Next() {
		c.JSON(404, response.Error(404, "未找到从起点到终点的路径"))
		return
	}

	record := result.Record()

	// 提取路径中的节点和关系
	pathNodesValue, _ := record.Get("pathNodes")
	pathRelationshipsValue, _ := record.Get("pathRelationships")

	pathNodes := pathNodesValue.([]interface{})
	pathRelationships := pathRelationshipsValue.([]interface{})

	// 构造学习路径
	learningPath := constructLearningPath(pathNodes, pathRelationships)

	// 返回结果
	c.JSON(200, response.Success(learningPath))
}

// 构造学习路径
func constructLearningPath(nodes []interface{}, relationships []interface{}) map[string]interface{} {
	// 将节点名称提取为字符串数组
	nodeNames := make([]string, len(nodes))
	for i, node := range nodes {
		nodeNames[i] = node.(string)
	}

	// 将关系类型提取为字符串数组
	relationTypes := make([]string, len(relationships))
	for i, rel := range relationships {
		relationTypes[i] = rel.(string)
	}

	// 构造路径描述
	pathDescription := []string{}
	for i := 0; i < len(nodeNames)-1; i++ {
		pathDescription = append(pathDescription, fmt.Sprintf("%s -> (%s) -> %s",
			nodeNames[i], relationTypes[i], nodeNames[i+1]))
	}

	return map[string]interface{}{
		"nodes":         nodeNames,
		"relationships": relationTypes,
		"description":   pathDescription,
	}
}
