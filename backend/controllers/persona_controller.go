package controllers

import (
	"encoding/json"
	"log"
	"time"

	"chatbot/models"
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
	ID              int     `json:"id"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	SystemPrompt    string  `json:"system_prompt"`
	Tone            string  `json:"tone"`
	Style           string  `json:"style"`
	Expertise       string  `json:"expertise"`
	Temperature     float32 `json:"temperature"`
	MaxTokens       int     `json:"max_tokens"`
	Model           string  `json:"model"`
	LanguageSetting string  `json:"language_setting"`
	Guardrails      string  `json:"guardrails"`
	Icon            string  `json:"icon"`
	IsActive        bool    `json:"is_active"`
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
			ID:              persona.ID,
			Name:            persona.Name,
			Description:     persona.Description,
			SystemPrompt:    persona.SystemPrompt,
			Tone:            persona.Tone,
			Style:           persona.Style,
			Expertise:       persona.Expertise,
			Temperature:     persona.Temperature,
			MaxTokens:       persona.MaxTokens,
			Model:           persona.Model,
			LanguageSetting: persona.LanguageSetting,
			Guardrails:      persona.Guardrails,
			Icon:            persona.Icon,
			IsActive:        persona.IsActive,
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

// CreatePersonaRequest represents request body for creating a persona
type CreatePersonaRequest struct {
	Name            string                 `json:"name" validate:"required"`
	Description     string                 `json:"description" validate:"required"`
	SystemPrompt    string                 `json:"system_prompt" validate:"required"`
	Tone            string                 `json:"tone"`
	Style           string                 `json:"style"`
	Expertise       string                 `json:"expertise"`
	Temperature     float32                `json:"temperature"`
	MaxTokens       int                    `json:"max_tokens"`
	Model           string                 `json:"model"`
	LanguageSetting LanguageSettingRequest `json:"language_setting"`
	Guardrails      GuardrailsRequest      `json:"guardrails"`
	Icon            string                 `json:"icon"`
}

type LanguageSettingRequest struct {
	DefaultLanguage string `json:"default_language"`
	ResponseStyle   string `json:"response_style"`
	LanguageCode    string `json:"language_code"`
}

type GuardrailsRequest struct {
	BlockProfanity    bool     `json:"block_profanity"`
	BlockSensitive    bool     `json:"block_sensitive"`
	AllowedTopics     []string `json:"allowed_topics"`
	BlockedTopics     []string `json:"blocked_topics"`
	MaxResponseLength int      `json:"max_response_length"`
	RequireModeration bool     `json:"require_moderation"`
}

// CreatePersona handles POST /api/personas endpoint
func (ctrl *PersonaController) CreatePersona(c *fiber.Ctx) error {
	var req CreatePersonaRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		log.Printf("❌ Failed to parse request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	// Validate required fields
	if req.Name == "" || req.Description == "" || req.SystemPrompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, description, and system_prompt are required",
		})
	}

	// Validate field lengths
	if len(req.Name) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Name must be less than 100 characters",
		})
	}
	if len(req.Tone) > 200 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tone must be less than 200 characters",
		})
	}
	if len(req.Style) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Style must be less than 500 characters",
		})
	}
	if len(req.Expertise) > 500 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Expertise must be less than 500 characters",
		})
	}
	if len(req.Icon) > 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Icon must be less than 10 characters",
		})
	}

	// Set defaults if not provided
	if req.Temperature == 0 {
		req.Temperature = 0.7
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = 2000
	}
	if req.Model == "" {
		req.Model = "gpt-4o-mini"
	}
	if req.Icon == "" {
		req.Icon = "🤖"
	}
	if req.Tone == "" {
		req.Tone = "friendly"
	}
	if req.Style == "" {
		req.Style = "conversational"
	}
	if req.Expertise == "" {
		req.Expertise = "general"
	}

	// Validate model name
	validModels := map[string]bool{
		"gpt-4o-mini":   true,
		"gpt-4o":        true,
		"gpt-4":         true,
		"gpt-3.5-turbo": true,
		"apac.anthropic.claude-sonnet-4-20250514-v1:0": true,
		"claude-sonnet-4": true,
		"claude-3-opus":   true,
		"claude-3-sonnet": true,
	}

	if !validModels[req.Model] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":        "Invalid model name",
			"valid_models": []string{"gpt-4o-mini", "gpt-4o", "gpt-4", "gpt-3.5-turbo", "claude-sonnet-4", "claude-3-opus", "claude-3-sonnet"},
			"received":     req.Model,
		})
	}

	// Validate temperature range
	if req.Temperature < 0 || req.Temperature > 2.0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Temperature must be between 0.0 and 2.0",
		})
	}

	// Marshal language setting to JSON
	languageSettingJSON, err := json.Marshal(req.LanguageSetting)
	if err != nil {
		log.Printf("❌ Failed to marshal language_setting: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process language settings",
		})
	}

	// Marshal guardrails to JSON
	guardrailsJSON, err := json.Marshal(req.Guardrails)
	if err != nil {
		log.Printf("❌ Failed to marshal guardrails: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to process guardrails",
		})
	}

	// Create persona model
	persona := &models.Persona{
		Name:            req.Name,
		Description:     req.Description,
		SystemPrompt:    req.SystemPrompt,
		Tone:            req.Tone,
		Style:           req.Style,
		Expertise:       req.Expertise,
		Temperature:     req.Temperature,
		MaxTokens:       req.MaxTokens,
		Model:           req.Model,
		LanguageSetting: string(languageSettingJSON),
		Guardrails:      string(guardrailsJSON),
		Icon:            req.Icon,
		IsActive:        true, // New personas are active by default
	}

	// Save to database
	if err := ctrl.personaRepo.Create(persona); err != nil {
		log.Printf("❌ Failed to create persona: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create persona",
		})
	}

	log.Printf("✅ Persona created successfully: id=%d, name=%s", persona.ID, persona.Name)

	// Return created persona
	response := PersonaResponse{
		ID:              persona.ID,
		Name:            persona.Name,
		Description:     persona.Description,
		Tone:            persona.Tone,
		Style:           persona.Style,
		Expertise:       persona.Expertise,
		Temperature:     persona.Temperature,
		MaxTokens:       persona.MaxTokens,
		Model:           persona.Model,
		LanguageSetting: persona.LanguageSetting,
		Guardrails:      persona.Guardrails,
		Icon:            persona.Icon,
		IsActive:        persona.IsActive,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// DeletePersona handles DELETE /api/persona/:id endpoint
func (ctrl *PersonaController) DeletePersona(c *fiber.Ctx) error {
	// Parse ID from URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid persona ID",
		})
	}

	// Check if persona exists
	persona, err := ctrl.personaRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Persona not found",
		})
	}

	// Count messages associated with this persona
	messageCount, err := ctrl.messageRepo.CountByPersonaID(id)
	if err != nil {
		log.Printf("⚠️ Failed to count messages for persona %d: %v", id, err)
		messageCount = 0
	}

	// Delete all messages associated with this persona first
	if messageCount > 0 {
		if err := ctrl.messageRepo.DeleteByPersonaID(id); err != nil {
			log.Printf("❌ Failed to delete messages for persona %d: %v", id, err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to delete persona messages",
			})
		}
		log.Printf("🗑️ Deleted %d messages for persona %d", messageCount, id)
	}

	// Delete persona from database
	if err := ctrl.personaRepo.Delete(id); err != nil {
		log.Printf("❌ Failed to delete persona: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete persona",
		})
	}

	log.Printf("✅ Persona deleted successfully: id=%d, name=%s (deleted %d messages)", id, persona.Name, messageCount)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":          "Persona deleted successfully",
		"id":               id,
		"messages_deleted": messageCount,
	})
}

// UpdatePersona handles PATCH /api/persona/:id endpoint
func (ctrl *PersonaController) UpdatePersona(c *fiber.Ctx) error {
	// Parse ID from URL parameter
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid persona ID",
		})
	}

	// Check if persona exists
	persona, err := ctrl.personaRepo.FindByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Persona not found",
		})
	}

	// Parse request body (partial update)
	type UpdatePersonaRequest struct {
		Name            *string                 `json:"name"`
		Description     *string                 `json:"description"`
		SystemPrompt    *string                 `json:"system_prompt"`
		Tone            *string                 `json:"tone"`
		Style           *string                 `json:"style"`
		Expertise       *string                 `json:"expertise"`
		Temperature     *float32                `json:"temperature"`
		MaxTokens       *int                    `json:"max_tokens"`
		Model           *string                 `json:"model"`
		LanguageSetting *LanguageSettingRequest `json:"language_setting"`
		Guardrails      *GuardrailsRequest      `json:"guardrails"`
		Icon            *string                 `json:"icon"`
		IsActive        *bool                   `json:"is_active"`
	}

	var req UpdatePersonaRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update only provided fields
	if req.Name != nil {
		if len(*req.Name) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name cannot be empty",
			})
		}
		if len(*req.Name) > 100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name must be less than 100 characters",
			})
		}
		persona.Name = *req.Name
	}

	if req.Description != nil {
		persona.Description = *req.Description
	}

	if req.SystemPrompt != nil {
		if len(*req.SystemPrompt) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "System prompt cannot be empty",
			})
		}
		persona.SystemPrompt = *req.SystemPrompt
	}

	if req.Tone != nil {
		if len(*req.Tone) > 200 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Tone must be less than 200 characters",
			})
		}
		persona.Tone = *req.Tone
	}

	if req.Style != nil {
		if len(*req.Style) > 500 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Style must be less than 500 characters",
			})
		}
		persona.Style = *req.Style
	}

	if req.Expertise != nil {
		if len(*req.Expertise) > 500 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Expertise must be less than 500 characters",
			})
		}
		persona.Expertise = *req.Expertise
	}

	if req.Temperature != nil {
		if *req.Temperature < 0 || *req.Temperature > 2.0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Temperature must be between 0.0 and 2.0",
			})
		}
		persona.Temperature = *req.Temperature
	}

	if req.MaxTokens != nil {
		persona.MaxTokens = *req.MaxTokens
	}

	if req.Model != nil {
		// Validate model name
		validModels := map[string]bool{
			"gpt-4o-mini":   true,
			"gpt-4o":        true,
			"gpt-4":         true,
			"gpt-3.5-turbo": true,
			"apac.anthropic.claude-sonnet-4-20250514-v1:0": true,
			"claude-sonnet-4": true,
			"claude-3-opus":   true,
			"claude-3-sonnet": true,
		}

		if !validModels[*req.Model] {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":        "Invalid model name",
				"valid_models": []string{"gpt-4o-mini", "gpt-4o", "gpt-4", "gpt-3.5-turbo", "claude-sonnet-4", "claude-3-opus", "claude-3-sonnet"},
				"received":     *req.Model,
			})
		}
		persona.Model = *req.Model
	}

	if req.Icon != nil {
		if len(*req.Icon) > 10 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Icon must be less than 10 characters",
			})
		}
		persona.Icon = *req.Icon
	}

	if req.IsActive != nil {
		persona.IsActive = *req.IsActive
	}

	// Update language setting if provided
	if req.LanguageSetting != nil {
		languageSettingJSON, err := json.Marshal(req.LanguageSetting)
		if err != nil {
			log.Printf("❌ Failed to marshal language_setting: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process language settings",
			})
		}
		persona.LanguageSetting = string(languageSettingJSON)
	}

	// Update guardrails if provided
	if req.Guardrails != nil {
		guardrailsJSON, err := json.Marshal(req.Guardrails)
		if err != nil {
			log.Printf("❌ Failed to marshal guardrails: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process guardrails",
			})
		}
		persona.Guardrails = string(guardrailsJSON)
	}

	// Save to database
	if err := ctrl.personaRepo.Update(persona); err != nil {
		log.Printf("❌ Failed to update persona: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update persona",
		})
	}

	log.Printf("✅ Persona updated successfully: id=%d, name=%s", persona.ID, persona.Name)

	// Return updated persona
	response := PersonaResponse{
		ID:              persona.ID,
		Name:            persona.Name,
		Description:     persona.Description,
		SystemPrompt:    persona.SystemPrompt,
		Tone:            persona.Tone,
		Style:           persona.Style,
		Expertise:       persona.Expertise,
		Temperature:     persona.Temperature,
		MaxTokens:       persona.MaxTokens,
		Model:           persona.Model,
		LanguageSetting: persona.LanguageSetting,
		Guardrails:      persona.Guardrails,
		Icon:            persona.Icon,
		IsActive:        persona.IsActive,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
