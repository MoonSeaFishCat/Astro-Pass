package controllers

import (
	"net/http"

	"astro-pass/internal/services"
	"github.com/gin-gonic/gin"
)

type OAuth2Controller struct {
	oauth2Service *services.OAuth2Service
}

func NewOAuth2Controller() *OAuth2Controller {
	return &OAuth2Controller{
		oauth2Service: services.NewOAuth2Service(),
	}
}

// AuthorizeRequest 授权请求
type AuthorizeRequest struct {
	ResponseType        string `form:"response_type" binding:"required"`
	ClientID            string `form:"client_id" binding:"required"`
	RedirectURI         string `form:"redirect_uri" binding:"required"`
	Scope               string `form:"scope"`
	State               string `form:"state"`
	CodeChallenge       string `form:"code_challenge"`
	CodeChallengeMethod string `form:"code_challenge_method"`
}

// TokenRequest 令牌请求
type TokenRequest struct {
	GrantType    string `form:"grant_type" binding:"required"`
	Code         string `form:"code"`
	RedirectURI  string `form:"redirect_uri"`
	ClientID     string `form:"client_id" binding:"required"`
	ClientSecret string `form:"client_secret" binding:"required"`
	CodeVerifier string `form:"code_verifier"`
}

// Authorize OAuth2授权端点
func (c *OAuth2Controller) Authorize(ctx *gin.Context) {
	var req AuthorizeRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证response_type
	if req.ResponseType != "code" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的response_type",
		})
		return
	}

	// 检查用户是否已登录（这里简化处理，实际应该检查session）
	userID, exists := ctx.Get("user_id")
	if !exists {
		// 重定向到登录页面
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "请先登录",
			"login_url": "/api/auth/login",
		})
		return
	}

	// 生成授权码
	code, err := c.oauth2Service.GenerateAuthorizationCode(
		req.ClientID,
		userID.(uint),
		req.RedirectURI,
		req.Scope,
		req.CodeChallenge,
		req.CodeChallengeMethod,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	// 构建重定向URI
	redirectURL := req.RedirectURI + "?code=" + code
	if req.State != "" {
		redirectURL += "&state=" + req.State
	}

	ctx.Redirect(http.StatusFound, redirectURL)
}

// Token OAuth2令牌端点
func (c *OAuth2Controller) Token(ctx *gin.Context) {
	var req TokenRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	// 验证grant_type
	if req.GrantType != "authorization_code" && req.GrantType != "refresh_token" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "不支持的grant_type",
		})
		return
	}

	if req.GrantType == "authorization_code" {
		// 交换授权码
		accessToken, refreshToken, err := c.oauth2Service.ExchangeAuthorizationCode(
			req.Code,
			req.ClientID,
			req.ClientSecret,
			req.RedirectURI,
			req.CodeVerifier,
		)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"code":    400,
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken.Token,
			"token_type":    "Bearer",
			"expires_in":    900, // 15分钟
			"refresh_token": refreshToken,
			"scope":         accessToken.Scope,
		})
	} else {
		// 刷新令牌逻辑（可以复用auth_service的RefreshToken方法）
		ctx.JSON(http.StatusOK, gin.H{
			"message": "刷新令牌功能（待完善）",
		})
	}
}

// UserInfo OIDC用户信息端点
func (c *OAuth2Controller) UserInfo(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未提供访问令牌",
		})
		return
	}

	// 提取Bearer token
	tokenString := authHeader[7:] // 跳过"Bearer "
	userInfo, err := c.oauth2Service.GetUserInfo(tokenString)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, userInfo)
}

// JWKS JWKS端点（简化实现）
func (c *OAuth2Controller) JWKS(ctx *gin.Context) {
	// 这里应该返回实际的公钥，简化处理
	ctx.JSON(http.StatusOK, gin.H{
		"keys": []gin.H{},
	})
}


