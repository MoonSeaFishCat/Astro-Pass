package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// IDTokenClaims ID Token的声明
type IDTokenClaims struct {
	jwt.RegisteredClaims
	Name              string `json:"name,omitempty"`
	PreferredUsername string `json:"preferred_username,omitempty"`
	Email             string `json:"email,omitempty"`
	EmailVerified     bool   `json:"email_verified,omitempty"`
	Nonce             string `json:"nonce,omitempty"`
}

// GenerateIDToken 生成ID Token（使用RS256签名）
func GenerateIDToken(userID uint, username, email, nickname string, emailVerified bool, nonce string, issuer string, audience string) (string, error) {
	// 生成唯一的JTI
	jtiBytes := make([]byte, 16)
	if _, err := rand.Read(jtiBytes); err != nil {
		return "", err
	}
	jti := base64.URLEncoding.EncodeToString(jtiBytes)

	now := time.Now()
	claims := IDTokenClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   username,
			Audience:  jwt.ClaimStrings{audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ID:        jti,
		},
		Name:              nickname,
		PreferredUsername: username,
		Email:             email,
		EmailVerified:     emailVerified,
		Nonce:             nonce,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	
	// 使用RSA私钥签名
	privateKey := GetPrivateKey()
	if privateKey == nil {
		return "", jwt.ErrInvalidKey
	}

	return token.SignedString(privateKey)
}

// ParseIDToken 解析ID Token
func ParseIDToken(tokenString string) (*IDTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &IDTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return GetPublicKey(), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*IDTokenClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
