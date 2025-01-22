package openai

import (
	"context"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	volcUtils "github.com/volcengine/volcengine-go-sdk/service/arkruntime/utils"
	"open-copilot.dev/sidecar/pkg/util"
)

type Client struct {
	cli *openai.Client
}

func NewClient(apiKey string) *Client {
	cli := openai.NewClient(
		option.WithAPIKey(apiKey),
	)
	return &Client{
		cli: cli,
	}
}

func (c *Client) CreateChatCompletion(ctx context.Context, request volcModel.ChatCompletionRequest) (response volcModel.ChatCompletionResponse, err error) {
	messages := c.convertMessages(request.Messages)
	chatCompletion, err := c.cli.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(request.Model),
	})
	if err != nil {
		return volcModel.ChatCompletionResponse{}, err
	}
	return volcModel.ChatCompletionResponse{
		ID:         chatCompletion.ID,
		Object:     util.TryToJson(chatCompletion.Object),
		Created:    chatCompletion.Created,
		Model:      chatCompletion.Model,
		Choices:    c.convertChoices(chatCompletion.Choices),
		Usage:      volcModel.Usage{},
		HttpHeader: nil,
	}, nil
}

func (c *Client) CreateChatCompletionStream(ctx context.Context, request volcModel.ChatCompletionRequest) (*volcUtils.ChatCompletionStreamReader, error) {
	stream := c.cli.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(c.convertMessages(request.Messages)),
		Model:    openai.F(request.Model),
	})

	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			println(chunk.Choices[0].Delta.Content)
		}
	}

	if err := stream.Err(); err != nil {
		panic(err)
	}

}

func (c *Client) convertMessages(messages []*volcModel.ChatCompletionMessage) []openai.ChatCompletionMessageParamUnion {
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	for _, message := range messages {
		openaiMessage := c.convertMessage(message)
		openaiMessages = append(openaiMessages, openaiMessage)
	}
	return openaiMessages
}

func (c *Client) convertMessage(message *volcModel.ChatCompletionMessage) openai.ChatCompletionMessageParamUnion {
	var openaiMessage openai.ChatCompletionMessageParamUnion = nil
	switch message.Role {
	case volcModel.ChatMessageRoleUser:
		openaiMessage = openai.UserMessage(*message.Content.StringValue)
	case volcModel.ChatMessageRoleAssistant:
		openaiMessage = openai.AssistantMessage(*message.Content.StringValue)
	case volcModel.ChatMessageRoleSystem:
		openaiMessage = openai.SystemMessage(*message.Content.StringValue)
	}
	return openaiMessage
}

func (c *Client) convertChoices(choices []openai.ChatCompletionChoice) []*volcModel.ChatCompletionChoice {
	volcChoices := make([]*volcModel.ChatCompletionChoice, 0, len(choices))
	for _, choice := range choices {
		volcChoices = append(volcChoices, &volcModel.ChatCompletionChoice{
			Index:        int(choice.Index),
			Message:      c.convertChoiceMessage(choice.Message),
			FinishReason: c.convertChoiceFinishReason(choice.FinishReason),
			LogProbs:     c.convertChoiceLogProbs(choice.Logprobs),
		})
	}
	return volcChoices
}

func (c *Client) convertChoiceMessage(choiceMessage openai.ChatCompletionMessage) volcModel.ChatCompletionMessage {
	volcChoiceMessage := volcModel.ChatCompletionMessage{
		Role: string(choiceMessage.Role),
		Content: &volcModel.ChatCompletionMessageContent{
			StringValue: &choiceMessage.Content,
		},
		FunctionCall: &volcModel.FunctionCall{
			Name:      choiceMessage.FunctionCall.Name,
			Arguments: choiceMessage.FunctionCall.Arguments,
		},
		ToolCalls:  nil,
		ToolCallID: "",
	}
	if choiceMessage.ToolCalls != nil {
		volcChoiceMessage.ToolCalls = make([]*volcModel.ToolCall, 0, len(choiceMessage.ToolCalls))
		for _, toolCall := range choiceMessage.ToolCalls {
			volcChoiceMessage.ToolCalls = append(volcChoiceMessage.ToolCalls, &volcModel.ToolCall{
				ID:   toolCall.ID,
				Type: volcModel.ToolType(toolCall.Type),
				Function: volcModel.FunctionCall{
					Name:      toolCall.Function.Name,
					Arguments: toolCall.Function.Arguments,
				},
			})
		}
	}
	return volcChoiceMessage
}

func (c *Client) convertChoiceFinishReason(finishReason openai.ChatCompletionChoicesFinishReason) volcModel.FinishReason {
	return volcModel.FinishReason(finishReason)
}

func (c *Client) convertChoiceLogProbs(logProbs openai.ChatCompletionChoicesLogprobs) *volcModel.LogProbs {
	volcLogProbs := &volcModel.LogProbs{
		Content: make([]*volcModel.LogProb, 0, len(logProbs.Content)),
	}
	for _, content := range logProbs.Content {
		volcLogProbs.Content = append(volcLogProbs.Content, &volcModel.LogProb{
			Token:       content.Token,
			LogProb:     content.Logprob,
			Bytes:       c.convertRune(content.Bytes),
			TopLogProbs: c.convertTopLogprobs(content.TopLogprobs),
		})
	}
	return volcLogProbs
}

func (c *Client) convertTopLogprobs(logprobs []openai.ChatCompletionTokenLogprobTopLogprob) []*volcModel.TopLogProbs {
	volcTopLogprobs := make([]*volcModel.TopLogProbs, 0, len(logprobs))
	for _, logprob := range logprobs {
		volcTopLogprobs = append(volcTopLogprobs, &volcModel.TopLogProbs{
			Bytes:   c.convertRune(logprob.Bytes),
			LogProb: logprob.Logprob,
			Token:   logprob.Token,
		})
	}
	return volcTopLogprobs
}

func (c *Client) convertRune(bytes []int64) []rune {
	runeSlice := make([]rune, len(bytes))
	for i, v := range bytes {
		runeSlice[i] = rune(v)
	}
	return runeSlice
}
