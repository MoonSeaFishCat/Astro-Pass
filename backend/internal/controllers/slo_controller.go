package controllers

import (
	"net/http"
	"strconv"

	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type SLOController struct {
	sloService *services.SLOService
}

func NewSLOController() *SLOController {
	return &SLOController{
		sloService: services.NewSLOService(),
	}
}

// InitiateLogoutRequest 发起登出请求
type InitiateLogoutRequest struct {
	SessionID string `json:"session_id" binding:"required"`
}

// InitiateLogout 发起单点登出
// @Summary 发起单点登出
// @Description 发起单点登出，通知所有相关应用
// @Tags SSO
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body InitiateLogoutRequest true "登出请求"
// @Success 200 {object} map[string]interface{}
// @Router /api/sso/logout [post]
func (c *SLOController) InitiateLogout(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	var req InitiateLogoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	logoutRequest, err := c.sloService.InitiateLogout(req.SessionID, "user", userID.(uint))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "登出请求已发起", gin.H{
		"request_id":     logoutRequest.RequestID,
		"status":         logoutRequest.Status,
		"total_clients":  logoutRequest.TotalClients,
	})
}

// GetUserSessions 获取用户的SSO会话
// @Summary 获取用户SSO会话
// @Description 获取当前用户的所有活跃SSO会话
// @Tags SSO
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/sso/sessions [get]
func (c *SLOController) GetUserSessions(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	sessions, err := c.sloService.GetUserActiveSessions(userID.(uint))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// 格式化响应数据
	sessionData := make([]gin.H, 0, len(sessions))
	for _, session := range sessions {
		sessionData = append(sessionData, gin.H{
			"session_id":  session.SessionID,
			"client_id":   session.ClientID,
			"client_name": session.Client.ClientName,
			"created_at":  session.CreatedAt,
			"status":      session.Status,
		})
	}

	utils.Success(ctx, gin.H{
		"sessions": sessionData,
		"total":    len(sessionData),
	})
}

// GetLogoutStatus 获取登出状态
// @Summary 获取登出状态
// @Description 获取指定登出请求的处理状态
// @Tags SSO
// @Security BearerAuth
// @Produce json
// @Param request_id path string true "登出请求ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/sso/logout/{request_id}/status [get]
func (c *SLOController) GetLogoutStatus(ctx *gin.Context) {
	requestID := ctx.Param("request_id")
	if requestID == "" {
		utils.BadRequest(ctx, "请求ID不能为空")
		return
	}

	logoutRequest, notifications, err := c.sloService.GetLogoutStatus(requestID)
	if err != nil {
		utils.NotFound(ctx, err.Error())
		return
	}

	// 格式化通知数据
	notificationData := make([]gin.H, 0, len(notifications))
	for _, notification := range notifications {
		notificationData = append(notificationData, gin.H{
			"client_id":       notification.ClientID,
			"status":          notification.Status,
			"response_code":   notification.ResponseCode,
			"attempt_count":   notification.AttemptCount,
			"last_attempt_at": notification.LastAttemptAt,
		})
	}

	utils.Success(ctx, gin.H{
		"request_id":        logoutRequest.RequestID,
		"status":            logoutRequest.Status,
		"total_clients":     logoutRequest.TotalClients,
		"completed_clients": logoutRequest.CompletedClients,
		"failed_clients":    logoutRequest.FailedClients,
		"notifications":     notificationData,
		"created_at":        logoutRequest.CreatedAt,
		"updated_at":        logoutRequest.UpdatedAt,
	})
}

// AdminRevokeUserSessions 管理员撤销用户会话
// @Summary 管理员撤销用户会话
// @Description 管理员强制撤销指定用户的所有会话
// @Tags SSO
// @Security BearerAuth
// @Produce json
// @Param user_id path string true "用户ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/admin/sso/users/{user_id}/revoke-sessions [post]
func (c *SLOController) AdminRevokeUserSessions(ctx *gin.Context) {
	adminUserID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	userIDStr := ctx.Param("user_id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.BadRequest(ctx, "无效的用户ID")
		return
	}

	err = c.sloService.RevokeUserSessions(uint(userID), "admin", adminUserID.(uint))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "用户会话撤销请求已发起", nil)
}

// HandleLogoutCallback 处理客户端登出回调
// @Summary 处理登出回调
// @Description 接收客户端的登出确认回调
// @Tags SSO
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/sso/logout/callback [post]
func (c *SLOController) HandleLogoutCallback(ctx *gin.Context) {
	var req struct {
		RequestID string `json:"request_id" binding:"required"`
		ClientID  string `json:"client_id" binding:"required"`
		Status    string `json:"status" binding:"required"` // success, failed
		Message   string `json:"message"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	// 这里可以更新登出通知的状态
	// 实际实现中可能需要验证客户端身份
	
	utils.SuccessWithMessage(ctx, "登出回调已处理", nil)
}

// GetOIDCLogout OIDC标准登出端点
// @Summary OIDC登出端点
// @Description OpenID Connect标准的登出端点
// @Tags OIDC
// @Produce json
// @Param id_token_hint query string false "ID Token提示"
// @Param post_logout_redirect_uri query string false "登出后重定向URI"
// @Param state query string false "状态参数"
// @Success 302 {string} string "重定向到登出后URI"
// @Router /api/oidc/logout [get]
func (c *SLOController) GetOIDCLogout(ctx *gin.Context) {
	idTokenHint := ctx.Query("id_token_hint")
	postLogoutRedirectURI := ctx.Query("post_logout_redirect_uri")
	state := ctx.Query("state")

	// 如果提供了ID Token，验证并提取用户信息
	if idTokenHint != "" {
		claims, err := utils.ParseToken(idTokenHint)
		if err == nil && claims.UserID > 0 {
			// 撤销用户的所有会话
			c.sloService.RevokeUserSessions(claims.UserID, "oidc", 0)
		}
	}

	// 构建重定向URL
	if postLogoutRedirectURI != "" {
		redirectURL := postLogoutRedirectURI
		if state != "" {
			redirectURL += "?state=" + state
		}
		ctx.Redirect(http.StatusFound, redirectURL)
		return
	}

	// 如果没有重定向URI，返回成功响应
	utils.SuccessWithMessage(ctx, "登出成功", nil)
}