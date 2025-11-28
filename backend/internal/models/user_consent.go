package models

import (
	"time"

	"gorm.io/gorm"
)

// UserConsent 用户授权同意记录
type UserConsent struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	ClientID  string         `gorm:"not null;index" json:"client_id"`
	Scope     string         `gorm:"type:varchar(500)" json:"scope"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"-"`
}

// TableName 指定表名
func (UserConsent) TableName() string {
	return "user_consents"
}
