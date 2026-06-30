package models

import (
	"time"

	"gorm.io/gorm"
)

// FilterPolicy 筛选策略
type FilterPolicy struct {
	ID              uint           `json:"id" gorm:"primaryKey"`
	Name            string         `json:"name" gorm:"size:100;not null"`
	Description     string         `json:"description" gorm:"size:500"`
	DeviceID        uint           `json:"deviceId" gorm:"index"`
	DeviceGroupID   uint           `json:"deviceGroupId" gorm:"index"`
	ParseTemplateID uint           `json:"parseTemplateId" gorm:"index"`
	Conditions      string         `json:"conditions" gorm:"type:text"`                 // JSON 格式筛选条件
	ConditionLogic  string         `json:"conditionLogic" gorm:"size:10;default:'AND'"` // AND/OR
	Whitelist       string         `json:"whitelist" gorm:"type:text"`                  // JSON 格式白名单
	WhitelistField  string         `json:"whitelistField" gorm:"size:50;default:'sip'"`
	Action          string         `json:"action" gorm:"size:20;default:'keep'"` // keep/discard
	Priority        int            `json:"priority" gorm:"default:0"`
	IsActive        bool           `json:"isActive" gorm:"default:true"`
	DedupEnabled    bool           `json:"dedupEnabled" gorm:"default:true"`
	DedupWindow     int            `json:"dedupWindow" gorm:"default:60"`
	DropUnmatched   bool           `json:"dropUnmatched" gorm:"default:false"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	DeletedAt       gorm.DeletedAt `json:"-" gorm:"index"`
}

// FilterCondition 筛选条件
type FilterCondition struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

// WhitelistItem 白名单条目
type WhitelistItem struct {
	CIDR        string `json:"cidr"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// AlertPolicy 告警策略（关联筛选策略与机器人）
type AlertPolicy struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	Name             string         `json:"name" gorm:"size:100;not null"`
	Description      string         `json:"description" gorm:"size:500"`
	FilterPolicyID   uint           `json:"filterPolicyId" gorm:"index"`
	RobotID          uint           `json:"robotId" gorm:"index"`
	OutputTemplateID uint           `json:"outputTemplateId" gorm:"index"`
	DeviceID         uint           `json:"deviceId" gorm:"index"`
	DeviceGroupID    uint           `json:"deviceGroupId" gorm:"index"`
	IsActive         bool           `json:"isActive" gorm:"default:true"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}

// AlertRule 告警规则（机器人 × 筛选策略 × 输出模板）
type AlertRule struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	RobotID          uint           `json:"robotId" gorm:"index;not null"`
	FilterPolicyID   uint           `json:"filterPolicyId" gorm:"index;not null"`
	OutputTemplateID uint           `json:"outputTemplateId"`
	OutputFormat     string         `json:"outputFormat" gorm:"size:20;default:'json'"` // Syslog 输出格式
	IsActive         bool           `json:"isActive" gorm:"default:true"`
	CreatedAt        time.Time      `json:"createdAt"`
	UpdatedAt        time.Time      `json:"updatedAt"`
	DeletedAt        gorm.DeletedAt `json:"-" gorm:"index"`
}
