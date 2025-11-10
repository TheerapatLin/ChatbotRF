package whisper_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

// getBinaryPath คืนค่า path ของ binary ตามระบบปฏิบัติการ
func getBinaryPath() string {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join("..", "..", "whisper", "binary", "windows", "main.exe")
	case "darwin":
		return filepath.Join("..", "..", "whisper", "binary", "macos", "main")
	case "linux":
		return filepath.Join("..", "..", "whisper", "binary", "linux", "main")
	default:
		return filepath.Join("..", "..", "whisper", "binary", "linux", "main")
	}
}

// getModelPath คืนค่า path ของโมเดล
func getModelPath() string {
	return filepath.Join("..", "..", "whisper", "models", "ggml-small.bin")
}

// TestWhisperBinaryExists ตรวจสอบว่ามี whisper.cpp binary สำหรับระบบปฏิบัติการปัจจุบัน
func TestWhisperBinaryExists(t *testing.T) {
	binaryPath := getBinaryPath()

	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Fatalf("ไม่พบ whisper.cpp binary ที่ %s (OS: %s)", binaryPath, runtime.GOOS)
	}

	t.Logf("✓ พบ whisper.cpp binary ที่ %s (OS: %s)", binaryPath, runtime.GOOS)
}

// TestWhisperModelExists ตรวจสอบว่ามีไฟล์โมเดล GGML
func TestWhisperModelExists(t *testing.T) {
	modelPath := getModelPath()

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		t.Fatalf("ไม่พบโมเดล whisper ที่ %s", modelPath)
	}

	info, _ := os.Stat(modelPath)
	t.Logf("✓ พบโมเดล whisper ที่ %s (ขนาด: %d MB)", modelPath, info.Size()/1024/1024)
}

// TestWhisperVersion ตรวจสอบว่า whisper.cpp สามารถทำงานได้
func TestWhisperVersion(t *testing.T) {
	binaryPath := getBinaryPath()

	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ไม่สามารถรัน whisper.cpp: %v (OS: %s)", err, runtime.GOOS)
	}

	if len(output) == 0 {
		t.Fatal("whisper.cpp ไม่มี output")
	}

	t.Logf("✓ whisper.cpp ทำงานได้ปกติ (OS: %s)", runtime.GOOS)
}

// TestWhisperTranscribeThaiAudio ทดสอบการแปลงเสียงภาษาไทย
func TestWhisperTranscribeThaiAudio(t *testing.T) {
	binaryPath := getBinaryPath()
	modelPath := getModelPath()
	audioPath := filepath.Join("testdata", "audio", "th_audio.wav")

	// ข้าม test ถ้าไม่มีไฟล์เสียงทดสอบ
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		t.Skipf("ข้าม test: ไม่พบไฟล์เสียงทดสอบ %s", audioPath)
	}

	cmd := exec.Command(binaryPath, "-m", modelPath, "-f", audioPath, "-l", "th")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ไม่สามารถแปลงเสียงภาษาไทย: %v\nOutput: %s", err, string(output))
	}

	t.Logf("✓ แปลงเสียงภาษาไทยสำเร็จ\nOutput:\n%s", string(output))
}

// TestWhisperTranscribeEnglishAudio ทดสอบการแปลงเสียงภาษาอังกฤษ
func TestWhisperTranscribeEnglishAudio(t *testing.T) {
	binaryPath := getBinaryPath()
	modelPath := getModelPath()
	audioPath := filepath.Join("testdata", "audio", "en_audio.mp3")

	// ข้าม test ถ้าไม่มีไฟล์เสียงทดสอบ
	if _, err := os.Stat(audioPath); os.IsNotExist(err) {
		t.Skipf("ข้าม test: ไม่พบไฟล์เสียงทดสอบ %s", audioPath)
	}

	cmd := exec.Command(binaryPath, "-m", modelPath, "-f", audioPath, "-l", "en")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("ไม่สามารถแปลงเสียงภาษาอังกฤษ: %v\nOutput: %s", err, string(output))
	}

	t.Logf("✓ แปลงเสียงภาษาอังกฤษสำเร็จ\nOutput:\n%s", string(output))
}
