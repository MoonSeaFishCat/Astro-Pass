package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"astro-pass/internal/config"
)

// EncryptToken 加密令牌
func EncryptToken(plaintext string) string {
	if plaintext == "" {
		return ""
	}

	key := []byte(config.Cfg.JWT.Secret)
	if len(key) < 32 {
		// 如果密钥太短，使用填充
		paddedKey := make([]byte, 32)
		copy(paddedKey, key)
		key = paddedKey
	} else {
		key = key[:32]
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext // 如果加密失败，返回原文（不推荐，但为了兼容性）
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return plaintext
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.URLEncoding.EncodeToString(ciphertext)
}

// DecryptToken 解密令牌
func DecryptToken(ciphertext string) (string, error) {
	if ciphertext == "" {
		return "", errors.New("密文为空")
	}

	key := []byte(config.Cfg.JWT.Secret)
	if len(key) < 32 {
		paddedKey := make([]byte, 32)
		copy(paddedKey, key)
		key = paddedKey
	} else {
		key = key[:32]
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	data, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(data) < gcm.NonceSize() {
		return "", errors.New("密文太短")
	}

	nonce, ciphertextBytes := data[:gcm.NonceSize()], data[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}


