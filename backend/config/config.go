package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	// Server
	Port   string
	AppEnv string
	AppName string

	// Database
	DatabaseURL string

	// OpenAI
	OpenAIAPIKey     string
	OpenAIModel      string
	OpenAIMaxTokens  int
	OpenAITemperature float64

	// CORS
	CORSOrigin string

	// Rate Limiting
	RateLimitWindow string
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
	}

	// Validate required configs
	if config.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is required in .env.development")
	}

	if config.OpenAIAPIKey == "" {
		log.Fatal("OPENAI_API_KEY is required in .env.development")
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

// GetConfig returns the application configuration
func GetConfig() *Config {
	if AppConfig == nil {
		return LoadConfig()
	}
	return AppConfig
}
