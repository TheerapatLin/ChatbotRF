package repositories

import (
	"chatbot/models"
	"gorm.io/gorm"
)

// PersonaRepository handles database operations for personas
type PersonaRepository struct {
	db *gorm.DB
}

// NewPersonaRepository creates a new persona repository
func NewPersonaRepository(db *gorm.DB) *PersonaRepository {
	return &PersonaRepository{db: db}
}

// FindAll retrieves all personas
func (r *PersonaRepository) FindAll() ([]models.Persona, error) {
	var personas []models.Persona
	err := r.db.Find(&personas).Error
	return personas, err
}

// FindActive retrieves all active personas
func (r *PersonaRepository) FindActive() ([]models.Persona, error) {
	var personas []models.Persona
	err := r.db.Where("is_active = ?", true).Find(&personas).Error
	return personas, err
}

// FindByID retrieves a persona by its ID
func (r *PersonaRepository) FindByID(id int) (*models.Persona, error) {
	var persona models.Persona
	err := r.db.Where("id = ?", id).First(&persona).Error
	if err != nil {
		return nil, err
	}
	return &persona, nil
}