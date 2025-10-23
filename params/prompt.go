package params

type Message struct {
	Role    MessageRole
	Content string
}

type Prompt struct {
	SystemMessage  string
	Messages       []Message
	ResponseFormat interface{} // Pointer to struct for JSON schema response format. It cannot be a nil pointer, or it will be ignored.
}

func NewPrompt(
	systemMessage string,
	messages []Message,
	responseFormat interface{}) Prompt {
	return Prompt{
		SystemMessage:  systemMessage,
		Messages:       messages,
		ResponseFormat: responseFormat,
	}
}

func NewSimplePrompt(systemMessage string, userMessage string) Prompt {
	return Prompt{
		SystemMessage: systemMessage,
		Messages: []Message{
			{Role: MessageRoleUser, Content: userMessage},
		},
	}
}
