package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// FileAnalysis represents a file analysis record in the database
type FileAnalysis struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SessionID       string         `gorm:"type:varchar(100);index" json:"session_id,omitempty"` // NEW - Link to conversation session
	FileName        string         `gorm:"type:varchar(255);not null" json:"filename"`
	FileType        string         `gorm:"type:varchar(100)" json:"file_type"`
	FileSize        int64          `gorm:"not null" json:"file_size"`
	FilePath        string         `gorm:"type:varchar(500)" json:"file_path,omitempty"` // Optional: store file path if saving files
	AnalysisType    string         `gorm:"type:varchar(50)" json:"analysis_type"`        // summary, detail, qa, extract
	CustomPrompt    string         `gorm:"type:text" json:"custom_prompt,omitempty"`
	Language        string         `gorm:"type:varchar(10)" json:"language"` // th, en
	Analysis        string         `gorm:"type:text" json:"analysis"`
	KeyPoints       pq.StringArray `gorm:"type:text[]" json:"key_points"`
	Entities        pq.StringArray `gorm:"type:text[]" json:"entities,omitempty"`
	Sentiment       string         `gorm:"type:varchar(50)" json:"sentiment,omitempty"`
	TokensUsed      int            `gorm:"default:0" json:"tokens_used"`
	ProcessTimeMs   float64        `gorm:"column:process_time_ms" json:"process_time_ms"`
	CreatedAt       time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	ReanalysisCount int            `gorm:"default:0" json:"reanalysis_count"` // Track how many times re-analyzed
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