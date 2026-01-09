package database

import (
	"fmt"
	"strings"

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

func MigrateRoleEnums(db *gorm.DB) error {
	// 1. Criar o tipo com o nome CORRETO
	err := db.Exec("CREATE TYPE user_role AS ENUM ('anonymous', 'reader', 'author');").Error

	if err != nil && !isTypeExistsError(err) {
		return fmt.Errorf("failed to create enum type: %w", err)
	}

	fmt.Println("Enum type 'user_role' created or already exists")
	return nil
}

func isTypeExistsError(err error) bool {
	// Verificar ambos os nomes possíveis
	errorStr := err.Error()
	return strings.Contains(errorStr, "type \"user_role\" already exists") ||
		strings.Contains(errorStr, "SQLSTATE 42710")
}
