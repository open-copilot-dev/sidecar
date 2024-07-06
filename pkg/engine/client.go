package engine

import (
	"context"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	volcUtils "github.com/volcengine/volcengine-go-sdk/service/arkruntime/utils"
)

type Client interface {
	CreateChatCompletion(ctx context.Context, request volcModel.ChatCompletionRequest) (response volcModel.ChatCompletionResponse, err error)
	CreateChatCompletionStream(ctx context.Context, request volcModel.ChatCompletionRequest) (stream *volcUtils.ChatCompletionStreamReader, err error)
}
