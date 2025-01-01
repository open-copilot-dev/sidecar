package domain

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
