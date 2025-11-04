package controllers

import (
	"encoding/json"
	"fmt"
	"log"
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

	// Get file context if files provided
	if len(req.FileIDs) > 0 {
		var err error
		fileContext, err = bc.contextService.BuildFileContext(req.FileIDs)
		if err != nil {
			log.Printf("⚠️ Failed to build file context: %v", err)
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
