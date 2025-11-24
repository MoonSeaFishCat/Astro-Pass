package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
	"net/http"
	"github.com/go-webauthn/webauthn/webauthn"
)

type WebAuthnService struct {
	webauthn *webauthn.WebAuthn
}

func NewWebAuthnService() (*WebAuthnService, error) {
	cfg := config.Cfg.WebAuthn
	
	webAuthnConfig := &webauthn.Config{
		RPDisplayName: cfg.RPDisplayName,
		RPID:          cfg.RPID,
		RPOrigins:     cfg.RPOrigins, // RPOrigins 已经是字符串数组
		// 使用默认的挑战超时时间（60秒）
		Timeouts: webauthn.TimeoutsConfig{
			Login: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    time.Minute,
				TimeoutUVD: time.Minute,
			},
			Registration: webauthn.TimeoutConfig{
				Enforce:    true,
				Timeout:    time.Minute,
				TimeoutUVD: time.Minute,
			},
		},
	}

	w, err := webauthn.New(webAuthnConfig)
	if err != nil {
		return nil, fmt.Errorf("初始化WebAuthn失败: %w", err)
	}

	return &WebAuthnService{
		webauthn: w,
	}, nil
}

// BeginRegistration 开始注册流程
func (s *WebAuthnService) BeginRegistration(userID uint) (*webauthn.SessionData, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取用户已有的凭证
	var credentials []models.WebAuthnCredential
	database.DB.Where("user_id = ?", userID).Find(&credentials)

	// 转换为webauthn.Credential格式
	webauthnCredentials := make([]webauthn.Credential, len(credentials))
	for i, cred := range credentials {
		credentialID, _ := base64.RawURLEncoding.DecodeString(cred.CredentialID)
		publicKey, _ := base64.RawURLEncoding.DecodeString(cred.PublicKey)
		
		webauthnCredentials[i] = webauthn.Credential{
			ID:              credentialID,
			PublicKey:       publicKey,
			AttestationType: "",
			Authenticator: webauthn.Authenticator{
				AAGUID:    []byte(cred.AAGUID),
				SignCount: cred.Counter,
			},
		}
	}

	// 创建WebAuthn用户
	waUser := &WebAuthnUser{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.Nickname,
		Credentials: webauthnCredentials,
	}

	// 开始注册 - BeginRegistration 返回 (options, sessionData, error)
	options, sessionData, err := s.webauthn.BeginRegistration(waUser)
	if err != nil {
		return nil, fmt.Errorf("开始注册失败: %w", err)
	}
	_ = options // 暂时不使用options，只返回sessionData

	return sessionData, nil
}

// FinishRegistration 完成注册流程
func (s *WebAuthnService) FinishRegistration(userID uint, sessionData *webauthn.SessionData, r *http.Request) (*models.WebAuthnCredential, error) {
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取用户已有的凭证
	var existingCredentials []models.WebAuthnCredential
	database.DB.Where("user_id = ?", userID).Find(&existingCredentials)

	webauthnCredentials := make([]webauthn.Credential, len(existingCredentials))
	for i, cred := range existingCredentials {
		credentialID, _ := base64.RawURLEncoding.DecodeString(cred.CredentialID)
		publicKey, _ := base64.RawURLEncoding.DecodeString(cred.PublicKey)
		
		webauthnCredentials[i] = webauthn.Credential{
			ID:              credentialID,
			PublicKey:       publicKey,
			AttestationType: "",
			Authenticator: webauthn.Authenticator{
				AAGUID:    []byte(cred.AAGUID),
				SignCount: cred.Counter,
			},
		}
	}

	waUser := &WebAuthnUser{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.Nickname,
		Credentials: webauthnCredentials,
	}

	// 完成注册 - 使用 http.Request
	credential, err := s.webauthn.FinishRegistration(waUser, *sessionData, r)
	if err != nil {
		return nil, fmt.Errorf("完成注册失败: %w", err)
	}

	// 保存凭证到数据库
	credentialID := base64.RawURLEncoding.EncodeToString(credential.ID)
	publicKey := base64.RawURLEncoding.EncodeToString(credential.PublicKey)
	aaguid := string(credential.Authenticator.AAGUID)

	// 检测设备类型和名称
	deviceType := "cross-platform"
	deviceName := "安全密钥"
	if len(aaguid) > 0 {
		// 可以根据AAGUID判断设备类型
		deviceType = "platform"
	}

	webauthnCredential := &models.WebAuthnCredential{
		UserID:       userID,
		CredentialID: credentialID,
		PublicKey:    publicKey,
		Counter:      credential.Authenticator.SignCount,
		AAGUID:       aaguid,
		DeviceType:   deviceType,
		DeviceName:   deviceName,
	}

	if err := database.DB.Create(webauthnCredential).Error; err != nil {
		return nil, errors.New("保存凭证失败")
	}

	// 记录审计日志
	auditService := NewAuditService()
	userIDPtr := &userID
	auditService.CreateAuditLog(userIDPtr, "webauthn_register", "webauthn_credential", 
		fmt.Sprintf("%d", webauthnCredential.ID), "WebAuthn凭证注册成功", "success", "", "", nil)

	return webauthnCredential, nil
}

// BeginLogin 开始登录流程
func (s *WebAuthnService) BeginLogin(username string) (*webauthn.SessionData, []models.WebAuthnCredential, error) {
	var user models.User
	if err := database.DB.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return nil, nil, errors.New("用户不存在")
	}

	// 获取用户的所有凭证
	var credentials []models.WebAuthnCredential
	if err := database.DB.Where("user_id = ?", user.ID).Find(&credentials).Error; err != nil {
		return nil, nil, errors.New("获取凭证失败")
	}

	if len(credentials) == 0 {
		return nil, nil, errors.New("用户未注册WebAuthn凭证")
	}

	// 转换为webauthn.Credential格式
	webauthnCredentials := make([]webauthn.Credential, len(credentials))
	for i, cred := range credentials {
		credentialID, _ := base64.RawURLEncoding.DecodeString(cred.CredentialID)
		publicKey, _ := base64.RawURLEncoding.DecodeString(cred.PublicKey)
		
		webauthnCredentials[i] = webauthn.Credential{
			ID:              credentialID,
			PublicKey:       publicKey,
			AttestationType: "",
			Authenticator: webauthn.Authenticator{
				AAGUID:    []byte(cred.AAGUID),
				SignCount: cred.Counter,
			},
		}
	}

	waUser := &WebAuthnUser{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.Nickname,
		Credentials: webauthnCredentials,
	}

	// 开始登录
	options, sessionData, err := s.webauthn.BeginLogin(waUser)
	if err != nil {
		return nil, nil, fmt.Errorf("开始登录失败: %w", err)
	}
	_ = options // 暂时不使用options，只返回sessionData

	return sessionData, credentials, nil
}

// FinishLogin 完成登录流程
func (s *WebAuthnService) FinishLogin(username string, sessionData *webauthn.SessionData, r *http.Request) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("username = ? OR email = ?", username, username).First(&user).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 获取用户的所有凭证
	var credentials []models.WebAuthnCredential
	if err := database.DB.Where("user_id = ?", user.ID).Find(&credentials).Error; err != nil {
		return nil, errors.New("获取凭证失败")
	}

	// 转换为webauthn.Credential格式
	webauthnCredentials := make([]webauthn.Credential, len(credentials))
	credentialMap := make(map[string]*models.WebAuthnCredential)
	
	for i, cred := range credentials {
		credentialID, _ := base64.RawURLEncoding.DecodeString(cred.CredentialID)
		publicKey, _ := base64.RawURLEncoding.DecodeString(cred.PublicKey)
		
		webauthnCredentials[i] = webauthn.Credential{
			ID:              credentialID,
			PublicKey:       publicKey,
			AttestationType: "",
			Authenticator: webauthn.Authenticator{
				AAGUID:    []byte(cred.AAGUID),
				SignCount: cred.Counter,
			},
		}
		credentialMap[cred.CredentialID] = &credentials[i]
	}

	waUser := &WebAuthnUser{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.Nickname,
		Credentials: webauthnCredentials,
	}

	// 完成登录 - 使用 http.Request
	credential, err := s.webauthn.FinishLogin(waUser, *sessionData, r)
	if err != nil {
		return nil, fmt.Errorf("完成登录失败: %w", err)
	}

	// 更新凭证计数器
	credentialID := base64.RawURLEncoding.EncodeToString(credential.ID)
	if cred, ok := credentialMap[credentialID]; ok {
		cred.Counter = credential.Authenticator.SignCount
		now := time.Now()
		cred.LastUsedAt = &now
		database.DB.Save(cred)
	}

	// 记录审计日志
	auditService := NewAuditService()
	userIDPtr := &user.ID
	auditService.CreateAuditLog(userIDPtr, "webauthn_login", "user", user.UUID, "WebAuthn登录成功", "success", "", "", nil)

	return &user, nil
}

// GetUserCredentials 获取用户的WebAuthn凭证列表
func (s *WebAuthnService) GetUserCredentials(userID uint) ([]models.WebAuthnCredential, error) {
	var credentials []models.WebAuthnCredential
	if err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&credentials).Error; err != nil {
		return nil, err
	}
	return credentials, nil
}

// DeleteCredential 删除WebAuthn凭证
func (s *WebAuthnService) DeleteCredential(userID uint, credentialID uint) error {
	var credential models.WebAuthnCredential
	if err := database.DB.Where("id = ? AND user_id = ?", credentialID, userID).First(&credential).Error; err != nil {
		return errors.New("凭证不存在")
	}

	if err := database.DB.Delete(&credential).Error; err != nil {
		return errors.New("删除凭证失败")
	}

	// 记录审计日志
	auditService := NewAuditService()
	userIDPtr := &userID
	auditService.CreateAuditLog(userIDPtr, "webauthn_delete", "webauthn_credential", 
		fmt.Sprintf("%d", credentialID), "删除WebAuthn凭证", "success", "", "", nil)

	return nil
}

// WebAuthnUser 实现webauthn.User接口
type WebAuthnUser struct {
	ID          uint
	Username    string
	DisplayName string
	Credentials []webauthn.Credential
}

func (u *WebAuthnUser) WebAuthnID() []byte {
	return []byte(fmt.Sprintf("%d", u.ID))
}

func (u *WebAuthnUser) WebAuthnName() string {
	return u.Username
}

func (u *WebAuthnUser) WebAuthnDisplayName() string {
	if u.DisplayName != "" {
		return u.DisplayName
	}
	return u.Username
}

func (u *WebAuthnUser) WebAuthnCredentials() []webauthn.Credential {
	return u.Credentials
}

func (u *WebAuthnUser) WebAuthnIcon() string {
	return ""
}

// GenerateSessionToken 生成会话令牌（用于存储sessionData）
func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// StoreSessionData 存储会话数据（可以使用Redis或数据库）
func StoreSessionData(token string, sessionData *webauthn.SessionData) error {
	// 这里简化实现，实际应该使用Redis或数据库
	// 为了演示，我们使用内存存储（生产环境应使用Redis）
	sessionDataJSON, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}
	
	// 存储到数据库的临时表（实际应该使用Redis）
	// 这里简化处理，使用utils存储
	return utils.StoreSessionData(token, sessionDataJSON)
}

// GetSessionData 获取会话数据
func GetSessionData(token string) (*webauthn.SessionData, error) {
	sessionDataJSON, err := utils.GetSessionData(token)
	if err != nil {
		return nil, err
	}

	var sessionData webauthn.SessionData
	if err := json.Unmarshal(sessionDataJSON, &sessionData); err != nil {
		return nil, err
	}

	return &sessionData, nil
}

// DeleteSessionData 删除会话数据
func DeleteSessionData(token string) error {
	return utils.DeleteSessionData(token)
}

