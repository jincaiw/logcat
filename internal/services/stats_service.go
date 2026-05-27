package services

import (
	"errors"
	"runtime"
	"time"

	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
)

// StatsService handles statistics computation
type StatsService struct {
	startTime time.Time
}

// NewStatsService creates a new StatsService
func NewStatsService() *StatsService {
	return &StatsService{
		startTime: time.Now(),
	}
}

// FieldStats represents statistics for a field
type FieldStats struct {
	Field      string            `json:"field"`
	TotalCount int64             `json:"totalCount"`
	TopValues  []ValueCount      `json:"topValues"`
}

// ValueCount represents a value and its count
type ValueCount struct {
	Value string `json:"value"`
	Count int64  `json:"count"`
}

// GetFieldStats returns statistics for available fields
func (s *StatsService) GetFieldStats() ([]FieldStats, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	var stats []FieldStats

	// Field: event_type
	var eventTypeCount int64
	db.Model(&models.SyslogLog{}).Count(&eventTypeCount)
	var topEventTypes []ValueCount
	db.Model(&models.SyslogLog{}).
		Select("event_type as value, count(*) as count").
		Where("event_type != ''").
		Group("event_type").Order("count DESC").Limit(10).
		Scan(&topEventTypes)

	stats = append(stats, FieldStats{
		Field:      "event_type",
		TotalCount: eventTypeCount,
		TopValues:  topEventTypes,
	})

	// Field: severity
	var severityCount int64
	db.Model(&models.SyslogLog{}).Count(&severityCount)
	var topSeverities []ValueCount
	db.Model(&models.SyslogLog{}).
		Select("severity as value, count(*) as count").
		Where("severity != ''").
		Group("severity").Order("count DESC").Limit(10).
		Scan(&topSeverities)

	stats = append(stats, FieldStats{
		Field:      "severity",
		TotalCount: severityCount,
		TopValues:  topSeverities,
	})

	// Field: source_ip
	var sourceIPCount int64
	db.Model(&models.SyslogLog{}).Count(&sourceIPCount)
	var topSourceIPs []ValueCount
	db.Model(&models.SyslogLog{}).
		Select("source_ip as value, count(*) as count").
		Where("source_ip != ''").
		Group("source_ip").Order("count DESC").Limit(10).
		Scan(&topSourceIPs)

	stats = append(stats, FieldStats{
		Field:      "source_ip",
		TotalCount: sourceIPCount,
		TopValues:  topSourceIPs,
	})

	// Field: filter_status
	var filterCount int64
	db.Model(&models.SyslogLog{}).Count(&filterCount)
	var topFilterStatuses []ValueCount
	db.Model(&models.SyslogLog{}).
		Select("filter_status as value, count(*) as count").
		Where("filter_status != ''").
		Group("filter_status").Order("count DESC").Limit(10).
		Scan(&topFilterStatuses)

	stats = append(stats, FieldStats{
		Field:      "filter_status",
		TotalCount: filterCount,
		TopValues:  topFilterStatuses,
	})

	return stats, nil
}

// AvailableFields returns the list of available fields for stats
func (s *StatsService) AvailableFields() []string {
	return []string{
		"event_type",
		"severity",
		"source_ip",
		"destination_ip",
		"filter_status",
		"device_name",
		"facility",
	}
}

// DashboardStats holds all dashboard metrics
type DashboardStats struct {
	TotalLogs       int64            `json:"totalLogs"`
	TotalAlerts     int64            `json:"totalAlerts"`
	ActiveDevices   int64            `json:"activeDevices"`
	ActivePolicies  int64            `json:"activePolicies"`
	LogsByHour      []CountByTime    `json:"logsByHour"`
	AlertsByStatus  []ValueCount     `json:"alertsByStatus"`
	LogsBySeverity  []ValueCount     `json:"logsBySeverity"`
}

// CountByTime represents a time-based count
type CountByTime struct {
	Time  string `json:"time"`
	Count int64  `json:"count"`
}

// GetDashboardStats returns all dashboard metrics
func (s *StatsService) GetDashboardStats() (*DashboardStats, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	stats := &DashboardStats{}

	// Total logs
	db.Model(&models.SyslogLog{}).Count(&stats.TotalLogs)

	// Total alerts
	db.Model(&models.AlertRecord{}).Count(&stats.TotalAlerts)

	// Active devices
	db.Model(&models.Device{}).Where("enabled = ?", true).Count(&stats.ActiveDevices)

	// Active policies
	db.Model(&models.FilterPolicy{}).Where("enabled = ?", true).Count(&stats.ActivePolicies)

	// Alerts by status
	db.Model(&models.AlertRecord{}).
		Select("status as value, count(*) as count").
		Group("status").Scan(&stats.AlertsByStatus)

	// Logs by severity
	db.Model(&models.SyslogLog{}).
		Select("severity as value, count(*) as count").
		Where("severity != ''").
		Group("severity").Scan(&stats.LogsBySeverity)

	return stats, nil
}

// RuntimeMetrics returns Go runtime metrics
func (s *StatsService) RuntimeMetrics() map[string]interface{} {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	return map[string]interface{}{
		"goroutines":      runtime.NumGoroutine(),
		"heap_alloc_mb":   float64(memStats.HeapAlloc) / 1024 / 1024,
		"heap_total_mb":   float64(memStats.TotalAlloc) / 1024 / 1024,
		"heap_sys_mb":     float64(memStats.HeapSys) / 1024 / 1024,
		"num_gc":          memStats.NumGC,
		"uptime_seconds":  int64(time.Since(s.startTime).Seconds()),
	}
}

// GetSystemStatus returns overall system status
func (s *StatsService) GetSystemStatus() map[string]interface{} {
	db := database.GetDB()
	dbStatus := "connected"
	if db == nil {
		dbStatus = "disconnected"
	}

	var dbType string
	if db != nil {
		var result struct{ Result string }
		db.Raw("SELECT 'sqlite' as result").Scan(&result)
		dbType = "sqlite"

		// Try MySQL check
		var dbName string
		if db.Raw("SELECT DATABASE()").Scan(&dbName).Error == nil && dbName != "" {
			dbType = "mysql"
		}
	}

	var totalLogs, totalAlerts int64
	if db != nil {
		db.Model(&models.SyslogLog{}).Count(&totalLogs)
		db.Model(&models.AlertRecord{}).Count(&totalAlerts)
	}

	return map[string]interface{}{
		"status":          "running",
		"version":         "0.1.0",
		"database":        dbType,
		"database_status": dbStatus,
		"total_logs":      totalLogs,
		"total_alerts":    totalAlerts,
		"uptime_seconds":  int64(time.Since(s.startTime).Seconds()),
		"go_version":      runtime.Version(),
	}
}