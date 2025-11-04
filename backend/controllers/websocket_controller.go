package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"

	"chatbot/models"
	"chatbot/repositories"
	"chatbot/services"

	"github.com/gofiber/contrib/websocket"
	"github.com/sashabaranov/go-openai"
)

// WebSocketController handles WebSocket connections for streaming chat
type WebSocketController struct {
	messageRepo      *repositories.MessageRepository
	personaRepo      *repositories.PersonaRepository
	fileAnalysisRepo *repositories.FileAnalysisRepository
	openaiService    *services.OpenAIService
	bedrockService   *services.BedrockService
	contextService   *services.ContextService
}

// NewWebSocketController creates a new WebSocket controller
func NewWebSocketController(
	messageRepo *repositories.MessageRepository,
	personaRepo *repositories.PersonaRepository,
	fileAnalysisRepo *repositories.FileAnalysisRepository,
	openaiService *services.OpenAIService,
	bedrockService *services.BedrockService,
) *WebSocketController {
	return &WebSocketController{
		messageRepo:      messageRepo,
		personaRepo:      personaRepo,
		fileAnalysisRepo: fileAnalysisRepo,
		openaiService:    openaiService,
		bedrockService:   bedrockService,
		contextService:   services.NewContextService(messageRepo, fileAnalysisRepo),
	}
}

// WSMessage represents incoming WebSocket messages
type WSMessage struct {
	Type         string   `json:"type"`          // "message"
	Content      string   `json:"content"`       // User message
	PersonaID    *int     `json:"persona_id"`    // Persona ID
	SystemPrompt string   `json:"system_prompt"` // Optional custom system prompt
	SessionID    string   `json:"session_id"`    // Session ID for conversation history
	FileIDs      []string `json:"file_ids"`      // File IDs for current message only
}

// WSResponse represents outgoing WebSocket messages
type WSResponse struct {
	Type       string `json:"type"`                  // "chunk"
	Content    string `json:"content"`               // Chunk content
	Done       bool   `json:"done"`                  // Is streaming done?
	MessageID  string `json:"message_id,omitempty"`  // Message ID (when done)
	TokensUsed int    `json:"tokens_used,omitempty"` // Tokens used (when done)
}

// HandleStreamingChat handles WebSocket connections for streaming chat
func (ctrl *WebSocketController) HandleStreamingChat(c *websocket.Conn) {
	log.Printf("New WebSocket connection from %s", c.RemoteAddr())
	// ***c.RemoteAddr() คือ ฟังก์ชันที่ใช้ดึง IP address ของ client (ผู้ที่ส่ง request เข้ามา)

	// Track connection context for cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Clean up on disconnect
	defer func() {
		log.Printf("WebSocket connection closed from %s", c.RemoteAddr())
		c.Close()
	}()

	// Message loop
	for {
		var msg WSMessage

		// Read message from client
		if err := c.ReadJSON(&msg); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle message type
		if msg.Type == "message" {
			if err := ctrl.handleMessage(ctx, c, msg); err != nil {
				log.Printf("Error handling message: %v", err)
				ctrl.sendError(c, err.Error())
			}
		} else {
			ctrl.sendError(c, fmt.Sprintf("Unknown message type: %s", msg.Type))
		}
	}
}

// handleMessage processes a streaming chat request
func (ctrl *WebSocketController) handleMessage(ctx context.Context, c *websocket.Conn, msg WSMessage) error {
	// 1. Validate message
	if msg.Content == "" {
		return fmt.Errorf("content is required")
	}

	// 2. Get persona if persona_id provided (default to 1 if not specified)
	personaID := 1
	if msg.PersonaID != nil {
		personaID = *msg.PersonaID
	}

	persona, err := ctrl.personaRepo.FindByID(personaID)
	if err != nil {
		return fmt.Errorf("persona with ID %d not found", personaID)
	}

	// 3. Determine system prompt (append custom prompt to persona's base prompt)
	systemPrompt := persona.SystemPrompt
	if msg.SystemPrompt != "" {
		// Append custom system_prompt to persona's prompt (don't replace!)
		systemPrompt = systemPrompt + "\n\n--- Additional Instructions ---\n" + msg.SystemPrompt
	}

	// 4. Determine which AI provider to use (auto-detection)
	var streamingService services.StreamingChatService
	var providerName string

	if ctrl.openaiService != nil && ctrl.openaiService.IsAvailable() {
		streamingService = ctrl.openaiService
		providerName = "OpenAI"
		log.Printf("🟢 Using OpenAI for WebSocket streaming")
	} else if ctrl.bedrockService != nil && ctrl.bedrockService.IsAvailable() {
		streamingService = ctrl.bedrockService
		providerName = "Bedrock"
		log.Printf("🔵 Using AWS Bedrock for WebSocket streaming")
	} else {
		return fmt.Errorf("no AI provider available (check OPENAI_API_KEY or AWS credentials in .env)")
	}

	// 5. Build context with history and files
	var messages interface{}
	if len(msg.FileIDs) > 0 {
		if providerName == "Bedrock" {
			// Bedrock: Build messages manually with file context
			messages = ctrl.buildBedrockMessagesWithFiles(msg.SessionID, systemPrompt, msg.Content, msg.FileIDs)
		} else {
			// OpenAI: Use existing function that supports images
			openaiMessages, err := ctrl.contextService.BuildContextWithHistoryAndFiles(
				msg.SessionID,
				systemPrompt,
				msg.Content,
				10, // history limit
				msg.FileIDs,
			)
			if err != nil {
				log.Printf("⚠️  Failed to build context with files: %v", err)
				// Fallback to simple messages
				openaiMessages = []openai.ChatCompletionMessage{
					{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
					{Role: openai.ChatMessageRoleUser, Content: msg.Content},
				}
			}
			messages = openaiMessages
		}
	} else if msg.SessionID != "" {
		// Use context service to build messages with history only
		openaiMessages, err := ctrl.contextService.BuildContextWithHistory(
			msg.SessionID,
			systemPrompt,
			msg.Content,
			10, // history limit
		)
		if err != nil {
			log.Printf("⚠️  Failed to build context with history: %v", err)
			// Fallback to simple messages
			openaiMessages = []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: msg.Content},
			}
		}

		// Convert to provider-specific format
		if providerName == "Bedrock" {
			messages = ctrl.convertToClaudeMessages(openaiMessages, systemPrompt)
		} else {
			messages = openaiMessages
		}
	} else {
		// No session and no files - simple messages
		if providerName == "Bedrock" {
			messages = []services.ClaudeMessage{
				{Role: "user", Content: msg.Content},
			}
		} else {
			messages = []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: msg.Content},
			}
		}
	}

	// 6. Create streaming request using unified interface
	streamReq := services.StreamingChatRequest{
		Messages:     messages,
		SystemPrompt: systemPrompt,
		Temperature:  0.7,
		MaxTokens:    2000,
	}

	stream, err := streamingService.CreateStreamingChat(ctx, streamReq)
	if err != nil {
		return fmt.Errorf("failed to create streaming request: %w", err)
	}
	defer stream.Close()

	// 7. Stream the response to client
	fullContent := ""

	for {
		select {
		case <-ctx.Done():
			// Request was cancelled
			return nil

		default:
			chunk, err := stream.Recv()
			if err == io.EOF {
				// Stream completed
				goto streamDone
			}
			if err != nil {
				return fmt.Errorf("stream error: %w", err)
			}

			// Add chunk to full content
			if chunk != "" {
				fullContent += chunk

				// Send chunk to client
				if err := ctrl.sendChunk(c, chunk, false); err != nil {
					return err
				}
			}
		}
	}

streamDone:
	// 7. Save user message to database
	userMessage := &models.Message{
		SessionID: msg.SessionID,
		Role:      models.RoleUser,
		Content:   msg.Content,
		PersonaID: &personaID,
	}

	if err := ctrl.messageRepo.Create(userMessage); err != nil {
		log.Printf("Failed to save user message: %v", err)
	}

	// 8. Save assistant response to database
	// Calculate tokens (simplified - in real app would get from OpenAI response)
	tokensUsed := len(fullContent) / 4 // Rough estimate

	assistantMessage := &models.Message{
		SessionID:  msg.SessionID,
		Role:       models.RoleAssistant,
		Content:    fullContent,
		PersonaID:  &personaID,
		TokensUsed: &tokensUsed,
	}

	if err := ctrl.messageRepo.Create(assistantMessage); err != nil {
		log.Printf("Failed to save assistant message: %v", err)
	}

	// 9. Send completion message
	if err := ctrl.sendDone(c, assistantMessage.ID.String(), tokensUsed); err != nil {
		return err
	}

	return nil
}

// sendChunk sends a content chunk to the client
// ส่งข้อมูลเป็นส่วนๆ (chunk)
func (ctrl *WebSocketController) sendChunk(c *websocket.Conn, content string, done bool) error {
	return c.WriteJSON(WSResponse{
		Type:    "chunk",
		Content: content,
		Done:    done,
	})
}

// sendDone sends the final completion message
func (ctrl *WebSocketController) sendDone(c *websocket.Conn, messageID string, tokensUsed int) error {
	return c.WriteJSON(WSResponse{
		Type:       "chunk",
		Content:    "",
		Done:       true,
		MessageID:  messageID,
		TokensUsed: tokensUsed,
	})
}

// sendError sends an error message to the client
func (ctrl *WebSocketController) sendError(c *websocket.Conn, errorMsg string) error {
	return c.WriteJSON(map[string]interface{}{
		"type":  "error",
		"error": errorMsg,
	})
}

// convertToClaudeMessages converts OpenAI messages to Claude message format
// Claude uses separate system prompt field, so we filter out system messages
// Claude also requires alternating user/assistant messages, so we merge consecutive messages
func (ctrl *WebSocketController) convertToClaudeMessages(openaiMessages []openai.ChatCompletionMessage, systemPrompt string) []services.ClaudeMessage {
	claudeMessages := make([]services.ClaudeMessage, 0)
	var lastRole string
	var accumulatedContent string

	for _, msg := range openaiMessages {
		// Skip system messages (they're passed separately in systemPrompt)
		if msg.Role == openai.ChatMessageRoleSystem {
			continue
		}

		// Skip messages with empty content
		if msg.Content == "" {
			continue
		}

		// Convert role
		role := "user"
		if msg.Role == openai.ChatMessageRoleAssistant {
			role = "assistant"
		}

		// If same role as last message, merge content
		if lastRole == role {
			accumulatedContent += "\n\n" + msg.Content
		} else {
			// Save previous accumulated message
			if accumulatedContent != "" {
				claudeMessages = append(claudeMessages, services.ClaudeMessage{
					Role:    lastRole,
					Content: accumulatedContent,
				})
			}

			// Start new message
			lastRole = role
			accumulatedContent = msg.Content
		}
	}

	// Add final accumulated message
	if accumulatedContent != "" {
		claudeMessages = append(claudeMessages, services.ClaudeMessage{
			Role:    lastRole,
			Content: accumulatedContent,
		})
	}

	// Ensure alternating pattern and starts with user
	claudeMessages = ctrl.ensureAlternatingMessages(claudeMessages)

	return claudeMessages
}

// ensureAlternatingMessages ensures messages alternate between user and assistant
// Claude API requires this pattern
func (ctrl *WebSocketController) ensureAlternatingMessages(messages []services.ClaudeMessage) []services.ClaudeMessage {
	if len(messages) == 0 {
		return messages
	}

	// Must start with user message
	if messages[0].Role != "user" {
		log.Printf("⚠️  First message is not 'user', prepending empty user message")
		messages = append([]services.ClaudeMessage{{Role: "user", Content: "[conversation started]"}}, messages...)
	}

	return messages
}

// isImageFile checks if a MIME type is an image
func (ctrl *WebSocketController) isImageFile(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

// readImageFile reads an image file and returns base64 encoded data
func (ctrl *WebSocketController) readImageFile(filePath string, mimeType string) (string, string, error) {
	fmt.Println(`filePath =>`, filePath)
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %w", err)
	}

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(data)

	return encoded, mimeType, nil
}

// buildBedrockMessagesWithFiles builds Claude messages with file context (text files + images)
func (ctrl *WebSocketController) buildBedrockMessagesWithFiles(
	sessionID string,
	systemPrompt string,
	currentMessage string,
	fileIDs []string,
) []services.ClaudeMessage {
	claudeMessages := make([]services.ClaudeMessage, 0)

	// 1. Add conversation history
	if sessionID != "" {
		history, err := ctrl.messageRepo.GetRecentBySession(sessionID, 10)
		if err == nil && len(history) > 0 {
			for _, msg := range history {
				role := "user"
				if msg.Role == models.RoleAssistant {
					role = "assistant"
				}

				// Skip empty messages
				if msg.Content == "" {
					continue
				}

				claudeMessages = append(claudeMessages, services.ClaudeMessage{
					Role:    role,
					Content: msg.Content,
				})
			}
		}
	}

	// 2. Separate image files from text files
	imageFileIDs := make([]string, 0)
	textFileIDs := make([]string, 0)

	for _, fileID := range fileIDs {
		file, err := ctrl.fileAnalysisRepo.FindByID(fileID)
		if err != nil {
			log.Printf("⚠️  File not found: %s", fileID)
			continue
		}

		if ctrl.isImageFile(file.MimeType) {
			imageFileIDs = append(imageFileIDs, fileID)
		} else {
			textFileIDs = append(textFileIDs, fileID)
		}
	}

	// 3. Build file context for text files
	var fileContext string
	if len(textFileIDs) > 0 {
		var err error
		fileContext, err = ctrl.contextService.BuildFileContext(textFileIDs)
		if err != nil {
			log.Printf("⚠️  Failed to build file context: %v", err)
		}
	}

	// 4. Add current user message with file context
	var userMessageContent string
	if fileContext != "" {
		userMessageContent = fileContext + "\n\n--- User Question ---\n" + currentMessage
	} else {
		userMessageContent = currentMessage
	}

	// 5. If there are images, create multimodal content blocks
	if len(imageFileIDs) > 0 {
		contentBlocks := make([]services.ClaudeContentBlock, 0)

		// Add text content first
		contentBlocks = append(contentBlocks, services.ClaudeContentBlock{
			Type: "text",
			Text: userMessageContent,
		})

		// Add image blocks
		for _, fileID := range imageFileIDs {
			file, err := ctrl.fileAnalysisRepo.FindByID(fileID)
			if err != nil {
				log.Printf("⚠️  File not found: %s", fileID)
				continue
			}

			// Read image file and encode to base64
			imageData, mediaType, err := ctrl.readImageFile(file.StoragePath, file.MimeType)
			if err != nil {
				log.Printf("⚠️  Failed to read image: %v", err)
				continue
			}

			// Add image content block
			contentBlocks = append(contentBlocks, services.ClaudeContentBlock{
				Type: "image",
				Source: &services.ClaudeImageSource{
					Type:      "base64",
					MediaType: mediaType,
					Data:      imageData,
				},
			})

			log.Printf("✓ Added image to content: %s (%s)", file.FileName, mediaType)
		}

		claudeMessages = append(claudeMessages, services.ClaudeMessage{
			Role:    "user",
			Content: contentBlocks,
		})
	} else {
		// Simple text message
		claudeMessages = append(claudeMessages, services.ClaudeMessage{
			Role:    "user",
			Content: userMessageContent,
		})
	}

	// 6. Ensure alternating pattern and starts with user
	claudeMessages = ctrl.ensureAlternatingMessages(claudeMessages)

	return claudeMessages
}
