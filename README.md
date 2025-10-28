# LLM Gopher

This is a small library that acts as a router and adapter for different LLM endpoints.

It doesn't provide any big framework capabilities - it is simply meant as an easy way to create a router for different models.

It is partially inspired by the [LiteLLM Python SDK](https://github.com/BerriAI/litellm) but in a Golang context.

## Client types supported

- OpenAI
- Vertex AI

### Vertex AI Authentication

The Vertex AI client supports multiple authentication methods, tried in the following priority order:

1. **Pre-configured Credentials** - Pass an `*auth.Credentials` object directly
2. **Service Account JSON (as string)** - Pass JSON content as a string (useful for secrets managers)
3. **Service Account JSON File** - Path to a service account JSON file
4. **Application Default Credentials (ADC)** - Automatic fallback (recommended for production)

#### Examples

**Using a service account file (development):**

```go
client, err := client.NewClient(client.ClientConfig{
    ProjectID:             "my-project",
    Location:              "us-central1",
    VertexCredentialsPath: "/path/to/service-account.json",
}, client.ClientTypeVertex)
```

**Using Application Default Credentials (production):**

```go
// No credentials specified - automatically uses:
// - GOOGLE_APPLICATION_CREDENTIALS env var
// - gcloud auth application-default login
// - GCE/GKE service account
client, err := client.NewClient(client.ClientConfig{
    ProjectID: "my-project",
    Location:  "us-central1",
}, client.ClientTypeVertex)
```

**Using JSON content from environment variable:**

```go
client, err := client.NewClient(client.ClientConfig{
    ProjectID:             "my-project",
    Location:              "us-central1",
    VertexCredentialsJSON: os.Getenv("VERTEX_CREDS_JSON"),
}, client.ClientTypeVertex)
```

**Using pre-configured credentials:**

```go
creds, _ := credentials.DetectDefault(&credentials.DetectOptions{
    Scopes: []string{"https://www.googleapis.com/auth/cloud-platform"},
})
client, err := client.NewClient(client.ClientConfig{
    ProjectID:         "my-project",
    Location:          "us-central1",
    VertexCredentials: creds,
}, client.ClientTypeVertex)
```

See `examples/basic/auth_examples.go` for more detailed examples.

## Streaming Responses

The library supports streaming responses from both OpenAI and Vertex AI providers. Streaming allows you to receive the response incrementally as it's being generated, rather than waiting for the complete response.

### Basic Usage

```go
ctx := context.Background()
chunks, err := client.StreamMessage(ctx, prompt, settings)
if err != nil {
    log.Fatal(err)
}

// Process chunks as they arrive
for chunk := range chunks {
    if chunk.Error != nil {
        log.Printf("Error: %v\n", chunk.Error)
        break
    }

    if chunk.Done {
        fmt.Println("Stream complete!")
        break
    }

    // Print content without newline to show streaming effect
    fmt.Print(chunk.Content)
}
```

### Stream Chunk Structure

Each `StreamChunk` contains:

- `Content` (string): The incremental text content in this chunk
- `Done` (bool): Indicates if this is the final chunk in the stream
- `Error` (error): Contains any error that occurred during streaming

### Example

See `examples/streaming/main.go` for a complete working example with both OpenAI and Vertex AI.

### Limitations

- **OpenAI**: Response format (structured outputs) is not supported in streaming mode
- **Vertex AI**: Response format may have limitations in streaming mode

## Presets

A preset represents a combination of the model and its settings.
We use this internally since we have different use cases which might require a specific combination.

## Router

The router can be configured with a bunch of clients and presets.

1. When you send a prompt to the router, you specify a preset.
2. The preset's settings will be applied and an appropriate client will be selected for the model.
