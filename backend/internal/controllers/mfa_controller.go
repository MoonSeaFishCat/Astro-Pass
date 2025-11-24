package controllers

import (
	"net/http"

	"astro-pass/internal/services"
	"github.com/gin-gonic/gin"
)

type MFAController struct {
	mfaService *services.MFAService
}

func NewMFAController() *MFAController {
	return &MFAController{
		mfaService: services.NewMFAService(),
	}
}

// GenerateTOTPRequest 生成TOTP请求
type GenerateTOTPRequest struct {
	// 无需参数，从上下文获取用户ID
}

// EnableMFARequest 启用MFA请求
type EnableMFARequest struct {
	Code string `json:"code" binding:"required"`
}

// DisableMFARequest 禁用MFA请求
type DisableMFARequest struct {
	Code string `json:"code" binding:"required"`
}

// VerifyMFARequest 验证MFA请求
type VerifyMFARequest struct {
	Code string `json:"code" binding:"required"`
}

// GenerateTOTP 生成TOTP密钥和二维码
func (c *MFAController) GenerateTOTP(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
		})
		return
	}

	email, _ := ctx.Get("email")
	secret, qrCodeURL, err := c.mfaService.GenerateTOTPSecret(userID.(uint), email.(string))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "生成成功",
		"data": gin.H{
			"secret":     secret,
			"qr_code_url": qrCodeURL,
		},
	})
}

// EnableMFA 启用MFA
func (c *MFAController) EnableMFA(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
		})
		return
	}

	var req EnableMFARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	if err := c.mfaService.EnableMFA(userID.(uint), req.Code); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	// 获取恢复码
	recoveryCodes, _ := c.mfaService.GetRecoveryCodes(userID.(uint))

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "MFA启用成功",
		"data": gin.H{
			"recovery_codes": recoveryCodes,
		},
	})
}

// DisableMFA 禁用MFA
func (c *MFAController) DisableMFA(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
		})
		return
	}

	var req DisableMFARequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	if err := c.mfaService.DisableMFA(userID.(uint), req.Code); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "MFA禁用成功",
	})
}

// GetRecoveryCodes 获取恢复码
func (c *MFAController) GetRecoveryCodes(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
		})
		return
	}

	recoveryCodes, err := c.mfaService.GetRecoveryCodes(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "获取成功",
		"data": gin.H{
			"recovery_codes": recoveryCodes,
		},
	})
}


