package pre

import (
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"strings"
)

type FilterProcessor struct {
}

var stopChars = []string{";", ")", "]", "}"}

func (f *FilterProcessor) process(c *domain.CompletionContext) bool {
	lineText := strings.TrimSpace(c.Request.CompletionLine.LineText)
	for _, char := range stopChars {
		if strings.HasSuffix(lineText, char) {
			return false
		}
	}
	return true
}
