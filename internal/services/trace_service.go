package services

import (
	"errors"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// TraceService handles log trace collection
type TraceService struct{}

// NewTraceService creates a new TraceService
func NewTraceService() *TraceService {
	return &TraceService{}
}

// GetTraceByLogID returns the full trace information for a log entry
func (s *TraceService) GetTraceByLogID(logID string) (*models.LogTraceInfo, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var trace models.LogTraceInfo
	if err := db.Where("log_id = ?", logID).First(&trace).Error; err != nil {
		return nil, err
	}

	return &trace, nil
}

// CreateTrace creates a new trace entry for a log
func (s *TraceService) CreateTrace(trace *models.LogTraceInfo) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Create(trace).Error
}

// UpdateTrace updates an existing trace entry
func (s *TraceService) UpdateTrace(logID string, updates map[string]interface{}) error {
	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	return db.Model(&models.LogTraceInfo{}).Where("log_id = ?", logID).Updates(updates).Error
}

// GetOrCreateTrace returns an existing trace or creates a new one
func (s *TraceService) GetOrCreateTrace(logID string) (*models.LogTraceInfo, error) {
	trace, err := s.GetTraceByLogID(logID)
	if err == nil {
		return trace, nil
	}

	newTrace := &models.LogTraceInfo{
		LogID:         logID,
		ReceiveStatus: "received",
	}

	if err := s.CreateTrace(newTrace); err != nil {
		return nil, err
	}

	return newTrace, nil
}