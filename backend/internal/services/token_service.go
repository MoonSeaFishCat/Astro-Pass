package services

import (
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
	"errors"
	"time"
)

type TokenService struct{}

func NewTokenService() *TokenService {
	return &TokenService{}
}

// RefreshAccessToken 刷新访问令牌
func (s *TokenService) RefreshAccessToken(refreshTokenString string) (string, string, error) {
	// 查找刷新令牌
	var refreshToken models.RefreshToken
	if err := database.DB.Where("token = ? AND revoked = ?", refreshTokenString, false).First(&refreshToken).Error; err != nil {
		return "", "", errors.New("无效的刷新令牌")
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

	// 生成新的访问令牌
	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		return "", "", errors.New("生成访问令牌失败")
	}

	// 生成新的刷新令牌
	newRefreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return "", "", errors.New("生成刷新令牌失败")
	}

	// 撤销旧的刷新令牌（令牌轮换）
	refreshToken.Revoked = true
	database.DB.Save(&refreshToken)

	// 保存新的刷新令牌
	newRefreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     newRefreshToken,
		ClientID:  refreshToken.ClientID,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}
	database.DB.Create(newRefreshTokenModel)

	return newAccessToken, newRefreshToken, nil
}

// RevokeToken 撤销令牌（RFC 7009）
func (s *TokenService) RevokeToken(token string, tokenTypeHint string) error {
	if tokenTypeHint == "refresh_token" || tokenTypeHint == "" {
		// 尝试撤销刷新令牌
		var refreshToken models.RefreshToken
		if err := database.DB.Where("token = ?", token).First(&refreshToken).Error; err == nil {
			refreshToken.Revoked = true
			database.DB.Save(&refreshToken)
			
			// 同时撤销相关的会话
			database.DB.Model(&models.UserSession{}).
				Where("refresh_token = ?", token).
				Update("revoked", true)
			
			return nil
		}
	}

	if tokenTypeHint == "access_token" || tokenTypeHint == "" {
		// 尝试撤销访问令牌
		var accessToken models.AccessToken
		if err := database.DB.Where("token = ?", token).First(&accessToken).Error; err == nil {
			accessToken.Revoked = true
			database.DB.Save(&accessToken)
			return nil
		}
	}

	return errors.New("令牌不存在")
}

// IntrospectToken 令牌内省（RFC 7662）
func (s *TokenService) IntrospectToken(token string) (map[string]interface{}, error) {
	// 尝试解析为访问令牌
	claims, err := utils.ParseToken(token)
	if err == nil {
		// 检查令牌是否被撤销
		var accessToken models.AccessToken
		if err := database.DB.Where("token = ? AND revoked = ?", token, false).First(&accessToken).Error; err != nil {
			return map[string]interface{}{
				"active": false,
			}, nil
		}

		// 检查是否过期
		if time.Now().After(accessToken.ExpiresAt) {
			return map[string]interface{}{
				"active": false,
			}, nil
		}

		return map[string]interface{}{
			"active":    true,
			"scope":     accessToken.Scope,
			"client_id": accessToken.ClientID,
			"username":  claims.Username,
			"token_type": "Bearer",
			"exp":       accessToken.ExpiresAt.Unix(),
			"iat":       accessToken.CreatedAt.Unix(),
			"sub":       claims.Username,
		}, nil
	}

	// 尝试作为刷新令牌
	var refreshToken models.RefreshToken
	if err := database.DB.Where("token = ? AND revoked = ?", token, false).First(&refreshToken).Error; err != nil {
		return map[string]interface{}{
			"active": false,
		}, nil
	}

	// 检查是否过期
	if time.Now().After(refreshToken.ExpiresAt) {
		return map[string]interface{}{
			"active": false,
		}, nil
	}

	// 获取用户信息
	var user models.User
	if err := database.DB.First(&user, refreshToken.UserID).Error; err != nil {
		return map[string]interface{}{
			"active": false,
		}, nil
	}

	return map[string]interface{}{
		"active":     true,
		"client_id":  refreshToken.ClientID,
		"username":   user.Username,
		"token_type": "refresh_token",
		"exp":        refreshToken.ExpiresAt.Unix(),
		"iat":        refreshToken.CreatedAt.Unix(),
		"sub":        user.Username,
	}, nil
}
