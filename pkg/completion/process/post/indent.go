package post

import (
	"open-copilot.dev/sidecar/pkg/completion/context"
	"open-copilot.dev/sidecar/pkg/util"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////////////
// 后处理：缩进修正

type IndentPostProcessor struct {
}

func (m *IndentPostProcessor) process(c *context.CompletionContext, modelText string) string {
	lines := strings.Split(modelText, "\n")
	for i, line := range lines {
		if i == 0 {
			lineTextBeforeCursor := c.GetLineTextBeforeCursor()
			if util.IsBlank(lineTextBeforeCursor) && len(lineTextBeforeCursor) == util.CalcIndent(line) {
				// 首行的缩进是匹配的，则整体都不进行修正了
				lines[0] = lines[0][len(lineTextBeforeCursor):]
				return strings.Join(lines, "\n")
			}
			continue
		}
		lines[i] = strings.Repeat(" ", c.Request.CompletionLine.NextLineIndent) + line
	}
	return strings.Join(lines, "\n")
}
