package models

import "time"

// DeviceGroup groups devices logically
type DeviceGroup struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:128;not null" json:"name"`
	Description string    `gorm:"size:512" json:"description"`
	Color       string    `gorm:"size:20" json:"color"`
	SortOrder   int       `gorm:"default:0" json:"sortOrder"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName for DeviceGroup
func (DeviceGroup) TableName() string {
	return "device_groups"
}

// Device represents a log-sending device
type Device struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:128;not null" json:"name"`
	IPAddress        string    `gorm:"uniqueIndex;size:45;not null" json:"ipAddress"`
	GroupID          *uint     `json:"groupId"`
	TemplateID       *uint     `json:"templateId"`
	ParseTemplateID  *uint     `json:"parseTemplateId"`
	DeviceType       string    `gorm:"size:64" json:"deviceType"`
	Description      string    `gorm:"size:512" json:"description"`
	Enabled          bool      `gorm:"default:true" json:"enabled"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Associations
	Group         *DeviceGroup    `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	Template      *DeviceTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	ParseTemplate *ParseTemplate  `gorm:"foreignKey:ParseTemplateID" json:"parseTemplate,omitempty"`
}

// TableName for Device
func (Device) TableName() string {
	return "devices"
}

// DeviceTemplate represents a reusable device configuration template
type DeviceTemplate struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	Name                string    `gorm:"size:128;not null" json:"name"`
	DeviceType          string    `gorm:"size:64" json:"deviceType"`
	ParseTemplateID     *uint     `json:"parseTemplateId"`
	FieldMappingDocID   *uint     `json:"fieldMappingDocId"`
	RecommendedPolicy   string    `gorm:"type:text" json:"recommendedPolicy"`
	Enabled             bool      `gorm:"default:true" json:"enabled"`
	CreatedAt           time.Time `json:"createdAt"`
	UpdatedAt           time.Time `json:"updatedAt"`
}

// TableName for DeviceTemplate
func (DeviceTemplate) TableName() string {
	return "device_templates"
}