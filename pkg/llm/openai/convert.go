package openai

import (
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/ssestream"
	"io"
	"open-copilot.dev/sidecar/pkg/domain"
)

func convertMessages(messages []*domain.ChatCompletionMessage) []openai.ChatCompletionMessageParamUnion {
	openaiMessages := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))
	for _, message := range messages {
		openaiMessage := convertMessage(message)
		openaiMessages = append(openaiMessages, openaiMessage)
	}
	return openaiMessages
}

func convertMessage(message *domain.ChatCompletionMessage) openai.ChatCompletionMessageParamUnion {
	var openaiMessage openai.ChatCompletionMessageParamUnion = nil
	switch message.Role {
	case domain.ChatMessageRoleUser:
		openaiMessage = openai.UserMessage(*message.Content.StringValue)
	case domain.ChatMessageRoleAssistant:
		openaiMessage = openai.AssistantMessage(*message.Content.StringValue)
	case domain.ChatMessageRoleSystem:
		openaiMessage = openai.SystemMessage(*message.Content.StringValue)
	}
	return openaiMessage
}

func convertChoices(openaiChoices []openai.ChatCompletionChoice) []*domain.ChatCompletionChoice {
	choices := make([]*domain.ChatCompletionChoice, 0, len(openaiChoices))
	for _, choice := range openaiChoices {
		choices = append(choices, &domain.ChatCompletionChoice{
			Index:        int(choice.Index),
			Message:      convertChoiceMessage(choice.Message),
			FinishReason: string(choice.FinishReason),
		})
	}
	return choices
}

func convertChoiceMessage(choiceMessage openai.ChatCompletionMessage) *domain.ChatCompletionMessage {
	volcChoiceMessage := &domain.ChatCompletionMessage{
		Role: string(choiceMessage.Role),
		Content: &domain.ChatCompletionMessageContent{
			StringValue: &choiceMessage.Content,
		},
	}
	return volcChoiceMessage
}

type StreamReaderAdapter struct {
	stream *ssestream.Stream[openai.ChatCompletionChunk]
}

func (s *StreamReaderAdapter) Recv() (response *domain.ChatCompletionStreamResponse, err error) {
	if s.stream.Next() {
		return convertStreamResponse(s.stream.Current()), nil
	}
	if s.stream.Err() != nil {
		return nil, s.stream.Err()
	}
	return nil, io.EOF
}

func (s *StreamReaderAdapter) Close() error {
	return s.stream.Close()
}

func convertStreamReader(stream *ssestream.Stream[openai.ChatCompletionChunk]) (domain.ChatCompletionStreamReader, error) {
	if stream.Err() != nil {
		return nil, stream.Err()
	}
	return &StreamReaderAdapter{stream: stream}, nil
}

func convertStreamResponse(openaiStreamChunk openai.ChatCompletionChunk) *domain.ChatCompletionStreamResponse {
	return &domain.ChatCompletionStreamResponse{
		ID:      openaiStreamChunk.ID,
		Object:  string(openaiStreamChunk.Object),
		Created: openaiStreamChunk.Created,
		Model:   openaiStreamChunk.Model,
		Choices: convertStreamChoices(openaiStreamChunk.Choices),
		Usage:   convertUsage(openaiStreamChunk.Usage),
	}
}

func convertStreamChoices(openaiChoices []openai.ChatCompletionChunkChoice) []*domain.ChatCompletionStreamChoice {
	choices := make([]*domain.ChatCompletionStreamChoice, 0, len(openaiChoices))
	for _, choice := range openaiChoices {
		choices = append(choices, &domain.ChatCompletionStreamChoice{
			Index:        int(choice.Index),
			Delta:        &domain.ChatCompletionStreamChoiceDelta{Role: string(choice.Delta.Role), Content: choice.Delta.Content},
			FinishReason: string(choice.FinishReason),
		})
	}
	return choices
}

func convertUsage(openaiUsage openai.CompletionUsage) *domain.ChatPromptUsage {
	return &domain.ChatPromptUsage{
		PromptTokens:     int(openaiUsage.PromptTokens),
		CompletionTokens: int(openaiUsage.CompletionTokens),
		TotalTokens:      int(openaiUsage.TotalTokens),
	}
}
