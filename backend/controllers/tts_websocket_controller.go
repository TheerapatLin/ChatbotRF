package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	"chatbot/repositories"
	"chatbot/services"

	"github.com/gofiber/contrib/websocket"
)

// TTSWebSocketController handles WebSocket connections for TTS streaming
type TTSWebSocketController struct {
	ttsService     *services.TTSService
	personaRepo    *repositories.PersonaRepository
	activeSessions sync.Map // session_id -> context.CancelFunc
}

// NewTTSWebSocketController creates a new TTS WebSocket controller
func NewTTSWebSocketController(
	ttsService *services.TTSService,
	personaRepo *repositories.PersonaRepository,
) *TTSWebSocketController {
	return &TTSWebSocketController{
		ttsService:  ttsService,
		personaRepo: personaRepo,
	}
}

// TTSWebSocketRequest represents incoming TTS WebSocket request
type TTSWebSocketRequest struct {
	Type         string  `json:"type"`          // "tts" or "stop"
	SessionID    string  `json:"session_id"`    // Unique session ID
	Text         string  `json:"text"`          // Text + emotion/style instructions
	Voice        string  `json:"voice"`         // alloy, echo, fable, onyx, nova, shimmer (default: nova)
	Model        string  `json:"model"`         // tts-1, tts-1-hd, gpt-4o-mini-tts (default: gpt-4o-mini-tts)
	Speed        float64 `json:"speed"`         // 0.25 - 4.0 (default: 1.0)
	EmotionRange string  `json:"emotion_range"` // Optional: neutral, happy, sad, excited, calm, serious
	StyleHint    string  `json:"style_hint"`    // Optional: formal, casual, empathetic
	PersonaID    *int    `json:"persona_id"`    // Optional persona ID
}

// TTSAudioChunk represents audio chunk sent to client
type TTSAudioChunk struct {
	Type        string `json:"type"`         // "audio_chunk"
	AudioData   string `json:"audio_data"`   // Base64 encoded audio
	Format      string `json:"format"`       // "mp3"
	ChunkIndex  int    `json:"chunk_index"`  // Chunk index
	TotalChunks int    `json:"total_chunks"` // Total chunks (estimated)
}

// TTSCompleted represents completion message
type TTSCompleted struct {
	Type           string  `json:"type"`            // "completed"
	Duration       float64 `json:"duration"`        // Audio duration in seconds
	CharactersUsed int     `json:"characters_used"` // Characters processed
	Format         string  `json:"format"`          // Audio format
}

// TTSError represents error message
type TTSError struct {
	Type  string `json:"type"`  // "error"
	Error string `json:"error"` // Error message
}

// TTSStopped represents stop confirmation
type TTSStopped struct {
	Type string `json:"type"` // "stopped"
}

// HandleTTSWebSocket handles WebSocket connections for TTS streaming
func (ctrl *TTSWebSocketController) HandleTTSWebSocket(c *websocket.Conn) {
	log.Printf("🔊 New TTS WebSocket connection from %s", c.RemoteAddr())

	// Clean up on disconnect
	defer func() {
		log.Printf("🔇 TTS WebSocket connection closed from %s", c.RemoteAddr())
		c.Close()
	}()

	// Message loop
	for {
		var req TTSWebSocketRequest

		// Read message from client
		if err := c.ReadJSON(&req); err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("❌ TTS WebSocket error: %v", err)
			}
			break
		}

		// Handle message type
		switch req.Type {
		case "tts":
			if err := ctrl.handleTTS(c, req); err != nil {
				log.Printf("❌ TTS Error: %v", err)
				ctrl.sendError(c, err.Error())
			}

		case "stop":
			ctrl.handleStop(req.SessionID)
			ctrl.sendStopped(c)

		default:
			ctrl.sendError(c, fmt.Sprintf("Unknown message type: %s", req.Type))
		}
	}
}

// handleTTS processes TTS request and streams audio
func (ctrl *TTSWebSocketController) handleTTS(c *websocket.Conn, req TTSWebSocketRequest) error {
	// 1. Validate request
	if req.Text == "" {
		return fmt.Errorf("text is required")
	}

	if req.SessionID == "" {
		return fmt.Errorf("session_id is required")
	}

	// 2. Apply persona voice mapping if persona_id provided
	voice := req.Voice
	if voice == "" {
		voice = "nova" // Default voice
	}

	if req.PersonaID != nil {
		persona, err := ctrl.personaRepo.FindByID(*req.PersonaID)
		if err == nil && persona != nil {
			// Map emotion to voice (optional enhancement)
			voice = ctrl.mapEmotionToVoice(req.EmotionRange, persona.Tone)
		}
	}

	// 3. Set defaults
	speed := req.Speed
	if speed == 0 {
		speed = 1.0
	}

	// Set default model (allow user to choose)
	model := req.Model
	if model == "" {
		model = services.DefaultTTSModel // gpt-4o-mini-tts
	}

	// 4. Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	ctrl.activeSessions.Store(req.SessionID, cancel)
	defer func() {
		cancel()
		ctrl.activeSessions.Delete(req.SessionID)
	}()

	// 5. Call TTS Service
	ttsReq := services.TTSRequest{
		Text:           req.Text,
		Voice:          voice,
		Model:          model,
		Speed:          speed,
		ResponseFormat: "mp3",
	}

	ttsResp, err := ctrl.ttsService.TextToSpeech(ctx, ttsReq)
	if err != nil {
		return fmt.Errorf("TTS generation failed: %w", err)
	}

	// 6. Split audio into chunks and stream
	if err := ctrl.streamAudioChunks(c, ttsResp, req.SessionID); err != nil {
		return fmt.Errorf("streaming failed: %w", err)
	}

	// 7. Send completion message
	if err := ctrl.sendCompleted(c, ttsResp); err != nil {
		return fmt.Errorf("failed to send completion: %w", err)
	}

	log.Printf("✅ TTS completed for session %s (%d characters)", req.SessionID, ttsResp.CharactersUsed)
	return nil
}

// streamAudioChunks splits audio data into chunks and sends via WebSocket
func (ctrl *TTSWebSocketController) streamAudioChunks(c *websocket.Conn, ttsResp *services.TTSResponse, sessionID string) error {
	const chunkSize = 64 * 1024 // 64KB chunks

	audioData := ttsResp.AudioData
	totalSize := len(audioData)
	totalChunks := (totalSize + chunkSize - 1) / chunkSize

	log.Printf("📦 Streaming %d bytes in %d chunks for session %s", totalSize, totalChunks, sessionID)

	for i := 0; i < totalChunks; i++ {
		// Check if session was cancelled
		if _, exists := ctrl.activeSessions.Load(sessionID); !exists {
			log.Printf("🛑 Session %s cancelled", sessionID)
			return fmt.Errorf("session cancelled")
		}

		// Calculate chunk boundaries
		start := i * chunkSize
		end := start + chunkSize
		if end > totalSize {
			end = totalSize
		}

		// Extract chunk
		chunk := audioData[start:end]

		// Encode to base64
		encodedChunk := base64.StdEncoding.EncodeToString(chunk)

		// Send chunk
		chunkMsg := TTSAudioChunk{
			Type:        "audio_chunk",
			AudioData:   encodedChunk,
			Format:      ttsResp.Format,
			ChunkIndex:  i,
			TotalChunks: totalChunks,
		}

		if err := c.WriteJSON(chunkMsg); err != nil {
			return fmt.Errorf("failed to send chunk %d: %w", i, err)
		}

		log.Printf("📤 Sent chunk %d/%d (%d bytes) for session %s", i+1, totalChunks, len(chunk), sessionID)
	}

	return nil
}

// handleStop cancels an active TTS session
func (ctrl *TTSWebSocketController) handleStop(sessionID string) {
	if cancel, exists := ctrl.activeSessions.LoadAndDelete(sessionID); exists {
		if cancelFunc, ok := cancel.(context.CancelFunc); ok {
			cancelFunc()
			log.Printf("🛑 Stopped session %s", sessionID)
		}
	}
}

// sendError sends error message to client
func (ctrl *TTSWebSocketController) sendError(c *websocket.Conn, errorMsg string) {
	c.WriteJSON(TTSError{
		Type:  "error",
		Error: errorMsg,
	})
}

// sendCompleted sends completion message to client
func (ctrl *TTSWebSocketController) sendCompleted(c *websocket.Conn, ttsResp *services.TTSResponse) error {
	return c.WriteJSON(TTSCompleted{
		Type:           "completed",
		Duration:       ttsResp.Duration,
		CharactersUsed: ttsResp.CharactersUsed,
		Format:         ttsResp.Format,
	})
}

// sendStopped sends stop confirmation to client
func (ctrl *TTSWebSocketController) sendStopped(c *websocket.Conn) {
	c.WriteJSON(TTSStopped{
		Type: "stopped",
	})
}

// mapEmotionToVoice maps emotion and persona tone to appropriate voice
func (ctrl *TTSWebSocketController) mapEmotionToVoice(emotion, tone string) string {
	// Emotion-based voice mapping
	emotionMap := map[string]string{
		"happy":   "nova",    // Cheerful female voice
		"sad":     "onyx",    // Deep, somber voice
		"excited": "shimmer", // Energetic voice
		"calm":    "alloy",   // Neutral, calm voice
		"serious": "echo",    // Professional male voice
		"neutral": "nova",    // Default
	}

	if voice, exists := emotionMap[emotion]; exists {
		return voice
	}

	// Fallback: tone-based mapping
	switch tone {
	case "friendly":
		return "nova"
	case "professional":
		return "echo"
	case "empathetic":
		return "shimmer"
	default:
		return "nova"
	}
}
