package models

import "time"

// AlertRule defines a rule that triggers alerts
type AlertRule struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:128;not null" json:"name"`
	FilterPolicyID   *uint     `json:"filterPolicyId"`
	PushConfigID     *uint     `json:"pushConfigId"`
	OutputTemplateID *uint     `json:"outputTemplateId"`
	ChannelType      string    `gorm:"size:32" json:"channelType"`
	Enabled          bool      `gorm:"default:true" json:"enabled"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Associations
	FilterPolicy   *FilterPolicy   `gorm:"foreignKey:FilterPolicyID" json:"filterPolicy,omitempty"`
	PushConfig     *PushConfig     `gorm:"foreignKey:PushConfigID" json:"pushConfig,omitempty"`
	OutputTemplate *OutputTemplate `gorm:"foreignKey:OutputTemplateID" json:"outputTemplate,omitempty"`
}

// TableName for AlertRule
func (AlertRule) TableName() string {
	return "alert_rules"
}

// AlertRecord represents a single alert notification record
type AlertRecord struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	LogID              string     `gorm:"size:64;index" json:"logId"`
	AlertRuleID        *uint      `json:"alertRuleId"`
	PushConfigID       *uint      `json:"pushConfigId"`
	ChannelType        string     `gorm:"size:32" json:"channelType"`
	Status             string     `gorm:"size:20;default:success" json:"status"`
	RetryCount         int        `gorm:"default:0" json:"retryCount"`
	RequestSummary     string     `gorm:"type:text" json:"requestSummary"`
	ResponseStatusCode int        `gorm:"default:0" json:"responseStatusCode"`
	ResponseSummary    string     `gorm:"type:text" json:"responseSummary"`
	ErrorMessage       string     `gorm:"type:text" json:"errorMessage"`
	DispositionStatus  string     `gorm:"size:32" json:"dispositionStatus"`
	SentAt             *time.Time `json:"sentAt"`
	CreatedAt          time.Time  `json:"createdAt"`

	// Associations
	AlertRule *AlertRule `gorm:"foreignKey:AlertRuleID" json:"alertRule,omitempty"`
	PushConfig *PushConfig `gorm:"foreignKey:PushConfigID" json:"pushConfig,omitempty"`
}

// TableName for AlertRecord
func (AlertRecord) TableName() string {
	return "alert_records"
}

// AggregatedAlert represents an aggregated alert entry
type AggregatedAlert struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	AggregateKey    string    `gorm:"size:256;index" json:"aggregateKey"`
	AggregateType   string    `gorm:"size:64" json:"aggregateType"`
	SourceIP        string    `gorm:"size:45" json:"sourceIp"`
	DestinationIP   string    `gorm:"size:45" json:"destinationIp"`
	EventType       string    `gorm:"size:64" json:"eventType"`
	DeviceID        *uint     `json:"deviceId"`
	Severity        string    `gorm:"size:20" json:"severity"`
	Count           int       `gorm:"default:1" json:"count"`
	FirstSeenAt     time.Time `json:"firstSeenAt"`
	LastSeenAt      time.Time `json:"lastSeenAt"`
	Status          string    `gorm:"size:32;default:active" json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// TableName for AggregatedAlert
func (AggregatedAlert) TableName() string {
	return "aggregated_alerts"
}

// AlertDisposition represents a disposition record for an alert
type AlertDisposition struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	AlertRecordID    *uint     `json:"alertRecordId"`
	AggregatedAlertID *uint    `json:"aggregatedAlertId"`
	Status           string    `gorm:"size:32;default:unprocessed" json:"status"`
	Note             string    `gorm:"type:text" json:"note"`
	OperatorID       *uint     `json:"operatorId"`
	OperatorName     string    `gorm:"size:128" json:"operatorName"`
	OperatedAt       *time.Time `json:"operatedAt"`
	CreatedAt        time.Time `json:"createdAt"`
}

// TableName for AlertDisposition
func (AlertDisposition) TableName() string {
	return "alert_dispositions"
}