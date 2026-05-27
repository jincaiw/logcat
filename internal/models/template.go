package models

import "time"

// FieldMappingDoc represents a field mapping documentation entry
type FieldMappingDoc struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	DeviceType    string    `gorm:"size:64;not null;index" json:"deviceType"`
	StandardField string    `gorm:"size:128;not null" json:"standardField"`
	OriginalField string    `gorm:"size:128;not null" json:"originalField"`
	Description   string    `gorm:"size:512" json:"description"`
	FieldType     string    `gorm:"size:64" json:"fieldType"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

// TableName for FieldMappingDoc
func (FieldMappingDoc) TableName() string {
	return "field_mapping_docs"
}

// ParseTemplate represents a log parsing template
type ParseTemplate struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `gorm:"size:128;not null" json:"name"`
	DeviceType     string    `gorm:"size:64" json:"deviceType"`
	ParseType      string    `gorm:"size:32;not null" json:"parseType"`
	HeaderRegex    string    `gorm:"size:512" json:"headerRegex"`
	Delimiter      string    `gorm:"size:16" json:"delimiter"`
	FieldMapping   string    `gorm:"type:text" json:"fieldMapping"`
	ValueTransform string    `gorm:"type:text" json:"valueTransform"`
	SampleLog      string    `gorm:"type:text" json:"sampleLog"`
	SubTemplates   string    `gorm:"type:text" json:"subTemplates"`
	Enabled        bool      `gorm:"default:true" json:"enabled"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

// TableName for ParseTemplate
func (ParseTemplate) TableName() string {
	return "parse_templates"
}

// OutputTemplate represents an output formatting template
type OutputTemplate struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:128;not null" json:"name"`
	ChannelType string    `gorm:"size:32;not null" json:"channelType"`
	Content     string    `gorm:"type:text" json:"content"`
	Fields      string    `gorm:"type:text" json:"fields"`
	DeviceType  string    `gorm:"size:64" json:"deviceType"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName for OutputTemplate
func (OutputTemplate) TableName() string {
	return "output_templates"
}