package controllers

import (
	"time"

	"chatbot/services"

	"github.com/gofiber/fiber/v2"
)

// AudioController handles audio-related HTTP requests
type AudioController struct {
	openaiService *services.OpenAIService
}

// NewAudioController creates a new audio controller
func NewAudioController(openaiService *services.OpenAIService) *AudioController {
	return &AudioController{
		openaiService: openaiService,
	}
}

// TranscribeResponse represents the audio transcription response
type TranscribeResponse struct {
	Text       string    `json:"text"`
	Language   string    `json:"language"`
	Duration   float64   `json:"duration"`
	Confidence float64   `json:"confidence"`
	Timestamp  time.Time `json:"timestamp"`
}

// TranscribeAudio handles POST /api/audio/transcribe endpoint
func (ctrl *AudioController) TranscribeAudio(c *fiber.Ctx) error {
	// Get file from form data
	file, err := c.FormFile("audio")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "audio file is required",
		})
	}

	// Validate file size (25 MB max)
	const maxFileSize = 25 * 1024 * 1024 // 25 MB in bytes
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"error": "file size exceeds maximum allowed (25MB)",
		})
	}

	// Open the uploaded file
	fileData, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read uploaded file",
		})
	}
	defer fileData.Close()

	// Call OpenAI Whisper service
	transcription, err := ctrl.openaiService.TranscribeAudio(fileData, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to transcribe audio: " + err.Error(),
		})
	}

	// Create response
	response := TranscribeResponse{
		Text:       transcription.Text,
		Language:   transcription.Language,
		Duration:   transcription.Duration,
		Confidence: 0.95, // Mock confidence for now
		Timestamp:  time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
