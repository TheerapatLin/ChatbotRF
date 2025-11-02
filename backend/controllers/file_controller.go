package controllers

import (
	"chatbot/models"
	"chatbot/repositories"
	"chatbot/services"
	"chatbot/utils"
	"fmt"
	"log"
	"mime/multipart"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
)

// FileController handles file analysis HTTP requests
type FileController struct {
	fileService *services.FileService
	repository  *repositories.FileAnalysisRepository
	messageRepo *repositories.MessageRepository
}

// NewFileController creates a new file controller
func NewFileController(
	fileService *services.FileService,
	repository *repositories.FileAnalysisRepository,
	messageRepo *repositories.MessageRepository,
) *FileController {
	return &FileController{
		fileService: fileService,
		repository:  repository,
		messageRepo: messageRepo,
	}
}

// AnalyzeFileRequest represents the form data for file analysis
type AnalyzeFileRequest struct {
	AnalysisType string `form:"analysis_type"` // summary, detail, qa, extract
	Prompt       string `form:"prompt"`
	Language     string `form:"language"` // th, en
}

// AnalyzeFile handles POST /api/file/analyze endpoint
func (ctrl *FileController) AnalyzeFile(c *fiber.Ctx) error {
	// 1. Parse and validate file upload
	file, analysisType, prompt, language, sessionID, systemPrompt, useHistory, err := ctrl.parseFileRequest(c)
	if err != nil {
		return err
	}

	// 2. Handle image files with Vision API
	contentType := file.Header.Get("Content-Type")
	if strings.Contains(contentType, "image/") {
		return ctrl.analyzeImageFile(c, file, contentType, prompt, language, systemPrompt)
	}

	// 3. Analyze text-based file
	req := services.FileAnalysisRequest{
		File:         file,
		AnalysisType: analysisType,
		Prompt:       prompt,
		Language:     language,
		SystemPrompt: systemPrompt,
		SessionID:    sessionID,
		UseHistory:   useHistory,
	}

	result, err := ctrl.fileService.AnalyzeFile(c.Context(), req)
	if err != nil {
		return ctrl.handleAnalysisError(c, err)
	}

	// 4. Save to database with session_id
	if err := ctrl.saveFileAnalysis(result, analysisType, prompt, language, sessionID); err != nil {
		log.Printf("⚠️  Failed to save file analysis to database: %v", err)
	}

	// 5. Save messages to database (like HandleChat does)
	if sessionID != "" {
		if err := ctrl.saveFileAnalysisMessages(file.Filename, prompt, result, sessionID); err != nil {
			log.Printf("⚠️  Failed to save file analysis messages: %v", err)
		}
	}

	// 6. Build ChatResponse-style response
	response := fiber.Map{
		"message_id":  result.FileID,
		"session_id":  sessionID,
		"reply":       result.Analysis,
		"tokens_used": result.TokensUsed,
		"model":       "gpt-4o-mini",
		"timestamp":   time.Now(),
		"file_info": fiber.Map{
			"file_id":   result.FileID,
			"filename":  result.FileName,
			"file_type": result.FileType,
			"file_size": result.FileSize,
		},
	}

	return utils.SuccessJSON(c, response)
}

// parseFileRequest parses and validates file upload request
func (ctrl *FileController) parseFileRequest(c *fiber.Ctx) (*multipart.FileHeader, string, string, string, string, string, bool, error) {
	// Get file from multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return nil, "", "", "", "", "", false, utils.BadRequest(c, "file is required")
	}

	// Get session_id (optional)
	sessionID := c.FormValue("session_id")

	// Get analysis type with default
	analysisType := c.FormValue("analysis_type")
	if analysisType == "" {
		analysisType = "summary"
	}

	// Validate analysis type
	validTypes := []string{"summary", "detail", "qa", "extract"}
	isValid := false
	for _, vt := range validTypes {
		if analysisType == vt {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, "", "", "", "", "", false, utils.BadRequest(c, "invalid analysis_type. Allowed: summary, detail, qa, extract")
	}

	// Get optional parameters
	prompt := c.FormValue("prompt")
	language := c.FormValue("language")
	if language == "" {
		language = "th"
	}

	// Get optional system_prompt
	systemPrompt := c.FormValue("system_prompt")

	// Get optional use_history (default: false)
	useHistory := c.FormValue("use_history") == "true"

	return file, analysisType, prompt, language, sessionID, systemPrompt, useHistory, nil
}

// analyzeImageFile handles image file analysis using Vision API
func (ctrl *FileController) analyzeImageFile(c *fiber.Ctx, file *multipart.FileHeader, contentType, prompt, language, systemPrompt string) error {
	// Open file
	fileData, err := file.Open()
	if err != nil {
		return utils.InternalError(c, "failed to open file")
	}
	defer fileData.Close()

	// Get session_id from form
	sessionID := c.FormValue("session_id")

	// Analyze image (pass systemPrompt to service if needed)
	analysis, err := ctrl.fileService.AnalyzeImage(c.Context(), fileData, file.Filename, prompt, systemPrompt)
	if err != nil {
		return utils.InternalError(c, fmt.Sprintf("failed to analyze image: %v", err))
	}

	// Build internal response for saving
	fileID := generateFileID()

	// Save to database (optional for images)
	result := &services.FileAnalysisResponse{
		FileID:      fileID,
		FileName:    file.Filename,
		FileType:    contentType,
		FileSize:    file.Size,
		Analysis:    analysis,
		KeyPoints:   []string{},
		Language:    language,
		TokensUsed:  0,
		ProcessTime: 0,
	}

	if err := ctrl.saveFileAnalysis(result, "image_analysis", prompt, language, sessionID); err != nil {
		log.Printf("⚠️  Failed to save image analysis to database: %v", err)
	}

	// Save messages to database (like HandleChat does)
	if sessionID != "" {
		if err := ctrl.saveFileAnalysisMessages(file.Filename, prompt, result, sessionID); err != nil {
			log.Printf("⚠️  Failed to save image analysis messages: %v", err)
		}
	}

	// Build ChatResponse-style response
	response := fiber.Map{
		"message_id":  fileID,
		"session_id":  sessionID,
		"reply":       analysis,
		"tokens_used": 0,
		"model":       "gpt-4o-mini",
		"timestamp":   time.Now(),
		"file_info": fiber.Map{
			"file_id":   fileID,
			"filename":  file.Filename,
			"file_type": contentType,
			"file_size": file.Size,
		},
	}

	return utils.SuccessJSON(c, response)
}

// handleAnalysisError maps service errors to appropriate HTTP status codes
func (ctrl *FileController) handleAnalysisError(c *fiber.Ctx, err error) error {
	errMsg := err.Error()

	if strings.Contains(errMsg, "unsupported file type") {
		return c.Status(fiber.StatusUnsupportedMediaType).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if strings.Contains(errMsg, "file size exceeds") {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if strings.Contains(errMsg, "failed to parse") || strings.Contains(errMsg, "failed to extract") {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": fmt.Sprintf("failed to parse file: %v", err),
		})
	}

	return utils.InternalError(c, fmt.Sprintf("failed to analyze file: %v", err))
}

// saveFileAnalysis saves analysis result to database
func (ctrl *FileController) saveFileAnalysis(result *services.FileAnalysisResponse, analysisType, prompt, language, sessionID string) error {
	if ctrl.repository == nil {
		return nil
	}

	fileAnalysis := &models.FileAnalysis{
		SessionID:     sessionID, // NEW - Link to conversation session
		FileName:      sanitizeUTF8(result.FileName),
		FileType:      result.FileType,
		FileSize:      result.FileSize,
		AnalysisType:  analysisType,
		CustomPrompt:  sanitizeUTF8(prompt),
		Language:      language,
		Analysis:      sanitizeUTF8(result.Analysis),
		KeyPoints:     pq.StringArray(sanitizeStringArray(result.KeyPoints)),
		Entities:      pq.StringArray(sanitizeStringArray(result.Entities)),
		Sentiment:     result.Sentiment,
		TokensUsed:    result.TokensUsed,
		ProcessTimeMs: result.ProcessTime,
	}

	if err := ctrl.repository.Create(fileAnalysis); err != nil {
		return err
	}

	result.FileID = fileAnalysis.ID.String()
	if sessionID != "" {
		log.Printf("✅ File analysis saved to database with ID: %s (session: %s)", result.FileID, sessionID)
	} else {
		log.Printf("✅ File analysis saved to database with ID: %s", result.FileID)
	}
	return nil
}

// GetFileHistory handles GET /api/file/history endpoint
func (ctrl *FileController) GetFileHistory(c *fiber.Ctx) error {
	// Validate pagination
	limit, offset := utils.ValidatePagination(
		c.QueryInt("limit", 20),
		c.QueryInt("offset", 0),
		100, // max limit
	)

	fileType := c.Query("file_type", "all")

	// Get files from repository
	files, total, err := ctrl.getFilesByType(fileType, limit, offset)
	if err != nil {
		log.Printf("Failed to fetch file history: %v", err)
		return utils.InternalError(c, "failed to fetch file history")
	}

	// Build and return response
	response := ctrl.buildFileHistoryResponse(files, total, limit, offset)
	return utils.SuccessJSON(c, response)
}

// getFilesByType retrieves files based on file type filter
func (ctrl *FileController) getFilesByType(fileType string, limit, offset int) ([]models.FileAnalysis, int64, error) {
	if fileType != "" && fileType != "all" {
		return ctrl.repository.GetByFileType(fileType, limit, offset)
	}
	return ctrl.repository.GetAll(limit, offset)
}

// buildFileHistoryResponse builds the file history response
func (ctrl *FileController) buildFileHistoryResponse(files []models.FileAnalysis, total int64, limit, offset int) fiber.Map {
	response := make([]fiber.Map, len(files))
	for i, file := range files {
		response[i] = fiber.Map{
			"file_id":       file.ID.String(),
			"filename":      file.FileName,
			"file_type":     file.FileType,
			"file_size":     file.FileSize,
			"analysis_type": file.AnalysisType,
			"language":      file.Language,
			"tokens_used":   file.TokensUsed,
			"created_at":    file.CreatedAt,
		}
	}

	return fiber.Map{
		"files":  response,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	}
}

// ReanalyzeFile handles POST /api/file/:file_id/reanalyze endpoint
// This is a placeholder for future implementation
func (ctrl *FileController) ReanalyzeFile(c *fiber.Ctx) error {
	fileID := c.Params("file_id")
	if fileID == "" {
		return utils.BadRequest(c, "file_id is required")
	}

	// TODO: Implement re-analysis from stored files
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "re-analysis feature not yet implemented",
	})
}

// Helper functions

func generateFileID() string {
	// Simple file ID generator using timestamp
	// In production, use UUID library or database auto-increment
	return fmt.Sprintf("file_%d", time.Now().UnixNano())
}

// sanitizeUTF8 removes invalid UTF-8 sequences from a string
func sanitizeUTF8(s string) string {
	if utf8.ValidString(s) {
		return s
	}

	// Build a new string with only valid UTF-8 runes
	var builder strings.Builder
	builder.Grow(len(s))

	for _, r := range s {
		if r != utf8.RuneError {
			builder.WriteRune(r)
		}
	}

	return builder.String()
}

// sanitizeStringArray sanitizes all strings in an array
func sanitizeStringArray(arr []string) []string {
	if arr == nil {
		return nil
	}

	result := make([]string, len(arr))
	for i, s := range arr {
		result[i] = sanitizeUTF8(s)
	}
	return result
}

// saveFileAnalysisMessages saves user message and AI response to messages table
// Similar to ChatController.saveMessages
func (ctrl *FileController) saveFileAnalysisMessages(filename, prompt string, result *services.FileAnalysisResponse, sessionID string) error {
	if ctrl.messageRepo == nil {
		return fmt.Errorf("message repository not initialized")
	}

	// 1. Save user message (file upload request)
	content := fmt.Sprintf("อัปโหลดไฟล์: %s", filename)
	if prompt != "" {
		content = fmt.Sprintf("อัปโหลดไฟล์: %s\nคำสั่ง: %s", filename, prompt)
	}

	userMessage := &models.Message{
		SessionID: sessionID,
		Role:      models.RoleUser,
		Content:   content,
	}

	// Add file attachment info
	attachment := models.FileAttachment{
		FileID:   result.FileID,
		Filename: result.FileName,
		FileType: result.FileType,
		FileSize: result.FileSize,
	}

	if err := userMessage.SetFileAttachments([]models.FileAttachment{attachment}); err != nil {
		log.Printf("⚠️  Warning: Failed to set file attachments: %v", err)
	}

	if err := ctrl.messageRepo.Create(userMessage); err != nil {
		return fmt.Errorf("failed to save user message: %w", err)
	}

	// 2. Save AI response (analysis result)
	tokensUsed := result.TokensUsed
	assistantMessage := &models.Message{
		SessionID:  sessionID,
		Role:       models.RoleAssistant,
		Content:    result.Analysis,
		TokensUsed: &tokensUsed,
	}

	if err := ctrl.messageRepo.Create(assistantMessage); err != nil {
		return fmt.Errorf("failed to save assistant message: %w", err)
	}

	log.Printf("✅ Saved file analysis messages to conversation (session: %s, file: %s)", sessionID, filename)
	return nil
}
