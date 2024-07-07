package post

import (
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"regexp"
	"strings"
)

type MarkdownProcessor struct {
}

var re1 = regexp.MustCompile("```[a-zA-Z]+\n(.*)\n```")
var re2 = regexp.MustCompile("```\n(.*)\n```")
var re3 = regexp.MustCompile("```(.*)```")

func (m *MarkdownProcessor) process(c *domain.CompletionContext, modelText string) string {
	startIndex := strings.Index(modelText, "```")
	if startIndex == -1 {
		return modelText
	}
	endIndex := strings.Index(modelText[startIndex+3:], "```")
	if endIndex == -1 {
		return modelText
	}

	codeBlockText := modelText[startIndex : startIndex+3+endIndex+3]

	found := re1.FindAllStringSubmatch(codeBlockText, -1)
	if len(found) > 0 {
		return found[0][1]
	}
	found = re2.FindAllStringSubmatch(codeBlockText, -1)
	if len(found) > 0 {
		return found[0][1]
	}
	found = re3.FindAllStringSubmatch(codeBlockText, -1)
	if len(found) > 0 {
		return found[0][1]
	}
	return modelText
}
