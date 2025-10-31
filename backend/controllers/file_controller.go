package controllers

import (
	"chatbot/models"
	"chatbot/repositories"
	"chatbot/services"
	"fmt"
	"log"
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
}

// NewFileController creates a new file controller
func NewFileController(
	fileService *services.FileService,
	repository *repositories.FileAnalysisRepository,
) *FileController {
	return &FileController{
		fileService: fileService,
		repository:  repository,
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
	// Get file from multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file is required",
		})
	}

	// Get optional parameters
	analysisType := c.FormValue("analysis_type")
	if analysisType == "" {
		analysisType = "summary" // default
	}

	// Validate analysis type
	validTypes := []string{"summary", "detail", "qa", "extract"}
	isValidType := false
	for _, vt := range validTypes {
		if analysisType == vt {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid analysis_type. Allowed: summary, detail, qa, extract",
		})
	}

	prompt := c.FormValue("prompt")
	language := c.FormValue("language")
	if language == "" {
		language = "th" // default to Thai
	}

	// Build request
	req := services.FileAnalysisRequest{
		File:         file,
		AnalysisType: analysisType,
		Prompt:       prompt,
		Language:     language,
	}

	// Handle image files separately using Vision API
	contentType := file.Header.Get("Content-Type")
	if strings.Contains(contentType, "image/") {
		// Open file
		fileData, err := file.Open()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to open file",
			})
		}
		defer fileData.Close()

		// Analyze image
		analysis, err := ctrl.fileService.AnalyzeImage(c.Context(), fileData, file.Filename, prompt)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to analyze image: " + err.Error(),
			})
		}

		// Build response for image
		response := services.FileAnalysisResponse{
			FileID:      generateFileID(),
			FileName:    file.Filename,
			FileType:    contentType,
			FileSize:    file.Size,
			Analysis:    analysis,
			Summary:     extractFirstSentence(analysis),
			KeyPoints:   []string{},
			Language:    language,
			TokensUsed:  0, // Not available for Vision API in this implementation
			ProcessTime: 0,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}

	// Analyze file with regular text analysis
	result, err := ctrl.fileService.AnalyzeFile(c.Context(), req)
	if err != nil {
		// Check error type for appropriate status code
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
				"error": "failed to parse file: " + err.Error(),
			})
		}

		// Generic error
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to analyze file: " + err.Error(),
		})
	}

	// Save to database
	fileAnalysis := &models.FileAnalysis{
		FileName:      sanitizeUTF8(result.FileName),
		FileType:      result.FileType,
		FileSize:      result.FileSize,
		AnalysisType:  analysisType,
		CustomPrompt:  sanitizeUTF8(prompt),
		Language:      language,
		Analysis:      sanitizeUTF8(result.Analysis),
		Summary:       sanitizeUTF8(result.Summary),
		KeyPoints:     pq.StringArray(sanitizeStringArray(result.KeyPoints)),
		Entities:      pq.StringArray(sanitizeStringArray(result.Entities)),
		Sentiment:     result.Sentiment,
		TokensUsed:    result.TokensUsed,
		ProcessTimeMs: result.ProcessTime,
	}

	// Save to database (log error but don't fail the request)
	if ctrl.repository != nil {
		if err := ctrl.repository.Create(fileAnalysis); err != nil {
			log.Printf("⚠️  Failed to save file analysis to database: %v", err)
		} else {
			// Update response with database ID
			result.FileID = fileAnalysis.ID.String()
			log.Printf("✅ File analysis saved to database with ID: %s", result.FileID)
		}
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// GetFileHistory handles GET /api/file/history endpoint
func (ctrl *FileController) GetFileHistory(c *fiber.Ctx) error {
	// Parse query parameters
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)
	fileType := c.Query("file_type", "all")

	// Validate limit (max 100)
	if limit > 100 {
		limit = 100
	}
	if limit < 1 {
		limit = 20
	}

	// Validate offset
	if offset < 0 {
		offset = 0
	}

	var files []models.FileAnalysis
	var total int64
	var err error

	// Get files based on filter
	if fileType != "" && fileType != "all" {
		files, total, err = ctrl.repository.GetByFileType(fileType, limit, offset)
	} else {
		files, total, err = ctrl.repository.GetAll(limit, offset)
	}

	if err != nil {
		log.Printf("Failed to fetch file history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch file history",
		})
	}

	// Build response
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"files":  response,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// ReanalyzeFile handles POST /api/file/:file_id/reanalyze endpoint
// This is a placeholder for future implementation
func (ctrl *FileController) ReanalyzeFile(c *fiber.Ctx) error {
	fileID := c.Params("file_id")

	if fileID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file_id is required",
		})
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

func extractFirstSentence(text string) string {
	// Extract first sentence for summary
	sentences := strings.Split(text, ".")
	if len(sentences) > 0 && len(sentences[0]) > 10 {
		return strings.TrimSpace(sentences[0]) + "."
	}
	if len(text) > 200 {
		return text[:200] + "..."
	}
	return text
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
