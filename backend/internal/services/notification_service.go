package services

import (
	"encoding/json"
	"errors"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
)

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

// CreateNotification 创建通知
func (s *NotificationService) CreateNotification(userID *uint, notificationType, title, message string, metadata map[string]interface{}) error {
	metadataJSON := ""
	if metadata != nil {
		metadataBytes, _ := json.Marshal(metadata)
		metadataJSON = string(metadataBytes)
	}

	notification := models.Notification{
		UserID:   userID,
		Type:     notificationType,
		Title:    title,
		Message:  message,
		Read:     false,
		Metadata: metadataJSON,
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		return errors.New("创建通知失败")
	}

	return nil
}

// GetUserNotifications 获取用户通知
func (s *NotificationService) GetUserNotifications(userID uint, unreadOnly bool, limit int) ([]models.Notification, error) {
	var notifications []models.Notification
	query := database.DB.Where("user_id = ?", userID)

	if unreadOnly {
		query = query.Where("read = ?", false)
	}

	query = query.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&notifications).Error; err != nil {
		return nil, errors.New("获取通知失败")
	}

	return notifications, nil
}

// MarkAsRead 标记为已读
func (s *NotificationService) MarkAsRead(notificationID uint, userID uint) error {
	var notification models.Notification
	if err := database.DB.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		return errors.New("通知不存在")
	}

	now := time.Now()
	notification.Read = true
	notification.ReadAt = &now

	if err := database.DB.Save(&notification).Error; err != nil {
		return errors.New("更新通知状态失败")
	}

	return nil
}

// MarkAllAsRead 标记所有通知为已读
func (s *NotificationService) MarkAllAsRead(userID uint) error {
	now := time.Now()
	if err := database.DB.Model(&models.Notification{}).
		Where("user_id = ? AND read = ?", userID, false).
		Updates(map[string]interface{}{
			"read":    true,
			"read_at": now,
		}).Error; err != nil {
		return errors.New("更新通知状态失败")
	}

	return nil
}

// DeleteNotification 删除通知
func (s *NotificationService) DeleteNotification(notificationID uint, userID uint) error {
	if err := database.DB.Where("id = ? AND user_id = ?", notificationID, userID).Delete(&models.Notification{}).Error; err != nil {
		return errors.New("删除通知失败")
	}

	return nil
}

// NotifySecurityEvent 发送安全事件通知
func (s *NotificationService) NotifySecurityEvent(userID uint, eventType, message string) {
	metadata := map[string]interface{}{
		"event_type": eventType,
		"timestamp":  time.Now().Unix(),
	}
	_ = s.CreateNotification(&userID, "security", "安全提醒", message, metadata)
}

// NotifyActivityEvent 发送活动事件通知
func (s *NotificationService) NotifyActivityEvent(userID uint, activityType, message string) {
	metadata := map[string]interface{}{
		"activity_type": activityType,
		"timestamp":      time.Now().Unix(),
	}
	_ = s.CreateNotification(&userID, "activity", "账户活动", message, metadata)
}


