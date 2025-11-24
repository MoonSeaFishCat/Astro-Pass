package models

import (
	"time"

	"gorm.io/gorm"
)

// WebAuthnCredential WebAuthn凭证
type WebAuthnCredential struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	UserID          uint      `gorm:"not null;index" json:"user_id"`
	CredentialID    string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"credential_id"`
	PublicKey       string    `gorm:"type:text;not null" json:"public_key"`
	Counter         uint32    `gorm:"default:0" json:"counter"`
	AAGUID          string    `gorm:"type:varchar(36)" json:"aaguid"`
	DeviceName      string    `gorm:"type:varchar(100)" json:"device_name"`
	DeviceType      string    `gorm:"type:varchar(50)" json:"device_type"` // platform, cross-platform
	LastUsedAt      *time.Time `json:"last_used_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// SocialAuth 社交媒体认证
type SocialAuth struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"not null;index" json:"user_id"`
	Provider      string    `gorm:"type:varchar(50);not null;index" json:"provider"` // github, wechat, etc.
	ProviderID    string    `gorm:"type:varchar(255);not null" json:"provider_id"`
	ProviderEmail string    `gorm:"type:varchar(255)" json:"provider_email"`
	AccessToken   string    `gorm:"type:text" json:"-"` // 加密存储
	RefreshToken  string    `gorm:"type:text" json:"-"` // 加密存储
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// EmailVerification 邮箱验证
type EmailVerification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	Email     string    `gorm:"type:varchar(255);not null;index" json:"email"`
	Token     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	Verified  bool      `gorm:"default:false" json:"verified"`
	ExpiresAt time.Time `gorm:"not null;index" json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// PasswordPolicy 密码策略
type PasswordPolicy struct {
	ID                    uint      `gorm:"primaryKey" json:"id"`
	MinLength             int       `gorm:"default:8" json:"min_length"`
	RequireUppercase      bool      `gorm:"default:true" json:"require_uppercase"`
	RequireLowercase      bool      `gorm:"default:true" json:"require_lowercase"`
	RequireNumber         bool      `gorm:"default:true" json:"require_number"`
	RequireSpecialChar    bool      `gorm:"default:true" json:"require_special_char"`
	MaxAge                int       `gorm:"default:90" json:"max_age"` // 天数
	HistoryCount          int       `gorm:"default:5" json:"history_count"` // 不能重复最近N个密码
	LockoutThreshold      int       `gorm:"default:5" json:"lockout_threshold"`
	LockoutDuration       int       `gorm:"default:30" json:"lockout_duration"` // 分钟
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// PasswordHistory 密码历史
type PasswordHistory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"not null;index" json:"user_id"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// Notification 通知
type Notification struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"` // nil表示系统通知
	Type      string    `gorm:"type:varchar(50);not null;index" json:"type"` // security, activity, system
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Read      bool      `gorm:"default:false;index" json:"read"`
	ReadAt    *time.Time `json:"read_at,omitempty"`
	Metadata  string    `gorm:"type:text" json:"metadata,omitempty"` // JSON格式
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}


