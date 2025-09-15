package params

import "llm-gopher/enums/messagerole"

type Message struct {
	Role    messagerole.MessageRole
	Content string
}

type Prompt struct {
	SystemMessage      string
	Messages           []Message
	ResponseFormatName string
	ResponseFormat     interface{} // JSON schema for response format
}

func NewPrompt(
	systemMessage string,
	messages []Message,
	responseFormatName string,
	responseFormat interface{}) Prompt {
	return Prompt{
		SystemMessage:      systemMessage,
		Messages:           messages,
		ResponseFormatName: responseFormatName,
		ResponseFormat:     responseFormat,
	}
}

func NewSimplePrompt(systemMessage string, userMessage string) Prompt {
	return Prompt{
		SystemMessage: systemMessage,
		Messages: []Message{
			{Role: messagerole.User, Content: userMessage},
		},
	}
}
