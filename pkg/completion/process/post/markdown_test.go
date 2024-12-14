package post

import (
	"github.com/stretchr/testify/assert"
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"testing"
)

func TestMarkdownPostProcessor_process(t *testing.T) {
	markdownProcessor := MarkdownPostProcessor{}
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

	modelText = "```java\nif (args.length > 0) {\n    // 这里添加具体的处理逻辑\n}\n```"
	text = markdownProcessor.process(c, modelText)
	assert.Equal(t, "if (args.length > 0) {\n    // 这里添加具体的处理逻辑\n}", text)
}
