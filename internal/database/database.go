package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB is the global database instance
var DB *gorm.DB

// Initialize opens the database connection based on configuration
func Initialize(cfg *config.Config) error {
	var dialector gorm.Dialector

	switch cfg.Database.Type {
	case "mysql":
		dsn := cfg.DSN()
		dialector = mysql.Open(dsn)
	case "sqlite":
		// Ensure the directory for the SQLite file exists
		dbPath := cfg.Database.SQLite.Path
		dir := filepath.Dir(dbPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create SQLite data directory %s: %w", dir, err)
		}
		dialector = sqlite.Open(dbPath)
	default:
		return fmt.Errorf("unsupported database type: %s (must be sqlite or mysql)", cfg.Database.Type)
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(dialector, gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Enable WAL mode for SQLite
	if cfg.Database.Type == "sqlite" && cfg.Database.SQLite.WAL {
		if err := db.Exec("PRAGMA journal_mode=WAL").Error; err != nil {
			log.Printf("WARNING: Failed to enable WAL mode: %v", err)
		}
		// Additional SQLite performance pragmas
		db.Exec("PRAGMA synchronous=NORMAL")
		db.Exec("PRAGMA cache_size=-8000")
		db.Exec("PRAGMA busy_timeout=5000")
	}

	// Configure connection pool for MySQL
	if cfg.Database.Type == "mysql" {
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("failed to get underlying sql.DB: %w", err)
		}
		sqlDB.SetMaxOpenConns(cfg.Database.MySQL.MaxOpenConns)
		sqlDB.SetMaxIdleConns(cfg.Database.MySQL.MaxIdleConns)
	}

	DB = db
	log.Printf("Database connected: type=%s", cfg.Database.Type)
	return nil
}

// GetDB returns the global database instance
func GetDB() *gorm.DB {
	return DB
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// allModels returns all models for auto-migration
func allModels() []interface{} {
	return []interface{}{
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.DeviceGroup{},
		&models.Device{},
		&models.DeviceTemplate{},
		&models.FieldMappingDoc{},
		&models.ParseTemplate{},
		&models.OutputTemplate{},
		&models.FilterPolicy{},
		&models.PushConfig{},
		&models.AlertRule{},
		&models.SyslogLog{},
		&models.AlertRecord{},
		&models.AggregatedAlert{},
		&models.AlertDisposition{},
		&models.DesensitizeRule{},
		&models.AuditLog{},
		&models.ExportHistory{},
		&models.ImportHistory{},
		&models.SystemConfig{},
		&models.LogTraceInfo{},
		&models.MetricSnapshot{},
	}
}