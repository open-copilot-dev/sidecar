package volcengine

import (
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"open-copilot.dev/sidecar/pkg/domain"
)

type Client struct {
	model string
	cli   *arkruntime.Client
}

func NewClient(apiKey string, model string) *Client {
	cli := arkruntime.NewClientWithApiKey(apiKey)
	return &Client{
		cli:   cli,
		model: model,
	}
}

func (c *Client) CreateChatCompletion(ctx context.Context, request *domain.ChatCompletionRequest) (response *domain.ChatCompletionResponse, err error) {
	completion, err := c.cli.CreateChatCompletion(ctx, convertRequest(request, c.model))
	return convertResponse(completion), err
}

func (c *Client) CreateChatCompletionStream(ctx context.Context, request *domain.ChatCompletionRequest) (stream domain.ChatCompletionStreamReader, err error) {
	completionStream, err := c.cli.CreateChatCompletionStream(ctx, convertRequest(request, c.model))
	return convertStreamReader(completionStream), err
}
