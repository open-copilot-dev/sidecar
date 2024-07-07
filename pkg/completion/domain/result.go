package domain

// CompletionResult 代码补全结果
type CompletionResult struct {
	Choices []*CompletionChoice `json:"choices"`
}

type CompletionChoice struct {
	Edits []*CompletionEdit `json:"edits"`
}

const (
	CompletionEditTypeInsert = "INSERT"
	CompletionEditTypeDelete = "DELETE"
)

type CompletionEdit struct {
	StartOffset int    `json:"startOffset"`
	EndOffset   int    `json:"endOffset"`
	Text        string `json:"text"`
	Type        string `json:"type"`
}
