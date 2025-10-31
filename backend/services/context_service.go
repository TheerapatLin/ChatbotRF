package services

import (
	"chatbot/models"
	"chatbot/repositories"
	"fmt"
	"strings"

	"github.com/sashabaranov/go-openai"
)

// ContextService handles conversation context management
type ContextService struct {
	messageRepo *repositories.MessageRepository
}

// NewContextService creates a new context service
func NewContextService(messageRepo *repositories.MessageRepository) *ContextService {
	return &ContextService{
		messageRepo: messageRepo,
	}
}

// BuildContextWithHistory builds OpenAI messages array with conversation history
func (s *ContextService) BuildContextWithHistory(
	sessionID string,
	systemPrompt string,
	currentMessage string,
	historyLimit int,
) ([]openai.ChatCompletionMessage, error) {

	messages := []openai.ChatCompletionMessage{}

	// 1. Add system prompt
	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// 2. Retrieve conversation history
	if sessionID != "" && historyLimit > 0 {
		history, err := s.messageRepo.GetRecentBySession(sessionID, historyLimit)
		if err != nil {
			// Log error but don't fail - continue without history
			fmt.Printf("⚠️  Failed to retrieve conversation history: %v\n", err)
		} else {
			// Add historical messages
			for _, msg := range history {
				role := openai.ChatMessageRoleUser
				if msg.Role == models.RoleAssistant {
					role = openai.ChatMessageRoleAssistant
				} else if msg.Role == models.RoleSystem {
					role = openai.ChatMessageRoleSystem
				}

				messages = append(messages, openai.ChatCompletionMessage{
					Role:    role,
					Content: msg.Content,
				})
			}
		}
	}

	// 3. Add current message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: currentMessage,
	})

	return messages, nil
}

// GetConversationSummary returns a summary of the conversation
func (s *ContextService) GetConversationSummary(sessionID string) (string, error) {
	messages, err := s.messageRepo.GetAllBySession(sessionID)
	if err != nil {
		return "", err
	}

	if len(messages) == 0 {
		return "New conversation", nil
	}

	var summary strings.Builder
	summary.WriteString(fmt.Sprintf("Conversation with %d messages\n", len(messages)))

	// Count by role
	userCount := 0
	assistantCount := 0
	for _, msg := range messages {
		if msg.Role == models.RoleUser {
			userCount++
		} else if msg.Role == models.RoleAssistant {
			assistantCount++
		}
	}

	summary.WriteString(fmt.Sprintf("User messages: %d\n", userCount))
	summary.WriteString(fmt.Sprintf("Assistant messages: %d\n", assistantCount))

	return summary.String(), nil
}

// BuildEnhancedSystemPrompt creates an enhanced system prompt with context
func (s *ContextService) BuildEnhancedSystemPrompt(
	basePrompt string,
	sessionID string,
) string {
	var enhanced strings.Builder

	// Start with base prompt
	enhanced.WriteString(basePrompt)

	// Add conversation context if available
	if sessionID != "" {
		summary, err := s.GetConversationSummary(sessionID)
		if err == nil && summary != "" {
			enhanced.WriteString("\n\n--- Conversation Context ---\n")
			enhanced.WriteString(summary)
		}
	}

	// Add instruction for using history
	enhanced.WriteString("\n\nUse the conversation history to provide contextual and personalized responses.")

	return enhanced.String()
}