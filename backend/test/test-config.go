package main

import (
	"fmt"
	"log"
	"runtime"

	"chatbot/config"
)

func main() {
	fmt.Println("=== Testing Whisper Configuration ===")
	fmt.Printf("Current OS: %s\n\n", runtime.GOOS)

	// Load configuration
	cfg := config.LoadConfig()

	// Display Whisper configuration
	fmt.Println("Whisper.cpp Configuration:")
	fmt.Printf("  Binary Path:        %s\n", cfg.WhisperBinaryPath)
	fmt.Printf("  Model Path:         %s\n", cfg.WhisperModelPath)
	fmt.Printf("  Temp Directory:     %s\n", cfg.WhisperTempDir)
	fmt.Printf("  Language:           %s\n", cfg.WhisperLanguage)
	fmt.Printf("  Model Name:         %s\n", cfg.WhisperModelName)
	fmt.Printf("  Threads:            %d\n", cfg.WhisperThreads)
	fmt.Printf("  Processors:         %d\n", cfg.WhisperProcessors)
	fmt.Printf("  Max Length:         %d\n", cfg.WhisperMaxLen)
	fmt.Printf("  Beam Size:          %d\n", cfg.WhisperBeamSize)
	fmt.Printf("  Best Of:            %d\n", cfg.WhisperBestOf)
	fmt.Printf("  Word Timestamps:    %v\n", cfg.WhisperWordTimestamps)
	fmt.Printf("  Supported Languages: %s\n", cfg.WhisperSupportedLangs)

	log.Println("\n✓ Configuration test completed successfully")
}
