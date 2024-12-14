package post

import (
	"open-copilot.dev/sidecar/pkg/completion/domain"
	"strings"
)

type OverlapProcessor struct {
}

func (m *OverlapProcessor) process(c *domain.CompletionContext, modelText string) string {
	// 将光标之前的单词提取出来
	lineTextBeforeCursor := c.Request.CompletionLine.GetLineTextBeforeCursor()
	words := strings.Fields(lineTextBeforeCursor)
	if len(words) == 0 {
		return modelText
	}

	// 将大模型第一行的单词提取出来
	lines := strings.Split(modelText, "\n")
	modelWords := strings.Fields(lines[0])

	// 重合单词对比
	i := len(words) - 1
	j := 0
	for {
		if i < 0 || j >= len(modelWords) {
			// 走到这边，证明 modelWords 与 words 完全匹配，则大模型返回的第一行可以直接去掉
			lines[0] = ""
			break
		}
		if words[i] == modelWords[j] {
			// 单词完全重合，继续向中间判断
			i--
			j++
			continue
		} else if strings.HasPrefix(modelWords[j], words[i]) {
			// 单词部分重合，去除重合的地方，然后返回
			modelWords = append([]string{
				strings.TrimPrefix(modelWords[j], words[i]),
			}, modelWords[j+1:]...)
			lines[0] = strings.Join(modelWords, " ")
			break
		} else {
			// 不重合，将 modelWords 中重合的部分去掉，然后返回
			if j == 0 {
				break
			}
			modelWords = modelWords[j:]
			lines[0] = strings.Join(modelWords, " ")
			break
		}
	}

	return strings.Join(lines, "\n")
}
