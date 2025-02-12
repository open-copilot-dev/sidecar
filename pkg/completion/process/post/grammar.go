package post

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	sitter "github.com/smacker/go-tree-sitter"
	"open-copilot.dev/sidecar/pkg/ast"
	"open-copilot.dev/sidecar/pkg/completion/context"
)

/////////////////////////////////////////////////////////////////////////////////////////
// 后处理：语法合法性检验

type GrammarPostProcessor struct {
}

func (m *GrammarPostProcessor) process(c *context.CompletionContext, modelText string) string {
	lang := ast.GetLanguage(c.Request.Language)
	if lang == nil {
		return modelText
	}
	parser := sitter.NewParser()
	parser.SetLanguage(lang)

	sourceCode := []byte(c.Request.TextBeforeCursor + modelText + c.Request.TextAfterCursor)
	tree, err := parser.ParseCtx(c.Ctx, nil, sourceCode)
	if tree == nil || err != nil {
		hlog.CtxErrorf(c.Ctx, "skip grammar because tree is nil, err: %v", err)
		return modelText
	}

	startPos := len(c.Request.TextBeforeCursor)
	endPos := startPos + len(modelText)

	node := ast.GetClosestNodeContainsRange(tree.RootNode(), uint32(startPos), uint32(endPos))
	if node == nil {
		hlog.CtxErrorf(c.Ctx, "grammar check failed, node is nil")
		return ""
	}
	if isNodeHasGrammarError(node, uint32(startPos), uint32(endPos)) {
		hlog.CtxErrorf(c.Ctx, "grammar check not pass")
		return ""
	}

	return modelText
}

func isNodeHasGrammarError(node *sitter.Node, startPos uint32, endPos uint32) bool {
	if node.IsError() || node.HasError() {
		return true
	}
	if node.StartByte() == startPos && node.EndByte() == endPos {
		parent := node.Parent()
		if parent != nil {
			return parent.IsError()
		}
	}
	return false
}
