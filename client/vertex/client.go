package vertex

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jamesleeht/llm-gopher/params"

	"cloud.google.com/go/auth"
	"google.golang.org/genai"
)

type Client struct {
	internalClient *genai.Client
}

type ClientConfig struct {
	ProjectID string
	Location  string

	// path to a service account JSON file or JSON string. Used in development.
	CredentialsPath       string
	CredentialsJSONString string
}

// NewVertexAIClient creates a new Vertex AI client using the official SDK
func NewVertexAIClient(config ClientConfig) (*Client, error) {
	ctx := context.Background()

	// Determine which authentication method to use.
	// If no path or json string is provided, creds is set to nil.
	creds, err := getCredentials(config)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
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

// SendCompletionMessage is a generic function that sends a completion request and returns a typed response.
func SendCompletionMessage[T any](ctx context.Context, c *Client, prompt params.Prompt[T], settings params.Settings) (string, *T, error) {
	config, err := mapSettingsToVertexSettings[T](prompt, settings)
	if err != nil {
		return "", nil, fmt.Errorf("failed to map settings to vertex settings: %w", err)
	}

	messages := mapPromptToMessages(prompt)
	resp, err := c.internalClient.Models.GenerateContent(ctx,
		string(settings.ModelName),
		messages,
		config,
	)
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate content: %w", err)
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

	content := resp.Text()

	// Check if T is not 'any' by attempting to unmarshal
	// If T is a concrete type, unmarshal the response
	var parsed T
	var parsedPtr *T
	if err := json.Unmarshal([]byte(content), &parsed); err == nil {
		// Only set Parsed if T is not 'any'
		parsedPtr = &parsed
	} else {
		// If we expected a structured response but unmarshaling failed, return error
		if hasStructuredOutput[T]() {
			return "", nil, fmt.Errorf("failed to unmarshal response into specified format: %w", err)
		}
	}

	return content, parsedPtr, nil
}

func getCredentials(config ClientConfig) (*auth.Credentials, error) {
	if config.CredentialsPath != "" {
		return getCredsFromCredentialsPath(config.CredentialsPath)
	}
	if config.CredentialsJSONString != "" {
		return parseCredentialsJSON([]byte(config.CredentialsJSONString))
	}
	return nil, nil
}

func getCredsFromCredentialsPath(credentialsPath string) (*auth.Credentials, error) {
	credJSON, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	return parseCredentialsJSON(credJSON)
}

// parseCredentialsJSON parses service account JSON and creates credentials
func parseCredentialsJSON(credJSON []byte) (*auth.Credentials, error) {
	// Step 1: Parse the provided service account JSON key
	var sa struct {
		ClientEmail string `json:"client_email"`
		PrivateKey  string `json:"private_key"`
		TokenURI    string `json:"token_uri"`
		ProjectID   string `json:"project_id"`
	}
	if err := json.Unmarshal(credJSON, &sa); err != nil {
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
		JSON:          credJSON,
	})

	return creds, nil
}
