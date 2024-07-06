package completion

import (
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
	"strings"
)

type PromptBuilder struct {
	request *Request
}

func NewPromptBuilder(request *Request) *PromptBuilder {
	return &PromptBuilder{
		request: request,
	}
}

func (b *PromptBuilder) Build() []*volcModel.ChatCompletionMessage {
	prompt := strings.Builder{}
	prompt.WriteString("现给你提供如下信息：")
	prompt.WriteString("当前正在编辑的代码文件：" + b.request.DocPath + "\n")

	prompt.WriteString("代码内容如下：\n")
	prompt.WriteString("```" + b.request.Language + "\n")
	prompt.WriteString(b.request.TextBeforeCursor)
	prompt.WriteString("[##CURSOR##]")
	prompt.WriteString(b.request.TextAfterCursor)
	prompt.WriteString("\n```\n")
	prompt.WriteString("其中[##CURSOR##]代表当前光标位置。\n")

	prompt.WriteString("你的任务是请阅读上述信息，补全光标处的代码内容。要求如下：\n")
	prompt.WriteString("1. 补全出来的代码内容与前文、后文代码拼接之后，能够正确编译，并且符合逻辑。\n")
	prompt.WriteString("2. 请只返回光标处要补全的代码，以markdown形式返回。\n")

	return []*volcModel.ChatCompletionMessage{
		{
			Role: volcModel.ChatMessageRoleSystem,
			Content: &volcModel.ChatCompletionMessageContent{
				StringValue: volcengine.String("你是一个智能代码补全助手"),
			},
		},
		{
			Role: volcModel.ChatMessageRoleUser,
			Content: &volcModel.ChatCompletionMessageContent{
				StringValue: volcengine.String(prompt.String()),
			},
		},
	}
}
