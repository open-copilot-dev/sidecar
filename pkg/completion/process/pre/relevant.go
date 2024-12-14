package pre

import "open-copilot.dev/sidecar/pkg/completion/domain"

/////////////////////////////////////////////////////////////////////////////////////////
// 前处理：提取补全相关的信息，例如：相邻文件

type RelevantPreProcessor struct {
}

func (f *RelevantPreProcessor) process(c *domain.CompletionContext) State {
	return StateContinue
}
