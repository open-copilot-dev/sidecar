package domain

type ChatRequest struct {
	ChatID  string `json:"chatID"`
	Content string `json:"content"`
}

type ChatStreamResult struct {
	ChatID     string `json:"chatID"`
	MessageID  string `json:"messageID"`
	Index      int    `json:"index"`
	Content    string `json:"content"`
	IsFinished bool   `json:"isFinished"`
}
