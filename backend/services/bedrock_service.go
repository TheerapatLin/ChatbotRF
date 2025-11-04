package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"chatbot/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

// BedrockService handles AWS Bedrock API interactions
type BedrockService struct {
	client *bedrockruntime.Client
	config *config.Config
}

// NewBedrockService creates a new Bedrock service
func NewBedrockService(cfg *config.Config) (*BedrockService, error) {
	// Load AWS config with explicit credentials
	awsCfg, err := awsconfig.LoadDefaultConfig(
		context.Background(),
		awsconfig.WithRegion(cfg.AWSRegion),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AWSAccessKeyID,
				cfg.AWSSecretAccessKey,
				"", // session token (empty for static credentials)
			),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create Bedrock Runtime client
	client := bedrockruntime.NewFromConfig(awsCfg)

	log.Printf("✓ AWS Bedrock client initialized (region: %s, model: %s)",
		cfg.AWSRegion, cfg.BedrockModelID)

	return &BedrockService{
		client: client,
		config: cfg,
	}, nil
}

// ========================================
// Claude Messages API Format (Bedrock)
// ========================================

// ClaudeMessage represents a single message in the conversation
type ClaudeMessage struct {
	Role    string `json:"role"` // "user" or "assistant"
	Content string `json:"content"`
}

// ClaudeRequest represents the request body for Claude on Bedrock
type ClaudeRequest struct {
	AnthropicVersion string          `json:"anthropic_version"`
	MaxTokens        int             `json:"max_tokens"`
	Messages         []ClaudeMessage `json:"messages"`
	Temperature      float64         `json:"temperature,omitempty"`
	SystemPrompt     string          `json:"system,omitempty"`
}

// ClaudeResponse represents the response from Claude on Bedrock
type ClaudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model      string `json:"model"`
	StopReason string `json:"stop_reason"`
	Usage      struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// BedrockChatRequest represents a chat request with all options
type BedrockChatRequest struct {
	Messages     []ClaudeMessage
	SystemPrompt string
	Temperature  float64
	MaxTokens    int
}

// BedrockChatResponse represents the response from Bedrock
type BedrockChatResponse struct {
	Content    string
	TokensUsed int
	Model      string
	StopReason string
}

// SendChatRequest sends a chat request to AWS Bedrock (Claude)
func (s *BedrockService) SendChatRequest(req BedrockChatRequest) (*BedrockChatResponse, error) {
	// Use default values if not specified
	if req.Temperature == 0 {
		req.Temperature = s.config.BedrockTemperature
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = s.config.BedrockMaxTokens
	}

	// Build Claude request
	claudeReq := ClaudeRequest{
		AnthropicVersion: "bedrock-2023-05-31",
		MaxTokens:        req.MaxTokens,
		Messages:         req.Messages,
		Temperature:      req.Temperature,
	}

	// Add system prompt if provided
	if req.SystemPrompt != "" {
		claudeReq.SystemPrompt = req.SystemPrompt
	}

	// Marshal request to JSON
	requestBody, err := json.Marshal(claudeReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("🔵 Bedrock Request: model=%s, messages=%d, max_tokens=%d",
		s.config.BedrockModelID, len(req.Messages), req.MaxTokens)

	// Call Bedrock API
	output, err := s.client.InvokeModel(context.Background(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(s.config.BedrockModelID),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        requestBody,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to invoke Bedrock model: %w", err)
	}

	// Parse response
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(output.Body, &claudeResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Extract text from content
	var responseText string
	if len(claudeResp.Content) > 0 {
		responseText = claudeResp.Content[0].Text
	}

	totalTokens := claudeResp.Usage.InputTokens + claudeResp.Usage.OutputTokens

	log.Printf("✓ Bedrock Response: tokens=%d (in=%d, out=%d), stop_reason=%s",
		totalTokens, claudeResp.Usage.InputTokens, claudeResp.Usage.OutputTokens,
		claudeResp.StopReason)

	return &BedrockChatResponse{
		Content:    responseText,
		TokensUsed: totalTokens,
		Model:      s.config.BedrockModelID,
		StopReason: claudeResp.StopReason,
	}, nil
}

// BuildMessages builds message array from user input (simple version)
func (s *BedrockService) BuildMessages(userMessage string, systemPrompt string) []ClaudeMessage {
	return []ClaudeMessage{
		{
			Role:    "user",
			Content: userMessage,
		},
	}
}

// BuildMessagesWithHistory builds message array with conversation history
func (s *BedrockService) BuildMessagesWithHistory(
	userMessage string,
	history []ClaudeMessage,
) []ClaudeMessage {
	messages := make([]ClaudeMessage, 0, len(history)+1)

	// Add conversation history
	messages = append(messages, history...)

	// Add current user message
	messages = append(messages, ClaudeMessage{
		Role:    "user",
		Content: userMessage,
	})

	return messages
}

// BuildMessagesWithContext builds message array with file context
func (s *BedrockService) BuildMessagesWithContext(
	userMessage string,
	fileContext string,
) []ClaudeMessage {
	// If there's file context, prepend it to the user message
	fullMessage := userMessage
	if fileContext != "" {
		fullMessage = fileContext + "\n\n" + userMessage
	}

	return []ClaudeMessage{
		{
			Role:    "user",
			Content: fullMessage,
		},
	}
}
