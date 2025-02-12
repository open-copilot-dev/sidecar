package domain

type ChatRequest struct {
	LlmClientName string `json:"llmClientName"`
	ChatID        string `json:"chatID"`
	MessageID     string `json:"messageID"`
	Content       string `json:"content"`
}

type ChatStreamResult struct {
	ChatID     string `json:"chatID"`
	MessageID  string `json:"messageID"`
	ResponseID string `json:"responseID"`
	Index      int    `json:"index"`
	Content    string `json:"content"`
	IsFinished bool   `json:"isFinished"`
}

type Chat struct {
	ChatID       string         `json:"chatID"`       // chat id
	Title        string         `json:"title"`        // chat title
	LastChatTime string         `json:"lastChatTime"` // last chat time
	ProjectName  string         `json:"projectName"`  // project name
	System       string         `json:"system"`       // system prompt
	Messages     []*ChatMessage `json:"messages"`     // chat messages
}

type ChatMessage struct {
	MessageID string `json:"messageID"` // message id
	DateTime  string `json:"dateTime"`  // message create time
	Content   string `json:"content"`   // message content
	Role      string `json:"role"`      // user or assistant
}
