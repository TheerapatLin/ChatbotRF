package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"chatbot/config"
)

// ElevenLabsService handles Text-to-Speech operations using ElevenLabs API
type ElevenLabsService struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// NewElevenLabsService creates a new ElevenLabs service
func NewElevenLabsService(cfg *config.Config) *ElevenLabsService {
	return &ElevenLabsService{
		apiKey:     cfg.ElevenLabsAPIKey,
		baseURL:    "https://api.elevenlabs.io/v1",
		httpClient: &http.Client{},
	}
}

// VoiceSettings represents the voice configuration for TTS
type VoiceSettings struct {
	Stability       *float64 `json:"stability,omitempty"`        // 0.0-1.0: Controls emotional consistency
	SimilarityBoost *float64 `json:"similarity_boost,omitempty"` // 0.0-1.0: How much to match original voice
	Style           *float64 `json:"style,omitempty"`            // 0.0-1.0: Exaggeration of speaking style
	Speed           *float64 `json:"speed,omitempty"`            // 0.7-1.2: Speaking speed (1.0 = normal)
	UseSpeakerBoost *bool    `json:"use_speaker_boost,omitempty"`
}

// ElevenLabsTTSRequest represents a text-to-speech request to ElevenLabs
type ElevenLabsTTSRequest struct {
	Text          string         `json:"text" validate:"required"`
	ModelID       string         `json:"model_id,omitempty"`
	VoiceSettings *VoiceSettings `json:"voice_settings,omitempty"`
}

// ElevenLabsTTSResponse represents the response from ElevenLabs TTS API
type ElevenLabsTTSResponse struct {
	AudioData       []byte  `json:"audio_data"`
	Format          string  `json:"format"`
	CharactersUsed  int     `json:"characters_used"`
	ModelID         string  `json:"model_id"`
	VoiceID         string  `json:"voice_id"`
	DurationSeconds float64 `json:"duration_seconds"`
}

// ValidateVoiceSettings validates voice settings parameters
func ValidateVoiceSettings(vs *VoiceSettings) error {
	if vs == nil {
		return nil
	}

	if vs.Stability != nil {
		if *vs.Stability < 0.0 || *vs.Stability > 1.0 {
			return fmt.Errorf("stability must be between 0.0 and 1.0")
		}
	}

	if vs.SimilarityBoost != nil {
		if *vs.SimilarityBoost < 0.0 || *vs.SimilarityBoost > 1.0 {
			return fmt.Errorf("similarity_boost must be between 0.0 and 1.0")
		}
	}

	if vs.Style != nil {
		if *vs.Style < 0.0 || *vs.Style > 1.0 {
			return fmt.Errorf("style must be between 0.0 and 1.0")
		}
	}

	if vs.Speed != nil {
		if *vs.Speed < 0.7 || *vs.Speed > 1.2 {
			return fmt.Errorf("speed must be between 0.7 and 1.2")
		}
	}

	return nil
}

// TextToSpeech converts text to speech using ElevenLabs API
func (s *ElevenLabsService) TextToSpeech(ctx context.Context, voiceID string, req ElevenLabsTTSRequest) (*ElevenLabsTTSResponse, error) {
	// Validate API key
	if s.apiKey == "" {
		return nil, fmt.Errorf("ElevenLabs API key is not configured")
	}

	// Validate text
	if req.Text == "" {
		return nil, fmt.Errorf("text is required")
	}

	// Validate voice settings
	if err := ValidateVoiceSettings(req.VoiceSettings); err != nil {
		return nil, fmt.Errorf("invalid voice settings: %w", err)
	}

	// Set default model if not provided
	if req.ModelID == "" {
		req.ModelID = "eleven_multilingual_v2"
	}

	// Set default voice settings if not provided
	if req.VoiceSettings == nil {
		stability := 0.5
		similarityBoost := 0.75
		req.VoiceSettings = &VoiceSettings{
			Stability:       &stability,
			SimilarityBoost: &similarityBoost,
		}
	}

	// Default voice ID if not provided
	if voiceID == "" {
		voiceID = "21m00Tcm4TlvDq8ikWAM" // Rachel - Default ElevenLabs voice
	}

	// Prepare request payload
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/text-to-speech/%s", s.baseURL, voiceID)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("xi-api-key", s.apiKey)

	// Send request
	resp, err := s.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ElevenLabs API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	// Read audio data
	audioData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	// Calculate approximate duration (rough estimate based on text length)
	// Average speech rate: ~150 words per minute, ~5 characters per word
	words := float64(len(req.Text)) / 5.0
	speed := 1.0
	if req.VoiceSettings != nil && req.VoiceSettings.Speed != nil {
		speed = *req.VoiceSettings.Speed
	}
	durationSeconds := (words / 150.0) * 60.0 * (1.0 / speed)

	return &ElevenLabsTTSResponse{
		AudioData:       audioData,
		Format:          "mp3",
		CharactersUsed:  len(req.Text),
		ModelID:         req.ModelID,
		VoiceID:         voiceID,
		DurationSeconds: durationSeconds,
	}, nil
}

// GetVoices retrieves available voices from ElevenLabs
func (s *ElevenLabsService) GetVoices(ctx context.Context) ([]map[string]interface{}, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("ElevenLabs API key is not configured")
	}

	url := fmt.Sprintf("%s/voices", s.baseURL)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("xi-api-key", s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("ElevenLabs API error (status %d): %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Voices []map[string]interface{} `json:"voices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Voices, nil
}
