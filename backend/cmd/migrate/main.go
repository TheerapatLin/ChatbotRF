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

	fmt.Println("🔄 Running migration: Add file_attachments column...")

	// Run migration
	if err := addFileAttachmentsColumn(db); err != nil {
		log.Fatalf("❌ Migration failed: %v", err)
	}

	fmt.Println("✅ Migration completed successfully!")
}

func addFileAttachmentsColumn(db *gorm.DB) error {
	// Check if column already exists
	var exists bool
	err := db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.columns
			WHERE table_name = 'messages'
			AND column_name = 'file_attachments'
		)
	`).Scan(&exists).Error

	if err != nil {
		return fmt.Errorf("failed to check column existence: %w", err)
	}

	if exists {
		fmt.Println("⚠️  Column 'file_attachments' already exists. Skipping...")
		return nil
	}

	// Add file_attachments column
	if err := db.Exec(`
		ALTER TABLE messages
		ADD COLUMN file_attachments JSONB DEFAULT '[]'::jsonb
	`).Error; err != nil {
		return fmt.Errorf("failed to add column: %w", err)
	}

	fmt.Println("✓ Added file_attachments column")

	// Create GIN index
	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_messages_file_attachments
		ON messages USING GIN (file_attachments)
	`).Error; err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	fmt.Println("✓ Created GIN index")

	// Add comment
	if err := db.Exec(`
		COMMENT ON COLUMN messages.file_attachments IS
		'Array of file attachment objects with structure: [{"file_id": "uuid", "filename": "example.pdf", "file_type": "application/pdf", "file_size": 123456, "analysis_summary": "Brief summary"}]'
	`).Error; err != nil {
		return fmt.Errorf("failed to add comment: %w", err)
	}

	fmt.Println("✓ Added column comment")

	// Verify the change
	var columnInfo struct {
		ColumnName    string
		DataType      string
		ColumnDefault *string
	}

	if err := db.Raw(`
		SELECT column_name, data_type, column_default
		FROM information_schema.columns
		WHERE table_name = 'messages'
		AND column_name = 'file_attachments'
	`).Scan(&columnInfo).Error; err != nil {
		return fmt.Errorf("failed to verify column: %w", err)
	}

	fmt.Printf("\n📋 Column Info:\n")
	fmt.Printf("  Name: %s\n", columnInfo.ColumnName)
	fmt.Printf("  Type: %s\n", columnInfo.DataType)
	if columnInfo.ColumnDefault != nil {
		fmt.Printf("  Default: %s\n", *columnInfo.ColumnDefault)
	}

	return nil
}
