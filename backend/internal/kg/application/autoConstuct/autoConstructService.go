package auto

import (
	"fmt"

	"github.com/RMS_V3/log"
	"github.com/RMS_V3/middleware/neo4jUtils"
	"github.com/RMS_V3/pkg/response"
	"github.com/gin-gonic/gin"
)

// ExtractKnowledgeFromFile 处理文件上传和知识点抽取
func ExtractKnowledgeFromFile(c *gin.Context) {
	file, err := c.FormFile("file")
	log.Info(file)
	log.Infof("Headers:", c.Request.Header)
	log.Infof("Form:", c.Request.Form)
	if err != nil {
		c.JSON(400, response.Error(400, "文件上传失败"))
		return
	}

	// 解析JSON文件
	graph, err := parseFile(file)
	if err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}

	// 验证数据
	if err := ValidateKnowledgeGraph(graph); err != nil {
		c.JSON(400, response.Error(400, err.Error()))
		return
	}

	// 批量创建节点和关系
	err = batchCreateNodesAndRelations(graph.Nodes, graph.Relations)
	if err != nil {
		c.JSON(500, response.Error(500, "图谱构建失败"))
		return
	}

	c.JSON(200, response.Success("知识图谱构建成功"))
}

func batchCreateNodesAndRelations(nodes []Node, relations []Relation) error {
	session := neo4jUtils.GetSession()
	if session == nil {
		return fmt.Errorf("无法获取 Neo4j 会话")
	}
	defer session.Close()

	// 开始事务
	tx, err := session.BeginTransaction()
	if err != nil {
		return fmt.Errorf("开始事务失败: %s", err.Error())
	}
	defer tx.Close()

	// 1. 批量创建节点
	for _, node := range nodes {
		// 检查节点是否已存在
		checkQuery := fmt.Sprintf(
			"MATCH (n:%s {name: $name}) RETURN n",
			node.Type,
		)
		result, err := tx.Run(checkQuery, map[string]interface{}{
			"name": node.Name,
		})
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("检查节点存在性失败: %s", err.Error())
		}

		// 如果节点不存在，则创建
		if !result.Next() {
			createQuery := fmt.Sprintf(
				"CREATE (n:%s {name: $name, description: $description})",
				node.Type,
			)
			params := map[string]interface{}{
				"name":        node.Name,
				"description": node.Description,
			}

			log.Infof("Creating node: %s with params: %v", createQuery, params)

			_, err = tx.Run(createQuery, params)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("创建节点失败: %s", err.Error())
			}
		}
	}

	// 2. 批量创建关系
	for _, relation := range relations {
		// 检查关系是否已存在
		checkRelationQuery := fmt.Sprintf(`
			MATCH (a:%s {name: $sourceName})-[r:%s]->(b:%s {name: $targetName})
			RETURN r
		`, relation.SourceType, relation.Type, relation.TargetType)

		checkParams := map[string]interface{}{
			"sourceName": relation.SourceName,
			"targetName": relation.TargetName,
		}

		result, err := tx.Run(checkRelationQuery, checkParams)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("检查关系存在性失败: %s", err.Error())
		}

		// 如果关系不存在，则创建
		if !result.Next() {
			createRelationQuery := fmt.Sprintf(`
				MATCH (a:%s {name: $sourceName})
				MATCH (b:%s {name: $targetName})
				CREATE (a)-[r:%s]->(b)
			`, relation.SourceType, relation.TargetType, relation.Type)

			createParams := map[string]interface{}{
				"sourceName": relation.SourceName,
				"targetName": relation.TargetName,
			}

			log.Infof("Creating relation: %s with params: %v", createRelationQuery, createParams)

			_, err = tx.Run(createRelationQuery, createParams)
			if err != nil {
				tx.Rollback()
				return fmt.Errorf("创建关系失败: %s", err.Error())
			}
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %s", err.Error())
	}

	return nil
}
