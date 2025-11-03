package models

import (
	"time"
)

// LanguageSetting represents language configuration for a persona
type LanguageSetting struct {
	DefaultLanguage string `json:"default_language"` // e.g., "th", "en"
	ResponseStyle   string `json:"response_style"`   // e.g., "formal", "casual", "professional"
	LanguageCode    string `json:"language_code"`    // ISO 639-1 code
}

// Guardrails represents content filtering rules for a persona
type Guardrails struct {
	BlockProfanity     bool     `json:"block_profanity"`
	BlockSensitive     bool     `json:"block_sensitive"`
	AllowedTopics      []string `json:"allowed_topics"`
	BlockedTopics      []string `json:"blocked_topics"`
	MaxResponseLength  int      `json:"max_response_length"`
	RequireModeration  bool     `json:"require_moderation"`
}

// Persona represents an AI personality/character with comprehensive configuration
type Persona struct {
	ID             int        `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string     `gorm:"type:varchar(100);not null;unique" json:"name"`
	Description    string     `gorm:"type:text" json:"description"`
	SystemPrompt   string     `gorm:"type:text;not null" json:"system_prompt"`
	Tone           string     `gorm:"type:varchar(50)" json:"tone"`              // e.g., "friendly", "professional", "empathetic"
	Style          string     `gorm:"type:varchar(50)" json:"style"`             // e.g., "concise", "detailed", "conversational"
	Expertise      string     `gorm:"type:varchar(100)" json:"expertise"`        // e.g., "technology", "healthcare", "education"
	Temperature    float32    `gorm:"type:decimal(3,2);default:0.7" json:"temperature"` // 0.0 - 2.0
	MaxTokens      int        `gorm:"default:2000" json:"max_tokens"`
	Model          string     `gorm:"type:varchar(50);default:'gpt-4o-mini'" json:"model"` // e.g., "gpt-4o-mini", "gpt-4"
	LanguageSetting string    `gorm:"type:jsonb" json:"language_setting"` // JSON field for language settings
	Guardrails     string     `gorm:"type:jsonb" json:"guardrails"`       // JSON field for guardrails
	Icon           string     `gorm:"type:varchar(50)" json:"icon"`
	IsActive       bool       `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// TableName specifies the table name for Persona model
func (Persona) TableName() string {
	return "personas"
}