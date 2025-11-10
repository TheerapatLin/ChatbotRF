package config

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port    string
	AppEnv  string
	AppName string

	// Database
	DatabaseURL string

	// OpenAI
	OpenAIAPIKey      string
	OpenAIModel       string
	OpenAIMaxTokens   int
	OpenAITemperature float64

	// CORS
	CORSOrigin string

	// Rate Limiting
	RateLimitWindow string

	// AWS Bedrock
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	AWSRegion          string
	BedrockModelID     string
	BedrockMaxTokens   int
	BedrockTemperature float64

	// ElevenLabs
	ElevenLabsAPIKey string

	// Whisper.cpp Configuration
	WhisperBinaryPath       string
	WhisperModelPath        string
	WhisperTempDir          string
	WhisperLanguage         string
	WhisperModelName        string
	WhisperThreads          int
	WhisperProcessors       int
	WhisperMaxLen           int
	WhisperBeamSize         int
	WhisperBestOf           int
	WhisperWordTimestamps   bool
	WhisperSupportedLangs   string
}

var AppConfig *Config

// LoadConfig loads configuration from .env.development file
func LoadConfig() *Config {
	// Load .env.development file
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Println("Warning: .env.development file not found, using system environment variables")
	}

	// Parse OpenAI max tokens (default: 2000)
	maxTokens, err := strconv.Atoi(getEnv("OPENAI_MAX_TOKENS", "2000"))
	if err != nil {
		maxTokens = 2000
	}

	// Parse OpenAI temperature (default: 0.7)
	temperature, err := strconv.ParseFloat(getEnv("OPENAI_TEMPERATURE", "0.7"), 64)
	if err != nil {
		temperature = 0.7
	}

	// Parse Bedrock max tokens (default: 2000)
	bedrockMaxTokens, err := strconv.Atoi(getEnv("BEDROCK_MAX_TOKENS", "2000"))
	if err != nil {
		bedrockMaxTokens = 2000
	}

	// Parse Bedrock temperature (default: 0.7)
	bedrockTemperature, err := strconv.ParseFloat(getEnv("BEDROCK_TEMPERATURE", "0.7"), 64)
	if err != nil {
		bedrockTemperature = 0.7
	}

	config := &Config{
		// Server
		Port:    getEnv("PORT", "3000"),
		AppEnv:  getEnv("APP_ENV", "development"),
		AppName: getEnv("APP_NAME", "ChatBotAPI"),

		// Database
		DatabaseURL: getEnv("DATABASE_URL", ""),

		// OpenAI
		OpenAIAPIKey:      getEnv("OPENAI_API_KEY", ""),
		OpenAIModel:       getEnv("OPENAI_MODEL", "gpt-4o-mini"),
		OpenAIMaxTokens:   maxTokens,
		OpenAITemperature: temperature,

		// CORS
		CORSOrigin: getEnv("CORS_ORIGIN", "http://localhost:5173"),

		// Rate Limiting
		RateLimitWindow: getEnv("RATE_LIMIT_WINDOW", "1m"),

		// AWS Bedrock
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		AWSRegion:          getEnv("AWS_REGION", "ap-southeast-1"),
		BedrockModelID:     getEnv("BEDROCK_MODEL_ID", "apac.anthropic.claude-sonnet-4-20250514-v1:0"),
		BedrockMaxTokens:   bedrockMaxTokens,
		BedrockTemperature: bedrockTemperature,

		// ElevenLabs
		ElevenLabsAPIKey: getEnv("ELEVENLABS_API_KEY", ""),

		// Whisper.cpp
		WhisperBinaryPath:     getWhisperBinaryPath(),
		WhisperModelPath:      getEnv("WHISPER_MODEL_PATH", "./backend/whisper/models/ggml-small.bin"),
		WhisperTempDir:        getEnv("WHISPER_TEMP_DIR", "./backend/whisper/temp"),
		WhisperLanguage:       getEnv("WHISPER_LANGUAGE", "auto"),
		WhisperModelName:      getEnv("WHISPER_MODEL_NAME", "small"),
		WhisperThreads:        getEnvAsInt("WHISPER_THREADS", 4),
		WhisperProcessors:     getEnvAsInt("WHISPER_PROCESSORS", 1),
		WhisperMaxLen:         getEnvAsInt("WHISPER_MAX_LEN", 0),
		WhisperBeamSize:       getEnvAsInt("WHISPER_BEAM_SIZE", 5),
		WhisperBestOf:         getEnvAsInt("WHISPER_BEST_OF", 5),
		WhisperWordTimestamps: getEnvAsBool("WHISPER_WORD_TIMESTAMPS", false),
		WhisperSupportedLangs: getEnv("WHISPER_SUPPORTED_LANGUAGES", "th,en,auto"),
	}

	// Validate required configs
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required in .env.development")
	}

	if config.OpenAIAPIKey == "" {
		fmt.Println("OPENAI_API_KEY is required in .env.development")
	}

	AppConfig = config
	log.Printf("✓ Configuration loaded successfully (env: %s)", config.AppEnv)

	return config
}

// getEnv retrieves environment variable or returns default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt retrieves environment variable as integer or returns default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}

// getEnvAsBool retrieves environment variable as boolean or returns default value
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	// Support multiple boolean representations
	valueStr = strings.ToLower(strings.TrimSpace(valueStr))
	switch valueStr {
	case "true", "1", "yes", "on":
		return true
	case "false", "0", "no", "off":
		return false
	default:
		log.Printf("Warning: Invalid boolean value for %s: %s, using default: %v", key, valueStr, defaultValue)
		return defaultValue
	}
}

// getWhisperBinaryPath returns the correct Whisper binary path based on the OS
func getWhisperBinaryPath() string {
	var envKey string
	switch runtime.GOOS {
	case "windows":
		envKey = "WHISPER_BINARY_PATH_WINDOWS"
	case "darwin":
		envKey = "WHISPER_BINARY_PATH_MACOS"
	case "linux":
		envKey = "WHISPER_BINARY_PATH_LINUX"
	default:
		envKey = "WHISPER_BINARY_PATH_LINUX"
	}

	defaultPath := "./backend/whisper/binary/linux/main"
	if runtime.GOOS == "windows" {
		defaultPath = "./backend/whisper/binary/windows/main.exe"
	} else if runtime.GOOS == "darwin" {
		defaultPath = "./backend/whisper/binary/macos/main"
	}

	return getEnv(envKey, defaultPath)
}

// GetConfig returns the application configuration
func GetConfig() *Config {
	if AppConfig == nil {
		return LoadConfig()
	}
	return AppConfig
}
