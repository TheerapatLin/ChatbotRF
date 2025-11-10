package services

import (
	"io"
)

// TranscriptionService กำหนด interface สำหรับ speech-to-text services
// ช่วยให้สามารถใช้งาน STT implementation หลายแบบ (whisper.cpp, Google STT, OpenAI Whisper API, etc.)
// แต่ละ implementation ต้อง implement interface นี้
type TranscriptionService interface {
	// Transcribe แปลงไฟล์ audio เป็นข้อความ
	// Parameters:
	//   - audioFile: io.Reader ที่มีข้อมูล audio
	//   - filename: ชื่อไฟล์เดิม (ใช้เพื่อกำหนดรูปแบบ)
	//   - language: รหัสภาษา (เช่น "th", "en", "auto")
	// Returns:
	//   - transcription: ข้อความที่แปลงได้
	//   - confidence: คะแนนความมั่นใจ (0.0 - 1.0)
	//   - error: ข้อผิดพลาดที่เกิดขึ้น
	Transcribe(audioFile io.Reader, filename string, language string) (transcription string, confidence float64, err error)

	// TranscribeWithTimestamps คืนค่าการแปลงพร้อม timestamps ระดับคำ
	// Parameters:
	//   - audioFile: io.Reader ที่มีข้อมูล audio
	//   - filename: ชื่อไฟล์เดิม
	//   - language: รหัสภาษา
	// Returns:
	//   - segments: array ของ segments การแปลงพร้อม timestamps
	//   - error: ข้อผิดพลาดที่เกิดขึ้น
	TranscribeWithTimestamps(audioFile io.Reader, filename string, language string) (segments []TranscriptionSegment, err error)

	// IsAvailable ตรวจสอบว่า service ตั้งค่าถูกต้องและพร้อมใช้งาน
	// ตรวจสอบเช่น: binary file อยู่ที่ถูกต้อง, model file พร้อมใช้งาน
	IsAvailable() bool

	// GetSupportedFormats คืนรายการรูปแบบ audio ที่รองรับ
	// เช่น: ["wav", "mp3", "m4a", "ogg", "flac"]
	GetSupportedFormats() []string
}

// TranscriptionSegment แทน segment ของ audio ที่แปลงแล้วพร้อม timestamps
// ใช้สำหรับการแปลงที่ต้องการทราบเวลาแต่ละ segment/คำ
type TranscriptionSegment struct {
	StartTime float64 `json:"start_time"` // เวลาเริ่มต้นเป็นวินาที
	EndTime   float64 `json:"end_time"`   // เวลาสิ้นสุดเป็นวินาที
	Text      string  `json:"text"`       // ข้อความที่แปลงได้สำหรับ segment นี้
}

// TranscriptionResponse แทน API response สำหรับคำขอการแปลง
// ใช้สำหรับ API endpoint ที่ return ข้อมูลการแปลงเสียง
type TranscriptionResponse struct {
	Success       bool                   `json:"success"`                 // สถานะความสำเร็จ
	Transcription string                 `json:"transcription"`           // ข้อความที่แปลงได้
	Confidence    float64                `json:"confidence,omitempty"`    // คะแนนความมั่นใจ (0.0 - 1.0)
	Segments      []TranscriptionSegment `json:"segments,omitempty"`      // Segments พร้อม timestamps (ถ้ามี)
	Language      string                 `json:"language"`                // ภาษาที่ตรวจพบหรือที่กำหนด
	Duration      float64                `json:"duration,omitempty"`      // ระยะเวลา audio เป็นวินาที
	Error         string                 `json:"error,omitempty"`         // ข้อความ error (ถ้ามี)
	Model         string                 `json:"model,omitempty"`         // ชื่อ model ที่ใช้ (เช่น "whisper.cpp-small")
	ProcessTime   float64                `json:"process_time,omitempty"`  // เวลาประมวลผล (วินาที)
}
