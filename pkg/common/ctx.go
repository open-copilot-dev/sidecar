package common

import (
	"context"
	"sync/atomic"
)

type CancelableContext struct {
	context.Context
	cancelFlag atomic.Bool
}

func NewCancelableContext(parent context.Context) *CancelableContext {
	return &CancelableContext{
		Context:    parent,
		cancelFlag: atomic.Bool{},
	}
}

func (ctx *CancelableContext) IsCanceled() bool {
	return ctx.cancelFlag.Load()
}
func (ctx *CancelableContext) Cancel() {
	ctx.cancelFlag.Store(true)
}
