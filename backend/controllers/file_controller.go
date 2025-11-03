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

// UploadFiles handles POST /api/file/upload endpoint
// Supports multiple file uploads (max 5 files)
func (ctrl *FileController) UploadFiles(c *fiber.Ctx) error {
	// 1. Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to parse form data",
		})
	}

	// 2. Get uploaded files
	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "at least one file is required",
		})
	}

	// 3. Validate file count (max 5 files)
	if len(files) > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "maximum 5 files allowed per upload",
		})
	}

	// 4. Create uploads directory if not exists
	uploadDir := "./uploads"
	if err := ensureDir(uploadDir); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create upload directory",
		})
	}

	// 5. Process each file
	uploadedFiles := make([]fiber.Map, 0, len(files))
	failedFiles := make([]fiber.Map, 0)

	for _, file := range files {
		// Get MIME type
		mimeType := file.Header.Get("Content-Type")
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		// Generate unique filename
		fileID := uuid.New()
		storagePath := fmt.Sprintf("%s/%s_%s", uploadDir, fileID.String(), file.Filename)

		// Save file to disk
		if err := c.SaveFile(file, storagePath); err != nil {
			log.Printf("⚠️  Failed to save file to disk: %s - %v", file.Filename, err)
			failedFiles = append(failedFiles, fiber.Map{
				"file_name": file.Filename,
				"error":     "failed to save file to disk",
			})
			continue
		}

		// Save metadata to database
		fileAnalysis := &models.FileAnalysis{
			ID:          fileID,
			FileName:    file.Filename,
			StoragePath: storagePath,
			MimeType:    mimeType,
			FileSize:    file.Size,
		}

		if err := ctrl.repository.Create(fileAnalysis); err != nil {
			log.Printf("⚠️  Failed to save file metadata to database: %s - %v", file.Filename, err)
			// Try to delete the file from disk if DB save failed
			os.Remove(storagePath)
			failedFiles = append(failedFiles, fiber.Map{
				"file_name": file.Filename,
				"error":     "failed to save file metadata",
			})
			continue
		}

		// Add to successful uploads
		uploadedFiles = append(uploadedFiles, fiber.Map{
			"file_id":      fileAnalysis.ID.String(),
			"file_name":    fileAnalysis.FileName,
			"storage_path": fileAnalysis.StoragePath,
			"mime_type":    fileAnalysis.MimeType,
			"file_size":    fileAnalysis.FileSize,
			"uploaded_at":  fileAnalysis.UploadedAt,
		})

		log.Printf("✅ File uploaded successfully: %s (ID: %s)", file.Filename, fileAnalysis.ID.String())
	}

	// 6. Build response
	response := fiber.Map{
		"success":        len(uploadedFiles),
		"failed":         len(failedFiles),
		"total":          len(files),
		"uploaded_files": uploadedFiles,
	}

	// Include failed files info if any
	if len(failedFiles) > 0 {
		response["failed_files"] = failedFiles
	}

	// Return appropriate status code
	if len(uploadedFiles) == 0 {
		return c.Status(fiber.StatusInternalServerError).JSON(response)
	}

	if len(failedFiles) > 0 {
		return c.Status(fiber.StatusPartialContent).JSON(response)
	}

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

// DeleteAllFiles handles DELETE /api/file/uploads endpoint
func (ctrl *FileController) DeleteAllFiles(c *fiber.Ctx) error {
	// Delete all file records from database
	if err := ctrl.repository.DeleteAll(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file records",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "All file records deleted successfully",
	})
}

// ensureDir creates a directory if it doesn't exist
func ensureDir(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		return os.MkdirAll(dirPath, 0755)
	}
	return nil
}
