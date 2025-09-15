package client

import (
	"context"
	"llm-gopher/client/oai"
	"llm-gopher/client/vertex"
	"llm-gopher/params"
)

type Client struct {
	OpenAIClient   ProviderClient
	VertexAIClient ProviderClient
}

type ClientConfig struct {
	APIKey string

	// OpenAI only
	BaseURL string

	// Vertex AI
	ProjectID       string
	Location        string
	CredentialsPath string
}

type ProviderClient interface {
	SendCompletionMessage(ctx context.Context, prompt params.Prompt, settings params.Settings) (string, error)
}

func NewClient(config ClientConfig, isVertexClient bool) (*Client, error) {
	//nolint:exhaustruct
	openAIClient := oai.NewOpenAIClient(oai.ClientConfig{
		APIKey:  config.APIKey,
		BaseURL: config.BaseURL,
	})

	var vertexAIClient *vertex.Client
	if isVertexClient {
		vc, err := vertex.NewVertexAIClient(vertex.ClientConfig{
			ProjectID:       config.ProjectID,
			Location:        config.Location,
			CredentialsPath: config.CredentialsPath,
		})
		if err != nil {
			return nil, err
		}
		vertexAIClient = vc
	}

	return &Client{
		OpenAIClient:   openAIClient,
		VertexAIClient: vertexAIClient,
	}, nil
}

func (c *Client) SendMessage(ctx context.Context,
	prompt params.Prompt,
	settings params.Settings) (string, error) {

	if settings.ModelName.IsGemini() {
		return c.VertexAIClient.SendCompletionMessage(ctx, prompt, settings)
	}

	return c.OpenAIClient.SendCompletionMessage(ctx, prompt, settings)
}
