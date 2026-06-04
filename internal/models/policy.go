package models

import "time"

// FilterPolicy represents a log filtering policy
type FilterPolicy struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Name             string    `gorm:"size:128;not null" json:"name"`
	DeviceID         *uint     `json:"deviceId"`
	DeviceGroupID    *uint     `json:"deviceGroupId"`
	ParseTemplateID  *uint     `json:"parseTemplateId"`
	Conditions       string    `gorm:"type:text" json:"conditions"`
	ConditionLogic   string    `gorm:"size:10;default:AND" json:"conditionLogic"`
	WhitelistEnabled bool      `gorm:"default:false" json:"whitelistEnabled"`
	WhitelistField   string    `gorm:"size:128" json:"whitelistField"`
	WhitelistValues  string    `gorm:"type:text" json:"whitelistValues"`
	Action           string    `gorm:"size:20;default:keep" json:"action"`
	Priority         int       `gorm:"default:0" json:"priority"`
	DedupEnabled     bool      `gorm:"default:false" json:"dedupEnabled"`
	DedupWindow      int       `gorm:"default:300" json:"dedupWindow"`
	Enabled          bool      `gorm:"default:true" json:"enabled"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`

	// Associations
	Device        *Device        `gorm:"foreignKey:DeviceID" json:"device,omitempty"`
	DeviceGroup   *DeviceGroup   `gorm:"foreignKey:DeviceGroupID" json:"deviceGroup,omitempty"`
	ParseTemplate *ParseTemplate `gorm:"foreignKey:ParseTemplateID" json:"parseTemplate,omitempty"`
}

// TableName for FilterPolicy
func (FilterPolicy) TableName() string {
	return "filter_policies"
}