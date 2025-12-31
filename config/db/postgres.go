package db

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresConfig() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		return nil, errors.New("Error Loading .env File")
	}

	DATABASE_URL := os.Getenv("DATABASE_URL")
	dsn := DATABASE_URL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	log.Println("DATABASE CONNECTED")
	AutoMigrate(db)
	return db, nil
}
