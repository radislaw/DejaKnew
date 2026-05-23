package repository

import (
	"fmt"

	"flashcards/api/internal/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect opens a GORM connection to PostgreSQL.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("connect to postgres: %w", err)
	}
	return db, nil
}

// Migrate applies the schema via GORM AutoMigrate.
// New domain models must be registered here.
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.User{}, &domain.Deck{}, &domain.Card{}, &domain.StudyProgress{}, &domain.ReviewLog{}); err != nil {
		return fmt.Errorf("auto-migrate: %w", err)
	}
	return nil
}
