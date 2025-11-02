package controllers

import (
	"fmt"
	"time"

	"chatbot/models"
	"chatbot/repositories"
	"chatbot/services"
	"chatbot/utils"

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
	FileIDs      []string `json:"file_ids,omitempty"` // Array of file analysis UUIDs to attach
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

// ChatHistoryResponse represents the chat history API response
type ChatHistoryResponse struct {
	Messages []utils.MessageHistoryItem `json:"messages"`
	Total    int64                      `json:"total"`
	Limit    int                        `json:"limit"`
	Offset   int                        `json:"offset"`
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
		return utils.ServiceUnavailable(c, fmt.Sprintf("Failed to get AI response: %v", err))
	}

	// 6. Save messages to database
	if err := ctrl.saveMessages(req, sessionID, openaiResp); err != nil {
		return utils.InternalError(c, "Failed to save messages")
	}

	// 7. Build and return response
	response := ctrl.buildResponse(sessionID, openaiResp, personaInfo, req.UseHistory, historyCount)
	return utils.SuccessJSON(c, response)
}

// parseRequest parses and validates the incoming request
func (ctrl *ChatController) parseRequest(c *fiber.Ctx) (*ChatRequest, error) {
	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return nil, utils.BadRequest(c, "Invalid request body")
	}

	if req.Message == "" {
		return nil, utils.BadRequest(c, "Message is required")
	}

	// Validate file_ids limit (max 5 files per message)
	if len(req.FileIDs) > 5 {
		return nil, utils.BadRequest(c, "Maximum 5 files allowed per message")
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

	systemPrompt := req.SystemPrompt
	if systemPrompt == "" {
		systemPrompt = persona.SystemPrompt
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

// buildMessages builds OpenAI messages array with optional history and files
func (ctrl *ChatController) buildMessages(req *ChatRequest, sessionID, systemPrompt string) ([]openai.ChatCompletionMessage, int) {
	historyCount := 0

	// Auto-load files from session if not explicitly provided
	fileIDs := req.FileIDs
	if len(fileIDs) == 0 && sessionID != "" {
		// Try to load recent files from this session
		sessionFiles, err := ctrl.fileAnalysisRepo.GetRecentBySessionID(sessionID, 5)
		if err == nil && len(sessionFiles) > 0 {
			fileIDs = make([]string, len(sessionFiles))
			for i, f := range sessionFiles {
				fileIDs[i] = f.ID.String()
			}
			fmt.Printf("🔗 Auto-loaded %d files from session %s\n", len(fileIDs), sessionID)
		}
	}

	// Build context with history and files if enabled
	if req.UseHistory && sessionID != "" {
		messages, err := ctrl.contextService.BuildContextWithFiles(
			sessionID, systemPrompt, req.Message, fileIDs, 10,
		)
		if err != nil {
			fmt.Printf("⚠️  Failed to build context with history: %v\n", err)
		} else {
			historyCount = len(messages) - 1 // Subtract current message
			if systemPrompt != "" {
				historyCount--
			}
			// Also subtract file context message if files are provided
			if len(fileIDs) > 0 {
				historyCount--
			}
			return messages, historyCount
		}
	}

	// Fallback: simple message without history
	return []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleUser, Content: req.Message},
	}, 0
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

	// Add file attachments if provided
	if len(req.FileIDs) > 0 {
		attachments, err := ctrl.buildFileAttachments(req.FileIDs)
		if err != nil {
			fmt.Printf("⚠️  Warning: Failed to build file attachments: %v\n", err)
		} else if len(attachments) > 0 {
			if err := userMessage.SetFileAttachments(attachments); err != nil {
				fmt.Printf("⚠️  Warning: Failed to set file attachments: %v\n", err)
			}
		}
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

// buildFileAttachments fetches file analysis and builds FileAttachment array
func (ctrl *ChatController) buildFileAttachments(fileIDs []string) ([]models.FileAttachment, error) {
	attachments := make([]models.FileAttachment, 0, len(fileIDs))

	for _, fileID := range fileIDs {
		analysis, err := ctrl.fileAnalysisRepo.FindByID(fileID)
		if err != nil {
			fmt.Printf("⚠️  Warning: File not found for ID %s: %v\n", fileID, err)
			continue
		}

		attachments = append(attachments, models.FileAttachment{
			FileID:   fileID,
			Filename: analysis.FileName,
			FileType: analysis.FileType,
			FileSize: analysis.FileSize,
		})
	}

	return attachments, nil
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
	limit, offset := utils.ValidatePagination(
		c.QueryInt("limit", 50),
		c.QueryInt("offset", 0),
		100, // max limit
	)

	// Get messages from repository
	messages, total, err := ctrl.messageRepo.FindWithPagination(limit, offset)
	if err != nil {
		return utils.InternalError(c, "Failed to retrieve chat history")
	}

	// Build response
	response := ChatHistoryResponse{
		Messages: utils.ToMessageHistoryItems(messages),
		Total:    total,
		Limit:    limit,
		Offset:   offset,
	}

	return utils.SuccessJSON(c, response)
}