package client

type ClientType string

const (
	ClientTypeOpenAI ClientType = "openai"
	ClientTypeVertex ClientType = "vertex"
)

// Response contains the result of a completion request.
// T is the type of the parsed response. Use any if no structured parsing is needed.
type Response[T any] struct {
	Content string // Raw text content from the model
	Parsed  *T     // Unmarshalled struct if ResponseFormat was specified, nil otherwise
}
