package llm

import (
	"open-copilot.dev/sidecar/pkg/domain"
	"open-copilot.dev/sidecar/pkg/llm/openai"
	"open-copilot.dev/sidecar/pkg/llm/volcengine"
)

func NewClient(cfg *domain.LlmConfig) Client {
	switch cfg.Platform {
	case "openai":
		return openai.NewClient(cfg.ApiKey, cfg.Model)
	case "volcengine":
		return volcengine.NewClient(cfg.ApiKey, cfg.Model)
	default:
		return openai.NewClient(cfg.ApiKey, cfg.Model)
	}
}

func GetClient(name string) Client {
	// TODO
	return volcengine.NewClient("659b5a99-0614-48ee-a04c-bee4d96d2e83", "ep-20240703013553-wjlhr")
}
