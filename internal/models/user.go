package models

import (
	"time"
)

// UserStatus represents the status of a user account
type UserStatus string

const (
	UserStatusEnabled  UserStatus = "enabled"
	UserStatusDisabled UserStatus = "disabled"
	UserStatusLocked   UserStatus = "locked"
)

// User represents a user in the system
type User struct {
	ID                 uint       `gorm:"primaryKey" json:"id"`
	Username           string     `gorm:"uniqueIndex;size:64;not null" json:"username"`
	PasswordHash       string     `gorm:"size:255;not null" json:"-"`
	DisplayName        string     `gorm:"size:128" json:"displayName"`
	Email              string     `gorm:"size:255" json:"email"`
	Status             UserStatus `gorm:"size:20;not null;default:enabled" json:"status"`
	FailedLoginCount   int        `gorm:"default:0" json:"failedLoginCount"`
	LockedUntil        *time.Time `json:"lockedUntil"`
	LastLoginAt        *time.Time `json:"lastLoginAt"`
	MustChangePassword bool       `gorm:"default:false" json:"mustChangePassword"`
	CreatedAt          time.Time  `json:"createdAt"`
	UpdatedAt          time.Time  `json:"updatedAt"`

	// Associations
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName specifies the table name for User
func (User) TableName() string {
	return "users"
}

// Role represents a role with permissions
type Role struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:64;not null" json:"name"`
	Code        string    `gorm:"uniqueIndex;size:64;not null" json:"code"`
	Description string    `gorm:"size:255" json:"description"`
	BuiltIn     bool      `gorm:"default:false" json:"builtIn"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`

	// Associations
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// TableName specifies the table name for Role
func (Role) TableName() string {
	return "roles"
}

// Permission represents a permission entry
// PermissionType can be "menu", "button", or "api"
// PermissionType string `json:"type"`
// Action string `json:"action"` - e.g. "create", "read", "update", "delete"
type Permission struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:128;not null" json:"name"`
	Code      string    `gorm:"uniqueIndex;size:128;not null" json:"code"`
	Type      string    `gorm:"size:20;not null" json:"type"`
	Resource  string    `gorm:"size:128;not null" json:"resource"`
	Action    string    `gorm:"size:64;not null" json:"action"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName specifies the table name for Permission
func (Permission) TableName() string {
	return "permissions"
}

// UserRole is the join table for users and roles
type UserRole struct {
	UserID uint `gorm:"primaryKey" json:"userId"`
	RoleID uint `gorm:"primaryKey" json:"roleId"`
}

// TableName specifies the table name for UserRole
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission is the join table for roles and permissions
type RolePermission struct {
	RoleID       uint `gorm:"primaryKey" json:"roleId"`
	PermissionID uint `gorm:"primaryKey" json:"permissionId"`
}

// TableName specifies the table name for RolePermission
func (RolePermission) TableName() string {
	return "role_permissions"
}