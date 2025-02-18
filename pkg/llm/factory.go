package llm

import (
	"open-copilot.dev/sidecar/pkg/domain"
	"open-copilot.dev/sidecar/pkg/llm/openai"
	"open-copilot.dev/sidecar/pkg/llm/volcengine"
)

func NewClient(llmSetting *domain.LlmSetting) Client {
	switch llmSetting.Platform {
	case "openai":
		return openai.NewClient(llmSetting.ApiKey, llmSetting.Model)
	case "volcengine":
		return volcengine.NewClient(llmSetting.ApiKey, llmSetting.Model)
	default:
		return openai.NewClient(llmSetting.ApiKey, llmSetting.Model)
	}
}

func GetClient(name string) Client {
	// TODO
	return volcengine.NewClient("659b5a99-0614-48ee-a04c-bee4d96d2e83", "ep-20240703013553-wjlhr")
}
