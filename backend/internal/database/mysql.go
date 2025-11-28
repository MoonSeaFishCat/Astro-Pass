package database

import (
	"fmt"
	"time"

	"astro-pass/internal/config"
	"astro-pass/internal/models"
	"github.com/fatih/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Init() error {
	cfg := config.Cfg.Database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)

	// 使用 Silent 模式，不输出 SQL 日志
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		return fmt.Errorf("连接数据库失败: %w", err)
	}

	// 显示彩色进度信息
	color.New(color.FgCyan, color.Bold).Print("🔗 正在连接数据库...")
	time.Sleep(200 * time.Millisecond)
	color.New(color.FgGreen).Println(" ✓ 连接成功")

	// 自动迁移
	color.New(color.FgCyan, color.Bold).Print("📊 正在初始化数据库表...")
	if err := AutoMigrate(); err != nil {
		color.New(color.FgRed).Println(" ✗ 失败")
		return fmt.Errorf("数据库迁移失败: %w", err)
	}
	color.New(color.FgGreen).Println(" ✓ 完成")

	return nil
}

func AutoMigrate() error {
	// 分步迁移，先创建基础表，再创建有外键依赖的表
	// 第一组：基础表（无外键依赖）
	baseModels := []struct {
		model interface{}
		name  string
	}{
		{&models.User{}, "用户表"},
		{&models.Role{}, "角色表"},
		{&models.Permission{}, "权限表"},
		{&models.OAuth2Client{}, "OAuth2客户端表"},
	}

	// 第二组：有外键依赖的表
	dependentModels := []struct {
		model interface{}
		name  string
	}{
		{&models.RefreshToken{}, "刷新令牌表"},
		{&models.AuditLog{}, "审计日志表"},
		{&models.AuthorizationCode{}, "授权码表"},
		{&models.AccessToken{}, "访问令牌表"},
		{&models.UserSession{}, "用户会话表"},
		{&models.LoginAttempt{}, "登录尝试表"},
		{&models.WebAuthnCredential{}, "WebAuthn凭证表"},
		{&models.SocialAuth{}, "社交登录表"},
		{&models.EmailVerification{}, "邮箱验证表"},
		{&models.PasswordPolicy{}, "密码策略表"},
		{&models.PasswordHistory{}, "密码历史表"},
		{&models.Notification{}, "通知表"},
		{&models.BackupRecord{}, "备份记录表"},
		{&models.SystemConfig{}, "系统配置表"},
		{&models.UserConsent{}, "用户授权同意表"},
	}

	// 先迁移基础表
	color.New(color.FgYellow).Print("    ├─ 迁移基础表...")
	for i, item := range baseModels {
		if err := DB.AutoMigrate(item.model); err != nil {
			color.New(color.FgRed).Printf("\n    └─ ✗ %s 迁移失败\n", item.name)
			return fmt.Errorf("迁移基础表失败: %w", err)
		}
		if i < len(baseModels)-1 {
			color.New(color.FgGreen).Printf(" ✓ %s\n    ├─ ", item.name)
		} else {
			color.New(color.FgGreen).Printf(" ✓ %s\n", item.name)
		}
		time.Sleep(50 * time.Millisecond)
	}

	// 再迁移有外键依赖的表
	color.New(color.FgYellow).Print("    └─ 迁移依赖表...")
	for i, item := range dependentModels {
		if err := DB.AutoMigrate(item.model); err != nil {
			color.New(color.FgRed).Printf("\n       ✗ %s 迁移失败\n", item.name)
			return fmt.Errorf("迁移依赖表失败: %w", err)
		}
		if i < len(dependentModels)-1 {
			color.New(color.FgGreen).Printf(" ✓ %s\n    ├─ ", item.name)
		} else {
			color.New(color.FgGreen).Printf(" ✓ %s\n", item.name)
		}
		time.Sleep(50 * time.Millisecond)
	}

	return nil
}
