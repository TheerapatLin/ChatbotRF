package controllers

import (
	"context"
	"fmt"
	"io"
	"log"

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
	contextService   *services.ContextService
}

// NewWebSocketController creates a new WebSocket controller
func NewWebSocketController(
	messageRepo *repositories.MessageRepository,
	personaRepo *repositories.PersonaRepository,
	fileAnalysisRepo *repositories.FileAnalysisRepository,
	openaiService *services.OpenAIService,
) *WebSocketController {
	return &WebSocketController{
		messageRepo:      messageRepo,
		personaRepo:      personaRepo,
		fileAnalysisRepo: fileAnalysisRepo,
		openaiService:    openaiService,
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

	// 3. Determine system prompt (custom prompt takes precedence over persona prompt)
	systemPrompt := persona.SystemPrompt
	if msg.SystemPrompt != "" {
		// Use custom system prompt from frontend
		// Append to persona's base prompt for context
		systemPrompt = persona.SystemPrompt + "\n\nAdditional instructions: " + msg.SystemPrompt
	}

	// 4. Build context with history and files
	var messages []openai.ChatCompletionMessage
	if len(msg.FileIDs) > 0 {
		// Use new function that supports images
		messages, err = ctrl.contextService.BuildContextWithHistoryAndFiles(
			msg.SessionID,
			systemPrompt,
			msg.Content,
			10, // history limit
			msg.FileIDs,
		)
		if err != nil {
			log.Printf("⚠️  Failed to build context with files: %v", err)
			// Fallback to simple messages
			messages = []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: msg.Content},
			}
		}
	} else if msg.SessionID != "" {
		// Use context service to build messages with history only
		messages, err = ctrl.contextService.BuildContextWithHistory(
			msg.SessionID,
			systemPrompt,
			msg.Content,
			10, // history limit
		)
		if err != nil {
			log.Printf("⚠️  Failed to build context with history: %v", err)
			// Fallback to simple messages
			messages = []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: msg.Content},
			}
		}
	} else {
		// No session and no files - simple messages
		messages = []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
			{Role: openai.ChatMessageRoleUser, Content: msg.Content},
		}
	}

	fmt.Println(`messages => `, messages)

	// 5. Call OpenAI streaming API
	stream, err := ctrl.openaiService.CreateStreamingChatCompletion(ctx, messages)
	if err != nil {
		return fmt.Errorf("failed to create streaming request: %w", err)
	}
	defer stream.Close()

	// 6. Stream the response to client
	fullContent := ""

	for {
		select {
		case <-ctx.Done():
			// Request was cancelled
			return nil

		default:
			response, err := stream.Recv()
			if err == io.EOF {
				// Stream completed
				goto streamDone
			}
			if err != nil {
				return fmt.Errorf("stream error: %w", err)
			}

			// Get delta content
			if len(response.Choices) > 0 {
				delta := response.Choices[0].Delta.Content
				if delta != "" {
					fullContent += delta

					// Send chunk to client
					if err := ctrl.sendChunk(c, delta, false); err != nil {
						return err
					}
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
