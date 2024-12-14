package pre

import (
	"github.com/cloudwego/hertz/pkg/common/hlog"
	sitter "github.com/smacker/go-tree-sitter"
	"open-copilot.dev/sidecar/pkg/ast"
	"open-copilot.dev/sidecar/pkg/completion/domain"
)

/////////////////////////////////////////////////////////////////////////////////////////
// 前处理：语法树解析、光标场景解析等

type GrammarPreProcessor struct {
}

func (f *GrammarPreProcessor) process(c *domain.CompletionContext) State {
	lang := ast.GetLanguage(c.Request.Language)
	if lang == nil {
		return StateContinue
	}

	// parse tree
	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	sourceCode := []byte(c.Request.TextBeforeCursor + c.Request.TextAfterCursor)
	tree, err := parser.ParseCtx(c.Ctx, nil, sourceCode)
	if tree == nil || err != nil {
		hlog.CtxErrorf(c.Ctx, "skip pre grammar process because tree is nil, err: %v", err)
		return StateContinue
	}

	// get which node current completion belong to
	cursorOffset := len(c.Request.TextBeforeCursor)
	node := ast.GetClosestNodeContainsRange(tree.RootNode(), uint32(cursorOffset), uint32(cursorOffset))

	// get cursor type
	cursorType := getCursorType(node, cursorOffset)

	c.Ast = &domain.CompletionAst{
		CursorType: cursorType,
		Node:       node,
		Tree:       tree,
	}

	return true
}

func getCursorType(node *sitter.Node, offset int) domain.CursorGrammarType {
	// TODO
	return domain.CursorGrammarTypeUnknown
}
