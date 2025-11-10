package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"chatbot/config"
)

// WhisperCppService implements TranscriptionService interface using whisper.cpp
// ใช้ whisper.cpp binary เพื่อแปลงเสียงเป็นข้อความแบบ offline
type WhisperCppService struct {
	config *config.Config
}

// NewWhisperCppService creates a new WhisperCppService instance
// ตรวจสอบว่า binary และ model พร้อมใช้งาน และสร้าง temp directory
func NewWhisperCppService(cfg *config.Config) (*WhisperCppService, error) {
	service := &WhisperCppService{
		config: cfg,
	}

	// Validate binary exists
	if _, err := os.Stat(cfg.WhisperBinaryPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("whisper.cpp binary not found at: %s", cfg.WhisperBinaryPath)
	}

	// Validate model exists
	if _, err := os.Stat(cfg.WhisperModelPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("whisper model not found at: %s", cfg.WhisperModelPath)
	}

	// Ensure temp directory exists
	if err := os.MkdirAll(cfg.WhisperTempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	fmt.Printf("✓ WhisperCppService initialized (binary: %s, model: %s)\n",
		filepath.Base(cfg.WhisperBinaryPath), cfg.WhisperModelName)

	return service, nil
}

// IsAvailable ตรวจสอบว่า service พร้อมใช้งาน
// ตรวจสอบ binary file, model file, และ temp directory
func (s *WhisperCppService) IsAvailable() bool {
	// ตรวจสอบ binary
	if _, err := os.Stat(s.config.WhisperBinaryPath); os.IsNotExist(err) {
		return false
	}

	// ตรวจสอบ model
	if _, err := os.Stat(s.config.WhisperModelPath); os.IsNotExist(err) {
		return false
	}

	// ตรวจสอบว่าสามารถสร้าง temp directory ได้
	if err := os.MkdirAll(s.config.WhisperTempDir, 0755); err != nil {
		return false
	}

	return true
}

// GetSupportedFormats คืนรายการรูปแบบ audio ที่ whisper.cpp รองรับ
func (s *WhisperCppService) GetSupportedFormats() []string {
	return []string{"wav", "mp3", "m4a", "ogg", "flac", "opus", "webm"}
}

// GetModelName คืนชื่อโมเดลที่ใช้งานอยู่
func (s *WhisperCppService) GetModelName() string {
	return s.config.WhisperModelName
}

// GetSupportedModels คืนรายการ models ที่รองรับ
func (s *WhisperCppService) GetSupportedModels() []string {
	models := strings.Split(s.config.WhisperSupportedModels, ",")
	result := make([]string, 0, len(models))
	for _, model := range models {
		trimmed := strings.TrimSpace(model)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// GetModelPath คืน path ของ model ที่ระบุ
// ถ้า modelName ว่าง จะใช้ default model
// ถ้า modelName มีค่า จะหา model file ใน models directory
func (s *WhisperCppService) GetModelPath(modelName string) (string, error) {
	// ถ้าไม่ระบุ model ใช้ default
	if modelName == "" {
		return s.config.WhisperModelPath, nil
	}

	// ตรวจสอบว่า model ที่ระบุอยู่ใน supported list หรือไม่
	supportedModels := s.GetSupportedModels()
	isSupported := false
	for _, supported := range supportedModels {
		if strings.EqualFold(supported, modelName) {
			isSupported = true
			break
		}
	}
	if !isSupported {
		return "", fmt.Errorf("model '%s' is not supported. Supported models: %s",
			modelName, strings.Join(supportedModels, ", "))
	}

	// สร้าง filename จาก model name
	// Format: ggml-{modelName}.bin หรือ ggml-{modelName}-q5_1.bin
	var modelFilename string
	if strings.Contains(modelName, ".") {
		// ถ้ามี extension แล้ว (เช่น tiny.en) ให้เติม .bin
		modelFilename = fmt.Sprintf("ggml-%s.bin", modelName)
	} else {
		modelFilename = fmt.Sprintf("ggml-%s.bin", modelName)
	}

	// Check ใน models directory
	modelPath := filepath.Join(s.config.WhisperModelsDir, modelFilename)
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		// ลอง format อื่น: ggml-{modelName}-q5_1.bin
		altFilename := fmt.Sprintf("ggml-%s-q5_1.bin", strings.ReplaceAll(modelName, ".", "-"))
		altPath := filepath.Join(s.config.WhisperModelsDir, altFilename)
		if _, err := os.Stat(altPath); os.IsNotExist(err) {
			return "", fmt.Errorf("model file not found: tried %s and %s", modelFilename, altFilename)
		}
		modelPath = altPath
	}

	return modelPath, nil
}

// Transcribe แปลงไฟล์ audio เป็นข้อความ
// Parameters:
//   - audioFile: io.Reader ที่มีข้อมูล audio
//   - filename: ชื่อไฟล์เดิม (ใช้เพื่อกำหนดรูปแบบ)
//   - language: รหัสภาษา (เช่น "th", "en", "auto")
//
// Returns:
//   - transcription: ข้อความที่แปลงได้
//   - confidence: คะแนนความมั่นใจ (0.0 - 1.0)
//   - error: ข้อผิดพลาดที่เกิดขึ้น
func (s *WhisperCppService) Transcribe(audioFile io.Reader, filename string, language string) (string, float64, error) {
	startTime := time.Now()

	// Default language ถ้าไม่ระบุ
	if language == "" {
		language = s.config.WhisperLanguage
	}

	fmt.Printf("🔄 Starting whisper.cpp transcription (language=%s, model=%s)\n",
		language, s.config.WhisperModelName)

	// 1. Save audio to temp file
	tempFilePath, err := s.saveTempFile(audioFile)
	if err != nil {
		return "", 0.0, fmt.Errorf("failed to save temp file: %w", err)
	}
	defer s.cleanupTempFile(tempFilePath)

	// 2. Build command arguments
	args := s.buildWhisperArgs(tempFilePath, language, false)

	// 3. Execute whisper.cpp
	// ใช้ timeout ที่สั้นลง (1 นาที) เพื่อให้ทำงานกับ test timeout ได้
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	output, err := s.executeWhisper(ctx, args)
	if err != nil {
		fmt.Printf("❌ Transcription failed: %v\n", err)
		return "", 0.0, err
	}

	// 4. Parse output
	transcription := s.parseTextOutput(output)

	// 5. Calculate confidence
	confidence := s.calculateConfidence(transcription)

	duration := time.Since(startTime)
	fmt.Printf("✅ Transcription completed in %.2fs (confidence: %.2f)\n",
		duration.Seconds(), confidence)

	return transcription, confidence, nil
}

// TranscribeWithTimestamps แปลงไฟล์ audio เป็นข้อความพร้อม timestamps
// Parameters:
//   - audioFile: io.Reader ที่มีข้อมูล audio
//   - filename: ชื่อไฟล์เดิม
//   - language: รหัสภาษา
//
// Returns:
//   - segments: array ของ segments พร้อม timestamps
//   - error: ข้อผิดพลาดที่เกิดขึ้น
func (s *WhisperCppService) TranscribeWithTimestamps(audioFile io.Reader, filename string, language string) ([]TranscriptionSegment, error) {
	startTime := time.Now()

	// Default language ถ้าไม่ระบุ
	if language == "" {
		language = s.config.WhisperLanguage
	}

	fmt.Printf("🔄 Starting whisper.cpp transcription with timestamps (language=%s, model=%s)\n",
		language, s.config.WhisperModelName)

	// 1. Save audio to temp file
	tempFilePath, err := s.saveTempFile(audioFile)
	if err != nil {
		return nil, fmt.Errorf("failed to save temp file: %w", err)
	}
	defer s.cleanupTempFile(tempFilePath)

	// 2. Build command arguments (with timestamps enabled)
	args := s.buildWhisperArgs(tempFilePath, language, true)

	// 3. Execute whisper.cpp
	// ใช้ timeout ที่สั้นลง (1 นาที) เพื่อให้ทำงานกับ test timeout ได้
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err = s.executeWhisper(ctx, args)
	if err != nil {
		fmt.Printf("❌ Transcription with timestamps failed: %v\n", err)
		return nil, err
	}

	// 4. Read JSON output from file (whisper.cpp saves JSON to <input>.json)
	jsonFilePath := tempFilePath + ".json"
	defer os.Remove(jsonFilePath) // cleanup JSON file

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON output file: %w", err)
	}

	// 5. Parse JSON output
	segments, err := s.parseJSONOutput(string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("✅ Transcription with timestamps completed in %.2fs (%d segments)\n",
		duration.Seconds(), len(segments))

	return segments, nil
}

// TranscribeWithModel แปลงไฟล์ audio เป็นข้อความโดยระบุ model
// Parameters:
//   - audioFile: io.Reader ที่มีข้อมูล audio
//   - filename: ชื่อไฟล์เดิม
//   - language: รหัสภาษา
//   - modelName: ชื่อ model (เช่น "tiny.en", "small", "medium")
//
// Returns:
//   - transcription: ข้อความที่แปลงได้
//   - confidence: คะแนนความมั่นใจ
//   - error: ข้อผิดพลาด
func (s *WhisperCppService) TranscribeWithModel(audioFile io.Reader, filename string, language string, modelName string) (string, float64, error) {
	startTime := time.Now()

	// Default language ถ้าไม่ระบุ
	if language == "" {
		language = s.config.WhisperLanguage
	}

	// Get model path
	modelPath, err := s.GetModelPath(modelName)
	if err != nil {
		return "", 0.0, fmt.Errorf("model selection error: %w", err)
	}

	actualModelName := modelName
	if actualModelName == "" {
		actualModelName = s.config.WhisperModelName
	}

	fmt.Printf("🔄 Starting whisper.cpp transcription (language=%s, model=%s)\n",
		language, actualModelName)

	// 1. Save audio to temp file
	tempFilePath, err := s.saveTempFile(audioFile)
	if err != nil {
		return "", 0.0, fmt.Errorf("failed to save temp file: %w", err)
	}
	defer s.cleanupTempFile(tempFilePath)

	// 2. Build command arguments with custom model
	args := s.buildWhisperArgsWithModel(tempFilePath, language, false, modelPath)

	// 3. Execute whisper.cpp
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	output, err := s.executeWhisper(ctx, args)
	if err != nil {
		fmt.Printf("❌ Transcription failed: %v\n", err)
		return "", 0.0, err
	}

	// 4. Parse output
	transcription := s.parseTextOutput(output)

	// 5. Calculate confidence
	confidence := s.calculateConfidence(transcription)

	duration := time.Since(startTime)
	fmt.Printf("✅ Transcription completed in %.2fs (confidence: %.2f, model: %s)\n",
		duration.Seconds(), confidence, actualModelName)

	return transcription, confidence, nil
}

// TranscribeWithTimestampsAndModel แปลงไฟล์ audio เป็นข้อความพร้อม timestamps โดยระบุ model
// Parameters:
//   - audioFile: io.Reader ที่มีข้อมูล audio
//   - filename: ชื่อไฟล์เดิม
//   - language: รหัสภาษา
//   - modelName: ชื่อ model
//
// Returns:
//   - segments: array ของ segments พร้อม timestamps
//   - error: ข้อผิดพลาด
func (s *WhisperCppService) TranscribeWithTimestampsAndModel(audioFile io.Reader, filename string, language string, modelName string) ([]TranscriptionSegment, error) {
	startTime := time.Now()

	// Default language ถ้าไม่ระบุ
	if language == "" {
		language = s.config.WhisperLanguage
	}

	// Get model path
	modelPath, err := s.GetModelPath(modelName)
	if err != nil {
		return nil, fmt.Errorf("model selection error: %w", err)
	}

	actualModelName := modelName
	if actualModelName == "" {
		actualModelName = s.config.WhisperModelName
	}

	fmt.Printf("🔄 Starting whisper.cpp transcription with timestamps (language=%s, model=%s)\n",
		language, actualModelName)

	// 1. Save audio to temp file
	tempFilePath, err := s.saveTempFile(audioFile)
	if err != nil {
		return nil, fmt.Errorf("failed to save temp file: %w", err)
	}
	defer s.cleanupTempFile(tempFilePath)

	// 2. Build command arguments with custom model (with timestamps enabled)
	args := s.buildWhisperArgsWithModel(tempFilePath, language, true, modelPath)

	// 3. Execute whisper.cpp
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	_, err = s.executeWhisper(ctx, args)
	if err != nil {
		fmt.Printf("❌ Transcription with timestamps failed: %v\n", err)
		return nil, err
	}

	// 4. Read JSON output from file
	jsonFilePath := tempFilePath + ".json"
	defer os.Remove(jsonFilePath)

	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON output file: %w", err)
	}

	// 5. Parse JSON output
	segments, err := s.parseJSONOutput(string(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON output: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("✅ Transcription with timestamps completed in %.2fs (%d segments, model: %s)\n",
		duration.Seconds(), len(segments), actualModelName)

	return segments, nil
}

// ========================================
// Helper Functions
// ========================================

// saveTempFile บันทึก audio data จาก io.Reader ไปยัง temp file
func (s *WhisperCppService) saveTempFile(audioFile io.Reader) (string, error) {
	// สร้าง temp file
	tempFile, err := os.CreateTemp(s.config.WhisperTempDir, "whisper-audio-*.wav")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	// Copy audio data
	written, err := io.Copy(tempFile, audioFile)
	if err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to write audio data: %w", err)
	}

	if written == 0 {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("audio file is empty")
	}

	return tempFile.Name(), nil
}

// buildWhisperArgs สร้าง command-line arguments สำหรับ whisper.cpp
func (s *WhisperCppService) buildWhisperArgs(audioPath string, language string, withTimestamps bool) []string {
	args := []string{
		"-m", s.config.WhisperModelPath, // model file
		"-f", audioPath, // audio file
		"-t", fmt.Sprintf("%d", s.config.WhisperThreads), // threads
		"-l", language, // language (th, en, auto)
	}

	// เพิ่ม beam size ถ้ากำหนด
	if s.config.WhisperBeamSize > 0 {
		args = append(args, "-bs", fmt.Sprintf("%d", s.config.WhisperBeamSize))
	}

	// เพิ่ม best_of ถ้ากำหนด
	if s.config.WhisperBestOf > 0 {
		args = append(args, "-bo", fmt.Sprintf("%d", s.config.WhisperBestOf))
	}

	// เพิ่ม processors ถ้ากำหนด
	if s.config.WhisperProcessors > 1 {
		args = append(args, "-p", fmt.Sprintf("%d", s.config.WhisperProcessors))
	}

	// เพิ่ม max_len ถ้ากำหนด
	if s.config.WhisperMaxLen > 0 {
		args = append(args, "-ml", fmt.Sprintf("%d", s.config.WhisperMaxLen))
	}

	// ถ้าต้องการ timestamps, ใช้ JSON output
	if withTimestamps {
		args = append(args, "-oj") // output as JSON
		if s.config.WhisperWordTimestamps {
			args = append(args, "-ml", "1") // max line length = 1 for word-level timestamps
		}
	} else {
		args = append(args, "-nt") // no timestamps in text output
	}

	return args
}

// buildWhisperArgsWithModel สร้าง arguments สำหรับ whisper.cpp โดยระบุ model path
func (s *WhisperCppService) buildWhisperArgsWithModel(audioPath string, language string, withTimestamps bool, modelPath string) []string {
	args := []string{
		"-m", modelPath, // custom model file
		"-f", audioPath, // audio file
		"-t", fmt.Sprintf("%d", s.config.WhisperThreads), // threads
		"-l", language, // language (th, en, auto)
	}

	// เพิ่ม beam size ถ้ากำหนด
	if s.config.WhisperBeamSize > 0 {
		args = append(args, "-bs", fmt.Sprintf("%d", s.config.WhisperBeamSize))
	}

	// เพิ่ม best_of ถ้ากำหนด
	if s.config.WhisperBestOf > 0 {
		args = append(args, "-bo", fmt.Sprintf("%d", s.config.WhisperBestOf))
	}

	// เพิ่ม processors ถ้ากำหนด
	if s.config.WhisperProcessors > 1 {
		args = append(args, "-p", fmt.Sprintf("%d", s.config.WhisperProcessors))
	}

	// เพิ่ม max_len ถ้ากำหนด
	if s.config.WhisperMaxLen > 0 {
		args = append(args, "-ml", fmt.Sprintf("%d", s.config.WhisperMaxLen))
	}

	// ถ้าต้องการ timestamps, ใช้ JSON output
	if withTimestamps {
		args = append(args, "-oj") // output as JSON
		if s.config.WhisperWordTimestamps {
			args = append(args, "-ml", "1") // max line length = 1 for word-level timestamps
		}
	} else {
		args = append(args, "-nt") // no timestamps in text output
	}

	return args
}

// executeWhisper รัน whisper.cpp binary และ return output
func (s *WhisperCppService) executeWhisper(ctx context.Context, args []string) (string, error) {
	cmd := exec.CommandContext(ctx, s.config.WhisperBinaryPath, args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// รัน command
	err := cmd.Run()

	// ตรวจสอบ timeout
	if ctx.Err() == context.DeadlineExceeded {
		return "", fmt.Errorf("whisper.cpp execution timeout after 1 minute")
	}

	// ตรวจสอบ error จากการรัน command
	if err != nil {
		stderrStr := stderr.String()
		return "", fmt.Errorf("whisper.cpp execution failed: %w, stderr: %s", err, stderrStr)
	}

	return stdout.String(), nil
}

// parseTextOutput parse plain text output จาก whisper.cpp
func (s *WhisperCppService) parseTextOutput(output string) string {
	// แยกบรรทัดและหา transcription
	lines := strings.Split(output, "\n")

	var transcription strings.Builder
	foundTranscription := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// ข้าม empty lines และ log lines
		if line == "" || strings.HasPrefix(line, "whisper_") ||
			strings.HasPrefix(line, "system_info") ||
			strings.Contains(line, "processing") {
			continue
		}

		// เริ่มต้น transcription หลังจากเจอ audio info
		if strings.Contains(line, "[") && strings.Contains(line, "-->") {
			foundTranscription = true
			// Extract text after timestamp (format: [00:00.000 --> 00:02.000] text)
			parts := strings.SplitN(line, "]", 2)
			if len(parts) == 2 {
				text := strings.TrimSpace(parts[1])
				if text != "" {
					if transcription.Len() > 0 {
						transcription.WriteString(" ")
					}
					transcription.WriteString(text)
				}
			}
			continue
		}

		// ถ้าเจอ transcription แล้ว เก็บ text
		if foundTranscription && !strings.HasPrefix(line, "[") {
			if transcription.Len() > 0 {
				transcription.WriteString(" ")
			}
			transcription.WriteString(line)
		}
	}

	result := strings.TrimSpace(transcription.String())

	// ถ้าไม่เจอ transcription จาก timestamp format, ลองเอาทั้ง output
	if result == "" {
		result = strings.TrimSpace(output)
	}

	return result
}

// parseJSONOutput parse JSON output จาก whisper.cpp
func (s *WhisperCppService) parseJSONOutput(output string) ([]TranscriptionSegment, error) {
	// whisper.cpp JSON format:
	// {
	//   "systeminfo": "...",
	//   "model": {...},
	//   "params": {...},
	//   "result": {...},
	//   "transcription": [
	//     {"timestamps": {...}, "offsets": {...}, "text": "..."},
	//     ...
	//   ]
	// }

	type WhisperTimestamps struct {
		From string `json:"from"`
		To   string `json:"to"`
	}

	type WhisperTranscriptionItem struct {
		Timestamps WhisperTimestamps `json:"timestamps"`
		Text       string            `json:"text"`
	}

	type WhisperJSONOutput struct {
		Transcription []WhisperTranscriptionItem `json:"transcription"`
	}

	// Parse JSON
	var whisperOutput WhisperJSONOutput
	if err := json.Unmarshal([]byte(output), &whisperOutput); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Convert to TranscriptionSegment
	segments := make([]TranscriptionSegment, 0, len(whisperOutput.Transcription))

	for _, item := range whisperOutput.Transcription {
		text := strings.TrimSpace(item.Text)
		if text == "" {
			continue
		}

		// Parse timestamps (format: "00:00:00,000")
		startTime := parseTimestamp(item.Timestamps.From)
		endTime := parseTimestamp(item.Timestamps.To)

		segments = append(segments, TranscriptionSegment{
			StartTime: startTime,
			EndTime:   endTime,
			Text:      text,
		})
	}

	return segments, nil
}

// parseTimestamp แปลง timestamp string เป็นวินาที
// Format: "00:00:00,000" หรือ "00:00:00.000"
func parseTimestamp(timestamp string) float64 {
	// Replace comma with dot
	timestamp = strings.ReplaceAll(timestamp, ",", ".")

	// Parse HH:MM:SS.mmm
	parts := strings.Split(timestamp, ":")
	if len(parts) != 3 {
		return 0.0
	}

	var hours, minutes float64
	var seconds float64

	fmt.Sscanf(parts[0], "%f", &hours)
	fmt.Sscanf(parts[1], "%f", &minutes)
	fmt.Sscanf(parts[2], "%f", &seconds)

	return hours*3600 + minutes*60 + seconds
}

// calculateConfidence ประมาณค่า confidence score จากคุณภาพ transcription
func (s *WhisperCppService) calculateConfidence(transcription string) float64 {
	transcription = strings.TrimSpace(transcription)

	// ถ้าว่างเปล่า → confidence ต่ำมาก
	if len(transcription) == 0 {
		return 0.1
	}

	// ถ้ามี [BLANK_AUDIO] หรือ error markers → confidence ต่ำ
	lowercased := strings.ToLower(transcription)
	if strings.Contains(lowercased, "[blank") ||
		strings.Contains(lowercased, "[inaudible]") ||
		strings.Contains(lowercased, "...") && len(transcription) < 10 {
		return 0.3
	}

	// คำนวณจำนวนคำ
	wordCount := len(strings.Fields(transcription))

	// Text ที่มีความยาวปกติ → confidence สูง
	if wordCount > 10 {
		return 0.92
	} else if wordCount > 5 {
		return 0.88
	} else if wordCount > 2 {
		return 0.82
	} else if wordCount > 0 {
		return 0.75
	}

	return 0.65
}

// cleanupTempFile ลบ temp file
func (s *WhisperCppService) cleanupTempFile(filePath string) {
	if err := os.Remove(filePath); err != nil {
		// Log warning but don't fail
		fmt.Printf("⚠️ Failed to cleanup temp file %s: %v\n", filePath, err)
	}
}
