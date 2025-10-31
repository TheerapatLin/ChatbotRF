package controllers

import (
	"encoding/base64"
	"time"

	"chatbot/services"

	"github.com/gofiber/fiber/v2"
)

// AudioController handles audio-related HTTP requests
type AudioController struct {
	openaiService *services.OpenAIService
	ttsService    *services.TTSService
}

// NewAudioController creates a new audio controller
func NewAudioController(openaiService *services.OpenAIService, ttsService *services.TTSService) *AudioController {
	return &AudioController{
		openaiService: openaiService,
		ttsService:    ttsService,
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
			"details": err.Error(),
		})
	}

	// Validate file size (25 MB max)
	const maxFileSize = 25 * 1024 * 1024 // 25 MB in bytes
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"error": "file size exceeds maximum allowed (25MB)",
		})
	}

	// Validate file size is not zero
	if file.Size == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "audio file is empty",
		})
	}

	// Open the uploaded file
	fileData, err := file.Open()
	if err != nil {
		println("❌ Failed to open file:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read uploaded file",
			"details": err.Error(),
		})
	}
	defer fileData.Close()

	// Call OpenAI Whisper service
	println("🔄 Calling OpenAI Whisper API...")
	transcription, err := ctrl.openaiService.TranscribeAudio(fileData, file.Filename)
	if err != nil {
		println("❌ Transcription failed:", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to transcribe audio",
			"details": err.Error(),
		})
	}

	println("✅ Transcription successful:", transcription.Text)

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

// TTSRequest represents the TTS request body
type TTSRequest struct {
	Text           string   `json:"text" validate:"required,max=4096"`
	Voice          string   `json:"voice"`
	Model          string   `json:"model"`
	ResponseFormat string   `json:"response_format"`
	Speed          *float64 `json:"speed"` // Use pointer to allow null/omitted values
}

// TTSResponse represents the TTS response with audio data
type TTSResponse struct {
	AudioData      string    `json:"audio_data"`
	AudioURL       string    `json:"audio_url,omitempty"`
	Format         string    `json:"format"`
	Duration       float64   `json:"duration"`
	CharactersUsed int       `json:"characters_used"`
	Voice          string    `json:"voice"`
	Timestamp      time.Time `json:"timestamp"`
}

// TextToSpeech handles POST /api/audio/tts endpoint
func (ctrl *AudioController) TextToSpeech(c *fiber.Ctx) error {
	// Parse request body
	var req TTSRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Validate text field
	if req.Text == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "text is required",
		})
	}

	if len(req.Text) > 4096 {
		return c.Status(fiber.StatusRequestEntityTooLarge).JSON(fiber.Map{
			"error": "text is too long (max 4096 characters)",
		})
	}

	// Create TTS service request
	speed := 1.0
	if req.Speed != nil {
		speed = *req.Speed
	}

	ttsReq := services.TTSRequest{
		Text:           req.Text,
		Voice:          req.Voice,
		Model:          req.Model,
		ResponseFormat: req.ResponseFormat,
		Speed:          speed,
	}

	// Call TTS service
	ttsRes, err := ctrl.ttsService.TextToSpeech(c.Context(), ttsReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if client wants binary audio or JSON response
	accept := c.Get("Accept")

	// If client accepts audio/* or application/octet-stream, return binary
	if accept == "audio/*" || accept == "audio/mpeg" || accept == "application/octet-stream" {
		// Set content type based on format
		contentType := "audio/mpeg"
		switch ttsRes.Format {
		case "opus":
			contentType = "audio/opus"
		case "aac":
			contentType = "audio/aac"
		case "flac":
			contentType = "audio/flac"
		case "wav":
			contentType = "audio/wav"
		case "pcm":
			contentType = "audio/pcm"
		}

		c.Set("Content-Type", contentType)
		c.Set("X-Audio-Duration", string(rune(int(ttsRes.Duration))))
		c.Set("X-Characters-Used", string(rune(ttsRes.CharactersUsed)))

		return c.Send(ttsRes.AudioData)
	}

	// Default: Return JSON with base64-encoded audio
	audioBase64 := base64.StdEncoding.EncodeToString(ttsRes.AudioData)

	response := TTSResponse{
		AudioData:      audioBase64,
		Format:         ttsRes.Format,
		Duration:       ttsRes.Duration,
		CharactersUsed: ttsRes.CharactersUsed,
		Voice:          ttsRes.Voice,
		Timestamp:      time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
