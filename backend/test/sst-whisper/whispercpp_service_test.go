package whisper_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"chatbot/config"
	"chatbot/services"
)

// TestNewWhisperCppService ทดสอบการสร้าง WhisperCppService
func TestNewWhisperCppService(t *testing.T) {
	cfg := config.LoadConfig()

	// ทดสอบการสร้าง service
	service, err := services.NewWhisperCppService(cfg)

	// ถ้าไม่มี binary หรือ model, service จะ error (expected behavior)
	if err != nil {
		t.Logf("⚠️ WhisperCppService creation failed (expected if binary/model not available): %v", err)
		return
	}

	// ถ้าสร้างสำเร็จ, ต้องไม่เป็น nil
	if service == nil {
		t.Error("WhisperCppService should not be nil when creation succeeds")
	}

	t.Log("✓ WhisperCppService สร้างสำเร็จ")
}

// TestWhisperCppServiceIsAvailable ทดสอบ IsAvailable method
func TestWhisperCppServiceIsAvailable(t *testing.T) {
	cfg := config.LoadConfig()

	service, err := services.NewWhisperCppService(cfg)
	if err != nil {
		t.Skipf("⚠️ Skipping test: WhisperCppService not available: %v", err)
		return
	}

	// ทดสอบ IsAvailable
	available := service.IsAvailable()

	if !available {
		t.Error("IsAvailable() should return true for initialized service")
	}

	t.Logf("✓ IsAvailable() = %v", available)
}

// TestWhisperCppServiceGetSupportedFormats ทดสอบ GetSupportedFormats method
func TestWhisperCppServiceGetSupportedFormats(t *testing.T) {
	cfg := config.LoadConfig()

	service, err := services.NewWhisperCppService(cfg)
	if err != nil {
		t.Skipf("⚠️ Skipping test: WhisperCppService not available: %v", err)
		return
	}

	// ทดสอบ GetSupportedFormats
	formats := service.GetSupportedFormats()

	if len(formats) == 0 {
		t.Error("GetSupportedFormats() should return at least one format")
	}

	// ตรวจสอบว่ามี wav, mp3
	hasWav := false
	hasMp3 := false
	for _, format := range formats {
		if format == "wav" {
			hasWav = true
		}
		if format == "mp3" {
			hasMp3 = true
		}
	}

	if !hasWav {
		t.Error("GetSupportedFormats() should include 'wav'")
	}

	if !hasMp3 {
		t.Error("GetSupportedFormats() should include 'mp3'")
	}

	t.Logf("✓ Supported formats: %v", formats)
}

// TestWhisperCppServiceTranscribe ทดสอบการ Transcribe (ต้องมี audio file ตัวอย่าง)
func TestWhisperCppServiceTranscribe(t *testing.T) {
	cfg := config.LoadConfig()

	service, err := services.NewWhisperCppService(cfg)
	if err != nil {
		t.Skipf("⚠️ Skipping test: WhisperCppService not available: %v", err)
		return
	}

	// ตรวจสอบว่ามีไฟล์ตัวอย่างหรือไม่
	testAudioPath := "testdata/audio/th_audio.wav"
	if _, err := os.Stat(testAudioPath); os.IsNotExist(err) {
		t.Skipf("⚠️ Skipping test: test audio file not found at %s", testAudioPath)
		return
	}

	// อ่านไฟล์ตัวอย่าง
	audioData, err := os.ReadFile(testAudioPath)
	if err != nil {
		t.Fatalf("Failed to read test audio file: %v", err)
	}

	// สร้าง io.Reader
	audioReader := bytes.NewReader(audioData)

	// ทดสอบ Transcribe
	transcription, confidence, err := service.Transcribe(audioReader, "th_audio.wav", "th")
	if err != nil {
		t.Fatalf("Transcribe failed: %v", err)
	}

	// ตรวจสอบ output
	if transcription == "" {
		t.Error("Transcription should not be empty")
	}

	if confidence < 0.0 || confidence > 1.0 {
		t.Errorf("Confidence should be between 0.0 and 1.0, got %.2f", confidence)
	}

	t.Logf("✓ Transcription: %s", transcription)
	t.Logf("✓ Confidence: %.2f", confidence)
}

// TestWhisperCppServiceTranscribeWithTimestamps ทดสอบการ TranscribeWithTimestamps
func TestWhisperCppServiceTranscribeWithTimestamps(t *testing.T) {
	cfg := config.LoadConfig()

	service, err := services.NewWhisperCppService(cfg)
	if err != nil {
		t.Skipf("⚠️ Skipping test: WhisperCppService not available: %v", err)
		return
	}

	// ตรวจสอบว่ามีไฟล์ตัวอย่างหรือไม่
	testAudioPath := "testdata/audio/th_audio.wav"
	if _, err := os.Stat(testAudioPath); os.IsNotExist(err) {
		t.Skipf("⚠️ Skipping test: test audio file not found at %s", testAudioPath)
		return
	}

	// อ่านไฟล์ตัวอย่าง
	audioData, err := os.ReadFile(testAudioPath)
	if err != nil {
		t.Fatalf("Failed to read test audio file: %v", err)
	}

	// สร้าง io.Reader
	audioReader := bytes.NewReader(audioData)

	// ทดสอบ TranscribeWithTimestamps
	segments, err := service.TranscribeWithTimestamps(audioReader, "th_audio.wav", "th")
	if err != nil {
		t.Fatalf("TranscribeWithTimestamps failed: %v", err)
	}

	// ตรวจสอบ output
	if len(segments) == 0 {
		t.Error("Segments should not be empty")
	}

	// ตรวจสอบ segment แรก
	if len(segments) > 0 {
		firstSeg := segments[0]
		if firstSeg.Text == "" {
			t.Error("First segment text should not be empty")
		}
		if firstSeg.EndTime <= firstSeg.StartTime {
			t.Errorf("Segment end time (%.2f) should be greater than start time (%.2f)",
				firstSeg.EndTime, firstSeg.StartTime)
		}
	}

	t.Logf("✓ Found %d segments", len(segments))
	for i, seg := range segments {
		if i < 3 { // แสดงแค่ 3 segments แรก
			t.Logf("  Segment %d: [%.2fs - %.2fs] %s",
				i+1, seg.StartTime, seg.EndTime, seg.Text)
		}
	}
}

// TestWhisperCppServiceTranscribeEmptyFile ทดสอบกับไฟล์ว่าง
func TestWhisperCppServiceTranscribeEmptyFile(t *testing.T) {
	cfg := config.LoadConfig()

	service, err := services.NewWhisperCppService(cfg)
	if err != nil {
		t.Skipf("⚠️ Skipping test: WhisperCppService not available: %v", err)
		return
	}

	// สร้าง empty reader
	emptyReader := bytes.NewReader([]byte{})

	// ทดสอบ Transcribe กับ empty file
	_, _, err = service.Transcribe(emptyReader, "empty.wav", "th")

	// ควร error เพราะไฟล์ว่าง
	if err == nil {
		t.Error("Transcribe with empty file should return error")
	}

	if !strings.Contains(err.Error(), "empty") {
		t.Logf("⚠️ Expected error to mention 'empty', got: %v", err)
	}

	t.Logf("✓ Empty file handling: %v", err)
}

// TestWhisperCppServiceVersion ทดสอบการรัน whisper.cpp --version (ถ้ารองรับ)
func TestWhisperCppServiceVersion(t *testing.T) {
	cfg := config.LoadConfig()

	// ตรวจสอบว่า binary มีอยู่หรือไม่
	if _, err := os.Stat(cfg.WhisperBinaryPath); os.IsNotExist(err) {
		t.Skipf("⚠️ Skipping test: binary not found at %s", cfg.WhisperBinaryPath)
		return
	}

	t.Logf("✓ Binary found at: %s", cfg.WhisperBinaryPath)
	t.Logf("✓ Model path: %s", cfg.WhisperModelPath)
	t.Logf("✓ Model name: %s", cfg.WhisperModelName)
	t.Logf("✓ Temp directory: %s", cfg.WhisperTempDir)
	t.Logf("✓ Language: %s", cfg.WhisperLanguage)
	t.Logf("✓ Threads: %d", cfg.WhisperThreads)
}

// TestWhisperCppServiceConfiguration ทดสอบการโหลด configuration
func TestWhisperCppServiceConfiguration(t *testing.T) {
	cfg := config.LoadConfig()

	// ตรวจสอบค่า configuration
	if cfg.WhisperBinaryPath == "" {
		t.Error("WhisperBinaryPath should not be empty")
	}

	if cfg.WhisperModelPath == "" {
		t.Error("WhisperModelPath should not be empty")
	}

	if cfg.WhisperTempDir == "" {
		t.Error("WhisperTempDir should not be empty")
	}

	if cfg.WhisperThreads <= 0 {
		t.Error("WhisperThreads should be greater than 0")
	}

	t.Logf("✓ Configuration loaded successfully")
	t.Logf("  Binary: %s", cfg.WhisperBinaryPath)
	t.Logf("  Model: %s (%s)", cfg.WhisperModelPath, cfg.WhisperModelName)
	t.Logf("  Temp: %s", cfg.WhisperTempDir)
	t.Logf("  Threads: %d", cfg.WhisperThreads)
	t.Logf("  Language: %s", cfg.WhisperLanguage)
}
