package oai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jamesleeht/llm-gopher/params"

	"github.com/openai/openai-go/v2"
	"github.com/openai/openai-go/v2/option"
)

type ClientConfig struct {
	Name    string `json:"name"`
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url,omitempty"`
}

type Client struct {
	internalClient *openai.Client
}

func NewOpenAIClient(config ClientConfig) *Client {
	opts := []option.RequestOption{
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.BaseURL),
	}

	internalClient := openai.NewClient(opts...)
	return &Client{
		internalClient: &internalClient,
	}
}

// SendCompletionMessage is a generic function that sends a completion request and returns a typed response.
func SendCompletionMessage[T any](ctx context.Context, c *Client, prompt params.Prompt[T], settings params.Settings) (string, *T, error) {
	params := mapSettingsToParams(settings)
	messages := mapPromptToMessages(prompt)
	params.Messages = messages

	if rf := mapPromptToResponseFormat[T](); rf != nil {
		params.ResponseFormat = *rf
	}

	completion, err := c.internalClient.Chat.Completions.New(ctx, params)

	if err != nil {
		return "", nil, fmt.Errorf("failed to send completion message: %w", err)
	}

	content := completion.Choices[0].Message.Content

	// Check if T is not 'any' by attempting to unmarshal
	// If T is a concrete type, unmarshal the response
	var parsed T
	var parsedPtr *T
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		// Only set Parsed if T is not 'any' - we can check by seeing if unmarshaling worked
		// For 'any' type, parsed will be nil which is expected
		parsedPtr = &parsed
	} else {
		// If we expected a structured response but unmarshaling failed, return error
		if rf := mapPromptToResponseFormat[T](); rf != nil {
			return "", nil, fmt.Errorf("failed to unmarshal response into specified format: %w", err)
		}
	}

	return content, parsedPtr, nil
}
