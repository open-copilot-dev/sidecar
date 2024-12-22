package domain

type Chat struct {
	ChatID   string         `json:"chatID"`
	System   string         `json:"system"`
	Messages []*ChatMessage `json:"messages"`
}

type ChatMessage struct {
	Content string `json:"content"`
	Role    string `json:"role"`
}
