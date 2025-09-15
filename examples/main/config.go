package main

// In production, these values should be loaded from the environment variables
type appConfig struct {
	apiKey                  string
	novitaKey               string
	vertexAIProjectID       string
	vertexAILocation        string
	vertexAICredentialsPath string
}

func getAppConfig() appConfig {
	return appConfig{
		apiKey:                  "mock-api-key",
		novitaKey:               "mock-novita-key",
		vertexAIProjectID:       "mock-project-id",
		vertexAILocation:        "us-central1",
		vertexAICredentialsPath: "/path/to/mock/credentials.json",
	}
}
