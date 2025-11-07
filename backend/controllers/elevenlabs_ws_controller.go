package controllers

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"

	"chatbot/services"
	"chatbot/utils"

	"github.com/gofiber/contrib/websocket"
)

// ElevenLabsWSController - ควบคุมการทำงานของ WebSocket สำหรับ ElevenLabs TTS
type ElevenLabsWSController struct {
	elevenLabsService *services.ElevenLabsService
}

// NewElevenLabsWSController - สร้าง controller instance ใหม่สำหรับจัดการ ElevenLabs WebSocket
func NewElevenLabsWSController(elevenLabsService *services.ElevenLabsService) *ElevenLabsWSController {
	return &ElevenLabsWSController{
		elevenLabsService: elevenLabsService,
	}
}

// WebSocketMessage - โครงสร้างของ message ที่รับจาก client
type WebSocketMessage struct {
	Type            string   `json:"type"`                       // ประเภทของ message (tts, stop)
	SessionID       string   `json:"session_id,omitempty"`       // ID ของ session
	Text            string   `json:"text,omitempty"`             // ข้อความที่ต้องการแปลงเป็นเสียง
	VoiceID         string   `json:"voice_id,omitempty"`         // ID ของเสียงที่ใช้
	ModelID         string   `json:"model_id,omitempty"`         // ID ของ model (eleven_multilingual_v2, etc.)
	Stability       *float64 `json:"stability,omitempty"`        // ความคงที่ของเสียง (0.0-1.0)
	SimilarityBoost *float64 `json:"similarity_boost,omitempty"` // ความคล้ายกับเสียงต้นฉบับ (0.0-1.0)
	Style           *float64 `json:"style,omitempty"`            // สไตล์การพูด (0.0-1.0)
	Speed           *float64 `json:"speed,omitempty"`            // ความเร็วในการพูด (0.7-1.2)
}

// AudioChunkResponse - โครงสร้างของ response ที่ส่งกลับไปยัง client
type AudioChunkResponse struct {
	Type        string `json:"type"`         // ประเภทของ response (audio_chunk)
	ChunkIndex  int    `json:"chunk_index"`  // ลำดับของ chunk (เริ่มจาก 0)
	TotalChunks int    `json:"total_chunks"` // จำนวน chunk ทั้งหมด
	AudioData   string `json:"audio_data"`   // ข้อมูลเสียงในรูปแบบ base64
	Text        string `json:"text"`         // ข้อความของ chunk นี้
	Format      string `json:"format"`       // รูปแบบของไฟล์เสียง (mp3)
}

// HandleElevenLabsWebSocket - จัดการ WebSocket connection สำหรับ ElevenLabs TTS
// รับ message จาก client และส่งต่อไปยัง handler ที่เหมาะสม
func (ctrl *ElevenLabsWSController) HandleElevenLabsWebSocket(c *websocket.Conn) {
	log.Printf("🔌 ElevenLabs WebSocket connected: %s", c.RemoteAddr())

	// ปิด connection และ log เมื่อจบการทำงาน
	defer func() {
		c.Close()
		log.Printf("🔌 ElevenLabs WebSocket disconnected: %s", c.RemoteAddr())
	}()

	// วนลูปรับ message จาก client
	for {
		var msg WebSocketMessage
		err := c.ReadJSON(&msg)
		if err != nil {
			log.Printf("❌ Error reading message: %v", err)
			break
		}

		// แยก message ตาม type ไปยัง handler ที่เหมาะสม
		switch msg.Type {
		case "tts":
			ctrl.handleTTSRequest(c, msg)
		case "stop":
			ctrl.handleStopRequest(c, msg)
		default:
			log.Printf("⚠️ Unknown message type: %s", msg.Type)
		}
	}
}

// handleTTSRequest - ประมวลผล TTS request โดยแบ่งข้อความเป็น chunks และส่งไปยัง ElevenLabs API
func (ctrl *ElevenLabsWSController) handleTTSRequest(c *websocket.Conn, msg WebSocketMessage) {
	log.Printf("📝 TTS Request received - SessionID: %s", msg.SessionID)

	// ตรวจสอบว่ามีข้อความหรือไม่
	if msg.Text == "" {
		ctrl.sendError(c, "Text is required")
		return
	}

	// แบ่งข้อความเป็น chunks ด้วย text chunker
	chunks := utils.ChunkText(msg.Text)
	totalChunks := len(chunks)

	// log ข้อความที่ถูกแบ่งแล้ว (เฉพาะข้อความรวม ไม่ใช่แต่ละ chunk)
	log.Printf("📦 Text chunked into %d chunks: %s", totalChunks, msg.Text)

	// วนลูปประมวลผลแต่ละ chunk
	for i, chunkText := range chunks {
		log.Printf("🎤 Processing chunk %d/%d: '%s'", i+1, totalChunks, chunkText)

		// เตรียม voice settings จาก message ที่ได้รับ
		voiceSettings := &services.VoiceSettings{
			Stability:       msg.Stability,
			SimilarityBoost: msg.SimilarityBoost,
			Style:           msg.Style,
			Speed:           msg.Speed,
		}

		// สร้าง TTS request สำหรับ chunk นี้
		ttsReq := services.ElevenLabsTTSRequest{
			Text:          chunkText,
			ModelID:       msg.ModelID,
			VoiceSettings: voiceSettings,
		}

		// เรียก ElevenLabs API เพื่อแปลงข้อความเป็นเสียง
		ttsRes, err := ctrl.elevenLabsService.TextToSpeech(
			context.Background(),
			msg.VoiceID,
			ttsReq,
		)
		if err != nil {
			log.Printf("❌ ElevenLabs API error for chunk %d: %v", i+1, err)
			ctrl.sendError(c, fmt.Sprintf("TTS error: %v", err))
			return
		}

		// log ขนาดของ audio data ที่ได้รับ
		log.Printf("✅ Chunk %d/%d processed successfully (%d bytes)",
			i+1, totalChunks, len(ttsRes.AudioData))

		// แปลง audio data เป็น base64 เพื่อส่งผ่าน WebSocket
		audioBase64 := base64.StdEncoding.EncodeToString(ttsRes.AudioData)
		response := AudioChunkResponse{
			Type:        "audio_chunk",
			ChunkIndex:  i,
			TotalChunks: totalChunks,
			AudioData:   audioBase64,
			Text:        chunkText,
			Format:      ttsRes.Format,
		}

		// ส่ง audio chunk กลับไปยัง client ผ่าน WebSocket
		if err := c.WriteJSON(response); err != nil {
			log.Printf("❌ Error sending chunk %d: %v", i+1, err)
			return
		}
	}

	// ส่ง completion message เมื่อประมวลผลเสร็จสมบูรณ์
	log.Printf("🎉 TTS completed - SessionID: %s", msg.SessionID)
	c.WriteJSON(map[string]interface{}{
		"type":         "completed",
		"session_id":   msg.SessionID,
		"total_chunks": totalChunks,
	})
}

// sendError - ส่ง error message กลับไปยัง client ผ่าน WebSocket
func (ctrl *ElevenLabsWSController) sendError(c *websocket.Conn, message string) {
	log.Printf("⚠️ Sending error: %s", message)
	c.WriteJSON(map[string]interface{}{
		"type":  "error",
		"error": message,
	})
}

// handleStopRequest - จัดการ stop request จาก client
func (ctrl *ElevenLabsWSController) handleStopRequest(c *websocket.Conn, msg WebSocketMessage) {
	log.Printf("🛑 Stop request - SessionID: %s", msg.SessionID)

	// ส่ง stopped message กลับไปยัง client
	c.WriteJSON(map[string]interface{}{
		"type":       "stopped",
		"session_id": msg.SessionID,
	})
}
