package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase establishes connection to PostgreSQL database
func ConnectDatabase(config *Config) (*gorm.DB, error) {
	// Configure GORM logger
	var gormLogger logger.Interface
	if config.AppEnv == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Warn)
	}

	// Open database connection
	db, err := gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10) // กำหนดจำนวนการเชื่อมต่อ สูงสุด ที่สามารถอยู่ในสถานะ "ว่าง" (Idle) และพร้อมให้ใช้งานใน Pool
	sqlDB.SetMaxOpenConns(100) // กำหนดจำนวนการเชื่อมต่อ สูงสุดทั้งหมด ที่สามารถเปิดใช้งานได้พร้อมกัน (ทั้งที่ใช้งานอยู่และที่อยู่ในสถานะ Idle)
	sqlDB.SetConnMaxLifetime(time.Hour) // กำหนดระยะเวลา สูงสุด ที่การเชื่อมต่อฐานข้อมูลแต่ละครั้งจะสามารถเปิดใช้งานอยู่ได้ ค่าที่ตั้ง: time.Hour (หนึ่งชั่วโมง)

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db
	log.Println("✓ Database connected successfully")

	return db, nil
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
