package database

import (
	"errors"
	"log"
)

// AutoMigrate runs GORM auto-migration for all models (idempotent)
func AutoMigrate() error {
	if DB == nil {
		return errors.New("database not initialized")
	}

	log.Println("Running auto-migration...")
	if err := DB.AutoMigrate(allModels()...); err != nil {
		return err
	}
	log.Println("Auto-migration completed successfully")
	return nil
}