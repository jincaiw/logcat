package services

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/logcat/logcat/internal/config"
	"github.com/logcat/logcat/internal/database"
	"github.com/logcat/logcat/internal/models"
	logsyslog "github.com/logcat/logcat/internal/syslog"
	"gorm.io/gorm"
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

// QueryStatsItem represents one row in the statistics result set.
type QueryStatsItem struct {
	Field      string  `json:"field"`
	Value      string  `json:"value"`
	Count      int64   `json:"count"`
	Percentage float64 `json:"percentage"`
	LastSeenAt string  `json:"lastSeenAt"`
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
	ServiceStatus      string             `json:"serviceStatus"`
	ReceiveRate        float64            `json:"receiveRate"`
	TodayTotal         int64              `json:"todayTotal"`
	TodayAlerts        int64              `json:"todayAlerts"`
	PushSuccessRate    float64            `json:"pushSuccessRate"`
	QueueBacklog       []QueueMetric      `json:"queueBacklog"`
	RecentPushFailures []PushFailureItem  `json:"recentPushFailures"`
	HealthStatus       RuntimeHealthStats `json:"healthStatus"`
}

// CountByTime represents a time-based count
type CountByTime struct {
	Time  string `json:"time"`
	Count int64  `json:"count"`
}

// QueueMetric is the queue status shown on the dashboard and status page.
type QueueMetric struct {
	Name        string  `json:"name"`
	Size        int     `json:"size"`
	Capacity    int     `json:"capacity"`
	EnqueueRate float64 `json:"enqueueRate"`
	DequeueRate float64 `json:"dequeueRate"`
}

// PushFailureItem represents a recent failed push summary.
type PushFailureItem struct {
	Time    string `json:"time"`
	Channel string `json:"channel"`
	Error   string `json:"error"`
	LogID   string `json:"logId"`
}

// RuntimeHealthStats is a lightweight runtime health snapshot.
type RuntimeHealthStats struct {
	Service    string  `json:"service"`
	CPU        float64 `json:"cpu"`
	Memory     float64 `json:"memory"`
	DiskUsage  float64 `json:"diskUsage"`
	NetworkIn  float64 `json:"networkIn"`
	NetworkOut float64 `json:"networkOut"`
}

func normalizeStatsField(field string) (string, bool) {
	switch strings.TrimSpace(field) {
	case "sourceIp", "source_ip":
		return "source_ip", true
	case "destinationIp", "destination_ip", "destIp", "dest_ip":
		return "destination_ip", true
	case "deviceName", "device_name":
		return "device_name", true
	case "eventType", "event_type":
		return "event_type", true
	case "severity":
		return "severity", true
	case "filterStatus", "filter_status":
		return "filter_status", true
	case "facility":
		return "facility", true
	default:
		return "", false
	}
}

func parseFlexibleTime(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}

	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02 15:04",
		"2006-01-02",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("invalid time format: %s", value)
}

func applyTimeRange(query *gorm.DB, startTime, endTime string) (*gorm.DB, error) {
	start, err := parseFlexibleTime(startTime)
	if err != nil {
		return nil, err
	}
	end, err := parseFlexibleTime(endTime)
	if err != nil {
		return nil, err
	}

	if start != nil {
		query = query.Where("received_at >= ?", *start)
	}
	if end != nil {
		query = query.Where("received_at <= ?", *end)
	}

	return query, nil
}

// QueryFieldStats returns field aggregation results for a chosen field and time window.
func (s *StatsService) QueryFieldStats(field, startTime, endTime string, topN int) ([]QueryStatsItem, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	column, ok := normalizeStatsField(field)
	if !ok {
		return nil, errors.New("unsupported stats field")
	}
	if topN <= 0 {
		topN = 20
	}

	query := db.Model(&models.SyslogLog{})
	query, err := applyTimeRange(query, startTime, endTime)
	if err != nil {
		return nil, err
	}

	var total int64
	if err := query.Where(fmt.Sprintf("%s <> ''", column)).Count(&total).Error; err != nil {
		return nil, err
	}
	if total == 0 {
		return []QueryStatsItem{}, nil
	}

	type row struct {
		Value      string
		Count      int64
		LastSeenAt time.Time
	}
	rows := make([]row, 0, topN)
	if err := query.
		Select(fmt.Sprintf("%s as value, count(*) as count, max(received_at) as last_seen_at", column)).
		Where(fmt.Sprintf("%s <> ''", column)).
		Group(column).
		Order("count DESC").
		Limit(topN).
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]QueryStatsItem, 0, len(rows))
	for _, item := range rows {
		result = append(result, QueryStatsItem{
			Field:      field,
			Value:      item.Value,
			Count:      item.Count,
			Percentage: float64(item.Count) / float64(total),
			LastSeenAt: item.LastSeenAt.Format(time.RFC3339),
		})
	}

	return result, nil
}

// GetIPList returns distinct source IPs within the time range.
func (s *StatsService) GetIPList(startTime, endTime string) ([]string, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	query := db.Model(&models.SyslogLog{})
	query, err := applyTimeRange(query, startTime, endTime)
	if err != nil {
		return nil, err
	}

	ips := make([]string, 0)
	if err := query.Where("source_ip <> ''").Distinct("source_ip").Order("source_ip ASC").Pluck("source_ip", &ips).Error; err != nil {
		return nil, err
	}

	return ips, nil
}

// GetDashboardStats returns all dashboard metrics
func (s *StatsService) GetDashboardStats() (*DashboardStats, error) {
	db := database.GetDB()
	if db == nil {
		return nil, errors.New("database not available")
	}

	stats := &DashboardStats{
		ServiceStatus: "running",
		QueueBacklog: []QueueMetric{
			{Name: "log-receive", Capacity: 0},
			{Name: "push", Capacity: 0},
		},
		RecentPushFailures: []PushFailureItem{},
		HealthStatus: RuntimeHealthStats{
			Service: "running",
		},
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	db.Model(&models.SyslogLog{}).Where("received_at >= ?", startOfDay).Count(&stats.TodayTotal)
	db.Model(&models.AlertRecord{}).Where("created_at >= ?", startOfDay).Count(&stats.TodayAlerts)

	var totalAlerts int64
	var successAlerts int64
	db.Model(&models.AlertRecord{}).Count(&totalAlerts)
	db.Model(&models.AlertRecord{}).Where("status = ?", "success").Count(&successAlerts)
	if totalAlerts > 0 {
		stats.PushSuccessRate = float64(successAlerts) / float64(totalAlerts)
	}

	receiver := logsyslog.GetGlobalReceiver()
	if receiver != nil {
		metrics := receiver.Metrics()
		uptimeSeconds := time.Since(s.startTime).Seconds()
		if uptimeSeconds > 0 {
			stats.ReceiveRate = float64(metrics.UDPReceived+metrics.TCPReceived) / uptimeSeconds
		}

		queueCapacity := 0
		if cfg := config.Get(); cfg != nil {
			queueCapacity = cfg.Queue.Capacity
		}
		stats.QueueBacklog = []QueueMetric{
			{
				Name:        "log-receive",
				Size:        int(metrics.ChannelDropped),
				Capacity:    queueCapacity,
				EnqueueRate: stats.ReceiveRate,
				DequeueRate: stats.ReceiveRate,
			},
			{
				Name:        "push",
				Size:        0,
				Capacity:    queueCapacity,
				EnqueueRate: 0,
				DequeueRate: 0,
			},
		}

		stats.HealthStatus.NetworkIn = float64(metrics.UDPReceived + metrics.TCPReceived)
	}

	type failedRecord struct {
		LogID       string
		ChannelType string
		Error       string
		CreatedAt   time.Time
	}
	failures := make([]failedRecord, 0, 5)
	if err := db.Model(&models.AlertRecord{}).
		Select("log_id, channel_type, error_message as error, created_at").
		Where("status = ?", "failed").
		Order("created_at DESC").
		Limit(5).
		Scan(&failures).Error; err == nil {
		for _, item := range failures {
			stats.RecentPushFailures = append(stats.RecentPushFailures, PushFailureItem{
				Time:    item.CreatedAt.Format(time.RFC3339),
				Channel: item.ChannelType,
				Error:   item.Error,
				LogID:   item.LogID,
			})
		}
	}

	runtimeMetrics := s.RuntimeMetrics()
	if heapMB, ok := runtimeMetrics["heap_alloc_mb"].(float64); ok {
		stats.HealthStatus.Memory = heapMB
	}

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

// GetSystemStatus returns overall system status.
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

	cfg := config.Get()
	queueCapacity := 0
	parseWorkers := 0
	filterWorkers := 0
	pushWorkers := 0
	listeners := make([]map[string]interface{}, 0, 2)
	if cfg != nil {
		queueCapacity = cfg.Queue.Capacity
		parseWorkers = cfg.Worker.ParseWorkers
		filterWorkers = cfg.Worker.FilterWorkers
		pushWorkers = cfg.Worker.PushWorkers
		if cfg.Syslog.UDPPort > 0 {
			listeners = append(listeners, map[string]interface{}{
				"protocol": "udp",
				"port":     cfg.Syslog.UDPPort,
				"address":  "0.0.0.0",
				"status":   "listening",
			})
		}
		if cfg.Syslog.TCPPort > 0 {
			listeners = append(listeners, map[string]interface{}{
				"protocol": "tcp",
				"port":     cfg.Syslog.TCPPort,
				"address":  "0.0.0.0",
				"status":   "listening",
			})
		}
	}

	receiver := logsyslog.GetGlobalReceiver()
	running := receiver != nil && receiver.IsRunning()
	tcpConnections := int64(0)
	lastReceiveAt := ""
	receiverMetrics := logsyslog.ReceiverMetrics{}
	if receiver != nil {
		receiverMetrics = receiver.Metrics()
		tcpConnections = receiverMetrics.TCPConnections
		if receiverMetrics.LastReceiveAt > 0 {
			lastReceiveAt = time.Unix(0, receiverMetrics.LastReceiveAt).Format(time.RFC3339)
		}
	}

	workers := []map[string]interface{}{
		{
			"id":             1,
			"stage":          "parse",
			"count":          parseWorkers,
			"status":         workerPoolStatus(parseWorkers, running),
			"processedCount": 0,
			"errorCount":     0,
			"lastActiveAt":   lastReceiveAt,
		},
		{
			"id":             2,
			"stage":          "filter",
			"count":          filterWorkers,
			"status":         workerPoolStatus(filterWorkers, running),
			"processedCount": 0,
			"errorCount":     0,
			"lastActiveAt":   lastReceiveAt,
		},
		{
			"id":             3,
			"stage":          "push",
			"count":          pushWorkers,
			"status":         workerPoolStatus(pushWorkers, running),
			"processedCount": 0,
			"errorCount":     0,
			"lastActiveAt":   lastReceiveAt,
		},
	}

	for i := range listeners {
		if !running {
			listeners[i]["status"] = "stopped"
		}
	}

	return map[string]interface{}{
		"status":         "running",
		"serviceRunning": running,
		"startedAt":      s.startTime.Format(time.RFC3339),
		"uptime":         formatDuration(time.Since(s.startTime)),
		"uptimeSeconds":  int64(time.Since(s.startTime).Seconds()),
		"version":        "0.1.0",
		"database":       dbType,
		"databaseStatus": dbStatus,
		"goVersion":      runtime.Version(),
		"listeners":      listeners,
		"connections": map[string]interface{}{
			"total":  tcpConnections,
			"active": tcpConnections,
			"idle":   0,
			"closed": 0,
		},
		"queue": map[string]interface{}{
			"name":        "log-receive",
			"size":        0,
			"capacity":    queueCapacity,
			"enqueueRate": 0,
			"dequeueRate": 0,
		},
		"workers": workers,
		"receiverMetrics": map[string]interface{}{
			"udpReceived":    receiverMetrics.UDPReceived,
			"tcpReceived":    receiverMetrics.TCPReceived,
			"udpErrors":      receiverMetrics.UDPErrors,
			"tcpErrors":      receiverMetrics.TCPErrors,
			"parseErrors":    receiverMetrics.ParseErrors,
			"channelDropped": receiverMetrics.ChannelDropped,
			"tcpConnections": receiverMetrics.TCPConnections,
			"lastReceiveAt":  lastReceiveAt,
		},
	}
}

func workerPoolStatus(count int, running bool) string {
	if count <= 0 {
		return "disabled"
	}
	if running {
		return "running"
	}
	return "stopped"
}

func formatDuration(d time.Duration) string {
	totalSeconds := int64(d.Seconds())
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60
	seconds := totalSeconds % 60

	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}
