package application

import (
	"fmt"
	"strconv"

	"github.com/RMS_V3/internal/kg/domain"
	"github.com/RMS_V3/log"
	"github.com/RMS_V3/middleware/neo4jUtils"
	"github.com/RMS_V3/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
)

func QueryChapterNodesAndRelations(ctx *gin.Context) {
	nodeLabel := "chapter" // 指定要查询的节点标签

	session := neo4jUtils.GetSession()
	if session == nil {
		ctx.JSON(500, response.Error(500, "无法获取 Neo4j 会话"))
		return
	}
	defer session.Close()

	// 查询所有带有指定标签的节点
	queryNodes := fmt.Sprintf("MATCH (n:%s) RETURN n", nodeLabel)
	resultNodes, err := session.Run(queryNodes, map[string]interface{}{})
	if err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询节点失败: %s", err.Error())))
		return
	}

	graph := domain.Graph{}
	nodeMap := make(map[int64]*domain.Node)

	for resultNodes.Next() {
		record := resultNodes.Record()
		beforeProcNode, ok := record.Get("n")
		if !ok {
			ctx.JSON(500, response.Error(500, "类型断言失败: 无法将记录转换为 Node"))
			return
		}
		node := beforeProcNode.(dbtype.Node)
		id := (node.Id)
		nodeMap[node.Id] = &domain.Node{
			ID:   id,
			Name: node.Props["name"].(string),
			Type: nodeLabel,
		}
		if desc, ok := node.Props["description"]; ok {
			nodeMap[node.Id].Description = desc.(string)
		}
		graph.Nodes = append(graph.Nodes, *nodeMap[node.Id])
	}

	if err = resultNodes.Err(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理节点结果失败: %s", err.Error())))
		return
	}

	// 查询节点间的关系
	queryRelations := fmt.Sprintf("MATCH (a:%s)-[r]->(b:%s) RETURN r", nodeLabel, nodeLabel)
	resultRelations, err := session.Run(queryRelations, map[string]interface{}{})
	if err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询关系失败: %s", err.Error())))
		return
	}

	for resultRelations.Next() {
		record := resultRelations.Record()
		// startNode, _ := record.Get("a").(dbtype.Node)
		beforeProcRelation, _ := record.Get("r")
		// endNode, _ := record.Get("b").(dbtype.Node)
		relation := beforeProcRelation.(dbtype.Relationship)
		sourceId := int(relation.StartId)
		targetId := int(relation.EndId)

		graph.Links = append(graph.Links, domain.Link{
			Source: sourceId,
			Target: targetId,
			Type:   relation.Type,
		})
	}

	if err = resultRelations.Err(); err != nil {
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理关系结果失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, graph)
}
func QuerySectionNodesAndRelations(ctx *gin.Context) {
	nodeLabel := "section" // 指定要查询的节点标签

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

	// 查询所有带有指定标签的节点
	queryNodes := fmt.Sprintf("MATCH (n:%s) RETURN n", nodeLabel)
	resultNodes, err := tx.Run(queryNodes, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询节点失败: %s", err.Error())))
		return
	}

	graph := domain.Graph{}
	nodeMap := make(map[int64]*domain.Node)

	for resultNodes.Next() {
		record := resultNodes.Record()
		beforeProcNode, ok := record.Get("n")
		if !ok {
			tx.Rollback()
			ctx.JSON(500, response.Error(500, "类型断言失败: 无法将记录转换为 Node"))
			return
		}
		node := beforeProcNode.(dbtype.Node)
		id := node.Id
		props := node.Props
		nodeMap[node.Id] = &domain.Node{
			ID:   id,
			Name: props["name"].(string),
			Type: nodeLabel,
		}
		if desc, exists := props["description"]; exists {
			nodeMap[node.Id].Description = desc.(string)
		}
		graph.Nodes = append(graph.Nodes, *nodeMap[node.Id])
	}

	if err = resultNodes.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理节点结果失败: %s", err.Error())))
		return
	}

	// 查询节点间的关系
	queryRelations := fmt.Sprintf("MATCH (a:%s)-[r]->(b:%s) RETURN r", nodeLabel, nodeLabel)
	resultRelations, err := tx.Run(queryRelations, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询关系失败: %s", err.Error())))
		return
	}

	for resultRelations.Next() {
		record := resultRelations.Record()
		beforeProcRelation, _ := record.Get("r")
		relation := beforeProcRelation.(dbtype.Relationship)
		sourceId := int(relation.StartId)
		targetId := int(relation.EndId)

		graph.Links = append(graph.Links, domain.Link{
			Source: sourceId,
			Target: targetId,
			Type:   relation.Type,
		})
	}

	if err = resultRelations.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理关系结果失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, graph)
}
func QueryPointNodesAndRelations(ctx *gin.Context) {
	nodeLabel := "point" // 指定要查询的节点标签

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

	// 查询所有带有指定标签的节点
	queryNodes := fmt.Sprintf("MATCH (n:%s) RETURN n", nodeLabel)
	resultNodes, err := tx.Run(queryNodes, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询节点失败: %s", err.Error())))
		return
	}

	graph := domain.Graph{}
	nodeMap := make(map[int64]*domain.Node)

	for resultNodes.Next() {
		record := resultNodes.Record()
		beforeProcNode, ok := record.Get("n")
		if !ok {
			tx.Rollback()
			ctx.JSON(500, response.Error(500, "类型断言失败: 无法将记录转换为 Node"))
			return
		}
		node := beforeProcNode.(dbtype.Node)
		id := node.Id
		props := node.Props
		nodeMap[node.Id] = &domain.Node{
			ID:   id,
			Name: props["name"].(string),
			Type: nodeLabel,
		}
		if desc, exists := props["description"]; exists {
			nodeMap[node.Id].Description = desc.(string)
		}
		graph.Nodes = append(graph.Nodes, *nodeMap[node.Id])
	}

	if err = resultNodes.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理节点结果失败: %s", err.Error())))
		return
	}

	// 查询节点间的关系
	queryRelations := fmt.Sprintf("MATCH (a:%s)-[r]->(b:%s) RETURN r", nodeLabel, nodeLabel)
	resultRelations, err := tx.Run(queryRelations, map[string]interface{}{})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询关系失败: %s", err.Error())))
		return
	}

	for resultRelations.Next() {
		record := resultRelations.Record()
		beforeProcRelation, _ := record.Get("r")
		relation := beforeProcRelation.(dbtype.Relationship)
		sourceId := int(relation.StartId)
		targetId := int(relation.EndId)

		graph.Links = append(graph.Links, domain.Link{
			Source: sourceId,
			Target: targetId,
			Type:   relation.Type,
		})
	}

	if err = resultRelations.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理关系结果失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, graph)
}
func QuerySectionsByChapterId(ctx *gin.Context) {
	chapterId := ctx.Query("chapter_id") // 章节ID作为查询参数
	if chapterId == "" {
		ctx.JSON(400, response.Error(400, "参数不完整: 必须提供 chapter_id"))
		return
	}
	// int类型才能查询
	chapterNodeId, err := strconv.Atoi(chapterId)
	if err != nil {
		ctx.JSON(400, response.Error(400, "无效的章节ID"))
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

	graph := domain.Graph{}
	nodeMap := make(map[int64]*domain.Node)

	// 查询与章节节点存在包含关系的所有section节点
	querySections := `
	MATCH (c) WHERE id(c) = $chapterNodeId
	MATCH (c)-[:包含]->(s:section)
	RETURN s`
	resultSections, err := tx.Run(querySections, map[string]interface{}{
		"chapterNodeId": chapterNodeId,
	})
	log.Infow("查询section节点", "chapterNodeId", chapterId)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询section节点失败: %s", err.Error())))
		return
	}

	for resultSections.Next() {
		record := resultSections.Record()
		beforeProcNode, ok := record.Get("s")
		if !ok {
			tx.Rollback()
			ctx.JSON(500, response.Error(500, "类型断言失败: 无法将记录转换为 Node"))
			return
		}
		node := beforeProcNode.(dbtype.Node)
		id := node.Id
		nodeMap[id] = &domain.Node{
			ID:   id,
			Name: node.Props["name"].(string),
			Type: "section",
		}
		if desc, ok := node.Props["description"]; ok {
			nodeMap[id].Description = desc.(string)
		}
		graph.Nodes = append(graph.Nodes, *nodeMap[id])
	}

	if err = resultSections.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理section结果失败: %s", err.Error())))
		return
	}

	// 如果没有找到任何section节点，则直接返回
	if len(graph.Nodes) == 0 {
		ctx.JSON(200, graph)
		return
	}

	// 查询section节点之间的关系
	querySectionRelations := `
		MATCH (a:section)-[r]->(b:section)
		WHERE id(a) IN $sectionIds AND id(b) IN $sectionIds
		RETURN r`
	sectionIds := make([]interface{}, 0, len(nodeMap))
	for id := range nodeMap {
		sectionIds = append(sectionIds, id)
	}

	resultSectionRelations, err := tx.Run(querySectionRelations, map[string]interface{}{
		"sectionIds": sectionIds,
	})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询section间关系失败: %s", err.Error())))
		return
	}

	for resultSectionRelations.Next() {
		record := resultSectionRelations.Record()
		beforeProcRelation, _ := record.Get("r")
		relation := beforeProcRelation.(dbtype.Relationship)
		sourceId := int(relation.StartId)
		targetId := int(relation.EndId)

		graph.Links = append(graph.Links, domain.Link{
			Source: sourceId,
			Target: targetId,
			Type:   relation.Type,
		})
	}

	if err = resultSectionRelations.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理section间关系结果失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, graph)
}
func QueryPointsBySectionId(ctx *gin.Context) {
	sectionIdStr := ctx.Query("section_id") // section ID作为查询参数
	if sectionIdStr == "" {
		ctx.JSON(400, response.Error(400, "参数不完整: 必须提供 section_id"))
		return
	}

	sectionNodeId, err := strconv.Atoi(sectionIdStr)
	if err != nil {
		ctx.JSON(400, response.Error(400, "无效的section_id"))
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

	graph := domain.Graph{}
	nodeMap := make(map[int64]*domain.Node)

	// 查询与section节点存在包含关系的所有point节点
	queryPoints := `
        MATCH (s:section) WHERE id(s) = $sectionNodeId
        MATCH (s)-[:包含]->(p:point)
        RETURN p`
	resultPoints, err := tx.Run(queryPoints, map[string]interface{}{
		"sectionNodeId": sectionNodeId,
	})
	log.Infow("查询point节点", "sectionNodeId", sectionIdStr)
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询point节点失败: %s", err.Error())))
		return
	}

	for resultPoints.Next() {
		record := resultPoints.Record()
		beforeProcNode, ok := record.Get("p")
		if !ok {
			tx.Rollback()
			ctx.JSON(500, response.Error(500, "类型断言失败: 无法将记录转换为 Node"))
			return
		}
		node := beforeProcNode.(dbtype.Node)
		id := node.Id
		nodeMap[id] = &domain.Node{
			ID:   id,
			Name: node.Props["name"].(string),
			Type: "point",
		}
		if desc, ok := node.Props["description"]; ok {
			nodeMap[id].Description = desc.(string)
		}
		graph.Nodes = append(graph.Nodes, *nodeMap[id])
	}

	if err = resultPoints.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理point结果失败: %s", err.Error())))
		return
	}

	// 如果没有找到任何point节点，则直接返回
	if len(graph.Nodes) == 0 {
		ctx.JSON(200, graph)
		return
	}

	// 查询point节点之间的关系
	queryPointRelations := `
        MATCH (a:point)-[r]->(b:point)
        WHERE id(a) IN $pointIds AND id(b) IN $pointIds
        RETURN r`
	pointIds := make([]interface{}, 0, len(nodeMap))
	for id := range nodeMap {
		pointIds = append(pointIds, id)
	}

	resultPointRelations, err := tx.Run(queryPointRelations, map[string]interface{}{
		"pointIds": pointIds,
	})
	if err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("查询point间关系失败: %s", err.Error())))
		return
	}

	for resultPointRelations.Next() {
		record := resultPointRelations.Record()
		beforeProcRelation, _ := record.Get("r")
		relation := beforeProcRelation.(dbtype.Relationship)
		sourceId := int(relation.StartId)
		targetId := int(relation.EndId)

		graph.Links = append(graph.Links, domain.Link{
			Source: sourceId,
			Target: targetId,
			Type:   relation.Type,
		})
	}

	if err = resultPointRelations.Err(); err != nil {
		tx.Rollback()
		ctx.JSON(500, response.Error(500, fmt.Sprintf("处理point间关系结果失败: %s", err.Error())))
		return
	}

	ctx.JSON(200, graph)
}
