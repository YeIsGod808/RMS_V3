package auto

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type KnowledgeGraph struct {
	Nodes     []Node     `json:"nodes"`
	Relations []Relation `json:"relations"`
}

// parseFile 只处理JSON文件
func parseFile(file *multipart.FileHeader) (*KnowledgeGraph, error) {
	// 检查文件扩展名
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".json" {
		return nil, fmt.Errorf("只支持JSON文件格式")
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	// 读取文件内容
	content, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析JSON
	var graph KnowledgeGraph
	if err := json.Unmarshal(content, &graph); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	return &graph, nil
}
