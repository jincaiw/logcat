package models

import "time"

// SyslogLog 接收到的 Syslog 日志
type SyslogLog struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	DeviceID        uint      `json:"deviceId" gorm:"index"`
	DeviceName      string    `json:"deviceName" gorm:"size:100;index"`
	SourceIP        string    `json:"sourceIp" gorm:"size:50;index"`
	RawMessage      string    `json:"rawMessage" gorm:"type:text"`
	ParsedData      string    `json:"parsedData" gorm:"type:text"`
	ParsedFields    string    `json:"parsedFields" gorm:"type:text"`
	FilterStatus    string    `json:"filterStatus" gorm:"size:20;default:'pending'"`
	MatchedPolicyID uint      `json:"matchedPolicyId" gorm:"index"`
	AlertStatus     string    `json:"alertStatus" gorm:"size:20;default:'none'"`
	AlertPolicyID   uint      `json:"alertPolicyId" gorm:"index"`
	Priority        string    `json:"priority" gorm:"size:10"`
	Facility        int       `json:"facility"`
	Severity        int       `json:"severity"`
	Timestamp       time.Time `json:"timestamp" gorm:"index"`
	ReceivedAt      time.Time `json:"receivedAt" gorm:"index"`
	IsProcessed     bool      `json:"isProcessed" gorm:"default:false"`
	IsAlerted       bool      `json:"isAlerted" gorm:"default:false"`
}

// SystemConfig 系统配置
type SystemConfig struct {
	ID                    uint   `json:"id" gorm:"primaryKey"`
	ListenPort            int    `json:"listenPort" gorm:"default:5140"`
	Protocol              string `json:"protocol" gorm:"size:10;default:'udp'"`
	LogRetention          int    `json:"logRetention" gorm:"default:7"`
	MaxLogSize            int64  `json:"maxLogSize" gorm:"default:524288000"`
	AutoStart             bool   `json:"autoStart" gorm:"default:false"`
	MinimizeToTray        bool   `json:"minimizeToTray" gorm:"default:true"`
	AlertEnabled          bool   `json:"alertEnabled" gorm:"default:true"`
	AlertInterval         int    `json:"alertInterval" gorm:"default:60"`
	UnmatchedLogRetention int    `json:"unmatchedLogRetention" gorm:"default:7"`
	UnmatchedLogAlert     bool   `json:"unmatchedLogAlert" gorm:"default:true"`
	DefaultFilterAction   string `json:"defaultFilterAction" gorm:"size:20;default:'keep'"`
	Theme                 string `json:"theme" gorm:"size:20;default:'dark'"`
	Language              string `json:"language" gorm:"size:10;default:'zh-CN'"`
	DataDir               string `json:"dataDir" gorm:"size:500"`
	ConfigDir             string `json:"configDir" gorm:"size:500"`
}
