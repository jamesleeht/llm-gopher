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

func (c *Client) SendCompletionMessage(ctx context.Context, prompt params.Prompt, settings params.Settings) (*params.Response, error) {
	chatParams := mapSettingsToParams(settings)
	messages := mapPromptToMessages(prompt)
	chatParams.Messages = messages

	if rf := mapPromptToResponseFormat(prompt); rf != nil {
		chatParams.ResponseFormat = *rf
	}

	completion, err := c.internalClient.Chat.Completions.New(ctx, chatParams)

	if err != nil {
		return nil, fmt.Errorf("failed to send completion message: %w", err)
	}

	content := completion.Choices[0].Message.Content

	response := &params.Response{
		Content: content,
		Parsed:  nil,
	}

	// If response format is specified, unmarshal into that type
	if prompt.ResponseFormat != nil {
		// ResponseFormat must be a pointer to unmarshal into
		responseType := reflect.TypeOf(prompt.ResponseFormat)

		if responseType.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("response format must be a pointer, got %v", responseType.Kind())
		}

		// Unmarshal directly into the pointer provided by the user
		if err := json.Unmarshal([]byte(content), prompt.ResponseFormat); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response into specified format: %w", err)
		}

		// Set the parsed field to the populated pointer
		response.Parsed = prompt.ResponseFormat
	}

	return response, nil
}
