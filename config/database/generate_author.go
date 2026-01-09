package database

import (
	"log"
	"os"

	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"github.com/clemilsonazevedo/blog/internal/domain/enums"
	"github.com/clemilsonazevedo/blog/internal/http/auth"
	"gorm.io/gorm"
)

func CreateAuthor(db *gorm.DB) error {
	authorName := os.Getenv("AUTHOR_NAME")
	authorEmail := os.Getenv("AUTHOR_EMAIL")
	authorPassword := os.Getenv("AUTHOR_PASSWORD")
	if authorEmail == "" || authorPassword == "" || authorName == "" {
		log.Println("UserName or Email or Password not configured!")
		return nil
	}

	var count int64
	db.Model(entities.User{}).Where("email = ?", authorEmail).Count(&count)

	if count > 0 {
		log.Println("Author already exists")
		return nil
	}

	hashpassword, err := auth.HashPassword(authorPassword)
	if err != nil {
		return err
	}

	author := entities.User{
		UserName: authorName,
		Email:    authorEmail,
		Password: hashpassword,
		Role:     enums.Author,
	}

	if err := db.Create(&author).Error; err != nil {
		return err
	}

	log.Printf("Author created: %s with success", authorEmail)
	return nil
}
