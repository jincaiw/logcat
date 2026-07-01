package models

import (
	"time"

	"gorm.io/gorm"
)

// Robot 推送机器人/通道（飞书/邮箱/Syslog/HTTP）
type Robot struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	Name        string `json:"name" gorm:"size:100;not null"`
	Platform    string `json:"platform" gorm:"size:20;default:'feishu'"`
	Description string `json:"description" gorm:"size:500"`
	IsActive    bool   `json:"isActive" gorm:"default:true"`
	// 飞书配置
	FeishuWebhookURL string `json:"feishuWebhookUrl" gorm:"size:500"`
	FeishuSecret     string `json:"feishuSecret" gorm:"size:200"`
	// 邮箱配置
	SMTPHost     string `json:"smtpHost" gorm:"size:100"`
	SMTPPort     int    `json:"smtpPort"`
	SMTPUsername string `json:"smtpUsername" gorm:"size:100"`
	SMTPPassword string `json:"smtpPassword" gorm:"size:100"`
	SMTPFrom     string `json:"smtpFrom" gorm:"size:100"`
	SMTPTo       string `json:"smtpTo" gorm:"size:500"`
	// Syslog 推送配置
	SyslogHost     string `json:"syslogHost" gorm:"size:100"`
	SyslogPort     int    `json:"syslogPort"`
	SyslogProtocol string `json:"syslogProtocol" gorm:"size:10;default:'udp'"`
	SyslogFormat   string `json:"syslogFormat" gorm:"size:20;default:'json'"`
	// HTTP 推送配置
	HTTPURL        string         `json:"httpUrl" gorm:"size:500"`
	HTTPTimeout    int            `json:"httpTimeout" gorm:"default:3"`
	HTTPRetryCount int            `json:"httpRetryCount" gorm:"default:3"`
	HTTPRetryDelay int            `json:"httpRetryDelay" gorm:"default:2"`
	HTTPNotesIDs   string         `json:"httpNotesIds" gorm:"size:500"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// AlertRecord 告警发送记录
type AlertRecord struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	LogID         uint      `json:"logId" gorm:"index"`
	RobotID       uint      `json:"robotId" gorm:"index"`
	AlertPolicyID uint      `json:"alertPolicyId" gorm:"index"`
	DeviceName    string    `json:"deviceName" gorm:"size:100"`
	Message       string    `json:"message" gorm:"type:text"`
	Status        string    `json:"status" gorm:"size:20;index"`
	ErrorMsg      string    `json:"errorMsg" gorm:"type:text"`
	SentAt        time.Time `json:"sentAt" gorm:"index"`
}
