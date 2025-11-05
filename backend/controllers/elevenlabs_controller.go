package controllers

import (
	"encoding/base64"
	"time"

	"chatbot/services"

	"github.com/gofiber/fiber/v2"
)

// ElevenLabsController handles ElevenLabs TTS HTTP requests
type ElevenLabsController struct {
	elevenLabsService *services.ElevenLabsService
}

// NewElevenLabsController creates a new ElevenLabs controller
func NewElevenLabsController(elevenLabsService *services.ElevenLabsService) *ElevenLabsController {
	return &ElevenLabsController{
		elevenLabsService: elevenLabsService,
	}
}

// ElevenLabsTTSRequest represents the request body for ElevenLabs TTS
type ElevenLabsTTSRequest struct {
	Text            string                  `json:"text" validate:"required"`
	VoiceID         string                  `json:"voice_id,omitempty"`
	ModelID         string                  `json:"model_id,omitempty"`
	Stability       *float64                `json:"stability,omitempty"`
	SimilarityBoost *float64                `json:"similarity_boost,omitempty"`
	Style           *float64                `json:"style,omitempty"`
	Speed           *float64                `json:"speed,omitempty"`
	UseSpeakerBoost *bool                   `json:"use_speaker_boost,omitempty"`
	VoiceSettings   *services.VoiceSettings `json:"voice_settings,omitempty"`
}

// ElevenLabsTTSResponse represents the response from ElevenLabs TTS API
type ElevenLabsTTSResponse struct {
	AudioData       string    `json:"audio_data"`
	Format          string    `json:"format"`
	CharactersUsed  int       `json:"characters_used"`
	ModelID         string    `json:"model_id"`
	VoiceID         string    `json:"voice_id"`
	DurationSeconds float64   `json:"duration_seconds"`
	Timestamp       time.Time `json:"timestamp"`
}

// TextToSpeech handles POST /audio/elevenlabs/tts endpoint
func (ctrl *ElevenLabsController) TextToSpeech(c *fiber.Ctx) error {
	// Parse request body
	var req ElevenLabsTTSRequest
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

	// Build voice settings from individual parameters or use provided voice_settings
	var voiceSettings *services.VoiceSettings
	if req.VoiceSettings != nil {
		voiceSettings = req.VoiceSettings
	} else {
		// Build from individual parameters
		voiceSettings = &services.VoiceSettings{
			Stability:       req.Stability,
			SimilarityBoost: req.SimilarityBoost,
			Style:           req.Style,
			Speed:           req.Speed,
			UseSpeakerBoost: req.UseSpeakerBoost,
		}
	}

	// Create service request
	serviceReq := services.ElevenLabsTTSRequest{
		Text:          req.Text,
		ModelID:       req.ModelID,
		VoiceSettings: voiceSettings,
	}

	// Call ElevenLabs service
	ttsRes, err := ctrl.elevenLabsService.TextToSpeech(c.Context(), req.VoiceID, serviceReq)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Check if client wants binary audio or JSON response
	accept := c.Get("Accept")

	// If client accepts audio/* or application/octet-stream, return binary
	if accept == "audio/*" || accept == "audio/mpeg" || accept == "application/octet-stream" {
		c.Set("Content-Type", "audio/mpeg")
		c.Set("X-Audio-Duration", string(rune(int(ttsRes.DurationSeconds))))
		c.Set("X-Characters-Used", string(rune(ttsRes.CharactersUsed)))
		c.Set("X-Voice-ID", ttsRes.VoiceID)
		c.Set("X-Model-ID", ttsRes.ModelID)

		return c.Send(ttsRes.AudioData)
	}

	// Default: Return JSON with base64-encoded audio
	audioBase64 := base64.StdEncoding.EncodeToString(ttsRes.AudioData)

	response := ElevenLabsTTSResponse{
		AudioData:       audioBase64,
		Format:          ttsRes.Format,
		CharactersUsed:  ttsRes.CharactersUsed,
		ModelID:         ttsRes.ModelID,
		VoiceID:         ttsRes.VoiceID,
		DurationSeconds: ttsRes.DurationSeconds,
		Timestamp:       time.Now(),
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetVoices handles GET /audio/elevenlabs/voices endpoint
func (ctrl *ElevenLabsController) GetVoices(c *fiber.Ctx) error {
	voices, err := ctrl.elevenLabsService.GetVoices(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"voices": voices,
		"count":  len(voices),
	})
}
