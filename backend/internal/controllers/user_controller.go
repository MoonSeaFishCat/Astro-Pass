package controllers

import (
	"net/http"
	"strconv"
	"astro-pass/internal/config"
	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController() *UserController {
	return &UserController{
		userService: services.NewUserService(),
	}
}

// UpdateProfileRequest 更新资料请求
type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ForgotPasswordRequest 忘记密码请求
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// UpdateProfile 更新用户资料
func (c *UserController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
		})
		return
	}

	var req UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	user, err := c.userService.UpdateProfile(userID.(uint), req.Nickname)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
		"data": gin.H{
			"user": gin.H{
				"id":       user.ID,
				"uuid":     user.UUID,
				"username": user.Username,
				"email":    user.Email,
				"nickname": user.Nickname,
			},
		},
	})
}

// ChangePassword 修改密码
func (c *UserController) ChangePassword(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
		})
		return
	}

	var req ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	if err := c.userService.ChangePassword(userID.(uint), req.OldPassword, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码修改成功",
	})
}

// ForgotPassword 忘记密码
func (c *UserController) ForgotPassword(ctx *gin.Context) {
	var req ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	userService := services.NewUserService()
	user, err := userService.GetUserByEmail(req.Email)
	if err != nil {
		// 为了安全，即使用户不存在也返回成功
		ctx.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "如果该邮箱存在，重置链接已发送",
		})
		return
	}

	// 生成重置令牌
	token, err := userService.GeneratePasswordResetToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成重置令牌失败",
		})
		return
	}

	// 发送密码重置邮件
	emailService := services.NewEmailService()
	if err := emailService.SendPasswordResetEmail(user.Email, token); err != nil {
		// 如果邮件发送失败，记录日志但不暴露给用户（安全考虑）
		utils.Warn("发送密码重置邮件失败: %v", err)
		// 开发环境可以返回token，生产环境不应该
		if config.Cfg.Server.Mode == "debug" {
			utils.SuccessWithMessage(ctx, "重置链接已生成（开发模式）", gin.H{
				"reset_token": token,
			})
			return
		}
	}

	utils.SuccessWithMessage(ctx, "如果该邮箱存在，重置链接已发送到您的邮箱", nil)
}

// ResetPassword 重置密码
func (c *UserController) ResetPassword(ctx *gin.Context) {
	var req ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	if err := c.userService.ResetPassword(req.Token, req.NewPassword); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "密码重置成功",
	})
}

// GetAllUsers 获取所有用户列表（管理员功能）
func (c *UserController) GetAllUsers(ctx *gin.Context) {
	page := 1
	pageSize := 20
	search := ""

	if p := ctx.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}
	if ps := ctx.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}
	if s := ctx.Query("search"); s != "" {
		search = s
	}

	users, total, err := c.userService.GetAllUsers(page, pageSize, search)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 格式化用户数据（隐藏敏感信息）
	userList := make([]gin.H, 0, len(users))
	for _, user := range users {
		roles := make([]gin.H, 0, len(user.Roles))
		for _, role := range user.Roles {
			roles = append(roles, gin.H{
				"id":          role.ID,
				"name":        role.Name,
				"display_name": role.DisplayName,
			})
		}

		userList = append(userList, gin.H{
			"id":             user.ID,
			"uuid":           user.UUID,
			"username":       user.Username,
			"email":          user.Email,
			"nickname":       user.Nickname,
			"status":         user.Status,
			"email_verified": user.EmailVerified,
			"mfa_enabled":    user.MFAEnabled,
			"last_login_at":  user.LastLoginAt,
			"created_at":     user.CreatedAt,
			"roles":          roles,
		})
	}

	utils.Success(ctx, gin.H{
		"users": userList,
		"pagination": gin.H{
			"page":       page,
			"page_size":  pageSize,
			"total":      total,
			"total_pages": (int(total) + pageSize - 1) / pageSize,
		},
	})
}

// GetUser 获取单个用户信息（管理员功能）
func (c *UserController) GetUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的用户ID")
		return
	}

	user, err := c.userService.GetUserByID(uint(userID))
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, err.Error())
		return
	}

	// 格式化用户数据
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
		"status":         user.Status,
		"email_verified": user.EmailVerified,
		"mfa_enabled":    user.MFAEnabled,
		"last_login_at":  user.LastLoginAt,
		"last_login_ip":  user.LastLoginIP,
		"created_at":     user.CreatedAt,
		"updated_at":     user.UpdatedAt,
		"roles":          roles,
	})
}

// UpdateUser 更新用户信息（管理员功能）
type UpdateUserRequest struct {
	Nickname      *string `json:"nickname"`
	Email         *string `json:"email"`
	Status        *string `json:"status"`
	EmailVerified *bool   `json:"email_verified"`
}

func (c *UserController) UpdateUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	updates := make(map[string]interface{})
	if req.Nickname != nil {
		updates["nickname"] = *req.Nickname
	}
	if req.Email != nil {
		updates["email"] = *req.Email
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.EmailVerified != nil {
		updates["email_verified"] = *req.EmailVerified
	}

	user, err := c.userService.UpdateUser(uint(userID), updates)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "更新成功", gin.H{
		"user": user,
	})
}

// DeleteUser 删除用户（管理员功能）
func (c *UserController) DeleteUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的用户ID")
		return
	}

	if err := c.userService.DeleteUser(uint(userID)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// AssignRoleToUser 为用户分配角色（管理员功能）
type AssignRoleToUserRequest struct {
	RoleName string `json:"role_name" binding:"required"`
}

func (c *UserController) AssignRoleToUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req AssignRoleToUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := c.userService.AssignRoleToUser(uint(userID), req.RoleName); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "角色分配成功", nil)
}

// RemoveRoleFromUser 移除用户角色（管理员功能）
type RemoveRoleFromUserRequest struct {
	RoleName string `json:"role_name" binding:"required"`
}

func (c *UserController) RemoveRoleFromUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的用户ID")
		return
	}

	var req RemoveRoleFromUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	if err := c.userService.RemoveRoleFromUser(uint(userID), req.RoleName); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "角色移除成功", nil)
}

// GetUserStats 获取用户统计信息（管理员功能）
func (c *UserController) GetUserStats(ctx *gin.Context) {
	stats, err := c.userService.GetUserStats()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	utils.Success(ctx, stats)
}

