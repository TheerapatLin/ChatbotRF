package main

import (
	"fmt"
	"log"

	"chatbot/config"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("🔄 Running migration: Add session_id to file_analyses...")

	if err := addSessionIDColumn(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Println("✅ Migration completed successfully!")
}

func addSessionIDColumn(db *gorm.DB) error {
	// Check if column already exists
	var columnExists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_name = 'file_analyses'
			AND column_name = 'session_id'
		)
	`).Scan(&columnExists).Error

	if err != nil {
		return fmt.Errorf("failed to check column existence: %w", err)
	}

	if columnExists {
		fmt.Println("⚠️  Column 'session_id' already exists, skipping...")
		return nil
	}

	// Add session_id column
	if err := db.Exec(`
		ALTER TABLE file_analyses
		ADD COLUMN session_id VARCHAR(100)
	`).Error; err != nil {
		return fmt.Errorf("failed to add session_id column: %w", err)
	}

	fmt.Println("✅ Added session_id column")

	// Create index
	if err := db.Exec(`
		CREATE INDEX idx_file_analyses_session_id
		ON file_analyses (session_id)
	`).Error; err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	fmt.Println("✅ Created index idx_file_analyses_session_id")

	// Verify
	var count int64
	if err := db.Raw(`
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name = 'file_analyses'
		AND column_name = 'session_id'
	`).Scan(&count).Error; err != nil {
		return fmt.Errorf("failed to verify column: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("verification failed: column was not created")
	}

	fmt.Println("✅ Verified: session_id column exists")

	return nil
}
