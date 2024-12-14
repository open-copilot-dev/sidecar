package pre

import (
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"strings"
)

type FilterProcessor struct {
}

// 光标在这些字符后面时，停止补全
var stopChars = []string{";", ")", "]", "}"}

func (f *FilterProcessor) process(c *domain.CompletionContext) bool {
	if c.IsCanceled() {
		return false
	}
	lineTextBeforeCursor := c.Request.CompletionLine.GetLineTextBeforeCursor()
	lineTextBeforeCursor = strings.TrimSpace(lineTextBeforeCursor)
	for _, char := range stopChars {
		if strings.HasSuffix(lineTextBeforeCursor, char) {
			return false
		}
	}
	return true
}
