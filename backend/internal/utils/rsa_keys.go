package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

// InitRSAKeys 初始化RSA密钥对
func InitRSAKeys() error {
	// 尝试从文件加载密钥
	if err := loadKeysFromFile(); err == nil {
		return nil
	}

	// 如果文件不存在，生成新密钥
	return generateAndSaveKeys()
}

// generateAndSaveKeys 生成并保存RSA密钥对
func generateAndSaveKeys() error {
	// 生成2048位RSA密钥对
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	privateKey = key
	publicKey = &key.PublicKey

	// 保存私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	if err := os.WriteFile("private_key.pem", privateKeyPEM, 0600); err != nil {
		return err
	}

	// 保存公钥
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	if err := os.WriteFile("public_key.pem", publicKeyPEM, 0644); err != nil {
		return err
	}

	return nil
}

// loadKeysFromFile 从文件加载RSA密钥对
func loadKeysFromFile() error {
	// 加载私钥
	privateKeyPEM, err := os.ReadFile("private_key.pem")
	if err != nil {
		return err
	}

	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return errors.New("failed to decode private key PEM")
	}

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return err
	}

	// 加载公钥
	publicKeyPEM, err := os.ReadFile("public_key.pem")
	if err != nil {
		return err
	}

	block, _ = pem.Decode(publicKeyPEM)
	if block == nil {
		return errors.New("failed to decode public key PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	var ok bool
	publicKey, ok = pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}

	return nil
}

// GetPrivateKey 获取私钥
func GetPrivateKey() *rsa.PrivateKey {
	return privateKey
}

// GetPublicKey 获取公钥
func GetPublicKey() *rsa.PublicKey {
	return publicKey
}
