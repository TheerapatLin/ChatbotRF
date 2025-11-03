package routes

import (
	"chatbot/config"
	"chatbot/controllers"
	"chatbot/repositories"
	"chatbot/services"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// SetupRoutes sets up all application routes
func SetupRoutes(app *fiber.App, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	messageRepo := repositories.NewMessageRepository(db)
	personaRepo := repositories.NewPersonaRepository(db)
	fileAnalysisRepo := repositories.NewFileAnalysisRepository(db)

	// Initialize services
	openaiService := services.NewOpenAIService(cfg)
	ttsService := services.NewTTSService(cfg)
	contextService := services.NewContextService(messageRepo, fileAnalysisRepo)
	fileService := services.NewFileService(openaiService.GetClient(), contextService)

	// Initialize controllers
	chatCtrl := controllers.NewChatController(messageRepo, personaRepo, fileAnalysisRepo, openaiService)
	personaCtrl := controllers.NewPersonaController(personaRepo, messageRepo)
	audioCtrl := controllers.NewAudioController(openaiService, ttsService)
	wsCtrl := controllers.NewWebSocketController(messageRepo, personaRepo, fileAnalysisRepo, openaiService)
	fileCtrl := controllers.NewFileController(fileService, fileAnalysisRepo, messageRepo)

	// API group
	api := app.Group("/api")

	// Health check endpoint
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "ChatBot API is running",
			"env":     cfg.AppEnv,
		})
	})

	// Personas endpoints
	api.Get("/personas", personaCtrl.GetAllPersonas)
	api.Get("/personas/:id", personaCtrl.GetPersonaByID)

	// Chat endpoints
	api.Post("/chat", chatCtrl.HandleChat)
	api.Get("/chats", chatCtrl.GetChatHistory)
	api.Delete("/chats", chatCtrl.DeleteAllMessages)

	// Audio endpoints
	api.Post("/audio/transcribe", audioCtrl.TranscribeAudio)
	api.Post("/audio/tts", audioCtrl.TextToSpeech)

	// File upload endpoints
	api.Post("/file/uploads", fileCtrl.UploadFiles)
	api.Get("/file/history", fileCtrl.GetFileHistory)
	api.Delete("/file/uploads", fileCtrl.DeleteAllFiles)

	// WebSocket upgrade middleware: ตรวจสอบ request จาก client
	app.Use("/api/chat/stream", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client requested upgrade to the WebSocket protocol
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true) // อนุญาติให้ใช้ WebSocket
			return c.Next()           // ส่งต่อไปยัง WebSocket Endpoint
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint for streaming chat: Upgrade HTTP เป็น WebSocket
	app.Get("/api/chat/stream", websocket.New(wsCtrl.HandleStreamingChat))
}
