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
	messageRepo      *repositories.MessageRepository
	personaRepo      *repositories.PersonaRepository
	fileAnalysisRepo *repositories.FileAnalysisRepository
	openaiService    *services.OpenAIService
	contextService   *services.ContextService
}

// NewChatController creates a new chat controller
func NewChatController(
	messageRepo *repositories.MessageRepository,
	personaRepo *repositories.PersonaRepository,
	fileAnalysisRepo *repositories.FileAnalysisRepository,
	openaiService *services.OpenAIService,
) *ChatController {
	return &ChatController{
		messageRepo:      messageRepo,
		personaRepo:      personaRepo,
		fileAnalysisRepo: fileAnalysisRepo,
		openaiService:    openaiService,
		contextService:   services.NewContextService(messageRepo, fileAnalysisRepo),
	}
}

// ChatRequest represents the incoming chat request
type ChatRequest struct {
	Message      string   `json:"message" validate:"required"`
	SessionID    string   `json:"session_id,omitempty"`
	PersonaID    *int     `json:"persona_id,omitempty"`
	SystemPrompt string   `json:"system_prompt,omitempty"`
	Temperature  float32  `json:"temperature,omitempty"`
	MaxTokens    int      `json:"max_tokens,omitempty"`
	Model        string   `json:"model,omitempty"`
	UseHistory   bool     `json:"use_history,omitempty"`
	FileIDs      []string `json:"file_ids,omitempty"` // File IDs for current message only
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
	MessageID    string       `json:"message_id"`
	SessionID    string       `json:"session_id"`
	Reply        string       `json:"reply"`
	PersonaUsed  *PersonaInfo `json:"persona,omitempty"`
	TokensUsed   int          `json:"tokens_used"`
	Model        string       `json:"model"`
	Timestamp    time.Time    `json:"timestamp"`
	HistoryUsed  bool         `json:"history_used"`
	HistoryCount int          `json:"history_count"`
}

// MessageHistoryItem represents a message in history
type MessageHistoryItem struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	PersonaID *int      `json:"persona_id,omitempty"`
	SessionID string    `json:"session_id,omitempty"`
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
	req, err := ctrl.parseRequest(c)
	if err != nil {
		return err
	}

	// 2. Get persona and system prompt
	systemPrompt, personaInfo, err := ctrl.getPersonaInfo(req)
	if err != nil {
		return err
	}

	// 3. Generate or use session ID
	sessionID := ctrl.getOrGenerateSessionID(req)

	// 4. Build context with conversation history
	messages, historyCount := ctrl.buildMessages(req, sessionID, systemPrompt)

	// 5. Call OpenAI service
	openaiResp, err := ctrl.callOpenAI(req, messages, systemPrompt)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to get AI response: %v", err),
		})
	}

	// 6. Save messages to database
	if err := ctrl.saveMessages(req, sessionID, openaiResp); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save messages",
		})
	}

	// 7. Build and return response
	response := ctrl.buildResponse(sessionID, openaiResp, personaInfo, req.UseHistory, historyCount)
	return c.Status(fiber.StatusOK).JSON(response)
}

// parseRequest parses and validates the incoming request
func (ctrl *ChatController) parseRequest(c *fiber.Ctx) (*ChatRequest, error) {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Message == "" {
		return nil, c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Message is required",
		})
	}

	return &req, nil
}

// getPersonaInfo retrieves persona information and determines system prompt
func (ctrl *ChatController) getPersonaInfo(req *ChatRequest) (string, *PersonaInfo, error) {
	if req.PersonaID == nil {
		return req.SystemPrompt, nil, nil
	}

	persona, err := ctrl.personaRepo.FindByID(*req.PersonaID)
	if err != nil {
		return "", nil, fmt.Errorf("persona with ID %d not found", *req.PersonaID)
	}

	// Start with persona's system prompt
	systemPrompt := persona.SystemPrompt

	// Append custom system_prompt if provided (don't replace!)
	if req.SystemPrompt != "" {
		systemPrompt = systemPrompt + "\n\n--- Additional Instructions ---\n" + req.SystemPrompt
	}

	personaInfo := &PersonaInfo{
		ID:          persona.ID,
		Name:        persona.Name,
		Expertise:   persona.Expertise,
		Icon:        persona.Icon,
		Description: persona.Description,
	}

	return systemPrompt, personaInfo, nil
}

// getOrGenerateSessionID returns existing session ID or generates a new one
func (ctrl *ChatController) getOrGenerateSessionID(req *ChatRequest) string {
	if req.SessionID != "" || !req.UseHistory {
		return req.SessionID
	}
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// buildMessages builds OpenAI messages array with optional history and current files
func (ctrl *ChatController) buildMessages(req *ChatRequest, sessionID, systemPrompt string) ([]openai.ChatCompletionMessage, int) {
	historyCount := 0

	// Build context with history if enabled
	if req.UseHistory && sessionID != "" {
		messages, err := ctrl.contextService.BuildContextWithHistory(
			sessionID, systemPrompt, req.Message, 10,
		)
		if err != nil {
			fmt.Printf("⚠️  Failed to build context with history: %v\n", err)
		} else {
			historyCount = len(messages) - 1 // Subtract current message
			if systemPrompt != "" {
				historyCount--
			}

			// Add file context for current message if files provided
			if len(req.FileIDs) > 0 {
				fileContext, err := ctrl.contextService.BuildFileContext(req.FileIDs)
				if err == nil && fileContext != "" {
					// Insert file context before the last user message
					messages = append(
						messages[:len(messages)-1],
						openai.ChatCompletionMessage{
							Role:    openai.ChatMessageRoleSystem,
							Content: fileContext,
						},
						messages[len(messages)-1],
					)
				}
			}

			return messages, historyCount
		}
	}

	// Fallback: simple message with optional file context
	messages := []openai.ChatCompletionMessage{}

	if systemPrompt != "" {
		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleSystem,
			Content: systemPrompt,
		})
	}

	// Add file context if files provided
	if len(req.FileIDs) > 0 {
		fileContext, err := ctrl.contextService.BuildFileContext(req.FileIDs)
		if err == nil && fileContext != "" {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleSystem,
				Content: fileContext,
			})
		}
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: req.Message,
	})

	return messages, 0
}

// callOpenAI sends request to OpenAI API
func (ctrl *ChatController) callOpenAI(req *ChatRequest, messages []openai.ChatCompletionMessage, systemPrompt string) (*services.ChatResponse, error) {
	// Only include system prompt if not using history (to avoid duplication)
	systemPromptForAPI := ""
	if !req.UseHistory || len(messages) == 1 {
		systemPromptForAPI = systemPrompt
	}

	// Debug log
	fmt.Printf("🔍 Debug callOpenAI:\n")
	fmt.Printf("  - UseHistory: %v\n", req.UseHistory)
	fmt.Printf("  - Messages count: %d\n", len(messages))
	fmt.Printf("  - SystemPrompt from param: %q\n", systemPrompt)
	fmt.Printf("  - SystemPrompt for API: %q\n", systemPromptForAPI)

	openaiReq := services.ChatRequest{
		Messages:     messages,
		Model:        req.Model,
		Temperature:  req.Temperature,
		MaxTokens:    req.MaxTokens,
		SystemPrompt: systemPromptForAPI,
	}

	return ctrl.openaiService.SendChatRequest(openaiReq)
}

// saveMessages saves user message and AI response to database
func (ctrl *ChatController) saveMessages(req *ChatRequest, sessionID string, openaiResp *services.ChatResponse) error {
	// Save user message
	userMessage := &models.Message{
		SessionID: sessionID,
		Role:      models.RoleUser,
		Content:   req.Message,
		PersonaID: req.PersonaID,
	}

	if err := ctrl.messageRepo.Create(userMessage); err != nil {
		return err
	}

	// Save AI response
	tokensUsed := openaiResp.TokensUsed
	assistantMessage := &models.Message{
		SessionID:  sessionID,
		Role:       models.RoleAssistant,
		Content:    openaiResp.Content,
		PersonaID:  req.PersonaID,
		TokensUsed: &tokensUsed,
	}
	return ctrl.messageRepo.Create(assistantMessage)
}

// buildResponse builds the final chat response
func (ctrl *ChatController) buildResponse(sessionID string, openaiResp *services.ChatResponse, personaInfo *PersonaInfo, useHistory bool, historyCount int) ChatResponse {
	return ChatResponse{
		MessageID:    "", // Will be set after save
		SessionID:    sessionID,
		Reply:        openaiResp.Content,
		PersonaUsed:  personaInfo,
		TokensUsed:   openaiResp.TokensUsed,
		Model:        openaiResp.Model,
		Timestamp:    time.Now(),
		HistoryUsed:  useHistory && historyCount > 0,
		HistoryCount: historyCount,
	}
}

// GetChatHistory handles GET /api/chat/history endpoint
func (ctrl *ChatController) GetChatHistory(c *fiber.Ctx) error {
	// Parse and validate pagination
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// Validate pagination
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Get messages from repository
	messages, total, err := ctrl.messageRepo.FindWithPagination(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve chat history",
		})
	}

	// Convert messages to response items
	items := make([]MessageHistoryItem, len(messages))
	for i, msg := range messages {
		items[i] = MessageHistoryItem{
			ID:        msg.ID.String(),
			Role:      string(msg.Role),
			Content:   msg.Content,
			PersonaID: msg.PersonaID,
			SessionID: msg.SessionID,
			CreatedAt: msg.CreatedAt,
		}
	}

	// Build response
	response := ChatHistoryResponse{
		Messages: items,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// DeleteAllMessages handles DELETE /api/chat endpoint
func (ctrl *ChatController) DeleteAllMessages(c *fiber.Ctx) error {
	// Delete all messages from database
	if err := ctrl.messageRepo.DeleteAll(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete messages",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All messages deleted successfully",
	})
}

// GetChatHistoryBySession handles GET /api/chats/session/:sessionId endpoint
func (ctrl *ChatController) GetChatHistoryBySession(c *fiber.Ctx) error {
	// Get session ID from URL parameter
	sessionID := c.Params("sessionId")
	if sessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}

	// Parse and validate pagination
	limit := c.QueryInt("limit", 50)
	offset := c.QueryInt("offset", 0)

	// Validate pagination
	if limit <= 0 {
		limit = 50
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	// Get messages by session ID from repository
	messages, err := ctrl.messageRepo.GetAllBySession(sessionID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve chat history for session",
		})
	}

	// Apply pagination manually
	total := int64(len(messages))
	start := offset
	end := offset + limit

	if start > len(messages) {
		start = len(messages)
	}
	if end > len(messages) {
		end = len(messages)
	}

	paginatedMessages := messages[start:end]

	// Convert messages to response items
	items := make([]MessageHistoryItem, len(paginatedMessages))
	for i, msg := range paginatedMessages {
		items[i] = MessageHistoryItem{
			ID:        msg.ID.String(),
			Role:      string(msg.Role),
			Content:   msg.Content,
			PersonaID: msg.PersonaID,
			CreatedAt: msg.CreatedAt,
		}
	}

	// Build response
	response := ChatHistoryResponse{
		Messages: items,
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
