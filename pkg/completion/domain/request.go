package domain

// CompletionRequest 代码补全请求
type CompletionRequest struct {
	UUID             string          `json:"uuid"`
	ProjectPath      string          `json:"projectPath"`
	DocPath          string          `json:"docPath"`
	DocModifySeq     int             `json:"docModifySeq"`
	Language         string          `json:"language"`
	TextBeforeCursor string          `json:"textBeforeCursor"`
	TextAfterCursor  string          `json:"textAfterCursor"`
	CompletionLine   *CompletionLine `json:"completionLine"`
	TriggerType      string          `json:"triggerType"`
}

type CompletionLine struct {
	CurrentLineNum         int    `json:"currentLineNum"`
	TotalLineCount         int    `json:"totalLineCount"`
	CurrentLineStartOffset int    `json:"currentLineStartOffset"`
	CurrentCursorOffset    int    `json:"currentCursorOffset"`
	LineText               string `json:"lineText"`
	NextLineIndent         int    `json:"nextLineIndent"`
}
