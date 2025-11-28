package controllers

import (
	"astro-pass/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ConsentController struct {
	consentService *services.ConsentService
	oauth2Service  *services.OAuth2Service
}

func NewConsentController() *ConsentController {
	return &ConsentController{
		consentService: services.NewConsentService(),
		oauth2Service:  services.NewOAuth2Service(),
	}
}

// GetConsentInfo 获取授权同意信息
func (cc *ConsentController) GetConsentInfo(c *gin.Context) {
	clientID := c.Query("client_id")
	scope := c.Query("scope")

	if clientID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "缺少client_id参数",
		})
		return
	}

	// 获取客户端信息
	var client struct {
		ClientName string
		ClientURI  string
		LogoURI    string
	}

	// 这里应该从数据库查询客户端信息
	// 简化处理，返回模拟数据
	client.ClientName = "示例应用"
	client.ClientURI = "https://example.com"
	client.LogoURI = "https://example.com/logo.png"

	// 解析scope
	scopes := parseScopeDescriptions(scope)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"client_name": client.ClientName,
			"client_uri":  client.ClientURI,
			"logo_uri":    client.LogoURI,
			"scopes":      scopes,
		},
	})
}

// ApproveConsent 批准授权
func (cc *ConsentController) ApproveConsent(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Scope    string `json:"scope" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	// 保存授权
	err := cc.consentService.SaveConsent(userID, req.ClientID, req.Scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "保存授权失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "授权成功",
	})
}

// DenyConsent 拒绝授权
func (cc *ConsentController) DenyConsent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "已拒绝授权",
	})
}

// GetUserConsents 获取用户的所有授权
func (cc *ConsentController) GetUserConsents(c *gin.Context) {
	userID := c.GetUint("user_id")

	consents, err := cc.consentService.GetUserConsents(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "获取授权列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": consents,
	})
}

// RevokeConsent 撤销授权
func (cc *ConsentController) RevokeConsent(c *gin.Context) {
	userID := c.GetUint("user_id")
	clientID := c.Param("client_id")

	err := cc.consentService.RevokeConsent(userID, clientID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "撤销授权成功",
	})
}

// parseScopeDescriptions 解析scope并返回描述
func parseScopeDescriptions(scopeString string) []map[string]string {
	scopeDescriptions := map[string]string{
		"openid":  "访问您的基本身份信息",
		"profile": "访问您的个人资料（昵称、头像等）",
		"email":   "访问您的邮箱地址",
		"phone":   "访问您的手机号码",
		"address": "访问您的地址信息",
	}

	var result []map[string]string
	if scopeString == "" {
		return result
	}

	scopes := splitBySpace(scopeString)
	for _, scope := range scopes {
		if scope != "" {
			description := scopeDescriptions[scope]
			if description == "" {
				description = "访问 " + scope + " 权限"
			}
			result = append(result, map[string]string{
				"scope":       scope,
				"description": description,
			})
		}
	}

	return result
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
