package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Message represents a chat message (simplified for learning project)
type Message struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	SessionID  string         `gorm:"type:varchar(100);index" json:"session_id"`                                      // Session identifier for grouping conversations
	Role       string         `gorm:"type:varchar(20);not null;check:role IN ('user', 'assistant', 'system')" json:"role"`
	Content    string         `gorm:"type:text;not null" json:"content"`
	PersonaID  *int           `json:"persona_id,omitempty"`
	TokensUsed *int           `json:"tokens_used,omitempty"`
	CreatedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Metadata   datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"metadata,omitempty"`

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

// MessageRole constants
const (
	RoleUser      = "user"
	RoleAssistant = "assistant"
	RoleSystem    = "system"
)
