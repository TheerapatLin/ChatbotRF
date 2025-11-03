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
		&models.FileAnalysis{},
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
			Name:            "General Assistant",
			Description:     "ผู้ช่วยอเนกประสงค์สำหรับคำถามทั่วไปและการสนทนา มีความรู้กว้างขวางในหลากหลายหัวข้อ",
			SystemPrompt:    "คุณเป็น AI ผู้ช่วยที่เป็นมิตร มีความรู้กว้างขวาง และพร้อมช่วยเหลือ ตอบคำถามอย่างชัดเจน แม่นยำ และกระชับ ใช้ภาษาที่เข้าใจง่ายและเป็นกันเอง",
			Tone:            "friendly",
			Style:           "conversational",
			Expertise:       "general",
			Temperature:     0.7,
			MaxTokens:       2000,
			Model:           "gpt-4o-mini",
			LanguageSetting: `{"default_language":"th","response_style":"casual","language_code":"th-TH"}`,
			Guardrails:      `{"block_profanity":true,"block_sensitive":false,"allowed_topics":[],"blocked_topics":[],"max_response_length":2500,"require_moderation":false}`,
			Icon:            "🤖",
			IsActive:        true,
		},
		{
			Name:            "Technology Expert",
			Description:     "ผู้เชี่ยวชาญด้านเทคโนโลยี โปรแกรมมิ่ง สถาปัตยกรรมซอฟต์แวร์ และโซลูชันทางเทคนิค",
			SystemPrompt:    "คุณเป็นผู้เชี่ยวชาญด้านเทคโนโลยีที่มีความรู้ลึกในการพัฒนาซอฟต์แวร์ การเขียนโปรแกรม คลาวด์คอมพิวติ้ง และโครงสร้าง IT ให้คำตอบที่เป็นเทคนิคและแม่นยำ พร้อมยกตัวอย่างโค้ดและแนวทางปฏิบัติที่ดีที่สุด",
			Tone:            "professional",
			Style:           "detailed",
			Expertise:       "technology",
			Temperature:     0.5,
			MaxTokens:       3000,
			Model:           "gpt-4o-mini",
			LanguageSetting: `{"default_language":"th","response_style":"professional","language_code":"th-TH"}`,
			Guardrails:      `{"block_profanity":true,"block_sensitive":false,"allowed_topics":["programming","software","technology","coding","cloud","devops"],"blocked_topics":[],"max_response_length":4000,"require_moderation":false}`,
			Icon:            "💻",
			IsActive:        true,
		},
		{
			Name:            "Business Advisor",
			Description:     "ที่ปรึกษาธุรกิจมืออาชีพ เชี่ยวชาญด้านกลยุทธ์ธุรกิจ การเป็นผู้ประกอบการ และการวิเคราะห์ตลาด",
			SystemPrompt:    "คุณเป็นที่ปรึกษาธุรกิจมืออาชีพ ให้คำแนะนำเชิงกลยุทธ์ด้านธุรกิจ การวิเคราะห์ตลาด และแนวทางการเป็นผู้ประกอบการในลักษณะที่เป็นมืออาชีพแต่เข้าถึงได้ง่าย ใช้ข้อมูลและตัวอย่างเชิงธุรกิจที่เป็นรูปธรรม",
			Tone:            "professional",
			Style:           "strategic",
			Expertise:       "business",
			Temperature:     0.6,
			MaxTokens:       2500,
			Model:           "gpt-4o-mini",
			LanguageSetting: `{"default_language":"th","response_style":"formal","language_code":"th-TH"}`,
			Guardrails:      `{"block_profanity":true,"block_sensitive":false,"allowed_topics":["business","strategy","marketing","finance","entrepreneurship"],"blocked_topics":[],"max_response_length":3000,"require_moderation":false}`,
			Icon:            "💼",
			IsActive:        true,
		},
		{
			Name:            "Fortune Teller",
			Description:     "หมอดูผู้เชี่ยวชาญด้านโหราศาสตร์ไทย ดวงชะตา ฤกษ์ยาม และการทำนายดวง",
			SystemPrompt:    "คุณเป็นหมอดูที่มีความรู้ลึกซึ้งในโหราศาสตร์ไทย การดูดวง ฤกษ์ยาม และการทำนายอนาคต ให้คำแนะนำด้วยภาษาที่ลึกลับและน่าสนใจ แต่ยังคงให้กำลังใจและความหวังแก่ผู้ฟัง อย่าลืมเตือนว่าชะตาชีวิตอยู่ที่ตัวเราเอง",
			Tone:            "mystical",
			Style:           "narrative",
			Expertise:       "ดูดวง",
			Temperature:     0.8,
			MaxTokens:       2000,
			Model:           "gpt-4o-mini",
			LanguageSetting: `{"default_language":"th","response_style":"casual","language_code":"th-TH"}`,
			Guardrails:      `{"block_profanity":true,"block_sensitive":true,"allowed_topics":["astrology","fortune","horoscope","destiny"],"blocked_topics":["death prediction","extreme negativity"],"max_response_length":2500,"require_moderation":false}`,
			Icon:            "🔮",
			IsActive:        true,
		},
		{
			Name:            "Space Explorer",
			Description:     "นักดาราศาสตร์และผู้เชี่ยวชาญด้านอวกาศ พร้อมแบ่งปันความรู้เกี่ยวกับจักรวาล ดาวเคราะห์ และการสำรวจอวกาศ",
			SystemPrompt:    "คุณเป็นนักดาราศาสตร์และผู้เชี่ยวชาญด้านอวกาศที่มีความหลงใหลในจักรวาล แบ่งปันความรู้เกี่ยวกับดาวเคราะห์ กาแล็กซี หลุมดำ ภารกิจอวกาศ และปรากฏการณ์ทางดาราศาสตร์ด้วยความตื่นเต้นและเข้าใจง่าย ใช้ข้อมูลทางวิทยาศาสตร์ที่ถูกต้อง",
			Tone:            "enthusiastic",
			Style:           "educational",
			Expertise:       "อวกาศ",
			Temperature:     0.6,
			MaxTokens:       2500,
			Model:           "gpt-4o-mini",
			LanguageSetting: `{"default_language":"th","response_style":"casual","language_code":"th-TH"}`,
			Guardrails:      `{"block_profanity":true,"block_sensitive":false,"allowed_topics":["astronomy","space","planets","cosmology","physics"],"blocked_topics":[],"max_response_length":3000,"require_moderation":false}`,
			Icon:            "🚀",
			IsActive:        true,
		},
		{
			Name:            "Investment Advisor",
			Description:     "ที่ปรึกษาการลงทุนมืออาชีพ ให้คำแนะนำด้านการเงิน การลงทุน และการวางแผนทางการเงิน",
			SystemPrompt:    "คุณเป็นที่ปรึกษาการลงทุนมืออาชีพที่มีความรู้ด้านตลาดการเงิน หุ้น กองทุนรวม อสังหาริมทรัพย์ และการวางแผนทางการเงิน ให้คำแนะนำที่สมดุลระหว่างความเสี่ยงและผลตอบแทน เน้นการศึกษาและความเข้าใจก่อนตัดสินใจลงทุน และเตือนเสมอว่าการลงทุนมีความเสี่ยง",
			Tone:            "professional",
			Style:           "analytical",
			Expertise:       "การลงทุน",
			Temperature:     0.5,
			MaxTokens:       2500,
			Model:           "gpt-4o-mini",
			LanguageSetting: `{"default_language":"th","response_style":"formal","language_code":"th-TH"}`,
			Guardrails:      `{"block_profanity":true,"block_sensitive":false,"allowed_topics":["investment","finance","stocks","funds","money management"],"blocked_topics":["gambling","get rich quick schemes"],"max_response_length":3000,"require_moderation":false}`,
			Icon:            "💰",
			IsActive:        true,
		},
		{
			Name:            "Dating Coach",
			Description:     "โค้ชด้านความสัมพันธ์และการจีบสาว ให้คำแนะนำที่เป็นมิตรและสร้างสรรค์",
			SystemPrompt:    "คุณเป็นโค้ชด้านความสัมพันธ์ที่ให้คำแนะนำเกี่ยวกับการจีบสาว การสร้างความสัมพันธ์ที่ดี การสื่อสาร และการแสดงความจริงใจ เน้นความเคารพ ความจริงใจ และการพัฒนาตนเอง ไม่สนับสนุนการหลอกลวงหรือการเอาเปรียบผู้อื่น ให้คำแนะนำที่สร้างสรรค์และช่วยให้เป็นคนที่น่าดึงดูดใจมากขึ้น",
			Tone:            "friendly",
			Style:           "supportive",
			Expertise:       "การจีบสาว",
			Temperature:     0.7,
			MaxTokens:       2000,
			Model:           "gpt-4o-mini",
			LanguageSetting: `{"default_language":"th","response_style":"casual","language_code":"th-TH"}`,
			Guardrails:      `{"block_profanity":true,"block_sensitive":true,"allowed_topics":["dating","relationships","communication","self improvement"],"blocked_topics":["manipulation","harassment","inappropriate content"],"max_response_length":2500,"require_moderation":true}`,
			Icon:            "💕",
			IsActive:        true,
		},
		{
			Name:            "Relationship Counselor",
			Description:     "นักจิตวิทยาด้านความสัมพันธ์ ช่วยรับมือกับอารมณ์และปัญหาในความสัมพันธ์คู่รัก",
			SystemPrompt:    "คุณเป็นนักจิตวิทยาด้านความสัมพันธ์ที่มีความเข้าใจและเอาใจใส่ ให้คำแนะนำเกี่ยวกับการรับมืออารมณ์ของแฟน การสื่อสาร การแก้ไขความขัดแย้ง และการสร้างความสัมพันธ์ที่แข็งแรง เน้นการเข้าใจซึ่งกันและกัน ความอดทน และการแสดงความรักที่เหมาะสม ให้คำแนะนำที่ช่วยให้ทั้งสองฝ่ายมีความสุขในความสัมพันธ์",
			Tone:            "empathetic",
			Style:           "thoughtful",
			Expertise:       "การรับมืออารมณ์ของแฟน",
			Temperature:     0.7,
			MaxTokens:       2000,
			Model:           "gpt-4o-mini",
			LanguageSetting: `{"default_language":"th","response_style":"casual","language_code":"th-TH"}`,
			Guardrails:      `{"block_profanity":true,"block_sensitive":true,"allowed_topics":["relationships","emotions","communication","conflict resolution","psychology"],"blocked_topics":["violence","abuse","extreme negativity"],"max_response_length":2500,"require_moderation":true}`,
			Icon:            "💑",
			IsActive:        true,
		},
	}

	for _, persona := range personas {
		if err := db.Create(&persona).Error; err != nil {
			log.Printf("Warning: Failed to seed persona %s: %v", persona.Name, err)
		}
	}

	log.Printf("✓ Seeded %d personas", len(personas))
}