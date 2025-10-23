package params

type Message struct {
	Role    MessageRole
	Content string
}

// Prompt is a generic type where T is the expected response format type.
// For unstructured text responses, use Prompt[any].
// For structured JSON responses, use Prompt[YourStructType] where YourStructType is your desired output struct.
type Prompt[T any] struct {
	SystemMessage string
	Messages      []Message
}

func NewPrompt[T any](
	systemMessage string,
	messages []Message) Prompt[T] {
	return Prompt[T]{
		SystemMessage: systemMessage,
		Messages:      messages,
	}
}

func NewSimplePrompt(systemMessage string, userMessage string) Prompt[any] {
	return Prompt[any]{
		SystemMessage: systemMessage,
		Messages: []Message{
			{Role: MessageRoleUser, Content: userMessage},
		},
	}
}
