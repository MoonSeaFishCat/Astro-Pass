package controllers

import (
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
)

type WebAuthnController struct {
	webauthnService *services.WebAuthnService
}

func NewWebAuthnController() *WebAuthnController {
	webauthnService, err := services.NewWebAuthnService()
	if err != nil {
		// 如果初始化失败，返回nil（可以在路由中检查）
		return nil
	}
	return &WebAuthnController{
		webauthnService: webauthnService,
	}
}

// BeginRegistrationRequest 开始注册请求
type BeginRegistrationRequest struct {
	UserID uint `json:"user_id" binding:"required"`
}

// BeginRegistration 开始WebAuthn注册
func (c *WebAuthnController) BeginRegistration(ctx *gin.Context) {
	if c.webauthnService == nil {
		utils.InternalError(ctx, "WebAuthn服务未初始化")
		return
	}

	// 从JWT中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未授权")
		return
	}

	sessionData, err := c.webauthnService.BeginRegistration(userID.(uint))
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// 生成会话令牌
	sessionToken, err := services.GenerateSessionToken()
	if err != nil {
		utils.InternalError(ctx, "生成会话令牌失败")
		return
	}

	// 存储会话数据
	if err := services.StoreSessionData(sessionToken, sessionData); err != nil {
		utils.InternalError(ctx, "存储会话数据失败")
		return
	}

	// 将SessionData转换为PublicKeyCredentialCreationOptions格式返回给前端
	// SessionData已经包含了所有需要的信息，我们需要将其转换为前端可用的格式
	utils.Success(ctx, gin.H{
		"session_token": sessionToken,
		"options":       sessionData, // 直接返回SessionData，它包含了所有需要的字段
	})
}

// FinishRegistrationRequest 完成注册请求
type FinishRegistrationRequest struct {
	SessionToken string                          `json:"session_token" binding:"required"`
	Response     protocol.CredentialCreationResponse `json:"response" binding:"required"`
}

// FinishRegistration 完成WebAuthn注册
func (c *WebAuthnController) FinishRegistration(ctx *gin.Context) {
	if c.webauthnService == nil {
		utils.InternalError(ctx, "WebAuthn服务未初始化")
		return
	}

	var req FinishRegistrationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 获取会话数据
	sessionData, err := services.GetSessionData(req.SessionToken)
	if err != nil {
		utils.BadRequest(ctx, "会话无效或已过期")
		return
	}

	// 从JWT中获取用户ID
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未授权")
		return
	}

	// 完成注册 - 传递 http.Request
	credential, err := c.webauthnService.FinishRegistration(userID.(uint), sessionData, ctx.Request)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// 删除会话数据
	services.DeleteSessionData(req.SessionToken)

	utils.SuccessWithMessage(ctx, "WebAuthn凭证注册成功", gin.H{
		"credential": gin.H{
			"id":          credential.ID,
			"device_name": credential.DeviceName,
			"device_type": credential.DeviceType,
			"created_at":  credential.CreatedAt,
		},
	})
}

// BeginLoginRequest 开始登录请求
type BeginLoginRequest struct {
	Username string `json:"username" binding:"required"`
}

// BeginLogin 开始WebAuthn登录
func (c *WebAuthnController) BeginLogin(ctx *gin.Context) {
	if c.webauthnService == nil {
		utils.InternalError(ctx, "WebAuthn服务未初始化")
		return
	}

	var req BeginLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	sessionData, credentials, err := c.webauthnService.BeginLogin(req.Username)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// 生成会话令牌
	sessionToken, err := services.GenerateSessionToken()
	if err != nil {
		utils.InternalError(ctx, "生成会话令牌失败")
		return
	}

	// 存储会话数据
	if err := services.StoreSessionData(sessionToken, sessionData); err != nil {
		utils.InternalError(ctx, "存储会话数据失败")
		return
	}

	// 转换凭证ID为base64格式
	credentialIDs := make([]string, len(credentials))
	for i, cred := range credentials {
		credentialIDs[i] = cred.CredentialID
	}

	// 将SessionData转换为PublicKeyCredentialRequestOptions格式返回给前端
	utils.Success(ctx, gin.H{
		"session_token": sessionToken,
		"options":       sessionData, // 直接返回SessionData，它包含了所有需要的字段
	})
}

// FinishLoginRequest 完成登录请求
type FinishLoginRequest struct {
	SessionToken string                       `json:"session_token" binding:"required"`
	Username     string                       `json:"username" binding:"required"`
	Response     protocol.CredentialAssertionResponse `json:"response" binding:"required"`
}

// FinishLogin 完成WebAuthn登录
func (c *WebAuthnController) FinishLogin(ctx *gin.Context) {
	if c.webauthnService == nil {
		utils.InternalError(ctx, "WebAuthn服务未初始化")
		return
	}

	var req FinishLoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 获取会话数据
	sessionData, err := services.GetSessionData(req.SessionToken)
	if err != nil {
		utils.BadRequest(ctx, "会话无效或已过期")
		return
	}

	// 完成登录 - 传递 http.Request
	user, err := c.webauthnService.FinishLogin(req.Username, sessionData, ctx.Request)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	// 删除会话数据
	services.DeleteSessionData(req.SessionToken)

	// 生成JWT令牌
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		utils.InternalError(ctx, "生成访问令牌失败")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.InternalError(ctx, "生成刷新令牌失败")
		return
	}

	utils.SuccessWithMessage(ctx, "WebAuthn登录成功", gin.H{
		"user": gin.H{
			"id":       user.ID,
			"uuid":     user.UUID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GetCredentials 获取用户的WebAuthn凭证列表
func (c *WebAuthnController) GetCredentials(ctx *gin.Context) {
	if c.webauthnService == nil {
		utils.InternalError(ctx, "WebAuthn服务未初始化")
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未授权")
		return
	}

	credentials, err := c.webauthnService.GetUserCredentials(userID.(uint))
	if err != nil {
		utils.InternalError(ctx, "获取凭证列表失败")
		return
	}

	utils.Success(ctx, gin.H{
		"credentials": credentials,
	})
}

// DeleteCredentialRequest 删除凭证请求
type DeleteCredentialRequest struct {
	CredentialID uint `json:"credential_id" binding:"required"`
}

// DeleteCredential 删除WebAuthn凭证
func (c *WebAuthnController) DeleteCredential(ctx *gin.Context) {
	if c.webauthnService == nil {
		utils.InternalError(ctx, "WebAuthn服务未初始化")
		return
	}

	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未授权")
		return
	}

	var req DeleteCredentialRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	if err := c.webauthnService.DeleteCredential(userID.(uint), req.CredentialID); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "凭证删除成功", nil)
}

