package vertex

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jamesleeht/llm-gopher/params"

	"cloud.google.com/go/auth"
	"cloud.google.com/go/auth/credentials"
	"google.golang.org/genai"
)

// Client wraps the official Vertex AI SDK
type Client struct {
	internalClient *genai.Client
}

// ClientConfig holds the configuration for Vertex AI
//
// Authentication Methods (in priority order):
//
// 1. Pre-configured Credentials:
//
//	creds := ... // obtain credentials from somewhere
//	config := ClientConfig{
//		ProjectID:   "my-project",
//		Location:    "us-central1",
//		Credentials: creds,
//	}
//
// 2. Service Account JSON (as string):
//
//	jsonContent := `{"type": "service_account", ...}`
//	config := ClientConfig{
//		ProjectID:       "my-project",
//		Location:        "us-central1",
//		CredentialsJSON: jsonContent,
//	}
//
// 3. Service Account JSON File (path):
//
//	config := ClientConfig{
//		ProjectID:       "my-project",
//		Location:        "us-central1",
//		CredentialsPath: "/path/to/service-account.json",
//	}
//
// 4. Application Default Credentials (ADC) - automatic fallback:
//
//	config := ClientConfig{
//		ProjectID: "my-project",
//		Location:  "us-central1",
//		// No auth fields - will use ADC automatically
//		// ADC sources (in order):
//		// - GOOGLE_APPLICATION_CREDENTIALS env var
//		// - gcloud auth application-default credentials
//		// - Compute Engine/GKE service account
//	}
type ClientConfig struct {
	ProjectID string
	Location  string

	// Credentials is a pre-configured auth.Credentials object (priority 1)
	Credentials *auth.Credentials
	// CredentialsJSON is the service account JSON content as a string (priority 2)
	CredentialsJSON string
	// CredentialsPath is the path to a service account JSON file (priority 3)
	CredentialsPath string
	// If none of the above are provided, Application Default Credentials (ADC) will be used (priority 4)
}

// NewVertexAIClient creates a new Vertex AI client using the official SDK
func NewVertexAIClient(config ClientConfig) (*Client, error) {
	ctx := context.Background()

	// Determine which authentication method to use
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

// getCredentials resolves credentials from the config in priority order
func getCredentials(config ClientConfig) (*auth.Credentials, error) {
	// Priority 1: Pre-configured credentials object
	if config.Credentials != nil {
		return config.Credentials, nil
	}

	// Priority 2: Direct JSON content
	if config.CredentialsJSON != "" {
		return parseCredentialsJSON([]byte(config.CredentialsJSON))
	}

	// Priority 3: Credentials file path
	if config.CredentialsPath != "" {
		return parseCredentialsFile(config.CredentialsPath)
	}

	// Priority 4: Application Default Credentials (ADC)
	// This will use GOOGLE_APPLICATION_CREDENTIALS env var or gcloud auth
	return getApplicationDefaultCredentials()
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

// parseCredentialsFile reads and parses a service account JSON file
func parseCredentialsFile(credentialsPath string) (*auth.Credentials, error) {
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

// getApplicationDefaultCredentials uses Application Default Credentials (ADC)
// This will automatically use:
// 1. GOOGLE_APPLICATION_CREDENTIALS environment variable
// 2. gcloud auth application-default credentials
// 3. Compute Engine/GKE service account when running on GCP
func getApplicationDefaultCredentials() (*auth.Credentials, error) {
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to detect default credentials: %w", err)
	}
	return creds, nil
}
