package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/models"
)

type EmailVerificationService struct {
	emailService *EmailService
}

func NewEmailVerificationService() *EmailVerificationService {
	return &EmailVerificationService{
		emailService: NewEmailService(),
	}
}

// SendVerificationEmail 发送验证邮件
func (s *EmailVerificationService) SendVerificationEmail(userID uint, email string) error {
	// 生成验证令牌
	token, err := s.generateToken()
	if err != nil {
		return errors.New("生成验证令牌失败")
	}

	// 创建或更新验证记录
	var verification models.EmailVerification
	if err := database.DB.Where("user_id = ? AND email = ?", userID, email).First(&verification).Error; err != nil {
		verification = models.EmailVerification{
			UserID:    userID,
			Email:     email,
			Token:     token,
			Verified:  false,
			ExpiresAt: time.Now().Add(24 * time.Hour),
		}
		if err := database.DB.Create(&verification).Error; err != nil {
			return errors.New("创建验证记录失败")
		}
	} else {
		verification.Token = token
		verification.ExpiresAt = time.Now().Add(24 * time.Hour)
		verification.Verified = false
		if err := database.DB.Save(&verification).Error; err != nil {
			return errors.New("更新验证记录失败")
		}
	}

	// 发送验证邮件
	verificationURL := config.Cfg.App.FrontendURL + "/verify-email?token=" + token
	return s.emailService.SendVerificationEmail(email, verificationURL)
}

// VerifyEmail 验证邮箱
func (s *EmailVerificationService) VerifyEmail(token string) error {
	var verification models.EmailVerification
	if err := database.DB.Where("token = ?", token).First(&verification).Error; err != nil {
		return errors.New("无效的验证令牌")
	}

	if verification.Verified {
		return errors.New("邮箱已验证")
	}

	if time.Now().After(verification.ExpiresAt) {
		return errors.New("验证令牌已过期")
	}

	// 标记为已验证
	verification.Verified = true
	if err := database.DB.Save(&verification).Error; err != nil {
		return errors.New("更新验证状态失败")
	}

	// 更新用户邮箱验证状态
	var user models.User
	if err := database.DB.First(&user, verification.UserID).Error; err != nil {
		return errors.New("用户不存在")
	}

	if user.Email == verification.Email {
		user.EmailVerified = true
		if err := database.DB.Save(&user).Error; err != nil {
			return errors.New("更新用户验证状态失败")
		}
	}

	return nil
}

func (s *EmailVerificationService) generateToken() (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(tokenBytes), nil
}


