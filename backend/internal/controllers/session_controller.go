package controllers

import (
	"fmt"

	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type SessionController struct {
	sessionService *services.SessionService
}

func NewSessionController() *SessionController {
	return &SessionController{
		sessionService: services.NewSessionService(),
	}
}

// GetSessions 获取用户的所有活跃会话
func (c *SessionController) GetSessions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	sessions, err := c.sessionService.GetUserSessions(userID.(uint))
	if err != nil {
		utils.InternalError(ctx, err.Error())
		return
	}

	utils.Success(ctx, sessions)
}

// RevokeSession 撤销指定会话
func (c *SessionController) RevokeSession(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	sessionID := ctx.Param("id")
	var sessionIDUint uint
	if _, err := fmt.Sscanf(sessionID, "%d", &sessionIDUint); err != nil {
		utils.BadRequest(ctx, "无效的会话ID")
		return
	}

	if err := c.sessionService.RevokeSession(sessionIDUint, userID.(uint)); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "会话已撤销", nil)
}

// RevokeAllSessions 撤销所有会话（除了当前会话）
func (c *SessionController) RevokeAllSessions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	// 从请求头获取当前token
	authHeader := ctx.GetHeader("Authorization")
	currentToken := ""
	if len(authHeader) > 7 {
		currentToken = authHeader[7:] // 跳过 "Bearer "
	}

	if err := c.sessionService.RevokeAllSessions(userID.(uint), currentToken); err != nil {
		utils.InternalError(ctx, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "所有其他会话已撤销", nil)
}

