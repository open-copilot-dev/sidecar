package pre

import (
	"open-copilot.dev/sidecar/pkg/completion/context"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////////////
// 前处理：过滤掉不必要的补全场景

type FilterPreProcessor struct {
}

// 光标之前是这些字符时，停止补全
var cursorBeforeStopChars = []string{";", ")", "]", "}"}

// 光标之后如果有字符，并且是这些字符时，才允许补全
var cursorAfterAllowChars = []string{")", "]"}

func (f *FilterPreProcessor) process(c *context.CompletionContext) State {
	if c.IsCanceled() {
		return StateStop
	}

	// 光标前字符检查
	lineTextBeforeCursor := c.GetLineTextBeforeCursor()
	lineTextBeforeCursor = strings.TrimSpace(lineTextBeforeCursor)
	for _, char := range cursorBeforeStopChars {
		if strings.HasSuffix(lineTextBeforeCursor, char) {
			return StateStop
		}
	}

	// 光标后字符检查
	lineTextAfterCursor := c.GetLineTextAfterCursor()
	lineTextAfterCursor = strings.TrimSpace(lineTextAfterCursor)
	if len(lineTextAfterCursor) > 0 {
		for _, char := range cursorAfterAllowChars {
			if strings.HasPrefix(lineTextAfterCursor, char) {
				return StateContinue
			}
		}
		return StateStop
	}

	return StateContinue
}
