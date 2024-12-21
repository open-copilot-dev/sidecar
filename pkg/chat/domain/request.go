package domain

type ChatRequest struct {
	UUID     string         `json:"uuid"`
	Messages []*ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}

type ChatStreamResult struct {
	Index      int    `json:"index"`
	Content    string `json:"content"`
	IsFinished bool   `json:"is_finished"`
}
