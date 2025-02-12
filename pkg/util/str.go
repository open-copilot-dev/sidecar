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

// TruncateString 截取字符串前n个字符
func TruncateString(content string, n int) string {
	if len(content) > n {
		return content[:n]
	}
	return content
}

// TryToJson 将对象转换为json字符串，出错时返回空字符串
func TryToJson(object any) string {
	bytes, err := json.Marshal(object)
	if err != nil {
		return ""
	}
	return string(bytes)
}
