package models

import (
	"time"
	"gorm.io/gorm"
)

// SAMLConfig SAML配置模型
type SAMLConfig struct {
	ID                    uint           `json:"id" gorm:"primaryKey"`
	EntityID              string         `json:"entity_id" gorm:"type:varchar(255);uniqueIndex;not null"`     // 实体ID
	Type                  string         `json:"type" gorm:"type:varchar(50);not null"`                      // idp, sp
	Name                  string         `json:"name" gorm:"type:varchar(255);not null"`                      // 配置名称
	Description           string         `json:"description" gorm:"type:text"`                               // 描述
	Status                string         `json:"status" gorm:"type:varchar(50);default:'active'"`           // active, inactive
	
	// IdP配置（作为身份提供者）
	IDPCertificate        string         `json:"idp_certificate" gorm:"type:text"`                          // IdP证书
	IDPPrivateKey         string         `json:"idp_private_key" gorm:"type:text"`                          // IdP私钥
	IDPSSOServiceURL      string         `json:"idp_sso_service_url" gorm:"type:varchar(500)"`                      // SSO服务URL
	IDPSLOServiceURL      string         `json:"idp_slo_service_url" gorm:"type:varchar(500)"`                      // SLO服务URL
	
	// SP配置（作为服务提供者）
	SPEntityID            string         `json:"sp_entity_id" gorm:"type:varchar(255)"`                             // SP实体ID
	SPAssertionConsumerURL string        `json:"sp_assertion_consumer_url" gorm:"type:varchar(500)"`               // 断言消费URL
	SPSingleLogoutURL     string         `json:"sp_single_logout_url" gorm:"type:varchar(500)"`                    // 单点登出URL
	SPCertificate         string         `json:"sp_certificate" gorm:"type:text"`                          // SP证书
	
	// 属性映射
	AttributeMapping      string         `json:"attribute_mapping" gorm:"type:text"`                        // JSON格式的属性映射
	
	// 安全设置
	SignAssertions        bool           `json:"sign_assertions" gorm:"default:true"`      // 签名断言
	EncryptAssertions     bool           `json:"encrypt_assertions" gorm:"default:false"`  // 加密断言
	SignRequests          bool           `json:"sign_requests" gorm:"default:false"`       // 签名请求
	
	CreatedAt             time.Time      `json:"created_at"`
	UpdatedAt             time.Time      `json:"updated_at"`
	DeletedAt             gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// SAMLRequest SAML请求模型
type SAMLRequest struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	RequestID     string         `json:"request_id" gorm:"type:varchar(255);uniqueIndex;not null"`  // SAML请求ID
	Type          string         `json:"type" gorm:"type:varchar(50);not null"`                    // AuthnRequest, LogoutRequest
	EntityID      string         `json:"entity_id" gorm:"type:varchar(255);not null"`               // 发起方实体ID
	UserID        uint           `json:"user_id"`                                 // 关联用户ID（如果已认证）
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	RelayState    string         `json:"relay_state" gorm:"type:varchar(500)"`                             // 中继状态
	RequestData   string         `json:"request_data" gorm:"type:text"`           // 原始请求数据
	Status        string         `json:"status" gorm:"type:varchar(50);default:'pending'"`         // pending, processed, expired
	ExpiresAt     time.Time      `json:"expires_at"`                              // 过期时间
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

// SAMLAssertion SAML断言模型
type SAMLAssertion struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	AssertionID   string         `json:"assertion_id" gorm:"type:varchar(255);uniqueIndex;not null"` // 断言ID
	RequestID     string         `json:"request_id" gorm:"type:varchar(255);not null"`               // 关联的请求ID
	UserID        uint           `json:"user_id" gorm:"not null"`                  // 用户ID
	User          User           `json:"user" gorm:"foreignKey:UserID"`
	EntityID      string         `json:"entity_id" gorm:"type:varchar(255);not null"`                // 目标实体ID
	AssertionData string         `json:"assertion_data" gorm:"type:text"`          // 断言数据
	Status        string         `json:"status" gorm:"type:varchar(50);default:'active'"`           // active, consumed, expired
	ExpiresAt     time.Time      `json:"expires_at"`                               // 过期时间
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}