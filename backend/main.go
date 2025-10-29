package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"chatbot/config"
	"chatbot/models"
	"chatbot/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create required PostgreSQL extensions
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Fatal("Failed to create uuid-ossp extension:", err)
	}
	log.Println("✓ PostgreSQL extensions initialized")

	// Auto-migrate database models
	err = db.AutoMigrate(
		&models.Persona{},
		&models.Message{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✓ Database migration completed")

	// Seed personas if empty
	seedPersonas(db)

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: cfg.AppName,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path} ${latency}\n",
	}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSOrigin,
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Setup routes
	routes.SetupRoutes(app, db, cfg)

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Gracefully shutting down...")
		_ = app.Shutdown()
		config.CloseDatabase()
		os.Exit(0)
	}()

	// Start server
	log.Printf("🚀 Server starting on port %s (env: %s)", cfg.Port, cfg.AppEnv)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// seedPersonas inserts default personas if the table is empty
func seedPersonas(db *gorm.DB) {
	// Check if personas already exist
	var count int64
	db.Model(&models.Persona{}).Count(&count)
	if count > 0 {
		log.Println("✓ Personas already exist, skipping seed")
		return
	}

	personas := []models.Persona{
		{
			Name:         "General Assistant",
			SystemPrompt: "You are a helpful, friendly, and knowledgeable AI assistant. Provide clear, accurate, and concise responses to user questions.",
			Expertise:    "general",
			Description:  "A versatile AI assistant for general questions and conversations",
			Icon:         "🤖",
			IsActive:     true,
		},
		{
			Name:         "Technology Expert",
			SystemPrompt: "You are a technology expert with deep knowledge in software development, programming, cloud computing, and IT infrastructure. Provide technical and accurate responses.",
			Expertise:    "technology",
			Description:  "Specialized in programming, software architecture, and tech solutions",
			Icon:         "💻",
			IsActive:     true,
		},
		{
			Name:         "Business Advisor",
			SystemPrompt: "You are a professional business consultant. Provide strategic business advice, market analysis, and entrepreneurship guidance in a professional yet approachable manner.",
			Expertise:    "business",
			Description:  "Expert in business strategy, entrepreneurship, and market analysis",
			Icon:         "💼",
			IsActive:     true,
		},
	}

	for _, persona := range personas {
		if err := db.Create(&persona).Error; err != nil {
			log.Printf("Warning: Failed to seed persona %s: %v", persona.Name, err)
		}
	}

	log.Printf("✓ Seeded %d personas", len(personas))
}