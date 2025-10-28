package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jamesleeht/llm-gopher/client"
	"github.com/jamesleeht/llm-gopher/params"
)

func main() {
	// Example using OpenAI
	streamWithOpenAI()

	// Example using Vertex AI
	// streamWithVertex()
}

func streamWithOpenAI() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY environment variable is not set")
		return
	}

	// Create OpenAI client
	llmClient, err := client.NewClient(client.ClientConfig{
		APIKey: apiKey,
	}, client.ClientTypeOpenAI)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	// Create a simple prompt
	prompt := params.NewSimplePrompt(
		"You are a helpful assistant.",
		"Tell me a short story about a robot learning to paint.",
	)

	// Settings
	temperature := 0.7
	settings := params.Settings{
		ModelName:   "gpt-4o-mini",
		Temperature: &temperature,
	}

	// Stream the response
	ctx := context.Background()
	chunks, err := llmClient.StreamMessage(ctx, prompt, settings)
	if err != nil {
		fmt.Printf("Failed to start streaming: %v\n", err)
		return
	}

	fmt.Println("Streaming response:")
	fmt.Println("---")

	// Process chunks as they arrive
	for chunk := range chunks {
		if chunk.Error != nil {
			fmt.Printf("\nError: %v\n", chunk.Error)
			break
		}

		if chunk.Done {
			fmt.Println("\n---")
			fmt.Println("Stream complete!")
			break
		}

		// Print the content without newline to show streaming effect
		fmt.Print(chunk.Content)
	}
}

func streamWithVertex() {
	projectID := os.Getenv("VERTEX_PROJECT_ID")
	location := os.Getenv("VERTEX_LOCATION")
	credsPath := os.Getenv("VERTEX_CREDENTIALS_PATH")

	if projectID == "" || location == "" {
		fmt.Println("VERTEX_PROJECT_ID and VERTEX_LOCATION environment variables must be set")
		return
	}

	// Create Vertex AI client
	llmClient, err := client.NewClient(client.ClientConfig{
		ProjectID:             projectID,
		Location:              location,
		VertexCredentialsPath: credsPath,
	}, client.ClientTypeVertex)
	if err != nil {
		fmt.Printf("Failed to create client: %v\n", err)
		return
	}

	// Create a simple prompt
	prompt := params.NewSimplePrompt(
		"You are a helpful assistant.",
		"Tell me a short story about a robot learning to paint.",
	)

	// Settings
	temperature := 0.7
	settings := params.Settings{
		ModelName:   "gemini-2.0-flash-exp",
		Temperature: &temperature,
	}

	// Stream the response
	ctx := context.Background()
	chunks, err := llmClient.StreamMessage(ctx, prompt, settings)
	if err != nil {
		fmt.Printf("Failed to start streaming: %v\n", err)
		return
	}

	fmt.Println("Streaming response:")
	fmt.Println("---")

	// Process chunks as they arrive
	for chunk := range chunks {
		if chunk.Error != nil {
			fmt.Printf("\nError: %v\n", chunk.Error)
			break
		}

		if chunk.Done {
			fmt.Println("\n---")
			fmt.Println("Stream complete!")
			break
		}

		// Print the content without newline to show streaming effect
		fmt.Print(chunk.Content)
	}
}
