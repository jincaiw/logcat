package models

import "time"

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       *uint     `json:"userId"`
	Username     string    `gorm:"size:64" json:"username"`
	Action       string    `gorm:"size:128;not null;index" json:"action"`
	ResourceType string    `gorm:"size:64;index" json:"resourceType"`
	ResourceID   string    `gorm:"size:64" json:"resourceId"`
	ClientIP     string    `gorm:"size:45" json:"clientIp"`
	UserAgent    string    `gorm:"size:512" json:"userAgent"`
	Result       string    `gorm:"size:20;not null" json:"result"`
	Detail       string    `gorm:"type:text" json:"detail"`
	CreatedAt    time.Time `gorm:"index" json:"createdAt"`
}

// TableName for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}

// ExportHistory represents an export operation history record
type ExportHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       *uint     `json:"userId"`
	Username     string    `gorm:"size:64" json:"username"`
	ResourceType string    `gorm:"size:64;not null" json:"resourceType"`
	Count        int       `json:"count"`
	ClientIP     string    `gorm:"size:45" json:"clientIp"`
	Result       string    `gorm:"size:20;not null" json:"result"`
	Detail       string    `gorm:"type:text" json:"detail"`
	CreatedAt    time.Time `gorm:"index" json:"createdAt"`
}

func (ExportHistory) TableName() string {
	return "export_histories"
}

// ImportHistory represents an import operation history record
type ImportHistory struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       *uint     `json:"userId"`
	Username     string    `gorm:"size:64" json:"username"`
	ResourceType string    `gorm:"size:64;not null" json:"resourceType"`
	Created      int       `json:"created"`
	Updated      int       `json:"updated"`
	Failed       int       `json:"failed"`
	Errors       string    `gorm:"type:text" json:"errors"`
	ClientIP     string    `gorm:"size:45" json:"clientIp"`
	Result       string    `gorm:"size:20;not null" json:"result"`
	Detail       string    `gorm:"type:text" json:"detail"`
	CreatedAt    time.Time `gorm:"index" json:"createdAt"`
}

func (ImportHistory) TableName() string {
	return "import_histories"
}