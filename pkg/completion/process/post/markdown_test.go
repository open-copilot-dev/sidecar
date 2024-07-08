package post

import (
	"github.com/stretchr/testify/assert"
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"testing"
)

func TestMarkdownProcessor_process(t *testing.T) {
	markdownProcessor := MarkdownProcessor{}
	c := &domain.CompletionContext{}
	var modelText string
	var text string

	modelText = "`int a=1;`"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "`int a=1;` aaa `int b=1;`"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "```\nint a=1;\n```"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "```int a=1;```"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "```java\nint a=1;\n```"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "aaa```\nint a=1;\n```"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "aaa```int a=1;```bbb"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "aaa```java\nint a=1;\n```bbb"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "```\nint a=1;\n```\naaa\n```\nint a=2;\n```"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "```int a=1;```\naaa\n```\nint a=2;\n```"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)

	modelText = "```java\nint a=1;\n```\naaa\n```\nint a=2;\n```"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "int a=1;", text)
}
