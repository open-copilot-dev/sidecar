package volcengine

import (
	volcModel "github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	volcUtils "github.com/volcengine/volcengine-go-sdk/service/arkruntime/utils"
	"open-copilot.dev/sidecar/pkg/domain"
)

func convertRequest(request *domain.ChatCompletionRequest, model string) volcModel.ChatCompletionRequest {
	return volcModel.ChatCompletionRequest{
		Model:       model,
		Messages:    convertMessages(request.Messages),
		MaxTokens:   request.MaxTokens,
		Temperature: request.Temperature,
		TopP:        request.TopP,
		Stop:        request.Stop,
		N:           request.N,
	}
}

func convertMessages(messages []*domain.ChatCompletionMessage) []*volcModel.ChatCompletionMessage {
	volcMessages := make([]*volcModel.ChatCompletionMessage, 0, len(messages))
	for _, message := range messages {
		volcMessage := volcModel.ChatCompletionMessage{
			Role: message.Role,
			Content: &volcModel.ChatCompletionMessageContent{
				StringValue: message.Content.StringValue,
			},
		}
		volcMessages = append(volcMessages, &volcMessage)
	}
	return volcMessages
}

func convertResponse(volcCompletion volcModel.ChatCompletionResponse) *domain.ChatCompletionResponse {
	return &domain.ChatCompletionResponse{
		ID:      volcCompletion.ID,
		Object:  volcCompletion.Object,
		Created: volcCompletion.Created,
		Model:   volcCompletion.Model,
		Choices: convertChoices(volcCompletion.Choices),
		Usage:   convertUsage(&volcCompletion.Usage),
	}
}

func convertChoices(volcChoices []*volcModel.ChatCompletionChoice) []*domain.ChatCompletionChoice {
	choices := make([]*domain.ChatCompletionChoice, 0, len(volcChoices))
	for _, choice := range volcChoices {
		choices = append(choices, &domain.ChatCompletionChoice{
			Index:        choice.Index,
			Message:      convertChoiceMessage(choice.Message),
			FinishReason: string(choice.FinishReason),
		})
	}
	return choices
}

func convertChoiceMessage(volcMessage volcModel.ChatCompletionMessage) *domain.ChatCompletionMessage {
	return &domain.ChatCompletionMessage{
		Role:    volcMessage.Role,
		Content: &domain.ChatCompletionMessageContent{StringValue: volcMessage.Content.StringValue},
	}
}
func convertUsage(volcUsage *volcModel.Usage) *domain.ChatPromptUsage {
	return &domain.ChatPromptUsage{
		PromptTokens:     volcUsage.PromptTokens,
		CompletionTokens: volcUsage.CompletionTokens,
		TotalTokens:      volcUsage.TotalTokens,
	}
}

type StreamReaderAdapter struct {
	volcStream *volcUtils.ChatCompletionStreamReader
}

func (v *StreamReaderAdapter) Recv() (response *domain.ChatCompletionStreamResponse, err error) {
	volcStreamResp, err := v.volcStream.Recv()
	if err != nil {
		return response, err
	}
	return convertStreamResponse(volcStreamResp), err
}

func (v *StreamReaderAdapter) Close() error {
	return v.volcStream.Close()
}

func convertStreamReader(volcStream *volcUtils.ChatCompletionStreamReader) domain.ChatCompletionStreamReader {
	return &StreamReaderAdapter{
		volcStream: volcStream,
	}
}

func convertStreamResponse(volcStreamResp volcModel.ChatCompletionStreamResponse) *domain.ChatCompletionStreamResponse {
	return &domain.ChatCompletionStreamResponse{
		ID:      volcStreamResp.ID,
		Object:  volcStreamResp.Object,
		Created: volcStreamResp.Created,
		Model:   volcStreamResp.Model,
		Choices: convertStreamChoices(volcStreamResp.Choices),
		Usage:   convertUsage(volcStreamResp.Usage),
	}
}

func convertStreamChoices(volcStreamChoices []*volcModel.ChatCompletionStreamChoice) []*domain.ChatCompletionStreamChoice {
	choices := make([]*domain.ChatCompletionStreamChoice, 0, len(volcStreamChoices))
	for _, choice := range volcStreamChoices {
		choices = append(choices, &domain.ChatCompletionStreamChoice{
			Index:        choice.Index,
			Delta:        &domain.ChatCompletionStreamChoiceDelta{Role: choice.Delta.Role, Content: choice.Delta.Content},
			FinishReason: string(choice.FinishReason),
		})
	}
	return choices
}
