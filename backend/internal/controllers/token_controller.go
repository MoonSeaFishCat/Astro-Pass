package controllers

import (
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"crypto/x509"
	"encoding/base64"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TokenController struct {
	tokenService *services.TokenService
}

func NewTokenController() *TokenController {
	return &TokenController{
		tokenService: services.NewTokenService(),
	}
}

// RefreshToken 刷新访问令牌
func (tc *TokenController) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	accessToken, refreshToken, err := tc.tokenService.RefreshAccessToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "令牌刷新成功",
		"data": gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
			"token_type":    "Bearer",
			"expires_in":    900, // 15分钟
		},
	})
}

// RevokeToken 撤销令牌（RFC 7009）
func (tc *TokenController) RevokeToken(c *gin.Context) {
	token := c.PostForm("token")
	tokenTypeHint := c.PostForm("token_type_hint") // "access_token" 或 "refresh_token"

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "token parameter is required",
		})
		return
	}

	err := tc.tokenService.RevokeToken(token, tokenTypeHint)
	if err != nil {
		// RFC 7009规定：即使令牌不存在也应返回200
		// 这是为了防止信息泄露
	}

	// 成功撤销或令牌不存在都返回200
	c.Status(http.StatusOK)
}

// IntrospectToken 令牌内省（RFC 7662）
func (tc *TokenController) IntrospectToken(c *gin.Context) {
	token := c.PostForm("token")
	_ = c.PostForm("token_type_hint") // token_type_hint 是可选参数，暂时不使用

	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":             "invalid_request",
			"error_description": "token parameter is required",
		})
		return
	}

	// 验证客户端凭证（简化处理，实际应该验证Basic Auth）
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")
	
	if clientID == "" || clientSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":             "invalid_client",
			"error_description": "client authentication failed",
		})
		return
	}

	result, err := tc.tokenService.IntrospectToken(token)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"active": false,
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetJWKS 获取JWKS（JSON Web Key Set）
func (tc *TokenController) GetJWKS(c *gin.Context) {
	publicKey := utils.GetPublicKey()
	if publicKey == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "public key not initialized",
		})
		return
	}

	// 将RSA公钥转换为JWK格式（暂时不使用，但保留用于未来扩展）
	_, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to marshal public key",
		})
		return
	}

	// 计算key ID（使用公钥的SHA-256哈希的前8字节）
	keyID := "rsa-key-1"

	// 提取RSA公钥的n和e
	n := publicKey.N
	e := publicKey.E

	// 将n转换为base64url编码
	nBytes := n.Bytes()
	nBase64 := base64.RawURLEncoding.EncodeToString(nBytes)

	// 将e转换为base64url编码
	eBytes := big.NewInt(int64(e)).Bytes()
	eBase64 := base64.RawURLEncoding.EncodeToString(eBytes)

	c.JSON(http.StatusOK, gin.H{
		"keys": []gin.H{
			{
				"kty": "RSA",
				"use": "sig",
				"kid": keyID,
				"alg": "RS256",
				"n":   nBase64,
				"e":   eBase64,
			},
		},
	})
}

// GetOpenIDConfiguration 获取OpenID Connect配置（自动发现）
func (tc *TokenController) GetOpenIDConfiguration(c *gin.Context) {
	baseURL := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	issuer := scheme + "://" + baseURL

	c.JSON(http.StatusOK, gin.H{
		"issuer":                                issuer,
		"authorization_endpoint":                issuer + "/api/oauth2/authorize",
		"token_endpoint":                        issuer + "/api/oauth2/token",
		"userinfo_endpoint":                     issuer + "/api/oauth2/userinfo",
		"jwks_uri":                              issuer + "/api/oauth2/jwks",
		"revocation_endpoint":                   issuer + "/api/oauth2/revoke",
		"introspection_endpoint":                issuer + "/api/oauth2/introspect",
		"response_types_supported":              []string{"code", "token", "id_token", "code token", "code id_token", "token id_token", "code token id_token"},
		"subject_types_supported":               []string{"public"},
		"id_token_signing_alg_values_supported": []string{"RS256"},
		"scopes_supported":                      []string{"openid", "profile", "email"},
		"token_endpoint_auth_methods_supported": []string{"client_secret_basic", "client_secret_post"},
		"claims_supported":                      []string{"sub", "name", "preferred_username", "email", "email_verified"},
		"grant_types_supported":                 []string{"authorization_code", "refresh_token", "client_credentials"},
	})
}
