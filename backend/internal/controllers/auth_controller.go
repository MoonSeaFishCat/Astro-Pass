package controllers

import (
	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController() *AuthController {
	return &AuthController{
		authService: services.NewAuthService(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建新用户账户
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "注册信息"
// @Success 201 {object} map[string]interface{}
// @Router /api/auth/register [post]
func (c *AuthController) Register(ctx *gin.Context) {
	var req RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	user, err := c.authService.Register(req.Username, req.Email, req.Password, req.Nickname)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "注册成功", gin.H{
		"user": gin.H{
			"id":             user.ID,
			"uuid":           user.UUID,
			"username":       user.Username,
			"email":          user.Email,
			"nickname":       user.Nickname,
			"email_verified": user.EmailVerified,
		},
	})
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录获取访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body LoginRequest true "登录信息"
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/login [post]
func (c *AuthController) Login(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	ip := ctx.ClientIP()
	userAgent := ctx.GetHeader("User-Agent")

	user, accessToken, refreshToken, err := c.authService.Login(req.Username, req.Password, ip, userAgent)
	if err != nil {
		utils.Unauthorized(ctx, err.Error())
		return
	}

	// 格式化角色信息
	roles := make([]gin.H, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, gin.H{
			"id":           role.ID,
			"name":         role.Name,
			"display_name": role.DisplayName,
		})
	}

	utils.SuccessWithMessage(ctx, "登录成功", gin.H{
		"user": gin.H{
			"id":             user.ID,
			"uuid":           user.UUID,
			"username":       user.Username,
			"email":          user.Email,
			"nickname":       user.Nickname,
			"roles":          roles,
			"email_verified": user.EmailVerified,
			"mfa_enabled":    user.MFAEnabled,
			"status":         user.Status,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    900, // 15分钟
	})
}

// RefreshToken 刷新访问令牌
// @Summary 刷新访问令牌
// @Description 使用刷新令牌获取新的访问令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body RefreshTokenRequest true "刷新令牌"
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误: "+err.Error())
		return
	}

	accessToken, refreshToken, err := c.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.Unauthorized(ctx, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "令牌刷新成功", gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    900, // 15分钟
	})
}

// GetProfile 获取当前用户信息
// @Summary 获取当前用户信息
// @Description 获取当前登录用户的详细信息
// @Tags 认证
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/auth/profile [get]
func (c *AuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		utils.Unauthorized(ctx, "未认证")
		return
	}

	var user models.User
	if err := database.DB.Preload("Roles").First(&user, userID).Error; err != nil {
		utils.NotFound(ctx, "用户不存在")
		return
	}

	// 格式化角色信息
	roles := make([]gin.H, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, gin.H{
			"id":           role.ID,
			"name":         role.Name,
			"display_name": role.DisplayName,
		})
	}

	utils.Success(ctx, gin.H{
		"id":             user.ID,
		"uuid":           user.UUID,
		"username":       user.Username,
		"email":          user.Email,
		"nickname":       user.Nickname,
		"roles":          roles,
		"email_verified": user.EmailVerified,
		"mfa_enabled":    user.MFAEnabled,
		"status":         user.Status,
	})
}

