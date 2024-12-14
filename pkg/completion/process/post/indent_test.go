package post

import (
	"github.com/stretchr/testify/assert"
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"testing"
)

func TestIndentProcessor_process(t *testing.T) {
	m := &IndentProcessor{}
	c := &domain.CompletionContext{}
	c.Request = &domain.CompletionRequest{
		CompletionLine: &domain.CompletionLine{
			CurrentLineStartOffset: 0,
			CurrentCursorOffset:    8,
			LineText:               "        ",
			NextLineIndent:         8,
		},
	}

	modelText := "        case 0:\n            System.out.println(\"i 是 0\");\n            break;\n        case 1:\n            System.out.println(\"i 是 1\");\n            break;\n        case 2:\n            System.out.println(\"i 是 2 再次确认\");\n            break;\n        case 3:\n            System.out.println(\"i 是 3\");\n            break;\n        case 4:\n            System.out.println(\"i 是 4\");\n            break;\n        default:\n            System.out.println(\"i 不在 0 到 4 之间\");"
	processedModelText := m.process(c, modelText)
	wantedModelText := modelText
	assert.Equal(t, wantedModelText, processedModelText)
}
