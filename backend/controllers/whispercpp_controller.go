package controllers

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"chatbot/services"

	"github.com/gofiber/fiber/v2"
)

// WhisperCppController handles Whisper.cpp STT HTTP requests
type WhisperCppController struct {
	whisperService *services.WhisperCppService
}

// NewWhisperCppController creates a new WhisperCpp controller
func NewWhisperCppController(whisperService *services.WhisperCppService) *WhisperCppController {
	return &WhisperCppController{
		whisperService: whisperService,
	}
}

// WhisperCppTranscribeRequest represents the transcription request parameters
type WhisperCppTranscribeRequest struct {
	Language   string `form:"language"`   // "th", "en", "auto" (default: "th")
	Timestamps bool   `form:"timestamps"` // Return segments with timestamps (default: false)
	Model      string `form:"model"`      // Model name: "tiny.en", "small", "medium", "large-v2" (default: use config default)
}

// WhisperCppTranscribeResponse represents the basic transcription response
type WhisperCppTranscribeResponse struct {
	Success       bool    `json:"success"`
	Transcription string  `json:"transcription"`
	Confidence    float64 `json:"confidence"`
	Language      string  `json:"language"`
	Model         string  `json:"model"` // Model used for transcription
}

// WhisperCppTranscribeWithTimestampsResponse represents the transcription response with timestamps
type WhisperCppTranscribeWithTimestampsResponse struct {
	Success       bool                             `json:"success"`
	Transcription string                           `json:"transcription"`
	Segments      []services.TranscriptionSegment  `json:"segments"`
	Language      string                           `json:"language"`
	Duration      float64                          `json:"duration"`
	Model         string                           `json:"model"` // Model used for transcription
}

// WhisperCppStatusResponse represents the service status response
type WhisperCppStatusResponse struct {
	Service            string   `json:"service"`
	Available          bool     `json:"available"`
	SupportedFormats   []string `json:"supported_formats"`
	SupportedLanguages []string `json:"supported_languages"`
	SupportedModels    []string `json:"supported_models"`
	DefaultLanguage    string   `json:"default_language"`
	DefaultModel       string   `json:"default_model"`
	CurrentOS          string   `json:"current_os"`
}

// GetStatus handles GET /api/stt/whispercpp/status endpoint
func (ctrl *WhisperCppController) GetStatus(c *fiber.Ctx) error {
	// Check if service is available
	available := ctrl.whisperService.IsAvailable()

	// Get supported formats
	supportedFormats := ctrl.whisperService.GetSupportedFormats()

	// Get supported languages from config (th, en, auto)
	supportedLanguages := []string{"th", "en", "auto"}

	// Get supported models
	supportedModels := ctrl.whisperService.GetSupportedModels()

	// Get default model name from config
	defaultModel := ctrl.whisperService.GetModelName()

	response := WhisperCppStatusResponse{
		Service:            "whisper.cpp",
		Available:          available,
		SupportedFormats:   supportedFormats,
		SupportedLanguages: supportedLanguages,
		SupportedModels:    supportedModels,
		DefaultLanguage:    "auto",
		DefaultModel:       defaultModel,
		CurrentOS:          runtime.GOOS,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// TranscribeAudio handles POST /api/stt/whispercpp endpoint
func (ctrl *WhisperCppController) TranscribeAudio(c *fiber.Ctx) error {
	startTime := time.Now()

	// Parse form parameters
	var req WhisperCppTranscribeRequest
	if err := c.BodyParser(&req); err != nil {
		// Set defaults if parsing fails
		req.Language = "th"
		req.Timestamps = false
	}

	// Set default language if not provided
	if req.Language == "" {
		req.Language = "th"
	}

	// Validate language
	validLanguages := []string{"th", "en", "auto"}
	isValidLanguage := false
	for _, lang := range validLanguages {
		if req.Language == lang {
			isValidLanguage = true
			break
		}
	}
	if !isValidLanguage {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("invalid language: %s (supported: th, en, auto)", req.Language),
		})
	}

	// Get audio file from form data
	file, err := c.FormFile("audio")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "audio file is required",
			"details": err.Error(),
		})
	}

	// Validate file size (25 MB max - same as OpenAI Whisper API limit)
	const maxFileSize = 25 * 1024 * 1024 // 25 MB in bytes
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"success": false,
			"error":   "file size exceeds maximum allowed (25MB)",
		})
	}

	// Validate file size is not zero
	if file.Size == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "audio file is empty",
		})
	}

	// Validate file extension
	supportedFormats := ctrl.whisperService.GetSupportedFormats()
	fileExt := strings.ToLower(strings.TrimPrefix(strings.ToLower(file.Filename), "."))
	parts := strings.Split(file.Filename, ".")
	if len(parts) > 1 {
		fileExt = strings.ToLower(parts[len(parts)-1])
	}

	isValidFormat := false
	for _, format := range supportedFormats {
		if fileExt == format {
			isValidFormat = true
			break
		}
	}
	if !isValidFormat {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   fmt.Sprintf("unsupported audio format: %s (supported: %s)", fileExt, strings.Join(supportedFormats, ", ")),
		})
	}

	// Open the uploaded file
	fileData, err := file.Open()
	if err != nil {
		fmt.Printf("❌ Failed to open file: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to read uploaded file",
			"details": err.Error(),
		})
	}
	defer fileData.Close()

	// Determine which model to use
	modelName := req.Model
	actualModelName := modelName
	if actualModelName == "" {
		actualModelName = ctrl.whisperService.GetModelName()
	}

	// Check if timestamps are requested
	if req.Timestamps {
		// Call Whisper.cpp service with timestamps (and optional model)
		var segments []services.TranscriptionSegment
		var err error

		if modelName != "" {
			// Use specific model
			fmt.Printf("🔄 Transcribing with timestamps (language: %s, model: %s)...\n", req.Language, modelName)
			segments, err = ctrl.whisperService.TranscribeWithTimestampsAndModel(fileData, file.Filename, req.Language, modelName)
		} else {
			// Use default model
			fmt.Printf("🔄 Transcribing with timestamps (language: %s)...\n", req.Language)
			segments, err = ctrl.whisperService.TranscribeWithTimestamps(fileData, file.Filename, req.Language)
		}

		if err != nil {
			fmt.Printf("❌ Transcription failed: %v\n", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"success": false,
				"error":   "failed to transcribe audio",
				"details": err.Error(),
			})
		}

		// Build full transcription text from segments
		var fullText strings.Builder
		for i, segment := range segments {
			if i > 0 {
				fullText.WriteString(" ")
			}
			fullText.WriteString(segment.Text)
		}

		// Calculate total duration
		duration := 0.0
		if len(segments) > 0 {
			duration = segments[len(segments)-1].EndTime
		}

		fmt.Printf("✅ Transcription with timestamps successful (%.2fs)\n", time.Since(startTime).Seconds())

		// Create response with segments
		response := WhisperCppTranscribeWithTimestampsResponse{
			Success:       true,
			Transcription: fullText.String(),
			Segments:      segments,
			Language:      req.Language,
			Duration:      duration,
			Model:         actualModelName,
		}

		return c.Status(fiber.StatusOK).JSON(response)
	}

	// Call Whisper.cpp service for basic transcription (with optional model)
	var transcription string
	var confidence float64
	var transcribeErr error

	if modelName != "" {
		// Use specific model
		fmt.Printf("🔄 Transcribing audio (language: %s, model: %s)...\n", req.Language, modelName)
		transcription, confidence, transcribeErr = ctrl.whisperService.TranscribeWithModel(fileData, file.Filename, req.Language, modelName)
	} else {
		// Use default model
		fmt.Printf("🔄 Transcribing audio (language: %s)...\n", req.Language)
		transcription, confidence, transcribeErr = ctrl.whisperService.Transcribe(fileData, file.Filename, req.Language)
	}

	if transcribeErr != nil {
		fmt.Printf("❌ Transcription failed: %v\n", transcribeErr)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to transcribe audio",
			"details": transcribeErr.Error(),
		})
	}

	fmt.Printf("✅ Transcription successful (%.2fs): %s\n", time.Since(startTime).Seconds(), transcription)

	// Create response
	response := WhisperCppTranscribeResponse{
		Success:       true,
		Transcription: transcription,
		Confidence:    confidence,
		Language:      req.Language,
		Model:         actualModelName,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
