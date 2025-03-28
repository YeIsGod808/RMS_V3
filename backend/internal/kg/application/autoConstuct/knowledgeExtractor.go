package auto

import "fmt"

// Node 表示节点
type Node struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

// Relation 表示关系
type Relation struct {
	Type       string `json:"type"`
	SourceType string `json:"source_type"`
	TargetType string `json:"target_type"`
	SourceName string `json:"source_name"`
	TargetName string `json:"target_name"`
}

// ValidateKnowledgeGraph 验证知识图谱数据
func ValidateKnowledgeGraph(graph *KnowledgeGraph) error {
	// 验证节点
	for _, node := range graph.Nodes {
		if !isValidNodeType(node.Type) {
			return fmt.Errorf("无效的节点类型: %s", node.Type)
		}
		if node.Name == "" {
			return fmt.Errorf("节点名称不能为空")
		}
	}

	// 验证关系
	for _, relation := range graph.Relations {
		// 检查节点类型是否有效
		if !isValidNodeType(relation.SourceType) || !isValidNodeType(relation.TargetType) {
			return fmt.Errorf("关系中的节点类型无效")
		}

		// 如果是不同类型的节点之间的关系
		if relation.SourceType != relation.TargetType {
			// 检查是否为有效的层级关系
			if !isValidHierarchicalRelation(relation.SourceType, relation.TargetType) {
				return fmt.Errorf("无效的节点层级关系: %s -> %s", relation.SourceType, relation.TargetType)
			}
			// 不同类型节点间只能是包含关系
			if relation.Type != "包含" {
				return fmt.Errorf("不同类型节点间只支持包含关系，当前关系类型: %s", relation.Type)
			}
		} else {
			// 同类型节点间的关系验证
			if !isValidRelationType(relation.Type) {
				return fmt.Errorf("无效的关系类型: %s", relation.Type)
			}
		}
	}

	return nil
}

// isValidNodeType 检查节点类型是否有效
func isValidNodeType(nodeType string) bool {
	validTypes := map[string]bool{
		"chapter": true,
		"section": true,
		"point":   true,
	}
	return validTypes[nodeType]
}

// isValidHierarchicalRelation 检查层级关系是否有效
func isValidHierarchicalRelation(sourceType, targetType string) bool {
	validHierarchy := map[string]map[string]bool{
		"chapter": {"section": true},
		"section": {"point": true},
	}

	if targetTypes, exists := validHierarchy[sourceType]; exists {
		return targetTypes[targetType]
	}
	return false
}

// isValidRelationType 检查关系类型是否有效
func isValidRelationType(relationType string) bool {
	validTypes := map[string]bool{
		"包含": true,
		"前置": true,
		"相关": true,
		"扩展": true,
	}
	return validTypes[relationType]
}
