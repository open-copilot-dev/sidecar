package context

import (
	sitter "github.com/smacker/go-tree-sitter"
	"open-copilot.dev/sidecar/pkg/domain"
)

type CompletionContext struct {
	Ctx      *domain.CancelableContext
	Request  *domain.CompletionRequest
	Ast      *CompletionAst
	Relevant *CompletionRelevant
}

// IsCanceled returns true if the context has been canceled.
func (c *CompletionContext) IsCanceled() bool {
	return c.Ctx.IsCanceled()
}

// GetLineTextBeforeCursor returns the text before the cursor in the current line.
func (c *CompletionContext) GetLineTextBeforeCursor() string {
	lineCursorOffset := c.Request.CompletionLine.CurrentCursorOffset - c.Request.CompletionLine.CurrentLineStartOffset
	if lineCursorOffset < 0 || lineCursorOffset > len(c.Request.CompletionLine.LineText) {
		return ""
	}
	return c.Request.CompletionLine.LineText[:lineCursorOffset]
}

// GetLineTextAfterCursor returns the text after the cursor in the current line.
func (c *CompletionContext) GetLineTextAfterCursor() string {
	lineCursorOffset := c.Request.CompletionLine.CurrentCursorOffset - c.Request.CompletionLine.CurrentLineStartOffset
	if lineCursorOffset < 0 || lineCursorOffset > len(c.Request.CompletionLine.LineText) {
		return ""
	}
	return c.Request.CompletionLine.LineText[lineCursorOffset:]
}

type CompletionAst struct {
	// current completion tree
	Tree *sitter.Tree
	// current completion node
	Node *sitter.Node
	// current completion cursor grammar type
	CursorType domain.CursorGrammarType
}

type CompletionRelevant struct {
}
