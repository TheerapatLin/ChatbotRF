package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// FileAttachment represents a file attached to a message
type FileAttachment struct {
	FileID   string `json:"file_id"`
	Filename string `json:"filename"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}

// Message represents a chat message (simplified for learning project)
type Message struct {
	ID              uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SessionID       string         `gorm:"type:varchar(100);index" json:"session_id"` // Session identifier for grouping conversations
	Role            string         `gorm:"type:varchar(20);not null;check:role IN ('user', 'assistant', 'system')" json:"role"`
	Content         string         `gorm:"type:text;not null" json:"content"`
	PersonaID       *int           `json:"persona_id,omitempty"`
	TokensUsed      *int           `json:"tokens_used,omitempty"`
	FileAttachments datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"file_attachments"` // NEW - Array of FileAttachment
	CreatedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Metadata        datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`

	// Relationships
	Persona *Persona `gorm:"foreignKey:PersonaID" json:"persona,omitempty"`
}

// TableName specifies the table name for Message model
func (Message) TableName() string {
	return "messages"
}

// BeforeCreate hook to generate UUID before creating
func (m *Message) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}

// GetFileAttachments parses file_attachments JSONB into FileAttachment slice
func (m *Message) GetFileAttachments() ([]FileAttachment, error) {
	var attachments []FileAttachment

	if len(m.FileAttachments) == 0 {
		return attachments, nil
	}

	if err := json.Unmarshal(m.FileAttachments, &attachments); err != nil {
		return nil, err
	}

	return attachments, nil
}

// SetFileAttachments sets file attachments from slice
func (m *Message) SetFileAttachments(attachments []FileAttachment) error {
	if attachments == nil {
		attachments = []FileAttachment{}
	}

	data, err := json.Marshal(attachments)
	if err != nil {
		return err
	}

	m.FileAttachments = data
	return nil
}

// HasFileAttachments checks if message has any file attachments
func (m *Message) HasFileAttachments() bool {
	attachments, err := m.GetFileAttachments()
	if err != nil {
		return false
	}
	return len(attachments) > 0
}

// MessageRole constants
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)
