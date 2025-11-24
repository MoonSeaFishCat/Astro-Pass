package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"astro-pass/internal/config"
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/utils"
)

type SocialAuthService struct{}

func NewSocialAuthService() *SocialAuthService {
	return &SocialAuthService{}
}

// GitHubUserInfo GitHub用户信息
type GitHubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

// GetGitHubAuthURL 获取GitHub授权URL
func (s *SocialAuthService) GetGitHubAuthURL(state string) string {
	clientID := config.Cfg.SocialAuth.GitHubClientID
	redirectURI := config.Cfg.App.URL + "/api/auth/social/github/callback"
	return fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&state=%s&scope=user:email", clientID, url.QueryEscape(redirectURI), state)
}

// HandleGitHubCallback 处理GitHub回调
func (s *SocialAuthService) HandleGitHubCallback(code string) (*GitHubUserInfo, error) {
	// 交换访问令牌
	accessToken, err := s.exchangeGitHubToken(code)
	if err != nil {
		return nil, err
	}

	// 获取用户信息
	userInfo, err := s.getGitHubUserInfo(accessToken)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

// ExchangeGitHubToken 交换GitHub访问令牌（公开方法）
func (s *SocialAuthService) ExchangeGitHubToken(code string) (string, error) {
	return s.exchangeGitHubToken(code)
}

// GetGitHubUserInfo 获取GitHub用户信息（公开方法）
func (s *SocialAuthService) GetGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	return s.getGitHubUserInfo(accessToken)
}

// exchangeGitHubToken 交换GitHub访问令牌
func (s *SocialAuthService) exchangeGitHubToken(code string) (string, error) {
	clientID := config.Cfg.SocialAuth.GitHubClientID
	clientSecret := config.Cfg.SocialAuth.GitHubClientSecret
	redirectURI := config.Cfg.App.URL + "/api/auth/social/github/callback"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	resp, err := http.PostForm("https://github.com/login/oauth/access_token", data)
	if err != nil {
		return "", errors.New("获取访问令牌失败")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("读取响应失败")
	}

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", errors.New("解析响应失败")
	}

	accessToken := values.Get("access_token")
	if accessToken == "" {
		return "", errors.New("未获取到访问令牌")
	}

	return accessToken, nil
}

// getGitHubUserInfo 获取GitHub用户信息
func (s *SocialAuthService) getGitHubUserInfo(accessToken string) (*GitHubUserInfo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, errors.New("创建请求失败")
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.New("请求GitHub API失败")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("GitHub API返回错误")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("读取响应失败")
	}

	var userInfo GitHubUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, errors.New("解析用户信息失败")
	}

	// 如果邮箱为空，尝试获取邮箱列表
	if userInfo.Email == "" {
		emails, _ := s.getGitHubEmails(accessToken)
		if len(emails) > 0 {
			for _, email := range emails {
				if email.Primary {
					userInfo.Email = email.Email
					break
				}
			}
			if userInfo.Email == "" {
				userInfo.Email = emails[0].Email
			}
		}
	}

	return &userInfo, nil
}

// GitHubEmail GitHub邮箱信息
type GitHubEmail struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func (s *SocialAuthService) getGitHubEmails(accessToken string) ([]GitHubEmail, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("GitHub API返回错误")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var emails []GitHubEmail
	if err := json.Unmarshal(body, &emails); err != nil {
		return nil, err
	}

	return emails, nil
}

// LinkSocialAccount 关联社交媒体账户
func (s *SocialAuthService) LinkSocialAccount(userID uint, provider string, providerID string, providerEmail string, accessToken string) error {
	encryptedToken := utils.EncryptToken(accessToken)

	// 检查是否已关联
	var existing models.SocialAuth
	if err := database.DB.Where("user_id = ? AND provider = ?", userID, provider).First(&existing).Error; err == nil {
		// 更新现有关联
		existing.ProviderID = providerID
		existing.ProviderEmail = providerEmail
		existing.AccessToken = encryptedToken
		if err := database.DB.Save(&existing).Error; err != nil {
			return errors.New("更新关联失败")
		}
		return nil
	}

	// 创建新关联
	socialAuth := models.SocialAuth{
		UserID:        userID,
		Provider:      provider,
		ProviderID:    providerID,
		ProviderEmail: providerEmail,
		AccessToken:   encryptedToken,
	}

	if err := database.DB.Create(&socialAuth).Error; err != nil {
		return errors.New("创建关联失败")
	}

	return nil
}

// FindUserBySocialAccount 通过社交媒体账户查找用户
func (s *SocialAuthService) FindUserBySocialAccount(provider string, providerID string) (*models.User, error) {
	var socialAuth models.SocialAuth
	if err := database.DB.Where("provider = ? AND provider_id = ?", provider, providerID).First(&socialAuth).Error; err != nil {
		return nil, errors.New("未找到关联账户")
	}

	var user models.User
	if err := database.DB.First(&user, socialAuth.UserID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

