package llm

import (
	"context"
	"open-copilot.dev/sidecar/pkg/domain"
)

type Client interface {
	CreateChatCompletion(ctx context.Context, request *domain.ChatCompletionRequest) (response *domain.ChatCompletionResponse, err error)
	CreateChatCompletionStream(ctx context.Context, request *domain.ChatCompletionRequest) (stream domain.ChatCompletionStreamReader, err error)
}
