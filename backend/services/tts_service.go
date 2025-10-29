package services

import (
	"context"
	"fmt"
	"io"

	"chatbot/config"

	"github.com/sashabaranov/go-openai"
)

// TTSService handles Text-to-Speech operations using OpenAI TTS API
type TTSService struct {
	client *openai.Client
	config *config.Config
}

// NewTTSService creates a new TTS service
func NewTTSService(cfg *config.Config) *TTSService {
	client := openai.NewClient(cfg.OpenAIAPIKey)
	return &TTSService{
		client: client,
		config: cfg,
	}
}

// TTSRequest represents a text-to-speech request
type TTSRequest struct {
	Text           string  `json:"text" validate:"required,max=4096"`
	Voice          string  `json:"voice"`
	Model          string  `json:"model"`
	ResponseFormat string  `json:"response_format"`
	Speed          float64 `json:"speed"`
}

// TTSResponse represents the response containing audio data
type TTSResponse struct {
	AudioData       []byte  `json:"audio_data"`
	Format          string  `json:"format"`
	Duration        float64 `json:"duration"`
	CharactersUsed  int     `json:"characters_used"`
	Voice           string  `json:"voice"`
}

// ValidVoices contains all available OpenAI TTS voices
var ValidVoices = map[string]bool{
	"alloy":   true,
	"echo":    true,
	"fable":   true,
	"onyx":    true,
	"nova":    true,
	"shimmer": true,
}

// ValidModels contains all available OpenAI TTS models
var ValidModels = map[string]bool{
	"tts-1":    true,
	"tts-1-hd": true,
}

// ValidFormats contains all available audio formats
var ValidFormats = map[string]bool{
	"mp3":  true,
	"opus": true,
	"aac":  true,
	"flac": true,
	"wav":  true,
	"pcm":  true,
}

// TextToSpeech converts text to speech using OpenAI TTS API
func (s *TTSService) TextToSpeech(ctx context.Context, req TTSRequest) (*TTSResponse, error) {
	// Validate and set defaults
	if req.Text == "" {
		return nil, fmt.Errorf("text is required")
	}

	if len(req.Text) > 4096 {
		return nil, fmt.Errorf("text exceeds maximum length of 4096 characters")
	}

	// Set default voice
	if req.Voice == "" {
		req.Voice = "nova"
	}
	if !ValidVoices[req.Voice] {
		return nil, fmt.Errorf("invalid voice. Allowed: alloy, echo, fable, onyx, nova, shimmer")
	}

	// Set default model
	if req.Model == "" {
		req.Model = "tts-1"
	}
	if !ValidModels[req.Model] {
		return nil, fmt.Errorf("invalid model. Allowed: tts-1, tts-1-hd")
	}

	// Set default format
	if req.ResponseFormat == "" {
		req.ResponseFormat = "mp3"
	}
	if !ValidFormats[req.ResponseFormat] {
		return nil, fmt.Errorf("invalid response_format. Allowed: mp3, opus, aac, flac, wav, pcm")
	}

	// Set default speed
	if req.Speed == 0 {
		req.Speed = 1.0
	}
	if req.Speed < 0.25 || req.Speed > 4.0 {
		return nil, fmt.Errorf("speed must be between 0.25 and 4.0")
	}

	// Create TTS request
	ttsReq := openai.CreateSpeechRequest{
		Model:          openai.SpeechModel(req.Model),
		Input:          req.Text,
		Voice:          openai.SpeechVoice(req.Voice),
		ResponseFormat: openai.SpeechResponseFormat(req.ResponseFormat),
		Speed:          req.Speed,
	}

	// Call OpenAI TTS API
	response, err := s.client.CreateSpeech(ctx, ttsReq)
	if err != nil {
		return nil, fmt.Errorf("failed to generate speech: %w", err)
	}
	defer response.Close()

	// Read audio data
	audioData, err := io.ReadAll(response)
	if err != nil {
		return nil, fmt.Errorf("failed to read audio data: %w", err)
	}

	// Estimate duration (rough estimate: ~150 words per minute, ~5 chars per word)
	// This is approximate - actual duration depends on voice speed
	words := float64(len(req.Text)) / 5.0
	duration := (words / 150.0) * 60.0 * (1.0 / req.Speed)

	return &TTSResponse{
		AudioData:      audioData,
		Format:         req.ResponseFormat,
		Duration:       duration,
		CharactersUsed: len(req.Text),
		Voice:          req.Voice,
	}, nil
}

// TextToSpeechSimple is a simplified version for quick TTS generation
func (s *TTSService) TextToSpeechSimple(ctx context.Context, text string) ([]byte, error) {
	req := TTSRequest{
		Text:           text,
		Voice:          "nova",
		Model:          "tts-1",
		ResponseFormat: "mp3",
		Speed:          1.0,
	}

	response, err := s.TextToSpeech(ctx, req)
	if err != nil {
		return nil, err
	}

	return response.AudioData, nil
}
