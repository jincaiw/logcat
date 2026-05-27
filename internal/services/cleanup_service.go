package services

import (
	"errors"
	"log"
	"time"

	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// CleanupService handles periodic log and alert cleanup
type CleanupService struct {
	stopCh chan struct{}
}

// NewCleanupService creates a new CleanupService
func NewCleanupService() *CleanupService {
	return &CleanupService{
		stopCh: make(chan struct{}),
	}
}

// Start begins periodic cleanup
func (s *CleanupService) Start(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				s.runCleanup()
			case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop stops the cleanup service
func (s *CleanupService) Stop() {
	close(s.stopCh)
}

// RunCleanupNow runs cleanup immediately
func (s *CleanupService) RunCleanupNow() error {
	return s.runCleanup()
}

func (s *CleanupService) runCleanup() error {
	cfg := config.Get()
	if cfg == nil {
		return errors.New("configuration not loaded")
	}

	db := database.GetDB()
	if db == nil {
		return errors.New("database not available")
	}

	// Cleanup old syslog logs
	retentionDays := cfg.Log.RetentionDays
	cutoffTime := time.Now().AddDate(0, 0, -retentionDays)

	result := db.Where("received_at < ?", cutoffTime).Delete(&models.SyslogLog{})
	if result.Error != nil {
		log.Printf("Cleanup: failed to delete old syslog logs: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("Cleanup: deleted %d old syslog logs (older than %d days)", result.RowsAffected, retentionDays)
	}

	// Cleanup old unmatched logs (shorter retention)
	unmatchedRetention := cfg.Log.UnmatchedRetentionDays
	unmatchedCutoff := time.Now().AddDate(0, 0, -unmatchedRetention)

	unmatchedResult := db.Where("filter_status = ? AND received_at < ?", "unmatched", unmatchedCutoff).
		Delete(&models.SyslogLog{})
	if unmatchedResult.Error != nil {
		log.Printf("Cleanup: failed to delete unmatched logs: %v", unmatchedResult.Error)
	} else if unmatchedResult.RowsAffected > 0 {
		log.Printf("Cleanup: deleted %d unmatched logs (older than %d days)", unmatchedResult.RowsAffected, unmatchedRetention)
	}

	// Cleanup old alert records
	alertCutoff := time.Now().AddDate(0, 0, -retentionDays)
	alertResult := db.Where("created_at < ?", alertCutoff).Delete(&models.AlertRecord{})
	if alertResult.Error != nil {
		log.Printf("Cleanup: failed to delete old alert records: %v", alertResult.Error)
	} else if alertResult.RowsAffected > 0 {
		log.Printf("Cleanup: deleted %d old alert records", alertResult.RowsAffected)
	}

	// Cleanup old trace info
	traceResult := db.Where("created_at < ?", cutoffTime).Delete(&models.LogTraceInfo{})
	if traceResult.Error != nil {
		log.Printf("Cleanup: failed to delete old trace info: %v", traceResult.Error)
	} else if traceResult.RowsAffected > 0 {
		log.Printf("Cleanup: deleted %d old trace records", traceResult.RowsAffected)
	}

	// Cleanup old audit logs (keep for 365 days)
	auditCutoff := time.Now().AddDate(0, 0, -365)
	auditResult := db.Where("created_at < ?", auditCutoff).Delete(&models.AuditLog{})
	if auditResult.Error != nil {
		log.Printf("Cleanup: failed to delete old audit logs: %v", auditResult.Error)
	} else if auditResult.RowsAffected > 0 {
		log.Printf("Cleanup: deleted %d old audit logs", auditResult.RowsAffected)
	}

	// Cleanup old metric snapshots (keep for 30 days)
	metricCutoff := time.Now().AddDate(0, 0, -30)
	metricResult := db.Where("created_at < ?", metricCutoff).Delete(&models.MetricSnapshot{})
	if metricResult.Error != nil {
		log.Printf("Cleanup: failed to delete old metrics: %v", metricResult.Error)
	}

	return nil
}