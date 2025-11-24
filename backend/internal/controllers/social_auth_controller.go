package controllers

import (
	"fmt"
	"net/http"
	"time"

	"astro-pass/internal/database"
	"astro-pass/internal/models"
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type SocialAuthController struct {
	socialAuthService *services.SocialAuthService
}

func NewSocialAuthController() *SocialAuthController {
	return &SocialAuthController{
		socialAuthService: services.NewSocialAuthService(),
	}
}

// GetGitHubAuthURL 获取GitHub授权URL
func (c *SocialAuthController) GetGitHubAuthURL(ctx *gin.Context) {
	state := ctx.Query("state")
	if state == "" {
		state = utils.GenerateUUID()
	}

	authURL := c.socialAuthService.GetGitHubAuthURL(state)
	utils.Success(ctx, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// GitHubCallbackRequest GitHub回调请求
type GitHubCallbackRequest struct {
	Code  string `json:"code" binding:"required"`
	State string `json:"state"`
}

// HandleGitHubCallback 处理GitHub回调
func (c *SocialAuthController) HandleGitHubCallback(ctx *gin.Context) {
	var req GitHubCallbackRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(ctx, "请求参数错误")
		return
	}

	// 交换访问令牌
	accessToken, err := c.socialAuthService.ExchangeGitHubToken(req.Code)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "获取访问令牌失败")
		return
	}

	// 获取GitHub用户信息
	githubUser, err := c.socialAuthService.GetGitHubUserInfo(accessToken)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 查找或创建用户
	var user *models.User
	existingUser, err := c.socialAuthService.FindUserBySocialAccount("github", fmt.Sprintf("%d", githubUser.ID))
	if err != nil {
		// 用户不存在，检查邮箱是否已注册
		var existingUserByEmail models.User
		if err := database.DB.Where("email = ?", githubUser.Email).First(&existingUserByEmail).Error; err == nil {
			// 邮箱已存在，关联GitHub账户
			user = &existingUserByEmail
			if err := c.socialAuthService.LinkSocialAccount(
				user.ID,
				"github",
				fmt.Sprintf("%d", githubUser.ID),
				githubUser.Email,
				accessToken,
			); err != nil {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "关联账户失败")
				return
			}
		} else {
			// 创建新用户
			user = &models.User{
				UUID:          utils.GenerateUUID(),
				Username:      githubUser.Login,
				Email:         githubUser.Email,
				Nickname:      githubUser.Name,
				Status:        "active",
				EmailVerified: true, // GitHub已验证邮箱
			}
			if err := database.DB.Create(user).Error; err != nil {
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "创建用户失败")
				return
			}

			// 关联GitHub账户
			if err := c.socialAuthService.LinkSocialAccount(
				user.ID,
				"github",
				fmt.Sprintf("%d", githubUser.ID),
				githubUser.Email,
				accessToken,
			); err != nil {
				utils.ErrorResponse(ctx, http.StatusBadRequest, "关联账户失败")
				return
			}
		}
	} else {
		user = existingUser
	}

	// 生成JWT令牌
	jwtAccessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "生成访问令牌失败")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "生成刷新令牌失败")
		return
	}

	// 保存刷新令牌
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 7),
	}
	database.DB.Create(refreshTokenModel)

	// 格式化角色信息
	roles := make([]gin.H, 0)
	if user.Roles != nil {
		for _, role := range user.Roles {
			roles = append(roles, gin.H{
				"id":           role.ID,
				"name":         role.Name,
				"display_name": role.DisplayName,
			})
		}
	}

	utils.SuccessWithMessage(ctx, "登录成功", gin.H{
		"user": gin.H{
			"id":       user.ID,
			"uuid":     user.UUID,
			"username": user.Username,
			"email":    user.Email,
			"nickname": user.Nickname,
			"roles":    roles,
		},
		"access_token":  jwtAccessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    900,
	})
}

