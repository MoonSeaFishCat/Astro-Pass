package controllers

import (
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type OAuth2ClientController struct {
	oauth2Service *services.OAuth2Service
}

func NewOAuth2ClientController() *OAuth2ClientController {
	return &OAuth2ClientController{
		oauth2Service: services.NewOAuth2Service(),
	}
}

// CreateClientRequest 创建客户端请求
type CreateClientRequest struct {
	ClientName   string   `json:"client_name" binding:"required"`
	ClientURI    string   `json:"client_uri"`
	LogoURI      string   `json:"logo_uri"`
	RedirectURIs []string `json:"redirect_uris" binding:"required"`
}

// CreateClient 创建OAuth2客户端
func (c *OAuth2ClientController) CreateClient(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	var req CreateClientRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	client, err := c.oauth2Service.CreateClient(
		userID.(uint),
		req.ClientName,
		req.ClientURI,
		req.LogoURI,
		req.RedirectURIs,
	)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "客户端创建成功", gin.H{
		"client_id":    client.ClientID,
		"client_name":  client.ClientName,
		"client_uri":   client.ClientURI,
		"logo_uri":     client.LogoURI,
		"status":       client.Status,
		"client_secret": client.ClientSecret, // 只在创建时返回一次
	})
}

// GetUserClients 获取用户的OAuth2客户端列表
func (c *OAuth2ClientController) GetUserClients(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	clients, err := c.oauth2Service.GetUserClients(userID.(uint))
	if err != nil {
		utils.InternalError(ctx, err.Error())
		return
	}

	// 隐藏敏感信息
	clientList := make([]gin.H, 0, len(clients))
	for _, client := range clients {
		clientList = append(clientList, gin.H{
			"id":          client.ID,
			"client_id":   client.ClientID,
			"client_name": client.ClientName,
			"client_uri":  client.ClientURI,
			"logo_uri":    client.LogoURI,
			"status":      client.Status,
			"created_at":  client.CreatedAt,
		})
	}

	utils.Success(ctx, clientList)
}

// RevokeClient 撤销客户端
func (c *OAuth2ClientController) RevokeClient(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	clientID := ctx.Param("id")
	if err := c.oauth2Service.RevokeClient(clientID, userID.(uint)); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "客户端已撤销", nil)
}

