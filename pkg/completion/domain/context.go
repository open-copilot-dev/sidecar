package domain

import "context"

type CompletionContext struct {
	Ctx     context.Context
	Request *CompletionRequest
}
