package main

import (
	"log"
	"os"

	"cloud.google.com/go/auth/credentials"
	"github.com/jamesleeht/llm-gopher/client"
)

// Example functions demonstrating different authentication methods for Vertex AI

// Example 1: Using a service account JSON file path (most common for development)
func createVertexClientWithFilePath() (*client.Client, error) {
	return client.NewClient(client.ClientConfig{
		ProjectID:             "my-gcp-project",
		Location:              "us-central1",
		VertexCredentialsPath: "/path/to/service-account.json",
	}, client.ClientTypeVertex)
}

// Example 2: Using Application Default Credentials (ADC) - recommended for production
// This automatically uses credentials from:
// - GOOGLE_APPLICATION_CREDENTIALS environment variable
// - gcloud auth application-default login
// - Compute Engine/GKE service account when running on GCP
func createVertexClientWithADC() (*client.Client, error) {
	return client.NewClient(client.ClientConfig{
		ProjectID: "my-gcp-project",
		Location:  "us-central1",
		// No credentials specified - will automatically use ADC
	}, client.ClientTypeVertex)
}

// Example 3: Using service account JSON content directly (from environment variable or secret manager)
func createVertexClientWithJSONContent() (*client.Client, error) {
	// Read credentials from environment variable or secret manager
	jsonContent := os.Getenv("VERTEX_AI_CREDENTIALS_JSON")

	return client.NewClient(client.ClientConfig{
		ProjectID:             "my-gcp-project",
		Location:              "us-central1",
		VertexCredentialsJSON: jsonContent,
	}, client.ClientTypeVertex)
}

// Example 4: Using pre-configured credentials object (for custom auth flows)
func createVertexClientWithPreConfiguredCreds() (*client.Client, error) {
	// Create credentials using Application Default Credentials
	creds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		return nil, err
	}

	return client.NewClient(client.ClientConfig{
		ProjectID:         "my-gcp-project",
		Location:          "us-central1",
		VertexCredentials: creds,
	}, client.ClientTypeVertex)
}

// Example 5: Custom credentials with service account impersonation
func createVertexClientWithImpersonation() (*client.Client, error) {
	// First, get credentials for the base service account
	baseCreds, err := credentials.DetectDefault(&credentials.DetectOptions{
		Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
	})
	if err != nil {
		return nil, err
	}

	// You could use impersonation or other custom auth flows here
	// This is just a demonstration - actual impersonation would require
	// additional configuration

	return client.NewClient(client.ClientConfig{
		ProjectID:         "my-gcp-project",
		Location:          "us-central1",
		VertexCredentials: baseCreds,
	}, client.ClientTypeVertex)
}

// Example 6: Reading credentials from a file but passing as JSON content
func createVertexClientWithFileAsJSON() (*client.Client, error) {
	// Read the file content
	jsonBytes, err := os.ReadFile("/path/to/service-account.json")
	if err != nil {
		log.Printf("failed to read credentials file: %v", err)
		return nil, err
	}

	return client.NewClient(client.ClientConfig{
		ProjectID:             "my-gcp-project",
		Location:              "us-central1",
		VertexCredentialsJSON: string(jsonBytes),
	}, client.ClientTypeVertex)
}

// Priority order demonstration:
// The authentication methods are tried in this order:
// 1. VertexCredentials (pre-configured *auth.Credentials)
// 2. VertexCredentialsJSON (JSON content as string)
// 3. VertexCredentialsPath (file path)
// 4. Application Default Credentials (automatic fallback)
func demonstratePriority() {
	// If multiple are provided, only the highest priority one is used
	_, err := client.NewClient(client.ClientConfig{
		ProjectID:             "my-gcp-project",
		Location:              "us-central1",
		VertexCredentials:     nil,                   // If this was set, it would be used
		VertexCredentialsJSON: "",                    // This would be tried second
		VertexCredentialsPath: "/path/to/creds.json", // This is tried third
		// If all above are empty/nil, ADC would be used (fourth)
	}, client.ClientTypeVertex)

	if err != nil {
		log.Printf("failed to create client: %v", err)
	}
}
