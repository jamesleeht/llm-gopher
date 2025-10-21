package client

import (
	"context"
	"fmt"

	"cloud.google.com/go/auth"
	"github.com/jamesleeht/llm-gopher/client/oai"
	"github.com/jamesleeht/llm-gopher/client/vertex"
	"github.com/jamesleeht/llm-gopher/params"
)

type Client struct {
	OpenAIClient   ProviderClient
	VertexAIClient ProviderClient
	ClientType     ClientType
}

type ClientConfig struct {
	APIKey string

	// OpenAI only
	BaseURL string

	// Vertex AI - Authentication options (in priority order):
	// 1. VertexCredentials - Pre-configured credentials object
	// 2. VertexCredentialsJSON - Service account JSON content as string
	// 3. VertexCredentialsPath - Path to service account JSON file
	// 4. If none provided, falls back to Application Default Credentials (ADC)
	ProjectID             string
	Location              string
	VertexCredentials     interface{} // *auth.Credentials - using interface{} to avoid import
	VertexCredentialsJSON string
	VertexCredentialsPath string
}

type ProviderClient interface {
	SendCompletionMessage(ctx context.Context, prompt params.Prompt, settings params.Settings) (string, error)
}

func NewClient(config ClientConfig, clientType ClientType) (*Client, error) {
	var openAIClient *oai.Client
	var vertexAIClient *vertex.Client

	switch clientType {
	case ClientTypeOpenAI:
		openAIClient = oai.NewOpenAIClient(oai.ClientConfig{
			APIKey:  config.APIKey,
			BaseURL: config.BaseURL,
		})
	case ClientTypeVertex:
		var err error
		// Convert interface{} back to *auth.Credentials if provided
		var creds *auth.Credentials
		if config.VertexCredentials != nil {
			if c, ok := config.VertexCredentials.(*auth.Credentials); ok {
				creds = c
			}
		}
		if vertexAIClient, err = vertex.NewVertexAIClient(vertex.ClientConfig{
			ProjectID:       config.ProjectID,
			Location:        config.Location,
			Credentials:     creds,
			CredentialsJSON: config.VertexCredentialsJSON,
			CredentialsPath: config.VertexCredentialsPath,
		}); err != nil {
			return nil, err
		}
	}

	return &Client{
		OpenAIClient:   openAIClient,
		VertexAIClient: vertexAIClient,
		ClientType:     clientType,
	}, nil
}

func (c *Client) SendMessage(ctx context.Context,
	prompt params.Prompt,
	settings params.Settings) (string, error) {

	switch c.ClientType {
	case ClientTypeOpenAI:
		return c.OpenAIClient.SendCompletionMessage(ctx, prompt, settings)
	case ClientTypeVertex:
		return c.VertexAIClient.SendCompletionMessage(ctx, prompt, settings)
	}

	return "", fmt.Errorf("client type not supported")
}
