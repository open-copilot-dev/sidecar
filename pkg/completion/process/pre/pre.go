package pre

import "open-copilot.dev/sidecar/pkg/completion/domain"

type Processor interface {
	process(c *domain.CompletionContext) bool
}

var processors = []Processor{
	&FilterProcessor{},
}

func Process(c *domain.CompletionContext) bool {
	for _, processor := range processors {
		if !processor.process(c) {
			return false
		}
	}
	return true
}
