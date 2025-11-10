package whisper_test

import (
	"os"
	"runtime"
	"strings"
	"testing"

	"chatbot/config"
)

// TestWhisperConfigDefaults ทดสอบค่า default ของการตั้งค่า whisper
func TestWhisperConfigDefaults(t *testing.T) {
	// ล้างค่า env vars ของ whisper เพื่อทดสอบค่า default
	os.Unsetenv("WHISPER_BINARY_PATH_LINUX")
	os.Unsetenv("WHISPER_BINARY_PATH_WINDOWS")
	os.Unsetenv("WHISPER_BINARY_PATH_MACOS")
	os.Unsetenv("WHISPER_MODEL_PATH")
	os.Unsetenv("WHISPER_LANGUAGE")
	os.Unsetenv("WHISPER_THREADS")
	os.Unsetenv("WHISPER_BEAM_SIZE")
	os.Unsetenv("WHISPER_TEMP_DIR")

	cfg := config.LoadConfig()

	// ตรวจสอบค่า default ของภาษา (ควรเป็น auto)
	if cfg.WhisperLanguage != "auto" {
		t.Errorf("คาดหวัง default language 'auto' แต่ได้ '%s'", cfg.WhisperLanguage)
	}

	if cfg.WhisperThreads != 4 {
		t.Errorf("คาดหวัง default threads 4 แต่ได้ %d", cfg.WhisperThreads)
	}

	if cfg.WhisperBeamSize != 5 {
		t.Errorf("คาดหวัง default beam size 5 แต่ได้ %d", cfg.WhisperBeamSize)
	}

	// ตรวจสอบว่า WhisperTempDir เป็น absolute path และมี "/whisper/temp" อยู่
	if !strings.Contains(cfg.WhisperTempDir, "/whisper/temp") && !strings.Contains(cfg.WhisperTempDir, "\\whisper\\temp") {
		t.Errorf("คาดหวัง temp dir ที่มี 'whisper/temp' แต่ได้ '%s'", cfg.WhisperTempDir)
	}

	if cfg.WhisperModelName != "small" {
		t.Errorf("คาดหวัง default model name 'small' แต่ได้ '%s'", cfg.WhisperModelName)
	}

	if cfg.WhisperProcessors != 1 {
		t.Errorf("คาดหวัง default processors 1 แต่ได้ %d", cfg.WhisperProcessors)
	}

	if cfg.WhisperMaxLen != 0 {
		t.Errorf("คาดหวัง default max_len 0 แต่ได้ %d", cfg.WhisperMaxLen)
	}

	if cfg.WhisperBestOf != 5 {
		t.Errorf("คาดหวัง default best_of 5 แต่ได้ %d", cfg.WhisperBestOf)
	}

	if cfg.WhisperWordTimestamps != false {
		t.Errorf("คาดหวัง default word_timestamps false แต่ได้ %v", cfg.WhisperWordTimestamps)
	}

	t.Log("✓ ค่า default ทั้งหมดของ whisper config ถูกต้อง")
}

// TestWhisperBinaryPathByOS ทดสอบการเลือก binary path ตามระบบปฏิบัติการ
func TestWhisperBinaryPathByOS(t *testing.T) {
	// ล้างค่า env vars เพื่อใช้ default path
	os.Unsetenv("WHISPER_BINARY_PATH_LINUX")
	os.Unsetenv("WHISPER_BINARY_PATH_WINDOWS")
	os.Unsetenv("WHISPER_BINARY_PATH_MACOS")

	cfg := config.LoadConfig()

	// ตรวจสอบว่า WhisperBinaryPath เป็น absolute path และมี platform-specific path
	switch runtime.GOOS {
	case "windows":
		if !strings.Contains(cfg.WhisperBinaryPath, "windows") || !strings.HasSuffix(cfg.WhisperBinaryPath, ".exe") {
			t.Errorf("คาดหวัง Windows binary path แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
		if !strings.Contains(cfg.WhisperBinaryPath, "/whisper/binary/windows/") && !strings.Contains(cfg.WhisperBinaryPath, "\\whisper\\binary\\windows\\") {
			t.Errorf("คาดหวัง path ที่มี 'whisper/binary/windows/' แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
	case "darwin":
		if !strings.Contains(cfg.WhisperBinaryPath, "macos") {
			t.Errorf("คาดหวัง macOS binary path แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
		if !strings.Contains(cfg.WhisperBinaryPath, "/whisper/binary/macos/") {
			t.Errorf("คาดหวัง path ที่มี 'whisper/binary/macos/' แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
	case "linux":
		if !strings.Contains(cfg.WhisperBinaryPath, "linux") {
			t.Errorf("คาดหวัง Linux binary path แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
		if !strings.Contains(cfg.WhisperBinaryPath, "/whisper/binary/linux/") {
			t.Errorf("คาดหวัง path ที่มี 'whisper/binary/linux/' แต่ได้ '%s'", cfg.WhisperBinaryPath)
		}
	}

	t.Logf("✓ Binary path สำหรับ %s: %s", runtime.GOOS, cfg.WhisperBinaryPath)
}

// TestWhisperConfigOverride ทดสอบการ override ด้วย environment variables
func TestWhisperConfigOverride(t *testing.T) {
	// เก็บค่าเดิมไว้
	oldLanguage := os.Getenv("WHISPER_LANGUAGE")
	oldThreads := os.Getenv("WHISPER_THREADS")
	oldBeamSize := os.Getenv("WHISPER_BEAM_SIZE")

	// ตั้งค่าแบบกำหนดเอง
	os.Setenv("WHISPER_LANGUAGE", "en")
	os.Setenv("WHISPER_THREADS", "8")
	os.Setenv("WHISPER_BEAM_SIZE", "10")

	cfg := config.LoadConfig()

	if cfg.WhisperLanguage != "en" {
		t.Errorf("คาดหวัง language override 'en' แต่ได้ '%s'", cfg.WhisperLanguage)
	}

	if cfg.WhisperThreads != 8 {
		t.Errorf("คาดหวัง threads override 8 แต่ได้ %d", cfg.WhisperThreads)
	}

	if cfg.WhisperBeamSize != 10 {
		t.Errorf("คาดหวัง beam size override 10 แต่ได้ %d", cfg.WhisperBeamSize)
	}

	// รีเซ็ตค่ากลับ
	if oldLanguage != "" {
		os.Setenv("WHISPER_LANGUAGE", oldLanguage)
	} else {
		os.Unsetenv("WHISPER_LANGUAGE")
	}
	if oldThreads != "" {
		os.Setenv("WHISPER_THREADS", oldThreads)
	} else {
		os.Unsetenv("WHISPER_THREADS")
	}
	if oldBeamSize != "" {
		os.Setenv("WHISPER_BEAM_SIZE", oldBeamSize)
	} else {
		os.Unsetenv("WHISPER_BEAM_SIZE")
	}

	t.Log("✓ การ override ด้วย environment variables ทำงานถูกต้อง")
}

// TestWhisperSupportedLanguages ทดสอบการโหลดรายการภาษาที่รองรับ
func TestWhisperSupportedLanguages(t *testing.T) {
	os.Unsetenv("WHISPER_SUPPORTED_LANGUAGES")

	cfg := config.LoadConfig()

	// ตรวจสอบว่ามีภาษาไทย, อังกฤษ, และ auto
	supportedLangs := cfg.WhisperSupportedLangs
	langs := strings.Split(supportedLangs, ",")

	foundTh := false
	foundEn := false
	foundAuto := false

	for _, lang := range langs {
		lang = strings.TrimSpace(lang)
		if lang == "th" {
			foundTh = true
		}
		if lang == "en" {
			foundEn = true
		}
		if lang == "auto" {
			foundAuto = true
		}
	}

	if !foundTh {
		t.Error("ไม่พบภาษาไทย (th) ในรายการภาษาที่รองรับ")
	}

	if !foundEn {
		t.Error("ไม่พบภาษาอังกฤษ (en) ในรายการภาษาที่รองรับ")
	}

	if !foundAuto {
		t.Error("ไม่พบการตรวจจับอัตโนมัติ (auto) ในรายการภาษาที่รองรับ")
	}

	t.Logf("✓ รองรับภาษา: %s", supportedLangs)
}

// TestWhisperBooleanConfig ทดสอบการโหลดค่า boolean configuration
func TestWhisperBooleanConfig(t *testing.T) {
	// ทดสอบค่า true ในรูปแบบต่างๆ
	testCases := []struct {
		value    string
		expected bool
		name     string
	}{
		{"true", true, "string 'true'"},
		{"1", true, "string '1'"},
		{"yes", true, "string 'yes'"},
		{"on", true, "string 'on'"},
		{"false", false, "string 'false'"},
		{"0", false, "string '0'"},
		{"no", false, "string 'no'"},
		{"off", false, "string 'off'"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("WHISPER_WORD_TIMESTAMPS", tc.value)
			cfg := config.LoadConfig()

			if cfg.WhisperWordTimestamps != tc.expected {
				t.Errorf("ค่า '%s' ควรได้ %v แต่ได้ %v", tc.value, tc.expected, cfg.WhisperWordTimestamps)
			}
		})
	}

	// รีเซ็ตค่า
	os.Unsetenv("WHISPER_WORD_TIMESTAMPS")

	t.Log("✓ การโหลดค่า boolean configuration ทำงานถูกต้อง")
}

// TestWhisperIntegerConfig ทดสอบการโหลดค่า integer configuration
func TestWhisperIntegerConfig(t *testing.T) {
	// ทดสอบค่าที่ valid
	os.Setenv("WHISPER_THREADS", "16")
	os.Setenv("WHISPER_BEAM_SIZE", "10")
	os.Setenv("WHISPER_MAX_LEN", "1000")

	cfg := config.LoadConfig()

	if cfg.WhisperThreads != 16 {
		t.Errorf("คาดหวัง threads 16 แต่ได้ %d", cfg.WhisperThreads)
	}

	if cfg.WhisperBeamSize != 10 {
		t.Errorf("คาดหวัง beam_size 10 แต่ได้ %d", cfg.WhisperBeamSize)
	}

	if cfg.WhisperMaxLen != 1000 {
		t.Errorf("คาดหวัง max_len 1000 แต่ได้ %d", cfg.WhisperMaxLen)
	}

	// ทดสอบค่าที่ไม่ valid (ควรใช้ default)
	os.Setenv("WHISPER_THREADS", "invalid")
	cfg = config.LoadConfig()

	if cfg.WhisperThreads != 4 { // default
		t.Errorf("ค่า invalid ควรใช้ default (4) แต่ได้ %d", cfg.WhisperThreads)
	}

	// รีเซ็ตค่า
	os.Unsetenv("WHISPER_THREADS")
	os.Unsetenv("WHISPER_BEAM_SIZE")
	os.Unsetenv("WHISPER_MAX_LEN")

	t.Log("✓ การโหลดค่า integer configuration ทำงานถูกต้อง")
}

// TestWhisperCustomBinaryPath ทดสอบการกำหนด custom binary path
func TestWhisperCustomBinaryPath(t *testing.T) {
	customPath := "/custom/path/to/whisper"

	switch runtime.GOOS {
	case "windows":
		os.Setenv("WHISPER_BINARY_PATH_WINDOWS", customPath)
	case "darwin":
		os.Setenv("WHISPER_BINARY_PATH_MACOS", customPath)
	case "linux":
		os.Setenv("WHISPER_BINARY_PATH_LINUX", customPath)
	}

	cfg := config.LoadConfig()

	if cfg.WhisperBinaryPath != customPath {
		t.Errorf("คาดหวัง custom path '%s' แต่ได้ '%s'", customPath, cfg.WhisperBinaryPath)
	}

	// รีเซ็ตค่า
	os.Unsetenv("WHISPER_BINARY_PATH_WINDOWS")
	os.Unsetenv("WHISPER_BINARY_PATH_MACOS")
	os.Unsetenv("WHISPER_BINARY_PATH_LINUX")

	t.Logf("✓ Custom binary path ทำงานถูกต้อง: %s", customPath)
}
