package services

import (
	"errors"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
)

type AccountLockService struct{}

func NewAccountLockService() *AccountLockService {
	return &AccountLockService{}
}

const (
	MaxLoginAttempts = 5        // 最大登录尝试次数
	LockoutDuration  = 30 * time.Minute // 锁定持续时间
)

// RecordLoginAttempt 记录登录尝试
func (s *AccountLockService) RecordLoginAttempt(username, ip string, success bool, message string) error {
	attempt := &models.LoginAttempt{
		Username: username,
		IP:       ip,
		Success:  success,
		Message:  message,
	}

	if err := database.DB.Create(attempt).Error; err != nil {
		return errors.New("记录登录尝试失败")
	}

	return nil
}

// IsAccountLocked 检查账户是否被锁定
func (s *AccountLockService) IsAccountLocked(username, ip string) (bool, time.Time, error) {
	// 检查最近30分钟内的失败尝试次数
	cutoffTime := time.Now().Add(-LockoutDuration)
	
	var failedAttempts int64
	if err := database.DB.Model(&models.LoginAttempt{}).
		Where("username = ? AND ip = ? AND success = ? AND created_at > ?", username, ip, false, cutoffTime).
		Count(&failedAttempts).Error; err != nil {
		return false, time.Time{}, err
	}

	if failedAttempts >= MaxLoginAttempts {
		// 获取最后一次失败尝试的时间
		var lastAttempt models.LoginAttempt
		if err := database.DB.Where("username = ? AND ip = ? AND success = ?", username, ip, false).
			Order("created_at DESC").
			First(&lastAttempt).Error; err == nil {
			unlockTime := lastAttempt.CreatedAt.Add(LockoutDuration)
			if time.Now().Before(unlockTime) {
				return true, unlockTime, nil
			}
		}
	}

	return false, time.Time{}, nil
}

// ClearLoginAttempts 清除登录尝试记录（登录成功后）
func (s *AccountLockService) ClearLoginAttempts(username, ip string) error {
	// 删除30分钟前的记录
	cutoffTime := time.Now().Add(-time.Hour)
	if err := database.DB.Where("username = ? AND ip = ? AND created_at < ?", username, ip, cutoffTime).
		Delete(&models.LoginAttempt{}).Error; err != nil {
		return errors.New("清除登录尝试记录失败")
	}
	return nil
}


