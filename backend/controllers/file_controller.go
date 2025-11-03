package controllers

import (
	"chatbot/models"
	"chatbot/repositories"
	"chatbot/services"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

// UploadFileRequest represents the form data for file upload
type UploadFileRequest struct {
	// No additional fields needed - just the file
}

// AnalyzeFile handles POST /api/file/analyze endpoint
func (ctrl *FileController) AnalyzeFile(c *fiber.Ctx) error {
	// 1. Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file is required",
		})
	}

	// 2. Get MIME type
	mimeType := file.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream" // Default MIME type
	}

	// 3. Create uploads directory if not exists
	uploadDir := "./uploads"
	if err := ensureDir(uploadDir); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create upload directory",
		})
	}

	// 4. Generate unique filename
	fileID := uuid.New()
	storagePath := fmt.Sprintf("%s/%s_%s", uploadDir, fileID.String(), file.Filename)

	// 5. Save file to disk
	if err := c.SaveFile(file, storagePath); err != nil {
		log.Printf("⚠️  Failed to save file to disk: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save file",
		})
	}

	// 6. Save metadata to database
	fileAnalysis := &models.FileAnalysis{
		ID:          fileID,
		FileName:    file.Filename,
		StoragePath: storagePath,
		MimeType:    mimeType,
		FileSize:    file.Size,
	}

	if err := ctrl.repository.Create(fileAnalysis); err != nil {
		log.Printf("⚠️  Failed to save file metadata to database: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save file metadata",
		})
	}

	// 7. Build response
	response := fiber.Map{
		"file_id":      fileAnalysis.ID.String(),
		"file_name":    fileAnalysis.FileName,
		"storage_path": fileAnalysis.StoragePath,
		"mime_type":    fileAnalysis.MimeType,
		"file_size":    fileAnalysis.FileSize,
		"uploaded_at":  fileAnalysis.UploadedAt,
	}

	log.Printf("✅ File uploaded successfully: %s (ID: %s)", file.Filename, fileAnalysis.ID.String())
	return c.Status(fiber.StatusOK).JSON(response)
}

// parseFileRequest is no longer used - kept for backward compatibility
// Deprecated: Use the simplified AnalyzeFile function instead

// GetFileHistory handles GET /api/file/history endpoint
func (ctrl *FileController) GetFileHistory(c *fiber.Ctx) error {
	// Parse and validate pagination
	limit := c.QueryInt("limit", 20)
	offset := c.QueryInt("offset", 0)

	// Validate pagination
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	fileType := c.Query("file_type", "all")

	// Get files from repository
	files, total, err := ctrl.getFilesByType(fileType, limit, offset)
	if err != nil {
		log.Printf("Failed to fetch file history: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch file history",
		})
	}

	// Build and return response
	response := ctrl.buildFileHistoryResponse(files, total, limit, offset)
	return c.Status(fiber.StatusOK).JSON(response)
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
			"file_id":      file.ID.String(),
			"file_name":    file.FileName,
			"storage_path": file.StoragePath,
			"mime_type":    file.MimeType,
			"file_size":    file.FileSize,
			"uploaded_at":  file.UploadedAt,
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
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file_id is required",
		})
	}

	// TODO: Implement re-analysis from stored files
	return c.Status(fiber.StatusNotImplemented).JSON(fiber.Map{
		"error": "re-analysis feature not yet implemented",
	})
}

// ensureDir creates a directory if it doesn't exist
func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}
