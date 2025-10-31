package utils

import (
	"chatbot/models"
	"time"

	"github.com/google/uuid"
)

// PersonaResponse represents persona information in API responses
type PersonaResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Expertise   string `json:"expertise"`
	Icon        string `json:"icon"`
	Description string `json:"description"`
}

// MessageHistoryItem represents a message in history
type MessageHistoryItem struct {
	ID        string    `json:"id"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	PersonaID *int      `json:"persona_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// ToPersonaResponse converts a Persona model to PersonaResponse
func ToPersonaResponse(p *models.Persona) PersonaResponse {
	return PersonaResponse{
		ID:          p.ID,
		Name:        p.Name,
		Expertise:   p.Expertise,
		Icon:        p.Icon,
		Description: p.Description,
	}
}

// ToPersonaResponses converts a slice of Persona models to PersonaResponses
func ToPersonaResponses(personas []models.Persona) []PersonaResponse {
	responses := make([]PersonaResponse, len(personas))
	for i, persona := range personas {
		responses[i] = ToPersonaResponse(&persona)
	}
	return responses
}

// ToMessageHistoryItem converts a Message model to MessageHistoryItem
func ToMessageHistoryItem(m *models.Message) MessageHistoryItem {
	return MessageHistoryItem{
		ID:        m.ID.String(),
		Role:      string(m.Role),
		Content:   m.Content,
		PersonaID: m.PersonaID,
		CreatedAt: m.CreatedAt,
	}
}

// ToMessageHistoryItems converts a slice of Message models to MessageHistoryItems
func ToMessageHistoryItems(messages []models.Message) []MessageHistoryItem {
	items := make([]MessageHistoryItem, len(messages))
	for i, msg := range messages {
		items[i] = ToMessageHistoryItem(&msg)
	}
	return items
}

// UUIDToString safely converts UUID to string
func UUIDToString(id uuid.UUID) string {
	return id.String()
}