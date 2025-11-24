package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	JWT         JWTConfig
	OAuth2      OAuth2Config
	MFA         MFAConfig
	SMTP        SMTPConfig
	SocialAuth  SocialAuthConfig
	WebAuthn    WebAuthnConfig
	App         AppConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port string
	Mode string // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host      string
	Port      string
	User      string
	Password  string
	Name      string
	Charset   string
	ParseTime bool
	Loc       string
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret            string
	AccessTokenExpire time.Duration
	RefreshTokenExpire time.Duration
}

// OAuth2Config OAuth2配置
type OAuth2Config struct {
	AuthorizationCodeExpire time.Duration
	AccessTokenExpire       time.Duration
	RefreshTokenExpire      time.Duration
}

// MFAConfig MFA配置
type MFAConfig struct {
	Issuer string
}

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	FromName string
}

// SocialAuthConfig 社交登录配置
type SocialAuthConfig struct {
	GitHubClientID     string
	GitHubClientSecret string
}

// WebAuthnConfig WebAuthn配置
type WebAuthnConfig struct {
	RPID          string
	RPOrigins     []string
	RPDisplayName string
}

// AppConfig 应用配置
type AppConfig struct {
	Name        string
	URL         string
	FrontendURL string
}

var Cfg *Config

// Load 加载配置
func Load() {
	// 尝试加载 .env 文件（如果存在）
	_ = godotenv.Load()

	Cfg = &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:      getEnv("DB_HOST", "localhost"),
			Port:      getEnv("DB_PORT", "3306"),
			User:      getEnv("DB_USER", "root"),
			Password:  getEnv("DB_PASSWORD", ""),
			Name:      getEnv("DB_NAME", "astro_pass"),
			Charset:   getEnv("DB_CHARSET", "utf8mb4"),
			ParseTime: getEnvBool("DB_PARSE_TIME", true),
			Loc:       getEnv("DB_LOC", "Local"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", ""),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "your-secret-key-change-in-production-min-32-chars"),
			AccessTokenExpire: getEnvDuration("JWT_ACCESS_TOKEN_EXPIRE", 15*time.Minute),
			RefreshTokenExpire: getEnvDuration("JWT_REFRESH_TOKEN_EXPIRE", 168*time.Hour),
		},
		OAuth2: OAuth2Config{
			AuthorizationCodeExpire: getEnvDuration("OAUTH2_AUTHORIZATION_CODE_EXPIRE", 10*time.Minute),
			AccessTokenExpire:       getEnvDuration("OAUTH2_ACCESS_TOKEN_EXPIRE", 15*time.Minute),
			RefreshTokenExpire:      getEnvDuration("OAUTH2_REFRESH_TOKEN_EXPIRE", 168*time.Hour),
		},
		MFA: MFAConfig{
			Issuer: getEnv("MFA_ISSUER", "Astro-Pass"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnv("SMTP_PORT", "587"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
			FromName: getEnv("SMTP_FROM_NAME", "Astro-Pass"),
		},
		SocialAuth: SocialAuthConfig{
			GitHubClientID:     getEnv("GITHUB_CLIENT_ID", ""),
			GitHubClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
		},
		WebAuthn: WebAuthnConfig{
			RPID:          getEnv("WEBAUTHN_RP_ID", "localhost"),
			RPOrigins:     []string{getEnv("WEBAUTHN_RP_ORIGIN", "http://localhost:3000")},
			RPDisplayName: getEnv("WEBAUTHN_RP_DISPLAY_NAME", "Astro-Pass"),
		},
		App: AppConfig{
			Name:        getEnv("APP_NAME", "星穹通行证"),
			URL:         getEnv("APP_URL", "http://localhost:8080"),
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:3000"),
		},
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvBool 获取布尔类型环境变量
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return defaultValue
}

// getEnvInt 获取整数类型环境变量
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return defaultValue
}

// getEnvDuration 获取时间间隔类型环境变量
func getEnvDuration(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if d, err := time.ParseDuration(value); err == nil {
			return d
		}
	}
	return defaultValue
}
