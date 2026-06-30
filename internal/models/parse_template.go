package models

import (
	"time"

	"gorm.io/gorm"
)

// ParseTemplate 解析模板
type ParseTemplate struct {
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `json:"name" gorm:"size:100;not null"`
	Description    string         `json:"description" gorm:"size:500"`
	ParseType      string         `json:"parseType" gorm:"size:20;default:'json'"` // json, regex, kv, syslog_json, smart_delimiter, delimiter, keyvalue
	HeaderRegex    string         `json:"headerRegex" gorm:"type:text"`
	FieldMapping   string         `json:"fieldMapping" gorm:"type:text"`   // JSON 格式字段映射
	ValueTransform string         `json:"valueTransform" gorm:"type:text"` // 值转换规则
	SampleLog      string         `json:"sampleLog" gorm:"type:text"`
	DeviceType     string         `json:"deviceType" gorm:"size:50"`
	Delimiter      string         `json:"delimiter" gorm:"size:50"`
	TypeField      int            `json:"typeField" gorm:"default:0"`
	SubTemplates   string         `json:"subTemplates" gorm:"type:text"` // 子模板配置（JSON）
	IsActive       bool           `json:"isActive" gorm:"default:true"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"-" gorm:"index"`
}

// OutputTemplate 输出模板（告警消息格式）
type OutputTemplate struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	Platform    string         `json:"platform" gorm:"size:20;default:'feishu'"`
	Description string         `json:"description" gorm:"size:500"`
	Content     string         `json:"content" gorm:"type:text;not null"` // 模板内容，支持 {{变量}} 替换
	Fields      string         `json:"fields" gorm:"type:text"`           // 可用字段列表
	DeviceType  string         `json:"deviceType" gorm:"size:50"`
	IsActive    bool           `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// FieldMappingDoc 字段映射文档
type FieldMappingDoc struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"size:100;not null"`
	DeviceType    string         `json:"deviceType" gorm:"size:50;not null;index"`
	Description   string         `json:"description" gorm:"type:text"`
	FieldMappings string         `json:"fieldMappings" gorm:"type:text"`
	IsActive      bool           `json:"isActive" gorm:"default:true"`
	CreatedAt     time.Time      `json:"createdAt"`
	UpdatedAt     time.Time      `json:"updatedAt"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`
}
