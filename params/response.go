package params

// Response contains both the raw string content and the unmarshalled struct
type Response struct {
	// Content is the raw string response from the LLM
	Content string
	// Parsed is a pointer to the unmarshalled struct if ResponseFormat was specified, nil otherwise.
	// Type assert as a pointer when using: myStruct := response.Parsed.(*MyStructType)
	Parsed interface{}
}

// StreamChunk represents a single chunk of a streaming response
type StreamChunk struct {
	// Content is the incremental text content in this chunk
	Content string
	// Done indicates if this is the final chunk in the stream
	Done bool
	// Error contains any error that occurred during streaming
	Error error
}
