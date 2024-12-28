package chat

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"io"
	chatDomain "open-copilot.dev/sidecar/pkg/chat/domain"
	"open-copilot.dev/sidecar/pkg/common"
	"open-copilot.dev/sidecar/pkg/engine/volcengine"
	"strings"
)

func ProcessRequest(ctx *common.CancelableContext, request *chatDomain.ChatRequest,
	onStreamResult func(streamResult *chatDomain.ChatStreamResult)) error {
	if strings.TrimSpace(request.Content) == "" {
		return errors.New("empty content")
	}

	modelMessages := make([]*volcModel.ChatCompletionMessage, 0)
	modelMessages = append(modelMessages, &volcModel.ChatCompletionMessage{
		Role:    volcModel.ChatMessageRoleUser,
		Content: &volcModel.ChatCompletionMessageContent{StringValue: &request.Content},
	})

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
	var messageID string
	var index int
	for {
		if modelStreamResponse.IsFinished {
			break
		}
		if ctx.IsCanceled() {
			_ = modelStreamResponse.Close()
			break
		}
		resp, err := modelStreamResponse.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			hlog.CtxErrorf(ctx, "Failed to chat completion request: %v", err)
			return err
		}
		messageID = resp.ID
		for _, choice := range resp.Choices {
			index = choice.Index
			onStreamResult(&chatDomain.ChatStreamResult{
				ChatID:     request.ChatID,
				MessageID:  resp.ID,
				Index:      choice.Index,
				Content:    choice.Delta.Content,
				IsFinished: modelStreamResponse.IsFinished,
			})
		}
	}
	onStreamResult(&chatDomain.ChatStreamResult{
		ChatID:     request.ChatID,
		MessageID:  messageID,
		Index:      index + 1,
		Content:    "",
		IsFinished: true,
	})
	return nil
}
