package pre

import "open-copilot.dev/sidecar/pkg/completion/domain"

type Processor interface {
	process(c *domain.CompletionContext) State
}

var processors = []Processor{
	&FilterPreProcessor{},
	&GrammarPreProcessor{},
}

type State bool

var StateContinue State = true
var StateStop State = false

func Process(c *domain.CompletionContext) State {
	for _, processor := range processors {
		if !processor.process(c) {
			return false
		}
	}
	return true
}
