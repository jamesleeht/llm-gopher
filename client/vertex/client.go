package vertex

import (
	"context"
	"encoding/json"
	"fmt"
	"llm-gopher/params"
	"os"

	"cloud.google.com/go/auth"
	"google.golang.org/genai"
)

// Client wraps the official Vertex AI SDK
type Client struct {
	internalClient *genai.Client
}

// ClientConfig holds the configuration for Vertex AI
type ClientConfig struct {
	ProjectID       string
	Location        string
	CredentialsPath string
}

// NewVertexAIClient creates a new Vertex AI client using the official SDK
func NewVertexAIClient(config ClientConfig) (*Client, error) {
	ctx := context.Background()

	creds, err := parseCredentialsFile(config.CredentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		Project:     config.ProjectID,
		Location:    config.Location,
		Credentials: creds,
		Backend:     genai.BackendVertexAI,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create genai client: %w", err)
	}

	return &Client{
		internalClient: client,
	}, nil
}

func (c *Client) SendCompletionMessage(ctx context.Context, prompt params.Prompt, settings params.Settings) (string, error) {
	config, err := mapSettingsToVertexSettings(prompt, settings)
	if err != nil {
		return "", fmt.Errorf("failed to map settings to vertex settings: %w", err)
	}

	messages := mapPromptToMessages(prompt)
	resp, err := c.internalClient.Models.GenerateContent(ctx,
		string(settings.ModelName),
		messages,
		config,
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate content: %w", err)
	}

	// candidate := resp.Candidates[0]
	// groundingMetadata := candidate.GroundingMetadata
	// if groundingMetadata != nil {
	// 	fmt.Println("Grounding metadata:")
	// 	fmt.Println(groundingMetadata)
	// 	searchEntryPoint := groundingMetadata.SearchEntryPoint
	// 	if searchEntryPoint != nil {
	// 		fmt.Println("Search entry point:")
	// 		fmt.Println(searchEntryPoint.RenderedContent)
	// 	}
	// }

	return resp.Text(), nil
}

func parseCredentialsFile(credentialsPath string) (*auth.Credentials, error) {
	// First read credentials file
	credJson, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	// Step 1: Parse the provided service account JSON key
	var sa struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
		TokenURI    string `json:"token_uri"`
		ProjectID   string `json:"project_id"`
	}
	if err := json.Unmarshal(credJson, &sa); err != nil {
		return nil, fmt.Errorf("invalid service-account JSON: %w", err)
	}

	// Step 2: Create a 2-legged OAuth token provider using the email and private key
	tp, err := auth.New2LOTokenProvider(&auth.Options2LO{
		Email:      sa.ClientEmail,
		PrivateKey: []byte(sa.PrivateKey),
		TokenURL:   sa.TokenURI,
		Scopes:     []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create 2LO token provider: %w", err)
	}

	// Step 3: Create credentials from the token provider and original JSON
	creds := auth.NewCredentials(&auth.CredentialsOptions{
		TokenProvider: tp,
		JSON:          []byte(credJson),
	})

	return creds, nil
}
