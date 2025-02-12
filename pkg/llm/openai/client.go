package openai

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"open-copilot.dev/sidecar/pkg/domain"
	"open-copilot.dev/sidecar/pkg/util"
)

type Client struct {
	model string
	cli   *openai.Client
}

func NewClient(apiKey string, model string) *Client {
	cli := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &Client{
		model: model,
		cli:   cli,
	}
}

func (c *Client) CreateChatCompletion(ctx context.Context, request *domain.ChatCompletionRequest) (response *domain.ChatCompletionResponse, err error) {
	messages := convertMessages(request.Messages)
	chatCompletion, err := c.cli.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(c.model),
	})
	if err != nil {
		return nil, err
	}
	return &domain.ChatCompletionResponse{
		ID:      chatCompletion.ID,
		Object:  util.TryToJson(chatCompletion.Object),
		Created: chatCompletion.Created,
		Model:   chatCompletion.Model,
		Choices: convertChoices(chatCompletion.Choices),
	}, nil
}

func (c *Client) CreateChatCompletionStream(ctx context.Context, request *domain.ChatCompletionRequest) (domain.ChatCompletionStreamReader, error) {
	stream := c.cli.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(convertMessages(request.Messages)),
		Model:    openai.F(c.model),
	})

	return convertStreamReader(stream)
}
