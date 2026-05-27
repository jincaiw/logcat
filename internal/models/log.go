package models

import "time"

// SyslogLog represents a received syslog entry
type SyslogLog struct {
	ID                      uint      `gorm:"primaryKey" json:"id"`
	LogID                   string    `gorm:"uniqueIndex;size:64;not null" json:"logId"`
	SourceIP                string    `gorm:"size:45;index" json:"sourceIp"`
	DestinationIP           string    `gorm:"size:45" json:"destinationIp"`
	EventType               string    `gorm:"size:64;index" json:"eventType"`
	Severity                string    `gorm:"size:20;index" json:"severity"`
	Facility                string    `gorm:"size:20" json:"facility"`
	DeviceID                *uint     `json:"deviceId"`
	DeviceName              string    `gorm:"size:128" json:"deviceName"`
	RawMessage              string    `gorm:"type:text" json:"rawMessage"`
	ParsedData              string    `gorm:"type:text" json:"parsedData"`
	FilterStatus            string    `gorm:"size:32;default:pending;index" json:"filterStatus"`
	MatchedFilterPolicyID   *uint     `json:"matchedFilterPolicyId"`
	AlertStatus             string    `gorm:"size:32" json:"alertStatus"`
	AlertRuleID             *uint     `json:"alertRuleId"`
	AggregatedAlertID       *uint     `json:"aggregatedAlertId"`
	ReceivedAt              time.Time `gorm:"index" json:"receivedAt"`
	CreatedAt               time.Time `json:"createdAt"`
}

// TableName for SyslogLog
func (SyslogLog) TableName() string {
	return "syslog_logs"
}

// LogTraceInfo stores the trace information for a log entry
type LogTraceInfo struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	LogID             string    `gorm:"uniqueIndex;size:64;not null" json:"logId"`
	ReceiveStatus     string    `gorm:"size:32;default:success" json:"receiveStatus"`
	ParseStatus       string    `gorm:"size:32" json:"parseStatus"`
	ParseTemplateID   *uint     `json:"parseTemplateId"`
	ParseTemplateName string    `gorm:"size:128" json:"parseTemplateName"`
	ParseResult       string    `gorm:"type:text" json:"parseResult"`
	ParseError        string    `gorm:"type:text" json:"parseError"`
	FilterStatus      string    `gorm:"size:32" json:"filterStatus"`
	MatchedPolicyID   *uint     `json:"matchedPolicyId"`
	MatchedPolicyName string    `gorm:"size:128" json:"matchedPolicyName"`
	DedupResult       string    `gorm:"size:32" json:"dedupResult"`
	WhitelistResult   string    `gorm:"size:32" json:"whitelistResult"`
	AggregationResult string    `gorm:"type:text" json:"aggregationResult"`
	PushResults       string    `gorm:"type:text" json:"pushResults"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

// TableName for LogTraceInfo
func (LogTraceInfo) TableName() string {
	return "log_trace_infos"
}