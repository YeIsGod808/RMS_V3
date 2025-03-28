package neo4jUtils

import (
	"fmt"
	"strings"
)

// sanitizeLabel 确保标签名安全，避免潜在的注入风险
func SanitizeLabel(label string) (string, error) {
	// 确保标签以字母开头
	if len(label) == 0 || !strings.Contains("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ", string(label[0])) {
		return "", fmt.Errorf("标签必须以字母开头")
	}

	// 只允许字母、数字和下划线
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	var clean strings.Builder
	for _, char := range label {
		if strings.ContainsRune(allowedChars, char) {
			clean.WriteRune(char)
		}
	}

	cleanedLabel := clean.String()
	if cleanedLabel == "" {
		return "", fmt.Errorf("无效的标签名")
	}

	return cleanedLabel, nil
}
