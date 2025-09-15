package oai

import (
	"context"
	"fmt"
	"llm-gopher/params"

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

func (c *Client) SendCompletionMessage(ctx context.Context, prompt params.Prompt, settings params.Settings) (string, error) {
	params := mapSettingsToParams(settings)
	messages := mapPromptToMessages(prompt)
	params.Messages = messages

	if rf := mapPromptToResponseFormat(prompt); rf != nil {
		params.ResponseFormat = *rf
	}

	completion, err := c.internalClient.Chat.Completions.New(ctx, params)

	if err != nil {
		return "", fmt.Errorf("failed to send completion message: %w", err)
	}

	return completion.Choices[0].Message.Content, nil
}
