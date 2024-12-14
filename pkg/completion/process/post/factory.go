package post

import "open-copilot.dev/sidecar/pkg/completion/domain"

type Processor interface {
	process(c *domain.CompletionContext, modelText string) string
}

var processors = []Processor{
	&MarkdownPostProcessor{},
	&IndentPostProcessor{},
	&OverlapPostProcessor{},
	&GrammarPostProcessor{},
}

func Process(c *domain.CompletionContext, modelText string) string {
	for _, processor := range processors {
		modelText = processor.process(c, modelText)
	}
	return modelText
}
