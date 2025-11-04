package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"chatbot/models"
	"chatbot/repositories"
	"chatbot/services"

	"github.com/gofiber/fiber/v2"
)

// BedrockController handles Bedrock-specific chat requests
type BedrockController struct {
	bedrockService   *services.BedrockService
	personaRepo      *repositories.PersonaRepository
	messageRepo      *repositories.MessageRepository
	contextService   *services.ContextService
	fileAnalysisRepo *repositories.FileAnalysisRepository
}

// NewBedrockController creates a new Bedrock controller
func NewBedrockController(
	bedrockService *services.BedrockService,
	personaRepo *repositories.PersonaRepository,
	messageRepo *repositories.MessageRepository,
	contextService *services.ContextService,
	fileAnalysisRepo *repositories.FileAnalysisRepository,
) *BedrockController {
	return &BedrockController{
		bedrockService:   bedrockService,
		personaRepo:      personaRepo,
		messageRepo:      messageRepo,
		contextService:   contextService,
		fileAnalysisRepo: fileAnalysisRepo,
	}
}

// SendBedrockMessage handles POST /api/chat/bedrock
type BedrockMessageRequest struct {
	Message      string   `json:"message" validate:"required"`
	PersonaID    uint     `json:"persona_id" validate:"required"`
	Temperature  float64  `json:"temperature,omitempty"`
	MaxTokens    int      `json:"max_tokens,omitempty"`
	SystemPrompt string   `json:"system_prompt,omitempty"`
	SessionID    string   `json:"session_id,omitempty"`
	UseHistory   bool     `json:"use_history,omitempty"`
	FileIDs      []string `json:"file_ids,omitempty"` // File IDs for current message only
}

type BedrockMessageResponse struct {
	MessageID string `json:"message_id"`
	SessionID string `json:"session_id"`
	Reply     string `json:"reply"`
	Persona   struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Expertise string `json:"expertise"`
		Icon      string `json:"icon"`
	} `json:"persona"`
	TokensUsed int    `json:"tokens_used"`
	Model      string `json:"model"`
	Provider   string `json:"provider"`
	Timestamp  string `json:"timestamp"`
}

func (bc *BedrockController) SendBedrockMessage(c *fiber.Ctx) error {
	// Parse request
	var req BedrockMessageRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("❌ Failed to parse request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message is required",
		})
	}

	if req.PersonaID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Persona ID is required",
		})
	}

	// Generate session ID if not provided
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = fmt.Sprintf("bedrock-%d-%d", req.PersonaID, time.Now().Unix())
	}

	log.Printf("📨 Bedrock Chat Request: session=%s, persona_id=%d, use_history=%v, files=%d",
		sessionID, req.PersonaID, req.UseHistory, len(req.FileIDs))

	// Get persona from database
	persona, err := bc.personaRepo.FindByID(int(req.PersonaID))
	if err != nil {
		log.Printf("❌ Persona not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Persona not found",
		})
	}

	// Build system prompt (combine persona system_prompt + persona description)
	systemPrompt := fmt.Sprintf("%s\n\nYou are %s, a %s.",
		persona.SystemPrompt, persona.Name, persona.Description)

	if req.SystemPrompt != "" {
		systemPrompt += "\n\n" + req.SystemPrompt
	}

	// Build messages with history and file context
	var messages []services.ClaudeMessage
	var fileContext string
	var hasImages bool

	// Get file context if files provided
	if len(req.FileIDs) > 0 {
		// Separate image files from text files
		imageFileIDs := make([]string, 0)
		textFileIDs := make([]string, 0)

		for _, fileID := range req.FileIDs {
			file, err := bc.fileAnalysisRepo.FindByID(fileID)
			if err != nil {
				log.Printf("⚠️ File not found: %s", fileID)
				continue
			}

			if bc.isImageFile(file.MimeType) {
				imageFileIDs = append(imageFileIDs, fileID)
				hasImages = true
			} else {
				textFileIDs = append(textFileIDs, fileID)
			}
		}

		// Build text file context (for non-image files)
		if len(textFileIDs) > 0 {
			var err error
			fileContext, err = bc.contextService.BuildFileContext(textFileIDs)
			if err != nil {
				log.Printf("⚠️ Failed to build file context: %v", err)
			}
		}

		// Store image file IDs for later multimodal processing
		if hasImages {
			log.Printf("📸 Found %d image files", len(imageFileIDs))
		}
	}

	// Build messages based on history and file support
	if req.UseHistory && sessionID != "" {
		// Get conversation history
		history, err := bc.messageRepo.GetRecentBySession(sessionID, 10)
		if err == nil && len(history) > 0 {
			// Convert history to Claude messages
			claudeHistory := make([]services.ClaudeMessage, 0, len(history))
			for _, msg := range history {
				role := "user"
				if msg.Role == models.RoleAssistant {
					role = "assistant"
				}
				claudeHistory = append(claudeHistory, services.ClaudeMessage{
					Role:    role,
					Content: msg.Content,
				})
			}

			// Build with history
			if fileContext != "" {
				// Add file context to user message
				fullMessage := fileContext + "\n\n" + req.Message
				messages = bc.bedrockService.BuildMessagesWithHistory(fullMessage, claudeHistory)
			} else {
				messages = bc.bedrockService.BuildMessagesWithHistory(req.Message, claudeHistory)
			}
		} else {
			// No history found, build simple message
			if fileContext != "" {
				messages = bc.bedrockService.BuildMessagesWithContext(req.Message, fileContext)
			} else {
				messages = bc.bedrockService.BuildMessages(req.Message, systemPrompt)
			}
		}
	} else {
		// No history requested
		if fileContext != "" {
			messages = bc.bedrockService.BuildMessagesWithContext(req.Message, fileContext)
		} else {
			messages = bc.bedrockService.BuildMessages(req.Message, systemPrompt)
		}
	}

	// Convert to multimodal format if there are images
	if hasImages && len(req.FileIDs) > 0 {
		messages = bc.convertToMultimodalMessages(messages, req.FileIDs)
	}

	// Save user message to database
	personaIDInt := int(req.PersonaID)
	userMsg := &models.Message{
		SessionID: sessionID,
		Role:      "user",
		Content:   req.Message,
		PersonaID: &personaIDInt,
	}

	// Save file attachments if provided
	if len(req.FileIDs) > 0 {
		fileAttachments := make([]models.FileAttachment, 0, len(req.FileIDs))
		for _, fileID := range req.FileIDs {
			fileRecord, err := bc.fileAnalysisRepo.FindByID(fileID)
			if err == nil {
				fileAttachments = append(fileAttachments, models.FileAttachment{
					FileID:   fileID,
					Filename: fileRecord.FileName,
					FileType: fileRecord.MimeType,
					FileSize: fileRecord.FileSize,
				})
			}
		}
		if err := userMsg.SetFileAttachments(fileAttachments); err != nil {
			log.Printf("⚠️ Failed to set file attachments: %v", err)
		}
	}

	if err := bc.messageRepo.Create(userMsg); err != nil {
		log.Printf("⚠️ Failed to save user message: %v", err)
	}

	// Send request to Bedrock
	bedrockReq := services.BedrockChatRequest{
		Messages:     messages,
		SystemPrompt: systemPrompt,
		Temperature:  req.Temperature,
		MaxTokens:    req.MaxTokens,
	}

	bedrockResp, err := bc.bedrockService.SendChatRequest(bedrockReq)
	if err != nil {
		log.Printf("❌ Bedrock API error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get AI response from Bedrock",
			"details": err.Error(),
		})
	}

	// Save assistant message to database
	metadataJSON, _ := json.Marshal(map[string]interface{}{
		"model":       bedrockResp.Model,
		"provider":    "bedrock",
		"stop_reason": bedrockResp.StopReason,
		"use_history": req.UseHistory,
		"file_count":  len(req.FileIDs),
	})

	assistantMsg := &models.Message{
		SessionID:  sessionID,
		Role:       "assistant",
		Content:    bedrockResp.Content,
		PersonaID:  &personaIDInt,
		TokensUsed: &bedrockResp.TokensUsed,
		Metadata:   metadataJSON,
	}
	if err := bc.messageRepo.Create(assistantMsg); err != nil {
		log.Printf("⚠️ Failed to save assistant message: %v", err)
	}

	log.Printf("✅ Bedrock response sent: message_id=%s, tokens=%d",
		assistantMsg.ID.String(), bedrockResp.TokensUsed)

	// Build response
	response := BedrockMessageResponse{
		MessageID:  assistantMsg.ID.String(),
		SessionID:  sessionID,
		Reply:      bedrockResp.Content,
		TokensUsed: bedrockResp.TokensUsed,
		Model:      bedrockResp.Model,
		Provider:   "bedrock",
		Timestamp:  assistantMsg.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
	response.Persona.ID = uint(persona.ID)
	response.Persona.Name = persona.Name
	response.Persona.Expertise = persona.Expertise
	response.Persona.Icon = persona.Icon

	return c.JSON(response)
}

// isImageFile checks if a MIME type is an image
func (bc *BedrockController) isImageFile(mimeType string) bool {
	switch mimeType {
	case "image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp":
		return true
	default:
		return false
	}
}

// convertToMultimodalMessages converts Claude messages to support images
// For the last user message, if there are image files, convert to multimodal content blocks
func (bc *BedrockController) convertToMultimodalMessages(messages []services.ClaudeMessage, fileIDs []string) []services.ClaudeMessage {
	if len(fileIDs) == 0 || len(messages) == 0 {
		return messages
	}

	// Find the last user message
	lastUserIdx := -1
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == "user" {
			lastUserIdx = i
			break
		}
	}

	if lastUserIdx == -1 {
		return messages
	}

	// Get image files from fileIDs
	contentBlocks := make([]services.ClaudeContentBlock, 0)

	// Add text content first
	textContent, ok := messages[lastUserIdx].Content.(string)
	if ok && textContent != "" {
		contentBlocks = append(contentBlocks, services.ClaudeContentBlock{
			Type: "text",
			Text: textContent,
		})
	}

	// Add image blocks
	for _, fileID := range fileIDs {
		file, err := bc.fileAnalysisRepo.FindByID(fileID)
		if err != nil {
			log.Printf("⚠️  File not found: %s", fileID)
			continue
		}

		// Check if it's an image
		if !bc.isImageFile(file.MimeType) {
			continue
		}

		// Read image file and encode to base64
		imageData, mediaType, err := bc.readImageFile(file.StoragePath, file.MimeType)
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

	// If we have multiple content blocks (text + images), use multimodal format
	if len(contentBlocks) > 1 {
		messages[lastUserIdx].Content = contentBlocks
		log.Printf("✓ Converted to multimodal message: %d blocks (%d images)", len(contentBlocks), len(contentBlocks)-1)
	}

	return messages
}

// readImageFile reads an image file and returns base64 encoded data
func (bc *BedrockController) readImageFile(filePath string, mimeType string) (string, string, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to read file: %w", err)
	}

	// Encode to base64
	encoded := base64.StdEncoding.EncodeToString(data)

	return encoded, mimeType, nil
}
