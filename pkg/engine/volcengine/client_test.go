package volcengine

import (
	"context"
	"fmt"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"testing"
)

var client = NewClient("659b5a99-0614-48ee-a04c-bee4d96d2e83")

func TestClient_CreateChatCompletion(t *testing.T) {
	ctx := context.Background()

	fmt.Println("----- standard request -----")
	req := model.ChatCompletionRequest{
		Model: "ep-20240703013553-wjlhr",
		Messages: []*model.ChatCompletionMessage{
			{
				Role: model.ChatMessageRoleSystem,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("你是豆包，是由字节跳动开发的 AI 人工智能助手"),
				},
			},
			{
				Role: model.ChatMessageRoleUser,
				Content: &model.ChatCompletionMessageContent{
					StringValue: volcengine.String("常见的十字花科植物有哪些？"),
				},
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
