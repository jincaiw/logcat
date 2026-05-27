package models

import "time"

// SystemConfig stores system configuration key-value pairs
type SystemConfig struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ConfigKey   string    `gorm:"uniqueIndex;size:128;not null" json:"configKey"`
	ConfigValue string    `gorm:"type:text" json:"configValue"`
	Description string    `gorm:"size:512" json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName for SystemConfig
func (SystemConfig) TableName() string {
	return "system_configs"
}

// MetricSnapshot represents a system metric snapshot
type MetricSnapshot struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	CPUPercent    float64   `gorm:"type:decimal(5,2)" json:"cpuPercent"`
	MemoryMB      float64   `gorm:"type:decimal(10,2)" json:"memoryMb"`
	DiskMB        float64   `gorm:"type:decimal(10,2)" json:"diskMb"`
	Goroutines    int       `json:"goroutines"`
	DBConnections int       `json:"dbConnections"`
	UptimeSeconds int64     `json:"uptimeSeconds"`
	LogsReceived  int64     `json:"logsReceived"`
	LogsProcessed int64     `json:"logsProcessed"`
	AlertsSent    int64     `json:"alertsSent"`
	CreatedAt     time.Time `json:"createdAt"`
}

// TableName for MetricSnapshot
func (MetricSnapshot) TableName() string {
	return "metric_snapshots"
}