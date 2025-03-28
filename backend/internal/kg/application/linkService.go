package application

import (
	"fmt"

	"github.com/RMS_V3/log"
	"github.com/RMS_V3/middleware/neo4jUtils"
	"github.com/RMS_V3/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func AddRelationBetweenNodes(ctx *gin.Context) {
	// 参数解析
	sourceNodeName := ctx.Query("source_name")
	targetNodeName := ctx.Query("target_name")
	relationType := ctx.Query("relation_type")
	sourceLabel := ctx.Query("source_type") // 源节点标签，如 "chapter", "section", 或 "point"
	targetLabel := ctx.Query("target_type") // 目标节点标签，如 "chapter", "section", 或 "point"

	// 校验参数完整性
	if sourceNodeName == "" || targetNodeName == "" || relationType == "" || sourceLabel == "" || targetLabel == "" {
		ctx.JSON(400, response.Error(400, "参数不完整: 必须提供 source_name, target_name, relation_type, source_label 和 target_label"))
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

	// 检查源节点是否存在
	checkSourceQuery := fmt.Sprintf("MATCH (a:%s {name: $nodeName}) RETURN a", sourceLabel)
	result, err := tx.Run(checkSourceQuery, map[string]interface{}{"nodeName": sourceNodeName})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询源节点失败: %s", err.Error())))
		return
	}
	// 源节点不存在，返回400
	if !result.Next() {
		tx.Rollback()
		ctx.JSON(400, response.Error(400, fmt.Sprintf("源节点 '%s' 不存在", sourceNodeName)))
		return
	}

	// 检查目标节点是否存在
	checkTargetQuery := fmt.Sprintf("MATCH (b:%s {name: $nodeName}) RETURN b", targetLabel)
	result, err = tx.Run(checkTargetQuery, map[string]interface{}{"nodeName": targetNodeName})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询目标节点失败: %s", err.Error())))
		return
	}
	// 目标节点不存在，返回400
	if !result.Next() {
		tx.Rollback()
		ctx.JSON(400, response.Error(400, fmt.Sprintf("目标节点 '%s' 不存在", targetNodeName)))
		return
	}
	// 检查关系是否存在
	checkRelationQuery := fmt.Sprintf(`
		MATCH (a:%s {name: $sourceName})-[r:%s]->(b:%s {name: $targetName})
		RETURN r
	`, sourceLabel, relationType, targetLabel)

	checkParams := map[string]interface{}{
		"sourceName": sourceNodeName,
		"targetName": targetNodeName,
	}

	result, err = tx.Run(checkRelationQuery, checkParams)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("检查关系失败: %s", err.Error())))
		return
	}
	// 如果节点间关系已存在，则回滚事务，并且返回400,以及错误信息
	if result.Next() {
		tx.Rollback()
		ctx.JSON(400, response.Error(400, fmt.Sprintf("从节点 %s 到节点 %s 关系 '%s' 已存在", sourceNodeName, targetNodeName, relationType)))
		return
	}
	// 创建关系查询
	query := fmt.Sprintf(`
        MATCH (a:%s {name: $sourceName}), (b:%s {name: $targetName})
        CREATE (a)-[r:%s]->(b)
        RETURN r
    `, sourceLabel, targetLabel, relationType)

	params := map[string]interface{}{
		"sourceName": sourceNodeName,
		"targetName": targetNodeName,
	}

	log.Infof("Executing query: %s with params: %v", query, params)

	result, err = tx.Run(query, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("创建关系失败: %s", err.Error())))
		return
	}

	_, err = result.Single()
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("关系创建失败: %s", err.Error())))
		return
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("提交事务失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, response.Success(fmt.Sprintf("成功为 %s 节点 '%s' 和 %s 节点 '%s' 添加了关系 '%s'", sourceLabel, sourceNodeName, targetLabel, targetNodeName, relationType)))
}

func DeleteRelationBetweenNodes(ctx *gin.Context) {
	sourceNodeName := ctx.Query("source_name")
	targetNodeName := ctx.Query("target_name")
	sourceLabel := ctx.Query("source_type")    // 源节点标签
	targetLabel := ctx.Query("target_type")    // 目标节点标签
	relationType := ctx.Query("relation_type") // 关系类型

	if sourceNodeName == "" || targetNodeName == "" || sourceLabel == "" || targetLabel == "" || relationType == "" {
		ctx.JSON(400, response.Error(400, "参数不完整: 必须提供 source_name, target_name, source_label, target_label和relation_type"))
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

	query := fmt.Sprintf(`
        MATCH (a:%s {name: $sourceName})-[r:%s]->(b:%s {name: $targetName})
        RETURN r
    `, sourceLabel, relationType, targetLabel)

	params := map[string]interface{}{
		"sourceName": sourceNodeName,
		"targetName": targetNodeName,
	}

	log.Info("Executing query to find relations: %s with params: %v", query, params)

	result, err := tx.Run(query, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询关系失败: %s", err.Error())))
		return
	}

	var relationsDeleted int = 0
	for result.Next() {
		record := result.Record()
		relation, ok := record.Get("r")
		if !ok {
			tx.Rollback()
			ctx.JSON(500, response.Error(500, "类型断言失败: 无法将记录转换为 Relationship"))
			return
		}
		processRelation := relation.(dbtype.Relationship)
		deleteQuery := fmt.Sprintf(`
            MATCH (a:%s {name: $sourceName})-[r:%s]->(b:%s {name: $targetName})
            WHERE id(r) = $relationId
            DELETE r
        `, sourceLabel, relationType, targetLabel)

		deleteParams := map[string]interface{}{
			"sourceName":   sourceNodeName,
			"targetName":   targetNodeName,
			"relationId":   processRelation.Id,
			"relationType": relationType,
		}

		log.Info("Executing delete relation query: %s with params: %v", deleteQuery, deleteParams)

		deleteResult, err := tx.Run(deleteQuery, deleteParams)
		if err != nil {
			tx.Rollback()
			ctx.JSON(500, response.Error(500, fmt.Sprintf("删除关系失败: %s", err.Error())))
			return
		}
		if deleteResult.Err() != nil {
			tx.Rollback()
			ctx.JSON(500, response.Error(500, fmt.Sprintf("处理结果失败: %s", deleteResult.Err())))
			return
		}
		relationsDeleted++
	}

	if err = result.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理结果失败: %s", err.Error())))
		return
	}
	// 要删除的关系不存在时，返回400
	if relationsDeleted == 0 {
		ctx.JSON(400, response.Error(400, "没有找到要删除的关系,请检查节点名字和类型是否正确以及节点关系是否存在"))
		return
	}

	if err = tx.Commit(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("提交事务失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, response.Success(fmt.Sprintf("成功删除了 %s 节点 '%s' 和 %s 节点 '%s' 之间的 %d 个关系类型 '%s'", sourceLabel, sourceNodeName, targetLabel, targetNodeName, relationsDeleted, relationType)))
}

func UpdateRelationBetweenNodes(ctx *gin.Context) {
	sourceNodeName := ctx.Query("source_name")
	targetNodeName := ctx.Query("target_name")
	newRelationType := ctx.Query("new_relation_type")
	sourceLabel := ctx.Query("source_type")           // 源节点标签
	targetLabel := ctx.Query("target_type")           // 目标节点标签
	oldRelationType := ctx.Query("old_relation_type") // 旧关系类型

	if sourceNodeName == "" || targetNodeName == "" || newRelationType == "" || sourceLabel == "" || targetLabel == "" || oldRelationType == "" {
		ctx.JSON(400, response.Error(400, "参数不完整: 必须提供 source_name, target_name, new_relation_type, source_label, target_label和old_relation_type"))
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
	// 查询关系是否存在
	queryFind := fmt.Sprintf(`
		MATCH (a:%s {name: $sourceName})-[r:%s]->(b:%s {name: $targetName})
		RETURN r
	`, sourceLabel, oldRelationType, targetLabel)

	params := map[string]interface{}{
		"sourceName":   sourceNodeName,
		"targetName":   targetNodeName,
		"relationType": oldRelationType,
	}

	log.Infof("Executing query to find relations: %s with params: %v", queryFind, params)

	result, err := tx.Run(queryFind, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询关系失败: %s", err.Error())))
		return
	}

	hasRelation := result.Next() // 检查是否有匹配的结果

	if !hasRelation {
		tx.Rollback()
		ctx.JSON(400, response.Error(400, "指定的关系不存在，请检查节点名字和类型是否正确以及节点关系是否存在"))
		return
	}

	// 删除已存在的关系
	queryDelete := fmt.Sprintf(`
		MATCH (a:%s {name: $sourceName})-[r:%s]->(b:%s {name: $targetName})
		DELETE r
	`, sourceLabel, oldRelationType, targetLabel)

	log.Infof("Executing query to delete relation: %s with params: %v", queryDelete, params)

	_, err = tx.Run(queryDelete, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("删除关系失败: %s", err.Error())))
		return
	}

	// 创建新的关系
	queryCreateNewRelation := fmt.Sprintf(`
        MATCH (a:%s {name: $sourceName}), (b:%s {name: $targetName})
        CREATE (a)-[r:%s]->(b)
        RETURN r
    `, sourceLabel, targetLabel, newRelationType)

	log.Info("Executing query to create new relation: %s with params: %v", queryCreateNewRelation, params)

	result, err = tx.Run(queryCreateNewRelation, params)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("创建新关系失败: %s", err.Error())))
		return
	}

	_, err = result.Single()
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("获取新关系失败: %s", err.Error())))
		return
	}

	if err = tx.Commit(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("提交事务失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, response.Success(fmt.Sprintf("成功将 %s 节点 '%s' 和 %s 节点 '%s' 之间的关系从 '%s' 更新为 '%s'", sourceLabel, sourceNodeName, targetLabel, targetNodeName, oldRelationType, newRelationType)))
}
func QueryRelationsBetweenTypes(ctx *gin.Context) {
	sourceType := ctx.Query("source_type")
	targetType := ctx.Query("target_type")
	session := neo4jUtils.GetSession()
	if session == nil {
		ctx.JSON(500, gin.H{"error": "无法获取 Neo4j 会话"})
		return
	}
	defer session.Close()

	// 查询所有源节点与目标节点间的关系
	query := `
        MATCH (a:%s)-[r]->(b:%s)
        RETURN a.name AS sourceName, b.name AS targetName, type(r) AS relationType, r
    `
	result, err := session.Run(fmt.Sprintf(query, sourceType, targetType), map[string]interface{}{})
	if err != nil {
		ctx.JSON(500, gin.H{"error": fmt.Sprintf("查询关系失败: %s", err.Error())})
		return
	}

	var relations []map[string]interface{}
	for result.Next() {
		record := result.Record()
		sourceName, _ := record.Get("sourceName")
		targetName, _ := record.Get("targetName")
		relationType, _ := record.Get("relationType")

		relation := map[string]interface{}{
			"sourceName":   sourceName,
			"targetName":   targetName,
			"relationType": relationType,
		}
		relations = append(relations, relation)
	}

	if err = result.Err(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理结果失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, response.Success(relations))
}
