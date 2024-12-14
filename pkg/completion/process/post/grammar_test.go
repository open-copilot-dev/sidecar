package post

import (
	"github.com/stretchr/testify/assert"
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"testing"
)

func TestGrammarPostProcessor_process(t *testing.T) {
	m := &GrammarPostProcessor{}
	c := &domain.CompletionContext{}
	c.Request = &domain.CompletionRequest{
		Language:         "Java",
		TextBeforeCursor: "",
		TextAfterCursor:  "",
		CompletionLine: &domain.CompletionLine{
			CurrentLineStartOffset: 0,
			CurrentCursorOffset:    5,
			LineText:               "switc",
		},
	}

	modelText := "switch (i) {\n    case 0:\n        System.out.println(\\\"i 是 0\\\");\n        break;\n    case 1:\n        System.out.println(\\\"i 是 1\\\");\n        break;\n    case 3:\n        System.out.println(\\\"i 是 3\\\");\n        break;\n    case 4:\n        System.out.println(\\\"i 是 4\\\");\n        break;\n    default:\n        System.out.println(\\\"其他情况\\\");\n}\n"
	processedModelText := m.process(c, modelText)
	wantedModelText := modelText[5:]
	assert.Equal(t, wantedModelText, processedModelText)
}
