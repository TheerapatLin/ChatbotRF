package controllers

import (
	"fmt"
	"time"

	"chatbot/models"
	"chatbot/repositories"
	"chatbot/services"

	"github.com/gofiber/fiber/v2"
	"github.com/sashabaranov/go-openai"
)

// ChatController handles chat-related HTTP requests
type ChatController struct {
	messageRepo    *repositories.MessageRepository
	personaRepo    *repositories.PersonaRepository
	openaiService  *services.OpenAIService
	contextService *services.ContextService
}

// NewChatController creates a new chat controller
func NewChatController(
	messageRepo *repositories.MessageRepository,
	personaRepo *repositories.PersonaRepository,
	openaiService *services.OpenAIService,
) *ChatController {
	return &ChatController{
		messageRepo:    messageRepo,
		personaRepo:    personaRepo,
		openaiService:  openaiService,
		contextService: services.NewContextService(messageRepo),
	}
}

// ChatRequest represents the incoming chat request
type ChatRequest struct {
	Message      string  `json:"message" validate:"required"`
	SessionID    string  `json:"session_id,omitempty"`    // Session ID for conversation history
	PersonaID    *int    `json:"persona_id,omitempty"`
	SystemPrompt string  `json:"system_prompt,omitempty"`
	Temperature  float32 `json:"temperature,omitempty"`
	MaxTokens    int     `json:"max_tokens,omitempty"`
	Model        string  `json:"model,omitempty"`
	UseHistory   bool    `json:"use_history,omitempty"`   // Enable conversation history (default: false)
}

// PersonaInfo contains persona information in response
type PersonaInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Expertise   string `json:"expertise"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

// ChatResponse represents the API response
type ChatResponse struct {
	MessageID       string       `json:"message_id"`
	SessionID       string       `json:"session_id"`
	Reply           string       `json:"reply"`
	PersonaUsed     *PersonaInfo `json:"persona,omitempty"`
	TokensUsed      int          `json:"tokens_used"`
	Model           string       `json:"model"`
	Timestamp       time.Time    `json:"timestamp"`
	HistoryUsed     bool         `json:"history_used"`     // Indicates if history was used
	HistoryCount    int          `json:"history_count"`    // Number of historical messages used
}

// MessageHistoryItem represents a single message in history
type MessageHistoryItem struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	PersonaID *int      `json:"persona_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ChatHistoryResponse represents the chat history API response
type ChatHistoryResponse struct {
	Messages []MessageHistoryItem `json:"messages"`
	Total    int64                `json:"total"`
	Limit    int                  `json:"limit"`
	Offset   int                  `json:"offset"`
}

// HandleChat handles POST /api/chat endpoint
func (ctrl *ChatController) HandleChat(c *fiber.Ctx) error {
	// 1. Parse and validate request
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate message is not empty
	if req.Message == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Message is required",
		})
	}

	// 2. Get persona if persona_id provided
	var persona *models.Persona
	var systemPrompt string
	var personaInfo *PersonaInfo

	if req.PersonaID != nil {
		var err error
		persona, err = ctrl.personaRepo.FindByID(*req.PersonaID)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": fmt.Sprintf("Persona with ID %d not found", *req.PersonaID),
			})
		}

		// Use persona's system prompt if no custom system prompt provided
		if req.SystemPrompt == "" {
			systemPrompt = persona.SystemPrompt
		} else {
			systemPrompt = req.SystemPrompt
		}

		personaInfo = &PersonaInfo{
			ID:          persona.ID,
			Name:        persona.Name,
			Expertise:   persona.Expertise,
			Icon:        persona.Icon,
			Description: persona.Description,
		}
	} else {
		// Use custom system prompt if provided
		systemPrompt = req.SystemPrompt
	}

	// 3. Generate or use session ID
	sessionID := req.SessionID
	if sessionID == "" && req.UseHistory {
		// Generate new session ID if history is requested but no session provided
		sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
	}

	// 4. Build context with conversation history
	var messages []openai.ChatCompletionMessage
	historyCount := 0

	if req.UseHistory && sessionID != "" {
		// Use context service to build messages with history (10 messages limit)
		var err error
		messages, err = ctrl.contextService.BuildContextWithHistory(
			sessionID,
			systemPrompt,
			req.Message,
			10, // Retrieve last 10 messages
		)
		if err != nil {
			// Log error but continue without history
			fmt.Printf("⚠️  Failed to build context with history: %v\n", err)
		} else {
			// Count historical messages (exclude system prompt and current message)
			historyCount = len(messages) - 1 // Subtract current message
			if systemPrompt != "" {
				historyCount-- // Subtract system prompt
			}
		}
	}

	// If not using history or history failed, build simple message array
	if len(messages) == 0 {
		messages = []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: req.Message,
			},
		}
	}

	// 5. Call OpenAI service
	// Note: When using history, system prompt is already included in messages
	// Only send SystemPrompt if NOT using history
	systemPromptForAPI := ""
	if !req.UseHistory || len(messages) == 1 {
		// Not using history OR only has current message (history retrieval failed)
		systemPromptForAPI = systemPrompt
	}

	openaiReq := services.ChatRequest{
		Messages:     messages,
		Model:        req.Model,
		Temperature:  req.Temperature,
		MaxTokens:    req.MaxTokens,
		SystemPrompt: systemPromptForAPI,
	}

	openaiResp, err := ctrl.openaiService.SendChatRequest(openaiReq)
	if err != nil {
		return c.Status(503).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get AI response: %v", err),
		})
	}

	// 6. Save user message to database
	userMessage := &models.Message{
		SessionID: sessionID,
		Role:      models.RoleUser,
		Content:   req.Message,
		PersonaID: req.PersonaID,
	}

	if err := ctrl.messageRepo.Create(userMessage); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save user message",
		})
	}

	// 7. Save AI response to database
	tokensUsed := openaiResp.TokensUsed
	assistantMessage := &models.Message{
		SessionID:  sessionID,
		Role:       models.RoleAssistant,
		Content:    openaiResp.Content,
		PersonaID:  req.PersonaID,
		TokensUsed: &tokensUsed,
	}

	if err := ctrl.messageRepo.Create(assistantMessage); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to save assistant message",
		})
	}

	// 8. Return response
	response := ChatResponse{
		MessageID:    assistantMessage.ID.String(),
		SessionID:    sessionID,
		Reply:        openaiResp.Content,
		PersonaUsed:  personaInfo,
		TokensUsed:   openaiResp.TokensUsed,
		Model:        openaiResp.Model,
		Timestamp:    assistantMessage.CreatedAt,
		HistoryUsed:  req.UseHistory && historyCount > 0,
		HistoryCount: historyCount,
	}

	return c.Status(200).JSON(response)
}

// GetChatHistory handles GET /api/chat/history endpoint
func (ctrl *ChatController) GetChatHistory(c *fiber.Ctx) error {
	// Parse query parameters with defaults
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// Validate limit (max 100)
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 1
	}

	// Validate offset
	if offset < 0 {
		offset = 0
	}

	// Get messages from repository
	messages, total, err := ctrl.messageRepo.FindWithPagination(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve chat history",
		})
	}

	// Convert to response format
	historyItems := make([]MessageHistoryItem, len(messages))
	for i, msg := range messages {
		historyItems[i] = MessageHistoryItem{
			ID:        msg.ID.String(),
			Role:      string(msg.Role),
			Content:   msg.Content,
			PersonaID: msg.PersonaID,
			CreatedAt: msg.CreatedAt,
		}
	}

	// Return response
	response := ChatHistoryResponse{
		Messages: historyItems,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	return c.Status(200).JSON(response)
}
