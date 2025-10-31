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
	messageRepo      *repositories.MessageRepository
	fileAnalysisRepo *repositories.FileAnalysisRepository
}

// NewContextService creates a new context service
func NewContextService(
	messageRepo *repositories.MessageRepository,
	fileAnalysisRepo *repositories.FileAnalysisRepository,
) *ContextService {
	return &ContextService{
		messageRepo:      messageRepo,
		fileAnalysisRepo: fileAnalysisRepo,
	}
}

// BuildContextWithHistory builds OpenAI messages array with conversation history (existing method)
func (s *ContextService) BuildContextWithHistory(
	sessionID string,
	systemPrompt string,
	currentMessage string,
	historyLimit int,
) ([]openai.ChatCompletionMessage, error) {
	// Call the new method with empty file IDs
	return s.BuildContextWithFiles(sessionID, systemPrompt, currentMessage, []string{}, historyLimit)
}

// BuildContextWithFiles builds OpenAI messages with file context
func (s *ContextService) BuildContextWithFiles(
	sessionID string,
	systemPrompt string,
	userMessage string,
	fileIDs []string,
	historyLimit int,
) ([]openai.ChatCompletionMessage, error) {

	messages := []openai.ChatCompletionMessage{}

	// 1. Add system prompt if provided
	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// 2. Build and add file context if files are provided
	if len(fileIDs) > 0 {
		fileContext, err := s.buildFileContext(fileIDs)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to build file context: %v\n", err)
		} else if fileContext != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: fileContext,
			})
		}
	}

	// 3. Add conversation history
	if sessionID != "" && historyLimit > 0 {
		history, err := s.messageRepo.GetRecentBySession(sessionID, historyLimit)
		if err == nil && len(history) > 0 {
			for _, msg := range history {
				role := openai.ChatMessageRoleUser
				if msg.Role == models.RoleAssistant {
					role = openai.ChatMessageRoleAssistant
				}

				content := msg.Content

				// If message has file attachments, add file info to content
				if msg.HasFileAttachments() {
					attachments, _ := msg.GetFileAttachments()
					fileInfo := s.formatFileAttachmentsInfo(attachments)
					content = content + "\n\n" + fileInfo
				}

				messages = append(messages, openai.ChatCompletionMessage{
					Role:    role,
					Content: content,
				})
			}
		}
	}

	// 4. Add current user message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: userMessage,
	})

	return messages, nil
}

// buildFileContext fetches file analysis and builds context string
func (s *ContextService) buildFileContext(fileIDs []string) (string, error) {
	if len(fileIDs) == 0 {
		return "", nil
	}

	var contextParts []string
	contextParts = append(contextParts,
		fmt.Sprintf("📎 The user has provided %d file(s) for you to analyze:", len(fileIDs)))
	contextParts = append(contextParts, "")

	successCount := 0
	for i, fileID := range fileIDs {
		analysis, err := s.fileAnalysisRepo.FindByID(fileID)
		if err != nil {
			contextParts = append(contextParts,
				fmt.Sprintf("⚠️  File %d: [Error: File not found - ID: %s]", i+1, fileID))
			continue
		}

		successCount++
		contextParts = append(contextParts,
			fmt.Sprintf("--- 📄 File %d: %s ---", i+1, analysis.FileName))
		contextParts = append(contextParts,
			fmt.Sprintf("Type: %s", analysis.FileType))
		contextParts = append(contextParts,
			fmt.Sprintf("Size: %s", formatFileSize(analysis.FileSize)))
		contextParts = append(contextParts, "")
		contextParts = append(contextParts, "Analysis Result:")
		contextParts = append(contextParts, analysis.Analysis)
		contextParts = append(contextParts, "")
		contextParts = append(contextParts, strings.Repeat("-", 60))
		contextParts = append(contextParts, "")
	}

	if successCount == 0 {
		return "", fmt.Errorf("no valid files found")
	}

	contextParts = append(contextParts,
		"📌 Instructions: Please reference the above file(s) when answering the user's question. "+
			"Provide insights based on the file content and analysis.")

	return strings.Join(contextParts, "\n"), nil
}

// formatFileAttachmentsInfo formats file attachment info for conversation history
func (s *ContextService) formatFileAttachmentsInfo(attachments []models.FileAttachment) string {
	if len(attachments) == 0 {
		return ""
	}

	parts := []string{"[📎 Attached files:"}
	for _, att := range attachments {
		parts = append(parts,
			fmt.Sprintf("  • %s (%s, %s)",
				att.Filename,
				att.FileType,
				formatFileSize(att.FileSize)))
	}
	parts = append(parts, "]")

	return strings.Join(parts, "\n")
}

// formatFileSize formats bytes to human-readable size
func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(GB))
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(MB))
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(KB))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
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
