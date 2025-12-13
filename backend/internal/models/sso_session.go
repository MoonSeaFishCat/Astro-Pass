package models

import (
	"time"
	"gorm.io/gorm"
)

// SSOSession SSO会话模型
type SSOSession struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	SessionID   string         `json:"session_id" gorm:"type:varchar(255);uniqueIndex;not null"` // 全局会话ID
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	ClientID    string         `json:"client_id" gorm:"type:varchar(255);not null"`              // OAuth2客户端ID
	Client      OAuth2Client   `json:"client" gorm:"foreignKey:ClientID;references:ClientID"`
	AccessToken string         `json:"access_token" gorm:"type:text;not null"`           // 访问令牌
	LogoutURL   string         `json:"logout_url" gorm:"type:varchar(500)"`                             // 客户端登出URL
	Status      string         `json:"status" gorm:"type:varchar(50);default:'active'"`         // active, logged_out
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// LogoutRequest 登出请求模型
type LogoutRequest struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	RequestID        string         `json:"request_id" gorm:"type:varchar(255);uniqueIndex;not null"` // 登出请求ID
	SessionID        string         `json:"session_id" gorm:"type:varchar(255);not null"`             // 关联的会话ID
	InitiatorType    string         `json:"initiator_type" gorm:"type:varchar(50);not null"`         // user, admin, system
	InitiatorID      uint           `json:"initiator_id"`                           // 发起者ID
	Status           string         `json:"status" gorm:"type:varchar(50);default:'pending'"`        // pending, processing, completed, failed
	TotalClients     int            `json:"total_clients" gorm:"default:0"`         // 需要通知的客户端总数
	CompletedClients int            `json:"completed_clients" gorm:"default:0"`     // 已完成通知的客户端数
	FailedClients    int            `json:"failed_clients" gorm:"default:0"`        // 通知失败的客户端数
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// LogoutNotification 登出通知模型
type LogoutNotification struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	RequestID     string         `json:"request_id" gorm:"type:varchar(255);not null"`     // 关联的登出请求ID
	ClientID      string         `json:"client_id" gorm:"type:varchar(255);not null"`      // 客户端ID
	LogoutURL     string         `json:"logout_url" gorm:"type:varchar(500);not null"`     // 登出URL
	Status        string         `json:"status" gorm:"type:varchar(50);default:'pending'"` // pending, success, failed, timeout
	ResponseCode  int            `json:"response_code"`                  // HTTP响应码
	ResponseBody  string         `json:"response_body" gorm:"type:text"`                  // 响应内容
	AttemptCount  int            `json:"attempt_count" gorm:"default:0"` // 尝试次数
	LastAttemptAt *time.Time     `json:"last_attempt_at"`                // 最后尝试时间
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}