package post

import (
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"regexp"
	"strings"
)

type MarkdownProcessor struct {
}

var re1 = regexp.MustCompile("```[a-zA-Z]+\n([\\s\\S]*)\n```")
var re2 = regexp.MustCompile("```\n([\\s\\S]*)\n```")
var re3 = regexp.MustCompile("```([\\s\\S]*)```")
var re4 = regexp.MustCompile("`([\\s\\S]*)`")

func (m *MarkdownProcessor) process(c *domain.CompletionContext, modelText string) string {
	blockStartIndex := strings.Index(modelText, "```")
	if blockStartIndex == -1 {
		inlineStartIndex := strings.Index(modelText, "`")
		if inlineStartIndex == -1 {
			return modelText
		}
		inlineEndIndex := strings.Index(modelText[inlineStartIndex+1:], "`")
		if inlineEndIndex == -1 {
			return modelText
		}

		codeInlineText := modelText[inlineStartIndex : inlineStartIndex+1+inlineEndIndex+1]

		found := re4.FindAllStringSubmatch(codeInlineText, -1)
		if len(found) > 0 {
			return found[0][1]
		}
		return modelText
	}
	blockEndIndex := strings.Index(modelText[blockStartIndex+3:], "```")
	if blockEndIndex == -1 {
		return modelText
	}

	codeBlockText := modelText[blockStartIndex : blockStartIndex+3+blockEndIndex+3]

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
