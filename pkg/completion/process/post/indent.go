package post

import (
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"strings"
)

type IndentProcessor struct {
}

func (m *IndentProcessor) process(c *domain.CompletionContext, modelText string) string {
	if c.Request.CompletionLine.NextLineIndent <= 0 {
		return modelText
	}
	lines := strings.Split(modelText, "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}
		lines[i] = strings.Repeat(" ", c.Request.CompletionLine.NextLineIndent) + line
	}
	return strings.Join(lines, "\n")
}
