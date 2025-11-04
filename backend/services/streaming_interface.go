package services

import (
	"context"
	"io"
)

// StreamingChatService provides a unified interface for streaming chat across different AI providers
type StreamingChatService interface {
	// CreateStreamingChat creates a streaming chat session
	// Returns a StreamReader that can be used to read chunks
	CreateStreamingChat(ctx context.Context, req StreamingChatRequest) (StreamReader, error)

	// GetProviderName returns the name of the provider (e.g., "openai", "bedrock")
	GetProviderName() string

	// IsAvailable checks if the service is configured and available
	IsAvailable() bool
}

// StreamReader provides a unified interface for reading streaming responses
type StreamReader interface {
	// Recv receives the next chunk of content from the stream
	// Returns io.EOF when the stream is complete
	// Returns error if there's a problem with the stream
	Recv() (string, error)

	// Close closes the stream and releases resources
	Close() error
}

// StreamingChatRequest represents a unified request for streaming chat
type StreamingChatRequest struct {
	// Messages can be either []openai.ChatCompletionMessage or []ClaudeMessage
	// The service implementation will handle the conversion
	Messages interface{}

	// SystemPrompt is the system prompt for the conversation
	SystemPrompt string

	// Temperature controls randomness (0.0 - 1.0)
	Temperature float64

	// MaxTokens is the maximum number of tokens to generate
	MaxTokens int

	// Model is the specific model ID to use (optional, uses default if empty)
	Model string
}

// StreamingChatResponse represents a chunk of streaming response
type StreamingChatResponse struct {
	// Content is the text content of this chunk
	Content string

	// Index is the sequential index of this chunk
	Index int

	// Done indicates if this is the final chunk
	Done bool

	// Error contains any error that occurred
	Error error
}

// ErrStreamClosed is returned when trying to read from a closed stream
var ErrStreamClosed = io.ErrClosedPipe
