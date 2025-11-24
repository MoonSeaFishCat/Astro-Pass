package controllers

import (
	"net/http"
	"strconv"

	"astro-pass/internal/services"
	"astro-pass/internal/utils"
	"github.com/gin-gonic/gin"
)

type PermissionController struct {
	permissionService *services.PermissionService
}

func NewPermissionController() *PermissionController {
	service, _ := services.NewPermissionService()
	return &PermissionController{
		permissionService: service,
	}
}

// AssignRoleRequest 分配角色请求
type AssignRoleRequest struct {
	RoleName string `json:"role_name" binding:"required"`
}

// CreateRoleRequest 创建角色请求
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// CreatePermissionRequest 创建权限请求
type CreatePermissionRequest struct {
	Name        string `json:"name" binding:"required"`
	DisplayName string `json:"display_name"`
	Resource    string `json:"resource" binding:"required"`
	Action      string `json:"action" binding:"required"`
	Description string `json:"description"`
}

// AssignPermissionRequest 分配权限请求
type AssignPermissionRequest struct {
	Resource string `json:"resource" binding:"required"`
	Action   string `json:"action" binding:"required"`
}

// AssignRole 为用户分配角色
func (c *PermissionController) AssignRole(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
		})
		return
	}

	var req AssignRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if err := c.permissionService.AssignRole(userID.(uint), req.RoleName); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "角色分配成功",
	})
}

// GetUserRoles 获取用户角色
func (c *PermissionController) GetUserRoles(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "未认证",
		})
		return
	}

	roles, err := c.permissionService.GetUserRoles(userID.(uint))
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
		"data":    roles,
	})
}

// CreateRole 创建角色（需要管理员权限）
func (c *PermissionController) CreateRole(ctx *gin.Context) {
	var req CreateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	role, err := c.permissionService.CreateRole(req.Name, req.DisplayName, req.Description)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "角色创建成功",
		"data":    role,
	})
}

// CreatePermission 创建权限（需要管理员权限）
func (c *PermissionController) CreatePermission(ctx *gin.Context) {
	var req CreatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	permission, err := c.permissionService.CreatePermission(req.Name, req.DisplayName, req.Resource, req.Action, req.Description)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "权限创建成功",
		"data":    permission,
	})
}

// AssignPermissionToRole 为角色分配权限（需要管理员权限）
func (c *PermissionController) AssignPermissionToRole(ctx *gin.Context) {
	roleName := ctx.Param("role")
	if roleName == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "角色名不能为空",
		})
		return
	}

	var req AssignPermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
		})
		return
	}

	if err := c.permissionService.AssignPermissionToRole(roleName, req.Resource, req.Action); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "权限分配成功",
	})
}

// GetAllRoles 获取所有角色列表（管理员功能）
func (c *PermissionController) GetAllRoles(ctx *gin.Context) {
	roles, err := c.permissionService.GetAllRoles()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, roles)
}

// GetAllPermissions 获取所有权限列表（管理员功能）
func (c *PermissionController) GetAllPermissions(ctx *gin.Context) {
	permissions, err := c.permissionService.GetAllPermissions()
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.Success(ctx, permissions)
}

// UpdateRole 更新角色（管理员功能）
type UpdateRoleRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

func (c *PermissionController) UpdateRole(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req UpdateRoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	role, err := c.permissionService.UpdateRole(uint(roleID), req.DisplayName, req.Description)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "更新成功", role)
}

// DeleteRole 删除角色（管理员功能）
func (c *PermissionController) DeleteRole(ctx *gin.Context) {
	roleIDStr := ctx.Param("id")
	roleID, err := strconv.ParseUint(roleIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的角色ID")
		return
	}

	if err := c.permissionService.DeleteRole(uint(roleID)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

// UpdatePermission 更新权限（管理员功能）
type UpdatePermissionRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

func (c *PermissionController) UpdatePermission(ctx *gin.Context) {
	permissionIDStr := ctx.Param("id")
	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的权限ID")
		return
	}

	var req UpdatePermissionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "请求参数错误")
		return
	}

	permission, err := c.permissionService.UpdatePermission(uint(permissionID), req.DisplayName, req.Description)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "更新成功", permission)
}

// DeletePermission 删除权限（管理员功能）
func (c *PermissionController) DeletePermission(ctx *gin.Context) {
	permissionIDStr := ctx.Param("id")
	permissionID, err := strconv.ParseUint(permissionIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, "无效的权限ID")
		return
	}

	if err := c.permissionService.DeletePermission(uint(permissionID)); err != nil {
		utils.ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessWithMessage(ctx, "删除成功", nil)
}

