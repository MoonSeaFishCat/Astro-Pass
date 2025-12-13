package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
	"github.com/google/uuid"
)

type SLOService struct{}

func NewSLOService() *SLOService {
	return &SLOService{}
}

// CreateSSOSession 创建SSO会话
func (s *SLOService) CreateSSOSession(userID uint, clientID, accessToken, logoutURL string) (*models.SSOSession, error) {
	sessionID := uuid.New().String()
	
	session := &models.SSOSession{
		SessionID:   sessionID,
		UserID:      userID,
		ClientID:    clientID,
		AccessToken: accessToken,
		LogoutURL:   logoutURL,
		Status:      "active",
	}

	if err := database.DB.Create(session).Error; err != nil {
		return nil, fmt.Errorf("创建SSO会话失败: %v", err)
	}

	return session, nil
}

// GetUserActiveSessions 获取用户的活跃会话
func (s *SLOService) GetUserActiveSessions(userID uint) ([]models.SSOSession, error) {
	var sessions []models.SSOSession
	if err := database.DB.Where("user_id = ? AND status = ?", userID, "active").
		Preload("Client").
		Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("获取用户会话失败: %v", err)
	}
	return sessions, nil
}

// InitiateLogout 发起登出请求
func (s *SLOService) InitiateLogout(sessionID, initiatorType string, initiatorID uint) (*models.LogoutRequest, error) {
	// 获取会话信息
	var session models.SSOSession
	if err := database.DB.Where("session_id = ? AND status = ?", sessionID, "active").
		First(&session).Error; err != nil {
		return nil, fmt.Errorf("会话不存在或已失效")
	}

	// 获取用户的所有活跃会话
	sessions, err := s.GetUserActiveSessions(session.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取用户会话失败: %v", err)
	}

	// 创建登出请求
	requestID := uuid.New().String()
	logoutRequest := &models.LogoutRequest{
		RequestID:     requestID,
		SessionID:     sessionID,
		InitiatorType: initiatorType,
		InitiatorID:   initiatorID,
		Status:        "pending",
		TotalClients:  len(sessions),
	}

	if err := database.DB.Create(logoutRequest).Error; err != nil {
		return nil, fmt.Errorf("创建登出请求失败: %v", err)
	}

	// 创建登出通知
	for _, sess := range sessions {
		if sess.LogoutURL != "" {
			notification := &models.LogoutNotification{
				RequestID: requestID,
				ClientID:  sess.ClientID,
				LogoutURL: sess.LogoutURL,
				Status:    "pending",
			}
			database.DB.Create(notification)
		}
	}

	// 异步处理登出通知
	go s.processLogoutNotifications(requestID)

	return logoutRequest, nil
}

// processLogoutNotifications 处理登出通知
func (s *SLOService) processLogoutNotifications(requestID string) {
	// 更新请求状态为处理中
	database.DB.Model(&models.LogoutRequest{}).
		Where("request_id = ?", requestID).
		Update("status", "processing")

	// 获取所有待处理的通知
	var notifications []models.LogoutNotification
	database.DB.Where("request_id = ? AND status = ?", requestID, "pending").
		Find(&notifications)

	completedCount := 0
	failedCount := 0

	for _, notification := range notifications {
		success := s.sendLogoutNotification(&notification)
		if success {
			completedCount++
		} else {
			failedCount++
		}
	}

	// 更新登出请求状态
	status := "completed"
	if failedCount > 0 {
		status = "failed"
	}

	database.DB.Model(&models.LogoutRequest{}).
		Where("request_id = ?", requestID).
		Updates(map[string]interface{}{
			"status":            status,
			"completed_clients": completedCount,
			"failed_clients":    failedCount,
		})

	// 标记所有相关会话为已登出
	var logoutRequest models.LogoutRequest
	if err := database.DB.Where("request_id = ?", requestID).First(&logoutRequest).Error; err == nil {
		var session models.SSOSession
		if err := database.DB.Where("session_id = ?", logoutRequest.SessionID).First(&session).Error; err == nil {
			// 标记用户的所有会话为已登出
			database.DB.Model(&models.SSOSession{}).
				Where("user_id = ? AND status = ?", session.UserID, "active").
				Update("status", "logged_out")
		}
	}
}

// sendLogoutNotification 发送登出通知
func (s *SLOService) sendLogoutNotification(notification *models.LogoutNotification) bool {
	maxAttempts := 3
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// 准备请求数据
		logoutData := map[string]interface{}{
			"logout_request_id": notification.RequestID,
			"client_id":         notification.ClientID,
			"timestamp":         time.Now().Unix(),
		}

		jsonData, _ := json.Marshal(logoutData)
		
		// 发送POST请求
		resp, err := client.Post(notification.LogoutURL, "application/json", bytes.NewBuffer(jsonData))
		
		now := time.Now()
		notification.LastAttemptAt = &now
		notification.AttemptCount = attempt

		if err != nil {
			utils.Error("发送登出通知失败 (尝试 %d/%d): %v", attempt, maxAttempts, err)
			if attempt == maxAttempts {
				notification.Status = "failed"
				notification.ResponseBody = err.Error()
				database.DB.Save(notification)
				return false
			}
			time.Sleep(time.Duration(attempt) * time.Second) // 指数退避
			continue
		}

		notification.ResponseCode = resp.StatusCode
		
		// 读取响应
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		notification.ResponseBody = buf.String()
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			notification.Status = "success"
			database.DB.Save(notification)
			return true
		}

		utils.Warn("登出通知返回错误状态码 %d (尝试 %d/%d)", resp.StatusCode, attempt, maxAttempts)
		if attempt == maxAttempts {
			notification.Status = "failed"
			database.DB.Save(notification)
			return false
		}
		
		time.Sleep(time.Duration(attempt) * time.Second)
	}

	return false
}

// GetLogoutStatus 获取登出状态
func (s *SLOService) GetLogoutStatus(requestID string) (*models.LogoutRequest, []models.LogoutNotification, error) {
	var logoutRequest models.LogoutRequest
	if err := database.DB.Where("request_id = ?", requestID).First(&logoutRequest).Error; err != nil {
		return nil, nil, fmt.Errorf("登出请求不存在")
	}

	var notifications []models.LogoutNotification
	database.DB.Where("request_id = ?", requestID).Find(&notifications)

	return &logoutRequest, notifications, nil
}

// CleanupExpiredSessions 清理过期会话
func (s *SLOService) CleanupExpiredSessions() error {
	// 清理24小时前的已登出会话
	cutoff := time.Now().Add(-24 * time.Hour)
	
	result := database.DB.Where("status = ? AND updated_at < ?", "logged_out", cutoff).
		Delete(&models.SSOSession{})
	
	if result.Error != nil {
		return fmt.Errorf("清理过期会话失败: %v", result.Error)
	}

	utils.Info("清理了 %d 个过期会话", result.RowsAffected)
	return nil
}

// RevokeUserSessions 撤销用户的所有会话
func (s *SLOService) RevokeUserSessions(userID uint, initiatorType string, initiatorID uint) error {
	// 获取用户的活跃会话
	sessions, err := s.GetUserActiveSessions(userID)
	if err != nil {
		return err
	}

	if len(sessions) == 0 {
		return nil // 没有活跃会话
	}

	// 为每个会话发起登出
	for _, session := range sessions {
		_, err := s.InitiateLogout(session.SessionID, initiatorType, initiatorID)
		if err != nil {
			utils.Error("发起会话登出失败: %v", err)
		}
	}

	return nil
}