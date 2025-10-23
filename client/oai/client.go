package oai

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

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

func (c *Client) SendCompletionMessage(ctx context.Context, prompt params.Prompt, settings params.Settings) (interface{}, error) {
	params := mapSettingsToParams(settings)
	messages := mapPromptToMessages(prompt)
	params.Messages = messages

	if rf := mapPromptToResponseFormat(prompt); rf != nil {
		params.ResponseFormat = *rf
	}

	completion, err := c.internalClient.Chat.Completions.New(ctx, params)

	if err != nil {
		return nil, fmt.Errorf("failed to send completion message: %w", err)
	}

	content := completion.Choices[0].Message.Content

	// If response format is specified, unmarshal into that type
	if prompt.ResponseFormat != nil {
		// Create a new instance of the response format type
		responseType := reflect.TypeOf(prompt.ResponseFormat)
		responseValue := reflect.New(responseType).Interface()

		// Unmarshal the JSON content into the response format
		if err := json.Unmarshal([]byte(content), responseValue); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response into specified format: %w", err)
		}

		// Return the dereferenced value (not the pointer)
		return reflect.ValueOf(responseValue).Elem().Interface(), nil
	}

	// If no response format specified, return the raw string content
	return content, nil
}
