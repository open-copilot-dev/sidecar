package completion

type CompletionRequest struct {
	UUID           string         `json:"uuid"`
	DocModifySeq   int            `json:"docModifySeq"`
	Language       string         `json:"language"`
	CompletionLine CompletionLine `json:"completionLine"`
	TriggerType    string         `json:"triggerType"`
}

type CompletionLine struct {
	CurrentLineNum         int    `json:"currentLineNum"`
	TotalLineCount         int    `json:"totalLineCount"`
	CurrentLineStartOffset int    `json:"currentLineStartOffset"`
	CurrentCursorOffset    int    `json:"currentCursorOffset"`
	LineText               string `json:"lineText"`
	NextLineIndent         int    `json:"nextLineIndent"`
}

type CompletionResult struct {
}
