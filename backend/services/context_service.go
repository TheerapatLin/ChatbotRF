package services

import (
	"chatbot/models"
	"chatbot/repositories"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
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

// BuildContextWithHistory builds OpenAI messages array with conversation history
func (s *ContextService) BuildContextWithHistory(
	sessionID string,
	systemPrompt string,
	currentMessage string,
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

	// 2. Add conversation history
	if sessionID != "" && historyLimit > 0 {
		history, err := s.messageRepo.GetRecentBySession(sessionID, historyLimit)
		if err == nil && len(history) > 0 {
			for _, msg := range history {
				role := openai.ChatMessageRoleUser
				if msg.Role == models.RoleAssistant {
					role = openai.ChatMessageRoleAssistant
				}

				messages = append(messages, openai.ChatCompletionMessage{
					Role:    role,
					Content: msg.Content,
				})
			}
		}
	}

	// 3. Add current user message
	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: currentMessage,
	})

	return messages, nil
}

// BuildContextWithHistoryAndFiles builds OpenAI messages with file support (including images)
func (s *ContextService) BuildContextWithHistoryAndFiles(
	sessionID string,
	systemPrompt string,
	currentMessage string,
	historyLimit int,
	fileIDs []string,
) ([]openai.ChatCompletionMessage, error) {

	messages := []openai.ChatCompletionMessage{}

	// 1. Add system prompt if provided
	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// 2. Add conversation history
	if sessionID != "" && historyLimit > 0 {
		history, err := s.messageRepo.GetRecentBySession(sessionID, historyLimit)
		if err == nil && len(history) > 0 {
			for _, msg := range history {
				role := openai.ChatMessageRoleUser
				if msg.Role == models.RoleAssistant {
					role = openai.ChatMessageRoleAssistant
				}

				messages = append(messages, openai.ChatCompletionMessage{
					Role:    role,
					Content: msg.Content,
				})
			}
		}
	}

	// 3. Add file context as system message if files exist
	if len(fileIDs) > 0 {
		fileContext, err := s.BuildFileContext(fileIDs)
		if err == nil && fileContext != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: fileContext,
			})
		}
	}

	// 4. Add current user message with image support
	if len(fileIDs) > 0 && s.hasImageFiles(fileIDs) {
		// Build multipart message with text and images
		multiContent := []openai.ChatMessagePart{
			{
				Type: openai.ChatMessagePartTypeText,
				Text: currentMessage,
			},
		}

		// Add images
		for _, fileID := range fileIDs {
			fileRecord, err := s.fileAnalysisRepo.FindByID(fileID)
			if err != nil {
				continue
			}

			// Check if it's an image
			if strings.HasPrefix(fileRecord.MimeType, "image/") {
				// Read and encode image to base64
				imageData, err := os.ReadFile(fileRecord.StoragePath)
				if err != nil {
					continue
				}

				base64Image := base64.StdEncoding.EncodeToString(imageData)
				dataURL := fmt.Sprintf("data:%s;base64,%s", fileRecord.MimeType, base64Image)

				multiContent = append(multiContent, openai.ChatMessagePart{
					Type: openai.ChatMessagePartTypeImageURL,
					ImageURL: &openai.ChatMessageImageURL{
						URL:    dataURL,
						Detail: openai.ImageURLDetailAuto,
					},
				})
			}
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:         openai.ChatMessageRoleUser,
			MultiContent: multiContent,
		})
	} else {
		// Simple text message
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: currentMessage,
		})
	}

	return messages, nil
}

// hasImageFiles checks if any of the files are images
func (s *ContextService) hasImageFiles(fileIDs []string) bool {
	for _, fileID := range fileIDs {
		fileRecord, err := s.fileAnalysisRepo.FindByID(fileID)
		if err != nil {
			continue
		}
		if strings.HasPrefix(fileRecord.MimeType, "image/") {
			return true
		}
	}
	return false
}

// BuildFileContext builds file context string from file IDs (for current message only)
func (s *ContextService) BuildFileContext(fileIDs []string) (string, error) {
	if len(fileIDs) == 0 {
		return "", nil
	}

	var contextParts []string
	contextParts = append(contextParts,
		fmt.Sprintf("📎 The user has provided %d file(s) for you to analyze:", len(fileIDs)))
	contextParts = append(contextParts, "")

	successCount := 0
	for i, fileID := range fileIDs {
		fileRecord, err := s.fileAnalysisRepo.FindByID(fileID)
		if err != nil {
			contextParts = append(contextParts,
				fmt.Sprintf("⚠️  File %d: [Error: File not found - ID: %s]", i+1, fileID))
			continue
		}

		successCount++
		contextParts = append(contextParts,
			fmt.Sprintf("--- 📄 File %d: %s ---", i+1, fileRecord.FileName))
		contextParts = append(contextParts,
			fmt.Sprintf("MIME Type: %s", fileRecord.MimeType))
		contextParts = append(contextParts,
			fmt.Sprintf("Size: %s", formatFileSize(fileRecord.FileSize)))
		contextParts = append(contextParts, "")

		// Read and include file content for supported text-based files
		fileContent, err := s.readFileContent(fileRecord.StoragePath, fileRecord.MimeType)
		if err == nil && fileContent != "" {
			contextParts = append(contextParts, "📝 File Content:")
			contextParts = append(contextParts, "```")
			contextParts = append(contextParts, fileContent)
			contextParts = append(contextParts, "```")
		} else if err != nil {
			contextParts = append(contextParts, fmt.Sprintf("⚠️  Could not read file content: %v", err))
		} else {
			contextParts = append(contextParts, "Note: File is stored on server (binary/unsupported format)")
		}

		contextParts = append(contextParts, "")
		contextParts = append(contextParts, strings.Repeat("-", 60))
		contextParts = append(contextParts, "")
	}

	if successCount == 0 {
		return "", fmt.Errorf("no valid files found")
	}

	contextParts = append(contextParts,
		"📌 Instructions: Please analyze the file content above and provide insights based on what you see. "+
			"Answer the user's question using the information from these files.")

	return strings.Join(contextParts, "\n"), nil
}

// readFileContent reads the content of a file based on its MIME type
func (s *ContextService) readFileContent(filePath string, mimeType string) (string, error) {
	// Maximum file size to read (5MB for documents, 1MB for text)
	maxFileSize := int64(1024 * 1024) // 1MB default

	// Check file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to get file info: %w", err)
	}

	// Handle different file types
	switch mimeType {
	case "application/pdf":
		maxFileSize = 5 * 1024 * 1024 // 5MB for PDF
		if fileInfo.Size() > maxFileSize {
			return "", fmt.Errorf("PDF file too large (max 5MB)")
		}
		return s.extractPDFText(filePath)

	case "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		maxFileSize = 5 * 1024 * 1024 // 5MB for DOCX
		if fileInfo.Size() > maxFileSize {
			return "", fmt.Errorf("DOCX file too large (max 5MB)")
		}
		return s.extractDOCXText(filePath)

	case "text/plain", "text/markdown", "text/csv", "application/json",
		"text/html", "text/css", "text/javascript", "application/xml", "text/xml":
		// Text-based files
		if fileInfo.Size() > maxFileSize {
			return "", fmt.Errorf("file too large (max 1MB)")
		}
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read file: %w", err)
		}
		return string(content), nil

	default:
		// Check if it starts with "text/"
		if strings.HasPrefix(mimeType, "text/") {
			if fileInfo.Size() > maxFileSize {
				return "", fmt.Errorf("file too large (max 1MB)")
			}
			content, err := os.ReadFile(filePath)
			if err != nil {
				return "", fmt.Errorf("failed to read file: %w", err)
			}
			return string(content), nil
		}
		return "", nil // Not supported, not an error
	}
}

// extractPDFText extracts text content from a PDF file
func (s *ContextService) extractPDFText(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer f.Close()

	var textContent strings.Builder
	totalPages := r.NumPage()

	// Limit to first 50 pages to avoid huge content
	maxPages := totalPages
	if maxPages > 50 {
		maxPages = 50
	}

	for pageNum := 1; pageNum <= maxPages; pageNum++ {
		p := r.Page(pageNum)
		if p.V.IsNull() {
			continue
		}

		text, err := p.GetPlainText(nil)
		if err != nil {
			continue // Skip pages with errors
		}

		textContent.WriteString(fmt.Sprintf("\n--- Page %d ---\n", pageNum))
		textContent.WriteString(text)
	}

	if totalPages > maxPages {
		textContent.WriteString(fmt.Sprintf("\n\n[Note: Only first %d of %d pages shown]", maxPages, totalPages))
	}

	result := textContent.String()
	if result == "" {
		return "", fmt.Errorf("no text content extracted from PDF")
	}

	return result, nil
}

// extractDOCXText extracts text content from a DOCX file
func (s *ContextService) extractDOCXText(filePath string) (string, error) {
	doc, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open DOCX: %w", err)
	}
	defer doc.Close()

	docx := doc.Editable()
	text := docx.GetContent()

	if text == "" {
		return "", fmt.Errorf("no text content extracted from DOCX")
	}

	return text, nil
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
