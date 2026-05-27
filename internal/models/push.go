package models

import "time"

// PushConfig represents a push notification configuration
type PushConfig struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"size:128;not null" json:"name"`
	Type string `gorm:"size:20;not null" json:"type"`
	Enabled bool `gorm:"default:true" json:"enabled"`

	// HTTP fields
	URL                  string `gorm:"type:text" json:"url"`
	Method               string `gorm:"size:10;default:POST" json:"method"`
	Timeout              int    `gorm:"default:30" json:"timeout"`
	RetryCount           int    `gorm:"default:0" json:"retryCount"`
	RetryDelay           int    `gorm:"default:5" json:"retryDelay"`
	NotesIDs             string `gorm:"type:text" json:"notesIds"`
	Headers              string `gorm:"type:text" json:"headers"`
	BodyTemplate         string `gorm:"type:text" json:"bodyTemplate"`
	SuccessStatusCodes   string `gorm:"type:text" json:"successStatusCodes"`
	SuccessBodyKeyword   string `gorm:"size:255" json:"successBodyKeyword"`
	AuthType             string `gorm:"size:20;default:none" json:"authType"`
	Token                string `gorm:"size:512" json:"token"`
	ContentType          string `gorm:"size:128" json:"contentType"`
	RetryOnStatusCodes   string `gorm:"type:text" json:"retryOnStatusCodes"`
	MaxResponseLogSize   int    `gorm:"default:4096" json:"maxResponseLogSize"`

	// Email fields
	SMTPHost          string `gorm:"size:255" json:"smtpHost"`
	SMTPPort          int    `gorm:"default:587" json:"smtpPort"`
	SMTPUsername      string `gorm:"size:255" json:"smtpUsername"`
	SMTPPassword      string `gorm:"size:512" json:"smtpPassword"`
	FromAddress       string `gorm:"size:255" json:"fromAddress"`
	ToAddresses       string `gorm:"type:text" json:"toAddresses"`
	SubjectTemplate   string `gorm:"type:text" json:"subjectTemplate"`
	EmailBodyTemplate string `gorm:"type:text" json:"emailBodyTemplate"`

	// Syslog fields
	SyslogHost     string `gorm:"size:255" json:"syslogHost"`
	SyslogPort     int    `gorm:"default:514" json:"syslogPort"`
	SyslogProtocol string `gorm:"size:10;default:udp" json:"syslogProtocol"`
	SyslogFormat   string `gorm:"size:20;default:rfc3164" json:"syslogFormat"`
	SyslogFields   string `gorm:"type:text" json:"syslogFields"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName for PushConfig
func (PushConfig) TableName() string {
	return "push_configs"
}