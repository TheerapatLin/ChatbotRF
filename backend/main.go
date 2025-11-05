package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"chatbot/config"
	"chatbot/database"
	"chatbot/models"
	"chatbot/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
		&models.FileAnalysis{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("✓ Database migration completed")

	// Seed personas if empty
	database.SeedPersonas(db)

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
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS",
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