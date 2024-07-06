package completion

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"open-copilot.dev/sidecar/pkg/engine/volcengine"
)

func ProcessRequest(ctx context.Context, request *Request) (*Result, error) {
	// 组装Prompt
	promptBuilder := NewPromptBuilder(request)
	messages := promptBuilder.Build()

	client := volcengine.NewClient("659b5a99-0614-48ee-a04c-bee4d96d2e83")
	modelResponse, err := client.CreateChatCompletion(ctx, volcModel.ChatCompletionRequest{
		Model:    "ep-20240703013553-wjlhr",
		Messages: messages,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to execute completion request: %v", err)
		return nil, err
	}

	choices := make([]*Choice, 0, len(modelResponse.Choices))
	for _, modelChoice := range modelResponse.Choices {
		choice := convertModelChoice(modelChoice, request)
		if choice == nil {
			continue
		}
		choices = append(choices, choice)
	}

	return &Result{
		Choices: choices,
	}, nil
}

func convertModelChoice(modelChoice *volcModel.ChatCompletionChoice, request *Request) *Choice {
	if modelChoice == nil || modelChoice.Message.Content == nil || modelChoice.Message.Content.StringValue == nil {
		return nil
	}
	edit := &Edit{
		StartOffset: request.CompletionLine.CurrentCursorOffset,
		EndOffset:   request.CompletionLine.CurrentCursorOffset,
		Text:        *modelChoice.Message.Content.StringValue,
	}
	return &Choice{Edits: []*Edit{edit}}
}
