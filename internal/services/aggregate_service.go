package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// AggregateService handles log aggregation
type AggregateService struct{}

// NewAggregateService creates a new AggregateService
func NewAggregateService() *AggregateService {
	return &AggregateService{}
}

// AggregateKey generates a key for log aggregation
func (s *AggregateService) AggregateKey(sourceIP, destinationIP, eventType, severity string) string {
	return fmt.Sprintf("%s|%s|%s|%s", sourceIP, destinationIP, eventType, severity)
}

// Aggregate aggregates a log entry into an AggregatedAlert
func (s *AggregateService) Aggregate(sourceIP, destinationIP, eventType, severity string, deviceID *uint) (*models.AggregatedAlert, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	key := s.AggregateKey(sourceIP, destinationIP, eventType, severity)

	var existing models.AggregatedAlert
	result := db.Where("aggregate_key = ? AND status = ?", key, "active").First(&existing)

	if result.Error != nil {
		// Create new aggregated alert
		now := time.Now()
		alert := models.AggregatedAlert{
			AggregateKey:  key,
			AggregateType: "basic",
			SourceIP:      sourceIP,
			DestinationIP: destinationIP,
			EventType:     eventType,
			Severity:      severity,
			DeviceID:      deviceID,
			Count:         1,
			FirstSeenAt:   now,
			LastSeenAt:    now,
			Status:        "active",
		}
		if err := db.Create(&alert).Error; err != nil {
			return nil, err
		}
		return &alert, nil
	}

	// Update existing
	existing.Count++
	existing.LastSeenAt = time.Now()
	if err := db.Save(&existing).Error; err != nil {
		return nil, err
	}

	return &existing, nil
}

// ListAggregatedAlerts returns paginated aggregated alerts
func (s *AggregateService) ListAggregatedAlerts(page, pageSize int, status, severity, eventType string) ([]models.AggregatedAlert, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New("database not available")
	}

	var alerts []models.AggregatedAlert
	var total int64

	query := db.Model(&models.AggregatedAlert{})

	if status != "" {
		query = query.Where("status = ?", status)
	}
	if severity != "" {
		query = query.Where("severity = ?", severity)
	}
	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("last_seen_at DESC").Find(&alerts).Error; err != nil {
		return nil, 0, err
	}

	return alerts, total, nil
}

// GetAggregatedAlert returns a single aggregated alert
func (s *AggregateService) GetAggregatedAlert(id uint) (*models.AggregatedAlert, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var alert models.AggregatedAlert
	if err := db.First(&alert, id).Error; err != nil {
		return nil, err
	}
	return &alert, nil
}

// UpdateStatus updates the status of an aggregated alert
func (s *AggregateService) UpdateStatus(id uint, status string) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Model(&models.AggregatedAlert{}).Where("id = ?", id).Update("status", status).Error
}