package database

import (
	"log"
)

// AutoMigrate runs GORM auto-migration for all models (idempotent)
func AutoMigrate() error {
	if DB == nil {
		return nil
	}

	log.Println("Running auto-migration...")
	if err := DB.AutoMigrate(allModels()...); err != nil {
		return err
	}
	log.Println("Auto-migration completed successfully")
	return nil
}