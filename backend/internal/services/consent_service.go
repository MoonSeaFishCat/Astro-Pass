package services

import (
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"errors"
	"time"
)

type ConsentService struct{}

func NewConsentService() *ConsentService {
	return &ConsentService{}
}

// CheckConsent 检查用户是否已授权
func (s *ConsentService) CheckConsent(userID uint, clientID, scope string) (bool, error) {
	var consent models.UserConsent
	err := database.DB.Where("user_id = ? AND client_id = ?", userID, clientID).First(&consent).Error
	
	if err != nil {
		return false, nil // 未找到授权记录
	}

	// 检查是否过期
	if time.Now().After(consent.ExpiresAt) {
		return false, nil
	}

	// 检查scope是否匹配（简化处理：检查请求的scope是否是已授权scope的子集）
	if !isScopeSubset(scope, consent.Scope) {
		return false, nil
	}

	return true, nil
}

// SaveConsent 保存用户授权
func (s *ConsentService) SaveConsent(userID uint, clientID, scope string) error {
	// 检查是否已存在
	var consent models.UserConsent
	err := database.DB.Where("user_id = ? AND client_id = ?", userID, clientID).First(&consent).Error
	
	if err == nil {
		// 更新现有授权
		consent.Scope = scope
		consent.ExpiresAt = time.Now().Add(365 * 24 * time.Hour) // 1年有效期
		return database.DB.Save(&consent).Error
	}

	// 创建新授权
	consent = models.UserConsent{
		UserID:    userID,
		ClientID:  clientID,
		Scope:     scope,
		ExpiresAt: time.Now().Add(365 * 24 * time.Hour),
	}

	return database.DB.Create(&consent).Error
}

// RevokeConsent 撤销用户授权
func (s *ConsentService) RevokeConsent(userID uint, clientID string) error {
	result := database.DB.Where("user_id = ? AND client_id = ?", userID, clientID).Delete(&models.UserConsent{})
	if result.Error != nil {
		return errors.New("撤销授权失败")
	}
	if result.RowsAffected == 0 {
		return errors.New("授权记录不存在")
	}
	return nil
}

// GetUserConsents 获取用户的所有授权
func (s *ConsentService) GetUserConsents(userID uint) ([]models.UserConsent, error) {
	var consents []models.UserConsent
	err := database.DB.Where("user_id = ?", userID).Find(&consents).Error
	return consents, err
}

// isScopeSubset 检查requestedScope是否是grantedScope的子集
func isScopeSubset(requestedScope, grantedScope string) bool {
	if requestedScope == "" {
		return true
	}
	if grantedScope == "" {
		return false
	}

	requestedScopes := make(map[string]bool)
	for _, s := range splitScopes(requestedScope) {
		requestedScopes[s] = true
	}

	grantedScopes := make(map[string]bool)
	for _, s := range splitScopes(grantedScope) {
		grantedScopes[s] = true
	}

	for scope := range requestedScopes {
		if !grantedScopes[scope] {
			return false
		}
	}

	return true
}
