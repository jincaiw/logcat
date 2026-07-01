package models

import (
	"time"

	"gorm.io/gorm"
)

// User 单用户认证模型。
// 系统仅维护一个用户记录（单用户模式），存储登录凭证与个人信息。
type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Username     string         `json:"username" gorm:"size:50;not null;uniqueIndex"`
	PasswordHash string         `json:"-" gorm:"size:200;not null"`
	Nickname     string         `json:"nickname" gorm:"size:50"`
	Email        string         `json:"email" gorm:"size:100"`
	Avatar       string         `json:"avatar" gorm:"size:500"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserView 用户信息视图（不含敏感字段），用于 API 响应。
type UserView struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// ToView 将 User 转换为 UserView。
func (u *User) ToView() UserView {
	return UserView{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
