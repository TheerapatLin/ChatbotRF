package routes

import (
	"log"

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
	elevenLabsService := services.NewElevenLabsService(cfg)
	contextService := services.NewContextService(messageRepo, fileAnalysisRepo)
	fileService := services.NewFileService(openaiService.GetClient(), contextService)

	// Initialize Bedrock service
	bedrockService, err := services.NewBedrockService(cfg)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to initialize Bedrock service: %v", err)
		log.Printf("   Bedrock endpoints will not be available")
	}

	// Initialize Whisper.cpp service
	whisperService, err := services.NewWhisperCppService(cfg)
	if err != nil {
		log.Printf("⚠️ Warning: Failed to initialize Whisper.cpp service: %v", err)
		log.Printf("   Whisper.cpp endpoints will not be available")
	}

	// Initialize controllers
	chatCtrl := controllers.NewChatController(messageRepo, personaRepo, fileAnalysisRepo, openaiService)
	personaCtrl := controllers.NewPersonaController(personaRepo, messageRepo)
	audioCtrl := controllers.NewAudioController(openaiService, ttsService)
	elevenLabsCtrl := controllers.NewElevenLabsController(elevenLabsService)
	wsCtrl := controllers.NewWebSocketController(messageRepo, personaRepo, fileAnalysisRepo, openaiService, bedrockService)
	ttsWSCtrl := controllers.NewTTSWebSocketController(ttsService, personaRepo)
	elevenLabsWSCtrl := controllers.NewElevenLabsWSController(elevenLabsService)
	fileCtrl := controllers.NewFileController(fileService, fileAnalysisRepo, messageRepo)

	// Initialize Bedrock controller
	var bedrockCtrl *controllers.BedrockController
	if bedrockService != nil {
		bedrockCtrl = controllers.NewBedrockController(bedrockService, personaRepo, messageRepo, contextService, fileAnalysisRepo)
	}

	// Initialize Whisper.cpp controller
	var whisperCtrl *controllers.WhisperCppController
	if whisperService != nil {
		whisperCtrl = controllers.NewWhisperCppController(whisperService)
	}

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
	api.Post("/persona", personaCtrl.CreatePersona)
	api.Patch("/persona/:id", personaCtrl.UpdatePersona)
	api.Delete("/persona/:id", personaCtrl.DeletePersona)

	// Chat endpoints (OpenAI)
	api.Post("/chat", chatCtrl.HandleChat)
	api.Get("/chats", chatCtrl.GetChatHistory)
	api.Get("/chats/session/:sessionId", chatCtrl.GetChatHistoryBySession)
	api.Delete("/chats", chatCtrl.DeleteAllMessages)
	api.Delete("/chats/session/:sessionId", chatCtrl.DeleteMessagesBySession)

	// Bedrock endpoints (AWS Bedrock)
	if bedrockCtrl != nil {
		api.Post("/chat/bedrock", bedrockCtrl.SendBedrockMessage)
	}

	// Audio endpoints
	api.Post("/audio/stt/openai", audioCtrl.TranscribeAudio)
	api.Post("/audio/tts/openai", audioCtrl.TextToSpeech)

	// ElevenLabs endpoints
	api.Post("/audio/elevenlabs/tts", elevenLabsCtrl.TextToSpeech)
	api.Get("/audio/elevenlabs/voices", elevenLabsCtrl.GetVoices)

	// ========================================
	// Speech-to-Text (STT) Routes
	// ========================================
	if whisperCtrl != nil {
		sttGroup := api.Group("/stt")

		// Whisper.cpp endpoints
		sttGroup.Get("/whispercpp/status", whisperCtrl.GetStatus)
		sttGroup.Post("/whispercpp", whisperCtrl.TranscribeAudio)

		log.Println("✅ Whisper.cpp endpoints registered:")
		log.Println("   GET  /api/stt/whispercpp/status")
		log.Println("   POST /api/stt/whispercpp")
	}

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

	// WebSocket upgrade middleware for TTS
	app.Use("/api/ws/tts", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint for TTS streaming
	app.Get("/api/ws/tts", websocket.New(ttsWSCtrl.HandleTTSWebSocket))
	log.Println("✅ TTS WebSocket endpoint registered at: ws://localhost:3001/api/ws/tts")

	// WebSocket upgrade middleware for ElevenLabs
	app.Use("/api/ws/elevenlabs", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	// WebSocket endpoint for ElevenLabs TTS streaming
	app.Get("/api/ws/elevenlabs", websocket.New(elevenLabsWSCtrl.HandleElevenLabsWebSocket))
	log.Println("✅ ElevenLabs WebSocket endpoint registered at: ws://localhost:3001/api/ws/elevenlabs")
}
