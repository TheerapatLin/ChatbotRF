package models

import (
	"time"
)

// Persona represents an AI personality/character
type Persona struct {
	ID           int       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string    `gorm:"type:varchar(100);not null" json:"name"`
	SystemPrompt string    `gorm:"type:text;not null" json:"system_prompt"`
	Expertise    string    `gorm:"type:varchar(100)" json:"expertise"`
	Description  string    `gorm:"type:text" json:"description"`
	Icon         string    `gorm:"type:varchar(50)" json:"icon"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
}

// TableName specifies the table name for Persona model
func (Persona) TableName() string {
	return "personas"
}