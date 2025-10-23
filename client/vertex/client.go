package vertex

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"

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

func (c *Client) SendCompletionMessage(ctx context.Context, prompt params.Prompt, settings params.Settings) (*params.Response, error) {
	config, err := mapSettingsToVertexSettings(prompt, settings)
	if err != nil {
		return nil, fmt.Errorf("failed to map settings to vertex settings: %w", err)
	}

	messages := mapPromptToMessages(prompt)
	resp, err := c.internalClient.Models.GenerateContent(ctx,
		string(settings.ModelName),
		messages,
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
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
