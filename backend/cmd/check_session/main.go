package main

import (
	"fmt"
	"log"

	"chatbot/config"
	"chatbot/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	fmt.Println("🔍 Checking file session data...")

	// Load environment and configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	// File IDs to check
	fileIDs := []string{
		"a3f95c98-884f-446f-aaef-26096ffa669d", // From latest upload
		"66c2abbe-4ce1-43a6-b251-06fa66d4dd804", // Used in chat request
	}

	fmt.Println("\n📋 Checking files:")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")

	for _, fileIDStr := range fileIDs {
		fileID, err := uuid.Parse(fileIDStr)
		if err != nil {
			fmt.Printf("❌ Invalid UUID: %s - %v\n", fileIDStr, err)
			continue
		}

		var analysis models.FileAnalysis
		err = db.Where("id = ?", fileID).First(&analysis).Error

		if err == gorm.ErrRecordNotFound {
			fmt.Printf("\n❌ File NOT FOUND: %s\n", fileIDStr)
			continue
		}

		if err != nil {
			fmt.Printf("\n❌ Error querying file %s: %v\n", fileIDStr, err)
			continue
		}

		// Display results
		fmt.Printf("\n✅ File Found: %s\n", fileIDStr)
		fmt.Printf("   📄 Filename: %s\n", analysis.FileName)
		fmt.Printf("   📅 Created: %s\n", analysis.CreatedAt.Format("2006-01-02 15:04:05"))

		if analysis.SessionID != "" {
			fmt.Printf("   🔗 Session ID: %s ✅\n", analysis.SessionID)
		} else {
			fmt.Printf("   ❌ Session ID: EMPTY/NULL\n")
		}

		// Show summary
		if len(analysis.Summary) > 100 {
			fmt.Printf("   📝 Summary: %s...\n", analysis.Summary[:100])
		} else {
			fmt.Printf("   📝 Summary: %s\n", analysis.Summary)
		}
	}

	fmt.Println("\n─────────────────────────────────────────────────────────────────────────────")

	// Also check all recent files with session_id = test_memory_001
	fmt.Println("\n🔍 All files with session 'test_memory_001':")
	fmt.Println("─────────────────────────────────────────────────────────────────────────────")

	var sessionFiles []models.FileAnalysis
	err = db.Where("session_id = ?", "test_memory_001").
		Order("created_at DESC").
		Find(&sessionFiles).Error

	if err != nil {
		fmt.Printf("❌ Error querying session files: %v\n", err)
	} else if len(sessionFiles) == 0 {
		fmt.Println("❌ No files found for session 'test_memory_001'")
	} else {
		fmt.Printf("✅ Found %d file(s) for session 'test_memory_001':\n\n", len(sessionFiles))
		for i, f := range sessionFiles {
			fmt.Printf("%d. ID: %s\n", i+1, f.ID.String())
			fmt.Printf("   📄 Filename: %s\n", f.FileName)
			fmt.Printf("   📅 Created: %s\n\n", f.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	fmt.Println("✅ Check complete!")
}
