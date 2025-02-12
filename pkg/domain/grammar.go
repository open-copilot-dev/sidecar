package domain

// 光标处的语法类型

type CursorGrammarType string

var (
	CursorGrammarTypeLineComment  CursorGrammarType = "lineComment"
	CursorGrammarTypeBlockComment CursorGrammarType = "blockComment"
	CursorGrammarTypeString       CursorGrammarType = "string"
	CursorGrammarTypeUnknown      CursorGrammarType = "unknown"
)
