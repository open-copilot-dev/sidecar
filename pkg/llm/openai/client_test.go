package openai

import (
	"context"
	"open-copilot.dev/sidecar/pkg/domain"
	"os"
	"testing"
)

var client = NewClient(os.Getenv("OPENAI_API_KEY"), "gpt-3.5-turbo")

func TestClient_CreateChatCompletion(t *testing.T) {
	request := &domain.ChatCompletionRequest{
		Messages: []*domain.ChatCompletionMessage{
			{
				Role:    domain.ChatMessageRoleUser,
				Content: domain.NewStringMessageContent("hello"),
			},
		},
		MaxTokens:   100,
		Temperature: 0.9,
		TopP:        1,
	}
	response, err := client.CreateChatCompletion(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(response)
}

func TestClient_CreateChatCompletionStream(t *testing.T) {

}
