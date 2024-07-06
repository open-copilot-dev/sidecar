package completion

// Request 代码补全请求
type Request struct {
	UUID             string `json:"uuid"`
	ProjectPath      string `json:"projectPath"`
	DocPath          string `json:"docPath"`
	DocModifySeq     int    `json:"docModifySeq"`
	Language         string `json:"language"`
	TextBeforeCursor string `json:"textBeforeCursor"`
	TextAfterCursor  string `json:"textAfterCursor"`
	CompletionLine   Line   `json:"completionLine"`
	TriggerType      string `json:"triggerType"`
}

type Line struct {
	CurrentLineNum         int    `json:"currentLineNum"`
	TotalLineCount         int    `json:"totalLineCount"`
	CurrentLineStartOffset int    `json:"currentLineStartOffset"`
	CurrentCursorOffset    int    `json:"currentCursorOffset"`
	LineText               string `json:"lineText"`
	NextLineIndent         int    `json:"nextLineIndent"`
}

// Result 代码补全结果
type Result struct {
	Choices []*Choice `json:"choices"`
}

type Choice struct {
	Edits []*Edit `json:"edits"`
}

type Edit struct {
	StartOffset int    `json:"startOffset"`
	EndOffset   int    `json:"endOffset"`
	Text        string `json:"text"`
}
