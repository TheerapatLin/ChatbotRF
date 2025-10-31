package controllers

import (
	"chatbot/services"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// FileController handles file analysis HTTP requests
type FileController struct {
	fileService *services.FileService
}

// NewFileController creates a new file controller
func NewFileController(fileService *services.FileService) *FileController {
	return &FileController{
		fileService: fileService,
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

	return c.Status(fiber.StatusOK).JSON(result)
}

// GetFileHistory handles GET /api/file/history endpoint
// This is a placeholder for future implementation
func (ctrl *FileController) GetFileHistory(c *fiber.Ctx) error {
	// TODO: Implement database storage for file analysis history

	// For now, return empty response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"files":  []interface{}{},
		"total":  0,
		"limit":  20,
		"offset": 0,
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
