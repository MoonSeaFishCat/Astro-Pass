package services

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
)

type OAuth2Service struct{}

func NewOAuth2Service() *OAuth2Service {
	return &OAuth2Service{}
}

// GenerateAuthorizationCode 生成授权码
func (s *OAuth2Service) GenerateAuthorizationCode(clientID string, userID uint, redirectURI, scope, codeChallenge, codeChallengeMethod string) (string, error) {
	// 生成随机授权码
	codeBytes := make([]byte, 32)
	if _, err := rand.Read(codeBytes); err != nil {
		return "", errors.New("生成授权码失败")
	}
	code := base64.URLEncoding.EncodeToString(codeBytes)

	// 查找客户端以获取 ID
	var client models.OAuth2Client
	if err := database.DB.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		return "", errors.New("无效的客户端")
	}

	// 保存授权码
	authCode := &models.AuthorizationCode{
		Code:              code,
		OAuth2ClientID:    client.ID, // 外键
		ClientID:          clientID,  // OAuth2 标准中的 client_id
		UserID:            userID,
		RedirectURI:       redirectURI,
		Scope:             scope,
		CodeChallenge:     codeChallenge,
		CodeChallengeMethod: codeChallengeMethod,
		ExpiresAt:         time.Now().Add(config.Cfg.OAuth2.AuthorizationCodeExpire),
	}

	if err := database.DB.Create(authCode).Error; err != nil {
		return "", errors.New("保存授权码失败")
	}

	return code, nil
}

// ClientCredentialsGrant 客户端凭证模式
func (s *OAuth2Service) ClientCredentialsGrant(clientID, clientSecret, scope string) (*models.AccessToken, error) {
	// 验证客户端
	var client models.OAuth2Client
	if err := database.DB.Where("client_id = ? AND status = ?", clientID, "active").First(&client).Error; err != nil {
		return nil, errors.New("无效的客户端")
	}

	// 验证客户端密钥
	if client.ClientSecret != clientSecret {
		return nil, errors.New("客户端密钥错误")
	}

	// 检查是否支持客户端凭证模式
	var grantTypes []string
	if err := json.Unmarshal([]byte(client.GrantTypes), &grantTypes); err != nil {
		return nil, errors.New("客户端配置错误")
	}

	supportsClientCredentials := false
	for _, gt := range grantTypes {
		if gt == "client_credentials" {
			supportsClientCredentials = true
			break
		}
	}

	if !supportsClientCredentials {
		return nil, errors.New("客户端不支持client_credentials授权类型")
	}

	// 生成访问令牌（客户端凭证模式不需要用户ID）
	accessTokenString, err := utils.GenerateAccessToken(0, "", "")
	if err != nil {
		return nil, errors.New("生成访问令牌失败")
	}

	// 保存访问令牌（UserID为nil）
	accessToken := &models.AccessToken{
		Token:          accessTokenString,
		OAuth2ClientID: client.ID,
		ClientID:       clientID,
		UserID:         nil, // 客户端凭证模式没有用户
		Scope:          scope,
		ExpiresAt:      time.Now().Add(config.Cfg.OAuth2.AccessTokenExpire),
	}
	database.DB.Create(accessToken)

	return accessToken, nil
}

// TokenResponse 令牌响应结构
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

// ExchangeAuthorizationCode 交换授权码获取令牌
func (s *OAuth2Service) ExchangeAuthorizationCode(code, clientID, clientSecret, redirectURI, codeVerifier string) (*TokenResponse, error) {
	// 查找授权码
	var authCode models.AuthorizationCode
	if err := database.DB.Where("code = ? AND used = ?", code, false).First(&authCode).Error; err != nil {
		return nil, errors.New("无效的授权码")
	}

	// 检查是否过期
	if time.Now().After(authCode.ExpiresAt) {
		return nil, errors.New("授权码已过期")
	}

	// 验证客户端
	var client models.OAuth2Client
	if err := database.DB.Where("client_id = ? AND status = ?", clientID, "active").First(&client).Error; err != nil {
		return nil, errors.New("无效的客户端")
	}

	// 验证客户端密钥
	if client.ClientSecret != clientSecret {
		return nil, errors.New("客户端密钥错误")
	}

	// 验证重定向URI
	var redirectURIs []string
	if err := json.Unmarshal([]byte(client.RedirectURIs), &redirectURIs); err != nil {
		return nil, errors.New("客户端配置错误")
	}

	redirectURIMatch := false
	for _, uri := range redirectURIs {
		if uri == redirectURI {
			redirectURIMatch = true
			break
		}
	}
	if !redirectURIMatch {
		return nil, errors.New("重定向URI不匹配")
	}

	// 验证PKCE（如果使用）
	if authCode.CodeChallenge != "" {
		if codeVerifier == "" {
			return nil, errors.New("需要提供code_verifier")
		}
		// 这里应该验证code_verifier，简化处理
	}

	// 标记授权码为已使用
	authCode.Used = true
	database.DB.Save(&authCode)

	// 获取用户信息
	var user models.User
	if err := database.DB.First(&user, authCode.UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	// 生成访问令牌
	accessTokenString, err := utils.GenerateAccessToken(authCode.UserID, user.Username, user.Email)
	if err != nil {
		return nil, errors.New("生成访问令牌失败")
	}

	// 生成刷新令牌
	refreshTokenString, err := utils.GenerateRefreshToken(authCode.UserID)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败")
	}

	// 生成ID Token（如果scope包含openid）
	var idTokenString string
	if containsScope(authCode.Scope, "openid") {
		issuer := config.Cfg.Server.AppURL
		idTokenString, err = utils.GenerateIDToken(
			user.ID,
			user.Username,
			user.Email,
			user.Nickname,
			user.EmailVerified,
			"", // nonce应该从授权请求中获取
			issuer,
			clientID,
		)
		if err != nil {
			return nil, errors.New("生成ID Token失败")
		}
	}

	// 保存访问令牌
	accessToken := &models.AccessToken{
		Token:          accessTokenString,
		OAuth2ClientID: client.ID, // 外键
		ClientID:       clientID,  // OAuth2 标准中的 client_id
		UserID:         &authCode.UserID,
		Scope:          authCode.Scope,
		ExpiresAt:      time.Now().Add(config.Cfg.OAuth2.AccessTokenExpire),
	}
	database.DB.Create(accessToken)

	// 保存刷新令牌
	refreshToken := &models.RefreshToken{
		UserID:    authCode.UserID,
		Token:     refreshTokenString,
		ClientID:  clientID,
		ExpiresAt: time.Now().Add(config.Cfg.OAuth2.RefreshTokenExpire),
	}
	database.DB.Create(refreshToken)

	return &TokenResponse{
		AccessToken:  accessTokenString,
		TokenType:    "Bearer",
		ExpiresIn:    int(config.Cfg.OAuth2.AccessTokenExpire.Seconds()),
		RefreshToken: refreshTokenString,
		IDToken:      idTokenString,
		Scope:        authCode.Scope,
	}, nil
}

// containsScope 检查scope字符串是否包含指定的scope
func containsScope(scopeString, targetScope string) bool {
	if scopeString == "" {
		return false
	}
	scopes := make(map[string]bool)
	for _, s := range splitScopes(scopeString) {
		scopes[s] = true
	}
	return scopes[targetScope]
}

// splitScopes 分割scope字符串
func splitScopes(scopeString string) []string {
	if scopeString == "" {
		return []string{}
	}
	var scopes []string
	for _, s := range splitBySpace(scopeString) {
		if s != "" {
			scopes = append(scopes, s)
		}
	}
	return scopes
}

// splitBySpace 按空格分割字符串
func splitBySpace(s string) []string {
	var result []string
	current := ""
	for _, c := range s {
		if c == ' ' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

// GetUserInfo 获取用户信息（OIDC）
func (s *OAuth2Service) GetUserInfo(accessToken string) (map[string]interface{}, error) {
	claims, err := utils.ParseToken(accessToken)
	if err != nil {
		return nil, errors.New("无效的访问令牌")
	}

	// 验证令牌是否被撤销
	var token models.AccessToken
	if err := database.DB.Where("token = ? AND revoked = ?", accessToken, false).First(&token).Error; err != nil {
		return nil, errors.New("访问令牌已撤销")
	}

	// 获取用户信息
	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return map[string]interface{}{
		"sub":                user.UUID,
		"name":               user.Nickname,
		"preferred_username": user.Username,
		"email":              user.Email,
		"email_verified":     user.EmailVerified,
	}, nil
}

// CreateClient 创建OAuth2客户端
func (s *OAuth2Service) CreateClient(userID uint, clientName, clientURI, logoURI string, redirectURIs []string) (*models.OAuth2Client, error) {
	// 生成客户端ID和密钥
	clientIDBytes := make([]byte, 16)
	if _, err := rand.Read(clientIDBytes); err != nil {
		return nil, errors.New("生成客户端ID失败")
	}
	clientID := base64.URLEncoding.EncodeToString(clientIDBytes)

	clientSecretBytes := make([]byte, 32)
	if _, err := rand.Read(clientSecretBytes); err != nil {
		return nil, errors.New("生成客户端密钥失败")
	}
	clientSecret := base64.URLEncoding.EncodeToString(clientSecretBytes)

	// 序列化重定向URI
	redirectURIsJSON, _ := json.Marshal(redirectURIs)

	// 默认授权类型和响应类型
	grantTypes := []string{"authorization_code", "refresh_token"}
	responseTypes := []string{"code"}
	grantTypesJSON, _ := json.Marshal(grantTypes)
	responseTypesJSON, _ := json.Marshal(responseTypes)

	client := &models.OAuth2Client{
		UserID:        userID,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		ClientName:    clientName,
		ClientURI:     clientURI,
		LogoURI:       logoURI,
		RedirectURIs:  string(redirectURIsJSON),
		GrantTypes:    string(grantTypesJSON),
		ResponseTypes: string(responseTypesJSON),
		Status:        "active",
	}

	if err := database.DB.Create(client).Error; err != nil {
		return nil, errors.New("创建客户端失败")
	}

	return client, nil
}

// GetUserClients 获取用户的所有OAuth2客户端
func (s *OAuth2Service) GetUserClients(userID uint) ([]models.OAuth2Client, error) {
	var clients []models.OAuth2Client
	if err := database.DB.Where("user_id = ? AND status != ?", userID, "revoked").
		Order("created_at DESC").
		Find(&clients).Error; err != nil {
		return nil, errors.New("获取客户端列表失败")
	}
	return clients, nil
}

// RevokeClient 撤销客户端
func (s *OAuth2Service) RevokeClient(clientID string, userID uint) error {
	var client models.OAuth2Client
	if err := database.DB.Where("client_id = ? AND user_id = ?", clientID, userID).First(&client).Error; err != nil {
		return errors.New("客户端不存在")
	}

	client.Status = "revoked"
	if err := database.DB.Save(&client).Error; err != nil {
		return errors.New("撤销客户端失败")
	}

	return nil
}

