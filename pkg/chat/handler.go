package chat

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	chatDomain "open-copilot.dev/sidecar/pkg/chat/domain"
	"open-copilot.dev/sidecar/pkg/common"
	"open-copilot.dev/sidecar/pkg/engine/volcengine"
)

func ProcessRequest(ctx *common.CancelableContext, request *chatDomain.ChatRequest,
	onStreamResult func(streamResult *chatDomain.ChatStreamResult)) error {
	modelMessages := make([]*volcModel.ChatCompletionMessage, 0, len(request.Messages))
	for _, message := range request.Messages {
		role := volcModel.ChatMessageRoleSystem
		if message.Role == "user" {
			role = volcModel.ChatMessageRoleUser
		}
		modelMessages = append(modelMessages, &volcModel.ChatCompletionMessage{
			Role:    role,
			Content: &volcModel.ChatCompletionMessageContent{StringValue: &message.Content},
		})
	}

	if ctx.IsCanceled() {
		return common.ErrCanceled
	}
	client := volcengine.NewClient("659b5a99-0614-48ee-a04c-bee4d96d2e83")
	modelStreamResponse, err := client.CreateChatCompletionStream(ctx, volcModel.ChatCompletionRequest{
		Model:    "ep-20240703013553-wjlhr",
		Messages: modelMessages,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to chat completion request: %v", err)
		return err
	}
	for {
		if modelStreamResponse.IsFinished {
			break
		}
		resp, err := modelStreamResponse.Recv()
		if err != nil {
			hlog.CtxErrorf(ctx, "Failed to chat completion request: %v", err)
			return err
		}
		for _, choice := range resp.Choices {
			onStreamResult(&chatDomain.ChatStreamResult{
				Index:      choice.Index,
				Content:    choice.Delta.Content,
				IsFinished: modelStreamResponse.IsFinished,
			})
		}
	}
	return nil
}
