package models

import (
	"time"

	"gorm.io/gorm"
)

// OAuth2Client OAuth2客户端模型
type OAuth2Client struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	UserID            uint           `gorm:"not null;index" json:"user_id"`
	ClientID          string         `gorm:"uniqueIndex;size:100;not null" json:"client_id"`
	ClientSecret      string         `gorm:"size:255;not null" json:"-"`
	ClientName        string         `gorm:"size:100;not null" json:"client_name"`
	ClientURI         string         `gorm:"size:255" json:"client_uri"`
	LogoURI           string         `gorm:"size:255" json:"logo_uri"`
	RedirectURIs      string         `gorm:"type:text;not null" json:"-"` // JSON格式存储多个重定向URI
	GrantTypes        string         `gorm:"type:text;not null" json:"-"` // JSON格式存储授权类型
	ResponseTypes     string         `gorm:"type:text;not null" json:"-"` // JSON格式存储响应类型
	Scope             string         `gorm:"size:255" json:"scope"`
	Status            string         `gorm:"size:20;default:active" json:"status"` // active, suspended, revoked
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	User              User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AuthorizationCodes []AuthorizationCode `gorm:"foreignKey:OAuth2ClientID" json:"-"`
	AccessTokens      []AccessToken  `gorm:"foreignKey:OAuth2ClientID" json:"-"`
}

// AuthorizationCode 授权码模型
type AuthorizationCode struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Code              string         `gorm:"uniqueIndex;size:255;not null" json:"code"`
	OAuth2ClientID   uint           `gorm:"not null;index" json:"oauth2_client_id"` // 外键引用 OAuth2Client.ID
	ClientID          string         `gorm:"not null;index" json:"client_id"` // OAuth2 标准中的 client_id（字符串）
	UserID            uint           `gorm:"not null;index" json:"user_id"`
	RedirectURI       string         `gorm:"size:255;not null" json:"redirect_uri"`
	Scope             string         `gorm:"size:255" json:"scope"`
	CodeChallenge     string         `gorm:"size:255" json:"-"` // PKCE支持
	CodeChallengeMethod string       `gorm:"size:20" json:"-"` // S256, plain
	ExpiresAt         time.Time      `gorm:"not null;index" json:"expires_at"`
	Used              bool           `gorm:"default:false" json:"used"`
	CreatedAt         time.Time      `json:"created_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Client            OAuth2Client   `gorm:"foreignKey:OAuth2ClientID" json:"client,omitempty"`
	User              User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// AccessToken 访问令牌模型
type AccessToken struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Token             string         `gorm:"uniqueIndex;size:500;not null" json:"token"`
	OAuth2ClientID   uint           `gorm:"not null;index" json:"oauth2_client_id"` // 外键引用 OAuth2Client.ID
	ClientID          string         `gorm:"not null;index" json:"client_id"` // OAuth2 标准中的 client_id（字符串）
	UserID            *uint          `gorm:"index" json:"user_id"` // 可为空，支持客户端凭证模式
	Scope             string         `gorm:"size:255" json:"scope"`
	ExpiresAt         time.Time      `gorm:"not null;index" json:"expires_at"`
	Revoked           bool           `gorm:"default:false" json:"revoked"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联关系
	Client            OAuth2Client   `gorm:"foreignKey:OAuth2ClientID" json:"client,omitempty"`
	User              *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

