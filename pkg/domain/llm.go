package domain

// chat message roles
const (
	ChatMessageRoleSystem    = "system"
	ChatMessageRoleUser      = "user"
	ChatMessageRoleAssistant = "assistant"
)

type ChatCompletionRequest struct {
	Messages    []*ChatCompletionMessage `json:"messages"`
	MaxTokens   int                      `json:"max_tokens,omitempty"`
	Temperature float32                  `json:"temperature,omitempty"`
	TopP        float32                  `json:"top_p,omitempty"`
	Stop        []string                 `json:"stop,omitempty"`
	N           int                      `json:"n,omitempty"`
}

type ChatCompletionMessage struct {
	Role    string                        `json:"role"`
	Content *ChatCompletionMessageContent `json:"content"`
}

type ChatCompletionMessageContent struct {
	StringValue *string
}

type ChatCompletionResponse struct {
	ID      string                  `json:"id"`
	Object  string                  `json:"object"`
	Created int64                   `json:"created"`
	Model   string                  `json:"llm"`
	Choices []*ChatCompletionChoice `json:"choices"`
	Usage   *ChatPromptUsage        `json:"usage"`
}

type ChatCompletionChoice struct {
	Index        int                    `json:"index"`
	Message      *ChatCompletionMessage `json:"message"`
	FinishReason string                 `json:"finish_reason"`
}

type ChatPromptUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ChatCompletionStreamReader interface {
	Recv() (response *ChatCompletionStreamResponse, err error)
	Close() error
}

type ChatCompletionStreamResponse struct {
	ID      string                        `json:"id"`
	Object  string                        `json:"object"`
	Created int64                         `json:"created"`
	Model   string                        `json:"llm"`
	Choices []*ChatCompletionStreamChoice `json:"choices"`
	// An optional field that will only be present when you set stream_options: {"include_usage": true} in your request.
	// When present, it contains a null value except for the last chunk which contains the token usage statistics
	// for the entire request.
	Usage *ChatPromptUsage `json:"usage,omitempty"`
}

type ChatCompletionStreamChoice struct {
	Index        int                              `json:"index"`
	Delta        *ChatCompletionStreamChoiceDelta `json:"delta"`
	FinishReason string                           `json:"finish_reason"`
}

type ChatCompletionStreamChoiceDelta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}

func NewStringMessageContent(s string) *ChatCompletionMessageContent {
	return &ChatCompletionMessageContent{
		StringValue: &s,
	}
}
