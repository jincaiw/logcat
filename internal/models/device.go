// Package models 定义数据模型，纯结构体，不包含业务逻辑。
// 保持 JSON tag 与原实现一致，确保 API 兼容。
package models

import (
	"time"

	"gorm.io/gorm"
)

// DeviceGroup 设备分组
type DeviceGroup struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null;unique"`
	Description string         `json:"description" gorm:"size:500"`
	Color       string         `json:"color" gorm:"size:20;default:'#409eff'"`
	SortOrder   int            `json:"sortOrder" gorm:"default:0"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

// Device 设备
type Device struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"size:100;not null"`
	IPAddress   string         `json:"ipAddress" gorm:"size:50;not null;uniqueIndex"`
	GroupID     uint           `json:"groupId" gorm:"index"`
	GroupName   string         `json:"groupName" gorm:"size:50;default:'default'"`
	Description string         `json:"description" gorm:"size:500"`
	IsActive    bool           `json:"isActive" gorm:"default:true"`
	CreatedAt   time.Time      `json:"createdAt"`
	UpdatedAt   time.Time      `json:"updatedAt"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}
