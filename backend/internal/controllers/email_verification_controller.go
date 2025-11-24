package controllers

import (
	"net/http"

	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type EmailVerificationController struct {
	emailVerificationService *services.EmailVerificationService
}

func NewEmailVerificationController() *EmailVerificationController {
	return &EmailVerificationController{
		emailVerificationService: services.NewEmailVerificationService(),
	}
}

// SendVerificationEmailRequest 发送验证邮件请求
type SendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// SendVerificationEmail 发送邮箱验证邮件
func (c *EmailVerificationController) SendVerificationEmail(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	var req SendVerificationEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误")
		return
	}

	if err := c.emailVerificationService.SendVerificationEmail(userID.(uint), req.Email); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "验证邮件已发送，请查收", nil)
}

// VerifyEmailRequest 验证邮箱请求
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
}

// VerifyEmail 验证邮箱
func (c *EmailVerificationController) VerifyEmail(ctx *gin.Context) {
	var req VerifyEmailRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误")
		return
	}

	if err := c.emailVerificationService.VerifyEmail(req.Token); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "邮箱验证成功", nil)
}


