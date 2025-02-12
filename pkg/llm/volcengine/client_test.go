package volcengine

import (
	"context"
	"fmt"
	"open-copilot.dev/sidecar/pkg/domain"
	"os"
	"testing"
)

var client = NewClient(os.Getenv("VOLC_API_KEY"), os.Getenv("VOLC_MODEL"))

func TestClient_CreateChatCompletion(t *testing.T) {
	ctx := context.Background()

	fmt.Println("----- standard request -----")
	req := &domain.ChatCompletionRequest{
		Messages: []*domain.ChatCompletionMessage{
			{
				Role:    domain.ChatMessageRoleSystem,
				Content: domain.NewStringMessageContent("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
			},
			{
				Role:    domain.ChatMessageRoleUser,
				Content: domain.NewStringMessageContent("常见的十字花科植物有哪些？"),
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("standard chat error: %v\n", err)
		return
	}
	fmt.Println(*resp.Choices[0].Message.Content.StringValue)

}

func TestClient_CreateChatCompletionStream(t *testing.T) {

}
