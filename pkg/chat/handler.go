package chat

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"io"
	chatDomain "open-copilot.dev/sidecar/pkg/chat/domain"
	"open-copilot.dev/sidecar/pkg/common"
	"open-copilot.dev/sidecar/pkg/engine/volcengine"
	"open-copilot.dev/sidecar/pkg/util"
	"path/filepath"
	"strings"
)

var chatStore Store = NewLocalStore(filepath.Join(common.BaseDir, "data/chats"))

func ProcessRequest(ctx *common.CancelableContext, request *chatDomain.ChatRequest,
	onStreamResult func(streamResult *chatDomain.ChatStreamResult)) error {
	if strings.TrimSpace(request.Content) == "" {
		return errors.New("empty content")
	}

	// 获取chat信息
	chat, err := chatStore.GetChat(request.ChatID)
	if err != nil {
		hlog.CtxErrorf(ctx, "get chat err: %v", err)
	}
	if chat == nil {
		chat = &chatDomain.Chat{
			ChatID:   request.ChatID,
			Messages: make([]*chatDomain.ChatMessage, 0),
		}
	}
	chat.Messages = append(chat.Messages, &chatDomain.ChatMessage{
		Content: request.Content,
		Role:    volcModel.ChatMessageRoleUser,
	})
	if chat.Title == "" {
		chat.Title = util.TruncateString(chat.Messages[0].Content, 20)
	}

	modelMessages := make([]*volcModel.ChatCompletionMessage, 0)
	for _, message := range chat.Messages {
		modelMessages = append(modelMessages, &volcModel.ChatCompletionMessage{
			Role:    message.Role,
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
	var content = ""
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
			content += choice.Delta.Content
		}
	}
	onStreamResult(&chatDomain.ChatStreamResult{
		ChatID:     request.ChatID,
		MessageID:  messageID,
		Index:      index + 1,
		Content:    "",
		IsFinished: true,
	})
	chat.Messages = append(chat.Messages, &chatDomain.ChatMessage{
		Content: content,
		Role:    volcModel.ChatMessageRoleAssistant,
	})
	err = chatStore.SaveChat(chat)
	if err != nil {
		hlog.CtxErrorf(ctx.Context, "save chat failed, err: %v", err)
	}
	return nil
}
