package pre

import (
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"strings"
)

type FilterProcessor struct {
}

// 光标之前是这些字符时，停止补全
var cursorBeforeStopChars = []string{";", ")", "]", "}"}

// 光标之后如果有字符，并且是这些字符时，才允许补全
var cursorAfterAllowChars = []string{")", "]"}

func (f *FilterProcessor) process(c *domain.CompletionContext) bool {
	if c.IsCanceled() {
		return false
	}

	// 光标前字符检查
	lineTextBeforeCursor := c.GetLineTextBeforeCursor()
	lineTextBeforeCursor = strings.TrimSpace(lineTextBeforeCursor)
	for _, char := range cursorBeforeStopChars {
		if strings.HasSuffix(lineTextBeforeCursor, char) {
			return false
		}
	}

	// 光标后字符检查
	lineTextAfterCursor := c.GetLineTextAfterCursor()
	lineTextAfterCursor = strings.TrimSpace(lineTextAfterCursor)
	if len(lineTextAfterCursor) > 0 {
		for _, char := range cursorAfterAllowChars {
			if strings.HasPrefix(lineTextAfterCursor, char) {
				return true
			}
		}
		return false
	}

	return true
}
