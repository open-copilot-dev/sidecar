package domain

import (
	"open-copilot.dev/sidecar/pkg/common"
)

type CompletionContext struct {
	Ctx     *common.CancelableContext
	Request *CompletionRequest
}

func (c *CompletionContext) IsCanceled() bool {
	return c.Ctx.IsCanceled()
}

func (c *CompletionContext) GetLineTextBeforeCursor() string {
	lineCursorOffset := c.Request.CompletionLine.CurrentCursorOffset - c.Request.CompletionLine.CurrentLineStartOffset
	if lineCursorOffset < 0 || lineCursorOffset > len(c.Request.CompletionLine.LineText) {
		return ""
	}
	return c.Request.CompletionLine.LineText[:lineCursorOffset]
}

func (c *CompletionContext) GetLineTextAfterCursor() string {
	lineCursorOffset := c.Request.CompletionLine.CurrentCursorOffset - c.Request.CompletionLine.CurrentLineStartOffset
	if lineCursorOffset < 0 || lineCursorOffset > len(c.Request.CompletionLine.LineText) {
		return ""
	}
	return c.Request.CompletionLine.LineText[lineCursorOffset:]
}
