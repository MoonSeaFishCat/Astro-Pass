package models

import (
	"time"

	"gorm.io/gorm"
)

// UserSession 用户会话模型
type UserSession struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Token     string         `gorm:"uniqueIndex;size:255;not null" json:"token"`
	IP        string         `gorm:"size:45" json:"ip"`
	UserAgent string         `gorm:"size:255" json:"user_agent"`
	Device    string         `gorm:"size:100" json:"device"` // 设备类型：desktop, mobile, tablet
	Location  string         `gorm:"size:100" json:"location"` // 地理位置（可选）
	LastActivity time.Time   `gorm:"not null;index" json:"last_activity"`
	ExpiresAt    time.Time   `gorm:"not null;index" json:"expires_at"`
	Revoked      bool        `gorm:"default:false" json:"revoked"`
	CreatedAt    time.Time   `json:"created_at"`
	UpdatedAt    time.Time   `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// LoginAttempt 登录尝试记录（用于账户锁定）
type LoginAttempt struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"size:100;not null;index" json:"username"`
	IP        string         `gorm:"size:45;not null;index" json:"ip"`
	Success   bool           `gorm:"default:false" json:"success"`
	Message   string         `gorm:"size:255" json:"message"`
	CreatedAt time.Time      `gorm:"index" json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


