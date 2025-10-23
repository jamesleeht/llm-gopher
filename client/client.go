package client

import (
	"context"
	"fmt"

	"github.com/jamesleeht/llm-gopher/client/oai"
	"github.com/jamesleeht/llm-gopher/client/vertex"
	"github.com/jamesleeht/llm-gopher/params"
)

type Client struct {
	OpenAIClient   *oai.Client
	VertexAIClient *vertex.Client
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

// SendMessage is a generic function that sends a message and returns a typed response.
// Use T to specify the expected response type:
//   - SendMessage[any](...) for unstructured text responses
//   - SendMessage[YourStruct](...) for structured JSON responses
func SendMessage[T any](ctx context.Context,
	client *Client,
	prompt params.Prompt[T],
	settings params.Settings) (*Response[T], error) {

	var content string
	var parsed *T
	var err error

	switch client.ClientType {
	case ClientTypeOpenAI:
		content, parsed, err = oai.SendCompletionMessage[T](ctx, client.OpenAIClient, prompt, settings)
	case ClientTypeVertex:
		content, parsed, err = vertex.SendCompletionMessage[T](ctx, client.VertexAIClient, prompt, settings)
	default:
		return nil, fmt.Errorf("client type not supported")
	}

	if err != nil {
		return nil, err
	}

	return &Response[T]{
		Content: content,
		Parsed:  parsed,
	}, nil
}
