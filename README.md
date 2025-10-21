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

## Presets

A preset represents a combination of the model and its settings.
We use this internally since we have different use cases which might require a specific combination.

## Router

The router can be configured with a bunch of clients and presets.

1. When you send a prompt to the router, you specify a preset.
2. The preset's settings will be applied and an appropriate client will be selected for the model.
