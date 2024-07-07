package completion

import (
	"context"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"open-copilot.dev/sidecar/pkg/completion/process/post"
	"open-copilot.dev/sidecar/pkg/completion/prompt"
	"open-copilot.dev/sidecar/pkg/engine/volcengine"
)

func ProcessRequest(ctx context.Context, request *domain.CompletionRequest) (*domain.CompletionResult, error) {
	c := &domain.CompletionContext{
		Ctx:     ctx,
		Request: request,
	}
	// 组装Prompt
	messages := prompt.Build(c)

	// 发起调用
	client := volcengine.NewClient("659b5a99-0614-48ee-a04c-bee4d96d2e83")
	modelResponse, err := client.CreateChatCompletion(ctx, volcModel.ChatCompletionRequest{
		Model:    "ep-20240703013553-wjlhr",
		Messages: messages,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to execute completion request: %v", err)
		return nil, err
	}

	// 处理返回结果
	choices := make([]*domain.CompletionChoice, 0, len(modelResponse.Choices))
	for _, modelChoice := range modelResponse.Choices {
		choice := convertModelChoice(c, modelChoice)
		if choice == nil {
			continue
		}
		choices = append(choices, choice)
	}

	return &domain.CompletionResult{
		Choices: choices,
	}, nil
}

func convertModelChoice(c *domain.CompletionContext, modelChoice *volcModel.ChatCompletionChoice) *domain.CompletionChoice {
	if modelChoice == nil || modelChoice.Message.Content == nil || modelChoice.Message.Content.StringValue == nil {
		return nil
	}
	modelText := *modelChoice.Message.Content.StringValue
	modelText = post.Process(c, modelText)
	edit := &domain.CompletionEdit{
		StartOffset: c.Request.CompletionLine.CurrentCursorOffset,
		EndOffset:   c.Request.CompletionLine.CurrentCursorOffset,
		Text:        modelText,
		Type:        domain.CompletionEditTypeInsert,
	}
	return &domain.CompletionChoice{Edits: []*domain.CompletionEdit{edit}}
}
