package services

import (
	"errors"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
)

type SessionService struct{}

func NewSessionService() *SessionService {
	return &SessionService{}
}

// CreateSession 创建会话
func (s *SessionService) CreateSession(userID uint, token, ip, userAgent, device string) (*models.UserSession, error) {
	session := &models.UserSession{
		UserID:       userID,
		Token:        token,
		IP:           ip,
		UserAgent:    userAgent,
		Device:       device,
		LastActivity: time.Now(),
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 7), // 7天
	}

	if err := database.DB.Create(session).Error; err != nil {
		return nil, errors.New("创建会话失败")
	}

	return session, nil
}

// GetUserSessions 获取用户的所有活跃会话
func (s *SessionService) GetUserSessions(userID uint) ([]models.UserSession, error) {
	var sessions []models.UserSession
	if err := database.DB.Where("user_id = ? AND revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Order("last_activity DESC").
		Find(&sessions).Error; err != nil {
		return nil, errors.New("获取会话列表失败")
	}
	return sessions, nil
}

// RevokeSession 撤销会话
func (s *SessionService) RevokeSession(sessionID, userID uint) error {
	var session models.UserSession
	if err := database.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&session).Error; err != nil {
		return errors.New("会话不存在")
	}

	session.Revoked = true
	if err := database.DB.Save(&session).Error; err != nil {
		return errors.New("撤销会话失败")
	}

	return nil
}

// RevokeAllSessions 撤销用户所有会话（除了当前会话）
func (s *SessionService) RevokeAllSessions(userID uint, currentToken string) error {
	if err := database.DB.Model(&models.UserSession{}).
		Where("user_id = ? AND token != ?", userID, currentToken).
		Update("revoked", true).Error; err != nil {
		return errors.New("撤销所有会话失败")
	}
	return nil
}

// UpdateSessionActivity 更新会话活动时间
func (s *SessionService) UpdateSessionActivity(token string) error {
	if err := database.DB.Model(&models.UserSession{}).
		Where("token = ? AND revoked = ?", token, false).
		Update("last_activity", time.Now()).Error; err != nil {
		return errors.New("更新会话活动时间失败")
	}
	return nil
}


