package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
	"gorm.io/gorm"
)

type AuthService struct{}

func NewAuthService() *AuthService {
	return &AuthService{}
}

// Register 用户注册
func (s *AuthService) Register(username, email, password, nickname string) (*models.User, error) {
	// 输入验证
	if !utils.ValidateUsername(username) {
		return nil, errors.New("用户名格式不正确（3-50个字符，只能包含字母、数字、下划线）")
	}
	if !utils.ValidateEmail(email) {
		return nil, errors.New("邮箱格式不正确")
	}
	if !utils.ValidatePassword(password) {
		return nil, errors.New("密码长度至少为6位")
	}
	if nickname != "" && !utils.ValidateNickname(nickname) {
		return nil, errors.New("昵称长度应在1-50个字符之间")
	}

	// 清理输入
	username = utils.SanitizeInput(username)
	email = utils.SanitizeInput(email)
	nickname = utils.SanitizeInput(nickname)

	// 检查用户名是否已存在
	var existingUser models.User
	if err := database.DB.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, errors.New("用户名或邮箱已存在")
	}

	// 加密密码
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	// 创建用户
	user := &models.User{
		UUID:         utils.GenerateUUID(),
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Nickname:     nickname,
		Status:       "active",
	}

	if err := database.DB.Create(user).Error; err != nil {
		return nil, errors.New("用户创建失败")
	}

	// 记录审计日志
	s.createAuditLog(user.ID, "register", "user", user.UUID, "用户注册成功", nil)

	// 发送欢迎邮件
	emailService := NewEmailService()
	go emailService.SendWelcomeEmail(user.Email, user.Username)

	return user, nil
}

	// Login 用户登录
func (s *AuthService) Login(username, password, ip, userAgent string) (*models.User, string, string, error) {
	// 检查账户是否被锁定
	lockService := NewAccountLockService()
	locked, unlockTime, err := lockService.IsAccountLocked(username, ip)
	if err == nil && locked {
		remainingTime := unlockTime.Sub(time.Now()).Minutes()
		return nil, "", "", fmt.Errorf("账户已被锁定，请在 %.0f 分钟后重试", remainingTime)
	}

	var user models.User
	if err := database.DB.Preload("Roles").Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			lockService.RecordLoginAttempt(username, ip, false, "用户不存在")
			s.createAuditLog(0, "login", "user", username, "登录失败：用户不存在", map[string]interface{}{"ip": ip})
			return nil, "", "", errors.New("用户名或密码错误")
		}
		return nil, "", "", err
	}

	// 检查用户状态
	if user.Status != "active" {
		lockService.RecordLoginAttempt(username, ip, false, "用户状态异常")
		s.createAuditLog(user.ID, "login", "user", user.UUID, "登录失败：用户状态异常", map[string]interface{}{"ip": ip})
		return nil, "", "", errors.New("账户已被暂停或删除")
	}

	// 验证密码
	if !utils.CheckPassword(password, user.PasswordHash) {
		lockService.RecordLoginAttempt(username, ip, false, "密码错误")
		s.createAuditLog(user.ID, "login", "user", user.UUID, "登录失败：密码错误", map[string]interface{}{"ip": ip})
		return nil, "", "", errors.New("用户名或密码错误")
	}

	// 登录成功，清除登录尝试记录
	lockService.ClearLoginAttempts(username, ip)

	// 生成Token
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", "", errors.New("生成访问令牌失败")
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", errors.New("生成刷新令牌失败")
	}

	// 保存刷新令牌
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7), // 7天
	}
	database.DB.Create(refreshTokenModel)

	// 更新最后登录信息
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = ip
	database.DB.Save(&user)

	// 创建会话
	sessionService := NewSessionService()
	device := s.detectDevice(userAgent)
	sessionService.CreateSession(user.ID, refreshToken, ip, userAgent, device)

	// 记录审计日志
	s.createAuditLog(user.ID, "login", "user", user.UUID, "登录成功", map[string]interface{}{"ip": ip, "user_agent": userAgent})

	return &user, accessToken, refreshToken, nil
}

// detectDevice 检测设备类型
func (s *AuthService) detectDevice(userAgent string) string {
	ua := strings.ToLower(userAgent)
	if strings.Contains(ua, "mobile") || strings.Contains(ua, "android") || strings.Contains(ua, "iphone") {
		return "mobile"
	}
	if strings.Contains(ua, "tablet") || strings.Contains(ua, "ipad") {
		return "tablet"
	}
	return "desktop"
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(refreshTokenString string) (string, string, error) {
	// 解析刷新令牌
	_, err := utils.ParseToken(refreshTokenString)
	if err != nil {
		return "", "", errors.New("无效的刷新令牌")
	}

	// 查找刷新令牌记录
	var refreshToken models.RefreshToken
	if err := database.DB.Where("token = ? AND revoked = ?", refreshTokenString, false).First(&refreshToken).Error; err != nil {
		return "", "", errors.New("刷新令牌不存在或已撤销")
	}

	// 检查是否过期
	if time.Now().After(refreshToken.ExpiresAt) {
		return "", "", errors.New("刷新令牌已过期")
	}

	// 获取用户信息
	var user models.User
	if err := database.DB.First(&user, refreshToken.UserID).Error; err != nil {
		return "", "", errors.New("用户不存在")
	}

	// 生成新的访问令牌和刷新令牌
	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		return "", "", errors.New("生成访问令牌失败")
	}

	newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", errors.New("生成刷新令牌失败")
	}

	// 撤销旧的刷新令牌
	refreshToken.Revoked = true
	database.DB.Save(&refreshToken)

	// 保存新的刷新令牌
	newRefreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	database.DB.Create(newRefreshTokenModel)

	return newAccessToken, newRefreshToken, nil
}

// createAuditLog 创建审计日志
func (s *AuthService) createAuditLog(userID uint, action, resource, resourceID, message string, metadata map[string]interface{}) {
	auditService := NewAuditService()
	var userIDPtr *uint
	if userID > 0 {
		userIDPtr = &userID
	}
	
	// 从metadata中提取IP和UserAgent
	var ip, userAgent string
	if metadata != nil {
		if ipVal, ok := metadata["ip"].(string); ok {
			ip = ipVal
		}
		if uaVal, ok := metadata["user_agent"].(string); ok {
			userAgent = uaVal
		}
	}
	
	_ = auditService.CreateAuditLog(userIDPtr, action, resource, resourceID, message, "success", ip, userAgent, metadata)
}

