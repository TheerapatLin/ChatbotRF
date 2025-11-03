package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FileAnalysis represents a file upload record in the database
type FileAnalysis struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	FileName    string         `gorm:"type:varchar(255);not null" json:"file_name"`
	StoragePath string         `gorm:"type:varchar(500);not null" json:"storage_path"` // Path where file is stored
	MimeType    string         `gorm:"type:varchar(100);not null" json:"mime_type"`    // e.g., application/pdf, image/jpeg
	FileSize    int64          `gorm:"not null" json:"file_size"`
	UploadedAt  time.Time      `gorm:"autoCreateTime" json:"uploaded_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for FileAnalysis model
func (FileAnalysis) TableName() string {
	return "file_analyses"
}

// BeforeCreate will set a UUID rather than numeric ID
func (f *FileAnalysis) BeforeCreate(tx *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}