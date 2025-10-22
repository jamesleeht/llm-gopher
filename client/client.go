package client

import (
	"context"
	"fmt"

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

	// Vertex only
	ProjectID             string
	Location              string
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
		if vertexAIClient, err = vertex.NewVertexAIClient(vertex.ClientConfig{
			ProjectID:             config.ProjectID,
			Location:              config.Location,
			CredentialsPath:       config.VertexCredentialsPath,
			CredentialsJSONString: config.VertexCredentialsJSON,
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
