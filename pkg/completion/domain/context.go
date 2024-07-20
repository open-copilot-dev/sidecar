package domain

import (
	"open-copilot.dev/sidecar/pkg/common"
)

type CompletionContext struct {
	Ctx     *common.CancelableContext
	Request *CompletionRequest
}

func (c *CompletionContext) IsCanceled() bool {
	return c.Ctx.IsCanceled()
}
