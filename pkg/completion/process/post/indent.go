package post

import (
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"open-copilot.dev/sidecar/pkg/util"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////////////
// 缩进修正

type IndentProcessor struct {
}

func (m *IndentProcessor) process(c *domain.CompletionContext, modelText string) string {
	if c.Request.CompletionLine.NextLineIndent <= 0 {
		return modelText
	}
	lines := strings.Split(modelText, "\n")
	for i, line := range lines {
		if i == 0 {
			lineTextBeforeCursor := c.GetLineTextBeforeCursor()
			if util.IsBlank(lineTextBeforeCursor) && len(lineTextBeforeCursor) == util.CalcIndent(line) {
				// 首行的缩进是匹配的，则整体都不进行修正了
				return modelText
			}
			continue
		}
		lines[i] = strings.Repeat(" ", c.Request.CompletionLine.NextLineIndent) + line
	}
	return strings.Join(lines, "\n")
}
