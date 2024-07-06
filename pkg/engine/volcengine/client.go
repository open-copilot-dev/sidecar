package volcengine

import (
	"context"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	volcUtils "github.com/volcengine/volcengine-go-sdk/service/arkruntime/utils"
)

type Client struct {
	cli *arkruntime.Client
}

func NewClient(apiKey string) *Client {
	cli := arkruntime.NewClientWithApiKey(apiKey)
	return &Client{
		cli: cli,
	}
}

func (c *Client) CreateChatCompletion(ctx context.Context, request volcModel.ChatCompletionRequest) (response volcModel.ChatCompletionResponse, err error) {
	completion, err := c.cli.CreateChatCompletion(ctx, request)
	return completion, err
}

func (c *Client) CreateChatCompletionStream(ctx context.Context, request volcModel.ChatCompletionRequest) (stream *volcUtils.ChatCompletionStreamReader, err error) {
	completionStream, err := c.cli.CreateChatCompletionStream(ctx, request)
	return completionStream, err
}
