package chat

import (
	"errors"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"io"
	"open-copilot.dev/sidecar/pkg/domain"
	"open-copilot.dev/sidecar/pkg/llm"
	"open-copilot.dev/sidecar/pkg/util"
	"path/filepath"
	"strings"
	"time"
)

var chatStore Store = NewLocalStore(filepath.Join(domain.BaseDir, "data/chats"))

func ProcessRequest(ctx *domain.CancelableContext, request *domain.ChatRequest,
	onStreamResult func(streamResult *domain.ChatStreamResult)) error {
	if strings.TrimSpace(request.Content) == "" {
		return errors.New("empty content")
	}

	// 获取chat信息
	chat, err := chatStore.GetChat(request.ChatID)
	if err != nil {
		hlog.CtxErrorf(ctx, "get chat err: %v", err)
	}
	if chat == nil {
		chat = &domain.Chat{
			ChatID:   request.ChatID,
			Messages: make([]*domain.ChatMessage, 0),
		}
	}
	chat.Messages = append(chat.Messages, &domain.ChatMessage{
		MessageID: request.MessageID,
		DateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Content:   request.Content,
		Role:      domain.ChatMessageRoleUser,
	})
	if chat.Title == "" {
		chat.Title = util.TruncateString(chat.Messages[0].Content, 50)
	}

	modelMessages := make([]*domain.ChatCompletionMessage, 0)
	for _, message := range chat.Messages {
		modelMessages = append(modelMessages, &domain.ChatCompletionMessage{
			Role:    message.Role,
			Content: &domain.ChatCompletionMessageContent{StringValue: &message.Content},
		})
	}

	if ctx.IsCanceled() {
		return domain.ErrCanceled
	}
	client := llm.GetClient(request.LlmClientName)
	modelStreamResponse, err := client.CreateChatCompletionStream(ctx, &domain.ChatCompletionRequest{
		Messages: modelMessages,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "Failed to chat completion request: %v", err)
		return err
	}
	var responseID = ""
	var content = ""
	var index int
	for {
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
		responseID = resp.ID
		for _, choice := range resp.Choices {
			index = choice.Index
			onStreamResult(&domain.ChatStreamResult{
				ChatID:     request.ChatID,
				MessageID:  request.MessageID,
				ResponseID: responseID,
				Index:      choice.Index,
				Content:    choice.Delta.Content,
				IsFinished: false,
			})
			content += choice.Delta.Content
		}
	}
	onStreamResult(&domain.ChatStreamResult{
		ChatID:     request.ChatID,
		MessageID:  request.MessageID,
		ResponseID: responseID,
		Index:      index + 1,
		Content:    "",
		IsFinished: true,
	})
	chat.Messages = append(chat.Messages, &domain.ChatMessage{
		MessageID: request.MessageID,
		DateTime:  time.Now().Format("2006-01-02 15:04:05"),
		Content:   content,
		Role:      volcModel.ChatMessageRoleAssistant,
	})
	chat.LastChatTime = chat.Messages[len(chat.Messages)-1].DateTime
	err = chatStore.SaveChat(chat)
	if err != nil {
		hlog.CtxErrorf(ctx.Context, "save chat failed, err: %v", err)
	}
	return nil
}

func ProcessDetailRequest(ctx *domain.CancelableContext, chatID string) (*domain.Chat, error) {
	return chatStore.GetChat(chatID)
}

func ProcessListRequest(ctx *domain.CancelableContext) ([]*domain.Chat, error) {
	_, chats, err := chatStore.ListChats(1, 100)
	if err != nil {
		hlog.CtxErrorf(ctx.Context, "list chats failed, err: %v", err)
		return nil, err
	}
	// 将不需要的内容置空，减少传输量
	for _, chat := range chats {
		chat.System = ""
		chat.Messages = nil
	}
	return chats, nil
}

func ProcessDeleteMessageRequest(ctx *domain.CancelableContext, chatID string, messageID string) error {
	return chatStore.DeleteChatMessage(chatID, messageID)
}

func ProcessDeleteRequest(ctx *domain.CancelableContext, chatID string) error {
	return chatStore.DeleteChat(chatID)
}

func ProcessDeleteAllRequest(ctx *domain.CancelableContext) error {
	return chatStore.DeleteAllChats()
}
