package completion

import (
	"encoding/json"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"open-copilot.dev/sidecar/pkg/completion/context"
	"open-copilot.dev/sidecar/pkg/completion/process/post"
	"open-copilot.dev/sidecar/pkg/completion/process/pre"
	"open-copilot.dev/sidecar/pkg/completion/prompt"
	"open-copilot.dev/sidecar/pkg/domain"
	"open-copilot.dev/sidecar/pkg/llm"
)

func ProcessRequest(ctx *domain.CancelableContext, request *domain.CompletionRequest) (*domain.CompletionResult, error) {
	c := &context.CompletionContext{
		Ctx:     ctx,
		Request: request,
	}
	// 前置处理
	if c.IsCanceled() {
		return nil, domain.ErrCanceled
	}
	if !pre.Process(c) {
		return nil, domain.ErrIgnored
	}

	// 组装Prompt
	if c.IsCanceled() {
		return nil, domain.ErrCanceled
	}
	messages := prompt.Build(c)

	// 发起调用
	if c.IsCanceled() {
		return nil, domain.ErrCanceled
	}
	client := llm.GetClient(request.LlmClientName)
	modelResponse, err := client.CreateChatCompletion(ctx, &domain.ChatCompletionRequest{
		Messages: messages,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to execute completion request: %v", err)
		return nil, err
	}
	modelResponseBytes, _ := json.Marshal(modelResponse)
	hlog.CtxInfof(ctx, "modelResponse: %s", modelResponseBytes)

	// 处理返回结果
	if c.IsCanceled() {
		return nil, domain.ErrCanceled
	}
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

func convertModelChoice(c *context.CompletionContext, modelChoice *domain.ChatCompletionChoice) *domain.CompletionChoice {
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
