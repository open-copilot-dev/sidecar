package domain

type Chat struct {
	ChatID   string         `json:"chatID"`   // chat id
	Title    string         `json:"title"`    // chat title
	System   string         `json:"system"`   // system prompt
	Messages []*ChatMessage `json:"messages"` // chat messages
}

type ChatMessage struct {
	Content string `json:"content"` // message content
	Role    string `json:"role"`    // user or assistant
}
