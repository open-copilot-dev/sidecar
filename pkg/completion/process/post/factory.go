package post

import (
	"open-copilot.dev/sidecar/pkg/completion/context"
)

type Processor interface {
	process(c *context.CompletionContext, modelText string) string
}

var processors = []Processor{
	&MarkdownPostProcessor{},
	&IndentPostProcessor{},
	&OverlapPostProcessor{},
	&GrammarPostProcessor{},
}

func Process(c *context.CompletionContext, modelText string) string {
	for _, processor := range processors {
		modelText = processor.process(c, modelText)
	}
	return modelText
}
