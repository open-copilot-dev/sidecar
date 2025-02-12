package pre

import (
	"open-copilot.dev/sidecar/pkg/completion/context"
)

type Processor interface {
	process(c *context.CompletionContext) State
}

var processors = []Processor{
	&FilterPreProcessor{},
	&GrammarPreProcessor{},
}

type State bool

var StateContinue State = true
var StateStop State = false

func Process(c *context.CompletionContext) State {
	for _, processor := range processors {
		if !processor.process(c) {
			return false
		}
	}
	return true
}
