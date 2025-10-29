package controllers

import (
	"time"

	"chatbot/repositories"

	"github.com/gofiber/fiber/v2"
)

// PersonaController handles persona-related HTTP requests
type PersonaController struct {
	personaRepo *repositories.PersonaRepository
	messageRepo *repositories.MessageRepository
}

// NewPersonaController creates a new persona controller
func NewPersonaController(personaRepo *repositories.PersonaRepository, messageRepo *repositories.MessageRepository) *PersonaController {
	return &PersonaController{
		personaRepo: personaRepo,
		messageRepo: messageRepo,
	}
}

// PersonaResponse represents a persona in API response
type PersonaResponse struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	SystemPrompt string    `json:"system_prompt"`
	Expertise    string    `json:"expertise"`
	Description  string    `json:"description"`
	Icon         string    `json:"icon"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

// PersonasListResponse represents the list of personas response
type PersonasListResponse struct {
	Personas []PersonaResponse `json:"personas"`
}

// PersonaStats represents statistics for a persona
type PersonaStats struct {
	TotalMessages   int64  `json:"total_messages"`
	AvgResponseTime string `json:"avg_response_time"`
}

// PersonaDetailResponse represents a single persona with stats
type PersonaDetailResponse struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	SystemPrompt string       `json:"system_prompt"`
	Expertise    string       `json:"expertise"`
	Description  string       `json:"description"`
	Icon         string       `json:"icon"`
	IsActive     bool         `json:"is_active"`
	CreatedAt    time.Time    `json:"created_at"`
	Stats        PersonaStats `json:"stats"`
}

// GetAllPersonas handles GET /api/personas endpoint
func (ctrl *PersonaController) GetAllPersonas(c *fiber.Ctx) error {
	// Get all active personas from repository
	personas, err := ctrl.personaRepo.FindActive()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve personas",
		})
	}

	// Convert to response format
	personaResponses := make([]PersonaResponse, len(personas))
	for i, persona := range personas {
		personaResponses[i] = PersonaResponse{
			ID:           persona.ID,
			Name:         persona.Name,
			SystemPrompt: persona.SystemPrompt,
			Expertise:    persona.Expertise,
			Description:  persona.Description,
			Icon:         persona.Icon,
			IsActive:     persona.IsActive,
			CreatedAt:    persona.CreatedAt,
		}
	}

	// Return response
	response := PersonasListResponse{
		Personas: personaResponses,
	}

	return c.Status(200).JSON(response)
}

// GetPersonaByID handles GET /api/personas/:id endpoint
func (ctrl *PersonaController) GetPersonaByID(c *fiber.Ctx) error {
	// Parse ID from URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid persona ID",
		})
	}

	// Get persona from repository
	persona, err := ctrl.personaRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "persona with id " + c.Params("id") + " not found",
		})
	}

	// Get message count for this persona
	messageCount, err := ctrl.messageRepo.CountByPersonaID(id)
	if err != nil {
		messageCount = 0 // Default to 0 if error
	}

	// Create response with stats
	response := PersonaDetailResponse{
		ID:           persona.ID,
		Name:         persona.Name,
		SystemPrompt: persona.SystemPrompt,
		Expertise:    persona.Expertise,
		Description:  persona.Description,
		Icon:         persona.Icon,
		IsActive:     persona.IsActive,
		CreatedAt:    persona.CreatedAt,
		Stats: PersonaStats{
			TotalMessages:   messageCount,
			AvgResponseTime: "2.3s", // Mock value for now
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
