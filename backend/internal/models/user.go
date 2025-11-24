package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UUID              string         `gorm:"uniqueIndex;size:36;not null" json:"uuid"`
	Username          string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email             string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash      string         `gorm:"size:255;not null" json:"-"`
	Nickname          string         `gorm:"size:50" json:"nickname"`
	Avatar            string         `gorm:"size:255" json:"avatar"`
	EmailVerified     bool           `gorm:"default:false" json:"email_verified"`
	MFAEnabled        bool           `gorm:"default:false" json:"mfa_enabled"`
	MFASecret         string         `gorm:"size:255" json:"-"`
	MFARecoveryCodes  string         `gorm:"type:text" json:"-"` // JSON格式存储恢复码
	Status            string         `gorm:"size:20;default:active" json:"status"` // active, suspended, deleted
	LastLoginAt       *time.Time     `json:"last_login_at"`
	LastLoginIP       string         `gorm:"size:45" json:"last_login_ip"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt  `gorm:"index" json:"-"`

	// 关联关系
	Roles             []Role         `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	OAuth2Clients     []OAuth2Client `gorm:"foreignKey:UserID" json:"oauth2_clients,omitempty"`
	RefreshTokens     []RefreshToken `gorm:"foreignKey:UserID" json:"-"`
	AuditLogs         []AuditLog     `gorm:"foreignKey:UserID" json:"-"`
}

// Role 角色模型
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Users       []User         `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission   `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission 权限模型
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:100;not null" json:"name"`
	DisplayName string         `gorm:"size:100" json:"display_name"`
	Resource    string         `gorm:"size:100;not null" json:"resource"`
	Action      string         `gorm:"size:50;not null" json:"action"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Roles       []Role         `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// RefreshToken 刷新令牌模型
type RefreshToken struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Token     string         `gorm:"uniqueIndex;size:255;not null" json:"token"`
	ClientID  string         `gorm:"size:100" json:"client_id"`
	ExpiresAt time.Time      `gorm:"not null;index" json:"expires_at"`
	Revoked   bool           `gorm:"default:false" json:"revoked"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User      User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AuditLog 审计日志模型
type AuditLog struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      *uint          `gorm:"index" json:"user_id"`
	Action      string         `gorm:"size:50;not null;index" json:"action"` // login, logout, register, update_profile, etc.
	Resource    string         `gorm:"size:100" json:"resource"`
	ResourceID  string         `gorm:"size:100" json:"resource_id"`
	IP          string         `gorm:"size:45" json:"ip"`
	UserAgent   string         `gorm:"size:255" json:"user_agent"`
	Status      string         `gorm:"size:20;default:success" json:"status"` // success, failed
	Message     string         `gorm:"type:text" json:"message"`
	Metadata    string         `gorm:"type:text" json:"metadata"` // JSON格式存储额外信息
	CreatedAt   time.Time      `gorm:"index" json:"created_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User       *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}


