package services

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"

	"astro-pass/internal/config"
)

type EmailService struct{}

func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendEmail 发送邮件
func (s *EmailService) SendEmail(to, subject, body string) error {
	cfg := config.Cfg.SMTP

	// 如果未配置SMTP，跳过发送（开发环境）
	if cfg.Host == "" {
		return nil
	}

	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		cfg.From, to, subject, body))

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	if err := smtp.SendMail(addr, auth, cfg.Username, []string{to}, msg); err != nil {
		return fmt.Errorf("发送邮件失败: %w", err)
	}

	return nil
}

// SendPasswordResetEmail 发送密码重置邮件
func (s *EmailService) SendPasswordResetEmail(to, resetToken string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", config.Cfg.App.FrontendURL, resetToken)

	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>密码重置</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #AEC6E4;">✨ 星穹通行证 - 密码重置</h2>
		<p>您好！</p>
		<p>我们收到了您的密码重置请求。请点击下面的链接来重置您的密码：</p>
		<p style="text-align: center; margin: 30px 0;">
			<a href="{{.ResetURL}}" style="background-color: #AEC6E4; color: white; padding: 12px 24px; text-decoration: none; border-radius: 8px; display: inline-block;">重置密码</a>
		</p>
		<p>或者复制以下链接到浏览器中打开：</p>
		<p style="word-break: break-all; color: #666;">{{.ResetURL}}</p>
		<p style="color: #999; font-size: 12px; margin-top: 30px;">
			此链接将在24小时后过期。如果您没有请求重置密码，请忽略此邮件。
		</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
		<p style="color: #999; font-size: 12px; text-align: center;">
			星穹通行证团队
		</p>
	</div>
</body>
</html>
	`

	t, err := template.New("reset").Parse(tmpl)
	if err != nil {
		return errors.New("解析邮件模板失败")
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]string{"ResetURL": resetURL}); err != nil {
		return errors.New("生成邮件内容失败")
	}

	return s.SendEmail(to, "星穹通行证 - 密码重置", buf.String())
}

// SendWelcomeEmail 发送欢迎邮件
func (s *EmailService) SendWelcomeEmail(to, username string) error {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>欢迎加入星穹通行证</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #AEC6E4;">✨ 欢迎加入星穹学院！</h2>
		<p>您好，{{.Username}}！</p>
		<p>恭喜您成功注册星穹通行证！</p>
		<p>现在您可以：</p>
		<ul>
			<li>使用您的账户登录系统</li>
			<li>设置多因素认证（MFA）增强账户安全</li>
			<li>管理您的个人资料和权限</li>
		</ul>
		<p>如果您有任何问题，请随时联系我们。</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
		<p style="color: #999; font-size: 12px; text-align: center;">
			星穹通行证团队
		</p>
	</div>
</body>
</html>
	`

	t, err := template.New("welcome").Parse(tmpl)
	if err != nil {
		return errors.New("解析邮件模板失败")
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]string{"Username": username}); err != nil {
		return errors.New("生成邮件内容失败")
	}

	return s.SendEmail(to, "欢迎加入星穹通行证", buf.String())
}

// SendVerificationEmail 发送邮箱验证邮件
func (s *EmailService) SendVerificationEmail(to, verificationURL string) error {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>验证您的邮箱</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
	<div style="max-width: 600px; margin: 0 auto; padding: 20px;">
		<h2 style="color: #AEC6E4;">✨ 星穹通行证 - 邮箱验证</h2>
		<p>您好！</p>
		<p>感谢您注册星穹通行证！请点击下面的链接来验证您的邮箱地址：</p>
		<p style="text-align: center; margin: 30px 0;">
			<a href="{{.VerificationURL}}" style="background-color: #AEC6E4; color: white; padding: 12px 24px; text-decoration: none; border-radius: 8px; display: inline-block;">验证邮箱</a>
		</p>
		<p>或者复制以下链接到浏览器中打开：</p>
		<p style="word-break: break-all; color: #666;">{{.VerificationURL}}</p>
		<p style="color: #999; font-size: 12px; margin-top: 30px;">
			此链接将在24小时后过期。如果您没有注册账户，请忽略此邮件。
		</p>
		<hr style="border: none; border-top: 1px solid #eee; margin: 30px 0;">
		<p style="color: #999; font-size: 12px; text-align: center;">
			星穹通行证团队
		</p>
	</div>
</body>
</html>
	`

	t, err := template.New("verification").Parse(tmpl)
	if err != nil {
		return errors.New("解析邮件模板失败")
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]string{"VerificationURL": verificationURL}); err != nil {
		return errors.New("生成邮件内容失败")
	}

	return s.SendEmail(to, "星穹通行证 - 验证您的邮箱", buf.String())
}

