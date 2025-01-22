package util

import (
	"encoding/json"
	"strings"
)

// CalcIndent 计算字符串开头的空白字符个数
func CalcIndent(str string) int {
	return len(str) - len(strings.TrimLeft(str, " \t"))
}

// IsBlank 判断字符串是否为空
func IsBlank(str string) bool {
	return strings.TrimSpace(str) == ""
}

func TruncateString(content string, cnt int) string {
	if len(content) > cnt {
		return content[:cnt]
	}
	return content
}

func TryToJson(object any) string {
	bytes, err := json.Marshal(object)
	if err != nil {
		return ""
	}
	return string(bytes)
}
