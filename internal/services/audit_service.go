package services

import (
	"errors"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// AuditService handles audit log writing and querying
type AuditService struct{}

// NewAuditService creates a new AuditService
func NewAuditService() *AuditService {
	return &AuditService{}
}

// Create creates an audit log entry
func (s *AuditService) Create(audit *models.AuditLog) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Create(audit).Error
}

// List returns paginated audit log entries
func (s *AuditService) List(page, pageSize int, username, action, result, resourceType string) ([]models.AuditLog, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not available")
	}

	var logs []models.AuditLog
	var total int64

	query := db.Model(&models.AuditLog{})

	if username != "" {
		query = query.Where("username = ?", username)
	}
	if action != "" {
		query = query.Where("action = ?", action)
	}
	if result != "" {
		query = query.Where("result = ?", result)
	}
	if resourceType != "" {
		query = query.Where("resource_type = ?", resourceType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetByID returns a single audit log entry
func (s *AuditService) GetByID(id uint) (*models.AuditLog, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var log models.AuditLog
	if err := db.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}