package application

import (
	"fmt"

	"github.com/RMS_V3/internal/kg/domain"
	"github.com/RMS_V3/log"
	"github.com/RMS_V3/middleware/neo4jUtils"
	"github.com/RMS_V3/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func AddNode(ctx *gin.Context) {
	// 参数解析
	nodeName := ctx.Query("name")
	nodeType := ctx.Query("type")
	description := ctx.Query("description")

	// 校验节点数据完整性
	if nodeType == "" || nodeName == "" {
		ctx.JSON(400, response.Error(400, "节点数据不完整: Type 和 Name 必须提供"))
		return
	}

	// 获取数据库会话
	session := neo4jUtils.GetSession()
	if session == nil {
		ctx.JSON(500, response.Error(500, "无法获取 Neo4j 会话"))
		return
	}
	defer session.Close()

	// 开始事务
	tx, err := session.BeginTransaction()
	if err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("开始事务失败: %s", err.Error())))
		return
	}
	defer tx.Close()

	// 构造查询语句以检查节点是否存在
	sanitizedType, err := neo4jUtils.SanitizeLabel(nodeType)
	if err != nil {
		ctx.JSON(400, response.Error(400, "无效的标签名: "+err.Error()))
		return
	}

	checkQuery := fmt.Sprintf("MATCH (n:%s {name: $name}) RETURN n", sanitizedType)
	params := map[string]interface{}{
		"name": nodeName,
	}

	log.Infof("Executing check query: %s with params: %v", checkQuery, params)

	// 执行检查查询
	result, err := tx.Run(checkQuery, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("Neo4j 查询失败: %s", err.Error())))
		return
	}

	var existingNode *domain.Node

	// 处理检查查询结果
	for result.Next() {
		record := result.Record()
		nodeValue, ok := record.Get("n")
		if !ok {
			continue // 跳过无效记录
		}
		node, ok := nodeValue.(dbtype.Node)
		if !ok {
			continue // 跳过无效记录
		}
		id := node.Id
		nodeType := node.Labels[0]
		properties := node.Props

		name, ok := properties["name"].(string)
		if !ok {
			log.Warn("类型断言失败: 无法将 'name' 断言为 string")
			return
		}
		desc, ok := properties["description"].(string)
		if !ok {
			log.Warn("类型断言失败: 无法将 'description' 断言为 string")
			return
		}

		existingNode = &domain.Node{
			ID:          id,
			Name:        name,
			Type:        nodeType,
			Description: desc,
		}

		log.Infof("已存在的节点: ID=%d, Name=%s, Description=%s\n", id, name, desc)
	}
	if err := result.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("Neo4j 结果处理失败: %s", err.Error())))
		return
	}

	if existingNode != nil {
		// 如果节点已存在，则返回已存在的节点信息
		ctx.JSON(200, gin.H{
			"message": "节点已存在",
			"node":    existingNode,
		})
		return
	}

	// 如果不存在相同名称的节点，则创建新节点
	createQuery := fmt.Sprintf("CREATE (n:%s {name: $name, description: $description}) RETURN n", sanitizedType)
	createParams := map[string]interface{}{
		"name":        nodeName,
		"description": description,
	}
	log.Infof("Executing create query: %s with params: %v", createQuery, createParams)

	// 执行创建查询
	result, err = tx.Run(createQuery, createParams)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("Neo4j 创建查询失败: %s", err.Error())))
		return
	}

	// 处理创建查询结果
	var createdNode *domain.Node
	for result.Next() {
		record := result.Record()
		log.Info("Result has records")
		nodeValue, ok := record.Get("n")
		if !ok {
			continue // 跳过无效记录
		}
		node, ok := nodeValue.(dbtype.Node)
		if !ok {
			continue // 跳过无效记录
		}
		id := node.Id
		nodeType := node.Labels[0]
		properties := node.Props

		name, ok := properties["name"].(string)
		if !ok {
			log.Warn("类型断言失败: 无法将 'name' 断言为 string")
			return
		}
		desc, ok := properties["description"].(string)
		if !ok {
			log.Warn("类型断言失败: 无法将 'description' 断言为 string")
			return
		}
		// 新创建的节点的id就是自增id
		createdNode = &domain.Node{
			ID:          id,
			Name:        name,
			Type:        nodeType,
			Description: desc,
		}

		log.Infof("创建的节点: ID=%d, Name=%s, Description=%s\n", id, name, desc)
	}
	if err := result.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("Neo4j 创建结果处理失败: %s", err.Error())))
		return
	}

	if createdNode == nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, "未能成功创建节点"))
		return
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("提交事务失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, response.Success(createdNode))
}

func UpdateNode(ctx *gin.Context) {
	nodeName := ctx.Query("name")
	nodeLabel := ctx.Query("node_type") // 节点标签，如 "chapter", "section", 或 "point"
	propertyName := ctx.Query("property_name")
	newValue := ctx.Query("new_value")

	if nodeName == "" || nodeLabel == "" || propertyName == "" || newValue == "" {
		ctx.JSON(400, response.Error(400, "参数不完整: 必须提供 name, node_label, property_name 和 new_value"))
		return
	}
	if propertyName == "name" {
		ctx.JSON(400, response.Error(400, "不允许修改 name 属性"))
		return
	}

	session := neo4jUtils.GetSession()
	if session == nil {
		ctx.JSON(500, response.Error(500, "无法获取 Neo4j 会话"))
		return
	}
	defer session.Close()

	tx, err := session.BeginTransaction()
	if err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("开始事务失败: %s", err.Error())))
		return
	}
	defer tx.Close()

	// 检查节点是否存在
	checkQuery := fmt.Sprintf("MATCH (n:%s {name: $nodeName}) RETURN n", nodeLabel)
	checkParams := map[string]interface{}{
		"nodeName": nodeName,
	}

	log.Infof("Executing check query: %s with params: %v", checkQuery, checkParams)

	result, err := tx.Run(checkQuery, checkParams)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询节点失败: %s", err.Error())))
		return
	}

	if !result.Next() {
		tx.Rollback()
		ctx.JSON(400, response.Error(404, fmt.Sprintf("节点 '%s' 类型 '%s' 不存在", nodeName, nodeLabel)))
		return
	}

	// 节点存在，执行更新操作
	query := fmt.Sprintf("MATCH (n:%s {name: $nodeName}) SET n.%s = $newValue RETURN n", nodeLabel, propertyName)
	params := map[string]interface{}{
		"nodeName": nodeName,
		"newValue": newValue,
	}

	log.Infof("Executing update query: %s with params: %v", query, params)

	result, err = tx.Run(query, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("更新节点失败: %s", err.Error())))
		return
	}

	_, err = result.Single()
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("获取更新后的节点失败: %s", err.Error())))
		return
	}

	if err = tx.Commit(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("提交事务失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, response.Success(fmt.Sprintf("成功更新 %s 节点 '%s' 的属性 '%s'", nodeLabel, nodeName, propertyName)))
}
func DeleteNode(ctx *gin.Context) {
	nodeName := ctx.Query("name")
	nodeLabel := ctx.Query("node_type")

	if nodeName == "" || nodeLabel == "" {
		ctx.JSON(400, response.Error(400, "参数不完整: 必须提供 name 和 node_type"))
		return
	}

	session := neo4jUtils.GetSession()
	if session == nil {
		ctx.JSON(500, response.Error(500, "无法获取 Neo4j 会话"))
		return
	}
	defer session.Close()

	tx, err := session.BeginTransaction()
	if err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("开始事务失败: %s", err.Error())))
		return
	}
	defer tx.Close()

	// 1. 检查节点是否存在
	checkQuery := fmt.Sprintf("MATCH (n:%s {name: $nodeName}) RETURN n", nodeLabel)
	params := map[string]interface{}{
		"nodeName": nodeName,
	}

	log.Infof("Executing check query: %s with params: %v", checkQuery, params)

	result, err := tx.Run(checkQuery, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询节点失败: %s", err.Error())))
		return
	}

	if !result.Next() {
		tx.Rollback()
		ctx.JSON(400, response.Error(400, fmt.Sprintf("类型为 '%s' 的节点 '%s' 不存在", nodeLabel, nodeName)))
		return
	}

	// 2. 根据节点类型构建级联删除查询
	var deleteQuery string
	switch nodeLabel {
	case "chapter":
		// 删除章节及其下所有小节和知识点
		deleteQuery = `
			MATCH (c:chapter {name: $nodeName})
			OPTIONAL MATCH (c)-[:包含*]->(child)
			WITH c, child
			DETACH DELETE c, child`
	case "section":
		// 删除小节及其下所有知识点
		deleteQuery = `
			MATCH (s:section {name: $nodeName})
			OPTIONAL MATCH (s)-[:包含*]->(child)
			WITH s, child
			DETACH DELETE s, child`
	case "point":
		// 删除知识点及其关联关系
		deleteQuery = `
			MATCH (p:point {name: $nodeName})
			DETACH DELETE p`
	default:
		tx.Rollback()
		ctx.JSON(400, response.Error(400, fmt.Sprintf("不支持的节点类型: %s", nodeLabel)))
		return
	}

	// 3. 执行级联删除操作
	log.Infof("Executing cascade delete query: %s with params: %v", deleteQuery, params)
	result, err = tx.Run(deleteQuery, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("删除节点失败: %s", err.Error())))
		return
	}

	if result.Err() != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理结果失败: %s", result.Err())))
		return
	}

	// 4. 获取删除的节点数量（可选）
	countQuery := `
		MATCH (n)
		WHERE n.name = $nodeName
		RETURN count(n) as count`

	result, err = tx.Run(countQuery, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("验证删除失败: %s", err.Error())))
		return
	}

	// 5. 提交事务
	if err = tx.Commit(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("提交事务失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, response.Success(fmt.Sprintf("成功删除 %s 节点 '%s' 及其相关节点和关系", nodeLabel, nodeName)))
}

// 根据关键词搜索节点（修改没有返回id）
func SearchNodesByKeyword(ctx *gin.Context) {
	keyword := ctx.Query("keyword")
	if keyword == "" {
		ctx.JSON(400, response.Error(400, "搜索关键词不能为空"))
		return
	}
	session := neo4jUtils.GetSession()
	if session == nil {
		ctx.JSON(500, response.Error(500, "无法获取 Neo4j 会话"))
		return
	}
	defer session.Close()

	tx, err := session.BeginTransaction()
	if err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("开始事务失败: %s", err.Error())))
		return
	}
	defer tx.Close()

	// 修改查询，将 n.id 改为 ID(n)
	query := `
        MATCH (n)
        WHERE n.name CONTAINS $keyword
        RETURN labels(n) AS labels, ID(n) AS id, n.name AS name, n.description AS description
    `
	result, err := tx.Run(query, map[string]interface{}{
		"keyword": keyword,
	})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("查询节点失败: %s", err.Error())})
		return
	}

	var nodes []map[string]interface{}
	for result.Next() {
		record := result.Record()
		labels, _ := record.Get("labels")
		id, _ := record.Get("id")
		name, _ := record.Get("name")
		description, _ := record.Get("description")

		labelSlice, ok := labels.([]interface{})
		if !ok || len(labelSlice) == 0 {
			continue // 跳过没有标签或类型转换失败的记录
		}

		nodeType := labelSlice[0].(string) // 假设每个节点只有一个主要标签
		level := determineLevel(nodeType)

		node := map[string]interface{}{
			"id":          id,
			"name":        name,
			"type":        nodeType,
			"level":       level,
			"description": description,
		}
		nodes = append(nodes, node)
	}

	if err = result.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("获取节点失败: %s", err.Error())))
		return
	}
	ctx.JSON(200, response.Success(nodes))
}

// 根据节点类型确定其级别
func determineLevel(nodeType string) string {
	switch nodeType {
	case "chapter":
		return "chapter"
	case "section":
		return "section"
	case "point":
		return "point"
	default:
		return "unknown"
	}
}
