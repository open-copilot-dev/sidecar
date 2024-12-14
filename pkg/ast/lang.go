package ast

import (
	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"strings"
)

var supportLanguages = map[string]*sitter.Language{
	"golang": golang.GetLanguage(),
	"java":   java.GetLanguage(),
}

func GetLanguage(language string) *sitter.Language {
	lang, ok := supportLanguages[strings.ToLower(language)]
	if !ok {
		return nil
	}
	return lang
}
