package services

import (
	"context"
	"fmt"
	"io"

	"chatbot/config"

	"github.com/sashabaranov/go-openai"
)

// OpenAIService handles OpenAI API interactions
type OpenAIService struct {
	client *openai.Client
	config *config.Config
}

// NewOpenAIService creates a new OpenAI service
func NewOpenAIService(cfg *config.Config) *OpenAIService {
	client := openai.NewClient(cfg.OpenAIAPIKey)
	return &OpenAIService{
		client: client,
		config: cfg,
	}
}

// ChatRequest represents a chat request with all options
type ChatRequest struct {
	Messages      []openai.ChatCompletionMessage
	Model         string
	Temperature   float32
	MaxTokens     int
	SystemPrompt  string
}

// ChatResponse represents the response from OpenAI
type ChatResponse struct {
	Content     string
	TokensUsed  int
	Model       string
	FinishReason string
}

// SendChatRequest sends a chat request to OpenAI API
func (s *OpenAIService) SendChatRequest(req ChatRequest) (*ChatResponse, error) {
	// Use default model if not specified
	if req.Model == "" {
		req.Model = s.config.OpenAIModel
	}

	// Use default temperature if not specified
	if req.Temperature == 0 {
		req.Temperature = float32(s.config.OpenAITemperature)
	}

	// Use default max tokens if not specified
	if req.MaxTokens == 0 {
		req.MaxTokens = s.config.OpenAIMaxTokens
	}

	// Prepare messages - add system prompt if provided
	messages := req.Messages
	if req.SystemPrompt != "" {
		messages = append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: req.SystemPrompt,
			},
		}, messages...)
	}

	// Call OpenAI API
	resp, err := s.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:       req.Model,
			Messages:    messages,
			Temperature: req.Temperature,
			MaxTokens:   req.MaxTokens,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to call OpenAI API: %w", err)
	}

	// Check if response has choices
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	// Build response
	return &ChatResponse{
		Content:      resp.Choices[0].Message.Content,
		TokensUsed:   resp.Usage.TotalTokens,
		Model:        resp.Model,
		FinishReason: string(resp.Choices[0].FinishReason),
	}, nil
}

// BuildContextMessages builds context array from recent messages
func (s *OpenAIService) BuildContextMessages(userMessage string, recentMessages []string, systemPrompt string) []openai.ChatCompletionMessage {
	messages := []openai.ChatCompletionMessage{}

	// Add system prompt if provided
	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// Add recent messages as context (simplified - in real app, would alternate user/assistant)
	for _, msg := range recentMessages {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: msg,
		})
	}

	// Add current user message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	return messages
}

// TranscriptionResponse represents the response from Whisper API
type TranscriptionResponse struct {
	Text     string
	Language string
	Duration float64
}

// TranscribeAudio transcribes audio file using OpenAI Whisper API
func (s *OpenAIService) TranscribeAudio(file io.Reader, filename string) (*TranscriptionResponse, error) {
	ctx := context.Background()

	// Create audio transcription request
	req := openai.AudioRequest{
		Model:    openai.Whisper1,
		FilePath: filename,
		Reader:   file,
	}

	// Call Whisper API
	resp, err := s.client.CreateTranscription(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to transcribe audio: %w", err)
	}

	// Build response (duration is mock for now as API doesn't return it)
	return &TranscriptionResponse{
		Text:     resp.Text,
		Language: resp.Language,
		Duration: 0, // Whisper API doesn't return duration
	}, nil
}

// CreateStreamingChatCompletion creates a streaming chat completion
func (s *OpenAIService) CreateStreamingChatCompletion(
	ctx context.Context,
	messages []openai.ChatCompletionMessage,
) (*openai.ChatCompletionStream, error) {
	// Use default model
	model := s.config.OpenAIModel
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	// Use default temperature
	temperature := float32(s.config.OpenAITemperature)
	if temperature == 0 {
		temperature = 0.7
	}

	// Use default max tokens
	maxTokens := s.config.OpenAIMaxTokens
	if maxTokens == 0 {
		maxTokens = 1000
	}

	// Create streaming request
	stream, err := s.client.CreateChatCompletionStream(
		ctx,
		openai.ChatCompletionRequest{
			Model:       model,
			Messages:    messages,
			Temperature: temperature,
			MaxTokens:   maxTokens,
			Stream:      true,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create streaming chat completion: %w", err)
	}

	return stream, nil
}