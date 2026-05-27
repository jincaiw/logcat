package models

import "time"

// DesensitizeRule represents a data desensitization rule
type DesensitizeRule struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	FieldName   string    `gorm:"size:128;not null" json:"fieldName"`
	RuleType    string    `gorm:"size:32;not null" json:"ruleType"`
	RuleConfig  string    `gorm:"type:text" json:"ruleConfig"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// TableName for DesensitizeRule
func (DesensitizeRule) TableName() string {
	return "desensitize_rules"
}