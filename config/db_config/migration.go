package db_config

import (
	"github.com/clemilsonazevedo/blog/internal/domain/entities"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(entities.RetrieveAll()...)
	if err != nil {
		return
	}
}

func DropAll(db *gorm.DB) {
	db.Migrator().DropTable(entities.RetrieveAll()...)
}
