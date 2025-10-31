package repositories

import (
	"chatbot/models"
	"gorm.io/gorm"
)

// MessageRepository handles database operations for messages
type MessageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new message repository
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create saves a new message to database
func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

// FindByID retrieves a message by its ID
func (r *MessageRepository) FindByID(id string) (*models.Message, error) {
	var message models.Message
	err := r.db.Where("id = ?", id).First(&message).Error
	if err != nil {
		return nil, err
	}
	return &message, nil
}

// FindRecent retrieves the most recent N messages
func (r *MessageRepository) FindRecent(limit int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Order("created_at DESC").Limit(limit).Find(&messages).Error
	return messages, err
}

// Delete removes a message by ID
func (r *MessageRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Message{}).Error
}

// Count returns total number of messages
func (r *MessageRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).Count(&count).Error
	return count, err
}

// FindWithPagination retrieves messages with pagination support
func (r *MessageRepository) FindWithPagination(limit, offset int) ([]models.Message, int64, error) {
	var messages []models.Message
	var total int64

	// Get total count
	if err := r.db.Model(&models.Message{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated messages ordered by creation time (newest first)
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&messages).Error
	return messages, total, err
}

// CountByPersonaID counts total messages for a specific persona
func (r *MessageRepository) CountByPersonaID(personaID int) (int64, error) {
	var count int64
	err := r.db.Model(&models.Message{}).Where("persona_id = ?", personaID).Count(&count).Error
	return count, err
}

// GetRecentBySession retrieves recent N messages from a specific session
func (r *MessageRepository) GetRecentBySession(sessionID string, limit int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("session_id = ?", sessionID).
		Order("created_at ASC").  // Oldest first for proper conversation flow
		Limit(limit).
		Find(&messages).Error

	if err != nil {
		return nil, err
	}

	// If we have more messages than the limit, return only the most recent ones
	if len(messages) > limit {
		messages = messages[len(messages)-limit:]
	}

	return messages, nil
}

// GetAllBySession retrieves all messages from a specific session
func (r *MessageRepository) GetAllBySession(sessionID string) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Find(&messages).Error
	return messages, err
}

// DeleteSession removes all messages from a specific session
func (r *MessageRepository) DeleteSession(sessionID string) error {
	return r.db.Where("session_id = ?", sessionID).Delete(&models.Message{}).Error
}